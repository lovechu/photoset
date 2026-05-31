# 消息、通知与关注 API 文档

> 基础路径：`/api/community/messages` `/api/community/notifications` `/api/community/users/:id/follow`
> 认证方式：Bearer Token（需登录）

---

## 一、关注系统

### 1.1 关注用户
```
POST /api/community/users/:id/follow
Authorization: Bearer <token>
```

### 1.2 取消关注
```
DELETE /api/community/users/:id/follow
Authorization: Bearer <token>
```

### 1.3 检查是否已关注
```
GET /api/community/users/:id/follow/check
Authorization: Bearer <token>
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "is_following": true
  }
}
```

### 1.4 批量检查关注状态
```
POST /api/community/users/batch-follow-check
Authorization: Bearer <token>
Content-Type: application/json
```
**请求体：**
```json
{
  "user_ids": [1, 2, 3]
}
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "1": true,
    "2": false,
    "3": true
  }
}
```

### 1.5 关注列表
```
GET /api/community/users/:id/following?page=1&page_size=20
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "list": [{
      "id": 2,
      "name": "用户B",
      "avatar": "https://..."
    }],
    "total": 50
  }
}
```

### 1.6 粉丝列表
```
GET /api/community/users/:id/followers?page=1&page_size=20
```
**响应：** 同 1.5

---

## 二、系统通知

### 2.1 通知列表
```
GET /api/community/notifications?page=1&page_size=20&type=all
Authorization: Bearer <token>
```
| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码 |
| page_size | int | 每页条数 |
| type | string | `all` / `system` / `like` / `reply` / `follow` / `mention` |

**响应：**
```json
{
  "code": 0,
  "data": {
    "list": [{
      "id": 1,
      "type": "like",
      "title": "新点赞",
      "content": "用户A 赞了你的帖子",
      "sender": {"id": 2, "name": "用户A", "avatar": "https://..."},
      "target_type": "post",
      "target_id": 10,
      "is_read": false,
      "created_at": "2026-05-31T10:00:00Z"
    }],
    "total": 20,
    "unread_count": 5
  }
}
```

### 2.2 未读通知数量
```
GET /api/community/notifications/unread-count
Authorization: Bearer <token>
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "count": 5
  }
}
```

### 2.3 标记通知为已读
```
PUT /api/community/notifications/:id/read
Authorization: Bearer <token>
```

### 2.4 全部标记已读
```
PUT /api/community/notifications/read-all
Authorization: Bearer <token>
```

### 2.5 删除通知
```
DELETE /api/community/notifications/:id
Authorization: Bearer <token>
```

### 2.6 清空所有通知
```
DELETE /api/community/notifications
Authorization: Bearer <token>
```

---

## 三、私信系统

### 3.1 会话列表
```
GET /api/community/messages/conversations?page=1&page_size=20
Authorization: Bearer <token>
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "list": [{
      "user_id": 2,
      "user_name": "用户B",
      "user_avatar": "https://...",
      "last_message": "你好！",
      "last_message_time": "2026-05-31T10:00:00Z",
      "unread_count": 3
    }],
    "total": 10
  }
}
```

### 3.2 获取与某用户的聊天记录
```
GET /api/community/messages/conversations/:user_id?page=1&page_size=20
Authorization: Bearer <token>
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "list": [{
      "id": 1,
      "from_user_id": 1,
      "to_user_id": 2,
      "content": "你好！",
      "is_read": true,
      "created_at": "2026-05-31T10:00:00Z"
    }],
    "total": 50
  }
}
```

### 3.3 发送私信
```
POST /api/community/messages
Authorization: Bearer <token>
Content-Type: application/json
```
**请求体：**
```json
{
  "to_user_id": 2,
  "content": "你好！"
}
```

### 3.4 私信未读数
```
GET /api/community/messages/unread-count
Authorization: Bearer <token>
```

### 3.5 标记单条消息已读
```
PUT /api/community/messages/:id/read
Authorization: Bearer <token>
```

### 3.6 标记整个会话已读
```
PUT /api/community/messages/conversations/:user_id/read
Authorization: Bearer <token>
```

### 3.7 删除消息
```
DELETE /api/community/messages/:id
Authorization: Bearer <token>
```

---

*文档更新时间：2026-05-31*
