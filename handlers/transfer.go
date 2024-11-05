/**
 * 转账处理
 */

package handlers

import (
	"math"
	"net/http"
	"sync"

	"wallet-service/models"
	"wallet-service/utils"

	"github.com/gin-gonic/gin"
)

// 转账请求结构体
type TransferRequest struct {
	SenderID   string  `json:"sender_id" validate:"required"`   // 发送方用户ID
	ReceiverID string  `json:"receiver_id" validate:"required"` // 接收方用户ID
	Amount     float64 `json:"amount" validate:"required,gt=0"` // 金额
}

// 转账
func Transfer(c *gin.Context) {
	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}

	// 验证参数
	if err := validate.Struct(req); err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	// 获取或创建用户锁
	senderMutex, _ := userLocks.LoadOrStore(req.SenderID, &sync.Mutex{})
	receiverMutex, _ := userLocks.LoadOrStore(req.ReceiverID, &sync.Mutex{})

	// 确保发送方和接收方的锁都锁定
	senderMutex.(*sync.Mutex).Lock()
	defer senderMutex.(*sync.Mutex).Unlock()

	receiverMutex.(*sync.Mutex).Lock()
	defer receiverMutex.(*sync.Mutex).Unlock()

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				HandleError(c, rbErr, http.StatusInternalServerError)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				HandleError(c, cmErr, http.StatusInternalServerError)
			}
		}
	}()

	// 获取发送方和接收方的钱包信息
	senderWallet, err := models.GetWalletByUserID(tx, req.SenderID)
	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	receiverWallet, err := models.GetWalletByUserID(tx, req.ReceiverID)
	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	// 检查发送方余额是否足够
	if senderWallet.Balance < req.Amount {
		HandleError(c, models.ErrInsufficientFunds, http.StatusInsufficientStorage)
		return
	}

	// 更新发送方和接收方的余额并四舍五入保留两位小数
	newSenderBalance := math.Round((senderWallet.Balance-req.Amount)*100) / 100
	newReceiverBalance := math.Round((receiverWallet.Balance+req.Amount)*100) / 100

	if err := models.UpdateWalletBalance(tx, req.SenderID, newSenderBalance); err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	if err := models.UpdateWalletBalance(tx, req.ReceiverID, newReceiverBalance); err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	// 使用通道和协程异步创建交易记录
	errChan := make(chan error, 2)
	go func() {
		err := models.CreateTransaction(tx, req.SenderID, req.SenderID, -req.Amount)
		errChan <- err
	}()
	go func() {
		err := models.CreateTransaction(tx, req.SenderID, req.ReceiverID, req.Amount)
		errChan <- err
	}()

	// 等待创建交易记录的结果
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			HandleError(c, err, http.StatusInternalServerError)
			return
		}
	}

	// 构建响应数据
	response := gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"sender_id":            req.SenderID,
			"receiver_id":          req.ReceiverID,
			"new_sender_balance":   newSenderBalance,
			"new_receiver_balance": newReceiverBalance,
		},
		"message": "successful",
	}

	c.JSON(http.StatusOK, response)
}
