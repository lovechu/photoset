<template>
  <div class="creator-dashboard">
    <div class="page-header">
      <div>
        <h2>创作者中心</h2>
        <p class="subtitle">查看你的作品数据与表现</p>
      </div>
      <div>
        <el-button @click="$router.push('/my-photosets')">我的套图</el-button>
        <el-button type="primary" :icon="Plus" @click="$router.push('/create')">创建套图</el-button>
      </div>
    </div>

    <div v-loading="loading">
      <!-- 数据概览 -->
      <el-row :gutter="16" v-if="stats">
        <el-col :xs="12" :sm="8" :md="6" v-for="card in statCards" :key="card.label">
          <el-card class="stat-card" shadow="hover">
            <div class="stat-icon" :style="{ color: card.color }"><el-icon :size="24"><component :is="card.icon" /></el-icon></div>
            <div class="stat-value">{{ formatNumber(card.value) }}</div>
            <div class="stat-label">{{ card.label }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 趋势图 -->
      <el-card class="section-card" v-if="dailyStats.length">
        <template #header>
          <div class="header-row">
            <h3>浏览趋势</h3>
            <el-select v-model="days" size="small" style="width:110px" @change="loadDaily">
              <el-option label="7天" :value="7" />
              <el-option label="30天" :value="30" />
              <el-option label="90天" :value="90" />
            </el-select>
          </div>
        </template>
        <div class="chart-bars">
          <div v-for="d in dailyStats" :key="d.date" class="bar-col" :title="`${d.date}: ${d.views} 浏览`">
            <div class="bar" :style="{ height: barHeight(d.views) }" />
            <div class="bar-label">{{ d.date.slice(5) }}</div>
          </div>
        </div>
      </el-card>

      <!-- 热门套图 -->
      <el-card class="section-card" v-if="photosetStats.length">
        <template #header><h3>套图表现</h3></template>
        <el-table :data="photosetStats" stripe>
          <el-table-column label="套图" min-width="200">
            <template #default="{ row }">
              <div class="photoset-cell">
                <el-image v-if="row.cover_image" :src="row.cover_image" fit="cover" class="thumb" />
                <span>{{ row.title }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="view_count" label="浏览" width="100" />
          <el-table-column prop="download_count" label="下载" width="100" />
          <el-table-column prop="fav_count" label="收藏" width="100" />
          <el-table-column prop="order_count" label="订单" width="100" />
          <el-table-column label="收入" width="100">
            <template #default="{ row }">¥{{ Number(row.revenue || 0).toFixed(2) }}</template>
          </el-table-column>
        </el-table>
      </el-card>

      <el-empty v-if="!loading && !stats" description="暂无统计数据" :image-size="80" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Plus, View, Download, Star, Money, Document } from '@element-plus/icons-vue'
import { getCreatorStats, getCreatorDailyStats, getCreatorPhotosetStats } from '@/api'

const loading = ref(false)
const stats = ref(null)
const dailyStats = ref([])
const photosetStats = ref([])
const days = ref(30)

const statCards = computed(() => stats.value ? [
  { label: '套图总数', value: stats.value.total_photosets, icon: Document, color: '#409eff' },
  { label: '累计浏览', value: stats.value.total_views, icon: View, color: '#67c23a' },
  { label: '累计下载', value: stats.value.total_downloads, icon: Download, color: '#e6a23c' },
  { label: '累计收藏', value: stats.value.total_favorites, icon: Star, color: '#f56c6c' },
  { label: '今日浏览', value: stats.value.today_views, icon: View, color: '#909399' },
  { label: '本周浏览', value: stats.value.week_views, icon: View, color: '#9b59b6' },
  { label: '累计订单', value: stats.value.total_orders, icon: Document, color: '#16a085' },
  { label: '累计收入', value: stats.value.total_revenue, icon: Money, color: '#d35400', money: true }
] : [])

onMounted(() => {
  loadStats()
  loadDaily()
  loadPhotosets()
})

async function loadStats() {
  loading.value = true
  try {
    const res = await getCreatorStats()
    stats.value = res?.data || null
  } catch (e) { console.warn('加载统计失败:', e.message) }
  finally { loading.value = false }
}

async function loadDaily() {
  try {
    const res = await getCreatorDailyStats({ days: days.value })
    const data = res?.data
    dailyStats.value = Array.isArray(data) ? data : []
  } catch (e) { dailyStats.value = [] }
}

async function loadPhotosets() {
  try {
    const res = await getCreatorPhotosetStats({ limit: 10 })
    const data = res?.data
    photosetStats.value = Array.isArray(data) ? data : []
  } catch (e) { photosetStats.value = [] }
}

function formatNumber(v) {
  if (v === undefined || v === null) return '0'
  return Number(v).toLocaleString('zh-CN')
}

function barHeight(views) {
  const max = Math.max(...dailyStats.value.map(d => d.views || 0), 1)
  return `${Math.max(4, ((views || 0) / max) * 140)}px`
}
</script>

<style scoped>
.creator-dashboard { max-width: 1100px; margin: 0 auto; padding: 24px 0 60px; }
.page-header { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 20px; }
.page-header h2 { margin: 0 0 4px; }
.subtitle { color: #999; margin: 0; font-size: 13px; }
.stat-card { text-align: center; margin-bottom: 16px; }
.stat-icon { margin-bottom: 6px; }
.stat-value { font-size: 24px; font-weight: 700; color: #303133; }
.stat-label { font-size: 12px; color: #999; margin-top: 4px; }
.section-card { margin-bottom: 16px; }
.section-card h3, .header-row h3 { margin: 0; font-size: 16px; }
.header-row { display: flex; align-items: center; justify-content: space-between; }
.chart-bars { display: flex; align-items: flex-end; gap: 4px; height: 180px; overflow-x: auto; padding-top: 10px; }
.bar-col { display: flex; flex-direction: column; align-items: center; min-width: 28px; flex: 1; }
.bar { width: 100%; background: linear-gradient(180deg, #409eff, #79bbff); border-radius: 4px 4px 0 0; }
.bar-label { font-size: 10px; color: #999; margin-top: 6px; white-space: nowrap; }
.photoset-cell { display: flex; align-items: center; gap: 8px; }
.thumb { width: 40px; height: 30px; border-radius: 4px; }
</style>
