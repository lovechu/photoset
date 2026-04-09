<template>
  <div class="content-review">
    <div class="action-bar">
      <el-radio-group v-model="currentStatus" @change="fetchList">
        <el-radio-button value="">全部</el-radio-button>
        <el-radio-button value="pending">待审核</el-radio-button>
        <el-radio-button value="published">已通过</el-radio-button>
        <el-radio-button value="draft">已拒绝</el-radio-button>
      </el-radio-group>

      <div class="batch-actions" v-if="selectedIds.length > 0">
        <el-text type="primary" size="small">已选择 {{ selectedIds.length }} 条</el-text>
        <el-button type="success" size="small" @click="handleBatchApprove" :loading="batchProcessing">
          批量通过
        </el-button>
        <el-button type="danger" size="small" @click="openBatchRejectDialog" :loading="batchProcessing">
          批量拒绝
        </el-button>
        <el-button type="danger" size="small" plain @click="handleBatchDelete" :loading="batchProcessing">
          批量删除
        </el-button>
        <el-button type="info" size="small" plain @click="clearSelection">
          取消选择
        </el-button>
      </div>
    </div>

    <div v-loading="loading">
      <el-row :gutter="20" v-if="list.length > 0">
        <el-col :span="6" v-for="item in list" :key="item.id" style="margin-bottom: 20px">
          <el-checkbox v-model="selectedIds" :label="item.id" class="select-checkbox" />
          <el-card shadow="hover" class="review-card" :body-style="{ padding: '0' }">
            <img :src="item.cover" class="review-card__cover" />
            <div class="review-card__body">
              <div class="review-card__title">{{ item.title }}</div>
              <div class="review-card__meta">
                <span>{{ item.user?.nickname || '未知' }}</span>
                <el-tag :type="statusTagType(item.status)" size="small">{{ statusLabel(item.status) }}</el-tag>
              </div>
              <div class="review-card__time">{{ formatTime(item.created_at) }}</div>
              <div class="review-card__actions" v-if="item.status === 'pending'">
                <el-button type="success" size="small" @click="handleApprove(item.id)">通过</el-button>
                <el-button type="danger" size="small" @click="openRejectDialog(item.id)">拒绝</el-button>
              </div>
              <div class="review-card__actions">
                <el-button
                  size="small"
                  plain
                  @click="$router.push(`/photoset/${item.id}/edit`)"
                >
                  编辑
                </el-button>
                <el-button
                  size="small"
                  type="danger"
                  plain
                  @click="handleDelete(item.id)"
                  :loading="deletingId === item.id"
                >
                  删除
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-empty v-else description="暂无数据" />
    </div>

    <el-dialog v-model="rejectDialogVisible" title="拒绝原因" width="460px">
      <el-input
        v-model="rejectReason"
        type="textarea"
        :rows="4"
        placeholder="请输入拒绝原因..."
      />
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="danger" :loading="rejectLoading" @click="handleReject">确认拒绝</el-button>
      </template>
    </el-dialog>

    <!-- 批量拒绝对话框 -->
    <el-dialog v-model="batchRejectDialogVisible" title="批量拒绝原因" width="460px">
      <div class="batch-reject-info">
        <el-text type="primary">已选择 {{ selectedIds.length }} 个套图</el-text>
        <el-text type="warning" size="small">同一原因将应用于所有选中的套图</el-text>
      </div>
      <el-input
        v-model="batchRejectReason"
        type="textarea"
        :rows="4"
        placeholder="请输入拒绝原因..."
        style="margin-top: 12px"
      />
      <template #footer>
        <el-button @click="batchRejectDialogVisible = false">取消</el-button>
        <el-button type="danger" :loading="batchProcessing" @click="handleBatchReject">批量拒绝</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPhotoSetsByStatus, approvePhotoSet, rejectPhotoSet, deletePhotoset } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const list = ref([])
const currentStatus = ref('')

const rejectDialogVisible = ref(false)
const rejectReason = ref('')
const rejectTargetId = ref(null)
const rejectLoading = ref(false)

// 删除相关状态
const deletingId = ref(null)

// 批量操作状态
const selectedIds = ref([])
const batchProcessing = ref(false)
const batchRejectDialogVisible = ref(false)
const batchRejectReason = ref('')

const statusMap = {
  pending: '待审核',
  published: '已通过',
  draft: '已拒绝'
}

