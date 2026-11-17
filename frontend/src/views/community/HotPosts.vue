<template>
  <div class="hot-posts-page">
    <div class="page-header">
      <h2>热门帖子</h2>
    </div>
    <div v-loading="loading">
      <el-empty v-if="!loading && posts.length === 0" description="暂无热门帖子" :image-size="80" />
      <PostCard v-for="(post, i) in posts" :key="post.id" :post="post" :rank="i + 1" />
      <div class="load-more" v-if="posts.length < total">
        <el-button @click="loadMore">加载更多</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import PostCard from '@/components/PostCard.vue'
import { getHotPosts } from '@/api'

const posts = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)

onMounted(() => load())
function reload() { page.value = 1; posts.value = []; load() }
async function load() {
  loading.value = true
  try {
    const res = await getHotPosts({ page: page.value, page_size: 20 })
    const data = res?.data || {}
    const items = data.posts || []
    posts.value = page.value === 1 ? items : [...posts.value, ...items]
    total.value = data.pagination?.total || 0
  } catch (e) { console.warn('加载热门帖子失败:', e.message) }
  finally { loading.value = false }
}
function loadMore() { page.value++; load() }
</script>

<style scoped>
.hot-posts-page {
  max-width: 1000px; margin: 0 auto; padding: 24px 20px 60px;
  color: #33373d; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}
.page-header {
  display: flex; align-items: center; justify-content: space-between;
  background: #fafbfc; border: 1px solid #e6e8eb; border-radius: 5px;
  padding: 14px 18px; margin-bottom: 16px;
}
.page-header h2 { margin: 0; font-size: 16px; font-weight: 700; color: #2b3038; }
.load-more { text-align: center; padding: 24px 0; }
.load-more :deep(.el-button) { border-color: #e0e3e7; color: #6b7280; border-radius: 4px; }
.load-more :deep(.el-button:hover) { border-color: #b8c0cc; color: #2b3038; background: #fff; }
</style>
