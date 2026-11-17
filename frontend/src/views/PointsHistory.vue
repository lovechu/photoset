<template>
  <div class="points-history-page">
    <el-card>
      <template #header>
        <div class="header-row">
          <h3>兑换记录</h3>
          <el-button size="small" @click="$router.push('/points-mall')">去商城</el-button>
        </div>
      </template>
      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="item_name" label="商品" min-width="160" />
        <el-table-column prop="points" label="消耗积分" width="120" />
        <el-table-column prop="item_type" label="类型" width="120" />
        <el-table-column prop="created_at" label="兑换时间" width="180" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'completed' ? 'success' : 'info'" size="small">
              {{ row.status === 'completed' ? '已完成' : row.status === 'refunded' ? '已退款' : row.status }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!loading && !list.length" description="暂无兑换记录" :image-size="60" />
      <div class="pagination" v-if="total > pageSize">
        <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="prev, pager, next" @current-change="load" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getPointsExchangeHistory } from '@/api'

const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)

onMounted(() => load())
async function load() {
  loading.value = true
  try {
    const res = await getPointsExchangeHistory({ page: page.value, page_size: pageSize })
    const data = res?.data || {}
    list.value = data.exchanges || []
    total.value = data.pagination?.total || 0
  } catch (e) { console.warn('加载兑换记录失败:', e.message) }
  finally { loading.value = false }
}
</script>

<style scoped>
.points-history-page { max-width: 900px; margin: 0 auto; padding: 24px 0 60px; }
.header-row { display: flex; align-items: center; justify-content: space-between; }
.header-row h3 { margin: 0; }
.pagination { display: flex; justify-content: center; margin-top: 20px; }
</style>
