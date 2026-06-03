<template>
  <div class="user-manage">
  <!-- 筛选栏 -->
  <div class="filter-bar">
    <el-button type="primary" @click="openCreateDialog">
      <el-icon><Plus /></el-icon>
      新增用户
    </el-button>
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
    <el-drawer v-model="detailVisible" title="用户详情" size="700px" destroy-on-close>
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

        <el-tabs v-model="activeTab" @tab-change="handleTabChange">
          <!-- 基本信息标签页 -->
          <el-tab-pane label="基本信息" name="basic">
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
          </el-tab-pane>

          <!-- 登录历史标签页 -->
          <el-tab-pane label="登录历史" name="loginHistory">
            <el-table :data="loginHistory" v-loading="loginHistoryLoading" stripe style="width: 100%">
              <el-table-column prop="ip" label="IP地址" width="140" />
              <el-table-column prop="ip_location" label="位置" min-width="120" show-overflow-tooltip />
              <el-table-column prop="device" label="设备" width="100" />
              <el-table-column prop="browser" label="浏览器" width="100" />
              <el-table-column prop="os" label="系统" width="100" />
              <el-table-column label="登录时间" width="165">
                <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
              </el-table-column>
              <el-table-column label="状态" width="80" align="center">
                <template #default="{ row }">
                  <el-tag :type="row.success ? 'success' : 'danger'" size="small">
                    {{ row.success ? '成功' : '失败' }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
            <el-pagination
              v-if="loginHistoryTotal > 10"
              v-model:current-page="loginHistoryPage"
              :total="loginHistoryTotal"
              :page-size="10"
              layout="total, prev, pager, next"
              @current-change="handleLoginHistoryPageChange"
              style="margin-top: 16px; justify-content: flex-end;"
            />
          </el-tab-pane>

          <!-- 设备管理标签页 -->
          <el-tab-pane label="设备管理" name="devices">
            <el-table :data="userDevices" v-loading="devicesLoading" stripe style="width: 100%">
              <el-table-column prop="device_name" label="设备名称" min-width="150" show-overflow-tooltip />
              <el-table-column prop="device_type" label="设备类型" width="100" />
              <el-table-column prop="os" label="操作系统" width="120" />
              <el-table-column prop="browser" label="浏览器" width="100" />
              <el-table-column prop="ip" label="IP地址" width="140" />
              <el-table-column label="最后活跃" width="165">
                <template #default="{ row }">{{ formatTime(row.last_active_at) }}</template>
              </el-table-column>
              <el-table-column label="状态" width="80" align="center">
                <template #default="{ row }">
                  <el-tag :type="row.is_active ? 'success' : 'info'" size="small">
                    {{ row.is_active ? '活跃' : '已停用' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="100" align="center">
                <template #default="{ row }">
                  <el-button
                    v-if="row.is_active"
                    size="small"
                    type="danger"
                    @click="handleDeactivateDevice(row.device_id)"
                  >
                    停用
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <!-- 隐私设置标签页 -->
          <el-tab-pane label="隐私设置" name="privacy">
            <div v-loading="privacyLoading" v-if="userPrivacySettings">
              <el-form label-width="140px">
                <el-form-item label="公开个人资料">
                  <el-switch
                    :model-value="userPrivacySettings.show_profile"
                    @change="(val) => handleUpdatePrivacy('show_profile', val)"
                  />
                  <span class="setting-desc">允许其他用户查看你的个人资料</span>
                </el-form-item>
                <el-form-item label="公开发布内容">
                  <el-switch
                    :model-value="userPrivacySettings.show_posts"
                    @change="(val) => handleUpdatePrivacy('show_posts', val)"
                  />
                  <span class="setting-desc">允许其他用户查看你发布的套图</span>
                </el-form-item>
                <el-form-item label="公开收藏列表">
                  <el-switch
                    :model-value="userPrivacySettings.show_favorites"
                    @change="(val) => handleUpdatePrivacy('show_favorites', val)"
                  />
                  <span class="setting-desc">允许其他用户查看你的收藏</span>
                </el-form-item>
                <el-form-item label="允许被搜索">
                  <el-switch
                    :model-value="userPrivacySettings.allow_search"
                    @change="(val) => handleUpdatePrivacy('allow_search', val)"
                  />
                  <span class="setting-desc">允许其他用户通过搜索找到你</span>
                </el-form-item>
                <el-form-item label="允许私信">
                  <el-switch
                    :model-value="userPrivacySettings.allow_message"
                    @change="(val) => handleUpdatePrivacy('allow_message', val)"
                  />
                  <span class="setting-desc">允许其他用户给你发送私信</span>
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>
        </el-tabs>
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

    <!-- 创建用户对话框 -->
    <el-dialog v-model="createDialogVisible" title="新增用户" width="480px" destroy-on-close>
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="80px">
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="createForm.nickname" placeholder="请输入昵称（2-50个字符）" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="createForm.email" placeholder="请输入邮箱地址" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="createForm.password" type="password" placeholder="请输入密码（至少6位）" show-password />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="createForm.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="访客" value="guest" />
            <el-option label="普通用户" value="user" />
            <el-option label="会员" value="member" />
            <el-option label="创作者" value="creator" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="createLoading" @click="handleCreateUser">确认创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getUserList, banUser, updateUserRole, getUserDetail, resetUserPassword, exportUsers, createUser, getUserLoginHistory, getUserDevices, deactivateUserDevice, getUserPrivacySettings, updateUserPrivacySettings } from '@/api'
import { ElMessage } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'

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

// 创建用户
const createDialogVisible = ref(false)
const createLoading = ref(false)
const createFormRef = ref(null)
const createForm = ref({
  nickname: '',
  email: '',
  password: '',
  role: 'user'
})
const createRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 50, message: '昵称长度在2-50个字符', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' },
  ],
}

