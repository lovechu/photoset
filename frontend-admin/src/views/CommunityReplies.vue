<template>
  <div class="community-replies">
    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-input
        v-model="filterPostId"
        placeholder="按帖子ID筛选"
        clearable
        @clear="handleFilterChange"
        @keyup.enter="handleFilterChange"
        style="width: 200px"
      />
      <el-button type="primary" plain @click="handleFilterChange" style="margin-left: 8px">
        搜索
      </el-button>
    </div>

    <!-- 回帖表格 -->
    <el-table :data="replyList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column label="内容" min-width="220">
        <template #default="{ row }">
          <span :title="row.content">{{ truncate(row.content, 50) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="作者" width="120" align="center">
        <template #default="{ row }">
          {{ row.author?.nickname || row.author?.username || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="所属帖子" min-width="180" show-overflow-tooltip>
        <template #default="{ row }">
          {{ row.post?.title || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="点赞" width="70" align="center">
        <template #default="{ row }">{{ row.like_count ?? 0 }}</template>
      </el-table-column>
      <el-table-column label="创建时间" width="170" align="center">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="100" align="center" fixed="right">
        <template #default="{ row }">
          <el-button
            type="danger"
            size="small"
            plain
            @click="handleDelete(row)"
            :loading="deletingId === row.id"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 空状态 -->
    <el-empty v-if="!loading && replyList.length === 0" description="暂无数据" />

    <!-- 分页 -->
    <div class="pagination-bar">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @current-change="fetchReplies"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCommunityReplies, deleteReply } from '@/api/community'
import { ElMessage, ElMessageBox } from 'element-plus'

// 列表状态
const loading = ref(false)
const replyList = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const filterPostId = ref('')

// 操作 loading 状态
const deletingId = ref(null)

function truncate(text, maxLen) {
  if (!text) return ''
  return text.length > maxLen ? text.slice(0, maxLen) + '...' : text
}

function formatTime(t) {
  if (!t) return ''
  const ts = Number(t)
  if (ts < 1e12) return new Date(ts * 1000).toLocaleString('zh-CN')
  return new Date(ts).toLocaleString('zh-CN')
}

function handleFilterChange() {
  page.value = 1
  fetchReplies()
}

function handleSizeChange() {
  page.value = 1
  fetchReplies()
}

async function fetchReplies() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filterPostId.value) params.post_id = filterPostId.value
    const res = await getCommunityReplies(params)
    replyList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(
      '确定要删除该回帖吗？此操作不可撤销。',
      '删除确认',
      { confirmButtonText: '确认删除', cancelButtonText: '取消', type: 'warning' }
    )
    deletingId.value = row.id
    await deleteReply(row.id)
    ElMessage.success('已删除')
    fetchReplies()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      // 错误已由拦截器处理
    }
  } finally {
    deletingId.value = null
  }
}

onMounted(fetchReplies)
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
