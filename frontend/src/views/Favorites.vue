<template>
  <div class="favorites-page">
    <div class="page-header">
      <h1>我的收藏</h1>
      <p>您收藏的套图都在这里</p>
    </div>
    <div class="photoset-grid" v-loading="loading">
      <div v-if="!loading && photosets.length === 0" class="empty-state">
        <el-empty description="暂无收藏">
          <el-button type="primary" @click="$router.push('/')">去发现套图</el-button>
        </el-empty>
      </div>
      <PhotosetCard
        v-for="item in photosets"
        :key="item.id"
        :data="item"
        class="card-hover"
      />
    </div>
    <div class="pagination-wrapper" v-if="total > 0">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[12, 20, 40]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadFavorites"
        @current-change="loadFavorites"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getFavorites } from '@/api'
import PhotosetCard from '@/components/PhotosetCard.vue'

const loading = ref(false)
const photosets = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(12)

const loadFavorites = async () => {
  loading.value = true
  try {
    const res = await getFavorites({
      page: currentPage.value,
      page_size: pageSize.value
    })
    photosets.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error('加载收藏失败', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadFavorites()
})
</script>

<style scoped>
.favorites-page { padding: 0 0 40px; }
.page-header { text-align: center; padding: 40px 0; }
.page-header h1 { font-size: 28px; font-weight: 600; color: #1a1a1a; margin-bottom: 8px; }
.page-header p { color: #666; font-size: 15px; }
.photoset-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}
.empty-state { grid-column: 1 / -1; padding: 60px 0; }
.pagination-wrapper { display: flex; justify-content: center; margin-top: 40px; }
@media (max-width: 768px) {
  .photoset-grid { grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 12px; }
}
</style>
