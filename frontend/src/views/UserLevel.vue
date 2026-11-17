<template>
  <div class="user-level-page">
    <div v-loading="loading">
      <!-- 等级卡片 -->
      <div v-if="levelInfo" class="level-hero">
        <div class="level-badge" :style="{ borderColor: levelInfo.level_color || '#409eff' }">
          <div class="level-num">{{ levelInfo.level }}</div>
          <div class="level-label">LEVEL</div>
        </div>
        <div class="level-info">
          <h2>{{ levelInfo.level_name || `Lv.${levelInfo.level}` }}</h2>
          <p class="points-line">当前积分 <strong>{{ levelInfo.points }}</strong></p>
          <div class="progress-wrap" v-if="levelInfo.points_to_next > 0">
            <el-progress
              :percentage="Math.min(100, Math.round((levelInfo.progress || 0) * 100))"
              :stroke-width="14"
              :show-text="false"
              color="linear-gradient(90deg, #667eea, #764ba2)"
            />
            <span class="progress-tip">距下一级还需 {{ levelInfo.points_to_next }} 积分</span>
          </div>
          <div v-else class="max-level-tip">已达到最高等级</div>
          <div class="level-stats">
            <span>🏆 {{ levelInfo.achievement_count || 0 }} 个成就</span>
          </div>
        </div>
        <div class="quick-links">
          <el-button @click="$router.push('/points-mall')">积分商城</el-button>
          <el-button @click="$router.push('/points-history')">兑换记录</el-button>
        </div>
      </div>

      <!-- 成就区 -->
      <el-card class="section-card">
        <template #header><h3>我的成就 <span class="count">{{ unlockedCount }}/{{ achievements.length }}</span></h3></template>
        <div class="achievement-grid">
          <div
            v-for="a in achievements"
            :key="a.id"
            :class="['achievement-item', { locked: !a.is_unlocked }]"
          >
            <div class="achievement-icon">{{ a.icon || '🏅' }}</div>
            <div class="achievement-name">{{ a.title || a.name }}</div>
            <div class="achievement-desc">{{ a.description }}</div>
            <div class="achievement-status">
              <el-tag v-if="a.is_unlocked" size="small" type="success">已解锁</el-tag>
              <el-tag v-else size="small" type="info">未解锁</el-tag>
              <span v-if="a.reward_points" class="reward">+{{ a.reward_points }} 积分</span>
            </div>
          </div>
        </div>
        <el-empty v-if="!achievements.length" description="暂无成就" :image-size="60" />
      </el-card>

      <!-- 等级权益 -->
      <el-card class="section-card">
        <template #header><h3>等级体系</h3></template>
        <el-table :data="levelConfigs" stripe v-loading="levelsLoading">
          <el-table-column label="等级" width="100">
            <template #default="{ row }">
              <span class="level-tag" :style="{ color: row.color || '#409eff' }">Lv.{{ row.level }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="name" label="名称" width="120" />
          <el-table-column label="积分范围" width="160">
            <template #default="{ row }">{{ row.min_points }} - {{ row.max_points || '∞' }}</template>
          </el-table-column>
          <el-table-column prop="description" label="说明" min-width="180" />
          <el-table-column label="当前" width="80">
            <template #default="{ row }">
              <el-tag v-if="row.level === levelInfo?.level" size="small" type="success">当前</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getLevelConfigs, getMyLevel, getMyAchievements } from '@/api'

const loading = ref(false)
const levelsLoading = ref(false)
const levelInfo = ref(null)
const achievements = ref([])
const levelConfigs = ref([])

const unlockedCount = computed(() => achievements.value.filter(a => a.is_unlocked).length)

onMounted(() => {
  loadLevel()
  loadAchievements()
  loadLevels()
})

async function loadLevel() {
  loading.value = true
  try {
    const res = await getMyLevel()
    levelInfo.value = res?.data || null
  } catch (e) { console.warn('加载等级信息失败:', e.message) }
  finally { loading.value = false }
}

async function loadAchievements() {
  try {
    const res = await getMyAchievements()
    const data = res?.data
    achievements.value = Array.isArray(data) ? data : []
  } catch (e) { achievements.value = [] }
}

async function loadLevels() {
  levelsLoading.value = true
  try {
    const res = await getLevelConfigs()
    const data = res?.data
    levelConfigs.value = Array.isArray(data) ? data : []
  } catch (e) { levelConfigs.value = [] }
  finally { levelsLoading.value = false }
}
</script>

<style scoped>
.user-level-page { max-width: 1000px; margin: 0 auto; padding: 24px 0 60px; }
.level-hero { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 16px; padding: 32px; color: #fff; display: flex; gap: 28px; align-items: center; margin-bottom: 20px; }
.level-badge { width: 90px; height: 90px; border: 3px solid #fff; border-radius: 50%; display: flex; flex-direction: column; align-items: center; justify-content: center; flex-shrink: 0; background: rgba(255,255,255,0.1); }
.level-num { font-size: 36px; font-weight: 700; line-height: 1; }
.level-label { font-size: 10px; letter-spacing: 2px; opacity: 0.8; }
.level-info { flex: 1; }
.level-info h2 { margin: 0 0 6px; font-size: 24px; }
.points-line { margin: 0 0 12px; opacity: 0.9; }
.points-line strong { font-size: 20px; }
.progress-wrap { max-width: 400px; }
.progress-tip { display: block; font-size: 12px; opacity: 0.8; margin-top: 6px; }
.max-level-tip { font-size: 14px; opacity: 0.85; }
.level-stats { margin-top: 10px; font-size: 13px; opacity: 0.9; }
.quick-links { display: flex; flex-direction: column; gap: 8px; }
.quick-links .el-button { margin: 0; color: #fff; border-color: rgba(255,255,255,0.5); background: rgba(255,255,255,0.1); }
.section-card { margin-bottom: 16px; }
.section-card h3 { margin: 0; font-size: 16px; }
.count { color: #999; font-size: 13px; font-weight: normal; }
.achievement-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 14px; }
.achievement-item { text-align: center; padding: 16px 10px; border: 1px solid #ebeef5; border-radius: 10px; }
.achievement-item.locked { opacity: 0.55; filter: grayscale(0.6); }
.achievement-icon { font-size: 30px; margin-bottom: 8px; }
.achievement-name { font-weight: 600; font-size: 14px; }
.achievement-desc { color: #999; font-size: 12px; margin: 6px 0; min-height: 30px; }
.achievement-status { display: flex; flex-direction: column; align-items: center; gap: 4px; }
.reward { color: #e6a23c; font-size: 11px; }
.level-tag { font-weight: 600; }
@media (max-width: 768px) { .level-hero { flex-direction: column; align-items: flex-start; } .quick-links { flex-direction: row; } }
</style>