// 用户详情标签页
const activeTab = ref('basic')
const loginHistory = ref([])
const loginHistoryLoading = ref(false)
const loginHistoryPage = ref(1)
const loginHistoryTotal = ref(0)

const userDevices = ref([])
const devicesLoading = ref(false)

const userPrivacySettings = ref(null)
const privacyLoading = ref(false)

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
  // 重置标签页状态
  activeTab.value = 'basic'
  loginHistory.value = []
  loginHistoryPage.value = 1
  loginHistoryTotal.value = 0
  userDevices.value = []
  userPrivacySettings.value = null
  try {
    const res = await getUserDetail(id)
    detail.value = res.data
  } catch { /* handled */ }
  finally { detailLoading.value = false }
}

// 标签页切换处理
async function handleTabChange(tab) {
  if (!detail.value) return
  const userId = detail.value.user?.id
  if (!userId) return

  if (tab === 'loginHistory' && loginHistory.value.length === 0) {
    await fetchLoginHistory(userId)
  } else if (tab === 'devices' && userDevices.value.length === 0) {
    await fetchUserDevices(userId)
  } else if (tab === 'privacy' && !userPrivacySettings.value) {
    await fetchUserPrivacySettings(userId)
  }
}

// 获取登录历史
async function fetchLoginHistory(userId) {
  loginHistoryLoading.value = true
  try {
    const res = await getUserLoginHistory(userId, {
      page: loginHistoryPage.value,
      page_size: 10
    })
    loginHistory.value = res.data?.list || []
    loginHistoryTotal.value = res.data?.total || 0
  } catch { /* handled */ }
  finally { loginHistoryLoading.value = false }
}

// 登录历史分页变化
function handleLoginHistoryPageChange(page) {
  loginHistoryPage.value = page
  if (detail.value?.user?.id) {
    fetchLoginHistory(detail.value.user.id)
  }
}

// 获取用户设备
async function fetchUserDevices(userId) {
  devicesLoading.value = true
  try {
    const res = await getUserDevices(userId)
    userDevices.value = res.data || []
  } catch { /* handled */ }
  finally { devicesLoading.value = false }
}

// 停用设备
async function handleDeactivateDevice(deviceId) {
  if (!detail.value?.user?.id) return
  try {
    await deactivateUserDevice(detail.value.user.id, deviceId)
    ElMessage.success('设备已停用')
    await fetchUserDevices(detail.value.user.id)
  } catch { /* handled */ }
}

// 获取隐私设置
async function fetchUserPrivacySettings(userId) {
  privacyLoading.value = true
  try {
    const res = await getUserPrivacySettings(userId)
    userPrivacySettings.value = res.data
  } catch { /* handled */ }
  finally { privacyLoading.value = false }
}

// 更新隐私设置
async function handleUpdatePrivacy(field, value) {
  if (!detail.value?.user?.id || !userPrivacySettings.value) return
  const updated = { ...userPrivacySettings.value, [field]: value }
  try {
    await updateUserPrivacySettings(detail.value.user.id, updated)
    userPrivacySettings.value = updated
    ElMessage.success('隐私设置已更新')
  } catch { /* handled */ }
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

// 创建用户
function openCreateDialog() {
  createForm.value = { nickname: '', email: '', password: '', role: 'user' }
  createDialogVisible.value = true
}

async function handleCreateUser() {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    createLoading.value = true
    try {
      await createUser(createForm.value)
      ElMessage.success('用户创建成功')
      createDialogVisible.value = false
      fetchUsers()
    } catch { /* handled */ }
    finally { createLoading.value = false }
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
.setting-desc {
  margin-left: 12px;
  color: #909399;
  font-size: 12px;
}
</style>
