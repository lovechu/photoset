<template>
  <div class="category-manage">
    <div class="header-bar">
      <el-button type="primary" @click="handleCreate" :icon="Plus">新增分类</el-button>
      <el-input
        v-model="filterKeyword"
        placeholder="搜索分类名称"
        clearable
        @clear="fetchData"
        @keyup.enter="fetchData"
        style="width: 240px"
      >
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
    </div>

    <el-table :data="list" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="name" label="分类名称" min-width="140" />
      <el-table-column prop="slug" label="标识(Slug)" min-width="130">
        <template #default="{ row }">
          <el-tag size="small" type="info">{{ row.slug }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip />
      <el-table-column prop="sort_order" label="排序" width="80" align="center" />
      <el-table-column label="创建时间" width="170" align="center">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="160" align="center" fixed="right">
        <template #default="{ row }">
          <div class="action-buttons">
            <el-button type="primary" size="small" plain @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定要删除这个分类吗？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" size="small" plain>删除</el-button>
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
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="fetchData"
        @current-change="fetchData"
      />
    </div>

    <!-- 新建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" @close="resetForm">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="如：自然风光" />
        </el-form-item>
        <el-form-item label="标识 Slug" prop="slug">
          <el-input v-model="form.slug" placeholder="如：nature（英文小写+连字符）" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" maxlength="200" show-word-limit placeholder="可选" />
        </el-form-item>
        <el-form-item label="排序权重">
          <el-input-number v-model="form.sort_order" :min="0" :max="9999" />
          <span style="margin-left:8px;color:#909399;font-size:12px">数值越大越靠前</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ isEditing ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { getCategoryList, createCategory, updateCategory, deleteCategory } from '@/api'
import { ElMessage } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const filterKeyword = ref('')

const dialogVisible = ref(false)
const isEditing = ref(false)
const submitting = ref(false)
const formRef = ref()

const form = reactive({ id: null, name: '', slug: '', description: '', sort_order: 0 })

const rules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度 2-50 个字符', trigger: 'blur' }
  ],
  slug: [
    { required: true, message: '请输入标识', trigger: 'blur' },
    { pattern: /^[a-z0-9-]+$/, message: '只能包含小写字母、数字和连字符', trigger: 'blur' }
  ]
}

const dialogTitle = computed(() => isEditing.value ? '编辑分类' : '新增分类')

function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getCategoryList({ page: currentPage.value, page_size: pageSize.value, keyword: filterKeyword.value })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    ElMessage.error('获取分类列表失败')
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
  Object.assign(form, { id: row.id, name: row.name, slug: row.slug, description: row.description || '', sort_order: row.sort_order })
  dialogVisible.value = true
}

async function handleDelete(row) {
  try {
    await deleteCategory(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

function resetForm() {
  Object.assign(form, { id: null, name: '', slug: '', description: '', sort_order: 0 })
  if (formRef.value) formRef.value.clearValidate()
}

async function submitForm() {
  if (!formRef.value) return
  try { await formRef.value.validate() } catch { return }
  submitting.value = true
  try {
    const data = { name: form.name.trim(), slug: form.slug.trim(), description: form.description.trim() || null, sort_order: form.sort_order }
    if (isEditing.value) {
      await updateCategory(form.id, data)
      ElMessage.success('更新成功')
    } else {
      await createCategory(data)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (e) {
    // 全局 request.js 会处理错误提示
  } finally {
    submitting.value = false
  }
}

watch([filterKeyword], () => { currentPage.value = 1; fetchData() })
onMounted(fetchData)
</script>

<style scoped>
.header-bar { margin-bottom: 20px; display: flex; justify-content: space-between; align-items: center; }
.pagination { margin-top: 24px; display: flex; justify-content: center; }
.action-buttons { display: flex; justify-content: center; gap: 8px; }
</style>
