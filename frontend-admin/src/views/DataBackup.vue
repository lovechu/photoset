<template>
  <div class="data-backup">
    <div class="backup-header">
      <h2>数据备份</h2>
      <el-button type="primary" @click="handleCreateBackup" :loading="creating">
        <el-icon><Download /></el-icon>
        创建备份
      </el-button>
    </div>

    <!-- 备份说明 -->
    <el-alert
      title="备份说明"
      type="info"
      :closable="false"
      show-icon
      style="margin-bottom: 20px"
    >
      <template #default>
        <ul style="margin: 0; padding-left: 20px;">
          <li>备份会导出整个数据库的结构和数据</li>
          <li>备份文件保存在服务器的 <code>./backups</code> 目录</li>
          <li>建议定期备份，特别是在系统升级或数据迁移前</li>
          <li>恢复备份需要通过命令行手动执行：<code>mysql -u用户名 -p密码 数据库名 < 备份文件.sql</code></li>
        </ul>
      </template>
    </el-alert>

    <!-- 备份列表 -->
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>备份列表</span>
          <el-button type="text" @click="fetchBackups" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <el-table :data="backupList" v-loading="loading" stripe style="width: 100%" border>
        <el-table-column prop="filename" label="文件名" min-width="200" show-overflow-tooltip />
        <el-table-column label="大小" width="120" align="center">
          <template #default="{ row }">
            {{ row.size_str }}
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180" align="center">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="handleDownload(row)">
              下载
            </el-button>
            <el-popconfirm
              title="确认删除该备份文件？删除后不可恢复。"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="backupList.length === 0 && !loading" class="empty-state">
        <el-empty description="暂无备份文件">
          <el-button type="primary" @click="handleCreateBackup">创建第一个备份</el-button>
        </el-empty>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { createBackup, getBackupList, downloadBackup, deleteBackup } from '@/api'
import { ElMessage } from 'element-plus'
import { Download, Refresh } from '@element-plus/icons-vue'

const loading = ref(false)
const creating = ref(false)
const backupList = ref([])

function formatTime(t) {
  if (!t) return ''
  const ts = Number(t)
  if (ts < 1e12) return new Date(ts * 1000).toLocaleString('zh-CN')
  return new Date(ts).toLocaleString('zh-CN')
}

// 获取备份列表
async function fetchBackups() {
  loading.value = true
  try {
    const res = await getBackupList()
    backupList.value = res.data?.list || []
  } catch {
    ElMessage.error('获取备份列表失败')
  } finally {
    loading.value = false
  }
}

// 创建备份
async function handleCreateBackup() {
  creating.value = true
  try {
    const res = await createBackup()
    ElMessage.success('备份创建成功')
    fetchBackups()
  } catch (error) {
    ElMessage.error('创建备份失败: ' + (error.message || '未知错误'))
  } finally {
    creating.value = false
  }
}

// 下载备份
async function handleDownload(row) {
  try {
    const res = await downloadBackup(row.filename)
    const blob = new Blob([res], { type: 'application/octet-stream' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = row.filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('下载成功')
  } catch {
    ElMessage.error('下载失败')
  }
}

// 删除备份
async function handleDelete(row) {
  try {
    await deleteBackup(row.filename)
    ElMessage.success('备份已删除')
    fetchBackups()
  } catch {
    ElMessage.error('删除失败')
  }
}

onMounted(fetchBackups)
</script>

<style scoped>
.data-backup {
  padding: 20px;
}

.backup-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.backup-header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-state {
  padding: 40px 0;
}

code {
  background-color: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  color: #606266;
}
</style>
