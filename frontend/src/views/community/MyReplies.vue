<template>
  <MyContentLayout title="我的回复" :fetcher="fetcher" empty-text="还没有回复">
    <template #default="{ items }">
      <div v-for="reply in items" :key="reply.id" class="reply-card" @click="goPost(reply.post_id)">
        <div class="reply-content">{{ reply.content }}</div>
        <div class="reply-meta">
          <span>{{ formatTime(reply.created_at) }}</span>
          <span>{{ reply.like_count || 0 }} 赞</span>
        </div>
      </div>
    </template>
  </MyContentLayout>
</template>

<script setup>
import MyContentLayout from './MyContentLayout.vue'
import { useRouter } from 'vue-router'
import { getMyReplies } from '@/api'
import { formatTime } from '@/utils/post'

const router = useRouter()
const fetcher = (params) => getMyReplies(params).then(res => {
  const data = res?.data || {}
  return { items: data.replies || [], total: data.pagination?.total || 0 }
})
function goPost(id) { if (id) router.push(`/community/post/${id}`) }
</script>

<style scoped>
.reply-card { background: #fff; border-radius: 10px; padding: 16px 20px; margin-bottom: 12px; cursor: pointer; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.reply-card:hover { box-shadow: 0 4px 14px rgba(0,0,0,0.08); }
.reply-content { color: #333; line-height: 1.6; }
.reply-meta { display: flex; gap: 16px; color: #999; font-size: 12px; margin-top: 8px; }
</style>
