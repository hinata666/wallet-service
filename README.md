### Wallet Service

这是一个简单的钱包服务，支持存款、取款、转账、获取指定用户的余额和交易历史等功能。

#### 实现的 RESTful API

1. 向指定用户的钱包存款<br>
2. 从指定用户的钱包取款<br>
3. 从一个用户向另一个用户转账<br>
4. 获取指定用户的余额<br>
5. 获取指定用户的交易历史<br>

#### 技术栈及描述

语言: Go<br>
数据库: PostgreSQL<br>
存款、取款、转账使用了协程、锁机制，事务保证并发安全。<br>

#### 目录结构

```sh
├── apitest
│   ├── deposit_test.go                # 存款处理单元测试
│   ├── get_balance_test.go            # 获取指定用户的余额单元测试
│   ├── get_transactions_test.go       # 获取指定用户的交易历史单元测试
│   ├── transfer_test.go               # 转账处理单元测试
│   └── withdraw_test.go               # 取款处理单元测试
├── handlers
│   ├── deposit.go                     # 存款处理
│   ├── get_balance.go                 # 获取指定用户的余额
│   ├── get_transactions.go            # 获取指定用户的交易历史
│   ├── transfer.go                    # 转账处理
│   ├── utils.go                       # 公共文件
│   └── withdraw.go                    # 取款处理
├── models
│   ├── transaction.go                 # 交易记录数据库操作
│   └── wallet.go                      # 钱包数据库操作
├── sql
│   └── create.sql                     # 在 PostgreSQL 数据库中添加库，创建表并添加字段注释，创建测试用户SQL脚本
├── utils
│   └── db.go                          # 数据库操作方法
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go                            #主程序入口
└── README.md
```

### 部署运行及测试

#### 1、添加库创建表

在 PostgreSQL 数据库中添加库，创建表并添加字段注释，创建测试数据

##### 1.1 创建库 wallet-db

```sql
-- 创建库 wallet-db
CREATE DATABASE wallet_db;
```

##### 1.2 创建钱包表

```sql
-- 创建钱包表
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0
);

COMMENT ON COLUMN wallets.id IS '钱包唯一标识';
COMMENT ON COLUMN wallets.user_id IS '用户唯一标识';
COMMENT ON COLUMN wallets.balance IS '用户余额';
```

##### 1.3 创建交易记录表

```sql
-- 创建交易记录表
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
```

##### 1.4 创建测试数据

```sql
-- 创建测试数据
INSERT INTO "wallets" VALUES (1, 'user1', 100.00);
INSERT INTO "wallets" VALUES (2, 'user2', 100.00);
INSERT INTO "wallets" VALUES (3, 'user3', 100.00);
INSERT INTO "wallets" VALUES (4, 'user4', 100.00);
INSERT INTO "wallets" VALUES (5, 'user5', 100.00);
```

#### 2、运行

##### 2.1 本地运行

1）安装依赖

```sh
go mod tidy
```

2）修改数据库配置 postgres://user:password@host:port/dbname?sslmode=disable 改成实际配置，文件路径 wallet-service/utils/db.go

```sh
cd wallet-service
vim ./utils/db.go

func init() {
	var err error
	// 连接字符串格式: postgres://user:password@host:port/dbname?sslmode=disable
	db, err = sql.Open("postgres", "postgres://user:password@localhost:5432/wallet_db?sslmode=disable")
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
```

3）运行服务

```sh
go run main.go
```

##### 2.2 使用 Docker 运行

###### 2.2.1 运行 Dockerfile

1）构建镜像

```sh
docker build -t wallet-service .
```

2）运行容器

```sh
docker run -d -p 8080:8080 wallet-service
```

###### 2.2.2 运行 docker-compose.yml

1）构建和启动服务

```sh
docker-compose up --build
```

2）仅启动服务

```sh
docker-compose up
```

3）停止服务

```sh
docker-compose down
```

#### 3、接口文档

##### 3.1 向指定用户的钱包存款

URL: /deposit<br>
Method: POST<br>
Request Body:<br>

```json
{
    "user_id": "user1",
    "amount": 50
}
```

Response:

```json
{
    "code": 200,
    "data": {
        "new_balance": 150,
        "user_id": "user1"
    },
    "message": "successful"
}
```

CURL 测试命令:

```sh
curl -d '{ "user_id": "user1", "amount": 50 }' http://127.0.0.1:8080/deposit
```

##### 3.2 从指定用户的钱包取款

URL: /withdraw<br>
Method: POST<br>
Request Body:<br>

```json
{
    "user_id": "user1",
    "amount": 50
}
```

Response:

```json
{
    "code": 200,
    "data": {
        "new_balance": 100,
        "user_id": "user1"
    },
    "message": "successful"
}
```

CURL 测试命令:

