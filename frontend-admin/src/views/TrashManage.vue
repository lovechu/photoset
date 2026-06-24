<template>
  <div class="trash-manage">
    <div class="page-header">
      <h2>回收站管理</h2>
      <p class="page-desc">管理已被删除的套图，可恢复或永久删除</p>
    </div>

    <el-card shadow="never">
      <div class="table-toolbar">
        <el-button type="primary" @click="fetchList" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button @click="clearAll" :disabled="list.length === 0" type="danger" plain>
          清空回收站
        </el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="list"
        stripe
        style="width: 100%; margin-top: 16px"
        empty-text="回收站为空"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="封面" width="100">
          <template #default="{ row }">
            <el-image
              v-if="row.cover"
              :src="row.cover"
              fit="cover"
              style="width: 60px; height: 60px; border-radius: 4px"
              preview-teleported
              :preview-src-list="[row.cover]"
            />
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ row.status || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_id" label="用户ID" width="80" />
        <el-table-column label="删除时间" width="180">
          <template #default="{ row }">
            {{ row.deleted_at ? new Date(row.deleted_at).toLocaleString() : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="photo_count" label="图片数" width="80" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button
              type="success"
              size="small"
              @click="handleRestore(row)"
              :loading="restoringId === row.id"
            >
              恢复
            </el-button>
            <el-popconfirm
              title="永久删除后不可恢复，确定要删除吗？"
              confirm-button-text="确定"
              cancel-button-text="取消"
              @confirm="handlePermanentDelete(row)"
            >
              <template #reference>
                <el-button
                  type="danger"
                  size="small"
                  :loading="deletingId === row.id"
                >
                  永久删除
                </el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchList"
          @size-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAdminTrashList, adminRestorePhotoset, adminPermanentDeletePhotoset } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'

const loading = ref(false)
const list = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const restoringId = ref(null)
const deletingId = ref(null)

async function fetchList() {
  loading.value = true
  try {
    const res = await getAdminTrashList({ page: page.value, page_size: pageSize.value })
    const data = res.data
    if (data?.list) {
      list.value = data.list
      total.value = data.total || 0
    } else if (Array.isArray(data)) {
      list.value = data
      total.value = data.length
    } else {
      list.value = []
      total.value = 0
    }
  } catch {
    ElMessage.error('获取回收站列表失败')
  } finally {
    loading.value = false
  }
}

async function handleRestore(row) {
  restoringId.value = row.id
  try {
    await adminRestorePhotoset(row.id)
    ElMessage.success('恢复成功')
    fetchList()
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '恢复失败')
  } finally {
    restoringId.value = null
  }
}

async function handlePermanentDelete(row) {
  deletingId.value = row.id
  try {
    await adminPermanentDeletePhotoset(row.id)
    ElMessage.success('永久删除成功')
    fetchList()
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '删除失败')
  } finally {
    deletingId.value = null
  }
}

async function clearAll() {
  try {
    await ElMessageBox.confirm(
      '清空回收站将永久删除所有已删除的套图，不可恢复！确定要继续吗？',
      '确认清空',
      { type: 'error', confirmButtonText: '确定清空', cancelButtonText: '取消' }
    )
  } catch {
    return
  }

  loading.value = true
  try {
    for (const item of list.value) {
      await adminPermanentDeletePhotoset(item.id)
    }
    ElMessage.success('回收站已清空')
    fetchList()
  } catch {
    ElMessage.error('清空失败，部分套图可能已被清除')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.trash-manage {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0 0 4px;
  font-size: 20px;
  color: #303133;
}

.page-desc {
  margin: 0;
  color: #909399;
  font-size: 13px;
}

.table-toolbar {
  display: flex;
  gap: 10px;
}

.pagination-wrap {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
