<template>
  <div class="performance-monitor">
    <div class="perf-header">
      <h2>性能分析</h2>
      <div class="perf-controls">
        <el-button-group style="margin-right: 12px">
          <el-button size="default" @click="downloadCPU" :loading="cpuDownloading">
            <el-icon><Download /></el-icon>CPU Profile (30s)
          </el-button>
          <el-button size="default" @click="downloadHeap" :loading="heapDownloading">
            <el-icon><Download /></el-icon>Heap Profile
          </el-button>
          <el-button size="default" @click="exportMetrics">
            <el-icon><DocumentCopy /></el-icon>导出 JSON
          </el-button>
        </el-button-group>
        <el-select v-model="refreshInterval" @change="handleIntervalChange" style="width: 120px">
          <el-option label="手动刷新" :value="0" />
          <el-option label="5秒刷新" :value="5" />
          <el-option label="15秒刷新" :value="15" />
          <el-option label="30秒刷新" :value="30" />
        </el-select>
        <el-button type="primary" @click="fetchMetrics" :loading="loading" style="margin-left: 12px">
          <el-icon><Refresh /></el-icon>刷新
        </el-button>
      </div>
    </div>

    <div v-loading="loading" class="perf-content">
      <!-- 关键指标卡片 6 个 -->
      <el-row :gutter="16" style="margin-bottom: 20px">
        <el-col :span="4" v-for="card in keyCards" :key="card.key">
          <el-card shadow="hover" class="key-card">
            <div class="key-icon" :style="{ backgroundColor: card.color + '15', color: card.color }">
              <el-icon :size="22"><component :is="card.icon" /></el-icon>
            </div>
            <div class="key-info">
              <div class="key-value">{{ card.value }}</div>
              <div class="key-label">{{ card.label }}</div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 路由延迟 Top10 + 响应时间分布 -->
      <el-row :gutter="16" style="margin-bottom: 20px">
        <el-col :span="13">
          <el-card shadow="hover">
            <template #header><span>路由延迟 Top10（平均响应时间）</span></template>
            <div ref="routeLatencyChartRef" style="height: 340px" />
          </el-card>
        </el-col>
        <el-col :span="11">
          <el-card shadow="hover">
            <template #header><span>响应时间分布</span></template>
            <div ref="latencyDistChartRef" style="height: 340px" />
          </el-card>
        </el-col>
      </el-row>

      <!-- QPS/错误率趋势 + 状态码分布 -->
      <el-row :gutter="16" style="margin-bottom: 20px">
        <el-col :span="14">
          <el-card shadow="hover">
            <template #header><span>QPS & 错误率趋势</span></template>
            <div ref="qpsChartRef" style="height: 320px" />
          </el-card>
        </el-col>
        <el-col :span="10">
          <el-card shadow="hover">
            <template #header><span>状态码分布</span></template>
            <div ref="statusCodeChartRef" style="height: 320px" />
          </el-card>
        </el-col>
      </el-row>

      <!-- 内存详情（4 个仪表盘） -->
      <el-row :gutter="16" style="margin-bottom: 20px">
        <el-col :span="6">
          <el-card shadow="hover">
            <template #header><span>内存使用率</span></template>
            <div ref="memAllocGaugeRef" style="width: 100%; height: 200px" />
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <template #header><span>GC 暂停</span></template>
            <div ref="gcPauseGaugeRef" style="width: 100%; height: 200px" />
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <template #header><span>堆对象数</span></template>
            <div ref="heapObjectsGaugeRef" style="width: 100%; height: 200px" />
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover">
            <template #header><span>栈内存</span></template>
            <div ref="stackGaugeRef" style="width: 100%; height: 200px" />
          </el-card>
        </el-col>
      </el-row>

      <!-- 内存数值明细表格 -->
      <el-row style="margin-bottom: 20px">
        <el-col :span="24">
          <el-card shadow="hover">
            <template #header><span>内存明细</span></template>
            <el-descriptions :column="4" border size="small">
              <el-descriptions-item label="Alloc">{{ (mem.alloc_mb || 0).toFixed(1) }} MB</el-descriptions-item>
              <el-descriptions-item label="Sys">{{ (mem.sys_mb || 0).toFixed(1) }} MB</el-descriptions-item>
              <el-descriptions-item label="Total Alloc">{{ (mem.total_alloc_mb || 0).toFixed(1) }} MB</el-descriptions-item>
              <el-descriptions-item label="Next GC">{{ (mem.next_gc_mb || 0).toFixed(1) }} MB</el-descriptions-item>
              <el-descriptions-item label="Heap Objects">{{ (mem.heap_objects || 0).toLocaleString() }}</el-descriptions-item>
              <el-descriptions-item label="Stack Inuse">{{ (mem.stack_inuse_mb || 0).toFixed(2) }} MB</el-descriptions-item>
              <el-descriptions-item label="GC Pause Total">{{ (mem.gc_pause_total_ms || 0).toFixed(2) }} ms</el-descriptions-item>
              <el-descriptions-item label="Last GC Pause">{{ (mem.last_gc_pause_ms || 0).toFixed(3) }} ms</el-descriptions-item>
              <el-descriptions-item label="GC Cycles">{{ mem.gc_cycles || 0 }}</el-descriptions-item>
              <el-descriptions-item label="Forced GC">{{ mem.num_forced_gc || 0 }}</el-descriptions-item>
              <el-descriptions-item label="Last GC">{{ mem.last_gc_seconds_ago || 0 }}s ago</el-descriptions-item>
              <el-descriptions-item label="Goroutines">{{ metrics?.num_goroutine || 0 }}</el-descriptions-item>
            </el-descriptions>
          </el-card>
        </el-col>
      </el-row>

      <!-- Goroutine 分析面板 -->
      <el-card shadow="hover">
        <template #header>
          <div style="display: flex; justify-content: space-between; align-items: center">
            <span>Goroutine 分析</span>
            <el-button size="small" @click="fetchGoroutines" :loading="goroutineLoading">
              <el-icon><Search /></el-icon>获取堆栈
            </el-button>
          </div>
        </template>
        <el-alert
          v-if="!goroutineData"
          type="info" :closable="false" show-icon
          title="点击「获取堆栈」查看所有 Goroutine 的运行状态，用于排查协程泄漏"
        />
        <template v-else>
          <div style="margin-bottom: 12px; display: flex; gap: 12px; align-items: center">
            <el-tag type="primary">总数: {{ goroutineData.goroutine_count }}</el-tag>
            <el-input
              v-model="goroutineFilter" placeholder="搜索 Goroutine（函数名/文件路径）" clearable
              style="width: 320px" size="small"
            />
            <el-button size="small" @click="goroutineFilter = ''">清除</el-button>
          </div>
          <div style="max-height: 500px; overflow-y: auto; background: #f5f7fa; border-radius: 6px; padding: 12px">
            <pre style="margin: 0; font-size: 12px; line-height: 1.6; white-space: pre-wrap; word-break: break-all">{{ filteredGoroutines }}</pre>
          </div>
        </template>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { getSystemMetrics, getGoroutines, downloadCPUProfile, downloadHeapProfile } from '@/api'
