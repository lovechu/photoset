# PhotoSet 摄影套图平台 - 前端用户端

摄影套图浏览平台 Vue 3 前端项目，**Phase 5 全部功能已完成**。

## 🚀 项目状态
- **开发状态**: 🔥 **生产就绪 (Phase 5 全部完成)**
- **技术栈**: Vue 3 + Element Plus + Vite
- **端口**: 3000
- **后端API**: Go + Gin (端口: 8080)

## 📋 核心功能概览

### ✅ **已完成的全部功能**

#### 1. **用户系统**
- 用户注册/登录/个人中心
- JWT Token 认证管理
- 三级权限系统（用户/创作者/管理员）

#### 2. **核心浏览功能**
- 首页瀑布流展示
- 智能搜索系统（支持中文全文检索）
- **新功能**: 搜索结果关键词高亮显示
- 标签筛选与分类浏览

#### 3. **套图管理**
- 创作者套图上传（支持Cloudflare R2存储）
- 套图编辑与删除（创作者管理自己的套图）
- 付费/免费内容智能展示

#### 4. **支付与订单系统**
- 集成支付宝网页支付
- 订单创建与管理
- 订单状态跟踪
- **新功能**: 48小时内用户自助退款

#### 5. **个人中心**
- 收藏管理功能
- 订单历史记录
- 个人偏好设置

#### 6. **系统增强 (Phase 5)**
- 性能优化：Redis 缓存系统对接
- 安全增强：URL 签名防盗链系统
- 搜索增强：FULLTEXT 中文全文搜索

## 📁 完整项目结构

```
frontend/
├── src/
│   ├── api/              # API 接口定义（已对接全部接口）
│   ├── assets/           # 静态资源
│   ├── components/       # 公共组件
│   │   ├── AppHeader.vue     # 顶部导航
│   │   ├── AppFooter.vue     # 底部信息
│   │   ├── PhotosetCard.vue  # 套图卡片（新增搜索高亮）
│   │   ├── OrderCard.vue     # 订单卡片组件
│   │   └── UploadProgress.vue # 上传进度组件
│   ├── router/           # 路由配置
│   ├── stores/           # Pinia 状态管理
│   │   ├── user.js       # 用户状态（三级权限）
│   │   └── photoset.js   # 套图全局状态
│   ├── utils/            # 工具函数
│   ├── views/            # 页面组件（完整实现）
│   │   ├── Home.vue              # 首页（套图列表 + 搜索）
│   │   ├── PhotosetDetail.vue    # 套图详情（购买/查看）
│   │   ├── Login.vue             # 登录页
│   │   ├── Register.vue          # 注册页
│   │   ├── CreatePhotoset.vue    # 创建套图（支持 R2上传）
│   │   ├── EditPhotoset.vue      # 编辑套图（创作者）
│   │   ├── Profile.vue           # 个人中心
│   │   ├── Favorites.vue         # 我的收藏
│   │   ├── Orders.vue            # 我的订单 + 退款操作
│   │   └── NotFound.vue          # 404 页面
│   ├── App.vue           # 根组件
│   └── main.js          # 入口文件
├── package.json         # 依赖配置
├── vite.config.js      # Vite 配置
└── README.md           # 本文档
```

## 🚀 快速开始

### 1. 环境准备
```bash
# 确保已安装 Node.js 18+ 和 npm 8+
node --version
npm --version
```

### 2. 安装依赖
```bash
cd frontend
npm install
```

### 3. 配置后端地址
项目默认使用 Vite 代理连接 `http://localhost:8080/api`（后端默认端口）。

如果需要修改，编辑 `vite.config.js`:
```js
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',  // 修改为你的后端地址
      changeOrigin: true
    }
  }
}
```

### 4. 启动开发服务器
```bash
npm run dev
```
访问 `http://localhost:3000`

### 5. 构建生产版本
```bash
npm run build
```

### 6. 预览构建结果
```bash
npm run preview
```

## 📱 页面路由与权限

### 完整页面路由

| 路径 | 说明 | 权限 | 页面组件 |
|------|------|------|----------|
| `/` | 首页套图列表 + 搜索 | **公开** | Home.vue |
| `/photoset/:id` | 套图详情 | **公开**（付费部分需购买） | PhotosetDetail.vue |
| `/create` | 创建套图 | **创作者/管理员** | CreatePhotoset.vue |
| `/edit/:id` | 编辑套图 | **创作者(本人)/管理员** | EditPhotoset.vue |
| `/login` | 用户登录 | **游客** | Login.vue |
| `/register` | 用户注册 | **游客** | Register.vue |
| `/profile` | 个人中心 | **已登录用户** | Profile.vue |
| `/favorites` | 我的收藏 | **已登录用户** | Favorites.vue |
| `/orders` | 我的订单（含退款） | **已登录用户** | Orders.vue |
| `/*` | 404 页面 | **公开** | NotFound.vue |

