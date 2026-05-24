# PhotoSet 社区功能 — PRD & 开发计划

> 版本：v1.1
> 日期：2026-05-24
> 状态：已评审完善
> 作者：许清楚（产品经理）

---

## 一、产品概述

在 PhotoSet 中新增**社区模块**，用户可发帖、回帖、点赞、举报，内容先发后审。
前端暂不开发（Web 端风险高），**仅做后端 API**，供后续 App 集成。

### 核心原则
- **模块独立性**：社区和套图是完全独立的两个模块，帖子不能出现在套图页面，套图页面不展示社区内容
- **先发后审**：敏感词替换为 `***` 后放行，不阻止发布
- **仅后端 API**：Web 前端不做，只出后端 API（供 App 用）
- **积分等级规则**：确认合理，按草案执行

---

## 二、用户故事

### 2.1 核心用户故事

1. **发帖**：作为一名用户，我想要发布帖子（含标题、内容、分类），以便分享经验或提问
2. **回帖**：作为一名用户，我想要回复帖子或回复其他回帖（楼中楼），以便参与讨论
3. **点赞**：作为一名用户，我想要对帖子或回帖点赞/取消点赞，以便表达认同
4. **浏览帖子**：作为一名用户（或游客），我想要浏览帖子列表和详情，以便获取信息
5. **举报**：作为一名用户，我想要举报不恰当的帖子或回帖，以便维护社区质量
6. **积分激励**：作为一名用户，我想要通过发帖/回帖/被点赞获得积分和等级，以便获得成就感
7. **敏感词过滤**：系统应自动过滤敏感词，替换为 `***`，以便维护社区文明

### 2.2 管理后台用户故事

1. **帖子管理**：作为管理员，我想要置顶/审核/删除帖子，以便管理社区内容
2. **回帖管理**：作为管理员，我想要删除不恰当的回帖
3. **敏感词管理**：作为管理员，我想要增删改敏感词库，以便动态调整过滤策略
4. **举报处理**：作为管理员，我想要查看和处理举报，以便响应社区反馈
5. **用户管理**：作为管理员，我想要查看用户积分/等级，并可以调整，以便运营干预

---

## 三、核心设计

### 3.1 数据模型

#### `posts` — 帖子主表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | 发帖人 |
| title | varchar(200) | 标题（必填，1-200字符） |
| content | text | 内容（必填，1-5000字符，含 [img:URL] 标记） |
| photoset_id | uint (nullable) | 关联套图（可选） |
| category | varchar(20) | 分类：discussion/qa/showcase/suggestion |
| visibility | varchar(20) | 可见范围：public/member/vip/admin |
| is_pinned | bool | 置顶（默认 false） |
| view_count | int | 浏览数（默认 0） |
| reply_count | int | 回帖数（默认 0，与 post_replies 表同步） |
| like_count | int | 点赞数（默认 0，与 post_likes 表同步） |
| status | varchar(20) | approved/pending/rejected（先发后审，默认 approved） |
| created_at / updated_at | datetime | |

**索引设计：**
- INDEX `idx_user_id` (user_id)
- INDEX `idx_category_visibility` (category, visibility)
- INDEX `idx_created_at` (created_at)
- INDEX `idx_is_pinned_created_at` (is_pinned, created_at)

#### `post_replies` — 回帖表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| post_id | uint | 所属帖子 |
| user_id | uint | 回帖人 |
| content | text | 内容（必填，1-2000字符） |
| parent_reply_id | uint (nullable) | 楼中楼（回复某条回复，null 表示直接回帖） |
| like_count | int | 点赞数（默认 0，与 post_reply_likes 表同步） |
| created_at / updated_at | datetime | |

**索引设计：**
- INDEX `idx_post_id` (post_id)
- INDEX `idx_user_id` (user_id)
- INDEX `idx_parent_reply_id` (parent_reply_id)

#### `post_likes` — 帖子点赞记录
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | |
| post_id | uint | |
| created_at | datetime | 联合唯一索引 (user_id, post_id) |