import { ElMessage } from 'element-plus'
import { Refresh, Timer, Connection, WarningFilled, DataAnalysis, Download, DocumentCopy, Search, TrendCharts, Coin } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

const loading = ref(false)
const refreshInterval = ref(15)
const metrics = ref(null)
let timer = null

// Profile download states
const cpuDownloading = ref(false)
const heapDownloading = ref(false)

// Goroutine analysis
const goroutineLoading = ref(false)
const goroutineData = ref(null)
const goroutineFilter = ref('')

// QPS tracking: store last snapshot to calculate delta
const lastTotalRequests = ref(0)
const lastFetchTime = ref(0)

// Chart refs
const routeLatencyChartRef = ref(null)
const latencyDistChartRef = ref(null)
const statusCodeChartRef = ref(null)
const qpsChartRef = ref(null)
const memAllocGaugeRef = ref(null)
const gcPauseGaugeRef = ref(null)
const heapObjectsGaugeRef = ref(null)
const stackGaugeRef = ref(null)
let routeLatencyChart = null
let latencyDistChart = null
let statusCodeChart = null
let qpsChart = null
let memAllocGauge = null
let gcPauseGauge = null
let heapObjectsGauge = null
let stackGauge = null

// History for QPS/error-rate trend
const qpsHistory = ref([])
const errorRateHistory = ref([])
const MAX_HISTORY = 40

