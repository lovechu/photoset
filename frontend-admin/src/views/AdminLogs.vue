<template>
  <div class="admin-logs">
    <div class="filter-bar">
      <el-select v-model="filterAction" placeholder="操作类型" clearable @change="handleSearch" style="width: 160px">
        <el-option label="全部操作" value="" />
        <el-option v-for="a in actionOptions" :key="a.value" :label="a.label" :value="a.value" />
      </el-select>
    </div>

    <el-table :data="logList" v-loading="loading" stripe border style="width: 100%">
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column label="管理员" width="120">
        <template #default="{ row }">{{ row.admin_name || '-' }}</template>
      </el-table-column>
      <el-table-column label="操作类型" width="110" align="center">
        <template #default="{ row }">
          <el-tag :type="actionTagType(row.action)" size="small" effect="plain">{{ actionLabel(row.action) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="target" label="操作对象" width="160" show-overflow-tooltip />
      <el-table-column prop="detail" label="详细信息" min-width="200" show-overflow-tooltip />
      <el-table-column prop="ip" label="IP" width="130" align="center" />
      <el-table-column label="时间" width="170" align="center">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
    </el-table>

    <div class="pagination-bar">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @current-change="fetchLogs"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAdminLogs } from '@/api'

const loading = ref(false)
const logList = ref([])
const filterAction = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const actionOptions = [
  { label: '封禁用户', value: 'ban_user' },
  { label: '解封用户', value: 'unban_user' },
  { label: '审核通过', value: 'approve' },
  { label: '审核拒绝', value: 'reject' },
  { label: '角色修改', value: 'role_change' },
  { label: '退款', value: 'refund' },
  { label: '删除套图', value: 'delete_photoset' },
  { label: '设置更新', value: 'settings_update' },
]

const actionMap = {}
actionOptions.forEach(a => { actionMap[a.value] = a.label })

const actionTagMap = {
  ban_user: 'danger', unban_user: 'success', approve: 'success',
  reject: 'warning', role_change: 'warning', refund: 'danger',
  delete_photoset: 'danger', settings_update: 'info',
}

function actionLabel(a) { return actionMap[a] || a }
function actionTagType(a) { return actionTagMap[a] || 'info' }

function formatTime(t) {
  if (!t) return ''
  const ts = Number(t)
  if (ts < 1e12) return new Date(ts * 1000).toLocaleString('zh-CN')
  return new Date(ts).toLocaleString('zh-CN')
}

function handleSearch() {
  page.value = 1
  fetchLogs()
}

function handleSizeChange() {
  page.value = 1
  fetchLogs()
}

async function fetchLogs() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filterAction.value) params.action = filterAction.value
    const res = await getAdminLogs(params)
    logList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {}
  finally { loading.value = false }
}

onMounted(fetchLogs)
</script>

<style scoped>
.filter-bar { margin-bottom: 20px; }
.pagination-bar { margin-top: 20px; display: flex; justify-content: flex-end; }
</style>
