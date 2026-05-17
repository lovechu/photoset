<template>
  <div class="membership-manage">
    <div class="header-bar">
      <el-button type="primary" @click="handleCreate" :icon="Plus">
        新增套餐
      </el-button>
    </div>

    <el-table :data="planList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="name" label="套餐名称" min-width="140" />
      <el-table-column label="时长" width="120" align="center">
        <template #default="{ row }">
          {{ formatDuration(row.duration) }}
        </template>
      </el-table-column>
      <el-table-column label="价格" width="120" align="right">
        <template #default="{ row }">
          ¥{{ row.price.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
            {{ row.status === 1 ? '上架' : '下架' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170" align="center">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" align="center" fixed="right">
        <template #default="{ row }">
          <div class="action-buttons">
            <el-button
              type="primary"
              size="small"
              plain
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              :type="row.status === 1 ? 'warning' : 'success'"
              size="small"
              plain
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 1 ? '下架' : '上架' }}
            </el-button>
            <el-popconfirm
              title="确定要删除这个套餐吗？"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button
                  type="danger"
                  size="small"
                  plain
                >
                  删除
                </el-button>
              </template>
            </el-popconfirm>
          </div>
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

    <!-- 新建/编辑套餐对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="520px"
      @close="resetForm"
    >
      <el-form ref="planFormRef" :model="planForm" :rules="planRules" label-width="80px">
        <el-form-item label="套餐名称" prop="name">
          <el-input v-model="planForm.name" placeholder="请输入套餐名称" maxlength="50" />
        </el-form-item>
        <el-form-item label="时长(天)" prop="duration">
          <el-input-number v-model="planForm.duration" :min="1" :max="3650" style="width: 100%" />
        </el-form-item>
        <el-form-item label="价格(元)" prop="price">
          <el-input-number v-model="planForm.price" :min="0" :precision="2" :step="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="planForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入套餐描述"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="planForm.status" :active-value="1" :inactive-value="0" active-text="上架" inactive-text="下架" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ isEditing ? '保存修改' : '创建套餐' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { getMembershipList, createMembership, updateMembership, deleteMembership } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const loading = ref(false)
const planList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const submitting = ref(false)
const dialogVisible = ref(false)
const isEditing = ref(false)

const planForm = reactive({
  id: null,
  name: '',
  duration: 30,
  price: 0,
  description: '',
  status: 1
})

const planRules = {
  name: [
    { required: true, message: '请输入套餐名称', trigger: 'blur' },
    { min: 2, max: 50, message: '套餐名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  duration: [
    { required: true, message: '请输入套餐时长', trigger: 'blur' }
  ],
  price: [
    { required: true, message: '请输入套餐价格', trigger: 'blur' }
  ]
}

const planFormRef = ref()

function formatDuration(days) {
  if (days >= 365) {
    const years = days / 365
    return years === Math.floor(years) ? `${Math.floor(years)}年` : `${days}天`
  }
  if (days >= 30) {
    const months = days / 30
    return months === Math.floor(months) ? `${Math.floor(months)}个月` : `${days}天`
  }
  return `${days}天`
}

function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

async function fetchPlans() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value
    }
    const res = await getMembershipList(params)
    planList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('获取套餐列表失败:', error)
    ElMessage.error('获取套餐列表失败')
  } finally {
    loading.value = false
  }
}

function handleSizeChange(val) {
  pageSize.value = val
  currentPage.value = 1
  fetchPlans()
}

function handleCurrentChange(val) {
  currentPage.value = val
  fetchPlans()
}

function handleCreate() {
  isEditing.value = false
  resetForm()
  dialogVisible.value = true
}

function handleEdit(plan) {
  isEditing.value = true
  Object.assign(planForm, {
    id: plan.id,
    name: plan.name,
    duration: plan.duration,
    price: plan.price,
    description: plan.description || '',
    status: plan.status
  })
  dialogVisible.value = true
}

async function handleToggleStatus(plan) {
  const newStatus = plan.status === 1 ? 0 : 1
  try {
    await updateMembership(plan.id, { status: newStatus })
    ElMessage.success(newStatus === 1 ? '已上架' : '已下架')
    fetchPlans()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

async function handleDelete(plan) {
  try {
    await deleteMembership(plan.id)
    ElMessage.success('删除成功')
    fetchPlans()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

function resetForm() {
  Object.assign(planForm, {
    id: null,
    name: '',
    duration: 30,
    price: 0,
    description: '',
    status: 1
  })
  if (planFormRef.value) {
    planFormRef.value.clearValidate()
  }
}

async function submitForm() {
  if (!planFormRef.value) return
  
  try {
    await planFormRef.value.validate()
    submitting.value = true

    const formData = {
      name: planForm.name.trim(),
      duration: planForm.duration,
      price: planForm.price,
      description: planForm.description.trim() || '',
      status: planForm.status
    }

    if (isEditing.value) {
      await updateMembership(planForm.id, formData)
      ElMessage.success('更新成功')
    } else {
      await createMembership(formData)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchPlans()
  } catch (error) {
    if (error.errorFields) {
      // 验证失败，不处理
    } else {
      ElMessage.error(isEditing.value ? '更新失败' : '创建失败')
    }
  } finally {
    submitting.value = false
  }
}

onMounted(fetchPlans)

const dialogTitle = computed(() => {
  return isEditing.value ? '编辑套餐' : '新增套餐'
})
</script>

<style scoped>
.header-bar {
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-start;
  align-items: center;
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 6px;
}
</style>
