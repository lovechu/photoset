# PhotoSet 社区功能 — PRD & 开发计划

> 版本：v1.0
> 日期：2026-05-24
> 状态：待用户确认后开工

---

## 一、产品概述

在 PhotoSet 中新增**社区模块**，用户可发帖、回帖、点赞，内容先发后审。
前端暂不开发（Web 端风险高），**仅做后端 API**，供后续 App 集成。

---

## 二、核心设计

### 2.1 数据模型

#### `posts` — 帖子主表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | 发帖人 |
| title | varchar(200) | 标题 |
| content | text | 内容（含 [img:URL] 标记） |
| photoset_id | uint (nullable) | 关联套图（可选） |
| category | varchar(20) | 分类：discussion/qa/showcase/suggestion |
| visibility | varchar(20) | 可见范围：public/member/vip/admin |
| is_pinned | bool | 置顶 |
| view_count | int | 浏览数 |
| reply_count | int | 回帖数 |
| like_count | int | 点赞数 |
| status | varchar(20) | approved/pending/rejected（先发后审，默认 approved） |
| created_at / updated_at | datetime | |

#### `post_replies` — 回帖表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| post_id | uint | 所属帖子 |
| user_id | uint | 回帖人 |
| content | text | 内容 |
| parent_reply_id | uint (nullable) | 楼中楼（回复某条回复） |
| like_count | int | 点赞数 |
| created_at / updated_at | datetime | |

#### `post_likes` — 帖子点赞记录
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | |
| post_id | uint | |
| created_at | datetime | 联合唯一索引 (user_id, post_id) |

#### `post_reply_likes` — 回帖点赞记录
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | |
| reply_id | uint | |
| created_at | datetime | 联合唯一索引 (user_id, reply_id) |

#### `user_points` — 用户积分/等级
| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | uint | 主键 |
| points | int | 总积分 |
| level | int | 等级 1-10 |
| updated_at | datetime | |

#### `sensitive_words` — 敏感词库
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| word | varchar(100) | 敏感词（唯一索引） |
| replacement | varchar(100) | 替换词，默认 `***` |
| is_active | bool | 是否启用 |
| created_at | datetime | |

#### `post_reports` — 举报记录
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| post_id | uint (nullable) | 举报帖子 |
| reply_id | uint (nullable) | 举报回帖 |
| reporter_id | uint | 举报人 |
| reason | varchar(500) | 举报原因 |
| status | varchar(20) | pending/resolved/rejected |
| created_at | datetime | |

---

### 2.2 积分/等级规则

| 行为 | 积分 | 说明 |
|------|------|------|
| 发帖 | +10 | 每日上限 50 |
| 回帖 | +5 | 每日上限 30 |
| 帖子被点赞 | +2 | 无上限 |
| 回帖被点赞 | +1 | 无上限 |
| 帖子被删除（违规） | -20 | |

**等级对照表：**

| 等级 | 所需积分 | 称号 |
|------|----------|------|
| L1 | 0 | 新手上路 |
| L2 | 100 | 活跃会员 |
| L3 | 500 | 资深会员 |
| L4 | 2000 | 金牌会员 |
| L5 | 5000 | 钻石会员 |
| L6 | 10000 | 至尊会员 |
| L7-L10 | 每 +10000 | 荣耀会员 |

---

### 2.3 敏感词过滤逻辑

1. 发帖/回帖时，对 `title` + `content` 执行过滤
2. 从 `sensitive_words` 表加载所有 `is_active=1` 的词
3. 匹配到的词替换为 `replacement`（默认 `***`）
4. 过滤后的内容存入数据库
5. 管理后台可增删敏感词

---

### 2.4 权限模型

**发帖权限：** 登录即可发帖

**可见范围（visibility）：**
- `public`：所有人可见（包括未登录）
- `member`：会员及以上可见
- `vip`：VIP 及以上可见
- `admin`：仅管理员可见

**API 层过滤逻辑：**
- 未登录：只返回 `visibility=public` 的帖子
- 登录用户：根据自己角色过滤

---

## 三、API 接口清单

### 3.1 社区公开 API（需登录态，部分支持匿名）

