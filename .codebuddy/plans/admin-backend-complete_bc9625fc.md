---
name: admin-backend-complete
overview: 完整完善管理后台：新增等级配置/成就/积分商城管理的后端Admin API + 前端管理页面，新增通知管理页面，Dashboard增加社区数据，修复路由重复Bug，站点设置增加Logo配置。
todos:
  - id: phase1-backend
    content: "Phase 1: 后端 - 新建 level_handler.go，实现等级配置/成就/积分商城商品/兑换记录的 Admin CRUD API"
    status: completed
  - id: phase1-routes
    content: "Phase 1: 后端 - 在 community_routes.go 的 adminGroup 中注册等级/成就/积分商城管理路由"
    status: completed
    dependencies:
      - phase1-backend
  - id: phase1-api
    content: "Phase 1: 前端 - 新建 api/level.js，实现等级/成就/积分商城的 Admin API 导出函数"
    status: completed
    dependencies:
      - phase1-backend
  - id: phase1-pages
    content: "Phase 1: 前端 - 新建 LevelManage.vue 等级配置管理页面（表格+编辑对话框）"
    status: completed
    dependencies:
      - phase1-api
  - id: phase1-achievement
    content: "Phase 1: 前端 - 新建 AchievementManage.vue 成就管理页面（表格+新增/编辑对话框+删除）"
    status: completed
    dependencies:
      - phase1-api
  - id: phase1-points
    content: "Phase 1: 前端 - 新建 PointsMallManage.vue 积分商城管理页面（表格+新增/编辑对话框+删除）"
    status: completed
    dependencies:
      - phase1-api
  - id: phase1-router
    content: "Phase 1: 前端 - 在 router/index.js 和 AdminLayout.vue 中添加等级/成就/积分商城路由和侧边栏入口"
    status: completed
    dependencies:
      - phase1-pages
      - phase1-achievement
      - phase1-points
  - id: phase2-notif-backend
    content: "Phase 2: 后端 - 在 community_handler.go 中新增通知管理方法（发送系统通知、通知列表查询）"
    status: completed
  - id: phase2-notif-route
    content: "Phase 2: 后端 - 在 adminGroup 中注册通知管理路由"
    status: completed
    dependencies:
      - phase2-notif-backend
  - id: phase2-notif-page
    content: "Phase 2: 前端 - 新建 NotificationManage.vue 通知管理页面 + API + 路由 + 侧边栏入口"
    status: completed
    dependencies:
      - phase2-notif-backend
  - id: phase2-dashboard
    content: "Phase 2: 前端 - 修改 Dashboard.vue，添加社区数据卡片（帖子数、回帖数、待处理举报等）"
    status: completed
  - id: phase3-router-bug
    content: "Phase 3: 前端 - 修复 router/index.js 中 community/categories 路由重复定义 Bug"
    status: completed
  - id: phase3-site-logo
    content: "Phase 3: 前端 - 修改 SiteSettings.vue，添加 logo_url 和 favicon_url 配置字段"
    status: completed
  - id: phase3-export
    content: "Phase 3: 前端 - 在套图管理页面添加 CSV 导出功能（前端 API + 按钮 + 后端路由）"
    status: completed
  - id: build-verify
    content: 验证 - 后端 go build 编译通过，前端无语法错误
    status: completed
    dependencies:
      - phase1-router
      - phase2-dashboard
      - phase3-router-bug
      - phase3-site-logo
      - phase3-export
---

## 需求概述

PhotoSet 管理后台功能完善，涉及后端 Go API 开发和前端 Vue3 管理页面开发，共三个阶段。

## 核心功能

### Phase 1: 用户等级/成就/积分商城管理（最大功能缺口）

- **等级配置管理页面** - 管理员可编辑 10 个等级的名称、图标、颜色、积分范围、权限配置（发帖/回复/上传限制等）、升级奖励
- **成就管理页面** - 管理员可创建/编辑/删除成就，设置条件类型（发帖数/回复数/获赞数/等级等）和奖励
- **积分商城管理页面** - 管理员可创建/编辑/删除商品（徽章/称号/VIP天数/特权），管理库存、定价、最低等级要求
- **兑换记录查看** - 管理员可查看所有用户的积分兑换历史
- 后端需要新增 Admin CRUD API（等级配置、成就、积分商城商品、兑换记录）
- 前端需要新增 3 个管理页面 + API 导出 + 路由 + 侧边栏入口

### Phase 2: 通知管理 + Dashboard 完善

- **系统通知管理** - 管理员可发送系统公告/通知给全部用户或指定用户，查看通知发送历史
- **Dashboard 社区数据卡片** - 在数据看板中添加社区统计（帖子数、回帖数、待处理举报等）

