<template>
  <div class="static-page">
    <el-card>
      <template #header>
        <div class="page-header">
          <h1>{{ pageTitle }}</h1>
        </div>
      </template>

      <div class="page-content" v-if="contentHtml">
        <div v-html="contentHtml"></div>
      </div>
      <div class="page-empty" v-else>
        <el-empty description="页面内容正在完善中" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import { getPageBySlug } from '@/api/pages'

const route = useRoute()
const pageSlug = computed(() => route.params.pageType || 'about')
const pageData = ref(null)
const pageTitle = ref('')
const contentHtml = ref('')
const loading = ref(false)

const fetchPage = async (slug) => {
  loading.value = true
  try {
    const res = await getPageBySlug(slug)
    pageData.value = res.data
    pageTitle.value = pageData.value?.title || slug
    const rawContent = pageData.value?.content_md || ''
    
    if (rawContent && rawContent.trim()) {
      // 如果内容是 Markdown，转换为 HTML
      if (rawContent.includes('#') || rawContent.includes('```') || rawContent.includes('- ')) {
        contentHtml.value = marked.parse(rawContent)
      } else {
        contentHtml.value = rawContent
      }
    } else {
      contentHtml.value = ''
    }
  } catch (e) {
    console.warn('页面加载失败:', e)
    contentHtml.value = ''
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchPage(pageSlug.value)
})

// 监听路由参数变化，点击链接时重新获取内容
watch(pageSlug, (newSlug) => {
  if (newSlug) {
    fetchPage(newSlug)
  }
})
</script>

<style scoped>
.static-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-content {
  min-height: 500px;
  line-height: 1.7;
}

.page-content :deep(h1) {
  font-size: 24px;
  margin: 20px 0 16px;
  color: #333;
}

.page-content :deep(h2) {
  font-size: 20px;
  margin: 18px 0 14px;
  color: #444;
}

.page-content :deep(h3) {
  font-size: 18px;
  margin: 16px 0 12px;
  color: #555;
}

.page-content :deep(p) {
  margin: 12px 0;
  color: #666;
}

.page-content :deep(ul),
.page-content :deep(ol) {
  margin: 12px 0;
  padding-left: 24px;
}

.page-content :deep(li) {
  margin: 6px 0;
}

.page-content :deep(code) {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
}

.page-content :deep(pre) {
  background: #f8f8f8;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 16px 0;
}

.page-content :deep(a) {
  color: #409eff;
  text-decoration: none;
}

.page-content :deep(a:hover) {
  text-decoration: underline;
}
</style>