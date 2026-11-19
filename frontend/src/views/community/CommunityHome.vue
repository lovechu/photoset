<template>
  <div class="community-home">
    <!-- 二级分类导航 -->
    <div class="subnav">
      <div class="subnav-inner">
        <a
          v-for="cat in navCats"
          :key="cat.slug || cat.name"
          href="javascript:;"
          :class="['subnav-link', { active: activeCategory === cat.slug }]"
          @click="activeCategory = cat.slug; reload()"
        >
          <span class="dot" :style="{ background: cat.color }"></span>{{ cat.name }}
        </a>
        <div class="subnav-spacer"></div>
        <button class="new-btn" @click="goCreate">＋ 新话题</button>
      </div>
    </div>

    <div class="container">
      <div class="layout">
        <main class="main">
          <div class="page-title">
            <div>
              <h1>社区</h1>
              <p>摄影作品分享与讨论</p>
            </div>
            <div class="view-tabs">
              <button
                v-for="t in sortTabs"
                :key="t.value"
                :class="{ active: sortBy === t.value }"
                @click="sortBy = t.value; reload()"
              >{{ t.label }}</button>
            </div>
          </div>

          <!-- 分类卡片 -->
          <div class="category-grid" v-if="categoryCards.length">
            <div
              v-for="(cat, i) in categoryCards"
              :key="cat.slug || cat.name"
              :class="['cat-card', 'c' + (i % 4 + 1)]"
              @click="activeCategory = cat.slug; reload()"
            >
              <div class="cat-name">{{ cat.name }}</div>
              <div class="cat-desc">{{ cat.description || defaultDesc(cat) }}</div>
              <div class="cat-stat" v-if="cat.post_count != null || cat.topic_count != null">
                <span><b>{{ formatNum(cat.post_count ?? cat.topic_count ?? 0) }}</b> 话题</span>
                <span v-if="cat.view_count != null"><b>{{ formatNum(cat.view_count) }}</b> 浏览</span>
              </div>
            </div>
          </div>

          <!-- 话题列表 -->
          <div class="topic-list" v-loading="loading">
            <div class="topic-head">
              <span class="c-topic">话题</span>
              <span class="c-stat">回复</span>
              <span class="c-stat">浏览</span>
              <span class="c-last">最后回复</span>
            </div>

            <div v-if="!loading && posts.length === 0" class="empty-state">
              <p>还没有话题</p>
              <button class="new-btn" @click="goCreate">发表第一帖</button>
            </div>

            <div v-for="post in posts" :key="post.id" class="topic" @click="goPost(post)">
              <div class="topic-main">
                <div class="topic-tags">
                  <span v-if="post.category" :class="['tlabel', tagColor(post.category)]">{{ post.category }}</span>
                  <span v-for="t in (post.topics || []).slice(0, 2)" :key="'tp' + t.id" class="tlabel orange">#{{ t.name }}</span>
                  <span v-for="t in (post.tags || []).slice(0, 1)" :key="'tg' + t.id" class="tlabel gray">{{ t.name }}</span>
                  <span v-if="post.is_pinned" class="tlabel red">置顶</span>
                  <span v-if="post.is_essence" class="tlabel gold">精华</span>
                </div>
                <div class="topic-title">
                  <router-link :to="`/community/post/${post.id}`" class="title-link" @click.stop>{{ post.title || plainText(post.content) || (postImages(post).length ? '[图片]' : '') }}</router-link>
                </div>
                <div v-if="postImages(post).length" class="topic-thumbs" @click.stop>
                  <div
                    v-for="(img, i) in postImages(post).slice(0, 3)"
                    :key="i"
                    class="topic-thumb"
                    @click.stop="previewImage(post, i)"
                  >
                    <img :src="img" :alt="post.title || '帖子图片'" loading="lazy" />
                    <span v-if="i === 2 && postImages(post).length > 3" class="thumb-more">+{{ postImages(post).length - 3 }}</span>
                  </div>
                </div>
                <div class="topic-sub">
                  <span class="author author-link" @click.stop="goProfile(post.author_id)">{{ post.author_name || '匿名用户' }}</span>
                  <span class="sep">·</span>{{ timeAgo(post.created_at) }}
                </div>
              </div>
              <div class="topic-stats">
                <div class="tstat"><span class="n">{{ post.reply_count || 0 }}</span><span :class="['heat', 'h' + heatLevel(post)]"></span></div>
                <div class="tstat"><span class="n">{{ formatNum(post.view_count || 0) }}</span></div>
              </div>
              <div class="topic-last">
                <div>{{ lastReplyTime(post) }}</div>
                <div>
                  <span class="u author-link" @click.stop="goProfile(lastReplyUserId(post))">{{ lastReplyName(post) }}</span>
                  <span class="last-av author-link" :style="{ background: avatarColor(lastReplyName(post)) }" @click.stop="goProfile(lastReplyUserId(post))">{{ lastReplyName(post).charAt(0) }}</span>
                </div>
              </div>
            </div>

            <div v-if="posts.length < total" class="load-more-row">
              <button @click="loadMore" :disabled="loadingMore">{{ loadingMore ? '加载中' : '加载更多话题' }}</button>
            </div>
          </div>
        </main>

        <aside class="side">
          <div class="card">
            <h3>关于社区</h3>
            <div class="card-body">
              一个专注摄影作品分享与交流的社区。在这里发现灵感、结交同好、一起进步。
              <div class="about-list">
                <div class="about-row"><span>已有话题</span><b>{{ total }}</b></div>
                <div class="about-row"><span>话题分类</span><b>{{ categories.length }}</b></div>
              </div>
              <router-link to="/community/guidelines" class="more-link">查看社区公约 →</router-link>
            </div>
          </div>

          <div class="card">
            <h3>热门话题</h3>
            <div class="card-body top-list">
              <div v-for="(p, i) in hotTopics" :key="p.id" class="top-link" @click="goPost(p)">
                <span class="rank">{{ i + 1 }}</span>
                <div>
                  <div class="tn">{{ p.title || plainText(p.content, 30) }}</div>
                  <div class="tm">{{ p.reply_count || 0 }} 回复 · {{ formatNum(p.view_count || 0) }} 浏览</div>
                </div>
              </div>
              <div v-if="!hotTopics.length" class="muted">暂无热门话题</div>
            </div>
          </div>

          <div class="card">
            <h3>社区公约</h3>
            <div class="card-body">
              尊重他人 · 友善交流<br />
              原创优先 · 注明出处<br />
              拒绝广告与违规内容
            </div>
          </div>
        </aside>
      </div>
    </div>

    <el-image-viewer
      v-if="showViewer"
      :url-list="viewerImages"
      :initial-index="viewerIndex"
      @close="showViewer = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useCommunityStore } from '@/stores/community'
