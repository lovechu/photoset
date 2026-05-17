<template>
  <div class="user-manage">
    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-button type="success" plain @click="handleExport" :loading="exporting">
        导出 CSV
      </el-button>
      <el-input
        v-model="keyword"
        placeholder="搜索昵称 / 邮箱"
        clearable
        style="width: 220px"
        @clear="fetchUsers"
        @keyup.enter="handleSearch"
      >
        <template #append>
          <el-button @click="handleSearch"><el-icon><Search /></el-icon></el-button>
        </template>
      </el-input>

      <el-select v-model="filterRole" placeholder="角色筛选" clearable @change="fetchUsers" style="width: 140px; margin-left: 12px">
        <el-option label="全部角色" value="" />
        <el-option label="访客" value="guest" />
        <el-option label="普通用户" value="user" />
        <el-option label="会员" value="member" />
        <el-option label="创作者" value="creator" />
        <el-option label="管理员" value="admin" />
      </el-select>

      <el-select v-model="filterStatus" placeholder="状态筛选" clearable @change="fetchUsers" style="width: 120px; margin-left: 12px">
        <el-option label="全部状态" value="" />
        <el-option label="正常" value="1" />
        <el-option label="已封禁" value="0" />
      </el-select>
    </div>

    <!-- 用户表格 -->
    <el-table :data="userList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="65" align="center" />
      <el-table-column prop="nickname" label="昵称" min-width="100">
        <template #default="{ row }">
          <el-link type="primary" @click="openDetail(row.id)">{{ row.nickname || '-' }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
      <el-table-column label="角色" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="roleTagType(row.role)" size="small" effect="plain">{{ roleLabel(row.role) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="85" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '正常' : '封禁' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="会员到期" width="115" align="center">
        <template #default="{ row }">
          <span v-if="row.membership_expires && row.membership_expires > 0" class="member-expire">
            {{ formatTime(row.membership_expires).split(' ')[0] }}
          </span>
          <span v-else class="text-muted">-</span>
        </template>
      </el-table-column>
      <el-table-column label="注册时间" width="165" align="center" sortable :sort-method="(a,b) => a.created_at - b.created_at">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="最后登录" width="165" align="center">
        <template #default="{ row }">{{ formatTime(row.last_login_at) || '-' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="280" align="center" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openDetail(row.id)">详情</el-button>
          <el-button size="small" type="warning" @click="openRoleDialog(row)">角色</el-button>
          <el-button size="small" type="info" @click="openPasswordDialog(row)">密码</el-button>
          <el-popconfirm
            :title="row.status === 1 ? '确认封禁该用户？' : '确认解封该用户？'"
            @confirm="handleBan(row)"
          >
            <template #reference>
              <el-button :type="row.status === 1 ? 'danger' : 'success'" size="small">
                {{ row.status === 1 ? '封号' : '解封' }}
              </el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-bar">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @current-change="fetchUsers"
        @size-change="handleSizeChange"
      />
    </div>

    <!-- 用户详情抽屉 -->
    <el-drawer v-model="detailVisible" title="用户详情" size="420px" destroy-on-close>
      <div v-loading="detailLoading" class="user-detail" v-if="detail">
        <div class="detail-header">
          <el-avatar :size="64">
            {{ detail.user?.nickname?.charAt(0) || 'U' }}
          </el-avatar>
          <div class="detail-header-info">
            <h3>{{ detail.user?.nickname || '未知用户' }}</h3>
            <el-tag :type="roleTagType(detail.user?.role)" size="small">{{ roleLabel(detail.user?.role) }}</el-tag>
            <el-tag :type="detail.user?.status === 1 ? 'success' : 'danger'" size="small" style="margin-left: 6px">
              {{ detail.user?.status === 1 ? '正常' : '封禁' }}
            </el-tag>
          </div>
        </div>

        <el-descriptions :column="1" border size="small" class="detail-desc">
          <el-descriptions-item label="用户ID">{{ detail.user?.id }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ detail.user?.email || '-' }}</el-descriptions-item>
          <el-descriptions-item label="角色">{{ roleLabel(detail.user?.role) }}</el-descriptions-item>
          <el-descriptions-item label="注册时间">{{ formatTime(detail.user?.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="最后登录">{{ formatTime(detail.user?.last_login_at) || '-' }}</el-descriptions-item>
          <el-descriptions-item label="会员状态">
            {{ detail.user?.membership_expires ? formatTime(detail.user?.membership_expires) : '非会员' }}
          </el-descriptions-item>
        </el-descriptions>

        <el-divider content-position="left">数据统计</el-divider>
        <el-row :gutter="16">
          <el-col :span="6">
            <el-statistic title="发布套图" :value="detail.photoset_count" />
          </el-col>
          <el-col :span="6">
            <el-statistic title="收藏数" :value="detail.favorite_count" />
          </el-col>
          <el-col :span="6">
            <el-statistic title="订单数" :value="detail.order_count" />
          </el-col>
          <el-col :span="6">
            <el-statistic title="消费总额" :value="detail.total_spent" :precision="2" prefix="¥" />
          </el-col>
        </el-row>
      </div>
    </el-drawer>

    <!-- 修改角色对话框 -->
    <el-dialog v-model="roleDialogVisible" title="修改用户角色" width="380px" destroy-on-close>
      <div v-if="roleTarget">
        <p>用户：<strong>{{ roleTarget.nickname }}</strong>（ID: {{ roleTarget.id }}）</p>
        <el-select v-model="newRole" placeholder="选择角色" style="width: 100%; margin-top: 12px">
          <el-option label="访客" value="guest" />
          <el-option label="普通用户" value="user" />
          <el-option label="会员" value="member" />
          <el-option label="创作者" value="creator" />
          <el-option label="管理员" value="admin" />
        </el-select>
      </div>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="roleLoading" @click="handleRoleChange">确认修改</el-button>
      </template>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog v-model="passwordDialogVisible" title="重置用户密码" width="400px" destroy-on-close>
      <div v-if="passwordTarget">
        <p>用户：<strong>{{ passwordTarget.nickname }}</strong>（{{ passwordTarget.email }}）</p>
        <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" style="margin-top: 16px">
          <el-form-item label="新密码" prop="newPassword">
            <el-input v-model="passwordForm.newPassword" type="password" placeholder="至少6位" show-password />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="再次输入新密码" show-password />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="passwordLoading" @click="handleResetPassword">确认重置</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getUserList, banUser, updateUserRole, getUserDetail, resetUserPassword, exportUsers } from '@/api'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'

// 列表状态
const loading = ref(false)
const userList = ref([])
const keyword = ref('')
const filterRole = ref('')
const filterStatus = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const exporting = ref(false)

// 详情
const detailVisible = ref(false)
const detailLoading = ref(false)
const detail = ref(null)

// 角色修改
const roleDialogVisible = ref(false)
const roleTarget = ref(null)
const newRole = ref('')
const roleLoading = ref(false)

const roleMap = {
  guest: '访客', user: '普通用户', member: '会员', creator: '创作者', admin: '管理员'
}
const roleTagMap = {
  guest: 'info', user: '', member: 'success', creator: 'warning', admin: 'danger'
}

function roleLabel(r) { return roleMap[r] || r || '-' }
function roleTagType(r) { return roleTagMap[r] || 'info' }

function formatTime(t) {
  if (!t) return ''
  const ts = Number(t)
  if (ts < 1e12) return new Date(ts * 1000).toLocaleString('zh-CN')
  return new Date(ts).toLocaleString('zh-CN')
}

function handleSearch() {
  page.value = 1
  fetchUsers()
}

function handleSizeChange() {
  page.value = 1
  fetchUsers()
}

async function fetchUsers() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filterRole.value) params.role = filterRole.value
    // status: 空字符串表示全部（传递-1），'1'表示正常，'0'表示已封禁
    params.status = filterStatus.value !== '' && filterStatus.value !== undefined ? parseInt(filterStatus.value) : -1
    if (keyword.value) params.keyword = keyword.value
    const res = await getUserList(params)
    userList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch { /* handled by interceptor */ }
  finally { loading.value = false }
}

async function handleBan(row) {
  const newStatus = Number(row.status) === 1 ? 0 : 1
  try {
    await banUser(row.id, newStatus)
    ElMessage.success(newStatus === 0 ? '已封禁' : '已解封')
    fetchUsers()
  } catch { /* handled */ }
}

async function openDetail(id) {
  detailVisible.value = true
  detailLoading.value = true
  detail.value = null
  try {
    const res = await getUserDetail(id)
    detail.value = res.data
  } catch { /* handled */ }
  finally { detailLoading.value = false }
}

function openRoleDialog(row) {
  roleTarget.value = row
  newRole.value = row.role
  roleDialogVisible.value = true
}

async function handleRoleChange() {
  if (!newRole.value) return
  roleLoading.value = true
  try {
    await updateUserRole(roleTarget.value.id, newRole.value)
    ElMessage.success('角色已更新')
    roleDialogVisible.value = false
    fetchUsers()
  } catch { /* handled */ }
  finally { roleLoading.value = false }
}

// 重置密码
const passwordDialogVisible = ref(false)
const passwordTarget = ref(null)
const passwordLoading = ref(false)
const passwordFormRef = ref(null)
const passwordForm = ref({ newPassword: '', confirmPassword: '' })
const passwordRules = {
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.value.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

function openPasswordDialog(row) {
  passwordTarget.value = row
  passwordForm.value = { newPassword: '', confirmPassword: '' }
  passwordDialogVisible.value = true
}

async function handleResetPassword() {
  if (!passwordFormRef.value) return
  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return
    passwordLoading.value = true
    try {
      await resetUserPassword(passwordTarget.value.id, passwordForm.value.newPassword)
      ElMessage.success('密码已重置')
      passwordDialogVisible.value = false
    } catch { /* handled */ }
    finally { passwordLoading.value = false }
  })
}

async function handleExport() {
  exporting.value = true
  try {
    const params = {}
    if (filterRole.value) params.role = filterRole.value
    if (filterStatus.value !== '' && filterStatus.value !== undefined) {
      params.status = filterStatus.value !== '' ? parseInt(filterStatus.value) : -1
    }
    if (keyword.value) params.keyword = keyword.value
    const res = await exportUsers(params)
    const blob = new Blob([res], { type: 'text/csv;charset=utf-8' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'users.csv'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

onMounted(fetchUsers)
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
.user-detail { padding: 0 8px; }
.detail-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}
.detail-header-info h3 {
  margin: 0 0 6px;
  font-size: 18px;
}
.detail-desc { margin-top: 16px; }
.text-muted { color: #999; }
.member-expire { font-size: 12px; color: #67c23a; }
</style>
