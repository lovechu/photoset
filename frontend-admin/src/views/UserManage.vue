<template>
  <div class="user-manage">
    <div class="filter-bar">
      <el-select v-model="filterRole" placeholder="角色筛选" clearable @change="fetchUsers" style="width: 160px">
        <el-option label="全部" value="" />
        <el-option label="普通用户" value="user" />
        <el-option label="创作者" value="creator" />
        <el-option label="管理员" value="admin" />
      </el-select>
    </div>

    <el-table :data="userList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="nickname" label="昵称" min-width="120" />
      <el-table-column prop="email" label="邮箱" min-width="180" />
      <el-table-column label="角色" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="roleTagType(row.role)" size="small">{{ roleLabel(row.role) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '正常' : '已封禁' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="注册时间" width="170" align="center">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="最后登录" width="170" align="center">
        <template #default="{ row }">
          {{ formatTime(row.last_login_at) || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" align="center" fixed="right">
        <template #default="{ row }">
          <el-popconfirm
            :title="row.status === 1 ? '确认封禁该用户？' : '确认解封该用户？'"
            @confirm="handleBan(row)"
          >
            <template #reference>
              <el-button
                :type="row.status === 1 ? 'danger' : 'success'"
                size="small"
              >
                {{ row.status === 1 ? '封号' : '解封' }}
              </el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-bar">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @current-change="fetchUsers"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getUserList, banUser } from '@/api'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const userList = ref([])
const filterRole = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const roleMap = {
  user: '普通用户',
  creator: '创作者',
  admin: '管理员'
}

const roleTagMap = {
  user: 'info',
  creator: 'warning',
  admin: 'danger'
}

function roleLabel(r) {
  return roleMap[r] || r
}

function roleTagType(r) {
  return roleTagMap[r] || 'info'
}

function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

async function fetchUsers() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filterRole.value) params.role = filterRole.value
    const res = await getUserList(params)
    userList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
}

function handleSizeChange() {
  page.value = 1
  fetchUsers()
}

async function handleBan(row) {
  const newStatus = row.status === 1 ? 0 : 1
  try {
    await banUser(row.id, newStatus)
    ElMessage.success(newStatus === 0 ? '已封禁' : '已解封')
    fetchUsers()
  } catch {
    // handled by interceptor
  }
}

onMounted(fetchUsers)
</script>

<style scoped>
.filter-bar {
  margin-bottom: 20px;
}

.pagination-bar {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