import { useUserStore } from '@/stores/user'
import { getPosts } from '@/api'
import { toPlainText, extractImages } from '@/utils/post'
import { useSeo } from '@/composables/useSeo'

useSeo({ title: '社区', description: 'PhotoSet 摄影社区，浏览用户分享的摄影作品、热门话题与原创讨论。', path: '/community' })

const router = useRouter()
const route = useRoute()
const communityStore = useCommunityStore()
const userStore = useUserStore()

const posts = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const loadingMore = ref(false)
const showViewer = ref(false)
const viewerImages = ref([])
const viewerIndex = ref(0)
const activeCategory = ref(route.query.category || '')
const sortBy = ref('latest')
const categories = ref([])

const sortTabs = [
  { label: '最新', value: 'latest' },
  { label: '热门', value: 'hot' },
  { label: '最多回复', value: 'reply_count' },
  { label: '最多点赞', value: 'like_count' }
]
const categoryColors = ['#e07a5f', '#4f6ef7', '#81b29a', '#9b5de5']

const navCats = computed(() => [
  { name: '全部', slug: '', color: '#6b7280' },
  ...categories.value.map((c, i) => ({ ...c, color: categoryColors[i % 4] }))
])
const categoryCards = computed(() =>
  categories.value.map((c, i) => ({ ...c, color: categoryColors[i % 4] }))
)
const hotTopics = computed(() =>
  [...posts.value].sort((a, b) => (b.view_count || 0) - (a.view_count || 0)).slice(0, 5)
)

