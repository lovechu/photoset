<template>
  <div class="post-detail-page" v-loading="loading">
    <!-- 面包屑 -->
    <nav class="crumb">
      <a class="crumb-link" @click="$router.push('/community')"><el-icon :size="13"><ArrowLeft /></el-icon> 社区</a>
      <span class="crumb-sep">›</span>
      <span class="crumb-cat">{{ post?.category || '话题' }}</span>
    </nav>

    <template v-if="post">
      <!-- 主题头 -->
      <header class="topic-header">
        <div class="topic-labels" v-if="post.category || post.topics?.length || post.tags?.length || post.is_pinned || post.is_essence">
          <span v-if="post.category" class="tlabel blue">{{ post.category }}</span>
          <span v-for="t in post.topics" :key="'t'+t.id" class="tlabel orange">#{{ t.name }}</span>
          <span v-for="t in post.tags" :key="'g'+t.id" class="tlabel gray">{{ t.name }}</span>
          <span v-if="post.is_pinned" class="tlabel red">置顶</span>
          <span v-if="post.is_essence" class="tlabel gold">精华</span>
        </div>
        <h1 class="topic-title" v-if="post.title">{{ post.title }}</h1>
        <div class="topic-meta">
          <span class="author">{{ post.author_name || '匿名用户' }}</span>
          <span>· {{ formatTime(post.created_at) }}</span>
          <span>· {{ post.reply_count || 0 }} 回复</span>
          <span>· {{ post.view_count || 0 }} 浏览</span>
        </div>
      </header>

      <!-- 主楼 -->
      <article class="post-floor">
        <div class="floor-author">
          <el-avatar :size="46" :src="post.author_avatar">{{ post.author_name?.charAt(0) || 'U' }}</el-avatar>
          <div class="floor-author-name">{{ post.author_name || '匿名用户' }}</div>
          <div class="floor-author-meta" v-if="post.author_ip_location">IP {{ post.author_ip_location }}</div>
          <button v-if="userStore.isLoggedIn && post.author_id !== userStore.user?.id" class="follow-btn" @click="toggleFollow">
            {{ post.is_following ? '取消关注' : '+ 关注' }}
          </button>
        </div>
        <div class="floor-body">
          <span class="floor-num">#1</span>
          <div class="post-content" v-html="renderedContent"></div>

          <div v-if="videos.length" class="post-videos">
            <video v-for="(v, i) in videos" :key="i" :src="v" controls class="post-video" />
          </div>

          <div class="floor-actions">
            <button :class="['act', { on: post.is_liked }]" @click="toggleLike"><el-icon :size="14"><Top /></el-icon> {{ post.is_liked ? '已赞' : '点赞' }} {{ post.like_count || 0 }}</button>
            <button :class="['act', { on: isFavorited }]" @click="toggleFavorite"><el-icon :size="14"><Star /></el-icon> {{ isFavorited ? '已收藏' : '收藏' }}</button>
            <button class="act" @click="sharePost"><el-icon :size="14"><Share /></el-icon> 分享</button>
            <button class="act" v-if="canEdit" @click="editPost"><el-icon :size="14"><Edit /></el-icon> 编辑</button>
            <button class="act danger" v-if="canEdit" @click="deleteCurrentPost"><el-icon :size="14"><Delete /></el-icon> 删除</button>
            <button class="act" @click="reportDialog = true"><el-icon :size="14"><Warning /></el-icon> 举报</button>
          </div>
        </div>
      </article>

      <!-- 回复区 -->
      <section class="replies-section">
        <h3 class="replies-heading">{{ pagination.total || 0 }} 条回复</h3>

        <div class="reply-compose" v-if="userStore.isLoggedIn">
          <el-avatar :size="40" :src="userStore.avatar || userStore.user?.avatar">{{ userStore.user?.nickname?.charAt(0) || '我' }}</el-avatar>
          <div class="compose-body">
            <el-input v-model="replyContent" type="textarea" :rows="3" maxlength="1000" show-word-limit placeholder="写下你的回复... 支持 Markdown 图片语法" />
            <div v-if="replyImages.length" class="reply-image-previews">
              <div v-for="(url, i) in replyImages" :key="i" class="preview-thumb">
                <img :src="url" :alt="`回复图片${i + 1}`" />
                <span class="preview-remove" @click="removeReplyImage(i)">×</span>
              </div>
            </div>
            <div class="compose-actions">
              <button class="attach-btn" :disabled="uploadingImage || replyImages.length >= 5" @click="imageInput?.click()">
                <el-icon :size="15"><Picture /></el-icon> {{ uploadingImage ? '上传中' : '添加图片' }}
              </button>
              <input ref="imageInput" type="file" accept="image/jpeg,image/png,image/webp" multiple style="display:none" @change="onReplyImagesPicked" />
              <button class="submit-btn" :disabled="submitting" @click="submitReply">{{ submitting ? '发表中' : '发表回复' }}</button>
            </div>
          </div>
        </div>
        <div class="login-tip" v-else>
          <button class="submit-btn" @click="$router.push({ name: 'Login', query: { redirect: $route.fullPath } })">登录后参与回复</button>
        </div>

        <div v-loading="repliesLoading">
          <el-empty v-if="!repliesLoading && replies.length === 0" description="暂无回复" :image-size="60" />
          <div v-for="(reply, i) in replies" :key="reply.id" class="post-floor reply-floor">
            <div class="floor-author">
              <el-avatar :size="40" :src="reply.user?.avatar">{{ reply.user?.nickname?.charAt(0) || 'U' }}</el-avatar>
              <div class="floor-author-name">{{ reply.user?.nickname || '用户' }}</div>
              <div class="floor-author-meta" v-if="reply.author_ip_location">IP {{ reply.author_ip_location }}</div>
            </div>
            <div class="floor-body">
              <div class="reply-floor-head">
                <span class="floor-num">#{{ i + 2 }}</span>
                <span class="reply-time">{{ formatTime(reply.created_at) }}</span>
              </div>
              <div class="reply-content" v-html="renderContent(reply.content)"></div>
              <div class="floor-actions">
                <button :class="['act', { on: reply.is_liked }]" @click="toggleReplyLike(reply)"><el-icon :size="14"><Top /></el-icon> {{ reply.like_count || 0 }}</button>
                <button class="act" @click="setReplyingTo(reply)"><el-icon :size="14"><ChatDotRound /></el-icon> 回复</button>
                <button class="act danger" v-if="canDeleteReply(reply)" @click="deleteReplyItem(reply)"><el-icon :size="14"><Delete /></el-icon> 删除</button>
              </div>
              <div v-if="reply.children?.length" class="children-replies">
                <div v-for="child in reply.children" :key="child.id" class="child-reply">
                  <span class="child-name">{{ child.user?.nickname || '用户' }}：</span>
                  <span v-html="renderContent(child.content)"></span>
                </div>
              </div>
              <div v-if="replyingTo === reply.id" class="inline-reply">
                <el-input v-model="childReplyContent" size="small" placeholder="回复这条回复..." @keyup.enter="submitChildReply(reply)" />
                <button class="attach-btn sm" :disabled="childUploading" @click="pickChildImage(reply.id)"><el-icon :size="14"><Picture /></el-icon> 图</button>
                <button class="submit-btn sm" @click="submitChildReply(reply)">发送</button>
                <button class="act" @click="replyingTo = null">取消</button>
                <input :ref="el => childImageInputs[reply.id] = el" type="file" accept="image/jpeg,image/png,image/webp" style="display:none" @change="e => onChildImagePicked(e, reply)" />
                <div v-if="childReplyImages[reply.id]?.length" class="reply-image-previews">
                  <div v-for="(url, i) in childReplyImages[reply.id]" :key="i" class="preview-thumb">
                    <img :src="url" :alt="`回复图片${i + 1}`" />
                    <span class="preview-remove" @click="removeChildImage(reply.id, i)">×</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="load-more" v-if="replies.length < pagination.total">
            <button class="load-btn" @click="loadMoreReplies">加载更多回复</button>
          </div>
        </div>
      </section>
    </template>

    <el-result v-else-if="!loading" icon="warning" title="帖子不存在" sub-title="可能已被删除">
      <template #extra><el-button type="primary" @click="$router.push('/community')">返回社区</el-button></template>
    </el-result>

    <!-- 举报弹窗 -->
    <el-dialog v-model="reportDialog" title="举报帖子" width="440px">
      <el-form label-width="70px">
        <el-form-item label="原因">
          <el-select v-model="reportReason" placeholder="请选择举报原因" style="width:100%">
            <el-option label="垃圾广告" value="spam" />
            <el-option label="违规内容" value="illegal" />
            <el-option label="恶意攻击" value="abuse" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="reportDetail" type="textarea" :rows="3" placeholder="补充说明（选填）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reportDialog = false">取消</el-button>
        <el-button type="primary" @click="submitReport">提交举报</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Top, Star, Share, Edit, Delete, Warning, ChatDotRound, Picture } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import {
  getPostDetail, getReplies, createReply, deleteReply, togglePostLike,
  toggleReplyLike as toggleReplyLikeApi, togglePostFavorite, checkPostFavorite, reportPost,
  deletePost, followUser, unfollowUser, recordShare, uploadCommunityImage
} from '@/api'
import { renderContent, extractVideos, formatTime } from '@/utils/post'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const postId = route.params.id

