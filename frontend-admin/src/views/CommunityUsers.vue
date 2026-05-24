<template>
  <div class="community-users">
    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-select v-model="filterLevel" placeholder="等级筛选" clearable @change="handleFilterChange" style="width: 140px">
        <el-option label="全部等级" value="" />
        <el-option v-for="lv in levels" :key="lv" :label="lv" :value="lv" />
      </el-select>
    </div>

    <!-- 用户表格 -->
    <el-table :data="userList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="username" label="用户名" min-width="130" />
      <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
      <el-table-column label="积分" width="100" align="center">
        <template #default="{ row }">
          <el-tag type="warning" size="small">{{ row.points ?? 0 }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="等级" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="levelTagType(row.level)" size="small">
            {{ row.level || 'L1' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="最后更新" width="170" align="center">
        <template #default="{ row }">{{ formatTime(row.updated_at) || formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="120" align="center" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" size="small" plain @click="openAdjustDialog(row)">
            调整积分
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 空状态 -->
    <el-empty v-if="!loading && userList.length === 0" description="暂无数据" />

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

    <!-- 调整积分对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="调整用户积分"
      width="450px"
      destroy-on-close
      @close="resetAdjustForm"
    >
      <div v-if="adjustTarget">
        <el-descriptions :column="1" border size="small" style="margin-bottom: 16px">
          <el-descriptions-item label="用户名">
            {{ adjustTarget.username || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="当前积分">
            {{ adjustTarget.points ?? 0 }}
          </el-descriptions-item>
          <el-descriptions-item label="当前等级">
            {{ adjustTarget.level || 'L1' }}
          </el-descriptions-item>
        </el-descriptions>
        <el-form ref="adjustFormRef" :model="adjustForm" :rules="adjustRules" label-width="80px">
          <el-form-item label="积分值" prop="points">
            <el-input-number
              v-model="adjustForm.points"
              :min="-999999"
              :max="999999"
              style="width: 100%"
              placeholder="正数增加，负数减少"
            />
          </el-form-item>
          <el-form-item label="调整原因" prop="reason">
            <el-input
              v-model="adjustForm.reason"
              type="textarea"
              :rows="2"
              placeholder="请输入调整原因"
              maxlength="200"
              show-word-limit
            />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitAdjust">
          确认调整
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { getCommunityUsers, adjustUserPoints } from '@/api/community'
import { ElMessage } from 'element-plus'

// 等级列表
const levels = ['L1', 'L2', 'L3', 'L4', 'L5', 'L6', 'L7', 'L8', 'L9', 'L10']

// 列表状态
const loading = ref(false)
const userList = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filterLevel = ref('')

// 操作状态
const submitting = ref(false)

// 对话框
const dialogVisible = ref(false)
const adjustTarget = ref(null)
const adjustFormRef = ref(null)
const adjustForm = reactive({
  points: 0,
  reason: ''
})

const adjustRules = {
  points: [{ required: true, message: '请输入积分值', trigger: 'blur' }],
  reason: [{ required: true, message: '请输入调整原因', trigger: 'blur' }]
}

// 等级标签颜色映射
const levelTagColors = { L1: 'info', L2: 'info', L3: '', L4: '', L5: 'success', L6: 'success', L7: 'warning', L8: 'warning', L9: 'danger', L10: 'danger' }
function levelTagType(lv) { return levelTagColors[lv] || 'info' }

function formatTime(t) {
  if (!t) return ''
  const ts = Number(t)
  if (ts < 1e12) return new Date(ts * 1000).toLocaleString('zh-CN')
  return new Date(ts).toLocaleString('zh-CN')
}

function handleFilterChange() {
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
    if (filterLevel.value) params.level = filterLevel.value
    const res = await getCommunityUsers(params)
    userList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}

function openAdjustDialog(row) {
  adjustTarget.value = row
  adjustForm.points = 0
  adjustForm.reason = ''
  if (adjustFormRef.value) {
    adjustFormRef.value.clearValidate()
  }
  dialogVisible.value = true
}

function resetAdjustForm() {
  adjustTarget.value = null
  adjustForm.points = 0
  adjustForm.reason = ''
  if (adjustFormRef.value) {
    adjustFormRef.value.clearValidate()
  }
}

async function submitAdjust() {
  if (!adjustFormRef.value || !adjustTarget.value) return
  try {
    await adjustFormRef.value.validate()
    submitting.value = true
    await adjustUserPoints(adjustTarget.value.id, {
      points: adjustForm.points,
      reason: adjustForm.reason.trim()
    })
    ElMessage.success('积分调整成功')
    dialogVisible.value = false
    fetchUsers()
  } catch (error) {
    if (error?.errorFields) {
      // 表单验证失败，不处理
    }
    // 其他错误已由拦截器处理
  } finally {
    submitting.value = false
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
</style>
