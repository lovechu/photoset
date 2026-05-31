# 用户等级与积分 API 文档

> 基础路径：`/api/community/levels` `/api/community/user/*` `/api/community/points-mall/*` `/api/community/points/*`

---

## 一、等级体系（公开）

### 1.1 全部等级配置
```
GET /api/community/levels
```
**响应：**
```json
{
  "code": 0,
  "data": [{
    "level": 1,
    "name": "新手",
    "min_points": 0,
    "max_points": 100,
    "icon": "https://...",
    "privileges": ["发帖", "评论"]
  }, {
    "level": 2,
    "name": "活跃者",
    "min_points": 100,
    "max_points": 500,
    "icon": "https://...",
    "privileges": ["发帖", "评论", "私信"]
  }]
}
```

### 1.2 查看指定用户等级
```
GET /api/community/users/:id/level
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "user_id": 1,
    "current_level": 3,
    "level_name": "达人",
    "current_points": 350,
    "next_level_points": 500,
    "progress_percent": 62.5
  }
}
```

### 1.3 查看指定用户成就
```
GET /api/community/users/:id/achievements
```
**响应：**
```json
{
  "code": 0,
  "data": [{
    "id": 1,
    "name": "初次发帖",
    "description": "发布第一篇文章",
    "icon": "https://...",
    "points_reward": 10,
    "achieved_at": "2026-05-01T00:00:00Z"
  }]
}
```

---

## 二、当前用户等级（需登录）

### 2.1 我的等级信息
```
GET /api/community/user/level
Authorization: Bearer <token>
```
**响应：** 同 1.2

### 2.2 我的成就
```
GET /api/community/user/achievements
Authorization: Bearer <token>
```
**响应：** 同 1.3

---

## 三、积分商城（公开/登录）

### 3.1 商品分类
```
GET /api/community/points-mall/categories
```
**响应：**
```json
{
  "code": 0,
  "data": [
    {"id": 1, "name": "虚拟道具", "sort": 1},
    {"id": 2, "name": "会员权益", "sort": 2}
  ]
}
```

### 3.2 商品列表
```
GET /api/community/points-mall/items?page=1&page_size=20&category_id=1
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "list": [{
      "id": 1,
      "name": "改名卡",
      "description": "修改用户名一次",
      "icon": "https://...",
      "points_price": 500,
      "stock": 100,
      "category_id": 1
    }],
    "total": 20
  }
}
```

### 3.3 兑换商品（需登录）
```
POST /api/community/points-mall/exchange
Authorization: Bearer <token>
Content-Type: application/json
```
**请求体：**
```json
{
  "item_id": 1
}
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "exchange_id": 100,
    "item_name": "改名卡",
    "points_spent": 500,
    "remaining_points": 150,
    "status": "success"
  }
}
```

### 3.4 兑换记录（需登录）
```
GET /api/community/points-mall/history?page=1&page_size=20
Authorization: Bearer <token>
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "list": [{
      "id": 100,
      "item_name": "改名卡",
      "points_spent": 500,
      "status": "completed",
      "created_at": "2026-05-31T10:00:00Z"
    }],
    "total": 5
  }
}
```

---

## 四、积分排行榜（需登录）

### 4.1 排行榜
```
GET /api/community/points/leaderboard?type=week&page=1&page_size=20
Authorization: Bearer <token>
```
| 参数 | 类型 | 说明 |
|------|------|------|
| type | string | `total` / `week` / `month` |
| page | int | 页码 |
| page_size | int | 每页条数 |

**响应：**
```json
{
  "code": 0,
  "data": {
    "list": [{
      "rank": 1,
      "user_id": 1,
      "user_name": "用户A",
      "avatar": "https://...",
      "points": 5000
    }],
    "my_rank": 15,
    "my_points": 350
  }
}
```

---

*文档更新时间：2026-05-31*
