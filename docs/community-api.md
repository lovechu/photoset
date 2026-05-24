# 社区功能 API 文档

> 基础路径：`/api/community`（公开/登录用户），`/api/admin/community`（管理员）
> 认证方式：Bearer Token（登录态请求传 `Authorization` 头）

---

## 一、公开路由（无需登录）

### 1.1 帖子列表
```
GET /api/community/posts?page=1&page_size=20&category=discuss&sort=latest
```
| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码，默认1 |
| page_size | int | 每页条数，默认20 |
| category | string | 分类过滤（可选）|
| sort | string | `latest` / `hot`，默认 latest |
| keyword | string | 关键词搜索（可选）|

**响应：**
```json
{
  "code": 0,
  "data": {
    "list": [{
      "id": 1,
      "title": "帖子标题",
      "content": "内容...",
      "category": "discuss",
      "author_id": 123,
      "author_name": "张三",
      "reply_count": 5,
      "like_count": 10,
      "is_pinned": false,
      "is_essence": false,
      "created_at": "2026-05-24T10:00:00Z"
    }],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 1.2 帖子详情
```
GET /api/community/posts/:id
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "title": "帖子标题",
    "content": "内容...",
    "category": "discuss",
    "author_id": 123,
    "author_name": "张三",
    "reply_count": 5,
    "like_count": 10,
    "is_pinned": false,
    "is_essence": false,
    "is_liked": false,
    "created_at": "2026-05-24T10:00:00Z"
  }
}
```

---

### 1.3 帖子回帖列表
```
GET /api/community/posts/:id/replies?page=1&page_size=20&sort=latest
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "list": [{
      "id": 1,
      "post_id": 1,
      "content": "回复内容",
      "author_id": 124,
      "author_name": "李四",
      "like_count": 3,
      "is_liked": false,
      "parent_id": 0,
      "children": [],
      "created_at": "2026-05-24T10:05:00Z"
    }],
    "total": 5,
    "page": 1,
    "page_size": 20
  }
}
```

---

### 1.4 分类列表
```
GET /api/community/categories
```
**响应：**
```json
{
  "code": 0,
  "data": [
    { "key": "discuss", "label": "讨论" },
    { "key": "suggestion", "label": "建议" },
    { "key": "showcase", "label": "作品展示" },
    { "key": "question", "label": "问答" }
  ]
}
```

---

### 1.5 热门帖子
```
GET /api/community/hot?limit=10
```
**响应：** 同 1.1，按 `like_count + reply_count` 排序

---

## 二、登录用户路由（需 Bearer Token）

### 2.1 发帖
```
POST /api/community/posts
Content-Type: application/json
```
**请求体：**
```json
{
  "title": "帖子标题",
  "content": "帖子内容",
  "category": "discuss"
}
```
**响应：**
```json
{ "code": 0, "data": { "id": 1 }, "message": "发布成功" }
```

---

### 2.2 回帖
```
POST /api/community/posts/:id/reply
Content-Type: application/json
```
**请求体：**
```json
{
  "content": "回复内容",
  "parent_id": 0
}
```
> `parent_id=0` 为一级回帖；非0为楼中楼

**响应：**
```json
{ "code": 0, "data": { "id": 10 }, "message": "回复成功" }
```

---

### 2.3 点赞/取消点赞帖子
```
POST /api/community/posts/:id/like
```
> 幂等操作：已点赞则取消，未点赞则添加

**响应：**
```json
{ "code": 0, "data": { "is_liked": true, "like_count": 11 }, "message": "操作成功" }
```

---

### 2.4 点赞/取消点赞回帖
```
POST /api/community/replies/:id/like
```
**响应：** 同 2.3

---

### 2.5 举报帖子
```
POST /api/community/posts/:id/report
Content-Type: application/json
```
**请求体：**
```json
{ "reason": "垃圾信息" }
```

---

### 2.6 举报回帖
```
POST /api/community/replies/:id/report
Content-Type: application/json
```
**请求体：** 同 2.5

---

### 2.7 我的帖子
```
GET /api/community/me/posts?page=1&page_size=20
```

---

### 2.8 我的回帖
```
GET /api/community/me/replies?page=1&page_size=20
```

---

### 2.9 我的积分
```
GET /api/community/me/points
```
**响应：**
```json
{ "code": 0, "data": { "points": 150, "level": "活跃用户" } }
```

---

## 三、管理后台路由（需 Bearer Token + Admin 角色）

### 3.1 帖子管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/community/posts?page=1&page_size=20&status=all` | 帖子列表（含过滤）|
| PUT | `/api/admin/community/posts/:id/pin` | 置顶/取消置顶 |
| PUT | `/api/admin/community/posts/:id/status` | 审核通过/拒绝 `{"status":"approved"}` |
| DELETE | `/api/admin/community/posts/:id` | 删除帖子（级联删回帖/点赞/举报）|

