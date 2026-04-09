<template>
  <div class="tag-manage">
    <div class="header-bar">
      <el-button type="primary" @click="handleCreate" :icon="Plus">
        新增标签
      </el-button>
      <el-input
        v-model="filterKeyword"
        placeholder="搜索标签名称"
        clearable
        @clear="fetchTags"
        @keyup.enter="fetchTags"
        style="width: 240px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <el-table :data="tagList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="name" label="标签名称" min-width="180" />
      <el-table-column prop="slug" label="URL标识" min-width="150">
        <template #default="{ row }">
          <el-tag size="small">/tag/{{ row.slug }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="使用次数" width="100" align="center">
        <template #default="{ row }">
          <el-tag type="info" size="small">{{ row.photo_count || 0 }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column label="创建时间" width="170" align="center">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="更新时间" width="170" align="center">
        <template #default="{ row }">
          {{ formatTime(row.updated_at) || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" align="center" fixed="right">
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
            <el-popconfirm
              title="确定要删除这个标签吗？删除后可能影响套图标签。"
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

    <!-- 新建/编辑标签对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      @close="resetForm"
    >
      <el-form ref="tagFormRef" :model="tagForm" :rules="tagRules" label-width="80px">
        <el-form-item label="标签名称" prop="name">
          <el-input v-model="tagForm.name" placeholder="请输入标签名称" />
        </el-form-item>
        <el-form-item label="URL标识" prop="slug">
          <el-input v-model="tagForm.slug" placeholder="请输入URL标识（英文小写字母和连字符）" />
          <div class="slug-hint">
            <span class="slug-preview">预览: /tag/{{ tagForm.slug || 'slug' }}</span>
          </div>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="tagForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入标签描述"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ isEditing ? '保存修改' : '创建标签' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, reactive, computed } from 'vue'
import { getTagList, createTag, updateTag, deleteTag } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'

const loading = ref(false)
const tagList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const filterKeyword = ref('')
const submitting = ref(false)
const dialogVisible = ref(false)
const isEditing = ref(false)

const tagForm = reactive({
  id: null,
  name: '',
  slug: '',
  description: ''
})

const tagRules = {
  name: [
    { required: true, message: '请输入标签名称', trigger: 'blur' },
    { min: 2, max: 50, message: '标签名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  slug: [
    { required: true, message: '请输入URL标识', trigger: 'blur' },
    { pattern: /^[a-z0-9-]+$/, message: '只能包含小写字母、数字和连字符', trigger: 'blur' },
    { min: 2, max: 30, message: 'URL标识长度在 2 到 30 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 200, message: '描述不能超过 200 个字符', trigger: 'blur' }
  ]
}

const tagFormRef = ref()

function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

async function fetchTags() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value
    }

    if (filterKeyword.value) {
      params.keyword = filterKeyword.value
    }

    const res = await getTagList(params)
    tagList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('获取标签列表失败:', error)
    ElMessage.error('获取标签列表失败')
  } finally {
    loading.value = false
  }
}

function handleSizeChange(val) {
  pageSize.value = val
  currentPage.value = 1
  fetchTags()
}

function handleCurrentChange(val) {
  currentPage.value = val
  fetchTags()
}

function handleCreate() {
  isEditing.value = false
  resetForm()
  dialogVisible.value = true
}

function handleEdit(tag) {
  isEditing.value = true
  Object.assign(tagForm, {
    id: tag.id,
    name: tag.name,
    slug: tag.slug,
    description: tag.description || ''
  })
  dialogVisible.value = true
}

async function handleDelete(tag) {
  try {
    if ((tag.photo_count || 0) > 0) {
      const confirm = await ElMessageBox.confirm(
        `标签 "${tag.name}" 正在被 ${tag.photo_count} 个套图使用，确定要删除吗？这可能导致套图丢失标签。`,
        '警告',
        {
          confirmButtonText: '强制删除',
          cancelButtonText: '取消',
          type: 'warning',
          dangerouslyUseHTMLString: true
        }
      )
    }

    await deleteTag(tag.id)
    ElMessage.success('删除成功')
    fetchTags()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

function resetForm() {
  Object.assign(tagForm, {
    id: null,
    name: '',
    slug: '',
    description: ''
  })
  if (tagFormRef.value) {
    tagFormRef.value.clearValidate()
  }
}

async function submitForm() {
  if (!tagFormRef.value) return
  
  try {
    await tagFormRef.value.validate()
    submitting.value = true

    const formData = {
      name: tagForm.name.trim(),
      slug: tagForm.slug.trim(),
      description: tagForm.description.trim() || null
    }

    if (isEditing.value) {
      await updateTag(tagForm.id, formData)
      ElMessage.success('更新成功')
    } else {
      await createTag(formData)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchTags()
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

// 监听筛选条件变化
watch([filterKeyword], () => {
  currentPage.value = 1
  fetchTags()
})

onMounted(fetchTags)

const dialogTitle = computed(() => {
  return isEditing.value ? '编辑标签' : '新增标签'
})
</script>

<style scoped>
.header-bar {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 8px;
}
</style>