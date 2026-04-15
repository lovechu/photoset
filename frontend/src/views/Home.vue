<template>
  <div class="home-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>发现精彩套图</h1>
      <p>浏览精选摄影作品，感受视觉之美</p>
    </div>

    <!-- 搜索区域 -->
    <div class="search-section">
      <!-- 搜索框 -->
      <div class="search-bar">
        <el-input
          v-model="keyword"
          placeholder="搜索套图标题或描述..."
          clearable
          @keyup.enter="handleSearch"
          @clear="handleSearch"
          class="main-search"
        >
          <template #append>
            <el-button :icon="Search" @click="handleSearch" />
          </template>
        </el-input>
      </div>

      <!-- 高级搜索按钮 -->
      <div class="advanced-search-button">
        <el-button type="primary" @click="showAdvancedSearch = true">
          <template #icon>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16">
              <path fill="currentColor" d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm-2 16c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"/>
            </svg>
          </template>
          高级搜索
        </el-button>
      </div>
    </div>

    <!-- 标签筛选 -->
    <div class="tag-filter">
      <el-radio-group v-model="selectedTag" @change="handleTagChange">
        <el-radio-button value="">全部</el-radio-button>
        <el-radio-button v-for="tag in tags" :key="tag.id" :value="tag.name">
          {{ tag.name }}
        </el-radio-button>
      </el-radio-group>
    </div>

    <!-- 筛选状态显示 -->
    <div v-if="activeFilterText" class="filter-status">
      <div class="filter-status-content">
        <span class="filter-label">筛选条件:</span>
        <span class="filter-text">{{ activeFilterText }}</span>
        <el-button 
          v-if="hasAdvancedFilters" 
          type="danger" 
          size="small" 
          text 
          @click="clearAdvancedFilters"
          class="clear-filters-btn"
        >
          清除筛选
        </el-button>
      </div>
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

    <!-- 高级搜索抽屉 -->
    <AdvancedSearch
      v-model="showAdvancedSearch"
      :current-user="currentUser"
      @search="handleAdvancedSearch"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { getPhotosetList, getTags, getPhotosetListAdvanced, getCategories } from '@/api'
import { getCurrentUser } from '@/api'
import PhotosetCard from '@/components/PhotosetCard.vue'
import AdvancedSearch from '@/components/AdvancedSearch.vue'
import { Search } from '@element-plus/icons-vue'

const loading = ref(false)
const photosets = ref([])
const tags = ref([])
const categoryOptions = ref([])
const selectedTag = ref('')
const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const showAdvancedSearch = ref(false)
const currentUser = ref({})
const advancedFilters = ref({})

// 加载标签列表
const loadTags = async () => {
  try {
    const res = await getTags()
    tags.value = res.data || []
  } catch (e) {
    console.error('加载标签失败', e)
  }
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const res = await getCategories()
    categoryOptions.value = res.data || []
  } catch (e) {
    console.error('加载分类失败', e)
  }
}

// 加载用户信息
const loadCurrentUser = async () => {
  try {
    const res = await getCurrentUser()
    currentUser.value = res.data.user || {}
  } catch (e) {
    console.warn('获取用户信息失败或未登录')
    currentUser.value = {}
  }
}

// 检查是否使用高级搜索
const hasAdvancedFilters = computed(() => {
  return Object.keys(advancedFilters.value).length > 0
})

// 加载套图列表
const loadPhotosets = async () => {
  loading.value = true
  try {
    let res
    
    if (hasAdvancedFilters.value) {
      // 使用高级搜索
      const params = {
        page: currentPage.value,
        page_size: pageSize.value,
        tag: selectedTag.value || undefined,
        keyword: keyword.value || undefined,
        ...advancedFilters.value
      }
      
      // 移除undefined和null值
      Object.keys(params).forEach(key => {
        if (params[key] === undefined || params[key] === null) {
          delete params[key]
        }
      })
      
      res = await getPhotosetListAdvanced(params)
    } else {
      // 使用基础搜索
      res = await getPhotosetList({
        page: currentPage.value,
        page_size: pageSize.value,
        tag: selectedTag.value || undefined,
        keyword: keyword.value || undefined
      })
    }
    
    photosets.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error('加载套图列表失败', e)
  } finally {
    loading.value = false
  }
}

