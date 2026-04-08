# PhotoSet 前端 Web

摄影套图浏览平台 Vue 3 前端项目。

## 技术栈

- **框架**: Vue 3 + Composition API
- **构建工具**: Vite
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **UI 组件库**: Element Plus
- **HTTP 客户端**: Axios

## 项目结构

```
frontend-web/
├── src/
│   ├── api/              # API 接口定义
│   ├── assets/           # 静态资源
│   ├── components/       # 公共组件
│   │   ├── AppHeader.vue # 顶部导航
│   │   ├── AppFooter.vue # 底部信息
│   │   └── PhotosetCard.vue # 套图卡片
│   ├── router/           # 路由配置
│   ├── stores/           # Pinia 状态管理
│   ├── utils/            # 工具函数
│   ├── views/            # 页面组件
│   │   ├── Home.vue          # 首页（套图列表）
│   │   ├── PhotosetDetail.vue # 套图详情
│   │   ├── Login.vue         # 登录
│   │   ├── Register.vue      # 注册
│   │   ├── CreatePhotoset.vue # 创建套图
│   │   └── NotFound.vue      # 404
│   ├── App.vue           # 根组件
│   └── main.js          # 入口文件
├── public/               # 公共静态文件
├── index.html           # HTML 入口
├── vite.config.js       # Vite 配置
└── package.json         # 依赖配置
```

## 快速开始

### 1. 安装依赖

```bash
cd frontend-web
npm install
```

### 2. 配置后端地址

项目默认使用 Vite 代理连接 `http://localhost:8080/api`。

如果后端运行在其他端口，修改 `vite.config.js`:

```js
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:你的端口号',
      changeOrigin: true
    }
  }
}
```

### 3. 启动开发服务器

```bash
npm run dev
```

访问 `http://localhost:3000`

### 4. 构建生产版本

```bash
npm run build
```

## 功能说明

### 已实现功能

| 功能 | 说明 |
|------|------|
| 首页套图列表 | 瀑布流展示，支持分页、标签筛选 |
| 套图详情页 | 图片画廊，支持付费内容裁剪 |
| 用户注册/登录 | 邮箱注册，JWT Token 鉴权 |
| 创建套图 | 仅 creator/admin 可访问（后端权限控制） |

### 页面路由

| 路径 | 说明 | 权限 |
|------|------|------|
| `/` | 首页，套图列表 | 公开 |
| `/photoset/:id` | 套图详情 | 公开（部分内容需登录） |
| `/login` | 登录页 | 游客 |
| `/register` | 注册页 | 游客 |
| `/create` | 创建套图 | creator/admin |
| `/*` | 404 页面 | 公开 |

### 角色权限

| 角色 | 说明 | 可创建套图 |
|------|------|-----------|
| `user` | 普通注册用户 | ❌ |
| `member` | 付费会员 | ❌ |
| `creator` | 创作者 | ✅ |
| `admin` | 管理员 | ✅ |

## API 接口对接

### 已对接接口

| 接口 | 方法 | 鉴权 | 说明 |
|------|------|------|------|
| `/api/health` | GET | ❌ | 健康检查 |
| `/api/auth/register` | POST | ❌ | 用户注册 |
| `/api/auth/login` | POST | ❌ | 用户登录 |
| `/api/auth/me` | GET | ✅ | 获取当前用户 |
| `/api/photosets` | GET | ❌ | 套图列表 |
| `/api/photosets/:id` | GET | 可选 | 套图详情 |
| `/api/photosets` | POST | ✅ | 创建套图 |
| `/api/tags` | GET | ❌ | 标签列表 |

### 待实现接口（后端后续开发）

- 图片上传接口
- 收藏功能
- 会员订阅/订单支付
- 用户资料编辑
- 管理后台

## 付费内容展示策略

```
套图详情页判断逻辑：

if (is_free === 1) {
  // 免费套图 → 直接展示所有图片
} else if (photos.length === 0) {
  // 付费套图，无权限 → 展示封面 + 购买提示
} else {
  // 付费套图，有权限 → 展示所有图片
}
```

## 开发注意事项

1. **Token 管理**: 登录后 Token 自动存储在 localStorage，请求拦截器自动携带
2. **权限控制**: 前端路由守卫 + 后端 API 双重权限验证
3. **响应处理**: 统一使用 request.js 拦截器处理响应和错误
4. **Element Plus**: 使用按需引入的图标组件

## 浏览器兼容性

- Chrome 80+
- Firefox 75+
- Safari 13+
- Edge 80+

## 许可证

MIT
