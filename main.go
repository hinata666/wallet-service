/**
 * 主程序入口
 */
package main

import (
	"log"
	"wallet-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 日志记录
	r.Use(gin.Logger())

	// 路由
	r.POST("/deposit", handlers.Deposit)              // 存款
	r.POST("/withdraw", handlers.Withdraw)            // 取款
	r.POST("/transfer", handlers.Transfer)            // 转账
	r.POST("/balance", handlers.GetBalance)           // 获取指定用户的余额
	r.POST("/transactions", handlers.GetTransactions) // 获取指定用户的交易历史

	// 启动服务器
	log.Fatal(r.Run(":8080"))
}
