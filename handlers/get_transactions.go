/**
 * 获取指定用户的交易历史
 */

package handlers

import (
	"net/http"

	"wallet-service/models"
	"wallet-service/utils"

	"github.com/gin-gonic/gin"
)

// 获取交易历史请求结构体
type GetTransactionsRequest struct {
	UserID   string `json:"user_id" validate:"required"`                  // 用户ID
	Page     int    `json:"page" validate:"omitempty,min=1"`              // 当前分页
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"` // 分页大小
}

// 获取用户交易历史
func GetTransactions(c *gin.Context) {
	var req GetTransactionsRequest
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

	// 获取用户交易历史
	transactions, err := models.GetTransactionsByUserID(db, req.UserID, req.Page, req.PageSize)
	if err != nil {
		if err == models.ErrWalletNotFound {
			HandleError(c, err, http.StatusNotFound)
		} else {
			HandleError(c, err, http.StatusInternalServerError)
		}
		return
	}

	// 构建响应数据
	response := gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"user_id":      req.UserID,
			"transactions": transactions,
		},
		"message": "successful",
	}

	c.JSON(http.StatusOK, response)
}
