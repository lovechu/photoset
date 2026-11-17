<template>
  <div class="comment-section">
    <div class="section-header">
      <h3>评论 <span class="count">{{ commentCount }}</span></h3>
      <el-select v-model="sortBy" size="small" style="width: 110px" @change="reload">
        <el-option label="最新" value="newest" />
        <el-option label="最热" value="hottest" />
      </el-select>
    </div>

    <!-- 评论输入框 -->
    <div class="comment-compose" v-if="userStore.isLoggedIn">
      <el-avatar :size="36" :src="userStore.avatar">
        {{ userStore.user?.nickname?.charAt(0) }}
      </el-avatar>
      <div class="compose-body">
        <el-input
          v-model="newComment"
          type="textarea"
          :rows="2"
          maxlength="500"
          show-word-limit
          placeholder="说点什么吧..."
        />
        <div class="compose-actions">
          <el-button size="small" :icon="Picture" @click="openImagePicker">配图</el-button>
          <input ref="imageInput" type="file" accept="image/*" style="display:none" @change="onImagePicked" />
          <span v-if="imageUrl" class="image-preview">
            <el-image :src="imageUrl" fit="cover" style="width:40px;height:40px;border-radius:4px" />
            <el-icon class="remove-image" @click="imageUrl = ''"><Close /></el-icon>
          </span>
          <el-button type="primary" size="small" :loading="submitting" @click="submitComment">发表评论</el-button>
        </div>
      </div>
    </div>
    <div class="login-tip" v-else>
      <el-button type="primary" link @click="goLogin">登录后发表评论</el-button>
    </div>

    <!-- 评论列表 -->
    <div class="comment-list" v-loading="loading">
      <el-empty v-if="!loading && list.length === 0" description="暂无评论，来抢沙发吧" :image-size="60" />

      <div v-for="item in list" :key="item.id" class="comment-item">
        <el-avatar :size="36" :src="item.user?.avatar">
          {{ item.user?.nickname?.charAt(0) || 'U' }}
        </el-avatar>
        <div class="comment-main">
          <div class="comment-head">
            <span class="commenter">{{ item.user?.nickname || '匿名用户' }}</span>
            <span v-if="item.user?.ip_location" class="ip-loc">{{ item.user.ip_location }}</span>
            <span class="comment-time">{{ formatTime(item.created_at) }}</span>
          </div>
          <div class="comment-content">{{ item.content }}</div>
          <el-image
            v-if="item.image_url"
            :src="item.image_url"
            fit="cover"
            class="comment-image"
            :preview-src-list="[item.image_url]"
          />
          <div class="comment-actions">
            <el-button link size="small" :class="{ 'liked': item.is_liked }" :icon="Top" @click="toggleCommentLike(item)">
              {{ item.like_count || 0 }}
            </el-button>
            <el-button link size="small" :icon="ChatDotRound" @click="toggleReplies(item)">
              {{ item.reply_count || 0 }} 回复
            </el-button>
            <el-button
              v-if="canDelete(item)"
              link
              size="small"
              type="danger"
              :icon="Delete"
              @click="deleteCmt(item)"
            >删除</el-button>
          </div>

          <!-- 楼中楼回复 -->
          <div v-if="expandedReplies.has(item.id)" class="replies-box">
            <div v-for="reply in item.replies" :key="reply.id" class="reply-item">
              <el-avatar :size="26" :src="reply.user?.avatar">
                {{ reply.user?.nickname?.charAt(0) || 'U' }}
              </el-avatar>
              <div class="reply-main">
                <div class="reply-head">
                  <span class="replyer">{{ reply.user?.nickname || '匿名用户' }}</span>
                  <span class="comment-time">{{ formatTime(reply.created_at) }}</span>
                </div>
                <div class="reply-content">{{ reply.content }}</div>
                <el-button
                  v-if="canDelete(reply)"
                  link size="small" type="danger" :icon="Delete"
                  @click="deleteReply(item, reply)"
                >删除</el-button>
              </div>
            </div>
            <div v-if="item.replies?.length && item.reply_count > item.replies.length" class="more-replies">
              <el-button link size="small" @click="loadMoreReplies(item)">查看更多回复</el-button>
            </div>
            <div class="reply-compose" v-if="userStore.isLoggedIn">
              <el-input v-model="replyDrafts[item.id]" size="small" placeholder="回复这条评论..." @keyup.enter="submitReply(item)" />
              <el-button size="small" type="primary" @click="submitReply(item)">回复</el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="load-more" v-if="hasMore">
      <el-button link @click="loadMore">加载更多</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Picture, Close, Top, ChatDotRound, Delete
} from '@element-plus/icons-vue'
import {
  listComments, createComment, deleteComment, toggleCommentLike as toggleCommentLikeApi,
  getCommentReplies, uploadCommunityImage
} from '@/api'

