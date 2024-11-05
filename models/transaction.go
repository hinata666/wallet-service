/**
 * 交易记录数据库操作
 */

package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

// 自定义错误类型
var (
	ErrTransactionFailed = errors.New("transaction failed")
)

// 交易记录
type Transaction struct {
	ID         int       `json:"id"`           // 主键ID
	FromUserID string    `json:"from_user_id"` // 付款方用户ID
	ToUserID   string    `json:"to_user_id"`   // 收款方用户ID
	Amount     float64   `json:"amount"`       // 金额
	CreatedAt  time.Time `json:"created_at"`   // 创建时间
}

// 创建交易记录 事务
func CreateTransaction(tx *sql.Tx, fromUserID, toUserID string, amount float64) error {
	_, err := tx.Exec("INSERT INTO transactions (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)", fromUserID, toUserID, amount)
	if err != nil {
		logrus.Errorf("Error creating transaction: %v", err)
		return err
	}
	return nil
}

// 获取指定用户的交易记录
func GetTransactionsByUserID(db *sql.DB, userID string, page int, pageSize int) ([]Transaction, error) {
	var transactions []Transaction

	offset := (page - 1) * pageSize

	rows, err := db.Query(`
		SELECT id, from_user_id, to_user_id, amount, created_at
		FROM transactions
		WHERE from_user_id = $1 OR to_user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.FromUserID, &t.ToUserID, &t.Amount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, ErrWalletNotFound
	}

	return transactions, nil
}
