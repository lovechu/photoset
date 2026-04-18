# PhotoSet 前端用户端

摄影套图浏览平台 Vue 3 用户端，全部功能开发完成。

## 基本信息

- **技术栈**: Vue 3 + Element Plus + Vite + Pinia
- **端口**: 3000
- **后端 API**: Go + Gin (端口 8080)
- **状态**: ✅ 生产就绪

## 快速开始

```bash
cd frontend
npm install
npm run dev          # 开发模式 (http://localhost:3000)
npm run build        # 生产构建
npm run preview      # 预览构建结果
```

## 功能列表

### 用户系统
- 注册 / 登录 / 个人中心
- JWT Token 认证
- 多级权限（用户/创作者/管理员）

### 浏览与搜索
- 首页套图列表（瀑布流）
- 中文全文搜索（FULLTEXT ngram）
- 搜索关键词高亮
- 标签筛选 + 分类浏览

### 套图管理
- 创作者上传套图（支持 S3/R2 云存储）
- 编辑套图（标题/描述/标签/价格）
- 删除套图
- 付费/免费内容智能展示

### 订单与支付
- 购买套图
- 订单管理
- 48 小时内自助退款

### 个人中心
- 收藏管理
- 订单历史
- 个人资料编辑

### 会员系统
- 会员套餐展示
- 会员等级权益

### 静态页面
- 站点自定义页面（`/p/:slug`）

## 项目结构

```
frontend/
├── src/
│   ├── api/              # API 接口封装
│   ├── assets/           # 静态资源
│   ├── components/       # 公共组件
│   │   ├── AppHeader.vue
│   │   ├── AppFooter.vue
│   │   └── PhotosetCard.vue
│   ├── router/           # 路由配置
│   ├── stores/           # Pinia 状态管理
│   │   ├── user.js       # 用户状态
│   │   └── site.js       # 站点设置
│   ├── views/            # 页面
│   │   ├── Home.vue              # 首页
│   │   ├── PhotosetDetail.vue    # 套图详情
│   │   ├── Login.vue / Register.vue
│   │   ├── CreatePhotoset.vue    # 上传套图
│   │   ├── EditPhotoset.vue      # 编辑套图
│   │   ├── Profile.vue           # 个人中心
│   │   ├── Favorites.vue         # 收藏
│   │   ├── Orders.vue            # 订单
│   │   ├── Membership.vue        # 会员
│   │   ├── StaticPage.vue        # 静态页面
│   │   └── NotFound.vue
│   ├── App.vue
│   └── main.js
├── package.json
├── vite.config.js
└── README.md
```

## 页面路由

| 路径 | 说明 | 权限 |
|------|------|------|
| `/` | 首页（搜索/分类/标签） | 公开 |
| `/photoset/:id` | 套图详情 | 公开 |
| `/login` / `/register` | 登录/注册 | 公开 |
| `/create` | 上传套图 | 创作者+ |
| `/edit/:id` | 编辑套图 | 创作者（本人） |
| `/profile` | 个人中心 | 登录 |
| `/favorites` | 收藏 | 登录 |
| `/orders` | 订单 | 登录 |
| `/membership` | 会员套餐 | 公开 |
| `/p/:slug` | 静态页面 | 公开 |

## 开发配置

### 代理后端
默认 Vite 代理 `/api` → `http://localhost:8080`，修改 `vite.config.js` 即可。

### 安全特性
- Token 自动携带（Axios 拦截器）
- 路由守卫 + 后端双重权限校验
- URL 签名防盗链（HMAC-SHA256）
- XSS/CSRF 防护

---

**最后更新**: 2026-04-19
**状态**: ✅ 生产就绪
