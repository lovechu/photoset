# ========== Go 后端服务 ==========
# 多阶段构建：编译阶段 + 运行阶段

# 阶段1：Go 编译
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Alpine apk 源：自动探测官方源可达性，国内环境自动回退阿里云源。
#   默认 auto：5 秒内可访问 dl-cdn.alpinelinux.org 就用官方源（海外直连最快），
#             否则 sed 替换为 mirrors.aliyun.com（国内可用）。
#   强制指定：--build-arg APK_MIRROR=aliyun 或 =official
ARG APK_MIRROR=auto
RUN if [ "$APK_MIRROR" = "auto" ]; then \
      if wget -q -O /dev/null -T 5 -t 1 https://dl-cdn.alpinelinux.org/ 2>/dev/null; then \
        echo ">>> Alpine repo: official (dl-cdn.alpinelinux.org)"; \
      else \
        sed -i 's|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g' /etc/apk/repositories; \
        echo ">>> Alpine repo: China mirror (mirrors.aliyun.com)"; \
      fi; \
    elif [ "$APK_MIRROR" = "aliyun" ]; then \
      sed -i 's|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g' /etc/apk/repositories; \
      echo ">>> Alpine repo: China mirror (mirrors.aliyun.com) [forced]"; \
    fi && \
    apk add --no-cache git ca-certificates tzdata

# Go 依赖代理：构建时自动探测，选择当前网络最快的源。
#   默认 auto：5 秒内可访问官方 proxy.golang.org 就用它（海外直连最快），
#             否则自动回退 goproxy.cn（国内可用）。
#   强制指定：--build-arg GOPROXY=https://goproxy.cn,direct
# 校验保持默认 sum.golang.org（goproxy.cn 会代理 sumdb），确保依赖完整性可验证。
ARG GOPROXY=auto
RUN if [ "$GOPROXY" = "auto" ]; then \
      if wget -q -O /dev/null -T 5 -t 1 https://proxy.golang.org/github.com/gin-gonic/gin/@v/list 2>/dev/null; then \
        GOPROXY="https://proxy.golang.org,direct"; echo ">>> Go proxy: official (proxy.golang.org)"; \
      else \
        GOPROXY="https://goproxy.cn,direct"; echo ">>> Go proxy: China mirror (goproxy.cn)"; \
      fi; \
    fi; \
    echo ">>> GOPROXY=$GOPROXY"; go env -w GOPROXY="$GOPROXY"

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY internal/ ./internal/
COPY cmd/ ./cmd/
COPY scripts/ ./scripts/
COPY data/ ./data/
COPY docs/ ./docs/

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

# Alpine apk 源：与编译阶段同样的自动探测（官方源 / 阿里云源）。
ARG APK_MIRROR=auto
RUN if [ "$APK_MIRROR" = "auto" ]; then \
      if wget -q -O /dev/null -T 5 -t 1 https://dl-cdn.alpinelinux.org/ 2>/dev/null; then \
        echo ">>> Alpine repo: official (dl-cdn.alpinelinux.org)"; \
      else \
        sed -i 's|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g' /etc/apk/repositories; \
        echo ">>> Alpine repo: China mirror (mirrors.aliyun.com)"; \
      fi; \
    elif [ "$APK_MIRROR" = "aliyun" ]; then \
      sed -i 's|dl-cdn.alpinelinux.org|mirrors.aliyun.com|g' /etc/apk/repositories; \
      echo ">>> Alpine repo: China mirror (mirrors.aliyun.com) [forced]"; \
    fi && \
    apk add --no-cache ca-certificates tzdata wget && \
    adduser -D -s /bin/sh -h /app nonroot && \
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