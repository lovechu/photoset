<template>
  <div class="community-categories">
    <!-- 操作栏 -->
    <div class="header-bar">
      <div class="header-left">
        <el-button type="primary" @click="handleCreate" :icon="Plus">
          新增分类
        </el-button>
      </div>
    </div>

    <!-- 分类表格 -->
    <el-table :data="categoryList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="name" label="分类名称" min-width="120" />
      <el-table-column prop="key" label="标识" min-width="120" />
      <el-table-column label="颜色" width="80" align="center">
        <template #default="{ row }">
          <span
            class="color-dot"
            :style="{ backgroundColor: row.color || '#409EFF' }"
          />
        </template>
      </el-table-column>
      <el-table-column prop="sort_order" label="排序号" width="90" align="center" />
      <el-table-column label="帖子数" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.post_count > 0 ? 'warning' : 'info'" size="small">
            {{ row.post_count ?? 0 }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170" align="center">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="260" align="center" fixed="right">
        <template #default="{ row }">
          <div class="action-buttons">
            <el-button
              size="small"
              :icon="ArrowUp"
              :disabled="isFirst(row)"
              @click="moveUp(row)"
            />
            <el-button
              size="small"
              :icon="ArrowDown"
              :disabled="isLast(row)"
              @click="moveDown(row)"
            />
            <el-button type="primary" size="small" plain @click="handleEdit(row)">
              编辑
            </el-button>
            <el-popconfirm
              title="确定要删除该分类吗？"
              confirm-button-text="确认删除"
              cancel-button-text="取消"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button
                  type="danger"
                  size="small"
                  plain
                  :loading="deletingId === row.id"
                >
                  删除
                </el-button>
              </template>
            </el-popconfirm>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <!-- 空状态 -->
    <el-empty v-if="!loading && categoryList.length === 0" description="暂无分类" />

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEditing ? '编辑分类' : '新增分类'"
      width="500px"
      destroy-on-close
      @close="resetForm"
    >
      <el-form ref="categoryFormRef" :model="categoryForm" :rules="categoryRules" label-width="100px">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="categoryForm.name" placeholder="例如：讨论" />
        </el-form-item>
        <el-form-item label="标识" prop="key">
          <el-input
            v-model="categoryForm.key"
            placeholder="例如：discussion"
            :disabled="isEditing"
          />
          <span class="form-tip">
            {{ isEditing ? '标识不可修改' : '小写字母开头，仅含小写字母、数字和下划线' }}
          </span>
        </el-form-item>
        <el-form-item label="排序号" prop="sort_order">
          <el-input-number
            v-model="categoryForm.sort_order"
            :min="0"
            :max="9999"
            controls-position="right"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="颜色" prop="color">
          <el-color-picker v-model="categoryForm.color" show-alpha />
          <span style="margin-left: 12px; color: #909399; font-size: 12px;">{{ categoryForm.color }}</span>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="categoryForm.description"
            type="textarea"
            :rows="2"
            placeholder="分类描述（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">
          {{ isEditing ? '保存修改' : '确认新增' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import {
  getCommunityCategories,
  createCommunityCategory,
  updateCommunityCategory,
  deleteCommunityCategory,
  sortCommunityCategories
} from '@/api/community'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, ArrowUp, ArrowDown } from '@element-plus/icons-vue'

// 列表状态
const loading = ref(false)
const categoryList = ref([])

// 操作状态
const deletingId = ref(null)
const submitting = ref(false)

// 排序位置判断
const isFirst = (row) => {
  return categoryList.value.length > 0 && categoryList.value[0].id === row.id
}
const isLast = (row) => {
  const list = categoryList.value
  return list.length > 0 && list[list.length - 1].id === row.id
}

// 对话框
const dialogVisible = ref(false)
const isEditing = ref(false)
const categoryFormRef = ref(null)
const categoryForm = reactive({
  id: null,
  key: '',
  name: '',
  description: '',
  color: '#409EFF',
  sort_order: 0
})

const categoryRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 1, max: 128, message: '分类名称长度在 1 到 128 个字符', trigger: 'blur' }
  ],
  key: [
    { required: true, message: '请输入分类标识', trigger: 'blur' },
    { pattern: /^[a-z][a-z0-9_]*$/, message: '标识必须以小写字母开头，仅含小写字母、数字和下划线', trigger: 'blur' }
  ]
}

