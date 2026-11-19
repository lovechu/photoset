<template>
  <div class="aggregate-page">
    <div class="page-header">
      <h2>标签：{{ tagName }}</h2>
      <el-button @click="$router.back()" text>返回</el-button>
    </div>
    <div v-loading="loading">
      <el-empty v-if="!loading && posts.length === 0" description="该标签暂无帖子" :image-size="80" />
      <PostCard v-for="post in posts" :key="post.id" :post="post" />
      <div class="load-more" v-if="posts.length < total">
        <el-button @click="loadMore">加载更多</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import PostCard from '@/components/PostCard.vue'
import { getPostsByTagName } from '@/api'
import { useSeo } from '@/composables/useSeo'

const route = useRoute()
const tagName = decodeURIComponent(route.params.name || '')
useSeo(() => ({
  title: tagName ? `标签 ${tagName}` : '标签',
  description: `浏览 ${tagName} 标签下的社区帖子与讨论`,
  path: route.fullPath
}))
const posts = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)

onMounted(() => load())
async function load() {
  loading.value = true
  try {
    const res = await getPostsByTagName(tagName, { page: page.value, page_size: 20 })
    const data = res?.data || {}
    const items = data.posts || []
    posts.value = page.value === 1 ? items : [...posts.value, ...items]
    total.value = data.total ?? data.pagination?.total ?? 0
  } catch (e) { console.warn('加载标签帖子失败:', e.message) }
  finally { loading.value = false }
}
function loadMore() { page.value++; load() }
</script>

<style scoped>
.aggregate-page {
  max-width: 1000px; margin: 0 auto; padding: 24px 20px 60px;
  color: #33373d; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}
.page-header {
  display: flex; align-items: center; justify-content: space-between;
  background: #fafbfc; border: 1px solid #e6e8eb; border-radius: 5px;
  padding: 14px 18px; margin-bottom: 16px;
}
.page-header h2 { margin: 0; font-size: 16px; font-weight: 700; color: #2b3038; }
.load-more { text-align: center; padding: 20px; }
</style>
