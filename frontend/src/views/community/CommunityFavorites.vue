<template>
  <MyContentLayout title="社区收藏" :fetcher="fetcher" empty-text="还没有收藏的帖子">
    <template #default="{ items }">
      <PostCard v-for="post in items" :key="post.id" :post="post" />
    </template>
  </MyContentLayout>
</template>

<script setup>
import MyContentLayout from './MyContentLayout.vue'
import PostCard from '@/components/PostCard.vue'
import { getMyCommunityFavorites } from '@/api'

const fetcher = (params) => getMyCommunityFavorites(params).then(res => {
  const data = res?.data || {}
  return { items: data.posts || [], total: data.pagination?.total || 0 }
})
</script>
