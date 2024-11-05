/**
 * 数据库操作
 */

package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // 导入 PostgreSQL 驱动
)

var db *sql.DB

func init() {
	var err error
	// 连接字符串格式: postgres://user:password@host:port/dbname?sslmode=disable
	db, err = sql.Open("postgres", "postgres://postgres:Winer123!!@localhost:5432/wallet_db?sslmode=disable")
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Failed to ping database:", err)
		return
	}
	fmt.Println("Database connection established")
}

func GetDB() *sql.DB {
	return db
}
