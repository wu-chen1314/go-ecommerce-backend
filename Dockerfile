# ================= Build 阶段 =================
# 使用官方 Go 镜像作为构建环境
FROM golang:alpine AS builder

# 设置工作目录
WORKDIR /app

# 1. 先拷贝依赖文件 (利用 Docker 缓存机制，加速构建)
COPY go.mod go.sum ./
# 下载依赖 (如果有网络问题，可以在这里配置 GOPROXY)
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

# 2. 拷贝源代码
COPY . .

# 3. 编译成二进制文件 (名字叫 main)
RUN go build -o main .

# ================= Run 阶段 =================
# 使用最小的 Linux 镜像 (Alpine) 作为运行环境
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从 Build 阶段把编译好的二进制文件拷贝过来
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./main"]