const mem = computed(() => metrics.value?.memory || {})

const currentQps = computed(() => {
  if (lastFetchTime.value === 0) return 0
  const total = metrics.value?.total_requests || 0
  const delta = total - lastTotalRequests.value
  const timeDelta = (Date.now() - lastFetchTime.value) / 1000
  return timeDelta > 0 ? +(delta / timeDelta).toFixed(1) : 0
})

const currentErrorRate = computed(() => {
  const total = lastTotalRequests.value > 0
    ? (metrics.value?.total_requests || 0) - lastTotalRequests.value
    : 0
  const errors = metrics.value?.total_errors_5xx || 0
  return total > 0 ? +((errors / total) * 100).toFixed(2) : 0
})

const keyCards = computed(() => {
  const m = metrics.value || {}
  return [
    { key: 'total', label: '总请求数', value: (m.total_requests || 0).toLocaleString(), icon: Connection, color: '#409EFF' },
    { key: 'active', label: '活跃请求', value: m.active_requests || 0, icon: Timer, color: '#67C23A' },
    { key: 'errors', label: '5xx 错误', value: m.total_errors_5xx || 0, icon: WarningFilled, color: '#F56C6C' },
    { key: 'goroutines', label: 'Goroutine', value: m.num_goroutine || 0, icon: DataAnalysis, color: '#E6A23C' },
    { key: 'qps', label: 'QPS', value: currentQps.value, icon: TrendCharts, color: '#722ED1' },
    { key: 'errorRate', label: '错误率', value: currentErrorRate.value + '%', icon: Coin, color: '#EB2F96' },
  ]
})

const filteredGoroutines = computed(() => {
  if (!goroutineData.value?.goroutines) return ''
  if (!goroutineFilter.value) return goroutineData.value.goroutines
  const filter = goroutineFilter.value.toLowerCase()
  const lines = goroutineData.value.goroutines.split('\n')
  const result = []
  let keepBlock = true
  for (const line of lines) {
    if (line.startsWith('goroutine ')) {
      keepBlock = line.toLowerCase().includes(filter)
    } else if (!keepBlock && line.trim() === '') {
      keepBlock = true
    }
    if (keepBlock || line.trim() === '') {
      result.push(line)
    }
  }
  return result.join('\n')
})

// ============ Chart Rendering ============

function renderRouteLatencyChart() {
  if (!routeLatencyChartRef.value) return
  if (!routeLatencyChart) routeLatencyChart = echarts.init(routeLatencyChartRef.value)
  const routes = (metrics.value?.routes || [])
    .filter(r => r.count > 0)
    .sort((a, b) => b.avg_time_ms - a.avg_time_ms)
    .slice(0, 10)

  routeLatencyChart.setOption({
    tooltip: {
      trigger: 'axis', axisPointer: { type: 'shadow' },
      formatter: (p) => {
        const d = p[0]
        return `${d.name}<br/>平均: ${d.value.toFixed(1)}ms<br/>请求数: ${routes[d.dataIndex].count}`
      }
    },
    grid: { top: 5, right: 30, bottom: 80, left: 180 },
    xAxis: { type: 'value', name: 'ms', axisLabel: { formatter: '{value}ms' } },
    yAxis: {
      type: 'category', data: routes.map(r => r.path).reverse(), inverse: true,
      axisLabel: { fontSize: 11, width: 170, overflow: 'truncate' }
    },
    series: [{
      type: 'bar', data: routes.map(r => r.avg_time_ms).reverse(),
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
          { offset: 0, color: '#67C23A' }, { offset: 0.5, color: '#E6A23C' }, { offset: 1, color: '#F56C6C' }
        ])
      },
      label: { show: true, position: 'right', formatter: p => p.value.toFixed(1) + 'ms', fontSize: 10 }
    }]
  }, true)
}