const loading = ref(false)
const post = ref(null)
const replies = ref([])
const pagination = ref({ total: 0, page: 1, page_size: 20 })
const repliesLoading = ref(false)
const replyContent = ref('')
const childReplyContent = ref('')
const replyingTo = ref(null)
const submitting = ref(false)
const isFavorited = ref(false)
const reportDialog = ref(false)
const reportReason = ref('spam')
const reportDetail = ref('')

// 回帖图片：以 Markdown ![](url) 内嵌到 content，后端零改动
const imageInput = ref(null)
const replyImages = ref([])
const uploadingImage = ref(false)
const childReplyImages = reactive({})
const childImageInputs = reactive({})
const childUploading = ref(false)

const renderedContent = computed(() => renderContent(post.value?.content))
const videos = computed(() => extractVideos(post.value?.content))
const canEdit = computed(() => userStore.isLoggedIn &&
  (post.value?.author_id === userStore.user?.id || userStore.user?.role === 'admin'))

onMounted(() => {
  loadPost()
  loadReplies()
  if (userStore.isLoggedIn) checkFavorite()
})

async function loadPost() {
  loading.value = true
  try {
    const res = await getPostDetail(postId)
    post.value = res?.data?.post || res?.data || null
  } catch (e) {
    console.warn('加载帖子失败:', e.message)
  } finally {
    loading.value = false
  }
}

