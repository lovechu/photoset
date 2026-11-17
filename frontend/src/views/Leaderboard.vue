<template>
  <div class="leaderboard-page">
    <el-card>
      <template #header>
        <div class="header-row">
          <h3>积分排行榜</h3>
          <el-button size="small" @click="$router.push('/user-level')">我的等级</el-button>
        </div>
      </template>
      <el-alert type="info" :closable="false" show-icon title="排行榜功能即将开放，敬请期待" v-if="!list.length && !loading" />
      <div v-for="(u, i) in list" :key="u.id || u.user_id" class="rank-item">
        <span :class="['rank-num', `rank-${i + 1}`]">{{ i + 1 }}</span>
        <el-avatar :size="40" :src="u.avatar">{{ u.nickname?.charAt(0) }}</el-avatar>
        <div class="rank-info">
          <div class="rank-name">{{ u.nickname || u.username }}</div>
          <div class="rank-level">Lv.{{ u.level || 1 }}</div>
        </div>
        <div class="rank-points">{{ u.points ?? u.total_points ?? 0 }} <span>积分</span></div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPointsLeaderboard } from '@/api'

const list = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const res = await getPointsLeaderboard()
    const data = res?.data
    list.value = Array.isArray(data) ? data : (data?.list || [])
  } catch (e) { list.value = [] }
  finally { loading.value = false }
})
</script>

<style scoped>
.leaderboard-page { max-width: 720px; margin: 0 auto; padding: 24px 0 60px; }
.header-row { display: flex; align-items: center; justify-content: space-between; }
.header-row h3 { margin: 0; }
.rank-item { display: flex; align-items: center; gap: 14px; padding: 14px 0; border-bottom: 1px solid #f0f0f0; }
.rank-num { width: 32px; text-align: center; font-weight: 700; color: #909399; font-size: 16px; }
.rank-1 { color: #f5a623; }.rank-2 { color: #909399; }.rank-3 { color: #cd7f32; }
.rank-info { flex: 1; }
.rank-name { font-weight: 600; }
.rank-level { color: #999; font-size: 12px; }
.rank-points { color: #e6a23c; font-weight: 600; font-size: 16px; }
.rank-points span { color: #999; font-size: 12px; font-weight: normal; }
</style>
