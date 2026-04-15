<template>
  <div class="pages-admin">
    <el-card>
      <template #header>
        <div class="header">
          <span>页面管理</span>
          <el-button type="primary" @click="showCreateDialog">新建页面</el-button>
        </div>
      </template>

      <!-- 搜索和筛选 -->
      <div class="filters">
        <el-input
          v-model="keyword"
          placeholder="搜索标题/标识"
          style="width: 300px;"
          @input="handleSearch"
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select v-model="filterStatus" placeholder="状态" style="width: 120px;" @change="handleFilter">
          <el-option label="全部" value="" />
          <el-option label="已发布" value="published" />
          <el-option label="草稿" value="draft" />
        </el-select>
      </div>

      <!-- 表格 -->
      <el-table :data="pages" style="width: 100%; margin-top: 20px;" v-loading="loading">
        <el-table-column prop="slug" label="标识" width="180" />
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="editPage(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deletePage(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        style="margin-top: 20px; justify-content: center;"
      />
    </el-card>

    <!-- 编辑/新建弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="800px"
      :before-close="closeDialog"
    >
      <el-form :model="form" :rules="rules" label-width="80px" ref="formRef">
        <el-form-item label="页面标识" prop="slug">
          <el-input v-model="form.slug" placeholder="例：about, terms" />
          <div class="form-tip">URL 路径标识，必须为英文单词且唯一</div>
        </el-form-item>
        <el-form-item label="页面标题" prop="title">
          <el-input v-model="form.title" placeholder="页面标题" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio label="published">已发布</el-radio>
            <el-radio label="draft">草稿</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="内容" prop="content_md">
          <el-input
            v-model="form.content_md"
            type="textarea"
            :rows="15"
            placeholder="支持 Markdown 语法"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { getPageList, createPage, updatePage, deletePage as deletePageApi } from '@/api/pages'

// 状态
const pages = ref([])
const loading = ref(false)
const keyword = ref('')
const filterStatus = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 弹窗相关
const dialogVisible = ref(false)
const editingId = ref(null)
const submitting = ref(false)
const formRef = ref(null)

const form = ref({
  slug: '',
  title: '',
  content_md: '',
  status: 'published'
})

const rules = {
  slug: [
    { required: true, message: '请输入页面标识', trigger: 'blur' },
    { pattern: /^[a-z0-9_-]+$/, message: '只能包含小写字母、数字、下划线和中划线', trigger: 'blur' }
  ],
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'blur' }],
  content_md: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const dialogTitle = computed(() => editingId.value ? '编辑页面' : '新建页面')

// 加载数据
const loadPages = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value
    }
    if (keyword.value) params.keyword = keyword.value
    if (filterStatus.value) params.status = filterStatus.value

    const res = await getPageList(params)
    pages.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadPages()
})

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  loadPages()
}
const handleCurrentChange = (val) => {
  page.value = val
  loadPages()
}
const handleSearch = () => {
  page.value = 1
  loadPages()
}
const handleFilter = () => {
  page.value = 1
  loadPages()
}

// 格式化日期
const formatDate = (timestamp) => {
  if (!timestamp) return ''
  const d = new Date(timestamp)
  return d.getFullYear() + '-' + (d.getMonth() + 1).toString().padStart(2, '0') + '-' + d.getDate().toString().padStart(2, '0')
}

// 弹窗
const showCreateDialog = () => {
  editingId.value = null
  form.value = { slug: '', title: '', content_md: '', status: 'published' }
  dialogVisible.value = true
  if (formRef.value) formRef.value.clearValidate()
}

const editPage = (row) => {
  editingId.value = row.id
  form.value = {
    slug: row.slug,
    title: row.title,
    content_md: row.content_md || '',
    status: row.status
  }
  dialogVisible.value = true
  if (formRef.value) formRef.value.clearValidate()
}

const closeDialog = () => {
  dialogVisible.value = false
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (editingId.value) {
        await updatePage(editingId.value, form.value)
        ElMessage.success('更新成功')
      } else {
        await createPage(form.value)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadPages()
    } catch (e) {
      ElMessage.error(editingId.value ? '更新失败' : '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

const deletePage = (row) => {
  ElMessageBox.confirm(`确定删除页面 "${row.title}" 吗？`, '删除确认', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deletePageApi(row.id)
      ElMessage.success('删除成功')
      loadPages()
    } catch (e) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}
</script>

<style scoped>
.pages-admin {
  padding: 0;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filters {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>