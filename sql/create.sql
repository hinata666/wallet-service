-- 创建库 wallet-db
CREATE DATABASE wallet_db;

-- 创建钱包表
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0
);

COMMENT ON COLUMN wallets.id IS '钱包唯一标识';
COMMENT ON COLUMN wallets.user_id IS '用户唯一标识';
COMMENT ON COLUMN wallets.balance IS '用户余额';

-- 创建交易表
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_user_id VARCHAR(255),
    to_user_id VARCHAR(255),
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN transactions.id IS '交易唯一标识';
COMMENT ON COLUMN transactions.from_user_id IS '付款用户唯一标识';
COMMENT ON COLUMN transactions.to_user_id IS '收款用户唯一标识';
COMMENT ON COLUMN transactions.amount IS '交易金额';
COMMENT ON COLUMN transactions.created_at IS '交易创建时间';

-- 创建测试数据
INSERT INTO "wallets" VALUES (1, 'user1', 100.00);
INSERT INTO "wallets" VALUES (2, 'user2', 100.00);
INSERT INTO "wallets" VALUES (3, 'user3', 100.00);
INSERT INTO "wallets" VALUES (4, 'user4', 100.00);
INSERT INTO "wallets" VALUES (5, 'user5', 100.00);