**索引设计：**
- UNIQUE INDEX `idx_user_post` (user_id, post_id)
- INDEX `idx_post_id` (post_id)

#### `post_reply_likes` — 回帖点赞记录
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | |
| reply_id | uint | |
| created_at | datetime | 联合唯一索引 (user_id, reply_id) |

**索引设计：**
- UNIQUE INDEX `idx_user_reply` (user_id, reply_id)
- INDEX `idx_reply_id` (reply_id)

#### `user_points` — 用户积分/等级
| 字段 | 类型 | 说明 |
|------|------|------|
| user_id | uint | 主键，外键关联 users.id |
| points | int | 总积分（默认 0，范围 0-2147483647） |
| level | int | 等级 1-10（根据积分自动计算） |
| updated_at | datetime | |

**注意：** points 字段使用 int 类型，上限约 21 亿，正常情况下不会溢出。如担心溢出，可改用 bigint。

#### `sensitive_words` — 敏感词库
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| word | varchar(100) | 敏感词（唯一索引，存储时转小写） |
| replacement | varchar(100) | 替换词，默认 `***` |
| is_active | bool | 是否启用（默认 true） |
| created_at | datetime | |

**重要说明：**
- 敏感词匹配时**大小写不敏感**（存储时转小写，匹配时输入也转小写）
- 支持模糊匹配（可选，v1 先做精确匹配）

#### `post_reports` — 举报记录
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| post_id | uint (nullable) | 举报帖子（与 reply_id 二选一） |
| reply_id | uint (nullable) | 举报回帖（与 post_id 二选一） |
| reporter_id | uint | 举报人 |
| reason | varchar(500) | 举报原因（必填） |
| status | varchar(20) | pending/resolved/rejected（默认 pending） |
| handler_id | uint (nullable) | 处理人（admin） |
| handled_at | datetime (nullable) | 处理时间 |
| handle_note | varchar(500) | 处理备注 |
| created_at | datetime | |

**索引设计：**
- INDEX `idx_status` (status)
- INDEX `idx_post_id` (post_id)
- INDEX `idx_reply_id` (reply_id)

---

### 3.2 积分/等级规则

| 行为 | 积分 | 说明 |
|------|------|------|
| 发帖 | +10 | 每日上限 50（每日最多 5 篇） |
| 回帖 | +5 | 每日上限 30（每日最多 6 条） |
| 帖子被点赞 | +2 | 无上限 |
| 回帖被点赞 | +1 | 无上限 |
| 帖子被删除（违规） | -20 | 管理员删除时扣除 |

**等级对照表：**

| 等级 | 所需积分 | 称号 |
|------|----------|------|
| L1 | 0 | 新手上路 |
| L2 | 100 | 活跃会员 |
| L3 | 500 | 资深会员 |
| L4 | 2000 | 金牌会员 |
| L5 | 5000 | 钻石会员 |
| L6 | 10000 | 至尊会员 |
| L7 | 20000 | 荣耀会员 L7 |
| L8 | 30000 | 荣耀会员 L8 |
| L9 | 40000 | 荣耀会员 L9 |
| L10 | 50000 | 荣耀会员 L10 |

**积分规则说明：**
- 积分可以为负数（最低 -9999），但等级最低为 L1
- 积分溢出保护：points 字段使用 int 类型，上限 21 亿，正常情况下不会达到
- 每日上限按自然日计算（UTC+8 时区）

---

### 3.3 敏感词过滤逻辑

1. 发帖/回帖时，对 `title` + `content` 执行过滤
2. 从 `sensitive_words` 表加载所有 `is_active=true` 的词（启动时加载到内存，支持热加载）
3. **大小写不敏感匹配**（输入转小写后匹配）
4. 匹配到的词替换为 `replacement`（默认 `***`）
5. 过滤后的内容存入数据库
6. 管理后台可增删改敏感词，修改后自动热加载

**示例：**
- 敏感词：`badword`
- 输入：`This is a BadWord test`
- 输出：`This is a *** test`

