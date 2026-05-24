<template>
  <div class="community-keywords">
    <!-- 操作栏 -->
    <div class="header-bar">
      <div class="header-left">
        <el-button type="primary" @click="handleCreate" :icon="Plus">
          新增敏感词
        </el-button>
        <el-button
          type="warning"
          plain
          @click="handleReload"
          :loading="reloading"
          style="margin-left: 8px"
        >
          重新加载
        </el-button>
      </div>
      <div class="header-right">
        <el-select v-model="filterActive" placeholder="启用状态" clearable @change="handleFilterChange" style="width: 140px">
          <el-option label="全部" value="" />
          <el-option label="已启用" :value="true" />
          <el-option label="已禁用" :value="false" />
        </el-select>
      </div>
    </div>

    <!-- 敏感词表格 -->
    <el-table :data="keywordList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="word" label="敏感词" min-width="160" />
      <el-table-column prop="replacement" label="替换词" min-width="160">
        <template #default="{ row }">
          {{ row.replacement || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="是否启用" width="100" align="center">
        <template #default="{ row }">
          <el-switch
            :model-value="row.is_active"
            :loading="switchingId === row.id"
            @change="(val) => handleToggleActive(row, val)"
          />
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170" align="center">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="180" align="center" fixed="right">
        <template #default="{ row }">
          <div class="action-buttons">
            <el-button type="primary" size="small" plain @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              plain
              @click="handleDelete(row)"
              :loading="deletingId === row.id"
            >
              删除
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <!-- 空状态 -->
    <el-empty v-if="!loading && keywordList.length === 0" description="暂无数据" />

    <!-- 分页 -->
    <div class="pagination-bar">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @current-change="fetchKeywords"
        @size-change="handleSizeChange"
      />
    </div>

    <!-- 新增/编辑敏感词对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEditing ? '编辑敏感词' : '新增敏感词'"
      width="480px"
      destroy-on-close
      @close="resetForm"
    >
      <el-form ref="keywordFormRef" :model="keywordForm" :rules="keywordRules" label-width="80px">
        <el-form-item label="敏感词" prop="word">
          <el-input v-model="keywordForm.word" placeholder="请输入敏感词" />
        </el-form-item>
        <el-form-item label="替换词" prop="replacement">
          <el-input v-model="keywordForm.replacement" placeholder="请输入替换词（可选）" />
        </el-form-item>
        <el-form-item v-if="isEditing" label="是否启用" prop="is_active">
          <el-switch v-model="keywordForm.is_active" />
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
import { ref, onMounted, reactive } from 'vue'
import { getKeywords, addKeyword, updateKeyword, deleteKeyword, reloadKeywords } from '@/api/community'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

// 列表状态
const loading = ref(false)
const keywordList = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filterActive = ref('')

// 操作状态
const deletingId = ref(null)
const switchingId = ref(null)
const submitting = ref(false)
const reloading = ref(false)

// 对话框
const dialogVisible = ref(false)
const isEditing = ref(false)
const keywordFormRef = ref(null)
const keywordForm = reactive({
  id: null,
  word: '',
  replacement: '',
  is_active: true
})

const keywordRules = {
  word: [
    { required: true, message: '请输入敏感词', trigger: 'blur' },
    { min: 1, max: 100, message: '敏感词长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  replacement: [
    { max: 100, message: '替换词不能超过 100 个字符', trigger: 'blur' }
  ]
}

function formatTime(t) {
  if (!t) return ''
  const ts = Number(t)
  if (ts < 1e12) return new Date(ts * 1000).toLocaleString('zh-CN')
  return new Date(ts).toLocaleString('zh-CN')
}

function handleFilterChange() {
  page.value = 1
  fetchKeywords()
}

function handleSizeChange() {
  page.value = 1
  fetchKeywords()
}

async function fetchKeywords() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filterActive.value !== '' && filterActive.value !== undefined) {
      params.is_active = filterActive.value
    }
    const res = await getKeywords(params)
    keywordList.value = res.data?.list || []
    total.value = res.data?.total || 0
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
  Object.assign(keywordForm, {
    id: row.id,
    word: row.word || '',
    replacement: row.replacement || '',
    is_active: !!row.is_active
  })
  dialogVisible.value = true
}

function resetForm() {
  Object.assign(keywordForm, {
    id: null,
    word: '',
    replacement: '',
    is_active: true
  })
  if (keywordFormRef.value) {
    keywordFormRef.value.clearValidate()
  }
}

async function submitForm() {
  if (!keywordFormRef.value) return
  try {
    await keywordFormRef.value.validate()
    submitting.value = true

    const data = {
      word: keywordForm.word.trim(),
      replacement: keywordForm.replacement.trim() || ''
    }
    if (isEditing.value) {
      data.is_active = keywordForm.is_active
      await updateKeyword(keywordForm.id, data)
      ElMessage.success('更新成功')
    } else {
      await addKeyword(data)
      ElMessage.success('新增成功')
    }

    dialogVisible.value = false
    fetchKeywords()
  } catch (error) {
    if (error?.errorFields) {
      // 表单验证失败，不处理
    }
    // 其他错误已由拦截器处理
  } finally {
    submitting.value = false
  }
}

async function handleToggleActive(row, val) {
  switchingId.value = row.id
  try {
    await updateKeyword(row.id, {
      word: row.word,
      replacement: row.replacement || '',
      is_active: val
    })
    ElMessage.success(val ? '已启用' : '已禁用')
    fetchKeywords()
  } catch {
    // 错误已由拦截器处理
  } finally {
    switchingId.value = null
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(
      `确定要删除敏感词「${row.word}」吗？此操作不可撤销。`,
      '删除确认',
      { confirmButtonText: '确认删除', cancelButtonText: '取消', type: 'warning' }
    )
    deletingId.value = row.id
    await deleteKeyword(row.id)
    ElMessage.success('已删除')
    fetchKeywords()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      // 错误已由拦截器处理
    }
  } finally {
    deletingId.value = null
  }
}

async function handleReload() {
  reloading.value = true
  try {
    await reloadKeywords()
    ElMessage.success('敏感词已重新加载')
  } catch {
    // 错误已由拦截器处理
  } finally {
    reloading.value = false
  }
}

onMounted(fetchKeywords)
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
.header-right {
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
  gap: 8px;
}
</style>
