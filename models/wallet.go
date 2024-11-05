/**
 * 钱包数据库操作
 */

package models

import (
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
)

// 自定义错误类型
var (
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// 钱包结构体
type Wallet struct {
	ID      int     `json:"id"`      // 主键ID
	UserID  string  `json:"user_id"` // 用户ID
	Balance float64 `json:"balance"` // 余额
}

// 获取用户钱包余额等信息 事务
func GetWalletByUserID(tx *sql.Tx, userID string) (*Wallet, error) {
	var wallet Wallet
	err := tx.QueryRow("SELECT id, user_id, balance FROM wallets WHERE user_id = $1", userID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrWalletNotFound
		}
		logrus.Errorf("Error getting wallet by user ID: %v", err)
		return nil, err
	}
	return &wallet, nil
}

// 获取用户钱包余额等信息 非事务
func GetWalletByUserID2(db *sql.DB, userID string) (*Wallet, error) {
	var wallet Wallet
	err := db.QueryRow("SELECT id, user_id, balance FROM wallets WHERE user_id = $1", userID).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrWalletNotFound
		}
		logrus.Errorf("Error getting wallet by user ID: %v", err)
		return nil, err
	}
	return &wallet, nil
}

// 更新钱包余额 事务
func UpdateWalletBalance(tx *sql.Tx, userID string, newBalance float64) error {
	_, err := tx.Exec("UPDATE wallets SET balance = $1 WHERE user_id = $2", newBalance, userID)
	if err != nil {
		logrus.Errorf("Error updating wallet balance: %v", err)
		return err
	}
	return nil
}

// 创建用户钱包 事务
func CreateWallet(tx *sql.Tx, userID string) error {
	_, err := tx.Exec("INSERT INTO wallets (user_id, balance) VALUES ($1, 0)", userID)
	if err != nil {
		logrus.Errorf("Error creating wallet: %v", err)
		return err
	}
	return nil
}