const statusTagMap = {
  pending: 'warning',
  published: 'success',
  draft: 'info'
}

function statusLabel(s) {
  return statusMap[s] || s
}

function statusTagType(s) {
  return statusTagMap[s] || 'info'
}

function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

async function fetchList() {
  loading.value = true
  try {
    const params = {}
    if (currentStatus.value) params.status = currentStatus.value
    const res = await getPhotoSetsByStatus(params)
    list.value = res.data?.list || res.data || []
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
}

async function handleApprove(id) {
  try {
    await approvePhotoSet(id)
    ElMessage.success('审核通过')
    fetchList()
  } catch {
    // handled by interceptor
  }
}

function openRejectDialog(id) {
  rejectTargetId.value = id
  rejectReason.value = ''
  rejectDialogVisible.value = true
}

async function handleReject() {
  if (!rejectReason.value.trim()) {
    ElMessage.warning('请输入拒绝原因')
    return
  }
  rejectLoading.value = true
  try {
    await rejectPhotoSet(rejectTargetId.value, rejectReason.value)
    ElMessage.success('已拒绝')
    rejectDialogVisible.value = false
    fetchList()
  } catch {
    // handled by interceptor
  } finally {
    rejectLoading.value = false
  }
}

// 删除套图
async function handleDelete(id) {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个套图吗？删除后将不可恢复。',
      '警告',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    deletingId.value = id
    await deletePhotoset(id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  } finally {
    deletingId.value = null
  }
}

// 批量操作函数
function clearSelection() {
  selectedIds.value = []
}

function openBatchRejectDialog() {
  batchRejectReason.value = ''
  batchRejectDialogVisible.value = true
}

async function handleBatchApprove() {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请先选择套图')
    return
  }
  
  batchProcessing.value = true
  try {
    // 批量通过
    for (const id of selectedIds.value) {
      const item = list.value.find(item => item.id === id)
      if (item && item.status === 'pending') {
        await approvePhotoSet(id)
      }
    }
    
    ElMessage.success(`成功通过 ${selectedIds.value.length} 个套图`)
    clearSelection()
    fetchList()
  } catch (error) {
    ElMessage.error('批量操作失败')
  } finally {
    batchProcessing.value = false
  }
}

async function handleBatchReject() {
  if (selectedIds.value.length === 0) return
  
  if (!batchRejectReason.value.trim()) {
    ElMessage.warning('请输入拒绝原因')
    return
  }
  
  batchProcessing.value = true
  try {
    // 批量拒绝
    for (const id of selectedIds.value) {
      const item = list.value.find(item => item.id === id)
      if (item && item.status === 'pending') {
        await rejectPhotoSet(id, batchRejectReason.value)
      }
    }
    
    ElMessage.success(`成功拒绝 ${selectedIds.value.length} 个套图`)
    batchRejectDialogVisible.value = false
    clearSelection()
    fetchList()
  } catch (error) {
    ElMessage.error('批量操作失败')
  } finally {
    batchProcessing.value = false
  }
}

async function handleBatchDelete() {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请先选择套图')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要批量删除选中的 ${selectedIds.value.length} 个套图吗？删除后将不可恢复。`,
      '批量删除警告',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    batchProcessing.value = true
    // 批量删除
    for (const id of selectedIds.value) {
      await deletePhotoset(id)
    }
    
    ElMessage.success(`成功删除 ${selectedIds.value.length} 个套图`)
    clearSelection()
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  } finally {
    batchProcessing.value = false
  }
}

onMounted(fetchList)
</script>

<style scoped>
.action-bar {
  margin-bottom: 20px;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 16px;
}

.batch-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  background-color: #f5f7fa;
  border-radius: 8px;
  margin-left: 16px;
}

.batch-reject-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 8px;
}

.select-checkbox {
  position: absolute;
  top: 8px;
  left: 8px;
  z-index: 10;
}

.review-card {
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s;
  position: relative;
}

.review-card:hover {
  transform: translateY(-4px);
}

.review-card__cover {
  width: 100%;
  height: 180px;
  object-fit: cover;
  display: block;
}

.review-card__body {
  padding: 14px;
}

.review-card__title {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.review-card__meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
  font-size: 13px;
  color: #909399;
}

.review-card__time {
  font-size: 12px;
  color: #c0c4cc;
  margin-bottom: 10px;
}

.review-card__actions {
  display: flex;
  gap: 8px;
}
</style>