| 方法 | 路径 | 说明 | 登录要求 |
|------|------|------|----------|
| GET | `/api/community/posts` | 帖子列表 | 可选 |
| GET | `/api/community/posts/:id` | 帖子详情 + 浏览+1 | 可选 |
| GET | `/api/community/posts/:id/replies` | 回帖列表（楼中楼树形） | 可选 |
| POST | `/api/community/posts` | 发帖 | **必需** |
| POST | `/api/community/posts/:id/reply` | 回帖 | **必需** |
| POST | `/api/community/posts/:id/like` | 点赞帖子（toggle） | **必需** |
| POST | `/api/community/replies/:id/like` | 点赞回帖（toggle） | **必需** |
| POST | `/api/community/posts/:id/report` | 举报帖子 | **必需** |
| POST | `/api/community/replies/:id/report` | 举报回帖 | **必需** |
| GET | `/api/community/categories` | 分类列表 | 可选 |
| GET | `/api/community/hot` | 热门帖子 | 可选 |

### 3.2 管理后台 API（admin 权限）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/community/posts` | 帖子列表（含 status 筛选） |
| PUT | `/api/admin/community/posts/:id/pin` | 置顶/取消置顶 |
| PUT | `/api/admin/community/posts/:id/status` | 修改状态（approve/reject） |
| DELETE | `/api/admin/community/posts/:id` | 删除帖子 |
| GET | `/api/admin/community/replies` | 回帖列表 |
| DELETE | `/api/admin/community/replies/:id` | 删除回帖 |
| GET | `/api/admin/community/keywords` | 敏感词列表 |
| POST | `/api/admin/community/keywords` | 添加敏感词 |
| PUT | `/api/admin/community/keywords/:id` | 修改敏感词 |
| DELETE | `/api/admin/community/keywords/:id` | 删除敏感词 |
| GET | `/api/admin/community/reports` | 举报列表 |
| PUT | `/api/admin/community/reports/:id/resolve` | 处理举报 |

---

## 四、开发任务分解

### Phase 1：数据库 + 基础 Model（1 个任务）
- [ ] 创建迁移 SQL：`migrations/add_community_tables.sql`
  - 建表：`posts`, `post_replies`, `post_likes`, `post_reply_likes`, `user_points`, `sensitive_words`, `post_reports`
  - 初始敏感词种子数据
  - `user_points` 表，给用户表加外键可选

### Phase 2：敏感词过滤 + 积分服务（2 个任务）
- [ ] `internal/service/sensitive_filter.go`：敏感词加载 + 替换逻辑（支持热加载）
- [ ] `internal/service/point_service.go`：积分增减 + 等级计算

### Phase 3：帖子 CRUD + 回帖（3 个任务）
- [ ] `internal/domain/post.go` + `post_reply.go`：Model 定义 + GORM 标签
- [ ] `internal/repository/post_repository.go`：数据库操作层
- [ ] `internal/service/post_service.go`：业务逻辑层（含过滤、积分、浏览计数）

### Phase 4：Handler + 路由注册（1 个任务）
- [ ] `internal/http/handlers/community_handler.go`：所有社区 API handler
- [ ] `cmd/server/main.go` 或 `internal/http/routes.go`：注册路由

### Phase 5：点赞 + 举报 + 热门（2 个任务）
- [ ] 点赞 toggle 逻辑（posts/replies 分别处理）
- [ ] 举报 API + 热门帖子计算（按 reply_count + like_count 排序）

### Phase 6：管理后台 API（1 个任务）
- [ ] 帖子审核/置顶/删除 + 敏感词管理 + 举报处理

### Phase 7：测试 + 部署（1 个任务）
- [ ] 编译验证 + 推送到 GitHub + 服务器 Docker 重建

---

## 五、待确认事项

1. **帖子关联套图** — `photoset_id` 字段已设计，发帖时传 `photoset_id` 即可关联，详情接口会返回套图摘要。✅ 已确认
2. **举报处理** — 举报后是否自动隐藏帖子，还是仅通知 admin？→ 建议：举报后不自动隐藏，admin 后台处理
3. **热门帖子算法** — 按 `(reply_count * 2 + like_count + view_count / 10)` 排序，时间衰减（7天内）✅ 建议方案，待确认

---

## 六、不涉及范围（v1 不做）

- ❌ Web 前端页面（等 App）
- ❌ 富文本编辑器（纯文本 + [img:URL] 标记）
- ❌ 图片直接在帖子中上传（图片走已有 `/api/upload` 接口，URL 嵌入 content）
- ❌ 消息通知（回帖/点赞通知，v2 做）
- ❌ 帖子编辑功能（v2 做）

---

*计划就绪，等待用户确认后开工。*