// 处理高级搜索
const handleAdvancedSearch = (filters) => {
  advancedFilters.value = filters
  currentPage.value = 1 // 重置到第一页
  loadPhotosets()
}

// 清除高级筛选
const clearAdvancedFilters = () => {
  advancedFilters.value = {}
  currentPage.value = 1
  loadPhotosets()
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

// 显示当前筛选状态
const activeFilterText = computed(() => {
  const filters = []
  
  if (selectedTag.value) {
    filters.push(`标签: ${selectedTag.value}`)
  }
  
  if (keyword.value) {
    filters.push(`关键词: ${keyword.value}`)
  }
  
  // 高级筛选
  if (advancedFilters.value.category) {
    const foundCat = categoryOptions.value.find(c => c.slug === advancedFilters.value.category)
    filters.push(`分类: ${foundCat ? foundCat.name : advancedFilters.value.category}`)
  }
  
  if (advancedFilters.value.is_free !== undefined && advancedFilters.value.is_free !== null) {
    filters.push(advancedFilters.value.is_free ? '仅免费' : '付费作品')
  }
  
  if (advancedFilters.value.sort_by && advancedFilters.value.sort_by !== 'latest') {
    const sortMap = {
      'popular': '最受欢迎',
      'price_asc': '价格最低',
      'price_desc': '价格最高',
      'rating': '评分最高'
    }
    filters.push(`排序: ${sortMap[advancedFilters.value.sort_by] || advancedFilters.value.sort_by}`)
  }
  
  if (advancedFilters.value.time_range) {
    const timeMap = {
      'today': '今天',
      'week': '本周',
      'month': '本月',
      'quarter': '三个月内',
      'half_year': '半年内'
    }
    filters.push(`时间: ${timeMap[advancedFilters.value.time_range] || advancedFilters.value.time_range}`)
  }
  
  if (advancedFilters.value.only_mine) {
    filters.push('仅我的作品')
  }
  
  return filters.join(' · ')
})

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

// 初始化加载
onMounted(() => {
  loadCurrentUser()
  loadTags()
  loadCategories()
  loadPhotosets()
})

// 监听高级筛选变化
watch(advancedFilters, () => {
  if (hasAdvancedFilters.value) {
    currentPage.value = 1
    loadPhotosets()
  }
}, { deep: true })
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

.search-section {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  gap: 16px;
  margin: 0 auto 24px;
  max-width: 800px;
}

.search-bar {
  flex: 1;
  max-width: 560px;
}

.main-search {
  width: 100%;
}

.advanced-search-button {
  flex-shrink: 0;
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

.filter-status {
  margin-bottom: 24px;
  padding: 12px 16px;
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border-radius: 8px;
  border: 1px solid rgba(0, 150, 255, 0.1);
}

.filter-status-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.filter-label {
  font-size: 14px;
  font-weight: 600;
  color: #0369a1;
  white-space: nowrap;
}

.filter-text {
  flex: 1;
  font-size: 14px;
  color: #0c4a6e;
  line-height: 1.5;
}

.clear-filters-btn {
  flex-shrink: 0;
  font-size: 13px;
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

  .search-section {
    flex-direction: column;
    padding: 0 16px;
  }

  .search-bar, .advanced-search-button {
    width: 100%;
    max-width: 100%;
  }

  .advanced-search-button .el-button {
    width: 100%;
  }

  .filter-status-content {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .filter-text {
    width: 100%;
    word-break: break-word;
  }

  .clear-filters-btn {
    align-self: flex-end;
  }

  .photoset-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 12px;
  }
}

@media (max-width: 480px) {
  .tag-filter :deep(.el-radio-group) {
    flex-wrap: wrap;
  }

  .filter-status {
    padding: 10px 12px;
  }
}
</style>
