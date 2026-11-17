<template>
  <div class="review-section">
    <div class="section-header">
      <h3>套图评价</h3>
      <el-button v-if="userStore.isLoggedIn && !myReview" type="primary" size="small" @click="openDialog()">我要评价</el-button>
      <el-button v-else-if="userStore.isLoggedIn && myReview" size="small" @click="openDialog(myReview)">编辑我的评价</el-button>
    </div>

    <!-- 评分摘要 -->
    <div v-if="summary && summary.total_count > 0" class="review-summary">
      <div class="score-block">
        <div class="score-number">{{ averageDisplay }}</div>
        <el-rate :model-value="summary.average_rating" disabled show-score text-color="#ff9900" />
        <div class="score-count">{{ summary.total_count }} 人评价</div>
      </div>
      <div class="score-bars">
        <div v-for="star in 5" :key="star" class="score-bar-row">
          <span class="star-label">{{ star }}星</span>
          <el-progress
            :percentage="starPercent(star)"
            :stroke-width="8"
            :show-text="false"
            :color="star >= 4 ? '#67c23a' : star === 3 ? '#e6a23c' : '#f56c6c'"
          />
          <span class="star-count">{{ summary[`rating_${star}_count`] || 0 }}</span>
        </div>
      </div>
      <div v-if="summary.top_tags?.length" class="top-tags">
        <el-tag v-for="t in summary.top_tags" :key="t.tag" size="small" effect="plain">
          {{ t.tag }} ({{ t.count }})
        </el-tag>
      </div>
    </div>
    <el-empty v-else description="暂无评价" :image-size="60" />

    <!-- 评价列表 -->
    <div class="review-list" v-loading="loading">
      <div v-for="item in list" :key="item.id" class="review-item">
        <el-avatar :size="36" :src="item.user?.avatar">
          {{ item.is_anonymous ? '匿' : (item.user?.nickname?.charAt(0) || 'U') }}
        </el-avatar>
        <div class="review-main">
          <div class="review-head">
            <span class="reviewer">{{ item.is_anonymous ? '匿名用户' : (item.user?.nickname || '用户') }}</span>
            <el-rate :model-value="item.rating" disabled size="small" />
            <span class="review-time">{{ formatTime(item.created_at) }}</span>
          </div>
          <div v-if="item.tags" class="review-tags">
            <el-tag v-for="tag in item.tags.split(',').filter(Boolean)" :key="tag" size="small" type="info" effect="plain">{{ tag }}</el-tag>
          </div>
          <div class="review-content">{{ item.content }}</div>
          <div v-if="canEdit(item)" class="review-actions">
            <el-button link size="small" :icon="Edit" @click="openDialog(item)">编辑</el-button>
            <el-button link size="small" type="danger" :icon="Delete" @click="removeReview(item)">删除</el-button>
          </div>
        </div>
      </div>
      <div class="load-more" v-if="list.length < total">
        <el-button link @click="loadMore">加载更多</el-button>
      </div>
    </div>

    <!-- 评价编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑评价' : '发表评价'" width="480px">
      <div class="rating-input">
        <span>总体评分：</span>
        <el-rate v-model="form.rating" :colors="['#99A9BF', '#F7BA2A', '#FF9900']" />
      </div>
      <div class="form-item">
        <span>标签（逗号分隔）：</span>
        <el-input v-model="form.tags" placeholder="如：画面精美,内容优质" />
      </div>
      <div class="form-item">
        <span>评价内容：</span>
        <el-input v-model="form.content" type="textarea" :rows="3" maxlength="500" show-word-limit placeholder="分享你的看法..." />
      </div>
      <div class="form-item">
        <el-checkbox v-model="form.is_anonymous">匿名评价</el-checkbox>
      </div>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitReview">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Delete } from '@element-plus/icons-vue'
import {
  listReviews, getReviewSummary, getMyReview, createReview, updateReview, deleteReview
} from '@/api'

const props = defineProps({
  photosetId: { type: [Number, String], required: true }
})

const userStore = useUserStore()
const list = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const summary = ref(null)
const myReview = ref(null)
const sortBy = ref('newest')

const dialogVisible = ref(false)
const editingId = ref(null)
const submitting = ref(false)
const form = ref({ rating: 5, tags: '', content: '', is_anonymous: false })

const averageDisplay = computed(() => {
  if (!summary.value) return '0.0'
  return Number(summary.value.average_rating || 0).toFixed(1)
})

onMounted(() => loadAll())

