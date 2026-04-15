<template>
  <div v-loading="loading" class="dashboard">
    <!-- 统计数据卡片 -->
    <el-row :gutter="20" style="margin-bottom: 24px">
      <el-col :span="4" v-for="item in statCards" :key="item.key">
        <el-card shadow="hover" :class="['stat-card', { 'stat-card--alert': item.key === 'pending_reviews' && stats.pending_reviews > 0 }]">
          <div class="stat-card__icon" :style="{ backgroundColor: item.color + '15', color: item.color }">
            <el-icon :size="24"><component :is="item.icon" /></el-icon>
          </div>
          <el-statistic :value="statCardsValue(item)" :precision="item.precision || 0" :prefix="item.prefix || ''" />
          <div class="stat-card__label">{{ item.label }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 趋势折线图 -->
    <el-row :gutter="20" style="margin-bottom: 24px">
      <el-col :span="16">
        <el-card shadow="hover">
          <template #header>
            <div class="chart-header">
              <span>近 7 天趋势</span>
              <el-radio-group v-model="chartDays" size="small" @change="fetchTrend">
                <el-radio-button :value="7">7天</el-radio-button>
                <el-radio-button :value="14">14天</el-radio-button>
                <el-radio-button :value="30">30天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="lineChartRef" style="height: 320px" v-loading="trendLoading" />
        </el-card>
      </el-col>

      <!-- 收入概览 -->
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header><span>收入概览</span></template>
          <div class="revenue-summary">
            <div class="revenue-row">
              <span class="label">今日收入</span>
              <span class="value">¥{{ fmtMoney(todayRevenue) }}</span>
            </div>
            <div class="revenue-row">
              <span class="label">7日收入</span>
              <span class="value">¥{{ fmtMoney(weekRevenue) }}</span>
            </div>
            <div class="revenue-row">
              <span class="label">总收入</span>
              <span class="value accent">¥{{ fmtMoney(stats.total_revenue / 100) }}</span>
            </div>
            <el-divider />
            <div class="revenue-row">
              <span class="label">已完成订单</span>
              <span class="value">{{ completedOrders }}</span>
            </div>
            <div class="revenue-row">
              <span class="label">平均客单价</span>
              <span class="value">¥{{ fmtMoney(avgOrderAmount) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 近期活动 -->
    <el-row>
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header><span>最近订单</span></template>
          <el-table :data="recentOrders" style="width: 100%" stripe v-loading="ordersLoading">
            <el-table-column prop="id" label="订单号" width="80" align="center" />
            <el-table-column label="用户" width="120">
              <template #default="{ row }">{{ row.user?.nickname || '-' }}</template>
            </el-table-column>
            <el-table-column label="套图" min-width="160">
              <template #default="{ row }">{{ row.photoset?.title || row.order_type || '-' }}</template>
            </el-table-column>
            <el-table-column label="金额" width="100" align="right">
              <template #default="{ row }">¥{{ (row.amount / 100).toFixed(2) }}</template>
            </el-table-column>
            <el-table-column label="状态" width="90" align="center">
              <template #default="{ row }">
                <el-tag :type="orderStatusTag(row.status)" size="small">{{ orderStatusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="时间" width="165" align="center">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, nextTick, computed } from 'vue'
import { getStats, getStatsTrend, getOrderList } from '@/api'
import { User, PictureFilled, ShoppingCart, Money, Warning } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

const loading = ref(false)
const trendLoading = ref(false)
const ordersLoading = ref(false)
const chartDays = ref(7)
const lineChartRef = ref(null)
let chart = null

const stats = reactive({
  total_users: 0, total_photosets: 0, total_orders: 0, total_revenue: 0, pending_reviews: 0
})

const trendData = ref([])

const statCards = [
  { key: 'total_users', label: '总用户数', icon: User, color: '#409EFF' },
  { key: 'total_photosets', label: '总套图数', icon: PictureFilled, color: '#67C23A' },
  { key: 'total_orders', label: '总订单数', icon: ShoppingCart, color: '#E6A23C' },
  { key: 'total_revenue', label: '总收入', icon: Money, color: '#F56C6C', prefix: '¥', precision: 2 },
  { key: 'pending_reviews', label: '待审核', icon: Warning, color: '#F56C6C' }
]

function statCardsValue(item) {
  if (item.key === 'total_revenue') return (stats[item.key] || 0) / 100
  return stats[item.key] || 0
}

const recentOrders = ref([])

// 计算收入
const todayRevenue = computed(() => {
  if (!trendData.value.length) return 0
  return trendData.value[trendData.value.length - 1].revenue || 0
})
const weekRevenue = computed(() => {
  return trendData.value.slice(-7).reduce((s, d) => s + d.revenue, 0)
})
const completedOrders = computed(() => recentOrders.value.filter(o => o.status === 'completed').length)
const avgOrderAmount = computed(() => {
  const done = recentOrders.value.filter(o => o.status === 'completed')
  if (!done.length) return 0
  const sum = done.reduce((s, o) => s + (o.amount || 0), 0)
  return sum / done.length / 100
})

function fmtMoney(v) { return Number(v || 0).toFixed(2) }

function orderStatusLabel(s) {
  const m = { pending: '待支付', paid: '已支付', completed: '已完成', refunded: '已退款', cancelled: '已取消' }
  return m[s] || s
}
function orderStatusTag(s) {
  const m = { pending: 'warning', paid: 'primary', completed: 'success', refunded: 'danger', cancelled: 'info' }
  return m[s] || 'info'
}
function formatTime(t) {
  if (!t) return ''
  const ts = Number(t)
  if (ts < 1e12) return new Date(ts * 1000).toLocaleString('zh-CN')
  return new Date(ts).toLocaleString('zh-CN')
}

function renderChart() {
  if (!lineChartRef.value) return
  if (!chart) {
    chart = echarts.init(lineChartRef.value)
    window.addEventListener('resize', () => chart?.resize())
  }
  const dates = trendData.value.map(d => d.date)
  const option = {
    tooltip: { trigger: 'axis' },
    legend: { data: ['新用户', '新订单', '收入(元)'], bottom: 0 },
    grid: { top: 10, right: 60, bottom: 40, left: 50 },
    xAxis: { type: 'category', data: dates, boundaryGap: false },
    yAxis: [
      { type: 'value', name: '数量' },
      { type: 'value', name: '收入', position: 'right' }
    ],
    series: [
      { name: '新用户', type: 'line', smooth: true, data: trendData.value.map(d => d.new_users), itemStyle: { color: '#409EFF' }, areaStyle: { opacity: 0.08 } },
      { name: '新订单', type: 'line', smooth: true, data: trendData.value.map(d => d.new_orders), itemStyle: { color: '#67C23A' }, areaStyle: { opacity: 0.08 } },
      { name: '收入(元)', type: 'line', smooth: true, yAxisIndex: 1, data: trendData.value.map(d => d.revenue), itemStyle: { color: '#F56C6C' }, areaStyle: { opacity: 0.08 } }
    ]
  }
  chart.setOption(option, true)
}

async function fetchStats() {
  loading.value = true
  try {
    const res = await getStats()
    Object.assign(stats, res.data)
  } catch {}
  finally { loading.value = false }
}

async function fetchTrend() {
  trendLoading.value = true
  try {
    const res = await getStatsTrend(chartDays.value)
    trendData.value = res.data?.trend || []
    await nextTick()
    renderChart()
  } catch {}
  finally { trendLoading.value = false }
}

async function fetchRecentOrders() {
  ordersLoading.value = true
  try {
    const res = await getOrderList({ page: 1, size: 10 })
    recentOrders.value = res.data?.data || []
  } catch {}
  finally { ordersLoading.value = false }
}

onMounted(() => {
  fetchStats()
  fetchTrend()
  fetchRecentOrders()
})

onBeforeUnmount(() => {
  chart?.dispose()
  chart = null
})
</script>

<style scoped>
.dashboard { padding: 4px 0; }
.stat-card {
  text-align: center; border-radius: 12px; transition: all 0.3s;
  min-height: 120px; display: flex; flex-direction: column; justify-content: center;
}
.stat-card:hover { transform: translateY(-3px); box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
.stat-card--alert { border: 2px solid #F56C6C; }
.stat-card__icon {
  display: inline-flex; align-items: center; justify-content: center;
  width: 44px; height: 44px; border-radius: 10px; margin: 0 auto 10px;
}
.stat-card__label { margin-top: 6px; font-size: 13px; color: #909399; }
.chart-header { display: flex; justify-content: space-between; align-items: center; }
.revenue-summary { padding: 4px 0; }
.revenue-row { display: flex; justify-content: space-between; padding: 8px 0; font-size: 14px; }
.revenue-row .label { color: #606266; }
.revenue-row .value { font-weight: 600; color: #303133; }
.revenue-row .value.accent { color: #F56C6C; font-size: 16px; }
</style>