---

### 3.4 权限模型

**发帖权限：** 登录即可发帖（无积分/等级限制）

**可见范围（visibility）：**
- `public`：所有人可见（包括未登录）
- `member`：会员及以上可见（user.role >= member）
- `vip`：VIP 及以上可见（user.role >= vip）
- `admin`：仅管理员可见

**API 层过滤逻辑：**
- 未登录：只返回 `visibility=public` 且 `status=approved` 的帖子
- 登录用户：根据自己角色过滤 `visibility`，只返回 `status=approved` 的帖子
- 管理员：可查看所有状态的帖子

**注意：** 社区内容**不会**出现在套图页面，套图页面**不会**展示社区帖子。两个模块完全独立。

---

### 3.5 数据一致性保证

**计数字段与真实记录同步：**
- `posts.like_count` = COUNT(post_likes WHERE post_id = posts.id)
- `posts.reply_count` = COUNT(post_replies WHERE post_id = posts.id)
- `post_replies.like_count` = COUNT(post_reply_likes WHERE reply_id = post_replies.id)

**实现方式：**
- 每次点赞/取消点赞时，更新对应计数
- 每次回帖/删除回帖时，更新帖子回帖数
- 提供定时任务（可选），用于修复计数不一致的情况

---

## 四、API 接口清单

### 4.1 社区公开 API（需登录态，部分支持匿名）

| 方法 | 路径 | 说明 | 登录要求 | 备注 |
|------|------|------|----------|------|
| GET | `/api/community/posts` | 帖子列表（支持分页、分类、排序） | 可选 | 未登录只返回 public |
| GET | `/api/community/posts/:id` | 帖子详情 + 浏览+1 | 可选 | 未登录只返回 public |
| GET | `/api/community/posts/:id/replies` | 回帖列表（楼中楼树形结构） | 可选 | |
| POST | `/api/community/posts` | 发帖 | **必需** | 需敏感词过滤 |
| POST | `/api/community/posts/:id/reply` | 回帖 | **必需** | 需敏感词过滤 |
| POST | `/api/community/posts/:id/like` | 点赞帖子（toggle） | **必需** | 重复调用取消点赞 |
| POST | `/api/community/replies/:id/like` | 点赞回帖（toggle） | **必需** | 重复调用取消点赞 |
| POST | `/api/community/posts/:id/report` | 举报帖子 | **必需** | |
| POST | `/api/community/replies/:id/report` | 举报回帖 | **必需** | |
| GET | `/api/community/categories` | 分类列表 | 可选 | 返回固定 4 个分类 |
| GET | `/api/community/hot` | 热门帖子 | 可选 | 按算法排序 |
| GET | `/api/community/me/posts` | 我的帖子列表 | **必需** | 分页查询当前用户帖子 |
| GET | `/api/community/me/replies` | 我的回帖列表 | **必需** | 分页查询当前用户回帖 |
| GET | `/api/community/me/points` | 我的积分/等级信息 | **必需** | 返回积分、等级、下一等级所需积分 |

---

### 4.2 管理后台 API（admin 权限）

| 方法 | 路径 | 说明 | 备注 |
|------|------|------|------|
| GET | `/api/admin/community/posts` | 帖子列表（含 status 筛选） | 支持分页、状态筛选 |
| PUT | `/api/admin/community/posts/:id/pin` | 置顶/取消置顶 | toggle |
| PUT | `/api/admin/community/posts/:id/status` | 修改状态（approve/reject） | |
| DELETE | `/api/admin/community/posts/:id` | 删除帖子（硬删除） | 同时扣除积分 |
| GET | `/api/admin/community/replies` | 回帖列表 | 支持分页、帖子筛选 |
| DELETE | `/api/admin/community/replies/:id` | 删除回帖（硬删除） | |
| GET | `/api/admin/community/keywords` | 敏感词列表 | 支持分页、启用状态筛选 |
| POST | `/api/admin/community/keywords` | 添加敏感词 | |
| PUT | `/api/admin/community/keywords/:id` | 修改敏感词 | 修改后热加载 |
| DELETE | `/api/admin/community/keywords/:id` | 删除敏感词 | |
| GET | `/api/admin/community/reports` | 举报列表 | 支持状态筛选 |
| PUT | `/api/admin/community/reports/:id/resolve` | 处理举报（resolve/reject） | 可附加处理备注 |
| GET | `/api/admin/community/users` | 用户积分/等级列表 | 支持分页、等级筛选 |
| PUT | `/api/admin/community/users/:id/points` | 调整用户积分 | 管理员手动调整 |
| GET | `/api/admin/community/stats` | 社区数据统计 | 返回帖子数、回帖数、用户数等 |

