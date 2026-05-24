<template>
  <div class="community-posts">
    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-select v-model="filterStatus" placeholder="状态筛选" clearable @change="handleFilterChange" style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="已审核" value="approved" />
        <el-option label="待审核" value="pending" />
        <el-option label="已拒绝" value="rejected" />
      </el-select>
    </div>

    <!-- 帖子表格 -->
    <el-table :data="postList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="title" label="标题" min-width="180" show-overflow-tooltip />
      <el-table-column label="作者" width="120" align="center">
        <template #default="{ row }">
          {{ row.author?.nickname || row.author?.username || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="分类" width="100" align="center">
        <template #default="{ row }">
          <el-tag size="small" type="info">{{ row.category || '-' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="可见性" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.visibility === 'public' ? 'success' : 'info'" size="small">
            {{ row.visibility === 'public' ? '公开' : '私密' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small">
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="点赞" width="70" align="center">
        <template #default="{ row }">{{ row.like_count ?? 0 }}</template>
      </el-table-column>
      <el-table-column label="回复" width="70" align="center">
        <template #default="{ row }">{{ row.reply_count ?? 0 }}</template>
      </el-table-column>
      <el-table-column label="创建时间" width="170" align="center">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="280" align="center" fixed="right">
        <template #default="{ row }">
          <div class="action-buttons">
            <el-button
              :type="row.is_pinned ? 'warning' : 'primary'"
              size="small"
              plain
              @click="handleTogglePin(row)"
              :loading="pinningId === row.id"
            >
              {{ row.is_pinned ? '取消置顶' : '置顶' }}
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="success"
              size="small"
              plain
              @click="handleApprove(row)"
              :loading="approvingId === row.id"
            >
              通过
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="danger"
              size="small"
              plain
              @click="handleReject(row)"
              :loading="rejectingId === row.id"
            >
              拒绝
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
    <el-empty v-if="!loading && postList.length === 0" description="暂无数据" />

    <!-- 分页 -->
    <div class="pagination-bar">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @current-change="fetchPosts"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCommunityPosts, togglePostPin, updatePostStatus, deletePost } from '@/api/community'
import { ElMessage, ElMessageBox } from 'element-plus'

// 列表状态
const loading = ref(false)
const postList = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filterStatus = ref('')

// 操作 loading 状态（防重复点击）
const pinningId = ref(null)
const approvingId = ref(null)
const rejectingId = ref(null)
const deletingId = ref(null)

// 状态映射
const statusMap = { approved: '已审核', pending: '待审核', rejected: '已拒绝' }
const statusTagMap = { approved: 'success', pending: 'warning', rejected: 'danger' }

function statusLabel(s) { return statusMap[s] || s || '-' }
function statusTagType(s) { return statusTagMap[s] || 'info' }

function formatTime(t) {
  if (!t) return ''
  const ts = Number(t)
  if (ts < 1e12) return new Date(ts * 1000).toLocaleString('zh-CN')
  return new Date(ts).toLocaleString('zh-CN')
}

function handleFilterChange() {
  page.value = 1
  fetchPosts()
}

function handleSizeChange() {
  page.value = 1
  fetchPosts()
}

async function fetchPosts() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filterStatus.value) params.status = filterStatus.value
    const res = await getCommunityPosts(params)
    postList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}

async function handleTogglePin(row) {
  try {
    pinningId.value = row.id
    await togglePostPin(row.id)
    ElMessage.success(row.is_pinned ? '已取消置顶' : '已置顶')
    fetchPosts()
  } catch {
    // 错误已由拦截器处理
  } finally {
    pinningId.value = null
  }
}

async function handleApprove(row) {
  try {
    await ElMessageBox.confirm(
      `确定要审核通过帖子「${row.title}」吗？`,
      '审核确认',
      { confirmButtonText: '确认通过', cancelButtonText: '取消', type: 'success' }
    )
    approvingId.value = row.id
    await updatePostStatus(row.id, { status: 'approved' })
    ElMessage.success('已审核通过')
    fetchPosts()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      // 错误已由拦截器处理
    }
  } finally {
    approvingId.value = null
  }
}

async function handleReject(row) {
  try {
    await ElMessageBox.confirm(
      `确定要拒绝帖子「${row.title}」吗？`,
      '拒绝确认',
      { confirmButtonText: '确认拒绝', cancelButtonText: '取消', type: 'warning' }
    )
    rejectingId.value = row.id
    await updatePostStatus(row.id, { status: 'rejected' })
    ElMessage.success('已拒绝')
    fetchPosts()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      // 错误已由拦截器处理
    }
  } finally {
    rejectingId.value = null
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(
      `确定要删除帖子「${row.title}」吗？此操作不可撤销。`,
      '删除确认',
      { confirmButtonText: '确认删除', cancelButtonText: '取消', type: 'warning' }
    )
    deletingId.value = row.id
    await deletePost(row.id)
    ElMessage.success('已删除')
    fetchPosts()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      // 错误已由拦截器处理
    }
  } finally {
    deletingId.value = null
  }
}

onMounted(fetchPosts)
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
.action-buttons {
  display: flex;
  justify-content: center;
  gap: 6px;
  flex-wrap: wrap;
}
</style>
