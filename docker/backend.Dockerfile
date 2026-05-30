# ========== Go 后端服务 ==========
# 多阶段构建：编译阶段 + 运行阶段

# 阶段1：Go 编译
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 配置 Go 国内镜像代理
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w GOSUMDB=off

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY internal/ ./internal/
COPY cmd/ ./cmd/
COPY scripts/ ./scripts/

# 启用 CGO，设置编译参数
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# 编译 Go 应用
RUN go build -ldflags="-w -s" -trimpath -o /app/main ./cmd/main.go

# 阶段2：运行环境（alpine 镜像，轻量安全）
FROM alpine:3.19

WORKDIR /app

# 从编译阶段复制二进制文件
COPY --from=builder /app/main /app/main

# 复制需要的脚本和配置
COPY --from=builder /app/scripts/ ./scripts/
COPY .env.example /app/.env.example

# 创建上传目录和日志目录
RUN mkdir -p /app/uploads /app/logs

# 设置时区
ENV TZ=Asia/Shanghai

# 安装 CA 证书、wget，创建 nonroot 用户
RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -S nonroot && \
    adduser -S -G nonroot -s /bin/sh nonroot && \
    chown -R nonroot:nonroot /app && \
    chmod -R g+w /app

# 复制并设置入口脚本（在切换用户之前）
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod 755 /entrypoint.sh

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 使用非 root 用户运行
USER nonroot:nonroot

# 运行应用（使用 entrypoint 处理权限）
ENTRYPOINT ["/entrypoint.sh"]

# 开发模式的备选方案：
# 需要热重载时使用 air，但生产用多阶段构建
# ```
# FROM golang:1.22-alpine AS dev
# RUN go install github.com/cosmtrek/air@latest
# COPY . .
# CMD ["air", "-c", ".air.toml"]
# ```