---

## 五、边界情况处理

### 5.1 输入验证

| 场景 | 处理方式 |
|------|----------|
| 发帖标题为空 | 返回错误：标题不能为空 |
| 发帖标题超长（>200字符） | 返回错误：标题不能超过200字符 |
| 发帖内容为空 | 返回错误：内容不能为空 |
| 发帖内容超长（>5000字符） | 返回错误：内容不能超过5000字符 |
| 回帖内容为空 | 返回错误：回复内容不能为空 |
| 回帖内容超长（>2000字符） | 返回错误：回复内容不能超过2000字符 |
| 举报原因不能为空 | 返回错误：请填写举报原因 |

### 5.2 重复操作

| 场景 | 处理方式 |
|------|----------|
| 重复点赞帖子 | toggle 逻辑：再次调用取消点赞 |
| 重复点赞回帖 | toggle 逻辑：再次调用取消点赞 |
| 重复回帖（内容完全相同，1分钟内） | 返回错误：请勿重复提交 |
| 自己点赞自己的帖子/回帖 | 允许（无限制） |

### 5.3 积分溢出

| 场景 | 处理方式 |
|------|----------|
| 积分超过 INT 上限（21亿） | 封顶到 21亿（实际不会达到） |
| 积分为负数（最低 -9999） | 允许，但等级最低为 L1 |

### 5.4 敏感词大小写

| 场景 | 处理方式 |
|------|----------|
| 敏感词 `badword`，输入 `BadWord` / `BADWORD` | 大小写不敏感，全部替换为 `***` |

### 5.5 并发问题

| 场景 | 处理方式 |
|------|----------|
| 多个用户同时点赞同一帖子 | 使用数据库事务 + 行锁，保证计数准确 |
| 多个用户同时回帖 | 使用数据库事务，保证回帖数和计数一致 |

---

## 六、开发任务分解

### Phase 1：数据库 + 基础 Model（1 个任务）
- [ ] 创建迁移 SQL：`migrations/add_community_tables.sql`
  - 建表：`posts`, `post_replies`, `post_likes`, `post_reply_likes`, `user_points`, `sensitive_words`, `post_reports`
  - 初始敏感词种子数据（常见敏感词 50 个）
  - `user_points` 表，给用户表加外键
  - 所有索引创建

### Phase 2：敏感词过滤 + 积分服务（2 个任务）
- [ ] `internal/service/sensitive_filter.go`：敏感词加载 + 替换逻辑（支持热加载、大小写不敏感）
- [ ] `internal/service/point_service.go`：积分增减 + 等级计算（含每日上限、溢出保护）

### Phase 3：帖子 CRUD + 回帖（3 个任务）
- [ ] `internal/domain/post.go` + `post_reply.go`：Model 定义 + GORM 标签 + 输入验证
- [ ] `internal/repository/post_repository.go`：数据库操作层（含计数同步逻辑）
- [ ] `internal/service/post_service.go`：业务逻辑层（含过滤、积分、浏览计数、数据一致性保证）

### Phase 4：Handler + 路由注册（1 个任务）
- [ ] `internal/http/handlers/community_handler.go`：所有社区 API handler（含输入验证、边界情况处理）
- [ ] `cmd/server/main.go` 或 `internal/http/routes.go`：注册路由（公开 API + 管理后台 API）

