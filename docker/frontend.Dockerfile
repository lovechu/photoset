# ========== 用户端 Vue 3 前端 ==========
# 多阶段构建：构建阶段 + nginx 服务

# 阶段1：Node.js 构建
FROM node:18-alpine AS builder

ARG VITE_API_BASE_URL=http://backend:8080

WORKDIR /app

# 复制包管理文件并安装依赖
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci --only=production

# 复制源代码
COPY frontend/ ./

# 设置环境变量
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV NODE_ENV=production

# 构建项目
RUN npm run build

# 阶段2：Nginx 服务
FROM nginx:alpine

# 复制构建产物到 nginx 目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制自定义 nginx 配置
COPY docker/nginx-frontend.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:80 || exit 1

# Nginx 默认命令即可
CMD ["nginx", "-g", "daemon off;"]

# 补充说明：
# - 构建时可以通过 --build-arg 指定 VITE_API_BASE_URL
# - 例如：docker build --build-arg VITE_API_BASE_URL=https://api.example.com -t frontend .