<template>
  <div class="trash-page">
    <el-card>
      <template #header>
        <div class="header-row">
          <h3>回收站</h3>
          <el-button size="small" @click="$router.push('/my-photosets')">返回我的套图</el-button>
        </div>
      </template>
      <el-alert type="info" :closable="false" show-icon title="回收站中的套图可恢复，永久删除后无法找回" style="margin-bottom:16px" />
      <div v-loading="loading">
        <el-empty v-if="!loading && list.length === 0" description="回收站是空的" :image-size="80" />
        <div v-for="item in list" :key="item.id" class="trash-item">
          <el-image :src="item.cover" fit="cover" class="thumb" />
          <div class="item-info">
            <div class="item-title">{{ item.title }}</div>
            <div class="item-meta">
              <span>删除时间：{{ formatDate(item.deleted_at) }}</span>
              <span>{{ item.photo_count || 0 }} 张图片</span>
            </div>
          </div>
          <el-button type="primary" size="small" @click="restore(item)">恢复</el-button>
          <el-button type="danger" size="small" plain @click="permanentDelete(item)">永久删除</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTrashPhotosets, restorePhotoset, permanentDeletePhotoset } from '@/api'

const list = ref([])
const loading = ref(false)

onMounted(() => load())
async function load() {
  loading.value = true
  try {
    const res = await getTrashPhotosets()
    const data = res?.data
    list.value = Array.isArray(data) ? data : (data?.list || [])
  } catch (e) { console.warn('加载回收站失败:', e.message) }
  finally { loading.value = false }
}

async function restore(item) {
  try {
    await restorePhotoset(item.id)
    list.value = list.value.filter(p => p.id !== item.id)
    ElMessage.success('已恢复')
  } catch (e) { ElMessage.error('恢复失败') }
}

async function permanentDelete(item) {
  try {
    await ElMessageBox.confirm(`确定永久删除「${item.title}」？此操作不可恢复！`, '危险操作', { type: 'warning', confirmButtonText: '永久删除' })
    await permanentDeletePhotoset(item.id)
    list.value = list.value.filter(p => p.id !== item.id)
    ElMessage.success('已永久删除')
  } catch (e) { /* 取消 */ }
}

function formatDate(str) {
  if (!str) return '未知'
  return new Date(str).toLocaleString('zh-CN')
}
</script>

<style scoped>
.trash-page { max-width: 860px; margin: 0 auto; padding: 24px 0 60px; }
.header-row { display: flex; align-items: center; justify-content: space-between; }
.header-row h3 { margin: 0; }
.trash-item { display: flex; align-items: center; gap: 14px; padding: 14px 0; border-bottom: 1px solid #f0f0f0; }
.thumb { width: 80px; height: 56px; border-radius: 6px; flex-shrink: 0; }
.item-info { flex: 1; }
.item-title { font-weight: 600; }
.item-meta { display: flex; gap: 16px; color: #999; font-size: 12px; margin-top: 4px; }
</style>
