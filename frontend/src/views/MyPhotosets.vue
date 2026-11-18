<template>
  <div class="my-photosets-page">
    <div class="page-header">
      <div>
        <h2>我的套图</h2>
        <p class="subtitle">管理你发布的套图作品</p>
      </div>
      <div>
        <el-button @click="$router.push('/trash')">回收站</el-button>
        <el-button type="primary" :icon="Plus" @click="$router.push('/create')">创建套图</el-button>
      </div>
    </div>

    <el-tabs v-model="activeStatus" @tab-change="reload">
      <el-tab-pane label="全部" name="" />
      <el-tab-pane label="已发布" name="published" />
      <el-tab-pane label="待审核" name="pending" />
      <el-tab-pane label="草稿" name="draft" />
    </el-tabs>

    <div v-loading="loading">
      <el-empty v-if="!loading && list.length === 0" description="还没有套图，去创建第一个吧" :image-size="80">
        <el-button type="primary" size="small" @click="$router.push('/create')">创建套图</el-button>
      </el-empty>
      <div class="photoset-grid">
        <el-card v-for="item in list" :key="item.id" class="photoset-card" shadow="hover">
          <el-image :src="item.cover" fit="cover" class="cover" @click="$router.push(`/photoset/${item.id}`)" />
          <div class="card-body">
            <div class="card-title" @click="$router.push(`/photoset/${item.id}`)">{{ item.title }}</div>
            <div class="card-meta">
              <el-tag size="small" :type="statusType(item.status)">{{ statusLabel(item.status) }}</el-tag>
              <span :class="['price', { free: item.is_free }]">{{ item.is_free ? '免费' : `¥${Number(item.price||0).toFixed(2)}` }}</span>
            </div>
            <div class="card-stats">
              <span>{{ item.photo_count || 0 }} 图</span>
              <span>{{ item.view_count || 0 }} 浏览</span>
            </div>
            <div class="card-actions">
              <el-button link size="small" :icon="Edit" @click="$router.push(`/photoset/${item.id}/edit`)">编辑</el-button>
              <el-button link size="small" type="danger" :icon="Delete" @click="deletePhotoset(item)">删除</el-button>
            </div>
          </div>
        </el-card>
      </div>
      <div class="load-more" v-if="list.length < total">
        <el-button @click="loadMore">加载更多</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import { getPhotosetListAdvanced, deletePhotoset as apiDeletePhotoset } from '@/api'

const list = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const activeStatus = ref('')

onMounted(() => load())
function reload() { page.value = 1; list.value = []; load() }

async function load() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: 20, only_mine: true }
    if (activeStatus.value) params.status = activeStatus.value
    const res = await getPhotosetListAdvanced(params)
    const data = res?.data || {}
    const items = data.list || data.photosets || []
    list.value = page.value === 1 ? items : [...list.value, ...items]
    total.value = data.total || 0
  } catch (e) { console.warn('加载我的套图失败:', e.message) }
  finally { loading.value = false }
}

function loadMore() { page.value++; load() }

async function deletePhotoset(item) {
  try {
    await ElMessageBox.confirm(`确定删除套图「${item.title}」？将移入回收站。`, '提示', { type: 'warning' })
    await apiDeletePhotoset(item.id)
    list.value = list.value.filter(p => p.id !== item.id)
    total.value--
    ElMessage.success('已移入回收站')
  } catch (e) { /* 取消 */ }
}

function statusLabel(s) {
  return { published: '已发布', pending: '待审核', draft: '草稿', rejected: '已拒绝' }[s] || s || '未知'
}
function statusType(s) {
  return { published: 'success', pending: 'warning', draft: 'info', rejected: 'danger' }[s] || 'info'
}
</script>

<style scoped>
.my-photosets-page { max-width: 1100px; margin: 0 auto; padding: 24px 0 60px; }
.page-header { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 16px; }
.page-header h2 { margin: 0 0 4px; }
.subtitle { color: #999; margin: 0; font-size: 13px; }
.photoset-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(230px, 1fr)); gap: 16px; }
.photoset-card { overflow: hidden; }
.cover { width: 100%; height: 150px; cursor: pointer; }
.card-title { font-weight: 600; cursor: pointer; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.card-title:hover { color: #409eff; }
.card-meta { display: flex; align-items: center; justify-content: space-between; margin-top: 8px; }
.price { color: #f56c6c; font-weight: 600; font-size: 13px; }
.price.free { color: #67c23a; }
.card-stats { display: flex; gap: 12px; color: #999; font-size: 12px; margin-top: 6px; }
.card-actions { display: flex; gap: 12px; margin-top: 8px; padding-top: 8px; border-top: 1px solid #f0f0f0; }
.load-more { text-align: center; padding: 20px; }
</style>
