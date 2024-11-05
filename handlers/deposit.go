/**
 * 存款处理
 */

package handlers

import (
	"net/http"
	"sync"

	"wallet-service/models"
	"wallet-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 存款请求结构体
type DepositRequest struct {
	UserID string  `json:"user_id" validate:"required"`     // 用户ID
	Amount float64 `json:"amount" validate:"required,gt=0"` // 金额
}

// 存款
func Deposit(c *gin.Context) {
	var req DepositRequest
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

	// 更新钱包余额
	newBalance := wallet.Balance + req.Amount // 新余额 = 余额 + 存款金额
	if err := models.UpdateWalletBalance(tx, req.UserID, newBalance); err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	// 创建交易记录
	if err := models.CreateTransaction(tx, req.UserID, req.UserID, req.Amount); err != nil {
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