function renderLatencyDistChart() {
  if (!latencyDistChartRef.value) return
  if (!latencyDistChart) latencyDistChart = echarts.init(latencyDistChartRef.value)
  const routes = metrics.value?.routes || []
  const buckets = { '<10ms': 0, '10-50ms': 0, '50-100ms': 0, '100-500ms': 0, '>500ms': 0 }
  routes.forEach(r => {
    const t = r.avg_time_ms
    if (t < 10) buckets['<10ms']++
    else if (t < 50) buckets['10-50ms']++
    else if (t < 100) buckets['50-100ms']++
    else if (t < 500) buckets['100-500ms']++
    else buckets['>500ms']++
  })
  const data = Object.entries(buckets)
  const colors = ['#67C23A', '#409EFF', '#E6A23C', '#F56C6C', '#8B0000']
  latencyDistChart.setOption({
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' }, formatter: '{b}: {c} 条路由' },
    grid: { top: 10, right: 20, bottom: 40, left: 50 },
    xAxis: { type: 'category', data: data.map(d => d[0]), axisLabel: { fontSize: 11 } },
    yAxis: { type: 'value', name: '路由数' },
    series: [{
      type: 'bar', data: data.map((d, i) => ({ value: d[1], itemStyle: { color: colors[i] } })),
      label: { show: true, position: 'top', fontSize: 11 }
    }]
  }, true)
}

function renderStatusCodeChart() {
  if (!statusCodeChartRef.value) return
  if (!statusCodeChart) statusCodeChart = echarts.init(statusCodeChartRef.value)
  const routes = metrics.value?.routes || []
  const total = { '2xx': 0, '3xx': 0, '4xx': 0, '5xx': 0 }
  routes.forEach(r => {
    total['2xx'] += (r.status_2xx || 0)
    total['3xx'] += (r.status_3xx || 0)
    total['4xx'] += (r.status_4xx || 0)
    total['5xx'] += (r.status_5xx || 0)
  })
  const colors = { '2xx': '#67C23A', '3xx': '#409EFF', '4xx': '#E6A23C', '5xx': '#F56C6C' }
  const data = Object.entries(total).filter(([, v]) => v > 0).map(([k, v]) => ({ name: k, value: v, itemStyle: { color: colors[k] } }))
  statusCodeChart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    series: [{
      type: 'pie', radius: ['50%', '75%'], center: ['50%', '50%'],
      label: { fontSize: 13, formatter: '{b}\n{d}%' },
      emphasis: { label: { fontSize: 18, fontWeight: 'bold' } },
      data,
    }]
  }, true)
}

function renderQpsChart() {
  if (!qpsChartRef.value) return
  if (!qpsChart) qpsChart = echarts.init(qpsChartRef.value)
  const now = new Date().toLocaleTimeString()

  // Calculate real QPS (delta / timeDelta)
  const total = metrics.value?.total_requests || 0
  const errors = metrics.value?.total_errors_5xx || 0
  let qps = 0
  let errorRate = 0
  if (lastFetchTime.value > 0) {
    const delta = total - lastTotalRequests.value
    const timeDelta = (Date.now() - lastFetchTime.value) / 1000
    qps = timeDelta > 0 ? +(delta / timeDelta).toFixed(1) : 0
    errorRate = delta > 0 ? +((errors / delta) * 100).toFixed(2) : 0
  }

  qpsHistory.value.push({ time: now, value: qps })
  errorRateHistory.value.push({ time: now, value: errorRate })
  if (qpsHistory.value.length > MAX_HISTORY) qpsHistory.value.shift()
  if (errorRateHistory.value.length > MAX_HISTORY) errorRateHistory.value.shift()

  qpsChart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['QPS', '错误率'], bottom: 0 },
    grid: { top: 8, right: 60, bottom: 40, left: 55 },
    xAxis: { type: 'category', data: qpsHistory.value.map(h => h.time), axisLabel: { rotate: 45, fontSize: 10 } },
    yAxis: [
      { type: 'value', name: 'QPS', min: 0 },
      { type: 'value', name: '%', min: 0, max: 100 }
    ],
    series: [
      {
        name: 'QPS', type: 'line', smooth: true, symbol: 'circle', symbolSize: 4,
        data: qpsHistory.value.map(h => h.value),
        areaStyle: { color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(64,158,255,0.25)' }, { offset: 1, color: 'rgba(64,158,255,0.02)' }
        ]) },
        itemStyle: { color: '#409EFF' },
      },
      {
        name: '错误率', type: 'line', smooth: true, yAxisIndex: 1, symbol: 'diamond', symbolSize: 5,
        data: errorRateHistory.value.map(h => h.value),
        itemStyle: { color: '#F56C6C' },
        markLine: {
          silent: true,
          data: [{ yAxis: 5, lineStyle: { color: '#E6A23C', type: 'dashed' }, label: { formatter: '5% 告警线' } }]
        }
      }
    ]
  }, true)
}

