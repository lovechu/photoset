<template>
  <!-- 高级搜索侧边栏/抽屉组件 -->
  <el-drawer
    v-model="visible"
    title="高级搜索"
    size="400px"
    direction="rtl"
    :close-on-click-modal="false"
    >
    <div class="advanced-search-content">
      <!-- 价格区间 -->
      <div class="filter-section">
        <h3>价格</h3>
        <div class="price-range">
          <el-input-number
            v-model="filters.price_min"
            :min="0"
            :max="10000"
            :precision="2"
            placeholder="最低价格"
            clearable
            class="price-input"
          />
          <span class="range-separator">~</span>
          <el-input-number
            v-model="filters.price_max"
            :min="0"
            :max="10000"
            :precision="2"
            placeholder="最高价格"
            clearable
            class="price-input"
          />
        </div>
        <div class="price-options">
          <div class="checkbox-row">
            <el-checkbox v-model="filters.is_free" :true-label="true" :false-label="null" label="仅免费" />
          </div>
        </div>
      </div>

      <!-- 分类 -->
      <div class="filter-section">
        <h3>分类</h3>
        <el-select v-model="filters.category" placeholder="选择分类" clearable class="full-width">
          <el-option v-for="cat in categoryOptions" :key="cat.slug" :label="cat.name" :value="cat.slug" />
        </el-select>
      </div>

      <!-- 时间范围 -->
      <div class="filter-section">
        <h3>发布时间</h3>
        <el-select v-model="filters.time_range" placeholder="选择时间范围" clearable class="full-width">
          <el-option label="今天" value="today" />
          <el-option label="本周" value="week" />
          <el-option label="本月" value="month" />
          <el-option label="三个月内" value="quarter" />
          <el-option label="半年内" value="half_year" />
        </el-select>
      </div>

      <!-- 排序方式 -->
      <div class="filter-section">
        <h3>排序方式</h3>
        <el-select v-model="filters.sort_by" placeholder="选择排序" class="full-width">
          <el-option label="最新发布" value="latest" />
          <el-option label="最受欢迎" value="popular" />
          <el-option label="价格最低" value="price_asc" />
          <el-option label="价格最高" value="price_desc" />
          <el-option label="评分最高" value="rating" />
        </el-select>
      </div>

      <!-- 作者筛选 -->
      <div class="filter-section">
        <h3>作者</h3>
        <div class="checkbox-row">
          <el-checkbox v-model="filters.only_mine" label="仅显示我的作品" @change="handleOnlyMineChange" />
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button type="primary" :loading="loading" @click="handleApplyFilters" class="apply-btn">
          应用筛选
        </el-button>
        <el-button @click="handleResetFilters" :disabled="loading">
          重置
        </el-button>
        <el-button @click="handleCancel" plain>取消</el-button>
      </div>
      
      <!-- 已应用的筛选器显示 -->
      <div v-if="activeFilterCount > 0" class="active-filters">
        <div class="active-filters-header">
          <span>已应用的筛选条件 ({{ activeFilterCount }})</span>
          <el-button type="danger" size="small" text @click="handleClearAllFilters">清空</el-button>
        </div>
        <div class="filter-tags">
          <el-tag
            v-for="(filter, key) in activeFilterTags"
            :key="key"
            closable
            size="small"
            @close="handleRemoveFilter(key)"
            type="info"
          >
            {{ filter }}
          </el-tag>
        </div>
      </div>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, computed, watch, defineEmits, defineProps } from 'vue'
import { getCategories } from '@/api'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  currentUser: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits([
  'update:modelValue',
  'search'
])

const visible = ref(false)
const loading = ref(false)
const categoryOptions = ref([])

// 筛选器状态
const filters = ref({
  category: '',
  price_min: undefined,
  price_max: undefined,
  is_free: null,
  sort_by: 'latest',
  time_range: '',
  only_mine: false,
  user_id: undefined
})

// 监听外部visible变化
watch(() => props.modelValue, (val) => {
  visible.value = val
})

// 监听内部visible变化
watch(visible, (val) => {
  emit('update:modelValue', val)
  if (val && categoryOptions.value.length === 0) {
    loadCategories()
  }
})

// 动态加载分类
const loadCategories = async () => {
  try {
    const res = await getCategories()
    categoryOptions.value = res.data || []
  } catch (e) {
    console.error('加载分类失败', e)
  }
}

// 统计已激活的筛选条件数量
const activeFilterCount = computed(() => {
  let count = 0
  if (filters.value.category) count++
  if (filters.value.price_min || filters.value.price_max) count++
  if (filters.value.is_free !== null) count++
  if (filters.value.sort_by !== 'latest') count++
  if (filters.value.time_range) count++
  if (filters.value.only_mine) count++
  return count
})