async function loadReplies() {
  repliesLoading.value = true
  try {
    const res = await getReplies(postId, { page: pagination.value.page, page_size: 20 })
    const data = res?.data || {}
    const items = data.replies || []
    replies.value = pagination.value.page === 1 ? items : [...replies.value, ...items]
    pagination.value = data.pagination || { total: 0, page: 1, page_size: 20 }
  } catch (e) {
    console.warn('加载回复失败:', e.message)
  } finally {
    repliesLoading.value = false
  }
}

function loadMoreReplies() {
  pagination.value.page++
  loadReplies().then(() => {})
}

async function checkFavorite() {
  try {
    const res = await checkPostFavorite(postId)
    isFavorited.value = !!(res?.data?.is_favorited ?? res?.data?.favorited)
  } catch (e) { /* 静默 */ }
}

async function toggleLike() {
  if (!userStore.isLoggedIn) return ElMessage.info('请先登录')
  try {
    await togglePostLike(postId)
    post.value.is_liked = !post.value.is_liked
    post.value.like_count = (post.value.like_count || 0) + (post.value.is_liked ? 1 : -1)
  } catch (e) { ElMessage.error('操作失败') }
}

async function toggleFavorite() {
  if (!userStore.isLoggedIn) return ElMessage.info('请先登录')
  try {
    await togglePostFavorite(postId)
    isFavorited.value = !isFavorited.value
    ElMessage.success(isFavorited.value ? '已收藏' : '已取消收藏')
  } catch (e) { ElMessage.error('操作失败') }
}

async function toggleFollow() {
  if (!userStore.isLoggedIn) return ElMessage.info('请先登录')
  try {
    if (post.value.is_following) await unfollowUser(post.value.author_id)
    else await followUser(post.value.author_id)
    post.value.is_following = !post.value.is_following
    ElMessage.success(post.value.is_following ? '已关注' : '已取消关注')
  } catch (e) { ElMessage.error('操作失败') }
}

async function sharePost() {
  try { await recordShare(postId, {}) } catch (e) { /* 静默 */ }
  try {
    await navigator.clipboard.writeText(`${window.location.origin}/community/post/${postId}`)
    ElMessage.success('链接已复制')
  } catch (e) { ElMessage.info('分享链接已生成') }
}

function editPost() {
  router.push({ path: '/community/create', query: { edit: postId } })
}

async function deleteCurrentPost() {
  try {
    await ElMessageBox.confirm('确定删除这个帖子？', '提示', { type: 'warning' })
    await deletePost(postId)
    ElMessage.success('已删除')
    router.push('/community')
  } catch (e) { /* 取消 */ }
}

