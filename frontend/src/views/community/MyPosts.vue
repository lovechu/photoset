<template>
  <MyContentLayout title="我的帖子" :fetcher="fetcher" empty-text="还没有发帖">
    <template #default="{ items }">
      <PostCard v-for="post in items" :key="post.id" :post="post" />
    </template>
  </MyContentLayout>
</template>

<script setup>
import MyContentLayout from './MyContentLayout.vue'
import PostCard from '@/components/PostCard.vue'
import { getMyPosts } from '@/api'

const fetcher = (params) => getMyPosts(params).then(res => {
  const data = res?.data || {}
  return { items: data.posts || [], total: data.pagination?.total || 0 }
})
</script>
