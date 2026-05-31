# 站点设置与页面 API 文档

> 基础路径：`/api/settings` `/api/pages`
> 公开接口，无需登录

---

## 一、公开站点设置

### 1.1 获取站点设置
```
GET /api/settings
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "site_name": "PhotoSet",
    "site_description": "摄影套图浏览平台",
    "site_logo": "https://...",
    "seo_keywords": "摄影,套图,人像",
    "about_me": "关于我们...",
    "contact_email": "admin@example.com",
    "icp": "京ICP备12345678号",
    "footer_text": "Copyright 2026",
    "watermark_enabled": true,
    "watermark_text": "PhotoSet"
  }
}
```

---

## 二、公开页面

### 2.1 页面列表
```
GET /api/pages
```
**响应：**
```json
{
  "code": 0,
  "data": [{
    "id": 1,
    "title": "关于我们",
    "slug": "about",
    "is_published": true,
    "sort": 1
  }]
}
```

### 2.2 页面详情
```
GET /api/pages/:slug
```
**响应：**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "title": "关于我们",
    "slug": "about",
    "content": "<p>这里是关于我们的内容...</p>",
    "is_published": true,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-05-01T00:00:00Z"
  }
}
```

---

*文档更新时间：2026-05-31*