### 用户角色与权限（三级权限系统）

| 角色 | 说明 | 创建套图 | 编辑/删除 | 退款权限 | 管理后台 |
|------|------|----------|-----------|----------|----------|
| **用户 (User)** | 普通注册用户 | ❌ | ❌ | 48h自助退款 | ❌ |
| **创作者 (Creator)** | 内容创作者 | ✅ | 仅自己的 | 48h自助退款 | ❌ |
| **管理员 (Admin)** | 系统管理员 | ✅ | 任意套图 | 无限期退款 | ✅ |

## 🔧 核心功能详解

### 1. **智能搜索系统** 🔍
- **FULLTEXT 中文全文搜索**：支持中文分词检索
- **实时搜索结果**：边输入边展示结果
- **关键词高亮显示**：搜索结果中关键词高亮展示（黄色背景 + 红色文本）
- **搜索历史记忆**：自动记录最近的搜索历史

### 2. **套图管理功能**
- **Cloudflare R2 存储**：海外对象存储，避免国内存储限制
- **图片批量上传**：支持拖拽、多文件上传
- **套图编辑/删除**：创作者可以管理自己的套图
- **付费内容保护**：未购买用户只能查看封面和缩略图

### 3. **订单与支付**
- **支付宝网页支付**：集成支付接口
- **订单状态管理**：生成、支付、完成、退款全流程
- **48小时自助退款**：用户可在48小时内申请退款
- **管理员无限退款**：管理员可以随时处理订单退款

### 4. **个人中心功能**
- **套图收藏管理**：收藏、取消收藏
- **订单历史记录**：查看所有购买历史
- **创作者面板**：创作者专属的套图管理界面
- **用户偏好设置**：界面主题、通知设置等

### 5. **性能与安全增强 (Phase 5)**
- **Redis 缓存系统**：列表/详情数据缓存，静默降级机制
- **URL 签名防盗链**：HMAC-SHA256 签名系统，防止未授权访问
- **前端性能优化**：路由懒加载、组件缓存、图片懒加载
- **搜索关键词高亮**：前端正则匹配，实时高亮显示

## 🔌 完整的 API 接口对接

### ✅ **全部已对接的接口**

#### 用户认证模块
| 接口 | 方法 | 权限 | 说明 | 封装位置 |
|------|------|------|------|----------|
| `/api/auth/register` | POST | 游客 | 用户注册 | `api.register()` |
| `/api/auth/login` | POST | 游客 | 用户登录 | `api.login()` |
| `/api/auth/me` | GET | 登录用户 | 获取当前用户 | `api.getMe()` |
| `/api/auth/logout` | POST | 登录用户 | 用户退出 | `api.logout()` |

#### 套图模块
| 接口 | 方法 | 权限 | 说明 | 封装位置 |
|------|------|------|------|----------|
| `/api/photosets` | GET | 公开 | 套图列表（支持搜索） | `api.getPhotosetList()` |
| `/api/photosets/:id` | GET | 公开 | 套图详情 | `api.getPhotosetDetail()` |
| `/api/photosets` | POST | 创作者+ | 创建套图 | `api.createPhotoset()` |
| **新** `/api/photosets/:id` | PUT | 创作者(本人)/管理员 | 编辑套图 | `api.updatePhotoset()` |
| **新** `/api/photosets/:id` | DELETE | 创作者(本人)/管理员 | 删除套图 | `api.deletePhotoset()` |
| `/api/photosets/:id/favorite` | POST | 登录用户 | 收藏套图 | `api.favoritePhotoset()` |
| `/api/photosets/:id/favorite` | DELETE | 登录用户 | 取消收藏 | `api.unfavoritePhotoset()` |
| `/api/photosets/my` | GET | 创作者+ | 我的套图 | `api.getMyPhotosets()` |

#### 订单与支付模块
| 接口 | 方法 | 权限 | 说明 | 封装位置 |
|------|------|------|------|----------|
| `/api/orders` | POST | 登录用户 | 创建订单 | `api.createOrder()` |
| `/api/orders` | GET | 登录用户 | 获取我的订单 | `api.getMyOrders()` |
| `/api/orders/:id` | GET | 登录用户 | 订单详情 | `api.getOrderDetail()` |
| **新** `/api/orders/:id/refund` | POST | 购买者(48h内)/管理员 | 申请退款 | `api.refundOrder()` |

