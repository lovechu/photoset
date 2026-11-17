<template>
  <div class="collection-list-page">
    <div class="page-header">
      <h2>我的合集</h2>
      <el-button type="primary" :icon="Plus" @click="openCreateDialog">创建合集</el-button>
    </div>

    <div v-loading="loading">
      <el-empty v-if="!loading && list.length === 0" description="还没有合集，创建一个吧" :image-size="100" />
      <div v-else class="collection-grid">
        <el-card v-for="item in list" :key="item.id" class="collection-card" shadow="hover">
          <div class="cover" @click="goDetail(item)">
            <el-image v-if="item.cover_url" :src="item.cover_url" fit="cover" class="cover-img" />
            <div v-else class="cover-placeholder">
              <el-icon :size="40"><Folder /></el-icon>
            </div>
          </div>
          <div class="card-body">
            <div class="card-title" @click="goDetail(item)">{{ item.name }}</div>
            <div class="card-meta">
              <span>{{ item.item_count || 0 }} 个套图</span>
              <el-tag size="small" :type="item.is_public ? 'success' : 'info'">
                {{ item.is_public ? '公开' : '私密' }}
              </el-tag>
            </div>
            <p class="card-desc">{{ item.description || '暂无描述' }}</p>
            <div class="card-actions">
              <el-button link size="small" :icon="Edit" @click="openEditDialog(item)">编辑</el-button>
              <el-button link size="small" type="danger" :icon="Delete" @click="removeCollection(item)">删除</el-button>
            </div>
          </div>
        </el-card>
      </div>

      <div class="load-more" v-if="list.length < total">
        <el-button @click="loadMore" :loading="loadingMore">加载更多</el-button>
      </div>
    </div>

    <!-- 创建/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑合集' : '创建合集'" width="440px">
      <el-form label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" maxlength="100" placeholder="给合集起个名字" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" maxlength="500" placeholder="合集描述（选填）" />
        </el-form-item>
        <el-form-item label="公开">
          <el-switch v-model="form.is_public" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Folder, Edit, Delete } from '@element-plus/icons-vue'
import { listCollections, createCollection, updateCollection, deleteCollection } from '@/api'

const router = useRouter()
const list = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const loadingMore = ref(false)

const dialogVisible = ref(false)
const editingId = ref(null)
const submitting = ref(false)
const form = ref({ name: '', description: '', is_public: false })

onMounted(() => loadCollections())

async function loadCollections() {
  loading.value = true
  try {
    const res = await listCollections({ page: page.value, page_size: 20 })
    const data = res?.data || {}
    const items = data.list || []
    list.value = page.value === 1 ? items : [...list.value, ...items]
    total.value = data.total || 0
  } catch (e) {
    console.warn('加载合集失败:', e.message)
  } finally {
    loading.value = false
  }
}

function loadMore() {
  page.value++
  loadCollections()
}

function openCreateDialog() {
  editingId.value = null
  form.value = { name: '', description: '', is_public: false }
  dialogVisible.value = true
}

function openEditDialog(item) {
  editingId.value = item.id
  form.value = { name: item.name, description: item.description || '', is_public: !!item.is_public }
  dialogVisible.value = true
}

async function submitForm() {
  if (!form.value.name.trim()) {
    ElMessage.warning('请输入合集名称')
    return
  }
  submitting.value = true
  try {
    if (editingId.value) {
      await updateCollection(editingId.value, { ...form.value })
      const item = list.value.find(c => c.id === editingId.value)
      if (item) Object.assign(item, form.value)
      ElMessage.success('已更新')
    } else {
      await createCollection({ ...form.value })
      page.value = 1
      await loadCollections()
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
  } catch (e) {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

async function removeCollection(item) {
  try {
    await ElMessageBox.confirm(`确定删除合集「${item.name}」？`, '提示', { type: 'warning' })
    await deleteCollection(item.id)
    list.value = list.value.filter(c => c.id !== item.id)
    total.value--
    ElMessage.success('已删除')
  } catch (e) { /* 取消 */ }
}

function goDetail(item) {
  router.push(`/collection/${item.id}`)
}
</script>

<style scoped>
.collection-list-page { max-width: 1100px; margin: 0 auto; padding: 24px 0 60px; }
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 24px; }
.page-header h2 { margin: 0; font-size: 22px; }
.collection-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 20px; }
.collection-card { overflow: hidden; }
.cover { height: 150px; cursor: pointer; background: #f5f7fa; }
.cover-img { width: 100%; height: 100%; }
.cover-placeholder { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: #c0c4cc; }
.card-title { font-size: 16px; font-weight: 600; cursor: pointer; margin-bottom: 8px; }
.card-title:hover { color: #409eff; }
.card-meta { display: flex; align-items: center; justify-content: space-between; color: #999; font-size: 13px; }
.card-desc { color: #666; font-size: 13px; margin: 8px 0; min-height: 20px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.card-actions { display: flex; gap: 12px; border-top: 1px solid #f0f0f0; padding-top: 8px; }
.load-more { text-align: center; padding: 24px; }
</style>
