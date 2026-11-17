<template>
  <div class="login-history-page">
    <el-card>
      <template #header><h3>登录历史</h3></template>
      <el-table :data="list" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="login_at" label="登录时间" width="180" />
        <el-table-column prop="ip" label="IP 地址" width="140" />
        <el-table-column prop="ip_location" label="登录地点" min-width="140" />
        <el-table-column prop="device" label="设备" width="120" />
        <el-table-column prop="browser" label="浏览器" width="120" />
        <el-table-column prop="os" label="系统" width="120" />
        <el-table-column prop="login_type" label="方式" width="100" />
        <el-table-column label="结果" width="90">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'" size="small">
              {{ row.success ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination" v-if="total > pageSize">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadHistory"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getLoginHistory } from '@/api'

const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)

onMounted(() => loadHistory())

async function loadHistory() {
  loading.value = true
  try {
    const res = await getLoginHistory({ page: page.value, page_size: pageSize })
    const data = res?.data || {}
    list.value = data.list || []
    total.value = data.total || 0
  } catch (e) {
    console.warn('加载登录历史失败:', e.message)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-history-page { max-width: 1100px; margin: 0 auto; padding: 24px 0 60px; }
.el-card__header h3 { margin: 0; font-size: 18px; }
.pagination { display: flex; justify-content: center; margin-top: 20px; }
</style>
