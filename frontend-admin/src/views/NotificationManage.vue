<template>
  <div class="notification-manage">
    <div class="header-bar">
      <el-button type="primary" @click="handleSend" :icon="Plus">
        发送通知
      </el-button>
      <el-input
        v-model="filterKeyword"
        placeholder="搜索通知标题"
        clearable
        @clear="fetchNotifications"
        @keyup.enter="fetchNotifications"
        style="width: 240px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-number">{{ stats.total_notifications }}</div>
          <div class="stat-label">总通知数</div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-number" style="color: #E6A23C">{{ stats.unread_notifications }}</div>
          <div class="stat-label">未读通知</div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-number" style="color: #67C23A">{{ stats.today_notifications }}</div>
          <div class="stat-label">今日通知</div>
        </div>
      </el-card>
    </div>

    <el-table :data="filteredNotifications" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="title" label="标题" min-width="150" />
      <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
      <el-table-column label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getTypeTagType(row.type)" size="small">{{ getTypeName(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="接收用户" width="120" align="center">
        <template #default="{ row }">
          <span>{{ row.user?.nickname || '未知用户' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="发送者" width="120" align="center">
        <template #default="{ row }">
          <span>{{ row.sender?.nickname || '系统' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_read ? 'success' : 'warning'" size="small">
            {{ row.is_read ? '已读' : '未读' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="发送时间" width="170" align="center">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 发送通知对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="发送系统通知"
      width="600px"
      @close="resetForm"
    >
      <el-form ref="formRef" :model="sendForm" :rules="formRules" label-width="100px">
        <el-form-item label="通知标题" prop="title">
          <el-input v-model="sendForm.title" placeholder="请输入通知标题" />
        </el-form-item>
        <el-form-item label="通知内容" prop="content">
          <el-input
            v-model="sendForm.content"
            type="textarea"
            :rows="4"
            placeholder="请输入通知内容"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="发送范围" prop="scope">
          <el-radio-group v-model="sendForm.scope">
            <el-radio label="all">全部用户</el-radio>
            <el-radio label="custom">指定用户</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="sendForm.scope === 'custom'" label="用户ID" prop="user_ids">
          <el-input v-model="sendForm.user_ids_str" placeholder="请输入用户ID，多个用逗号分隔" />
          <div class="slug-hint">
            <span class="slug-preview">多个用户ID请用英文逗号分隔，如: 1,2,3</span>
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          发送通知
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, reactive, computed } from 'vue'
import { getNotifications, sendNotification, getNotificationStats } from '@/api/community'
import { ElMessage } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'

const loading = ref(false)
const notificationList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const filterKeyword = ref('')
const submitting = ref(false)
const dialogVisible = ref(false)

const stats = reactive({
  total_notifications: 0,
  unread_notifications: 0,
  today_notifications: 0
})

const sendForm = reactive({
  title: '',
  content: '',
  scope: 'all',
  user_ids_str: '',
  user_ids: []
})

const formRules = {
  title: [
    { required: true, message: '请输入通知标题', trigger: 'blur' },
    { min: 2, max: 200, message: '标题长度在 2 到 200 个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入通知内容', trigger: 'blur' },
    { min: 1, max: 500, message: '内容长度在 1 到 500 个字符', trigger: 'blur' }
  ],
  scope: [
    { required: true, message: '请选择发送范围', trigger: 'change' }
  ]
}

const formRef = ref()

const filteredNotifications = computed(() => {
  if (!filterKeyword.value) return notificationList.value
  const keyword = filterKeyword.value.toLowerCase()
  return notificationList.value.filter(item => 
    item.title.toLowerCase().includes(keyword) ||
    (item.content && item.content.toLowerCase().includes(keyword))
  )
})

function getTypeName(type) {
  const typeMap = {
    'like': '点赞',
    'reply': '回复',
    'follow': '关注',
    'mention': '提及',
    'system': '系统'
  }
  return typeMap[type] || type
}

function getTypeTagType(type) {
  const typeMap = {
    'like': 'warning',
    'reply': 'primary',
    'follow': 'success',
    'mention': 'info',
    'system': 'danger'
  }
  return typeMap[type] || ''
}

function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

async function fetchNotifications() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    const res = await getNotifications(params)
    notificationList.value = res.data?.notifications || []
    total.value = res.data?.pagination?.total || 0
  } catch (error) {
    console.error('获取通知列表失败:', error)
    ElMessage.error('获取通知列表失败')
  } finally {
    loading.value = false
  }
}

async function fetchStats() {
  try {
    const res = await getNotificationStats()
    Object.assign(stats, res.data?.stats || {})
  } catch (error) {
    console.error('获取通知统计失败:', error)
  }
}

function handleSizeChange(val) {
  pageSize.value = val
  currentPage.value = 1
  fetchNotifications()
}

function handleCurrentChange(val) {
  currentPage.value = val
  fetchNotifications()
}

function handleSend() {
  resetForm()
  dialogVisible.value = true
}

function resetForm() {
  Object.assign(sendForm, {
    title: '',
    content: '',
    scope: 'all',
    user_ids_str: '',
    user_ids: []
  })
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

async function submitForm() {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true

    const formData = {
      title: sendForm.title.trim(),
      content: sendForm.content.trim(),
      user_ids: sendForm.scope === 'all' ? [] : sendForm.user_ids
    }

    await sendNotification(formData)
    ElMessage.success('通知发送成功')
    dialogVisible.value = false
    fetchNotifications()
    fetchStats()
  } catch (error) {
    if (error.errorFields) {
      // 验证失败，不处理
    } else {
      ElMessage.error('通知发送失败')
    }
  } finally {
    submitting.value = false
  }
}

// 监听筛选条件变化
watch([filterKeyword], () => {
  // 由于使用计算属性，这里不需要额外操作
})

onMounted(() => {
  fetchNotifications()
  fetchStats()
})
</script>

<style scoped>
.header-bar {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 8px;
}

.stat-content {
  text-align: center;
}

.stat-number {
  font-size: 32px;
  font-weight: bold;
  color: #409EFF;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.slug-hint {
  margin-top: 6px;
  font-size: 12px;
}

.slug-preview {
  color: #909399;
  background: #f2f3f5;
  padding: 4px 8px;
  border-radius: 4px;
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}
</style>