### Phase 3: Bug 修复与小功能补充

- **修复路由重复 Bug** - 前端 router/index.js 中 community/categories 路由重复定义
- **站点设置增加 Logo 配置** - 表单添加 logo_url 和 favicon_url 字段
- **套图导出功能** - 在套图管理页面添加 CSV 导出按钮

## 技术栈

- **后端**: Go + Gin + GORM（已有项目技术栈）
- **前端**: Vue 3 + Element Plus + Vite（已有管理后台技术栈）
- **语言**: Go（后端）、JavaScript/Vue（前端）

## 实现方案

### 后端实现策略

**新增 Admin Handler**: 在 `internal/http/handlers/admin/` 目录下创建新文件 `level_handler.go`，集中处理等级配置、成就、积分商城商品和兑换记录的 Admin CRUD 操作。复用已有的 Repository 层方法（`UserLevelRepository`、`AchievementRepository`、`PointsMallRepository` 已有完备的 CRUD 方法）。

**通知管理 API**: 在 `internal/http/handlers/admin/community_handler.go` 中扩展通知相关的 Admin 方法，复用已有的 `NotificationRepository`。

**路由注册**: 在 `internal/http/routes/community_routes.go` 的 adminGroup 中新增等级/成就/积分商城/通知管理路由。在 `internal/http/routes/routes.go` 的 admin group 中新增套图导出路由。

**设计原则**:

- 复用已有 Repository 方法，不重复造轮子
- 遵循现有 Admin Handler 的编码风格（`recordLog` 记录操作日志）
- 等级配置不允许删除（预置 10 级），只允许编辑
- 成就和积分商城商品支持完整 CRUD
- 所有 Admin API 都需要 admin 权限中间件保护

### 前端实现策略

**新增 API 导出函数**: 在 `frontend-admin/src/api/level.js` 中创建等级/成就/积分商城的 Admin API。在 `frontend-admin/src/api/community.js` 中补充通知管理 API。在 `frontend-admin/src/api/index.js` 中补充套图导出 API。

**新增管理页面**: 在 `frontend-admin/src/views/` 下新增 4 个页面文件，复用现有 Element Plus 表格+对话框的 CRUD 模式。

**路由和侧边栏**: 在 `router/index.js` 中新增路由条目（同时修复重复路由 Bug）。在 `AdminLayout.vue` 侧边栏中添加等级管理、成就管理、积分商城、通知管理入口。

**修改现有页面**:

- `Dashboard.vue` - 添加社区数据卡片
- `SiteSettings.vue` - 添加 Logo 配置字段
- `ContentReview.vue`（或套图列表页面） - 添加导出按钮

## 关键文件

| 操作 | 文件路径 | 说明 |
| --- | --- | --- |
| [新建] | `internal/http/handlers/admin/level_handler.go` | 等级/成就/积分商城 Admin Handler |
| [修改] | `internal/http/routes/community_routes.go` | 新增等级/成就/积分商城/通知管理路由 |
| [修改] | `internal/http/routes/routes.go` | 新增套图导出路由 |
| [修改] | `internal/http/handlers/admin/community_handler.go` | 新增通知管理方法 |
| [修改] | `frontend-admin/src/api/level.js` | [新建] 等级/成就/积分商城 Admin API |
| [修改] | `frontend-admin/src/api/community.js` | 新增通知管理 API |
| [修改] | `frontend-admin/src/api/index.js` | 新增套图导出 API |
| [新建] | `frontend-admin/src/views/LevelManage.vue` | 等级配置管理页面 |
| [新建] | `frontend-admin/src/views/AchievementManage.vue` | 成就管理页面 |
| [新建] | `frontend-admin/src/views/PointsMallManage.vue` | 积分商城管理页面 |
| [新建] | `frontend-admin/src/views/NotificationManage.vue` | 通知管理页面 |
| [修改] | `frontend-admin/src/router/index.js` | 新增路由 + 修复重复路由 Bug |
| [修改] | `frontend-admin/src/layout/AdminLayout.vue` | 侧边栏添加新入口 |
| [修改] | `frontend-admin/src/views/Dashboard.vue` | 添加社区数据卡片 |
| [修改] | `frontend-admin/src/views/SiteSettings.vue` | 添加 Logo 配置字段 |


## Agent Extensions

- **[subagent:code-explorer]**: 在实现每个阶段前，使用 code-explorer 搜索现有代码模式（如 admin handler 的 CRUD 模式、前端页面的表格+对话框模式），确保新代码与现有代码风格一致。