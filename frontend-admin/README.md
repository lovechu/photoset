# PhotoSet 管理后台

PhotoSet 摄影套图浏览平台的管理后台前端系统。

## 🚀 项目概述

- **项目名称**: PhotoSet Admin Panel
- **技术栈**: Vue 3 + Element Plus + Vite
- **开发状态**: 🔥 **生产就绪 (Phase 5 全部完成)**
- **端口**: 3001
- **API 后端**: Go + Gin (端口: 8080)

## 📋 功能模块

### ✅ 已完成的功能

#### 1. **用户认证**
- 管理员登录页面 (`src/views/Login.vue`)
- Token 管理
- 权限验证拦截

#### 2. **仪表盘管理** (`src/views/Dashboard.vue`)
- 实时统计数据展示
- 套图上传量/订单量/销售额统计
- 用户增长趋势图表

#### 3. **内容审核** (`src/views/ContentReview.vue`)
- 套图列表审核
- 管理员删除功能（支持批量操作）
- 状态筛选（待审核/已审核/已拒绝）
- 快速预览套图详情

#### 4. **订单管理** (`src/views/OrderManage.vue`)
- 订单列表查看（支持筛选、分页）
- 订单状态管理
- **管理员退款功能**（支持 48 小时后无限期退款）
- 用户自助退款界面集成

#### 5. **标签管理** (`src/views/TagManage.vue`)
- 标签增删改查
- 批量操作支持
- 标签与套图关联管理

#### 6. **用户管理**
- 用户信息查看
- 权限管理（用户/创作者/管理员）
- 账号状态管理

### 🔧 Phase 5 新增功能

#### **1. 安全增强**
- URL 签名防盗链系统集成
- 管理员权限分级

#### **2. 性能优化**
- Redis 缓存系统对接
- 图片加载优化

#### **3. 搜索增强**
- FULLTEXT 中文全文搜索
- 前端搜索结果关键词高亮

## 🛠️ 开发指南

### 环境要求
- Node.js 18+
- npm 8+

### 安装依赖
```bash
npm install
```

### 开发模式
```bash
npm run dev
```
访问地址: http://localhost:3001

### 构建生产版本
```bash
npm run build
```

### 预览构建结果
```bash
npm run preview
```

## 📁 项目结构

```
frontend-admin/
├── src/
│   ├── api/            # API 接口调用
│   │   └── index.js    # API 封装（含退款/删除/标签等接口）
│   ├── assets/         # 静态资源
│   ├── components/     # 公共组件
│   ├── layout/         # 布局组件
│   │   └── AdminLayout.vue  # 管理后台主布局
│   ├── router/         # 路由配置
│   │   └── index.js    # 路由定义
│   ├── stores/         # Pinia 状态管理
│   │   └── user.js     # 用户状态
│   ├── views/          # 页面视图
│   │   ├── Dashboard.vue      # 仪表盘
│   │   ├── ContentReview.vue  # 内容审核
│   │   ├── OrderManage.vue    # 订单管理
│   │   ├── TagManage.vue      # 标签管理
│   │   └── Login.vue          # 登录页
│   └── App.vue         # 根组件
├── index.html          # HTML 入口
├── package.json        # 项目依赖
├── vite.config.js     # Vite 配置
└── README.md          # 项目说明（本文档）
```

## 🔌 API 接口文档

### 管理后台专用接口
已实现 6 个管理后台接口：

1. **GET /api/admin/orders** - 获取订单列表
2. **GET /api/admin/orders/:id** - 获取订单详情
3. **POST /api/admin/orders/:id/refund** - 管理员退款
4. **GET /api/admin/tags** - 获取标签列表
5. **POST /api/admin/tags** - 创建标签
6. **PUT /api/admin/tags/:id** - 更新标签
7. **DELETE /api/admin/tags/:id** - 删除标签
8. **DELETE /api/admin/photosets/:id** - 管理员删除套图
9. **GET /api/admin/stats** - 获取统计数据
10. **GET /api/admin/users** - 获取用户列表

### 接口封装位置
- `src/api/index.js` - 完整接口调用封装

## 🎨 技术特性

### 1. **组件化架构**
1. 基于 Vue 3 Composition API
2. Element Plus UI 组件库
3. SCSS 样式预处理

### 2. **状态管理**
1. Pinia 状态管理库
2. 响应式数据流
3. 模块化存储设计

### 3. **路由系统**
1. 基于 Vue Router
2. 路由守卫权限控制
3. 路由懒加载优化

### 4. **代码规范**
1. ESLint + Prettier 代码格式化
2. 组件解耦设计
3. 注释文档完善

### 5. **性能优化**
1. 路由懒加载
2. 组件按需引入
3. 图片懒加载
4. API 响应缓存

## 🔒 权限控制

### 用户角色
1. **普通用户** - 仅查看公开内容
2. **创作者 (Creator)** - 可以上传和管理自己的套图
3. **管理员 (Admin)** - 完整管理权限

### 访问控制
1. 路由级别权限验证
2. 页面元素权限控制
3. API 接口权限验证

## 🧪 测试与部署

### 开发测试
```bash
# 启动开发服务器
npm run dev

# 访问管理后台
http://localhost:3001
```

### 生产部署
1. 构建静态文件
   ```bash
   npm run build
   ```
2. 部署到 Web 服务器（Nginx、Apache 等）
3. 配置反向代理到后端 API

### 环境配置
```bash
# .env.development
VITE_API_BASE_URL=http://localhost:8080/api

# .env.production
VITE_API_BASE_URL=https://api.yourdomain.com/api
```

## 📈 开发进度

### Phase 功能完成情况
- ✅ **Phase 1** - 基础架构搭建（完成）
- ✅ **Phase 2** - 用户端核心功能（完成）
- ✅ **Phase 3** - 支付与订单系统（完成）
- ✅ **Phase 3.5** - Cloudflare R2 存储集成（完成）
- ✅ **Phase 4** - 管理后台基础功能（完成）
- ✅ **Phase 5** - 高级功能增强（完成）

### Phase 5 具体完成内容
1. **套图删除功能增强** - 管理员可以删除任意套图
2. **Redis 缓存系统** - 静默降级策略，提升性能
3. **FULLTEXT 搜索** - MySQL ngram 中文全文搜索
4. **URL 签名防盗链** - HMAC-SHA256 签名系统
5. **订单退款系统** - 用户48h自助 + 管理员无限期退款

## 🤝 开发指南

### 新增页面步骤
1. 在 `src/views/` 创建 Vue 组件
2. 在 `src/router/index.js` 配置路由
3. 在 `src/layout/AdminLayout.vue` 中添加菜单项
4. 在 `src/api/index.js` 中添加接口调用

### 代码规范
1. 组件使用 `<script setup>` 语法
2. CSS 使用 scoped 样式
3. API 调用统一使用导入的 `api` 对象
4. 状态管理优先使用 Pinia

## 📞 支持与维护

### 开发团队
- 后端：Go + Gin 服务
- 前端：Vue 3 + Element Plus
- 移动端：Flutter（可选）

### 常见问题
1. **跨域问题** - 后端已配置 CORS
2. **权限问题** - 检查 Token 和角色权限
3. **API 接口变更** - 前后端接口版本保持一致

## 📄 许可证

私有项目，仅供内部使用。

---

**最后更新**: 2026-04-09  
**项目状态**: ✅ **生产就绪 - Phase 5 全部完成**  
**维护者**: PhotoSet 开发团队