// ============ Memory Gauges ============

function renderGauge(chartRef, chartVar, config) {
  if (!chartRef?.value) return
  if (!chartVar.value) chartVar.value = echarts.init(chartRef.value)
  chartVar.value.setOption({
    series: [{
      type: 'gauge', startAngle: 200, endAngle: -20, center: ['50%', '58%'], radius: '88%',
      progress: { show: true, width: 8, itemStyle: { color: config.color } },
      axisLine: { lineStyle: { width: 8 } },
      axisTick: { show: false }, splitLine: { show: false }, axisLabel: { show: false },
      detail: {
        valueAnimation: true, formatter: config.formatter, fontSize: 18,
        offsetCenter: [0, '65%'], color: '#303133'
      },
      title: { offsetCenter: [0, '40%'], fontSize: 11, color: '#909399' },
      data: [{ value: config.value, name: config.name }]
    }]
  }, true)
}

const memAllocGaugeVar = ref(null)
const gcPauseGaugeVar = ref(null)
const heapObjectsGaugeVar = ref(null)
const stackGaugeVar = ref(null)

function renderMemoryGauges() {
  const m = mem.value
  if (!m) return

  const allocMB = m.alloc_mb || 0
  const sysMB = m.sys_mb || 512
  const pct = sysMB > 0 ? Math.min((allocMB / sysMB) * 100, 100) : 0
  let allocColor = '#67C23A'
  if (pct > 80) allocColor = '#F56C6C'
  else if (pct > 60) allocColor = '#E6A23C'

  const gcPauseMs = m.last_gc_pause_ms || 0
  let gcColor = '#67C23A'
  if (gcPauseMs > 10) gcColor = '#F56C6C'
  else if (gcPauseMs > 5) gcColor = '#E6A23C'

  const heapObj = m.heap_objects || 0
  let heapColor = '#67C23A'
  if (heapObj > 500000) heapColor = '#F56C6C'
  else if (heapObj > 100000) heapColor = '#E6A23C'

  const stackMB = m.stack_inuse_mb || 0
  let stackColor = '#67C23A'
  if (stackMB > 10) stackColor = '#E6A23C'

  renderGauge(memAllocGaugeRef, memAllocGaugeVar, { value: Math.round(pct), name: 'Alloc/Sys', color: allocColor, formatter: '{value}%' })
  renderGauge(gcPauseGaugeRef, gcPauseGaugeVar, { value: +gcPauseMs.toFixed(2), name: 'Last Pause', color: gcColor, formatter: '{value}ms' })
  renderGauge(heapObjectsGaugeRef, heapObjectsGaugeVar, { value: heapObj, name: 'Heap Objects', color: heapColor, formatter: v => v >= 1000 ? (v / 1000).toFixed(1) + 'k' : v })
  renderGauge(stackGaugeRef, stackGaugeVar, { value: +stackMB.toFixed(2), name: 'Stack', color: stackColor, formatter: '{value}MB' })
}

// ============ Data Fetching ============