#### 标签与分类模块
| 接口 | 方法 | 权限 | 说明 | 封装位置 |
|------|------|------|------|----------|
| `/api/tags` | GET | 公开 | 获取标签列表 | `api.getTags()` |
| `/api/tags/hot` | GET | 公开 | 热门标签 | `api.getHotTags()` |
| `/api/categories` | GET | 公开 | 分类列表 | `api.getCategories()` |
| **新** `AdvancedSearch.vue` | — | 公开 | 动态分类搜索组件 | — |

#### 个人中心功能
| 接口 | 方法 | 权限 | 说明 | 封装位置 |
|------|------|------|------|----------|
| `/api/user/favorites` | GET | 登录用户 | 获取收藏列表 | `api.getFavorites()` |
| `/api/user/profile` | PUT | 登录用户 | 更新用户资料 | `api.updateProfile()` |

#### 图片上传
| 接口 | 方法 | 权限 | 说明 | 封装位置 |
|------|------|------|------|----------|
| `/api/upload` | POST | 创作者+ | 图片上传 | `api.upload()` |

#### 配置说明
- **接口位置**: `src/api/index.js` - 完整封装所有 API 调用
- **鉴权机制**: 自动携带 JWT Token
- **错误处理**: 统一拦截器处理错误响应
- **代理配置**: Vite 代理到 `http://localhost:8080/api`

## 🎯 开发进度与版本历史

### Phase 完成情况
- ✅ **Phase 1** - 基础架构搭建（完成）
- ✅ **Phase 2** - 用户端核心功能（收藏、搜索、个人中心）
- ✅ **Phase 3** - 支付与订单系统（支付宝集成）
- ✅ **Phase 3.5** - Cloudflare R2 存储集成（海外存储）
- ✅ **Phase 4** - 用户端功能补齐（编辑、退款接口）
- ✅ **Phase 5** - 高级功能增强（全部完成）

### Phase 5 具体实现
1. **套图删除功能** - 创作者删除自己的套图
2. **Redis 缓存系统** - 接口数据缓存，静默降级
3. **FULLTEXT 搜索** - 中文全文搜索 + 关键词高亮
4. **URL 签名防盗链** - HMAC-SHA256 签名安全防护
5. **订单退款系统** - 48h自助退款 + 管理员无限退款

### 关键词高亮功能说明
```javascript
// PhotosetCard.vue 中实现的关键词高亮
function highlightText(text, keyword) {
  const regex = new RegExp(`(${escapeRegExp(keyword)})`, 'gi')
  return text.replace(regex, '<span class="highlight">$1</span>')
}
```
- **搜索时**: Home.vue 传递 `keyword` 到 PhotosetCard
- **未搜索时**: 正常显示标题和描述
- **高亮样式**: 黄色背景 + 红色文字，增强可读性

## 🛠️ 开发注意事项

### 核心开发规范
1. **Token 管理**: 登录后 Token 自动存储在 localStorage，Axios 拦截器自动携带
2. **权限控制**: 前端路由守卫 + 后端 API 双重验证 + Pinia 状态管理
3. **响应处理**: 统一拦截器处理 HTTP 状态码、错误提示和加载状态
4. **组件设计**: 使用 `<script setup>` 语法，Props 类型定义，组件解耦

### 性能优化要点
1. **路由懒加载**: 页面组件按需加载，减少初始包体积
2. **图片优化**: 懒加载 + 缩略图，提升页面加载速度
3. **API 缓存**: 使用 Redis 缓存系统，减少数据库查询压力
4. **状态管理**: Pinia 状态缓存，避免重复网络请求

### 安全注意事项
1. **文件上传**: 仅支持图片格式，Cloudflare R2 存储，无国内服务器
2. **URL 签名**: 图片资源使用 HMAC-SHA256 签名，防止盗链
3. **输入验证**: 前后端双重输入验证，防止 XSS/CSRF 攻击
4. **权限验证**: 角色权限检查贯穿整个应用流程

## 🌐 浏览器兼容性

- Chrome 85+
- Firefox 80+
- Safari 14+
- Edge 85+
- Mobile Chrome / Safari 适配完成

## 📄 许可证

私有项目，仅供内部使用。

---

**最后更新**: 2026-04-09  
**项目状态**: ✅ **生产就绪 - Phase 5 全部完成**  
**维护者**: PhotoSet 开发团队
