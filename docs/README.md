# PhotoSet API 文档索引

> 本文档目录包含 PhotoSet 后端所有 API 接口的调用说明。
> 基础路径：`https://tt.cy.mk/api`
> 认证方式：Bearer Token（`Authorization: Bearer <token>`）

---

## 文档列表

| 文档 | 说明 | 路径前缀 |
|------|------|----------|
| [认证 API](auth-api.md) | 注册/登录/验证码/找回密码/修改密码 | `/api/auth/*` |
| [套图 API](photosets-api.md) | 套图浏览/搜索/创建/编辑/删除/评论/收藏/分类/标签 | `/api/photosets/*` `/api/favorites/*` `/api/tags` `/api/categories` |
| [订单 API](orders-api.md) | 会员套餐/订单创建/支付/退款 | `/api/memberships` `/api/orders/*` |
| [社区 API](community-api.md) | 帖子/回帖/点赞/举报/草稿/分享/话题/标签/热门/搜索/推荐 | `/api/community/*` |
| [消息与通知 API](messages-api.md) | 私信/系统通知/关注/粉丝 | `/api/community/messages/*` `/api/community/notifications/*` `/api/community/users/*/follow` |
| [用户等级 API](user-api.md) | 用户等级/成就/积分商城/积分排行榜 | `/api/community/user/*` `/api/community/levels` `/api/community/points-mall/*` |
| [管理后台 API](admin-api.md) | 用户管理/套图审核/订单管理/标签管理/分类管理/站点设置/页面管理/社区管理 | `/api/admin/*` |
| [站点设置 API](settings-api.md) | 公开站点设置/公开页面 | `/api/settings` `/api/pages/*` |

---

## 通用响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | `0`=成功，其他=错误码 |
| message | string | 提示信息 |
| data | any | 响应数据 |

## 通用错误码

| code | 说明 |
|------|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 401 | 未登录 / Token 失效 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

---

*文档更新时间：2026-05-31*
