<template>
  <div class="history-page">
    <div class="page-header">
      <h2>浏览历史</h2>
      <el-button v-if="list.length" type="danger" plain :icon="Delete" @click="clearAll">清空历史</el-button>
    </div>

    <div v-loading="loading">
      <el-empty v-if="!loading && list.length === 0" description="暂无浏览记录" :image-size="100" />
      <div v-else class="history-list">
        <div v-for="item in list" :key="item.id" class="history-item">
          <el-image :src="item.cover_image" fit="cover" class="item-cover" @click="goDetail(item)" />
          <div class="item-info" @click="goDetail(item)">
            <div class="item-title">
              {{ item.title || '未命名套图' }}
              <el-tag v-if="item.is_vip" size="small" type="warning">VIP</el-tag>
            </div>
            <div class="item-meta">
              <span>{{ item.photo_count || 0 }} 张图片</span>
              <span v-if="item.price">¥{{ Number(item.price).toFixed(2) }}</span>
              <span class="view-time">{{ formatTime(item.viewed_at || item.created_at) }}</span>
            </div>
          </div>
          <el-button link type="danger" :icon="Close" @click="removeItem(item)">删除</el-button>
        </div>
      </div>

      <div class="load-more" v-if="list.length < total">
        <el-button @click="loadMore" :loading="loadingMore">加载更多</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Close } from '@element-plus/icons-vue'
import { listViewHistory, deleteViewHistory, clearViewHistory } from '@/api'

const router = useRouter()
const list = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const loadingMore = ref(false)

onMounted(() => loadHistory())

async function loadHistory() {
  loading.value = true
  try {
    const res = await listViewHistory({ page: page.value, page_size: 20 })
    const data = res?.data || {}
    const items = data.list || []
    list.value = page.value === 1 ? items : [...list.value, ...items]
    total.value = data.total || 0
  } catch (e) {
    console.warn('加载浏览历史失败:', e.message)
  } finally {
    loading.value = false
  }
}

function loadMore() {
  page.value++
  loadHistory()
}

async function removeItem(item) {
  try {
    await deleteViewHistory(item.id)
    list.value = list.value.filter(h => h.id !== item.id)
    total.value--
    ElMessage.success('已删除')
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

async function clearAll() {
  try {
    await ElMessageBox.confirm('确定清空全部浏览历史？', '提示', { type: 'warning' })
    await clearViewHistory()
    list.value = []
    total.value = 0
    ElMessage.success('已清空')
  } catch (e) { /* 取消 */ }
}

function goDetail(item) {
  router.push(`/photoset/${item.photoset_id}`)
}

function formatTime(str) {
  if (!str) return ''
  const d = new Date(str.replace(' ', 'T'))
  if (isNaN(d.getTime())) return str
  const diff = Date.now() - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins}分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}小时前`
  return d.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.history-page { max-width: 900px; margin: 0 auto; padding: 24px 0 60px; }
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.page-header h2 { margin: 0; }
.history-item { display: flex; align-items: center; gap: 16px; background: #fff; padding: 12px; border-radius: 8px; margin-bottom: 12px; }
.item-cover { width: 100px; height: 70px; border-radius: 6px; cursor: pointer; flex-shrink: 0; }
.item-info { flex: 1; cursor: pointer; min-width: 0; }
.item-title { font-size: 15px; font-weight: 600; display: flex; align-items: center; gap: 8px; }
.item-meta { display: flex; gap: 16px; color: #999; font-size: 12px; margin-top: 6px; }
.view-time { margin-left: auto; }
.load-more { text-align: center; padding: 20px; }
</style>
