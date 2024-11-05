/**
 * 取款处理
 */

package handlers

import (
	"math"
	"net/http"
	"sync"

	"wallet-service/models"
	"wallet-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 取款请求结构体
type WithdrawRequest struct {
	UserID string  `json:"user_id" validate:"required"`     // 用户ID
	Amount float64 `json:"amount" validate:"required,gt=0"` // 金额
}

// 取款
func Withdraw(c *gin.Context) {
	var req WithdrawRequest
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
	userMutex, _ := userLocks.LoadOrStore(req.UserID, &sync.Mutex{})
	userMutex.(*sync.Mutex).Lock()
	defer userMutex.(*sync.Mutex).Unlock()

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logrus.Errorf("Rollback error: %v", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				HandleError(c, cmErr, http.StatusInternalServerError)
			}
		}
	}()

	wallet, err := models.GetWalletByUserID(tx, req.UserID) // 获取用户钱包余额等信息
	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	// 检查用户余额是否足够
	if wallet.Balance < req.Amount {
		HandleError(c, models.ErrInsufficientFunds, http.StatusInsufficientStorage)
		return
	}

	// 更新钱包余额并四舍五入保留两位小数
	newBalance := math.Round((wallet.Balance-req.Amount)*100) / 100
	if err := models.UpdateWalletBalance(tx, req.UserID, newBalance); err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	// 使用通道和协程异步创建交易记录
	errChan := make(chan error, 1)
	go func() {
		err := models.CreateTransaction(tx, req.UserID, req.UserID, -req.Amount)
		errChan <- err
	}()

	// 等待创建交易记录的结果
	if err := <-errChan; err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	// 构建响应数据
	response := gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"user_id":     req.UserID,
			"new_balance": newBalance,
		},
		"message": "successful",
	}

	c.JSON(http.StatusOK, response)
}