onMounted(async () => {
  categories.value = await communityStore.fetchCategories()
  loadPosts()
})

function reload() {
  page.value = 1
  posts.value = []
  loadPosts()
}

async function loadPosts() {
  loading.value = true
  try {
    const res = await getPosts({
      page: page.value,
      page_size: 20,
      category: activeCategory.value || undefined,
      sort_by: sortBy.value
    })
    const data = res?.data || {}
    const items = data.posts || []
    posts.value = page.value === 1 ? items : [...posts.value, ...items]
    total.value = data.pagination?.total || 0
  } catch (e) {
    console.warn('加载帖子失败:', e.message)
  } finally {
    loading.value = false
  }
}

function loadMore() {
  if (loadingMore.value) return
  loadingMore.value = true
  page.value++
  loadPosts().finally(() => { loadingMore.value = false })
}

function goPost(post) {
  router.push(`/community/post/${post.id}`)
}

function goProfile(userId) {
  if (userId) router.push(`/community/user/${userId}`)
}

function lastReplyUserId(post) {
  return post.last_reply_user?.id || post.last_reply_user_id || post.author_id
}

function goCreate() {
  if (!userStore.isLoggedIn) {
    router.push({ name: 'Login', query: { redirect: '/community/create' } })
    return
  }
  router.push('/community/create')
}

function plainText(content, len = 180) {
  return toPlainText(content, len)
}

function postImages(post) {
  return extractImages(post?.content)
}

function previewImage(post, index) {
  viewerImages.value = postImages(post)
  viewerIndex.value = index
  showViewer.value = true
}

function defaultDesc(cat) {
  const map = {
    '作品展示': '分享你的摄影作品',
    '问答': '摄影技巧与器材求助',
    '讨论': '任意话题自由交流',
    '建议': '产品反馈与功能建议'
  }
  return map[cat.name] || '参与讨论，分享你的观点'
}

function formatNum(n) {
  const num = Number(n) || 0
  return num >= 1000 ? (num / 1000).toFixed(1).replace(/\.0$/, '') + 'k' : String(num)
}

function hashOf(str) {
  let h = 0
  for (let i = 0; i < (str || '').length; i++) h = (h * 31 + str.charCodeAt(i)) >>> 0
  return h
}

function tagColor(name) {
  const colors = ['green', 'blue', 'purple', 'orange']
  return colors[hashOf(name) % colors.length]
}

function avatarColor(name) {
  const palette = ['#4f6ef7', '#e07a5f', '#81b29a', '#9b5de5', '#3d5a80', '#d9902b']
  return palette[hashOf(name) % palette.length]
}

function heatLevel(post) {
  const n = (post.reply_count || 0) + Math.floor((post.view_count || 0) / 100)
  if (n >= 50) return 5
  if (n >= 30) return 4
  if (n >= 15) return 3
  if (n >= 5) return 2
  return 1
}

function lastReplyName(post) {
  return post.last_reply_user?.nickname || post.last_reply_nickname ||
    post.last_reply_username || post.author_name || '用户'
}

function lastReplyTime(post) {
  return post.last_reply_at ? timeAgo(post.last_reply_at) : timeAgo(post.created_at)
}

