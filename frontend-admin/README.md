# PhotoSet 管理后台

PhotoSet 摄影套图平台管理后台，全部功能开发完成。

## 基本信息

- **技术栈**: Vue 3 + Element Plus + Vite + Pinia + ECharts
- **端口**: 3001
- **后端 API**: Go + Gin (端口 8080)
- **状态**: ✅ 生产就绪

## 快速开始

```bash
cd frontend-admin
npm install
npm run dev          # 开发模式 (http://localhost:3001)
npm run build        # 生产构建
npm run preview      # 预览构建结果
```

## 功能模块

### 仪表盘 (Dashboard)
- ECharts 趋势折线图（7/14/30 天切换）
- 新用户 / 新订单 / 收入趋势
- 收入概览卡片（今日/7日/总收入、客单价）
- 最近订单表格

### 内容审核 (ContentReview)
- 套图列表审核
- 通过 / 拒绝 / 删除（支持批量操作）
- 状态筛选
- 编辑套图（标题/封面/图片/标签/分类/收费方式/状态）
- 回收站：恢复 / 永久删除（含批量）
- 导出 CSV

### 用户管理 (UserManage)
- 用户列表 + 搜索（昵称/邮箱）
- 角色筛选 + 状态筛选
- 封号 / 解封
- 角色修改（guest/user/member/creator/admin）
- 用户详情抽屉（发布套图数、订单数、消费总额）
- 用户积分调整

### 订单管理 (OrderManage)
- 订单列表 + 筛选 + 分页
- 订单详情
- 管理员无限期退款

### 标签管理 (TagManage)
- 标签 CRUD + 批量操作

### 分类管理 (CategoryManage)
- 分类 CRUD + 排序

### 社区管理（按独立页面拆分）
- 帖子管理（审核/置顶/加精/删除）
- 回帖管理（删除）
- 敏感词管理（CRUD + 热加载）
- 分类管理
- 举报处理（通过/警告/忽略）
- 用户积分管理（调整积分）
- 等级 / 成就 / 积分商城 / 通知管理

### 会员管理 (MembershipManage)
- 会员套餐 CRUD（价格 / 时长 / 权益）

### 开发者中心 (DeveloperCenter)
- API Key 管理
- API 文档 / 签名 URL 文档

### OAuth 应用管理 (OAuthManage)
- OAuth2 Client CRUD

### 站点设置 (SiteSettings)
- **基本信息**：站点名称、描述、Logo、备案号
- **SEO 设置**：关键词、描述
- **关于我**：个人简介、联系方式
- **邮件设置**：SMTP 配置
- **水印设置**：文字/图片水印
- **存储与 CDN**：存储类型切换（本地/S3/R2）、S3 配置、CDN 域名、测试连接
- **支付配置**：支付宝 / 微信支付连通性测试

### 系统监控 (SystemMonitor / PerformanceMonitor)
- 系统状态 / 健康检查 / 性能指标 / Goroutine 堆栈
- CPU / Heap Profile 下载
- 后端服务重启

### 回收站 (TrashManage) / 数据备份 (DataBackup) / IP 地理位置 (IPGeoManage)
- 软删除套图恢复 / 永久删除
- 数据库备份的创建 / 下载 / 删除
- IP 地理位置查询

### 页面管理 (Pages)
- 自定义静态页面 CRUD
- Slug / 标题 / 内容 / 状态管理

### 操作日志 (AdminLogs)
- 关键操作自动记录（封号/解封/审核/退款/角色修改）
- 日志列表 + 操作类型筛选 + 分页

## 项目结构