async function submitReport() {
  try {
    await reportPost(postId, { reason: reportReason, detail: reportDetail })
    reportDialog.value = false
    ElMessage.success('举报已提交')
  } catch (e) { ElMessage.error('提交失败') }
}

// ============ 回帖图片（Markdown 内嵌，后端零改动） ============
async function uploadImages(files, target) {
  const urls = []
  for (const file of files) {
    if (file.size > 5 * 1024 * 1024) {
      ElMessage.warning(`${file.name} 超过 5MB，已跳过`)
      continue
    }
    try {
      const res = await uploadCommunityImage(file)
      const url = res?.data?.url || res?.data?.image_url || res?.url
      if (url) urls.push(url)
    } catch (e) {
      ElMessage.error('图片上传失败')
    }
  }
  if (urls.length) {
    target.push(...urls)
    ElMessage.success(`已添加 ${urls.length} 张图片`)
  }
  return urls
}

async function onReplyImagesPicked(e) {
  const files = Array.from(e.target.files || [])
  if (!files.length) return
  uploadingImage.value = true
  try {
    const urls = await uploadImages(files, replyImages.value)
    if (urls.length) {
      replyContent.value += (replyContent.value ? '\n' : '') + urls.map(u => `![](${u})`).join('\n')
    }
  } finally {
    uploadingImage.value = false
    e.target.value = ''
  }
}

function removeReplyImage(i) {
  const url = replyImages.value[i]
  replyImages.value.splice(i, 1)
  replyContent.value = replyContent.value.replace(`![](${url})`, '').trim()
}

function pickChildImage(replyId) {
  childImageInputs[replyId]?.click()
}

async function onChildImagePicked(e, reply) {
  const files = Array.from(e.target.files || [])
  if (!files.length) return
  childUploading.value = true
  try {
    if (!childReplyImages[reply.id]) childReplyImages[reply.id] = []
    const urls = await uploadImages(files, childReplyImages[reply.id])
    if (urls.length) {
      childReplyContent.value += (childReplyContent.value ? '\n' : '') + urls.map(u => `![](${u})`).join('\n')
    }
  } finally {
    childUploading.value = false
    e.target.value = ''
  }
}

function removeChildImage(replyId, i) {
  const url = childReplyImages[replyId]?.[i]
  if (childReplyImages[replyId]) childReplyImages[replyId].splice(i, 1)
  childReplyContent.value = childReplyContent.value.replace(`![](${url})`, '').trim()
}

async function submitReply() {
  const content = replyContent.value.trim()
  if (!content && !replyImages.value.length) return ElMessage.warning('请输入回复内容')
  submitting.value = true
  try {
    const res = await createReply(postId, { content })
    const reply = res?.data?.reply
    if (reply) replies.value.unshift({ ...reply, user: reply.user || { id: userStore.user?.id, nickname: userStore.user?.nickname, avatar: userStore.avatar } })
    replyContent.value = ''
    replyImages.value = []
    post.value.reply_count = (post.value.reply_count || 0) + 1
    pagination.value.total = (pagination.value.total || 0) + 1
    ElMessage.success('回复成功')
  } catch (e) { ElMessage.error('回复失败') }
  finally { submitting.value = false }
}

function setReplyingTo(reply) {
  replyingTo.value = reply.id
  childReplyContent.value = ''
  childReplyImages[reply.id] = childReplyImages[reply.id] || []
}

async function submitChildReply(parent) {
  const content = childReplyContent.value.trim()
  if (!content && !(childReplyImages[parent.id]?.length)) return
  try {
    const res = await createReply(postId, { content, parent_reply_id: parent.id })
    const child = res?.data?.reply
    if (!parent.children) parent.children = []
    parent.children.push(child || { content, user: { nickname: userStore.user?.nickname } })
    childReplyContent.value = ''
    childReplyImages[parent.id] = []
    replyingTo.value = null
    ElMessage.success('回复成功')
  } catch (e) { ElMessage.error('回复失败') }
}

async function toggleReplyLike(reply) {
  if (!userStore.isLoggedIn) return ElMessage.info('请先登录')
  try {
    await toggleReplyLikeApi(reply.id)
    reply.is_liked = !reply.is_liked
    reply.like_count = (reply.like_count || 0) + (reply.is_liked ? 1 : -1)
  } catch (e) { ElMessage.error('操作失败') }
}

async function deleteReplyItem(reply) {
  try {
    await ElMessageBox.confirm('确定删除这条回复？', '提示', { type: 'warning' })
    await deleteReply(reply.id)
    replies.value = replies.value.filter(r => r.id !== reply.id)
    ElMessage.success('已删除')
  } catch (e) { /* 取消 */ }
}

