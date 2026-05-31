# 管理后台 API 文档

> 基础路径：`/api/admin`
> 认证方式：Bearer Token + Admin 角色

---

## 一、用户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/users?page=1&page_size=20&role=all&status=all` | 用户列表 |
| GET | `/api/admin/users/export` | 导出用户（CSV）|
| GET | `/api/admin/users/:id` | 用户详情 |
| PUT | `/api/admin/users/:id/ban` | 封禁/解封用户 |
| PUT | `/api/admin/users/:id/role` | 修改用户角色 `{"role":"creator"}` |
| PUT | `/api/admin/users/:id/password` | 重置用户密码 |

---

## 二、套图审核

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/photosets?status=all&page=1&page_size=20` | 套图列表 |
| POST | `/api/admin/photosets/:id/approve` | 审核通过 |
| POST | `/api/admin/photosets/:id/reject` | 审核拒绝 |
| POST | `/api/admin/photosets/batch/approve` | 批量通过 `{"ids":[1,2,3]}` |
| POST | `/api/admin/photosets/batch/reject` | 批量拒绝 |
| POST | `/api/admin/photosets/batch/delete` | 批量删除 |

---

## 三、订单管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/orders?page=1&page_size=20&status=all` | 订单列表 |
| GET | `/api/admin/orders/export` | 导出订单（CSV）|
| POST | `/api/admin/orders/:id/refund` | 管理员强制退款 |

---

## 四、标签管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/tags?page=1&page_size=20` | 标签列表 |
| POST | `/api/admin/tags` | 创建标签 `{"name":"标签名"}` |
| PUT | `/api/admin/tags/:id` | 更新标签 |
| DELETE | `/api/admin/tags/:id` | 删除标签 |

---

## 五、分类管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/categories` | 分类列表 |
| POST | `/api/admin/categories` | 创建分类 |
| PUT | `/api/admin/categories/:id` | 更新分类 |
| DELETE | `/api/admin/categories/:id` | 删除分类 |

---

## 六、站点设置

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/settings` | 获取站点设置 |
| PUT | `/api/admin/settings` | 更新站点设置 |

### 设置字段说明
```json
{
  "site_name": "PhotoSet",
  "site_description": "摄影套图浏览平台",
  "site_logo": "https://...",
  "seo_keywords": "摄影,套图,人像",
  "about_me": "关于我们...",
  "contact_email": "admin@example.com",
  "icp": "京ICP备12345678号",
  "footer_text": "Copyright 2026"
}
```

---

## 七、系统管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/admin/system/restart` | 重启后端服务 |
| POST | `/api/admin/mail/test-connection` | 测试邮件连接 |
| GET | `/api/admin/mail/config` | 获取邮件配置 |
| POST | `/api/admin/mail/send-test` | 发送测试邮件 |
| GET | `/api/admin/watermark/info` | 获取水印配置 |
| POST | `/api/admin/storage/test` | 测试存储连接 |
| GET | `/api/admin/storage/status` | 存储状态 |

---

## 八、页面管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/pages?page=1&page_size=20` | 页面列表 |
| POST | `/api/admin/pages` | 创建页面 |
| GET | `/api/admin/pages/:id` | 页面详情 |
| PUT | `/api/admin/pages/:id` | 更新页面 |
| DELETE | `/api/admin/pages/:id` | 删除页面 |

### 页面字段
```json
{
  "title": "关于我们",
  "slug": "about",
  "content": "<p>内容</p>",
  "is_published": true,
  "sort": 1
}
```

---

## 九、会员套餐管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/memberships` | 套餐列表 |
| POST | `/api/admin/memberships` | 创建套餐 |
| PUT | `/api/admin/memberships/:id` | 更新套餐 |
| DELETE | `/api/admin/memberships/:id` | 删除套餐 |

---

## 十、数据统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/stats` | 统计面板（用户数/套图数/订单数/收入）|
| GET | `/api/admin/stats/trend?days=7` | 趋势数据 |
| GET | `/api/admin/logs?page=1&page_size=20` | 管理员操作日志 |

---

## 十一、开发者中心

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/dev/api-keys` | API Key 列表 |
| POST | `/api/admin/dev/api-keys` | 创建 API Key |
| DELETE | `/api/admin/dev/api-keys/:id` | 删除 API Key |
| GET | `/api/admin/dev/api-docs` | API 文档说明 |
| GET | `/api/admin/dev/sign-url-docs` | 签名 URL 文档 |

---

## 十二、社区管理

详见 [社区 API 文档](community-api.md) 的 **三、管理后台路由** 部分。

---

*文档更新时间：2026-05-31*