```
frontend-admin/
├── src/
│   ├── api/
│   │   ├── index.js        # 核心 API 封装
│   │   ├── pages.js        # 页面管理 API
│   │   └── community.js    # 社区管理 API
│   ├── layout/
│   │   └── AdminLayout.vue  # 主布局（侧边栏菜单）
│   ├── router/
│   │   └── index.js
│   ├── stores/
│   │   └── user.js
│   ├── views/
│   │   ├── Login.vue           # 登录
│   │   ├── Dashboard.vue       # 仪表盘
│   │   ├── ContentReview.vue   # 内容审核
│   │   ├── EditPhotoset.vue    # 编辑套图
│   │   ├── UserManage.vue      # 用户管理
│   │   ├── OrderManage.vue     # 订单管理
│   │   ├── TagManage.vue       # 标签管理
│   │   ├── CategoryManage.vue  # 分类管理
│   │   ├── MembershipManage.vue # 会员套餐
│   │   ├── SiteSettings.vue    # 站点设置
│   │   ├── Pages.vue           # 页面管理
│   │   ├── DeveloperCenter.vue # 开发者中心
│   │   ├── OAuthManage.vue     # OAuth 应用
│   │   ├── AdminLogs.vue       # 操作日志
│   │   ├── SystemMonitor.vue / PerformanceMonitor.vue  # 系统监控
│   │   ├── TrashManage.vue / DataBackup.vue / IPGeoManage.vue  # 回收站/备份/IP
│   │   ├── LevelManage.vue / AchievementManage.vue / PointsMallManage.vue / NotificationManage.vue  # 社区增强
│   │   └── CommunityPosts.vue / CommunityReplies.vue / CommunityKeywords.vue
│   │       / CommunityCategories.vue / CommunityReports.vue / CommunityUsers.vue  # 社区管理
│   └── App.vue
├── package.json
├── vite.config.js
└── README.md
```

## API 接口

管理后台专用接口（`/api/admin/*`）：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/stats` | 统计概览 |
| GET | `/admin/stats/trend` | 趋势数据（7/14/30天） |
| GET | `/admin/users` | 用户列表（支持搜索/筛选） |
| GET | `/admin/users/:id` | 用户详情 |
| PUT | `/admin/users/:id/ban` | 封号/解封 |
| PUT | `/admin/users/:id/role` | 修改角色 |
| GET | `/admin/orders` | 订单列表 |
| GET | `/admin/orders/:id` | 订单详情 |
| POST | `/admin/orders/:id/refund` | 管理员退款 |
| GET | `/admin/photosets` | 套图列表 |
| PUT | `/admin/photosets/:id/review` | 审核 |
| DELETE | `/admin/photosets/:id` | 删除套图 |
| GET/POST/PUT/DELETE | `/admin/tags` | 标签 CRUD |
| GET/POST/PUT/DELETE | `/admin/categories` | 分类 CRUD |
| GET/PUT | `/admin/settings` | 站点设置 |
| POST | `/admin/storage/test` | 测试存储连接 |
| GET | `/admin/storage/status` | 存储状态 |
| GET/POST/PUT/DELETE | `/admin/pages` | 页面管理 |
| GET | `/admin/logs` | 操作日志 |
| GET/PUT | `/photosets/:id` | 套图详情 / 编辑（管理端复用公开接口） |
| GET | `/tags` / `/categories` | 标签 / 分类（编辑页选项） |
| POST | `/upload/image` | 图片上传（编辑页） |
| GET/POST/PUT/DELETE | `/admin/memberships` | 会员套餐管理 |
| GET/POST/DELETE | `/admin/dev/api-keys` | 开发者 API Key |
| GET | `/admin/dev/api-docs` / `/admin/dev/sign-url-docs` | API 文档 |
| GET/POST/PUT/DELETE | `/admin/oauth/clients` | OAuth 应用管理 |
| GET | `/admin/system/{status,health,metrics,goroutines}` | 系统监控 |
| POST | `/admin/system/restart` | 重启服务 |
| GET | `/admin/system/profile/{cpu,heap}` | Profile 下载 |
| GET/POST/DELETE | `/admin/backups` | 数据备份 |

### 社区管理接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/community/posts` | 帖子列表 |
| PUT | `/admin/community/posts/:id/pin` | 置顶/取消置顶 |
| PUT | `/admin/community/posts/:id/essence` | 加精/取消加精 |
| PUT | `/admin/community/posts/:id/status` | 审核状态 |
| DELETE | `/admin/community/posts/:id` | 删除帖子 |
| GET | `/admin/community/replies` | 回帖列表 |
| DELETE | `/admin/community/replies/:id` | 删除回帖 |
| GET/POST/PUT/DELETE | `/admin/community/keywords` | 敏感词管理 |
| PUT | `/admin/community/keywords/reload` | 热加载敏感词 |
| GET/PUT | `/admin/community/reports` | 举报处理 |
| GET/PUT | `/admin/community/users` | 社区用户积分管理 |
| GET | `/admin/community/stats` | 社区统计 |

---

**最后更新**: 2026-11-17
**状态**: ✅ 生产就绪（两个前端 `npm run build` 均通过）