function timeAgo(str) {
  if (!str) return ''
  const d = new Date(str)
  if (isNaN(d.getTime())) return str
  const mins = Math.floor((Date.now() - d.getTime()) / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins}分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}天前`
  return `${d.getMonth() + 1}月${d.getDate()}日`
}
</script>

<style scoped>
.community-home {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  color: #33373d;
  font-size: 14px;
  line-height: 1.5;
}

/* 二级分类导航 */
.subnav {
  background: #f7f8f9;
  border: 1px solid #e6e8eb;
  border-radius: 5px;
  height: 46px;
}
.subnav-inner {
  height: 100%;
  display: flex;
  align-items: center;
  gap: 26px;
  padding: 0 20px;
  overflow-x: auto;
}
.subnav-link {
  font-size: 13px;
  color: #5a606a;
  display: flex;
  align-items: center;
  gap: 7px;
  white-space: nowrap;
  height: 100%;
  cursor: pointer;
  text-decoration: none;
}
.subnav-link:hover { color: #2b3038; }
.subnav-link.active { color: #2b3038; font-weight: 600; }
.subnav-link .dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.subnav-spacer { flex: 1; }
.new-btn {
  background: #f5a623;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 7px 16px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
}
.new-btn:hover { background: #e8991c; }

.container { max-width: 1220px; margin: 0 auto; padding: 24px 6px 60px; }
.layout { display: grid; grid-template-columns: 1fr 290px; gap: 22px; }
.main { min-width: 0; }

/* 页面标题 */
.page-title { display: flex; align-items: flex-end; justify-content: space-between; margin-bottom: 18px; gap: 12px; }
.page-title h1 { font-size: 24px; color: #2b3038; font-weight: 700; }
.page-title p { color: #8d95a0; font-size: 13px; margin-top: 4px; }
.view-tabs { display: flex; gap: 4px; flex-wrap: wrap; }
.view-tabs button {
  padding: 6px 12px;
  border: 1px solid #e0e3e7;
  background: #fff;
  color: #6b7280;
  font-size: 12px;
  cursor: pointer;
  border-radius: 4px;
}
.view-tabs button.active { background: #eef1f5; color: #2b3038; border-color: #d5dae1; font-weight: 600; }

/* 分类卡片 */
.category-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 24px; }
.cat-card {
  background: #fff;
  border: 1px solid #e6e8eb;
  border-radius: 5px;
  padding: 14px 14px 12px;
  position: relative;
  overflow: hidden;
  cursor: pointer;
  transition: border-color 0.15s;
}
.cat-card:hover { border-color: #b8c0cc; }
.cat-card::before { content: ''; position: absolute; left: 0; top: 0; bottom: 0; width: 4px; }
.cat-card.c1::before { background: #e07a5f; }
.cat-card.c2::before { background: #4f6ef7; }
.cat-card.c3::before { background: #81b29a; }
.cat-card.c4::before { background: #9b5de5; }
.cat-name { font-size: 14px; font-weight: 700; color: #2b3038; margin-bottom: 4px; padding-left: 8px; }
.cat-desc { font-size: 11px; color: #9aa1ac; padding-left: 8px; min-height: 30px; line-height: 1.4; }
.cat-stat { display: flex; gap: 12px; font-size: 11px; color: #8d95a0; padding-left: 8px; margin-top: 8px; }
.cat-stat b { color: #5a606a; }

/* 话题列表 */
.topic-list { background: #fff; border: 1px solid #e6e8eb; border-radius: 5px; overflow: hidden; }
.topic-head {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  background: #fafbfc;
  border-bottom: 1px solid #eef0f2;
  color: #9aa1ac;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.topic-head .c-topic { flex: 1; }
.topic-head .c-stat { width: 74px; text-align: center; }
.topic-head .c-last { width: 140px; text-align: right; }
.topic {
  display: flex;
  align-items: center;
  padding: 13px 16px;
  border-bottom: 1px solid #f0f1f3;
  cursor: pointer;
  transition: background 0.15s;
}
.topic:last-child { border-bottom: none; }
.topic:hover { background: #f8fafc; }
.topic-main { flex: 1; min-width: 0; }
.topic-tags { display: flex; gap: 5px; margin-bottom: 4px; flex-wrap: wrap; }
.tlabel { padding: 1px 7px; border-radius: 3px; font-size: 10px; font-weight: 700; color: #fff; }
.tlabel.green { background: #5a9e6f; }
.tlabel.blue { background: #4f78d8; }
.tlabel.purple { background: #8b5cc7; }
.tlabel.orange { background: #d9902b; }
.tlabel.gray { background: #9aa1ac; }
.tlabel.red { background: #d04545; }
.tlabel.gold { background: #c9971f; }
.topic-title {
  font-size: 14px;
  font-weight: 600;
  color: #2b3038;
  margin-bottom: 3px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.topic:hover .topic-title { color: #4f6ef7; }
.title-link { color: inherit; text-decoration: none; }
.topic-thumbs { display: flex; gap: 5px; margin: 6px 0; max-width: 320px; }
.topic-thumb { position: relative; width: 72px; height: 54px; border-radius: 4px; overflow: hidden; background: #f2f3f5; cursor: zoom-in; flex-shrink: 0; }
.topic-thumb img { width: 100%; height: 100%; object-fit: cover; display: block; }
.thumb-more {
  position: absolute; inset: 0; background: rgba(0,0,0,0.5); color: #fff;
  display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 600;
}
.topic-sub { font-size: 11px; color: #9aa1ac; }
.topic-sub .author { color: #6b7280; font-weight: 600; }
.author-link { cursor: pointer; }
.author-link:hover { color: #4f6ef7; text-decoration: underline; }
.topic-sub .sep { margin: 0 4px; color: #c6cbd2; }
.topic-stats { display: flex; }
.tstat { width: 74px; text-align: center; font-size: 12px; color: #6b7280; }
.tstat .n { font-weight: 700; color: #3d444e; font-size: 13px; }
.heat { display: inline-block; width: 9px; height: 9px; border-radius: 2px; margin-left: 3px; vertical-align: middle; }
.h1 { background: #e8f0e8; }
.h2 { background: #b8d4b8; }
.h3 { background: #7fb87f; }
.h4 { background: #e08a3f; }
.h5 { background: #d04545; }
.topic-last { width: 140px; text-align: right; font-size: 11px; color: #9aa1ac; }
.topic-last .u { color: #6b7280; font-weight: 600; }
.last-av {
  width: 22px; height: 22px; border-radius: 50%;
  display: inline-flex; align-items: center; justify-content: center;
  color: #fff; font-size: 10px; font-weight: 600;
  margin-left: 6px; vertical-align: middle;
}

.empty-state { text-align: center; padding: 70px 0; color: #9aa1ac; font-size: 14px; }
.empty-state p { margin: 0 0 18px; }
.load-more-row { text-align: center; padding: 18px; border-top: 1px solid #f0f1f3; }
.load-more-row button {
  padding: 7px 24px;
  border: 1px solid #e0e3e7;
  background: #fff;
  color: #6b7280;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}
.load-more-row button:hover { border-color: #b8c0cc; color: #2b3038; }
.load-more-row button:disabled { opacity: 0.5; cursor: default; }

/* 侧栏 */
.side { display: flex; flex-direction: column; gap: 16px; }
.card { background: #fff; border: 1px solid #e6e8eb; border-radius: 5px; overflow: hidden; }
.card h3 {
  font-size: 12px; font-weight: 700; color: #3d444e;
  padding: 12px 16px;
  border-bottom: 1px solid #eef0f2;
  background: #fafbfc;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.card-body { padding: 12px 16px; font-size: 12px; color: #6b7280; line-height: 1.7; }
.about-list { margin-top: 10px; }
.about-row { display: flex; justify-content: space-between; padding: 5px 0; border-bottom: 1px solid #f2f3f5; }
.about-row:last-child { border-bottom: none; }
.about-row b { color: #3d444e; }
.more-link { display: inline-block; margin-top: 10px; color: #4f6ef7; font-size: 12px; text-decoration: none; }
.top-list { padding: 8px 16px; }
.top-link { display: flex; gap: 8px; padding: 8px 0; border-bottom: 1px solid #f2f3f5; cursor: pointer; }
.top-link:last-child { border-bottom: none; }
.top-link:hover .tn { color: #4f6ef7; }
.top-link .rank { color: #d9902b; font-weight: 700; font-size: 12px; width: 14px; flex-shrink: 0; }
.top-link .tn { font-size: 12px; color: #3d444e; line-height: 1.4; }
.top-link .tm { font-size: 10px; color: #9aa1ac; margin-top: 2px; }
.muted { color: #9aa1ac; }

@media (max-width: 1024px) {
  .layout { grid-template-columns: 1fr; }
  .side { display: grid; grid-template-columns: repeat(3, 1fr); }
}
@media (max-width: 768px) {
  .container { padding: 16px 4px 48px; }
  .category-grid { grid-template-columns: repeat(2, 1fr); }
  .side { grid-template-columns: 1fr; }
  .topic-head .c-last, .topic-last { display: none; }
  .topic-head .c-stat, .tstat { width: 52px; }
  .page-title h1 { font-size: 20px; }
}
</style>