async function fetchMetrics() {
  loading.value = true
  try {
    const res = await getSystemMetrics()
    const now = Date.now()
    lastTotalRequests.value = metrics.value?.total_requests || 0
    lastFetchTime.value = now
    metrics.value = res.data
    await nextTick()
    renderRouteLatencyChart()
    renderLatencyDistChart()
    renderStatusCodeChart()
    renderQpsChart()
    renderMemoryGauges()
  } catch {
    ElMessage.error('获取性能数据失败')
  } finally {
    loading.value = false
  }
}

async function downloadCPU() {
  cpuDownloading.value = true
  try {
    const res = await downloadCPUProfile(30)
    const url = URL.createObjectURL(new Blob([res.data]))
    const link = document.createElement('a')
    link.href = url
    link.download = `cpu_profile_${Date.now()}.prof`
    link.click()
    URL.revokeObjectURL(url)
    ElMessage.success('CPU Profile 下载完成')
  } catch {
    ElMessage.error('CPU Profile 下载失败')
  } finally {
    cpuDownloading.value = false
  }
}

async function downloadHeap() {
  heapDownloading.value = true
  try {
    const res = await downloadHeapProfile()
    const url = URL.createObjectURL(new Blob([res.data]))
    const link = document.createElement('a')
    link.href = url
    link.download = `heap_profile_${Date.now()}.prof`
    link.click()
    URL.revokeObjectURL(url)
    ElMessage.success('Heap Profile 下载完成')
  } catch {
    ElMessage.error('Heap Profile 下载失败')
  } finally {
    heapDownloading.value = false
  }
}

function exportMetrics() {
  if (!metrics.value) {
    ElMessage.warning('请先获取性能数据')
    return
  }
  const json = JSON.stringify(metrics.value, null, 2)
  const url = URL.createObjectURL(new Blob([json], { type: 'application/json' }))
  const link = document.createElement('a')
  link.href = url
  link.download = `metrics_${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.json`
  link.click()
  URL.revokeObjectURL(url)
  ElMessage.success('导出成功')
}

async function fetchGoroutines() {
  goroutineLoading.value = true
  try {
    const res = await getGoroutines()
    goroutineData.value = res.data
    goroutineFilter.value = ''
    ElMessage.success('Goroutine 堆栈获取完成')
  } catch {
    ElMessage.error('获取 Goroutine 堆栈失败')
  } finally {
    goroutineLoading.value = false
  }
}

function handleIntervalChange(interval) {
  if (timer) { clearInterval(timer); timer = null }
  if (interval > 0) timer = setInterval(fetchMetrics, interval * 1000)
}

onMounted(() => {
  fetchMetrics()
  if (refreshInterval.value > 0) {
    timer = setInterval(fetchMetrics, refreshInterval.value * 1000)
  }
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
  routeLatencyChart?.dispose(); routeLatencyChart = null
  latencyDistChart?.dispose(); latencyDistChart = null
  statusCodeChart?.dispose(); statusCodeChart = null
  qpsChart?.dispose(); qpsChart = null
  memAllocGaugeVar.value?.dispose(); memAllocGaugeVar.value = null
  gcPauseGaugeVar.value?.dispose(); gcPauseGaugeVar.value = null
  heapObjectsGaugeVar.value?.dispose(); heapObjectsGaugeVar.value = null
  stackGaugeVar.value?.dispose(); stackGaugeVar.value = null
})
</script>

<style scoped>
.performance-monitor { padding: 20px; }
.perf-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
.perf-header h2 { margin: 0; font-size: 20px; color: #303133; }
.perf-controls { display: flex; align-items: center; }
.key-card {
  display: flex; align-items: center; gap: 16px; border-radius: 12px; padding: 8px 0;
  transition: all 0.3s;
}
.key-card:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
.key-icon {
  width: 48px; height: 48px; border-radius: 12px;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.key-info { flex: 1; min-width: 0; }
.key-value { font-size: 24px; font-weight: 700; color: #303133; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.key-label { font-size: 13px; color: #909399; margin-top: 2px; }
:deep(.el-card__header) { padding: 12px 16px; border-bottom: 1px solid #EBEEF5; }
:deep(.el-card__header span) { font-weight: 600; font-size: 14px; color: #303133; }
</style>
