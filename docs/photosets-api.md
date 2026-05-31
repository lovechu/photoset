# 套图 API 文档

> 基础路径：`/api/photosets` `/api/favorites` `/api/tags` `/api/categories`

---

## 一、套图浏览（公开）

### 1.1 套图列表
```
GET /api/photosets?page=1&page_size=20&sort=latest&tag=tagname&category=categoryname
```
| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码，默认1 |
| page_size | int | 每页条数，默认20 |
| sort | string | `latest` / `popular` / `price_asc` / `price_desc` |
| tag | string | 标签过滤（可选）|
| category | string | 分类过滤（可选）|
| query | string | 关键词搜索（可选，FULLTEXT）|

**响应：**
```json
{
  "code": 0,
  "data": {
    "list": [{
      "id": 1,
      "title": "套图标题",
      "cover": "https://...",
      "price": 9.99,
      "original_price": 19.99,
      "tag_names": ["日系", "人像"],
      "category_name": "人像摄影",
      "author_name": "摄影师A",
      "photo_count": 50,
      "is_free": false,
      "created_at": "2026-05-01T00:00:00Z"
    }],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 1.2 高级搜索
```
GET /api/photosets/advanced?query=关键词&tag=标签&min_price=0&max_price=100&sort=latest&page=1
```
| 参数 | 类型 | 说明 |
|------|------|------|
| query | string | 关键词（标题/描述搜索）|
| tag | string | 标签过滤 |
| min_price | float | 最低价格 |
| max_price | float | 最高价格 |
| sort | string | 排序方式 |
| page | int | 页码 |
| page_size | int | 每页条数 |

### 1.3 套图详情
```
GET /api/photosets/:id
Authorization: Bearer <token>（可选，登录用户可显示购买状态）
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "title": "套图标题",
    "description": "套图描述...",
    "cover": "https://...",
    "price": 9.99,
    "original_price": 19.99,
    "tag_names": ["日系", "人像"],
    "category_name": "人像摄影",
    "author_name": "摄影师A",
    "author_id": 123,
    "photo_count": 50,
    "is_free": false,
    "is_purchased": false,
    "photos": [{
      "id": 1,
      "url": "https://...",
      "thumb_url": "https://...",
      "width": 1920,
      "height": 1080
    }],
    "created_at": "2026-05-01T00:00:00Z"
  }
}
```

### 1.4 下载套图（已购买用户）
```
GET /api/photosets/:id/download
Authorization: Bearer <token>
```
**响应：** 文件流下载

---

## 二、套图管理（创作者/管理员）

### 2.1 创建套图
```
POST /api/photosets
Authorization: Bearer <token>
Content-Type: application/json
```
**请求体：**
```json
{
  "title": "套图标题",
  "description": "套图描述",
  "price": 9.99,
  "original_price": 19.99,
  "tag_names": ["日系", "人像"],
  "category_name": "人像摄影",
  "cover": "https://...",
  "is_free": false,
  "photos": [
    {"url": "https://...", "thumb_url": "https://...", "width": 1920, "height": 1080}
  ]
}
```

### 2.2 更新套图
```
PUT /api/photosets/:id
Authorization: Bearer <token>
Content-Type: application/json
```
**请求体：** 同 2.1

### 2.3 删除套图
```
DELETE /api/photosets/:id
Authorization: Bearer <token>
```

---

## 三、评论（公开/登录）

### 3.1 评论列表
```
GET /api/photosets/:id/comments?page=1&page_size=20
Authorization: Bearer <token>（可选）
```

### 3.2 发表评论
```
POST /api/photosets/:id/comments
Authorization: Bearer <token>
Content-Type: application/json
```
**请求体：**
```json
{
  "content": "评论内容",
  "parent_id": 0
}
```

### 3.3 评论回复列表
```
GET /api/photosets/:id/comments/:commentId/replies?page=1&page_size=20
```

### 3.4 删除评论
```
DELETE /api/photosets/:id/comments/:commentId
Authorization: Bearer <token>
```

### 3.5 点赞/取消点赞评论
```
POST /api/photosets/:id/comments/:commentId/like
Authorization: Bearer <token>
```

---

## 四、收藏（需登录）

### 4.1 添加收藏
```
POST /api/favorites/:photosetId
Authorization: Bearer <token>
```

### 4.2 取消收藏
```
DELETE /api/favorites/:photosetId
Authorization: Bearer <token>
```

### 4.3 我的收藏列表
```
GET /api/favorites?page=1&page_size=20
Authorization: Bearer <token>
```
**响应：** 同套图列表格式

---

## 五、标签（公开）

### 5.1 标签列表
```
GET /api/tags
```
**响应：**
```json
{
  "code": 0,
  "data": [
    {"id": 1, "name": "日系", "count": 100},
    {"id": 2, "name": "人像", "count": 200}
  ]
}
```

---

## 六、分类（公开）

### 6.1 分类列表
```
GET /api/categories
```
**响应：**
```json
{
  "code": 0,
  "data": [
    {"id": 1, "name": "人像摄影", "sort": 1},
    {"id": 2, "name": "风景摄影", "sort": 2}
  ]
}
```

---

*文档更新时间：2026-05-31*
