# 使用官方的 Go 运行时作为构建环境
FROM golang:1.18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod go.sum ./

# 下载所有依赖项
RUN go mod download

# 将应用程序代码复制到工作目录
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -o wallet-service .

# 使用 Alpine Linux 作为最终的运行时环境
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制可执行文件
COPY --from=builder /app/wallet-service .

# 暴露应用程序的端口
EXPOSE 8080

# 运行应用程序
CMD ["./wallet-service"]