function canDeleteReply(reply) {
  return userStore.isLoggedIn && (reply.user_id === userStore.user?.id || userStore.user?.role === 'admin')
}
</script>

<style scoped>
.post-detail-page {
  max-width: 1080px; margin: 0 auto; padding: 0 20px 80px;
  color: #33373d; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  font-size: 14px; line-height: 1.5;
}

/* 面包屑 */
.crumb { display: flex; align-items: center; gap: 8px; padding: 18px 0 14px; font-size: 12px; color: #9aa1ac; }
.crumb-link { display: inline-flex; align-items: center; gap: 4px; color: #6b7280; cursor: pointer; }
.crumb-link:hover { color: #4f6ef7; }
.crumb-sep { color: #c6cbd2; }
.crumb-cat { color: #6b7280; }

/* 主题头 */
.topic-header {
  background: #fff; border: 1px solid #e6e8eb; border-radius: 5px 5px 0 0;
  padding: 20px 24px 18px; border-bottom: 1px solid #eef0f2;
}
.topic-labels { display: flex; gap: 6px; flex-wrap: wrap; margin-bottom: 10px; }
.tlabel { padding: 1px 8px; border-radius: 3px; font-size: 10px; font-weight: 700; color: #fff; }
.tlabel.blue { background: #4f78d8; }
.tlabel.orange { background: #d9902b; }
.tlabel.gray { background: #9aa1ac; }
.tlabel.red { background: #d04545; }
.tlabel.gold { background: #c9971f; }
.topic-title { font-size: 24px; font-weight: 700; line-height: 1.4; color: #2b3038; margin: 0 0 8px; }
.topic-meta { display: flex; align-items: center; gap: 6px; font-size: 12px; color: #9aa1ac; flex-wrap: wrap; }
.topic-meta .author { color: #6b7280; font-weight: 600; }

/* 楼层：左作者 + 右正文 */
.post-floor {
  display: flex; align-items: stretch;
  background: #fff; border: 1px solid #e6e8eb;
}
.topic-header + .post-floor { border-radius: 0 0 5px 5px; border-top: none; }
.floor-author {
  width: 160px; flex-shrink: 0;
  padding: 20px 16px;
  background: #fafbfc;
  border-right: 1px solid #eef0f2;
  display: flex; flex-direction: column; align-items: center; gap: 6px;
  text-align: center;
}
.floor-author :deep(.el-avatar) { background: #eef0f3; color: #6b7280; margin-bottom: 2px; }
.floor-author-name { font-size: 13px; font-weight: 600; color: #3d444e; word-break: break-all; }
.floor-author-meta { font-size: 11px; color: #9aa1ac; }
.follow-btn {
  margin-top: 6px; padding: 4px 12px; font-size: 11px; cursor: pointer;
  background: #fff; border: 1px solid #d5dae1; border-radius: 4px; color: #6b7280;
}
.follow-btn:hover { border-color: #4f6ef7; color: #4f6ef7; }
.floor-body { flex: 1; min-width: 0; padding: 20px 24px; position: relative; }
.floor-num { position: absolute; top: 14px; right: 18px; font-size: 12px; color: #b6bcc5; font-weight: 600; }
.post-content { font-size: 15px; line-height: 1.9; color: #33373d; word-break: break-word; }
.post-content :deep(.post-content-img) { max-width: 100%; border-radius: 4px; margin: 12px 0; }
.post-videos { margin-top: 16px; }
.post-video { width: 100%; max-height: 440px; border-radius: 4px; margin-bottom: 8px; }

/* 文字式操作栏 */
.floor-actions { display: flex; align-items: center; flex-wrap: wrap; gap: 4px; margin-top: 22px; padding-top: 14px; border-top: 1px solid #eef0f2; }
.act {
  display: inline-flex; align-items: center; gap: 4px;
  background: none; border: none; cursor: pointer;
  padding: 5px 9px; border-radius: 3px;
  font-size: 12px; color: #8d95a0;
  transition: background 0.15s, color 0.15s;
}
.act:hover { background: #f2f5f9; color: #3d444e; }
.act.on { color: #4f6ef7; }
.act.danger:hover { color: #d04545; background: #fdeeee; }

/* 回复区 */
.replies-section { padding: 28px 0 0; }
.replies-heading { margin: 0 0 16px; font-size: 15px; font-weight: 700; color: #2b3038; }
.reply-compose {
  display: flex; gap: 14px; margin-bottom: 20px;
  padding: 18px; background: #fff; border: 1px solid #e6e8eb; border-radius: 5px;
}
.reply-compose :deep(.el-avatar) { flex-shrink: 0; background: #eef0f3; color: #6b7280; }
.compose-body { flex: 1; min-width: 0; }
.reply-compose :deep(.el-textarea__inner) { border-radius: 4px; border-color: #e0e3e7; font-size: 14px; }
.reply-compose :deep(.el-textarea__inner:focus) { border-color: #4f6ef7; box-shadow: none; }
.compose-actions { display: flex; justify-content: space-between; align-items: center; margin-top: 10px; }
.attach-btn {
  display: inline-flex; align-items: center; gap: 5px; cursor: pointer;
  padding: 6px 10px; background: #fff; border: 1px solid #e0e3e7; border-radius: 4px;
  font-size: 12px; color: #6b7280;
}
.attach-btn:hover:not(:disabled) { border-color: #4f6ef7; color: #4f6ef7; }
.attach-btn:disabled { opacity: 0.6; cursor: default; }
.submit-btn {
  padding: 7px 18px; background: #f5a623; color: #fff; border: none; border-radius: 4px;
  font-size: 12px; font-weight: 600; cursor: pointer;
}
.submit-btn:hover:not(:disabled) { background: #e8991c; }
.submit-btn:disabled { opacity: 0.6; cursor: default; }
.submit-btn.sm { padding: 5px 14px; }
.attach-btn.sm { padding: 4px 9px; }
.reply-image-previews { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 10px; }
.preview-thumb { position: relative; width: 72px; height: 72px; }
.preview-thumb img { width: 100%; height: 100%; object-fit: cover; border-radius: 4px; border: 1px solid #e5e8eb; display: block; }
.preview-remove {
  position: absolute; top: -7px; right: -7px;
  width: 18px; height: 18px; border-radius: 50%;
  background: #2b3038; color: #fff;
  font-size: 13px; line-height: 18px; text-align: center;
  cursor: pointer;
}
.reply-content :deep(.post-content-img) { max-width: 260px; border-radius: 4px; margin: 6px 0; display: block; }
.child-reply :deep(.post-content-img) { max-width: 180px; border-radius: 4px; margin: 4px 0; display: block; }
.login-tip { padding: 4px 0 20px; }

.reply-floor { border-radius: 5px; margin-bottom: 10px; }
.reply-floor .floor-author { width: 140px; background: #fff; }
.reply-floor .floor-body { padding: 16px 22px; }
.reply-floor-head { display: flex; align-items: center; gap: 10px; margin-bottom: 8px; }
.reply-floor-head .floor-num { position: static; color: #6b7280; }
.reply-time { font-size: 11px; color: #9aa1ac; }
.reply-content { color: #444; font-size: 14px; line-height: 1.7; word-break: break-word; }
.children-replies { background: #f7f8f9; border-left: 3px solid #4f6ef7; border-radius: 0 4px 4px 0; padding: 8px 14px; margin-top: 12px; }
.child-reply { font-size: 13px; color: #555; padding: 4px 0; line-height: 1.6; }
.child-name { color: #2b3038; font-weight: 600; }
.inline-reply { display: flex; align-items: center; gap: 8px; margin-top: 12px; flex-wrap: wrap; }
.inline-reply :deep(.el-input) { flex: 1; min-width: 160px; }
.load-more { text-align: center; padding: 20px 0; }
.load-btn {
  padding: 8px 24px; border: 1px solid #e0e3e7; background: #fff; color: #6b7280;
  border-radius: 4px; font-size: 12px; cursor: pointer;
}
.load-btn:hover { border-color: #4f6ef7; color: #4f6ef7; }

@media (max-width: 768px) {
  .post-detail-page { padding: 0 10px 60px; }
  .topic-header { padding: 16px; }
  .topic-title { font-size: 20px; }
  .post-floor { flex-direction: column; }
  .floor-author, .reply-floor .floor-author {
    width: 100%; flex-direction: row; justify-content: flex-start; text-align: left;
    padding: 12px 16px; border-right: none; border-bottom: 1px solid #eef0f2;
  }
  .floor-author :deep(.el-avatar) { margin-bottom: 0; }
  .floor-body, .reply-floor .floor-body { padding: 16px; }
  .floor-num { top: 10px; right: 12px; }
  .reply-compose { padding: 12px; }
}
</style>