### Phase 5：点赞 + 举报 + 热门（2 个任务）
- [ ] 点赞 toggle 逻辑（posts/replies 分别处理，含并发控制）
- [ ] 举报 API + 热门帖子计算（按 `(reply_count * 2 + like_count + view_count / 10)` 排序，时间衰减（7天内））

### Phase 6：管理后台 API（1 个任务）
- [ ] 帖子审核/置顶/删除 + 敏感词管理 + 举报处理 + 用户积分管理 + 数据统计

### Phase 7：测试 + 部署（1 个任务）
- [ ] 单元测试（service 层）+ 集成测试（API 层）+ 编译验证 + 推送到 GitHub + 服务器 Docker 重建

---

## 七、待确认事项

1. ~~**帖子关联套图** — `photoset_id` 字段已设计，发帖时传 `photoset_id` 即可关联，详情接口会返回套图摘要。~~ ✅ 已确认
2. ~~**举报处理** — 举报后是否自动隐藏帖子，还是仅通知 admin？~~ → ✅ 已确认：举报后不自动隐藏，admin 后台处理
3. ~~**热门帖子算法** — 按 `(reply_count * 2 + like_count + view_count / 10)` 排序，时间衰减（7天内）~~ → ✅ 已确认
4. **楼中楼深度限制** — 是否限制楼中楼深度（建议：最多 3 层）？→ 待确认
5. **帖子删除方式** — 硬删除还是软删除（标记删除）？→ 建议：管理员硬删除，普通用户不能删除（v2 做编辑/删除功能）

---

## 八、不涉及范围（v1 不做）

- ❌ Web 前端页面（等 App）
- ❌ 富文本编辑器（纯文本 + [img:URL] 标记）
- ❌ 图片直接在帖子中上传（图片走已有 `/api/upload` 接口，URL 嵌入 content）
- ❌ 消息通知（回帖/点赞通知，v2 做）
- ❌ 帖子编辑功能（v2 做）
- ❌ 帖子删除功能（普通用户，v2 做）
- ❌ 搜索功能（v2 做）
- ❌ 关注功能（v2 做）
- ❌ 私信功能（v2 做）

---

## 九、上线后的运营建议

1. **敏感词库维护**：定期更新敏感词库，根据社区情况动态调整
2. **举报处理时效**：建议 24 小时内处理举报
3. **热门帖子算法调优**：根据上线后的数据，调整热门帖子排序算法
4. **积分等级调整**：根据上线后的用户活跃度，调整积分规则和等级门槛
5. **数据分析**：定期查看社区数据统计（帖子数、回帖数、用户活跃度等），优化产品

---

*PRD 已完善，可开工。*

---

## 附录：API 请求/响应示例

### 发帖
**请求：**
```json
POST /api/community/posts
{
  "title": "如何拍摄夜景？",
  "content": "大家好，我想请教一下夜景拍摄的技巧...",
  "category": "qa",
  "visibility": "public",
  "photoset_id": 123
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "title": "如何拍摄夜景？",
    "content": "大家好，我想请教一下夜景拍摄的技巧...",
    "category": "qa",
    "visibility": "public",
    "user_id": 456,
    "created_at": "2026-05-24T10:00:00Z"
  }
}
```

### 回帖
**请求：**
```json
POST /api/community/posts/1/reply
{
  "content": "我觉得可以用三脚架...",
  "parent_reply_id": null
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "post_id": 1,
    "user_id": 456,
    "content": "我觉得可以用三脚架...",
    "created_at": "2026-05-24T10:05:00Z"
  }
}
```

### 点赞帖子（toggle）
**请求：**
```
POST /api/community/posts/1/like
```

**响应（点赞）：**
```json
{
  "code": 0,
  "message": "liked",
  "data": {
    "like_count": 11
  }
}
```

**响应（取消点赞）：**
```json
{
  "code": 0,
  "message": "unliked",
  "data": {
    "like_count": 10
  }
}
```
