<template>
  <div class="home-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>发现精彩套图</h1>
      <p>浏览精选摄影作品，感受视觉之美</p>
    </div>

    <!-- 搜索框 -->
    <div class="search-bar">
      <el-input
        v-model="keyword"
        placeholder="搜索套图标题或描述..."
        clearable
        @keyup.enter="handleSearch"
        @clear="handleSearch"
      >
        <template #append>
          <el-button :icon="Search" @click="handleSearch" />
        </template>
      </el-input>
    </div>

    <!-- 标签筛选 -->
    <div class="tag-filter">
      <el-radio-group v-model="selectedTag" @change="handleTagChange">
        <el-radio-button label="">全部</el-radio-button>
        <el-radio-button v-for="tag in tags" :key="tag.id" :label="tag.name">
          {{ tag.name }}
        </el-radio-button>
      </el-radio-group>
    </div>

    <!-- 套图列表 -->
    <div class="photoset-grid" v-loading="loading">
      <div v-if="!loading && photosets.length === 0" class="empty-state">
        <el-empty description="暂无套图内容" />
      </div>
      <PhotosetCard
        v-for="item in photosets"
        :key="item.id"
        :data="item"
        :search-keyword="keyword"
        class="card-hover"
      />
    </div>

    <!-- 加载更多 / 分页 -->
    <div class="pagination-wrapper" v-if="total > 0">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[12, 20, 40, 80]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { getPhotosetList, getTags } from '@/api'
import PhotosetCard from '@/components/PhotosetCard.vue'
import { Search } from '@element-plus/icons-vue'

const loading = ref(false)
const photosets = ref([])
const tags = ref([])
const selectedTag = ref('')
const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

// 加载标签列表
const loadTags = async () => {
  try {
    const res = await getTags()
    tags.value = res.data || []
  } catch (e) {
    console.error('加载标签失败', e)
  }
}

// 加载套图列表
const loadPhotosets = async () => {
  loading.value = true
  try {
    const res = await getPhotosetList({
      page: currentPage.value,
      page_size: pageSize.value,
      tag: selectedTag.value || undefined,
      keyword: keyword.value || undefined
    })
    photosets.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error('加载套图列表失败', e)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  loadPhotosets()
}

// 标签切换
const handleTagChange = () => {
  currentPage.value = 1
  loadPhotosets()
}

// 分页大小变化
const handleSizeChange = () => {
  currentPage.value = 1
  loadPhotosets()
}

// 页码变化
const handlePageChange = () => {
  loadPhotosets()
  // 滚动到顶部
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  loadTags()
  loadPhotosets()
})
</script>

<style scoped>
.home-page {
  padding: 0 0 40px;
}

.page-header {
  text-align: center;
  padding: 40px 0;
}

.page-header h1 {
  font-size: 32px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.page-header p {
  color: #666;
  font-size: 16px;
}

.search-bar {
  max-width: 480px;
  margin: 0 auto 24px;
}

.tag-filter {
  margin-bottom: 24px;
  overflow-x: auto;
  padding-bottom: 8px;
}

.tag-filter :deep(.el-radio-group) {
  display: flex;
  flex-wrap: nowrap;
  gap: 8px;
}

.tag-filter :deep(.el-radio-button__inner) {
  white-space: nowrap;
}

.photoset-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.empty-state {
  grid-column: 1 / -1;
  padding: 60px 0;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

@media (max-width: 768px) {
  .page-header {
    padding: 24px 0;
  }

  .page-header h1 {
    font-size: 24px;
  }

  .photoset-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 12px;
  }
}
</style>
