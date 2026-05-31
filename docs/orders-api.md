# 订单与会员 API 文档

> 基础路径：`/api/memberships` `/api/orders`

---

## 一、会员套餐（公开）

### 1.1 会员套餐列表
```
GET /api/memberships
```
**响应：**
```json
{
  "code": 0,
  "data": [{
    "id": 1,
    "name": "月度会员",
    "description": "享受30天会员权益",
    "price": 29.99,
    "original_price": 39.99,
    "duration_days": 30,
    "benefits": ["无限下载", "免广告", "专属客服"],
    "is_active": true
  }]
}
```

---

## 二、订单（需登录）

### 2.1 我的订单列表
```
GET /api/orders?page=1&page_size=20&status=all
Authorization: Bearer <token>
```
| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码 |
| page_size | int | 每页条数 |
| status | string | `all` / `pending` / `paid` / `refunded` |

**响应：**
```json
{
  "code": 0,
  "data": {
    "list": [{
      "id": "ORD202605310001",
      "type": "photoset",
      "target_id": 1,
      "target_title": "套图标题",
      "amount": 9.99,
      "status": "paid",
      "payment_method": "alipay",
      "paid_at": "2026-05-31T10:00:00Z",
      "created_at": "2026-05-31T09:00:00Z"
    }],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

### 2.2 创建订单
```
POST /api/orders
Authorization: Bearer <token>
Content-Type: application/json
```
**请求体：**
```json
{
  "type": "photoset",
  "target_id": 1
}
```
**type 可选值：**
- `photoset`：购买套图
- `membership`：购买会员

**响应：**
```json
{
  "code": 0,
  "data": {
    "id": "ORD202605310001",
    "amount": 9.99,
    "status": "pending",
    "payment_url": "https://..."
  }
}
```

### 2.3 支付订单
```
POST /api/orders/:id/pay
Authorization: Bearer <token>
Content-Type: application/json
```
**请求体：**
```json
{
  "payment_method": "alipay"
}
```
**payment_method 可选值：**
- `alipay`：支付宝
- `wechat`：微信支付

**响应：**
```json
{
  "code": 0,
  "data": {
    "order_id": "ORD202605310001",
    "status": "paid",
    "paid_at": "2026-05-31T10:00:00Z"
  }
}
```

### 2.4 申请退款（用户48小时内）
```
POST /api/orders/:id/refund
Authorization: Bearer <token>
Content-Type: application/json
```
**请求体：**
```json
{
  "reason": "重复购买"
}
```
**响应：**
```json
{
  "code": 0,
  "message": "退款申请已提交"
}
```

---

*文档更新时间：2026-05-31*
