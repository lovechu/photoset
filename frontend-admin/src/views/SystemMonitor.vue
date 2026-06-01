<template>
  <div class="system-monitor">
    <div class="monitor-header">
      <h2>系统监控</h2>
      <div class="refresh-controls">
        <el-select v-model="refreshInterval" @change="handleIntervalChange" style="width: 120px">
          <el-option label="手动刷新" :value="0" />
          <el-option label="10秒刷新" :value="10" />
          <el-option label="30秒刷新" :value="30" />
          <el-option label="60秒刷新" :value="60" />
        </el-select>
        <el-button type="primary" @click="fetchStatus" :loading="loading" style="margin-left: 12px">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <div v-loading="loading" class="monitor-content">
      <!-- 服务器信息 -->
      <el-card class="monitor-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>服务器信息</span>
            <el-tag :type="status?.server ? 'success' : 'danger'" size="small">
              {{ status?.server ? '运行中' : '未知' }}
            </el-tag>
          </div>
        </template>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="运行时间">{{ status?.server?.uptime || '-' }}</el-descriptions-item>
          <el-descriptions-item label="Go 版本">{{ status?.server?.go_version || '-' }}</el-descriptions-item>
          <el-descriptions-item label="系统类型">{{ status?.server?.os || '-' }}</el-descriptions-item>
          <el-descriptions-item label="主机名">{{ status?.server?.hostname || '-' }}</el-descriptions-item>
          <el-descriptions-item label="协程数">{{ status?.server?.goroutines || 0 }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 内存使用 -->
      <el-card class="monitor-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>内存使用</span>
            <el-tag type="info" size="small">{{ status?.memory?.gc_count || 0 }} 次 GC</el-tag>
          </div>
        </template>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="当前分配">{{ status?.memory?.alloc_mb || 0 }} MB</el-descriptions-item>
          <el-descriptions-item label="总分配">{{ status?.memory?.total_alloc_mb || 0 }} MB</el-descriptions-item>
          <el-descriptions-item label="系统内存">{{ status?.memory?.sys_mb || 0 }} MB</el-descriptions-item>
          <el-descriptions-item label="GC 次数">{{ status?.memory?.gc_count || 0 }}</el-descriptions-item>
        </el-descriptions>
        <div class="memory-bar" style="margin-top: 16px">
          <div class="bar-label">内存使用率</div>
          <el-progress
            :percentage="memoryPercentage"
            :color="memoryColor"
            :stroke-width="20"
            :text-inside="true"
          />
        </div>
      </el-card>

      <!-- 数据库状态 -->
      <el-card class="monitor-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>数据库状态</span>
            <el-tag :type="dbStatusType" size="small">{{ dbStatusText }}</el-tag>
          </div>
        </template>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="连接状态">{{ dbStatusText }}</el-descriptions-item>
          <el-descriptions-item label="打开连接">{{ status?.database?.open_connections || 0 }}</el-descriptions-item>
          <el-descriptions-item label="使用中">{{ status?.database?.in_use || 0 }}</el-descriptions-item>
          <el-descriptions-item label="空闲">{{ status?.database?.idle || 0 }}</el-descriptions-item>
          <el-descriptions-item label="等待次数">{{ status?.database?.wait_count || 0 }}</el-descriptions-item>
          <el-descriptions-item label="等待时长">{{ status?.database?.wait_duration || '0s' }}</el-descriptions-item>
          <el-descriptions-item label="最大连接数">{{ status?.database?.max_open_conns || 0 }}</el-descriptions-item>
          <el-descriptions-item label="最大空闲">{{ status?.database?.max_idle_conns || 0 }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- Redis 状态 -->
      <el-card class="monitor-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>Redis 状态</span>
            <el-tag :type="redisStatusType" size="small">{{ redisStatusText }}</el-tag>
          </div>
        </template>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="连接状态">{{ redisStatusText }}</el-descriptions-item>
          <el-descriptions-item label="内存使用">{{ formatMemory(status?.redis?.used_memory_mb) }}</el-descriptions-item>
          <el-descriptions-item label="客户端数">{{ status?.redis?.connected_clients || 0 }}</el-descriptions-item>
          <el-descriptions-item label="运行时间">{{ formatUptime(status?.redis?.uptime_seconds) }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 磁盘使用 -->
      <el-card class="monitor-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>磁盘使用</span>
            <el-tag type="info" size="small">uploads 目录</el-tag>
          </div>
        </template>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="目录大小">{{ formatMemory(status?.disk?.uploads_size_mb) }}</el-descriptions-item>
          <el-descriptions-item label="文件数量">{{ status?.disk?.uploads_count || 0 }} 个</el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { getSystemStatus } from '@/api'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'

const loading = ref(false)
const status = ref(null)
const refreshInterval = ref(0)
let timer = null

// 内存使用百分比
const memoryPercentage = computed(() => {
  if (!status.value?.memory) return 0
  const { alloc_mb, sys_mb } = status.value.memory
  if (!sys_mb || sys_mb === 0) return 0
  return Math.round((alloc_mb / sys_mb) * 100)
})

const memoryColor = computed(() => {
  const pct = memoryPercentage.value
  if (pct < 60) return '#67c23a'
  if (pct < 80) return '#e6a23c'
  return '#f56c6c'
})

// 数据库状态
const dbStatusType = computed(() => {
  const s = status.value?.database?.status
  if (s === 'connected') return 'success'
  if (s === 'error') return 'danger'
  return 'info'
})

const dbStatusText = computed(() => {
  const s = status.value?.database?.status
  if (s === 'connected') return '已连接'
  if (s === 'error') return '连接错误'
  if (s === 'disconnected') return '未连接'
  return '未知'
})

// Redis 状态
const redisStatusType = computed(() => {
  const s = status.value?.redis?.status
  if (s === 'connected') return 'success'
  if (s === 'error') return 'danger'
  return 'info'
})

const redisStatusText = computed(() => {
  const s = status.value?.redis?.status
  if (s === 'connected') return '已连接'
  if (s === 'error') return '连接错误'
  if (s === 'disconnected') return '未连接'
  return '未知'
})

// 格式化内存
function formatMemory(mb) {
  if (!mb || mb === 0) return '0 MB'
  if (mb < 1) return `${(mb * 1024).toFixed(1)} KB`
  if (mb < 1024) return `${mb.toFixed(1)} MB`
  return `${(mb / 1024).toFixed(2)} GB`
}

// 格式化运行时间
function formatUptime(seconds) {
  if (!seconds || seconds === 0) return '-'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)

  if (days > 0) return `${days}天${hours}小时${minutes}分钟`
  if (hours > 0) return `${hours}小时${minutes}分钟`
  return `${minutes}分钟`
}

// 获取系统状态
async function fetchStatus() {
  loading.value = true
  try {
    const res = await getSystemStatus()
    status.value = res.data
  } catch {
    ElMessage.error('获取系统状态失败')
  } finally {
    loading.value = false
  }
}

// 切换刷新间隔
function handleIntervalChange(interval) {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
  if (interval > 0) {
    timer = setInterval(fetchStatus, interval * 1000)
  }
}

onMounted(() => {
  fetchStatus()
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.system-monitor {
  padding: 20px;
}

.monitor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.monitor-header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.refresh-controls {
  display: flex;
  align-items: center;
}

.monitor-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 20px;
}

.monitor-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.memory-bar {
  margin-top: 16px;
}

.bar-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 8px;
}

@media (max-width: 768px) {
  .monitor-content {
    grid-template-columns: 1fr;
  }
}
</style>