// 显示活跃的筛选标签
const activeFilterTags = computed(() => {
  const tags = {}
  
  if (filters.value.category) {
    const foundCat = categoryOptions.value.find(c => c.slug === filters.value.category)
    tags.category = `分类: ${foundCat ? foundCat.name : filters.value.category}`
  }
  
  if (filters.value.price_min || filters.value.price_max) {
    const min = filters.value.price_min ? `¥${filters.value.price_min}` : '不限'
    const max = filters.value.price_max ? `¥${filters.value.price_max}` : '不限'
    tags.price = `价格: ${min} - ${max}`
  }
  
  if (filters.value.is_free !== null) {
    tags.is_free = filters.value.is_free ? '仅免费' : '付费作品'
  }
  
  if (filters.value.sort_by !== 'latest') {
    const sortMap = {
      'latest': '最新发布',
      'popular': '最受欢迎',
      'price_asc': '价格最低',
      'price_desc': '价格最高',
      'rating': '评分最高'
    }
    tags.sort_by = `排序: ${sortMap[filters.value.sort_by] || filters.value.sort_by}`
  }
  
  if (filters.value.time_range) {
    const timeMap = {
      'today': '今天',
      'week': '本周',
      'month': '本月',
      'quarter': '三个月内',
      'half_year': '半年内'
    }
    tags.time_range = `时间: ${timeMap[filters.value.time_range] || filters.value.time_range}`
  }
  
  if (filters.value.only_mine) {
    tags.only_mine = '仅我的作品'
  }
  
  return tags
})

// "仅显示我的作品"变化处理
const handleOnlyMineChange = (value) => {
  if (value && props.currentUser?.id) {
    filters.value.user_id = props.currentUser.id
  } else {
    filters.value.user_id = undefined
  }
}

// 应用筛选
const handleApplyFilters = async () => {
  loading.value = true
  try {
    // 组装筛选参数
    const searchParams = {
      category: filters.value.category || undefined,
      price_min: filters.value.price_min > 0 ? filters.value.price_min : undefined,
      price_max: filters.value.price_max > 0 ? filters.value.price_max : undefined,
      is_free: filters.value.is_free,
      sort_by: filters.value.sort_by,
      time_range: filters.value.time_range || undefined,
      user_id: filters.value.user_id || undefined,
      only_mine: filters.value.only_mine ? true : undefined
    }
    
    // 移除undefined和null值
    Object.keys(searchParams).forEach(key => {
      if (searchParams[key] === undefined || searchParams[key] === null) {
        delete searchParams[key]
      }
    })
    
    emit('search', searchParams)
    visible.value = false
  } catch (error) {
    console.error('应用筛选失败:', error)
  } finally {
    loading.value = false
  }
}

// 重置筛选
const handleResetFilters = () => {
  filters.value = {
    category: '',
    price_min: undefined,
    price_max: undefined,
    is_free: null,
    sort_by: 'latest',
    time_range: '',
    only_mine: false,
    user_id: undefined
  }
}

// 取消
const handleCancel = () => {
  visible.value = false
}

// 清除特定筛选
const handleRemoveFilter = (key) => {
  switch (key) {
    case 'category':
      filters.value.category = ''
      break
    case 'price':
      filters.value.price_min = undefined
      filters.value.price_max = undefined
      break
    case 'is_free':
      filters.value.is_free = null
      break
    case 'sort_by':
      filters.value.sort_by = 'latest'
      break
    case 'time_range':
      filters.value.time_range = ''
      break
    case 'only_mine':
      filters.value.only_mine = false
      filters.value.user_id = undefined
      break
  }
}

// 清除所有筛选
const handleClearAllFilters = () => {
  handleResetFilters()
}
</script>

<style scoped>
.advanced-search-content {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.filter-section {
  margin-bottom: 20px;
}

.filter-section h3 {
  font-size: 14px;
  font-weight: 600;
  color: #606266;
  margin-bottom: 12px;
}

.price-range {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.price-input {
  flex: 1;
}

.range-separator {
  color: #c0c4cc;
  font-weight: 500;
}

.checkbox-row {
  margin-bottom: 8px;
}

.checkbox-row:last-child {
  margin-bottom: 0;
}

.full-width {
  width: 100%;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: auto;
  padding-top: 20px;
}

.apply-btn {
  width: 100%;
}

.active-filters {
  margin-top: 20px;
  padding: 16px;
  background: #f7f9fc;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.active-filters-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 13px;
  color: #606266;
}

.filter-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

:deep(.el-drawer__header) {
  margin-bottom: 0;
  padding: 20px;
  border-bottom: 1px solid #e4e7ed;
}

:deep(.el-drawer__body) {
  padding: 0;
}

@media (max-width: 768px) {
  .advanced-search-content {
    padding: 16px;
  }
}
</style>