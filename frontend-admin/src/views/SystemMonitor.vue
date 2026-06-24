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
        <el-button type="danger" @click="handleRestart" :loading="restarting" style="margin-left: 12px">
          <el-icon><RefreshRight /></el-icon>
          重启服务
        </el-button>
      </div>
    </div>

    <div v-loading="loading" class="monitor-content">
      <!-- 健康检查 -->
      <el-card class="monitor-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>健康检查</span>
            <el-tag :type="healthStatus === 'healthy' ? 'success' : 'danger'" size="small">
              {{ healthStatus === 'healthy' ? '健康' : '异常' }}
            </el-tag>
          </div>
        </template>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="API 状态">
            <el-tag :type="healthData?.status === 'ok' ? 'success' : 'danger'" size="small">
              {{ healthData?.status === 'ok' ? '正常' : '异常' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="响应时间">{{ healthLatency }}ms</el-descriptions-item>
          <el-descriptions-item label="数据库">
            <el-tag :type="healthData?.database === 'ok' ? 'success' : 'danger'" size="small">
              {{ healthData?.database === 'ok' ? '正常' : '异常' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="Redis">
            <el-tag :type="healthData?.redis === 'ok' ? 'success' : 'danger'" size="small">
              {{ healthData?.redis === 'ok' ? '正常' : '异常' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="最后检查">{{ lastHealthCheck || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 服务器信息 -->
      <el-card class="monitor-card" shadow="hover">
        <template #header><div class="card-header"><span>服务器信息</span></div></template>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="运行时间">{{ status?.server?.uptime || '-' }}</el-descriptions-item>
          <el-descriptions-item label="Go 版本">{{ status?.server?.go_version || '-' }}</el-descriptions-item>
          <el-descriptions-item label="系统">{{ status?.server?.os || '-' }}</el-descriptions-item>
          <el-descriptions-item label="主机名">{{ status?.server?.hostname || '-' }}</el-descriptions-item>
          <el-descriptions-item label="协程数">
            <el-tag :type="(status?.server?.goroutines || 0) > 500 ? 'warning' : 'success'" size="small">
              {{ status?.server?.goroutines || 0 }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 内存仪表盘 -->
      <el-card class="monitor-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>内存使用</span>
            <el-tag type="info" size="small">{{ status?.memory?.gc_count || 0 }} GC</el-tag>
          </div>
        </template>
        <div class="chart-row">
          <div ref="memoryGaugeRef" style="width: 180px; height: 180px" />
          <div class="memory-stats">
            <div class="mem-item">
              <span class="mem-label">当前分配</span>
              <span class="mem-value">{{ status?.memory?.alloc_mb || 0 }} MB</span>
            </div>
            <div class="mem-item">
              <span class="mem-label">总分配量</span>
              <span class="mem-value">{{ status?.memory?.total_alloc_mb || 0 }} MB</span>
            </div>
            <div class="mem-item">
              <span class="mem-label">系统内存</span>
              <span class="mem-value">{{ status?.memory?.sys_mb || 0 }} MB</span>
            </div>
            <div class="mem-item">
              <span class="mem-label">GC 次数</span>
              <span class="mem-value">{{ status?.memory?.gc_count || 0 }}</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 数据库连接池 -->
      <el-card class="monitor-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>数据库连接池</span>
            <el-tag :type="dbStatusType" size="small">{{ dbStatusText }}</el-tag>
          </div>
        </template>
        <div class="chart-row">
          <div ref="dbPoolGaugeRef" style="width: 220px; height: 180px" />
          <div class="memory-stats">
            <div class="mem-item">
              <span class="mem-label">打开连接</span>
              <span class="mem-value">{{ status?.database?.open_connections || 0 }}</span>
            </div>
            <div class="mem-item">
              <span class="mem-label">使用中 / 空闲</span>
              <span class="mem-value">{{ status?.database?.in_use || 0 }} / {{ status?.database?.idle || 0 }}</span>
            </div>
            <div class="mem-item">
              <span class="mem-label">最大连接</span>
              <span class="mem-value">{{ status?.database?.max_open_conns || 0 }}</span>
            </div>
            <div class="mem-item">
              <span class="mem-label">等待次数</span>
              <span class="mem-value" :class="{ 'mem-warn': (status?.database?.wait_count || 0) > 0 }">
                {{ status?.database?.wait_count || 0 }}
              </span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- Redis 状态 -->
      <el-card class="monitor-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>Redis 状态</span>
            <el-tag :type="redisStatusType" size="small">{{ redisStatusText }}</el-tag>
          </div>
        </template>
        <div class="chart-row">
          <div ref="redisMemoryGaugeRef" style="width: 180px; height: 180px" />
          <div class="memory-stats">
            <div class="mem-item">
              <span class="mem-label">内存使用</span>
              <span class="mem-value">{{ formatMemory(status?.redis?.used_memory_mb) }}</span>
            </div>
            <div class="mem-item">
              <span class="mem-label">客户端数</span>
              <span class="mem-value">{{ status?.redis?.connected_clients || 0 }}</span>
            </div>
            <div class="mem-item">
              <span class="mem-label">运行时间</span>
              <span class="mem-value">{{ formatUptime(status?.redis?.uptime_seconds) }}</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 磁盘使用 -->
      <el-card class="monitor-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>磁盘使用</span>
            <el-tag type="info" size="small">uploads 目录</el-tag>
          </div>
        </template>
        <div class="disk-info">
          <div class="disk-stat">
            <span class="disk-label">目录大小</span>
            <span class="disk-big">{{ formatMemory(status?.disk?.uploads_size_mb) }}</span>
          </div>
          <div class="disk-stat">
            <span class="disk-label">文件数量</span>
            <span class="disk-big">{{ status?.disk?.uploads_count || 0 }} 个</span>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { getSystemStatus, restartServer, healthCheck } from '@/api'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import { Refresh, RefreshRight } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

const loading = ref(false)
const restarting = ref(false)
const status = ref(null)
const refreshInterval = ref(0)
let timer = null

const memoryGaugeRef = ref(null)
const dbPoolGaugeRef = ref(null)
const redisMemoryGaugeRef = ref(null)
let memoryGauge = null
let dbPoolGauge = null
let redisMemoryGauge = null

const healthStatus = ref('checking')
const healthData = ref(null)
const healthLatency = ref(0)
const lastHealthCheck = ref('')
let healthTimer = null

const dbStatusType = computed(() => {
  const s = status.value?.database?.status
  return s === 'connected' ? 'success' : s === 'error' ? 'danger' : 'info'
})
const dbStatusText = computed(() => {
  const m = { connected: '已连接', error: '连接错误', disconnected: '未连接' }
  return m[status.value?.database?.status] || '未知'
})

const redisStatusType = computed(() => {
  const s = status.value?.redis?.status
  return s === 'connected' ? 'success' : s === 'error' ? 'danger' : 'info'
})
const redisStatusText = computed(() => {
  const m = { connected: '已连接', error: '连接错误', disconnected: '未连接' }
  return m[status.value?.redis?.status] || '未知'
})

function formatMemory(mb) {
  if (!mb || mb === 0) return '0 MB'
  if (mb < 1) return `${(mb * 1024).toFixed(1)} KB`
  if (mb < 1024) return `${mb.toFixed(1)} MB`
  return `${(mb / 1024).toFixed(2)} GB`
}
function formatUptime(seconds) {
  if (!seconds || seconds === 0) return '-'
  const d = Math.floor(seconds / 86400), h = Math.floor((seconds % 86400) / 3600), m = Math.floor((seconds % 3600) / 60)
  if (d > 0) return `${d}天${h}小时${m}分钟`
  if (h > 0) return `${h}小时${m}分钟`
  return `${m}分钟`
}

function renderGaugeChart(refEl, chartRef, title, value, max, unit) {
  if (!refEl.value) return
  if (!chartRef.value) chartRef.value = echarts.init(refEl.value)
  const pct = max > 0 ? Math.min((value / max) * 100, 100) : 0
  let color = '#67C23A'
  if (pct > 80) color = '#F56C6C'
  else if (pct > 60) color = '#E6A23C'
  chartRef.value.setOption({
    series: [{
      type: 'gauge',
      startAngle: 210, endAngle: -30,
      center: ['50%', '55%'],
      radius: '85%',
      min: 0, max,
      progress: { show: true, width: 12, itemStyle: { color } },
      axisLine: { lineStyle: { width: 12 } },
      axisTick: { show: false },
      splitLine: { show: false },
      axisLabel: { show: false },
      detail: {
        valueAnimation: true,
        formatter: `{value} ${unit}`,
        fontSize: 16,
        offsetCenter: [0, '75%'],
        color: '#303133'
      },
      title: { offsetCenter: [0, '45%'], fontSize: 12, color: '#909399' },
      data: [{ value: Math.round(value), name: title }]
    }]
  }, true)
}

function renderCharts() {
  const mem = status.value?.memory
  if (mem) {
    renderGaugeChart(memoryGaugeRef, { value: memoryGauge }, '内存使用率', mem.alloc_mb, mem.sys_mb || 512, 'MB')
  }
  const db = status.value?.database
  if (db) {
    renderGaugeChart(dbPoolGaugeRef, { value: dbPoolGauge }, '连接池负载', db.open_connections || 0, db.max_open_conns || 50, '个')
  }
  const redis = status.value?.redis
  if (redis) {
    renderGaugeChart(redisMemoryGaugeRef, { value: redisMemoryGauge }, 'Redis内存', redis.used_memory_mb || 0, 1024, 'MB')
  }
}

async function fetchStatus() {
  loading.value = true
  try {
    const res = await getSystemStatus()
    status.value = res.data
    await nextTick()
    renderCharts()
  } catch {
    ElMessage.error('获取系统状态失败')
  } finally {
    loading.value = false
  }
}

async function fetchHealth() {
  const start = Date.now()
  try {
    const res = await healthCheck()
    healthLatency.value = Date.now() - start
    healthData.value = res.data
    healthStatus.value = 'healthy'
    lastHealthCheck.value = new Date().toLocaleTimeString()
  } catch {
    healthLatency.value = Date.now() - start
    healthStatus.value = 'unhealthy'
    healthData.value = null
    lastHealthCheck.value = new Date().toLocaleTimeString()
  }
}

async function handleRestart() {
  try {
    await ElMessageBox.confirm(
      '确定要重启后端服务吗？<br><br><strong>注意：</strong><br>- 重启期间服务将暂时不可用（约 20-30 秒）<br>- 所有用户的连接将被中断<br>- 页面将在后端恢复后自动刷新',
      '重启确认',
      { type: 'warning', dangerouslyUseHTMLString: true, confirmButtonText: '确定重启', cancelButtonText: '取消' }
    )
  } catch { return }

  restarting.value = true
  try {
    const res = await restartServer()
    const delay = res.data?.delay || 5
    ElMessage.info(`后端正在重启，预计 ${delay + 15} 秒后恢复...`)
    await new Promise(r => setTimeout(r, (delay + 8) * 1000))
    const loadingInstance = ElLoading.service({ lock: true, text: '正在等待后端恢复...', background: 'rgba(0, 0, 0, 0.7)' })
    let retries = 0
    await new Promise((resolve) => {
      const poll = setInterval(async () => {
        retries++
        try {
          await healthCheck()
          clearInterval(poll)
          loadingInstance.close()
          ElMessage.success('后端重启成功！页面即将刷新...')
          setTimeout(() => window.location.reload(), 1500)
          resolve()
        } catch {
          if (retries >= 30) {
            clearInterval(poll); loadingInstance.close()
            ElMessage.warning('重启超时，请手动检查服务状态')
            restarting.value = false
            resolve()
          }
        }
      }, 2000)
    })
  } catch (err) {
    ElMessage.error(err.response?.data?.message || '重启失败')
    restarting.value = false
  }
}

function handleIntervalChange(interval) {
  if (timer) { clearInterval(timer); timer = null }
  if (interval > 0) timer = setInterval(fetchStatus, interval * 1000)
}

onMounted(() => {
  fetchStatus()
  fetchHealth()
  healthTimer = setInterval(fetchHealth, 30000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
  if (healthTimer) clearInterval(healthTimer)
  memoryGauge?.dispose(); memoryGauge = null
  dbPoolGauge?.dispose(); dbPoolGauge = null
  redisMemoryGauge?.dispose(); redisMemoryGauge = null
})
</script>

<style scoped>
.system-monitor { padding: 20px; }
.monitor-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.monitor-header h2 { margin: 0; font-size: 20px; color: #303133; }
.refresh-controls { display: flex; align-items: center; }
.monitor-content { display: grid; grid-template-columns: repeat(auto-fit, minmax(420px, 1fr)); gap: 20px; }
.monitor-card { height: 100%; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.chart-row { display: flex; align-items: center; gap: 20px; }
.memory-stats { flex: 1; display: flex; flex-direction: column; gap: 10px; }
.mem-item { display: flex; justify-content: space-between; padding: 6px 0; border-bottom: 1px dashed #ebeef5; }
.mem-label { color: #909399; font-size: 13px; }
.mem-value { font-weight: 600; font-size: 14px; color: #303133; }
.mem-warn { color: #F56C6C !important; }
.disk-info { display: flex; gap: 40px; padding: 10px 0; }
.disk-stat { display: flex; flex-direction: column; align-items: center; gap: 8px; }
.disk-label { color: #909399; font-size: 13px; }
.disk-big { font-size: 28px; font-weight: 700; color: #409EFF; }
</style>
