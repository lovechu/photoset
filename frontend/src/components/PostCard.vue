<template>
  <article class="post-card" @click="goDetail">
    <span v-if="rank" class="rank-badge" :class="{ top3: rank <= 3 }">{{ rank }}</span>
    <div class="topic-main">
      <div class="topic-tags" v-if="post.category || post.topics?.length || post.tags?.length || post.is_pinned || post.is_essence">
        <span v-if="post.category" :class="['tlabel', tagColor(post.category)]">{{ post.category }}</span>
        <span v-for="t in (post.topics || []).slice(0, 2)" :key="'tp' + t.id" class="tlabel orange" @click.stop="goTopic(t.name)">#{{ t.name }}</span>
        <span v-for="t in (post.tags || []).slice(0, 2)" :key="'tg' + t.id" class="tlabel gray" @click.stop="goTag(t.name)">{{ t.name }}</span>
        <span v-if="post.is_pinned" class="tlabel red">置顶</span>
        <span v-if="post.is_essence" class="tlabel gold">精华</span>
      </div>

      <h3 class="topic-title" v-if="post.title">
        <router-link :to="`/community/post/${post.id}`" class="title-link">{{ post.title }}</router-link>
      </h3>
      <p class="topic-excerpt">{{ plainText }}</p>

      <div v-if="images.length" class="topic-thumbs" :class="'n-' + Math.min(images.length, 3)">
        <div v-for="(img, i) in images.slice(0, 3)" :key="i" class="thumb" @click.stop="previewImage(i)">
          <img :src="img" :alt="post.title || '帖子图片'" loading="lazy" />
          <span v-if="i === 2 && images.length > 3" class="thumb-more">+{{ images.length - 3 }}</span>
        </div>
      </div>

      <div class="topic-meta">
        <span class="author author-link" @click.stop="goProfile">{{ post.author_name || '匿名用户' }}</span>
        <span class="sep">·</span>
        <span>{{ formatTime(post.created_at) }}</span>
        <span v-if="post.author_ip_location" class="sep">·</span>
        <span v-if="post.author_ip_location">{{ post.author_ip_location }}</span>
        <el-tag v-if="post.is_following" size="small" type="info" effect="plain" class="follow-tag">已关注</el-tag>
        <span class="mini-actions">
          <button class="mini-btn" :class="{ liked: post.is_liked }" @click.stop="onLike" title="点赞">
            <el-icon :size="13"><Top /></el-icon>{{ post.like_count || 0 }}
          </button>
          <button class="mini-btn" @click.stop="onShare" title="分享">
            <el-icon :size="13"><Share /></el-icon>分享
          </button>
        </span>
      </div>
    </div>

    <div class="topic-stats">
      <div class="tstat"><span class="n">{{ post.reply_count || 0 }}</span><span class="lbl">回复</span></div>
      <div class="tstat"><span class="n">{{ formatNum(post.view_count || 0) }}</span><span class="lbl">浏览</span></div>
    </div>

    <div class="topic-last">
      <div>{{ formatTime(post.created_at) }}</div>
      <div>
        <span class="u author-link" @click.stop="goProfile">{{ post.author_name || '用户' }}</span>
        <span class="last-av author-link" :style="{ background: avatarColor(post.author_name || 'U') }" @click.stop="goProfile">{{ (post.author_name || 'U').charAt(0) }}</span>
      </div>
    </div>

    <el-image-viewer
      v-if="showViewer"
      :url-list="images"
      :initial-index="viewerIndex"
      @close="showViewer = false"
    />
  </article>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Top, Share } from '@element-plus/icons-vue'
import { extractImages, toPlainText, formatTime } from '@/utils/post'
import { togglePostLike, recordShare } from '@/api'
import { useUserStore } from '@/stores/user'

const props = defineProps({
  post: { type: Object, required: true },
  rank: { type: Number, default: 0 }
})
const emit = defineEmits(['like'])

const router = useRouter()
const userStore = useUserStore()
const showViewer = ref(false)
const viewerIndex = ref(0)

const images = computed(() => extractImages(props.post.content))
const plainText = computed(() => toPlainText(props.post.content, 200))

