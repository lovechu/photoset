# PhotoSet 独立 SEO 渲染服务
# Node 20 + 系统 Chromium（不依赖 puppeteer 自动下载，避免构建期网络波动）
FROM node:20-bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive \
    PUPPETEER_SKIP_DOWNLOAD=1 \
    CHROME_PATH=/usr/bin/chromium \
    CACHE_DIR=/var/cache/seo \
    NODE_ENV=production

# Chromium 及其中文字体/运行库
RUN apt-get update && apt-get install -y --no-install-recommends \
    chromium \
    fonts-noto-cjk \
    fonts-liberation \
    ca-certificates \
    libnss3 libatk-bridge2.0-0 libgtk-3-0 libgbm1 libasound2 libxss1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 依赖单独一层，锁定版本（package-lock 一并 COPY）
COPY package.json package-lock.json* ./
RUN npm ci --omit=dev --no-audit --no-fund

COPY server.js render.js cache.js queue.js ./

RUN mkdir -p /var/cache/seo && chown -R node:node /var/cache/seo

USER node

EXPOSE 3000

CMD ["node", "server.js"]
