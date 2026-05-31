# 认证 API 文档

> 基础路径：`/api/auth`
> 所有认证相关接口

---

## 一、验证码

### 1.1 获取验证码
```
GET /api/auth/captcha
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "captcha_id": "abc123",
    "captcha_image": "data:image/png;base64,iVBORw0..."
  }
}
```

---

## 二、注册

### 2.1 用户注册
```
POST /api/auth/register
Content-Type: application/json
```
**请求体：**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "用户名",
  "captcha_id": "abc123",
  "captcha_code": "1234"
}
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "用户名",
    "role": "user",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

---

## 三、登录

### 3.1 用户登录
```
POST /api/auth/login
Content-Type: application/json
```
**请求体：**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "captcha_id": "abc123",
  "captcha_code": "1234"
}
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "用户名",
    "avatar": "https://...",
    "role": "user",
    "status": "active",
    "created_at": "2026-01-01T00:00:00Z",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 3.2 获取当前用户信息
```
GET /api/auth/me
Authorization: Bearer <token>
```
**响应：** 同 3.1

---

## 四、密码管理

### 4.1 修改密码
```
PUT /api/auth/password
Authorization: Bearer <token>
Content-Type: application/json
```
**请求体：**
```json
{
  "old_password": "oldpass123",
  "new_password": "newpass123"
}
```

### 4.2 忘记密码（发送重置邮件）
```
POST /api/auth/forgot-password
Content-Type: application/json
```
**请求体：**
```json
{
  "email": "user@example.com"
}
```

### 4.3 通过 Token 重置密码
```
POST /api/auth/reset-password
Content-Type: application/json
```
**请求体：**
```json
{
  "token": "reset-token-string",
  "new_password": "newpass123"
}
```

---

## 五、个人资料

### 5.1 更新个人资料
```
PUT /api/auth/profile
Authorization: Bearer <token>
Content-Type: application/json
```
**请求体：**
```json
{
  "name": "新用户名",
  "avatar": "https://...",
  "bio": "个人简介"
}
```

---

## 六、邮件配置检查

### 6.1 检查邮件服务配置
```
GET /api/auth/email-config
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "enabled": true,
    "smtp_host": "smtp.example.com"
  }
}
```

---

*文档更新时间：2026-05-31*
