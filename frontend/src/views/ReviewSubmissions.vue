<template>
  <div class="review-submissions-page">
    <el-card>
      <template #header>
        <div class="header-row">
          <h3>投稿审核</h3>
          <el-radio-group v-model="status" size="small" @change="reload">
            <el-radio-button value="pending">待审核</el-radio-button>
            <el-radio-button value="published">已通过</el-radio-button>
            <el-radio-button value="rejected">已拒绝</el-radio-button>
          </el-radio-group>
        </div>
      </template>

      <template v-if="!userStore.isAdmin">
        <el-alert type="info" :closable="false" show-icon title="以下为你提交的投稿状态，审核操作仅限管理员" style="margin-bottom:16px" />
      </template>

      <div v-loading="loading">
        <el-empty v-if="!loading && list.length === 0" :description="`暂无${statusLabel}的投稿`" :image-size="80" />
        <div v-for="item in list" :key="item.id" class="submission-item">
          <el-image :src="item.cover" fit="cover" class="thumb" />
          <div class="item-info">
            <div class="item-title">{{ item.title }}</div>
            <div class="item-meta">
              <span v-if="userStore.isAdmin && item.user">投稿人：{{ item.user.nickname }}</span>
              <span>{{ formatDate(item.created_at) }}</span>
              <el-tag size="small" :type="statusType(item.status)">{{ statusLabel(item.status) }}</el-tag>
            </div>
          </div>
          <template v-if="userStore.isAdmin && item.status === 'pending'">
            <el-button type="success" size="small" @click="approve(item)">通过</el-button>
            <el-button type="danger" size="small" plain @click="reject(item)">拒绝</el-button>
          </template>
          <el-button v-else size="small" @click="$router.push(`/photoset/${item.id}`)">查看</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import {
  getPhotosetListAdvanced, getAdminPhotosets, approvePhotoset, rejectPhotoset
} from '@/api'

const userStore = useUserStore()
const list = ref([])
const status = ref('pending')
const loading = ref(false)

onMounted(() => load())
function reload() { load() }

async function load() {
  loading.value = true
  try {
    if (userStore.isAdmin) {
      const res = await getAdminPhotosets({ status: status.value })
      const data = res?.data
      list.value = Array.isArray(data) ? data : (data?.list || [])
    } else {
      const res = await getPhotosetListAdvanced({ only_mine: true, status: status.value, page: 1, page_size: 50 })
      const data = res?.data || {}
      list.value = data.list || data.photosets || []
    }
  } catch (e) { console.warn('加载投稿列表失败:', e.message); list.value = [] }
  finally { loading.value = false }
}

async function approve(item) {
  try {
    await approvePhotoset(item.id)
    ElMessage.success('已通过')
    load()
  } catch (e) { ElMessage.error('操作失败') }
}

async function reject(item) {
  try {
    const { value } = await ElMessageBox.prompt('请输入拒绝原因', '拒绝投稿', {
      inputPlaceholder: '请说明拒绝原因', inputType: 'textarea'
    })
    await rejectPhotoset(item.id, { reason: value })
    ElMessage.success('已拒绝')
    load()
  } catch (e) { /* 取消 */ }
}

function statusLabel(s) {
  return { pending: '待审核', published: '已通过', rejected: '已拒绝', draft: '草稿' }[s] || s || '未知'
}
function statusType(s) {
  return { pending: 'warning', published: 'success', rejected: 'danger' }[s] || 'info'
}
function formatDate(str) {
  if (!str) return ''
  return new Date(str).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.review-submissions-page { max-width: 900px; margin: 0 auto; padding: 24px 0 60px; }
.header-row { display: flex; align-items: center; justify-content: space-between; }
.header-row h3 { margin: 0; }
.submission-item { display: flex; align-items: center; gap: 14px; padding: 14px 0; border-bottom: 1px solid #f0f0f0; }
.thumb { width: 76px; height: 54px; border-radius: 6px; flex-shrink: 0; }
.item-info { flex: 1; }
.item-title { font-weight: 600; }
.item-meta { display: flex; align-items: center; gap: 12px; color: #999; font-size: 12px; margin-top: 4px; }
</style>
