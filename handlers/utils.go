package handlers

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// 验证器实例
var validate = validator.New()

// 用户余额锁映射
var userLocks = sync.Map{}

// 处理错误并返回 HTTP 响应
func HandleError(c *gin.Context, err error, status int) {
	logrus.Errorf("Error: %v", err)
	response := gin.H{
		"code":    status,
		"message": err.Error(),
		"data":    nil,
	}
	c.JSON(status, response)
}
