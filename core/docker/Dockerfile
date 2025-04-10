# 第一阶段：构建Go二进制文件
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制Go模块文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目文件到容器
COPY . .

# 构建应用程序（注意main.go的路径）
RUN CGO_ENABLED=0 GOOS=linux go build -o haze-detection ./main

# 第二阶段：生成最终轻量级镜像
FROM alpine:3.19

# 安装必要的运行时依赖（如需要TZ/Ca证书）
RUN apk add --no-cache tzdata ca-certificates

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件和配置文件
COPY --from=builder /app/haze-detection .
COPY --from=builder /app/config ./config
COPY --from=builder /app/sql ./sql

# 暴露服务端口（根据你的实际端口修改）
EXPOSE 8248

# 启动应用程序
CMD ["./haze-detection"]