---

### 3.2 回帖管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/community/replies?page=1&page_size=20` | 回帖列表 |
| DELETE | `/api/admin/community/replies/:id` | 删除回帖 |

---

### 3.3 敏感词管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/community/keywords` | 敏感词列表 |
| POST | `/api/admin/community/keywords` | 添加敏感词 `{"word":"xxx","category":"politics"}` |
| PUT | `/api/admin/community/keywords/:id` | 修改敏感词 |
| DELETE | `/api/admin/community/keywords/:id` | 删除敏感词 |
| PUT | `/api/admin/community/keywords/reload` | 热加载敏感词到内存 |

---

### 3.4 举报处理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/community/reports?page=1&page_size=20&status=pending` | 举报列表 |
| PUT | `/api/admin/community/reports/:id/resolve` | 处理举报 `{"action":"delete_post"}` |

**action 可选值：**
- `delete_post`：删除被举报帖子
- `delete_reply`：删除被举报回帖
- `warn_user`：警告用户
- `dismiss`：忽略举报

---

### 3.5 用户积分管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/community/users?page=1&page_size=20` | 社区用户列表（含积分）|
| PUT | `/api/admin/community/users/:id/points` | 调整积分 `{"points":50,"reason":"优质内容奖励"}` |

---

### 3.6 统计面板

```
GET /api/admin/community/stats
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "total_posts": 1234,
    "total_replies": 5678,
    "total_users": 300,
    "reports_pending": 5,
    "posts_today": 12,
    "replies_today": 45
  }
}
```

---

## 四、错误码

| code | 说明 |
|------|------|
| 0 | 成功 |
| 401 | 未登录 / Token 失效 |
| 403 | 无权限（非管理员）|
| 404 | 资源不存在 |
| 400 | 参数错误 |
| 500 | 服务器内部错误 |

---

## 五、数据库表结构

### `posts`
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| title | varchar | 标题 |
| content | text | 内容（脱敏后存储）|
| category | varchar(20) | 分类 |
| author_id | bigint | 作者ID |
| reply_count | int | 回帖数 |
| like_count | int | 点赞数 |
| is_pinned | boolean | 是否置顶 |
| is_essence | boolean | 是否加精 |
| status | varchar(20) | pending/approved/rejected |
| created_at | datetime | 创建时间 |

### `post_replies`
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| post_id | bigint | 所属帖子 |
| author_id | bigint | 作者ID |
| content | text | 内容 |
| parent_id | bigint | 父回帖ID（楼中楼）|
| like_count | int | 点赞数 |
| created_at | datetime | 创建时间 |

### `post_likes` / `post_reply_likes`
| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | bigint | 用户ID |
| post_id / reply_id | bigint | 被点赞对象 |
| created_at | datetime | 点赞时间 |

### `sensitive_words`
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| word | varchar(100) | 敏感词 |
| category | varchar(50) | 分类 |
| is_active | boolean | 是否启用 |
| created_at | datetime | 创建时间 |

### `post_reports`
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| post_id | bigint | 被举报帖子（可选）|
| reply_id | bigint | 被举报回帖（可选）|
| reporter_id | bigint | 举报人 |
| reason | varchar(255) | 举报原因 |
| status | varchar(20) | pending/resolved |
| created_at | datetime | 创建时间 |

### `user_points`
| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | bigint | 用户ID |
| points | int | 当前积分 |
| level | varchar(50) | 等级 |

---

*文档生成时间：2026-05-24*