function goDetail() {
  router.push(`/community/post/${props.post.id}`)
}
function goProfile() {
  if (props.post.author_id) router.push(`/community/user/${props.post.author_id}`)
}
function goTopic(name) {
  router.push(`/community/topic/${encodeURIComponent(name)}`)
}
function goTag(name) {
  router.push(`/community/tag/${encodeURIComponent(name)}`)
}
function previewImage(i) {
  viewerIndex.value = i
  showViewer.value = true
}
async function onLike() {
  if (!userStore.isLoggedIn) {
    ElMessage.info('请先登录')
    return
  }
  try {
    await togglePostLike(props.post.id)
    props.post.is_liked = !props.post.is_liked
    props.post.like_count = (props.post.like_count || 0) + (props.post.is_liked ? 1 : -1)
    let liked = JSON.parse(localStorage.getItem('photoset_liked_posts') || '[]')
    if (props.post.is_liked) {
      if (!liked.includes(props.post.id)) liked.push(props.post.id)
    } else {
      liked = liked.filter(id => id !== props.post.id)
    }
    localStorage.setItem('photoset_liked_posts', JSON.stringify(liked))
    emit('like', props.post)
  } catch (e) {
    ElMessage.error('操作失败')
  }
}
async function onShare() {
  try {
    await recordShare(props.post.id, {})
    props.post.share_count = (props.post.share_count || 0) + 1
  } catch (e) { /* 静默 */ }
  try {
    await navigator.clipboard.writeText(`${window.location.origin}/community/post/${props.post.id}`)
    ElMessage.success('链接已复制')
  } catch (e) {
    ElMessage.info('分享链接已生成')
  }
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
</script>

<style scoped>
.post-card {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 14px 16px;
  background: #fff;
  border: 1px solid #e6e8eb;
  border-radius: 5px;
  margin-bottom: 10px;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  color: #33373d;
  font-size: 14px;
  line-height: 1.5;
}
.post-card:hover { background: #f8fafc; border-color: #cdd5e0; }
.rank-badge {
  flex-shrink: 0; width: 24px; height: 24px; border-radius: 4px;
  background: #eef1f5; color: #6b7280; font-size: 12px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
}
.rank-badge.top3 { background: #f5a623; color: #fff; }
.topic-main { flex: 1; min-width: 0; }
.topic-tags { display: flex; gap: 5px; margin-bottom: 5px; flex-wrap: wrap; }
.tlabel { padding: 1px 7px; border-radius: 3px; font-size: 10px; font-weight: 700; color: #fff; cursor: default; }
.tlabel.orange, .tlabel.gray { cursor: pointer; }
.tlabel.green { background: #5a9e6f; }
.tlabel.blue { background: #4f78d8; }
.tlabel.purple { background: #8b5cc7; }
.tlabel.orange { background: #d9902b; }
.tlabel.gray { background: #9aa1ac; }
.tlabel.red { background: #d04545; }
.tlabel.gold { background: #c9971f; }
.topic-title {
  font-size: 15px; font-weight: 600; color: #2b3038;
  margin: 0 0 4px; line-height: 1.45;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.title-link { color: inherit; text-decoration: none; }
.post-card:hover .topic-title { color: #4f6ef7; }
.topic-excerpt {
  color: #8d95a0; font-size: 12px; line-height: 1.5; margin: 0;
  display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden;
}
.topic-thumbs { display: grid; gap: 5px; margin-top: 10px; max-width: 300px; }
.topic-thumbs.n-1 { grid-template-columns: minmax(0, 120px); }
.topic-thumbs.n-2, .topic-thumbs.n-3 { grid-template-columns: repeat(3, 1fr); }
.thumb { position: relative; aspect-ratio: 1; border-radius: 3px; overflow: hidden; background: #f2f3f5; }
.thumb img { width: 100%; height: 100%; object-fit: cover; display: block; }
.thumb-more {
  position: absolute; inset: 0; background: rgba(0,0,0,0.55); color: #fff;
  display: flex; align-items: center; justify-content: center; font-size: 14px;
}
.topic-meta { display: flex; align-items: center; gap: 4px; margin-top: 10px; font-size: 11px; color: #9aa1ac; flex-wrap: wrap; }
.topic-meta .author { color: #6b7280; font-weight: 600; }
.author-link { cursor: pointer; }
.author-link:hover { color: #4f6ef7; text-decoration: underline; }
.topic-meta .sep { color: #c6cbd2; }
.follow-tag { margin-left: 4px; transform: scale(0.9); transform-origin: left center; }
.mini-actions { display: flex; gap: 12px; margin-left: auto; }
.mini-btn {
  display: inline-flex; align-items: center; gap: 3px;
  background: none; border: none; padding: 0;
  font-size: 11px; color: #9aa1ac; cursor: pointer;
}
.mini-btn:hover { color: #4f6ef7; }
.mini-btn.liked { color: #d04545; }
.topic-stats { display: flex; flex-shrink: 0; }
.tstat { width: 60px; text-align: center; }
.tstat .n { display: block; font-weight: 700; color: #3d444e; font-size: 14px; }
.tstat .lbl { display: block; font-size: 10px; color: #9aa1ac; margin-top: 1px; }
.topic-last {
  width: 120px; text-align: right; font-size: 11px; color: #9aa1ac;
  flex-shrink: 0; line-height: 1.6;
}
.topic-last .u { color: #6b7280; font-weight: 600; }
.last-av {
  width: 22px; height: 22px; border-radius: 50%;
  display: inline-flex; align-items: center; justify-content: center;
  color: #fff; font-size: 10px; font-weight: 600;
  margin-left: 6px; vertical-align: middle;
}

@media (max-width: 768px) {
  .post-card { gap: 10px; padding: 12px; }
  .topic-last { display: none; }
  .tstat { width: 46px; }
  .topic-title { white-space: normal; }
}
</style>
