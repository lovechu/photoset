<template>
  <div class="tags-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>浏览标签</h1>
      <p>选择感兴趣的标签，发现更多精彩套图</p>
    </div>

    <!-- 标签列表 -->
    <div class="tags-grid" v-loading="loading">
      <div v-if="!loading && tags.length === 0" class="empty-state">
        <el-empty description="暂无标签" />
      </div>
      <div
        v-for="tag in tags"
        :key="tag.id"
        class="tag-card"
        @click="handleTagClick(tag.name)"
      >
        <div class="tag-icon">
          <el-icon :size="32"><PriceTag /></el-icon>
        </div>
        <div class="tag-name">{{ tag.name }}</div>
        <div class="tag-count" v-if="tag.count">
          {{ tag.count }} 套图
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getTags } from '@/api'
import { PriceTag } from '@element-plus/icons-vue'

const router = useRouter()
const loading = ref(false)
const tags = ref([])

const loadTags = async () => {
  loading.value = true
  try {
    const res = await getTags()
    tags.value = res.data || []
  } catch (e) {
    console.error('加载标签失败', e)
  } finally {
    loading.value = false
  }
}

const handleTagClick = (tagName) => {
  // 跳转到首页并传递标签筛选参数
  router.push({ path: '/', query: { tag: tagName } })
}

onMounted(() => {
  loadTags()
})
</script>

<style scoped>
.tags-page {
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

.tags-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 20px;
  padding: 0 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.tag-card {
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  padding: 24px 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.tag-card:hover {
  border-color: #409eff;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.15);
  transform: translateY(-2px);
}

.tag-icon {
  color: #409eff;
  margin-bottom: 12px;
}

.tag-name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.tag-count {
  font-size: 13px;
  color: #909399;
}

.empty-state {
  grid-column: 1 / -1;
  padding: 60px 0;
}

@media (max-width: 768px) {
  .page-header {
    padding: 24px 0;
  }

  .page-header h1 {
    font-size: 24px;
  }

  .tags-grid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 12px;
    padding: 0 16px;
  }

  .tag-card {
    padding: 16px 12px;
  }
}
</style>