const props = defineProps({
  photosetId: { type: [Number, String], required: true }
})

const router = useRouter()
const userStore = useUserStore()

const list = ref([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = 20
const sortBy = ref('newest')
const newComment = ref('')
const imageUrl = ref('')
const submitting = ref(false)
const imageInput = ref(null)
const expandedReplies = reactive(new Set())
const replyDrafts = reactive({})
const loadedReplyPages = reactive({})

const commentCount = computed(() => total.value)
const hasMore = computed(() => list.value.length < total.value)

onMounted(() => loadComments())

function reload() {
  page.value = 1
  list.value = []
  loadComments()
}

async function loadComments() {
  loading.value = true
  try {
    const res = await listComments(props.photosetId, {
      page: page.value, page_size: pageSize
    })
    const data = res?.data || {}
    const items = data.list || []
    if (page.value === 1) {
      list.value = items
    } else {
      list.value = [...list.value, ...items]
    }
    total.value = data.total || 0
    // 初始化已展开回复的页码
    items.forEach(item => {
      if (!(item.id in loadedReplyPages)) loadedReplyPages[item.id] = 1
    })
  } catch (e) {
    console.warn('加载评论失败:', e.message)
  } finally {
    loading.value = false
  }
}

function loadMore() {
  page.value++
  loadComments()
}

async function loadMoreReplies(item) {
  const nextPage = (loadedReplyPages[item.id] || 1) + 1
  try {
    const res = await getCommentReplies(props.photosetId, item.id, {
      page: nextPage, page_size: 10
    })
    const data = res?.data || {}
    item.replies = [...(item.replies || []), ...(data.list || [])]
    loadedReplyPages[item.id] = nextPage
  } catch (e) {
    ElMessage.error('加载回复失败')
  }
}

function toggleReplies(item) {
  if (expandedReplies.has(item.id)) {
    expandedReplies.delete(item.id)
  } else {
    expandedReplies.add(item.id)
  }
}

function openImagePicker() {
  imageInput.value?.click()
}

async function onImagePicked(e) {
  const file = e.target.files?.[0]
  if (!file) return
  try {
    const res = await uploadCommunityImage(file)
    const url = res?.data?.url || res?.data?.image_url || res?.url
    if (url) {
      imageUrl.value = url
      ElMessage.success('配图上传成功')
    } else {
      ElMessage.error('上传返回地址为空')
    }
  } catch (err) {
    ElMessage.error('配图上传失败')
  }
  e.target.value = ''
}

async function submitComment() {
  const content = newComment.value.trim()
  if (!content) {
    ElMessage.warning('请输入评论内容')
    return
  }
  submitting.value = true
  try {
    const res = await createComment(props.photosetId, {
      content,
      image_url: imageUrl.value || undefined
    })
    const created = res?.data
    if (created) {
      list.value.unshift({
        ...created,
        user: created.user || { id: userStore.user?.id, nickname: userStore.user?.nickname, avatar: userStore.avatar },
        replies: []
      })
      total.value++
    }
    newComment.value = ''
    imageUrl.value = ''
    ElMessage.success('评论发表成功')
  } catch (e) {
    ElMessage.error('发表失败')
  } finally {
    submitting.value = false
  }
}

async function submitReply(item) {
  const content = (replyDrafts[item.id] || '').trim()
  if (!content) return
  try {
    await createComment(props.photosetId, { content, parent_id: item.id })
    replyDrafts[item.id] = ''
    item.reply_count = (item.reply_count || 0) + 1
    if (!item.replies) item.replies = []
    item.replies.push({
      id: Date.now(),
      content,
      user: { id: userStore.user?.id, nickname: userStore.user?.nickname, avatar: userStore.avatar },
      created_at: new Date().toISOString()
    })
    expandedReplies.add(item.id)
    ElMessage.success('回复成功')
  } catch (e) {
    ElMessage.error('回复失败')
  }
}

async function toggleCommentLike(item) {
  if (!userStore.isLoggedIn) {
    ElMessage.info('请先登录')
    return
  }
  try {
    await toggleCommentLikeApi(item.id)
    item.is_liked = !item.is_liked
    item.like_count = (item.like_count || 0) + (item.is_liked ? 1 : -1)
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

async function deleteCmt(item) {
  try {
    await ElMessageBox.confirm('确定删除这条评论？', '提示', { type: 'warning' })
    await deleteComment(props.photosetId, item.id)
    list.value = list.value.filter(c => c.id !== item.id)
    total.value--
    ElMessage.success('已删除')
  } catch (e) { /* 取消删除 */ }
}

async function deleteReply(item, reply) {
  try {
    await ElMessageBox.confirm('确定删除这条回复？', '提示', { type: 'warning' })
    await deleteComment(props.photosetId, reply.id)
    item.replies = item.replies.filter(r => r.id !== reply.id)
    item.reply_count = Math.max(0, (item.reply_count || 0) - 1)
    ElMessage.success('已删除')
  } catch (e) { /* 取消删除 */ }
}

function canDelete(item) {
  return userStore.isLoggedIn &&
    (item.user?.id === userStore.user?.id || userStore.user?.role === 'admin')
}

function goLogin() {
  router.push({ name: 'Login', query: { redirect: router.currentRoute.value.fullPath } })
}

function formatTime(str) {
  if (!str) return ''
  const d = new Date(str)
  if (isNaN(d.getTime())) return str
  const diff = Date.now() - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins}分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}天前`
  return d.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.comment-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-top: 24px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}
.section-header h3 { font-size: 18px; margin: 0; }
.count { color: #999; font-size: 14px; font-weight: normal; }
.comment-compose { display: flex; gap: 12px; margin-bottom: 24px; }
.compose-body { flex: 1; }
.compose-actions { display: flex; align-items: center; gap: 8px; margin-top: 8px; }
.image-preview { position: relative; display: inline-flex; }
.remove-image { position: absolute; top: -6px; right: -6px; cursor: pointer; color: #f56c6c; background: #fff; border-radius: 50%; }
.login-tip { padding: 12px 0 20px; }
.comment-list { min-height: 60px; }
.comment-item { display: flex; gap: 12px; padding: 16px 0; border-top: 1px solid #f0f0f0; }
.comment-item:first-child { border-top: none; }
.comment-main { flex: 1; min-width: 0; }
.comment-head { display: flex; align-items: center; gap: 8px; }
.commenter { font-weight: 600; color: #409eff; font-size: 14px; }
.ip-loc { font-size: 11px; color: #bbb; }
.comment-time { font-size: 12px; color: #999; margin-left: auto; }
.comment-content { margin: 6px 0; color: #333; line-height: 1.6; word-break: break-word; }
.comment-image { max-width: 200px; max-height: 150px; border-radius: 6px; margin: 6px 0; }
.comment-actions { display: flex; align-items: center; gap: 12px; margin-top: 4px; }
.comment-actions .liked { color: #f56c6c; }
.replies-box { background: #f8f8f8; border-radius: 8px; padding: 12px; margin-top: 10px; }
.reply-item { display: flex; gap: 8px; padding: 8px 0; }
.reply-main { flex: 1; }
.reply-head { display: flex; align-items: center; gap: 8px; }
.replyer { color: #409eff; font-size: 13px; font-weight: 600; }
.reply-content { font-size: 13px; color: #444; margin: 4px 0; word-break: break-word; }
.more-replies { padding: 4px 0; }
.reply-compose { display: flex; gap: 8px; margin-top: 8px; }
.load-more { text-align: center; padding: 12px; }
</style>
