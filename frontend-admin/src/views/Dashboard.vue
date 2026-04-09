<template>
  <div v-loading="loading || chartLoading" class="dashboard">
    <!-- 统计数据卡片 -->
    <el-row :gutter="20" style="margin-bottom: 24px">
      <el-col :span="3" v-for="item in statCards" :key="item.key">
        <el-card
          shadow="hover"
          :class="['stat-card', { 'stat-card--alert': item.key === 'pending_reviews' && stats.pending_reviews > 0 }]"
        >
          <div class="stat-card__icon" :style="{ backgroundColor: item.color + '15', color: item.color }">
            <el-icon :size="24"><component :is="item.icon" /></el-icon>
          </div>
          <el-statistic
            :value="formatStatValue(stats[item.key], item.format)"
            :precision="item.precision || 0"
            :prefix="item.prefix || ''"
          />
          <div class="stat-card__label">{{ item.label }}</div>
          <div class="stat-card__trend" v-if="item.trend && stats[item.key + '_trend'] !== undefined">
            <el-text :type="stats[item.key + '_trend'] >= 0 ? 'success' : 'danger'" size="small">
              <el-icon v-if="stats[item.key + '_trend'] >= 0" size="12"><CaretTop /></el-icon>
              <el-icon v-else size="12"><CaretBottom /></el-icon>
              {{ Math.abs(stats[item.key + '_trend']) }}%
            </el-text>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 详细统计数据 -->
    <el-row :gutter="20" style="margin-bottom: 24px">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span>收入分析</span>
          </template>
          <div class="detailed-stats">
            <div class="stat-row">
              <div class="stat-label">今日收入</div>
              <div class="stat-value">¥{{ formatCurrency(todayRevenue) }}</div>
              <div class="stat-trend" :class="trendClass(todayRevenueTrend)">
                <el-icon :size="12"><CaretTop v-if="todayRevenueTrend >= 0" /><CaretBottom v-else /></el-icon>
                <span>{{ Math.abs(todayRevenueTrend) }}%</span>
              </div>
            </div>
            <div class="stat-row">
              <div class="stat-label">本周收入</div>
              <div class="stat-value">¥{{ formatCurrency(weekRevenue) }}</div>
            </div>
            <div class="stat-row">
              <div class="stat-label">本月收入</div>
              <div class="stat-value">¥{{ formatCurrency(monthRevenue) }}</div>
            </div>
            <el-divider />
            <div class="stat-row">
              <div class="stat-label">套图购买收入</div>
              <div class="stat-value">¥{{ formatCurrency(photosetRevenue) }}</div>
              <div class="stat-ratio">{{ photosetPercentage }}%</div>
            </div>
            <div class="stat-row">
              <div class="stat-label">会员订阅收入</div>
              <div class="stat-value">¥{{ formatCurrency(membershipRevenue) }}</div>
              <div class="stat-ratio">{{ membershipPercentage }}%</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span>订单类型分布</span>
          </template>
          <div class="order-type-stats">
            <div class="order-type-item">
              <div class="order-type-info">
                <el-icon><PictureFilled /></el-icon>
                <span>套图购买</span>
              </div>
              <div class="order-type-value">
                <span class="count">{{ orderTypeStats.photoset }}</span>
                <span class="percentage">单 / {{ photosetOrderPercentage }}%</span>
              </div>
            </div>
            <div class="order-type-item">
              <div class="order-type-info">
                <el-icon><User /></el-icon>
                <span>会员订阅</span>
              </div>
              <div class="order-type-value">
                <span class="count">{{ orderTypeStats.membership }}</span>
                <span class="percentage">单 / {{ membershipOrderPercentage }}%</span>
              </div>
            </div>
            <el-divider />
            <div class="order-type-summary">
              <div class="summary-item">
                <span class="label">已完成订单</span>
                <span class="value">{{ completedOrders }}</span>
              </div>
              <div class="summary-item">
                <span class="label">退款订单</span>
                <span class="value">{{ refundedOrders }}</span>
              </div>
              <div class="summary-item">
                <span class="label">平均客单价</span>
                <span class="value">¥{{ formatCurrency(averageOrderAmount) }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 近期活动 -->
    <el-row>
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <span>近期活动</span>
          </template>
          <el-table :data="recentActivities" style="width: 100%" stripe>
            <el-table-column prop="time" label="时间" width="160">
              <template #default="{ row }">
                {{ row.time ? row.time : '--' }}
              </template>
            </el-table-column>
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag :type="activityTypeTag(row.type)" size="small">
                  {{ activityTypeLabel(row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="content" label="内容" min-width="200" />
            <el-table-column prop="user" label="用户" width="120">
              <template #default="{ row }">
                <el-text v-if="row.user">{{ row.user.nickname }}</el-text>
                <el-text v-else type="info">系统</el-text>
              </template>
            </el-table-column>
            <el-table-column prop="amount" label="金额" width="100" align="right">
              <template #default="{ row }">
                <template v-if="row.amount">
                  ¥{{ (row.amount / 100).toFixed(2) }}
                </template>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { getStats } from '@/api'
import { 
  User, 
  PictureFilled, 
  ShoppingCart, 
  Money, 
  Warning,
  CaretTop,
  CaretBottom
} from '@element-plus/icons-vue'

const loading = ref(false)

const stats = reactive({
  total_users: 0,
  total_photosets: 0,
  total_orders: 0,
  total_revenue: 0,
  pending_reviews: 0,
  // 趋势数据
  total_users_trend: 0,
  total_photosets_trend: 0,
  total_orders_trend: 0,
  total_revenue_trend: 0
})

const statCards = [
  { key: 'total_users', label: '总用户数', icon: User, color: '#409EFF', trend: true },
  { key: 'total_photosets', label: '总套图数', icon: PictureFilled, color: '#67C23A', trend: true },
  { key: 'total_orders', label: '总订单数', icon: ShoppingCart, color: '#E6A23C', trend: true },
  { key: 'total_revenue', label: '总收入', icon: Money, color: '#F56C6C', prefix: '¥', precision: 2, format: 'currency', trend: true },
  { key: 'pending_reviews', label: '待审核', icon: Warning, color: '#F56C6C' }
]

// 收入分析数据
const todayRevenue = ref(125000)  // 1250元
const todayRevenueTrend = ref(15)
const weekRevenue = ref(850000)   // 8500元
const monthRevenue = ref(3200000) // 32000元
const photosetRevenue = ref(1500000) // 15000元
const membershipRevenue = ref(1700000) // 17000元

// 订单类型统计
const orderTypeStats = reactive({
  photoset: 85,
  membership: 42
})

const completedOrders = ref(120)
const refundedOrders = ref(8)
const averageOrderAmount = ref(2450) // 24.5元

const recentActivities = ref([
  { time: '2026-04-09 10:30:00', type: 'order', content: '用户购买套图"夏日海边"', user: { nickname: '张三' }, amount: 1990 },
  { time: '2026-04-09 09:15:00', type: 'photoset', content: '创作者发布新套图"城市夜景"', user: { nickname: '李四' } },
  { time: '2026-04-08 16:45:00', type: 'refund', content: '订单#10086申请退款', user: { nickname: '王五' }, amount: 1990 },
  { time: '2026-04-08 14:20:00', type: 'order', content: '用户订阅年度会员', user: { nickname: '赵六' }, amount: 29900 },
  { time: '2026-04-08 11:10:00', type: 'register', content: '新用户注册', user: { nickname: '钱七' } }
])

// 计算属性
const photosetPercentage = computed(() => {
  const total = photosetRevenue.value + membershipRevenue.value
  return total > 0 ? Math.round((photosetRevenue.value / total) * 100) : 0
})

const membershipPercentage = computed(() => {
  const total = photosetRevenue.value + membershipRevenue.value
  return total > 0 ? Math.round((membershipRevenue.value / total) * 100) : 0
})

const totalOrders = computed(() => {
  return orderTypeStats.photoset + orderTypeStats.membership
})

const photosetOrderPercentage = computed(() => {
  return totalOrders.value > 0 ? Math.round((orderTypeStats.photoset / totalOrders.value) * 100) : 0
})

const membershipOrderPercentage = computed(() => {
  return totalOrders.value > 0 ? Math.round((orderTypeStats.membership / totalOrders.value) * 100) : 0
})

function formatStatValue(value, format) {
  if (format === 'currency') {
    return typeof value === 'number' ? value / 100 : 0
  }
  return value
}

function activityTypeLabel(type) {
  const map = {
    'order': '订单',
    'photoset': '套图',
    'refund': '退款', 
    'register': '注册',
    'review': '审核'
  }
  return map[type] || type
}

function activityTypeTag(type) {
  const map = {
    'order': 'success',
    'photoset': 'warning',
    'refund': 'danger',
    'register': 'info',
    'review': 'primary'
  }
  return map[type] || 'info'
}

function formatCurrency(value) {
  return (value / 100).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

function trendClass(trend) {
  return trend >= 0 ? 'trend-up' : 'trend-down'
}

async function fetchStats() {
  loading.value = true
  try {
    const res = await getStats()
    Object.assign(stats, res.data)
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
}

onMounted(fetchStats)
</script>

<style scoped>
.dashboard {
  padding: 4px 0;
}

.stat-card {
  text-align: center;
  border-radius: 12px;
  transition: all 0.3s;
  min-height: 140px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-card--alert {
  border: 2px solid #F56C6C;
}

.stat-card__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 12px;
  margin: 0 auto 12px;
}

.stat-card__label {
  margin-top: 8px;
  font-size: 13px;
  color: #909399;
  padding-bottom: 4px;
}

.stat-card__trend {
  margin-top: 4px;
  font-size: 12px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  position: relative;
}

.chart-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  opacity: 0.6;
}
</style>