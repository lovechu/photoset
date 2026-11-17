<template>
  <div class="my-likes-page">
    <div class="page-header">
      <h2>我的点赞</h2>
      <el-button @click="$router.push('/community/center')" text>社区中心</el-button>
    </div>
    <el-alert
      type="info"
      :closable="false"
      show-icon
      title="点赞记录会在浏览帖子时自动记录，当前展示最近点赞的帖子"
      style="margin-bottom: 16px"
    />
    <div v-loading="loading">
      <el-empty v-if="!loading && posts.length === 0" description="还没有点赞记录" :image-size="80">
        <el-button type="primary" size="small" @click="$router.push('/community')">去社区逛逛</el-button>
      </el-empty>
      <PostCard v-for="post in posts" :key="post.id" :post="post" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import PostCard from '@/components/PostCard.vue'
import { getPosts } from '@/api'

const posts = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    // 后端暂无「我的点赞」专属接口，拉取帖子列表并筛选本地记录过的点赞
    const liked = JSON.parse(localStorage.getItem('photoset_liked_posts') || '[]')
    if (liked.length) {
      const res = await getPosts({ page: 1, page_size: 50 })
      const all = res?.data?.posts || []
      posts.value = all.filter(p => liked.includes(p.id)).map(p => ({ ...p, is_liked: true }))
    }
  } catch (e) { console.warn('加载点赞记录失败:', e.message) }
  finally { loading.value = false }
})
</script>

<style scoped>
.my-likes-page { max-width: 900px; margin: 0 auto; padding: 24px 0 60px; }
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px; }
.page-header h2 { margin: 0; }
</style>