```sh
curl -d '{ "user_id": "user1", "amount": 50 }' http://127.0.0.1:8080/withdraw
```

##### 3.3 从一个用户向另一个用户转账

URL: /transfer<br>
Method: POST<br>
Request Body:<br>

```json
{
    "sender_id": "user2",
    "receiver_id": "user1",
    "amount": 50
}
```

Response:

```json
{
    "code": 200,
    "data": {
        "new_receiver_balance": 150,
        "new_sender_balance": 50,
        "receiver_id": "user1",
        "sender_id": "user2"
    },
    "message": "successful"
}
```

CURL 测试命令:

```sh
curl -d '{ "sender_id": "user2", "receiver_id": "user1", "amount": 50 }' http://127.0.0.1:8080/transfer
```

##### 3.4 获取指定用户的余额

URL: /balance<br>
Method: POST<br>
Request Body:<br>

```json
{
    "user_id": "user1"
}
```

Response:

```json
{
    "code": 200,
    "data": {
        "balance": 150,
        "user_id": "user1"
    },
    "message": "successful"
}
```

CURL 测试命令:

```sh
curl -d '{ "user_id": "user1" }' http://127.0.0.1:8080/balance
```

##### 3.5 获取指定用户的交易历史

URL: /transactions<br>
Method: POST<br>
Request Body:<br>

```json
{
    "user_id": "user1",
    "page": 1,
    "page_size": 2
}
```

Response:

```json
{
    "code": 200,
    "data": {
        "transactions": [
            {
                "id": 19,
                "from_user_id": "user2",
                "to_user_id": "user1",
                "amount": -50,
                "created_at": "2024-11-05T14:49:22.704832Z"
            },
            {
                "id": 20,
                "from_user_id": "user1",
                "to_user_id": "user2",
                "amount": 50,
                "created_at": "2024-11-05T14:49:22.704832Z"
            }
        ],
        "user_id": "user1"
    },
    "message": "successful"
}
```

CURL 测试命令:

```sh
curl -d '{ "user_id": "user1", "page": 1, "page_size": 2 }' http://127.0.0.1:8080/transactions
```

#### 4、单元测试

测试代码位于 wallet-service/apitest 测试方法如下

```sh
cd wallet-service/apitest
```

##### 4.1 向指定用户的钱包存款命令

```sh
# 参数 {"user_id": "user3", "amount": 100} 可根据业务需求进行修改
go test -run TestDeposit -v
```

##### 4.2 从指定用户的钱包取款命令

```sh
# 参数 {"user_id": "user3", "amount": 50} 可根据业务需求进行修改
go test -run TestWithdraw -v
```

##### 4.3 从一个用户向另一个用户转账命令

```sh
# 参数
# {
# 	 "sender_id": "user3",
# 	 "receiver_id": "user4",
# 	 "amount": 150
# }
# 可根据业务需求进行修改
go test -run TestTransfer -v
```

##### 4.4 获取指定用户的余额命令

```sh
# 参数 {"user_id": "user1"} 可根据业务需求进行修改
go test -run TestGetBalance -v
```

##### 4.5 获取指定用户的交易历史命令

```sh
# 参数
# {
# 	 "user_id": "user1",
# 	 "page": 1,
# 	 "page_size": 2
# }
# 可根据业务需求进行修改
go test -run TestGetTransactions -v
```

#### 5、其他事项

##### 5.1 如何设置和运行您的代码

可参考 README.md 文件中的 部署运行及测试

##### 5.2 解释您所做的任何决策

存款、取款、转账使用了协程、锁机制，事务保证并发安全。

##### 5.3 强调评审者应如何查看您的代码

可参考 README.md 文件中的 目录结构 及其注释

##### 5.4 花费的时间

24 小时

##### 5.5 单元测试

可参考 README.md 文件中的 3、接口文档 及 4、单元测试

##### 5.6 将使用 golangci-lint v1.61.0 进行代码检查

##### 5.7 日志记录

##### 5.8 小数处理

在写入库之前四舍五入已将小数点后保留两位，当然这只是一个简单的处理，实际业务中可根据业务需求进行修改，比如数据库保持小数点后 3 位，接口输出只展示小数点后两位。

##### 5.9 Dockerfile

可参考 Dockerfile 文档

##### 5.10 Docker-compose.yml

可参考 Docker-compose.yml 文档

##### 5.11 检查 goroutine 泄漏（使用 uber-go/goleak）

#### 6、加分事项

##### 6.1 GitHub Actions CI 管道（在 PR 上运行测试、构建镜像等）

##### 6.2 性能基准测试

##### 6.3 API 示例（首选 Postman）

可参考 README.md 文件中的 3、接口文档

##### 6.4 E2E API 测试（首选 Postman）

可参考 README.md 文件中的 3、接口文档
