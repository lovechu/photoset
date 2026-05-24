<template>
  <div class="community-reports">
    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-select v-model="filterStatus" placeholder="状态筛选" clearable @change="handleFilterChange" style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="待处理" value="pending" />
        <el-option label="已解决" value="resolved" />
        <el-option label="已驳回" value="rejected" />
      </el-select>
    </div>

    <!-- 举报表格 -->
    <el-table :data="reportList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column label="类型" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.target_type === 'post' ? '' : 'warning'" size="small">
            {{ row.target_type === 'post' ? '帖子' : '回复' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="被举报内容" min-width="200" show-overflow-tooltip>
        <template #default="{ row }">
          {{ row.target_title || row.target_content || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="举报人" width="120" align="center">
        <template #default="{ row }">
          {{ row.reporter?.nickname || row.reporter?.username || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="举报原因" min-width="160" show-overflow-tooltip>
        <template #default="{ row }">
          {{ row.reason || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small">
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170" align="center">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="240" align="center" fixed="right">
        <template #default="{ row }">
          <div v-if="row.status === 'pending'" class="action-buttons">
            <el-button
              type="success"
              size="small"
              plain
              @click="handleQuickResolve(row, 'resolved')"
              :loading="resolvingId === row.id"
            >
              已解决
            </el-button>
            <el-button
              type="danger"
              size="small"
              plain
              @click="handleQuickResolve(row, 'rejected')"
              :loading="resolvingId === row.id"
            >
              驳回
            </el-button>
            <el-button
              type="primary"
              size="small"
              plain
              @click="openResolveDialog(row)"
            >
              处理
            </el-button>
          </div>
          <el-tag v-else :type="statusTagType(row.status)" size="small">
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>

    <!-- 空状态 -->
    <el-empty v-if="!loading && reportList.length === 0" description="暂无数据" />

    <!-- 分页 -->
    <div class="pagination-bar">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @current-change="fetchReports"
        @size-change="handleSizeChange"
      />
    </div>

    <!-- 处理举报对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="处理举报"
      width="480px"
      destroy-on-close
    >
      <div v-if="resolveTarget">
        <el-descriptions :column="1" border size="small" style="margin-bottom: 16px">
          <el-descriptions-item label="举报类型">
            {{ resolveTarget.target_type === 'post' ? '帖子' : '回复' }}
          </el-descriptions-item>
          <el-descriptions-item label="举报人">
            {{ resolveTarget.reporter?.nickname || resolveTarget.reporter?.username || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="举报原因">
            {{ resolveTarget.reason || '-' }}
          </el-descriptions-item>
        </el-descriptions>
        <el-form :model="resolveForm" :rules="resolveRules" ref="resolveFormRef" label-width="80px">
          <el-form-item label="处理结果" prop="status">
            <el-radio-group v-model="resolveForm.status">
              <el-radio value="resolved">已解决</el-radio>
              <el-radio value="rejected">驳回</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="处理备注" prop="note">
            <el-input
              v-model="resolveForm.note"
              type="textarea"
              :rows="3"
              placeholder="请输入处理备注"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitResolve">
          确认处理
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { getReports, resolveReport } from '@/api/community'
import { ElMessage, ElMessageBox } from 'element-plus'

// 列表状态
const loading = ref(false)
const reportList = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filterStatus = ref('')

// 操作状态
const resolvingId = ref(null)
const submitting = ref(false)

// 对话框
const dialogVisible = ref(false)
const resolveTarget = ref(null)
const resolveFormRef = ref(null)
const resolveForm = reactive({
  status: 'resolved',
  note: ''
})

const resolveRules = {
  status: [{ required: true, message: '请选择处理结果', trigger: 'change' }]
}

// 状态映射
const statusMap = { pending: '待处理', resolved: '已解决', rejected: '已驳回' }
const statusTagMap = { pending: 'warning', resolved: 'success', rejected: 'danger' }

function statusLabel(s) { return statusMap[s] || s || '-' }
function statusTagType(s) { return statusTagMap[s] || 'info' }

function formatTime(t) {
  if (!t) return ''
  const ts = Number(t)
  if (ts < 1e12) return new Date(ts * 1000).toLocaleString('zh-CN')
  return new Date(ts).toLocaleString('zh-CN')
}

function handleFilterChange() {
  page.value = 1
  fetchReports()
}

function handleSizeChange() {
  page.value = 1
  fetchReports()
}

async function fetchReports() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filterStatus.value) params.status = filterStatus.value
    const res = await getReports(params)
    reportList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}

function openResolveDialog(row) {
  resolveTarget.value = row
  resolveForm.status = 'resolved'
  resolveForm.note = ''
  if (resolveFormRef.value) {
    resolveFormRef.value.clearValidate()
  }
  dialogVisible.value = true
}

async function submitResolve() {
  if (!resolveFormRef.value || !resolveTarget.value) return
  try {
    await resolveFormRef.value.validate()
    submitting.value = true
    await resolveReport(resolveTarget.value.id, {
      status: resolveForm.status,
      note: resolveForm.note.trim()
    })
    ElMessage.success('处理完成')
    dialogVisible.value = false
    fetchReports()
  } catch (error) {
    if (error?.errorFields) {
      // 表单验证失败，不处理
    }
    // 其他错误已由拦截器处理
  } finally {
    submitting.value = false
  }
}

async function handleQuickResolve(row, status) {
  const label = status === 'resolved' ? '已解决' : '驳回'
  try {
    await ElMessageBox.confirm(
      `确定将该举报标记为「${label}」吗？`,
      `${label}确认`,
      { confirmButtonText: `确认${label}`, cancelButtonText: '取消', type: status === 'resolved' ? 'success' : 'warning' }
    )
    resolvingId.value = row.id
    await resolveReport(row.id, { status, note: '' })
    ElMessage.success(`已${label}`)
    fetchReports()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      // 错误已由拦截器处理
    }
  } finally {
    resolvingId.value = null
  }
}

onMounted(fetchReports)
</script>

<style scoped>
.filter-bar {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}
.pagination-bar {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
.action-buttons {
  display: flex;
  justify-content: center;
  gap: 6px;
  flex-wrap: wrap;
}
</style>