function formatTime(t) {
  if (!t) return ''
  const ts = Number(t)
  if (ts < 1e12) return new Date(ts * 1000).toLocaleString('zh-CN')
  return new Date(ts).toLocaleString('zh-CN')
}

async function fetchCategories() {
  loading.value = true
  try {
    const res = await getCommunityCategories()
    categoryList.value = res.data?.categories || []
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  isEditing.value = false
  resetForm()
  dialogVisible.value = true
}

function handleEdit(row) {
  isEditing.value = true
  Object.assign(categoryForm, {
    id: row.id,
    key: row.key || '',
    name: row.name || '',
    description: row.description || '',
    color: row.color || '#409EFF',
    sort_order: row.sort_order ?? 0
  })
  dialogVisible.value = true
}

function resetForm() {
  Object.assign(categoryForm, {
    id: null,
    key: '',
    name: '',
    description: '',
    color: '#409EFF',
    sort_order: 0
  })
  if (categoryFormRef.value) {
    categoryFormRef.value.clearValidate()
  }
}

async function submitForm() {
  if (!categoryFormRef.value) return
  try {
    await categoryFormRef.value.validate()
    submitting.value = true

    const data = {
      name: categoryForm.name.trim(),
      description: categoryForm.description.trim(),
      color: categoryForm.color,
      sort_order: categoryForm.sort_order
    }

    if (isEditing.value) {
      await updateCommunityCategory(categoryForm.id, data)
      ElMessage.success('更新成功')
    } else {
      data.key = categoryForm.key.trim()
      await createCommunityCategory(data)
      ElMessage.success('新增成功')
    }

    dialogVisible.value = false
    fetchCategories()
  } catch (error) {
    if (error?.errorFields) {
      // 表单验证失败
    } else if (error?.response?.data?.message) {
      ElMessage.error(error.response.data.message)
    }
    // 其他错误已由拦截器处理
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row) {
  // Check for associated posts
  if (row.post_count > 0) {
    ElMessageBox.alert(
      `该分类「${row.name}」下有 ${row.post_count} 篇帖子，无法删除。请先将帖子迁移到其他分类。`,
      '无法删除',
      { confirmButtonText: '知道了', type: 'warning' }
    )
    return
  }

  deletingId.value = row.id
  try {
    await deleteCommunityCategory(row.id)
    ElMessage.success('已删除')
    fetchCategories()
  } catch (error) {
    if (error?.response?.data?.message) {
      ElMessage.error(error.response.data.message)
    }
  } finally {
    deletingId.value = null
  }
}

async function moveUp(row) {
  const list = categoryList.value
  const idx = list.findIndex((c) => c.id === row.id)
  if (idx <= 0) return

  // Swap sort_order values
  const currentSort = list[idx].sort_order
  const aboveSort = list[idx - 1].sort_order

  await sortCommunityCategories([
    { id: list[idx].id, sort_order: aboveSort },
    { id: list[idx - 1].id, sort_order: currentSort }
  ])
  ElMessage.success('排序已更新')
  fetchCategories()
}

async function moveDown(row) {
  const list = categoryList.value
  const idx = list.findIndex((c) => c.id === row.id)
  if (idx < 0 || idx >= list.length - 1) return

  // Swap sort_order values
  const currentSort = list[idx].sort_order
  const belowSort = list[idx + 1].sort_order

  await sortCommunityCategories([
    { id: list[idx].id, sort_order: belowSort },
    { id: list[idx + 1].id, sort_order: currentSort }
  ])
  ElMessage.success('排序已更新')
  fetchCategories()
}

onMounted(fetchCategories)
</script>

<style scoped>
.header-bar {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-left {
  display: flex;
  align-items: center;
}
.action-buttons {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 4px;
}
.color-dot {
  display: inline-block;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 1px solid #dcdfe6;
  vertical-align: middle;
}
.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}
</style>
