<template>
  <div class="points-mall-page">
    <div class="page-header">
      <div>
        <h2>积分商城</h2>
        <p v-if="myPoints" class="points-balance">我的积分：<strong>{{ myPoints.points }}</strong> <el-tag size="small">{{ myPoints.level_name }}</el-tag></p>
      </div>
      <div>
        <el-button @click="$router.push('/points-history')">兑换记录</el-button>
        <el-button @click="$router.push('/user-level')">我的等级</el-button>
      </div>
    </div>

    <div class="category-tabs">
      <el-radio-group v-model="activeCategory" @change="loadItems">
        <el-radio-button value="">全部</el-radio-button>
        <el-radio-button v-for="c in categories" :key="c.key" :value="c.key">{{ c.name }}</el-radio-button>
      </el-radio-group>
    </div>

    <div v-loading="loading">
      <el-empty v-if="!loading && items.length === 0" description="暂无可兑换商品" :image-size="80" />
      <div class="items-grid">
        <el-card v-for="item in items" :key="item.id" class="item-card" shadow="hover">
          <div class="item-image">
            <el-image v-if="item.image" :src="item.image" fit="cover" class="image" />
            <div v-else class="image-placeholder"><el-icon :size="36"><Present /></el-icon></div>
          </div>
          <div class="item-body">
            <div class="item-name">{{ item.name }}</div>
            <div class="item-desc">{{ item.description }}</div>
            <div class="item-meta">
              <span class="cost">{{ item.points_cost }} 积分</span>
              <span v-if="item.min_level" class="min-level">Lv.{{ item.min_level }} 起</span>
            </div>
            <div class="stock">
              <span v-if="item.is_unlimited">不限量</span>
              <span v-else>剩余 {{ Math.max(0, item.total_stock - item.used_stock) }} 件</span>
            </div>
            <el-button
              type="primary"
              size="small"
              style="width:100%;margin-top:10px"
              :disabled="!canExchange(item)"
              @click="exchange(item)"
            >{{ canExchange(item) ? '立即兑换' : '积分不足' }}</el-button>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Present } from '@element-plus/icons-vue'
import {
  getPointsMallItems, getPointsMallCategories, exchangePointsItem, getMyPoints
} from '@/api'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const items = ref([])
const categories = ref([])
const activeCategory = ref('')
const myPoints = ref(null)
const loading = ref(false)

onMounted(() => {
  loadCategories()
  loadItems()
  if (userStore.isLoggedIn) loadMyPoints()
})

async function loadCategories() {
  try {
    const res = await getPointsMallCategories()
    const data = res?.data
    categories.value = Array.isArray(data) ? data : []
  } catch (e) { categories.value = [] }
}

async function loadItems() {
  loading.value = true
  try {
    const res = await getPointsMallItems(activeCategory.value ? { category: activeCategory.value } : {})
    const data = res?.data
    items.value = Array.isArray(data) ? data : []
  } catch (e) { items.value = [] }
  finally { loading.value = false }
}

async function loadMyPoints() {
  try {
    const res = await getMyPoints()
    myPoints.value = res?.data || null
  } catch (e) { /* 静默 */ }
}

function canExchange(item) {
  if (!userStore.isLoggedIn) return false
  const points = myPoints.value?.points ?? userStore.user?.points ?? 0
  return points >= item.points_cost && (item.is_unlimited || item.total_stock - item.used_stock > 0) && (!item.min_level || (userStore.user?.level || 1) >= item.min_level)
}

async function exchange(item) {
  try {
    await ElMessageBox.confirm(`确定用 ${item.points_cost} 积分兑换「${item.name}」？`, '兑换确认', { type: 'warning' })
    await exchangePointsItem({ item_id: item.id })
    ElMessage.success('兑换成功')
    loadMyPoints()
    loadItems()
  } catch (e) { /* 取消或失败 */ }
}
</script>

<style scoped>
.points-mall-page { max-width: 1100px; margin: 0 auto; padding: 24px 0 60px; }
.page-header { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 16px; }
.page-header h2 { margin: 0 0 6px; }
.points-balance { margin: 0; color: #606266; }
.points-balance strong { color: #e6a23c; font-size: 18px; }
.category-tabs { margin-bottom: 20px; }
.items-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 16px; }
.item-card { overflow: hidden; }
.item-image { height: 130px; background: #f5f7fa; }
.image { width: 100%; height: 100%; }
.image-placeholder { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; color: #c0c4cc; background: linear-gradient(135deg, #f5f7fa, #e4e7ed); }
.item-name { font-weight: 600; margin-bottom: 4px; }
.item-desc { color: #999; font-size: 12px; min-height: 32px; line-height: 1.4; }
.item-meta { display: flex; align-items: center; justify-content: space-between; margin-top: 8px; }
.cost { color: #e6a23c; font-weight: 600; }
.min-level { color: #909399; font-size: 11px; }
.stock { color: #bbb; font-size: 11px; margin-top: 4px; }
</style>
