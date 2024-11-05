/**
 * 获取指定用户的余额
 */

package handlers

import (
	"net/http"

	"wallet-service/models"
	"wallet-service/utils"

	"github.com/gin-gonic/gin"
)

// 获取用户余额请求结构体
type GetBalanceRequest struct {
	UserID string `json:"user_id" validate:"required"` // 用户ID
}

// 获取用户余额
func GetBalance(c *gin.Context) {
	var req GetBalanceRequest
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

	// 获取用户钱包信息
	wallet, err := models.GetWalletByUserID2(db, req.UserID)
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
			"user_id": req.UserID,
			"balance": wallet.Balance,
		},
		"message": "successful",
	}

	c.JSON(http.StatusOK, response)
}