function loadAll() {
  loadSummary()
  loadMyReview()
  reload()
}

async function loadSummary() {
  try {
    const res = await getReviewSummary(props.photosetId)
    summary.value = res?.data || null
  } catch (e) {
    console.warn('加载评价摘要失败:', e.message)
  }
}

async function loadMyReview() {
  if (!userStore.isLoggedIn) return
  try {
    const res = await getMyReview(props.photosetId)
    myReview.value = res?.data || null
  } catch (e) {
    myReview.value = null
  }
}

async function reload() {
  page.value = 1
  list.value = []
  await loadReviews()
}

async function loadReviews() {
  loading.value = true
  try {
    const res = await listReviews(props.photosetId, {
      page: page.value, page_size: 10, sort_by: sortBy.value
    })
    const data = res?.data || {}
    const items = data.list || []
    list.value = page.value === 1 ? items : [...list.value, ...items]
    total.value = data.total || 0
  } catch (e) {
    console.warn('加载评价失败:', e.message)
  } finally {
    loading.value = false
  }
}

function loadMore() {
  page.value++
  loadReviews()
}

function starPercent(star) {
  if (!summary.value || !summary.value.total_count) return 0
  return Math.round((summary.value[`rating_${star}_count`] || 0) / summary.value.total_count * 100)
}

function openDialog(review) {
  if (review) {
    editingId.value = review.id
    form.value = {
      rating: review.rating,
      tags: review.tags || '',
      content: review.content || '',
      is_anonymous: !!review.is_anonymous
    }
  } else {
    editingId.value = null
    form.value = { rating: 5, tags: '', content: '', is_anonymous: false }
  }
  dialogVisible.value = true
}

async function submitReview() {
  if (!form.value.rating) {
    ElMessage.warning('请选择评分')
    return
  }
  submitting.value = true
  try {
    const payload = { ...form.value }
    if (editingId.value) {
      await updateReview(props.photosetId, editingId.value, payload)
      ElMessage.success('评价已更新')
    } else {
      await createReview(props.photosetId, payload)
      ElMessage.success('评价发表成功')
    }
    dialogVisible.value = false
    loadAll()
  } catch (e) {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

async function removeReview(item) {
  try {
    await ElMessageBox.confirm('确定删除这条评价？', '提示', { type: 'warning' })
    await deleteReview(props.photosetId, item.id)
    list.value = list.value.filter(r => r.id !== item.id)
    total.value--
    if (myReview.value?.id === item.id) myReview.value = null
    loadSummary()
    ElMessage.success('已删除')
  } catch (e) { /* 取消 */ }
}

function canEdit(item) {
  return userStore.isLoggedIn && item.user_id === userStore.user?.id
}

function formatTime(str) {
  if (!str) return ''
  const d = new Date(str)
  if (isNaN(d.getTime())) return str
  return d.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.review-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-top: 24px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}
.section-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.section-header h3 { font-size: 18px; margin: 0; }
.review-summary { display: flex; gap: 32px; flex-wrap: wrap; align-items: center; background: #fafafa; padding: 20px; border-radius: 8px; margin-bottom: 20px; }
.score-block { text-align: center; }
.score-number { font-size: 40px; font-weight: 700; color: #ff9900; line-height: 1; }
.score-count { color: #999; font-size: 12px; margin-top: 6px; }
.score-bars { flex: 1; min-width: 200px; }
.score-bar-row { display: flex; align-items: center; gap: 8px; margin: 4px 0; }
.star-label { font-size: 12px; color: #666; width: 28px; }
.star-count { font-size: 12px; color: #999; width: 30px; text-align: right; }
.top-tags { display: flex; flex-wrap: wrap; gap: 6px; }
.review-item { display: flex; gap: 12px; padding: 16px 0; border-top: 1px solid #f0f0f0; }
.review-item:first-child { border-top: none; }
.review-main { flex: 1; }
.review-head { display: flex; align-items: center; gap: 10px; }
.reviewer { font-weight: 600; color: #409eff; font-size: 14px; }
.review-time { font-size: 12px; color: #999; margin-left: auto; }
.review-tags { display: flex; flex-wrap: wrap; gap: 6px; margin: 6px 0; }
.review-content { color: #333; line-height: 1.6; margin: 6px 0; }
.review-actions { margin-top: 4px; }
.rating-input, .form-item { margin: 16px 0; }
.rating-input > span, .form-item > span { display: block; margin-bottom: 8px; color: #606266; font-size: 14px; }
.load-more { text-align: center; padding: 12px; }
</style>
