# PhotoSet 摄影套图浏览平台

> 专业的摄影套图内容浏览与会员订阅平台

---

## 项目概述

| 项目 | 说明 |
|------|------|
| **项目名称** | PhotoSet |
| **技术栈** | Go + Gin / Vue 3 / MySQL + Redis / Docker |
| **项目状态** | ✅ 已完成部署 |
| **访问地址** | https://tt.cy.mk |

---

## 技术架构

### 后端 (Go + Gin)
- RESTful API 设计
- JWT 认证
- MySQL 8.0 + GORM
- Redis 缓存
- FULLTEXT 中文全文搜索
- URL 签名防盗链
- Cloudflare R2 对象存储支持

### 前端用户端 (Vue 3)
- 响应式设计
- 高级搜索过滤器
- 会员订阅系统
- 收藏功能
- 搜索关键词高亮

### 管理后台 (Vue 3)
- 内容审核
- 用户管理
- 订单管理
- 标签管理
- Dashboard 统计

### Phase 6 功能 (站点设置)
- 基本信息管理
- SEO 关键词设置
- 关于页面管理
- 邮件配置
- 水印设置

---

## 功能模块

| 模块 | 功能 | 状态 |
|------|------|------|
| 用户系统 | 注册/登录/资料/角色 | ✅ |
| 套图管理 | CRUD/分类/标签/搜索 | ✅ |
| 会员系统 | 套餐/订单/支付/退款 | ✅ |
| 收藏功能 | 收藏/取消收藏 | ✅ |
| 高级搜索 | 价格/分类/时间/排序 | ✅ |
| 管理后台 | 审核/用户/统计 | ✅ |
| 站点设置 | 基本/SEO/关于/邮件/水印 | ✅ |
| Redis 缓存 | 列表5min/详情10min/标签30min | ✅ |
| FULLTEXT 搜索 | MySQL ngram 中文分词 | ✅ |
| URL 签名 | HMAC-SHA256 防盗链 | ✅ |
| 订单退款 | 用户48h自助 + admin无限 | ✅ |

---

## 部署信息

### 生产环境

### Docker 部署路径
```
/opt/photoset/
```

---

## 宝塔 Nginx 伪静态规则

> 将以下规则添加到宝塔面板 → 网站 → tt.cy.mk → 伪静态

### 1. 上传文件代理（uploads）

```nginx
location /uploads/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

### 2. 后端 API 代理

```nginx
location ~ ^/api/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

### 3. 管理后台静态资源

```nginx
location /admin/assets/ {
    proxy_pass http://127.0.0.1:3001/assets/;
    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

### 4. 管理后台主路由

```nginx
location /admin/ {
    proxy_pass http://127.0.0.1:3001/;
    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

### 规则优先级

```
1. /admin/assets/  (静态资源优先)
2. /admin/         (管理后台)
3. /uploads/       (上传文件)
4. /api/           (后端API)
5. /               (前端默认)
```

---

## Docker 运维命令

### 查看容器状态
```bash
docker ps
```

### 查看日志
```bash
docker logs -f photoset-backend
docker logs -f photoset-frontend
docker logs -f photoset-admin
```

### 重启容器
```bash
docker restart photoset-backend
docker restart photoset-frontend
docker restart photoset-admin
```

### 进入容器
```bash
docker exec -it photoset-backend sh
```

### 备份数据库
```bash
docker exec photoset-mysql mysqldump -u root -p<密码> photoset > backup_$(date +%Y%m%d).sql
```

---

## 开发环境

### 启动后端 (WSL)
```bash
cd ~/projects/20260408115223/backend/
go run cmd/main.go
# 端口: 8080
```

### 启动前端 (Windows)
```powershell
cd C:\Users\ichuy\WorkBuddy\20260408115223\backend\frontend
npm run dev
# 端口: 3000
```

### 启动管理后台 (Windows)
```powershell
cd C:\Users\ichuy\WorkBuddy\20260408115223\backend\frontend-admin
npm run dev
# 端口: 3001
```

---

## 项目进度

| Phase | 内容 | 状态 |
|-------|------|------|
| Phase 1 | 后端基础架构 + 用户认证 | ✅ |
| Phase 2 | 套图核心功能 API | ✅ |
| Phase 3 | Web 前端完整实现 | ✅ |
| Phase 4 | 会员支付系统 + 管理后台 | ✅ |
| Phase 4 补齐 | Cloudflare R2 + 套图编辑 | ✅ |
| Phase 5 | Redis缓存/FULLTEXT/退款 | ✅ |
| Phase 6 | 站点设置 + 页面管理 | ✅ |

---

## 项目里程碑

- ✅ Phase 1-5 全部完成
- ✅ 管理后台功能完善
- ✅ Phase 6 站点设置上线
- 🔄 Flutter 移动端开发（可选）

---

*最后更新: 2026-04-19*
