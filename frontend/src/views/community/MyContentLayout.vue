<template>
  <div class="my-content-page">
    <div class="page-header">
      <h2>{{ title }}</h2>
      <el-button @click="$router.push('/community/center')" text>社区中心</el-button>
    </div>
    <div v-loading="loading">
      <el-empty v-if="!loading && items.length === 0" :description="emptyText" :image-size="80" />
      <slot :items="items" />
      <div class="load-more" v-if="items.length < total">
        <el-button @click="loadMore" :loading="loadingMore">加载更多</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  title: { type: String, required: true },
  emptyText: { type: String, default: '暂无内容' },
  fetcher: { type: Function, required: true }
})

const items = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)
const loadingMore = ref(false)

onMounted(() => load())

async function load() {
  loading.value = true
  try {
    const res = await props.fetcher({ page: page.value, page_size: 20 })
    items.value = page.value === 1 ? res.items : [...items.value, ...res.items]
    total.value = res.total
  } catch (e) {
    console.warn(`${props.title}加载失败:`, e.message)
  } finally {
    loading.value = false
  }
}

async function loadMore() {
  loadingMore.value = true
  page.value++
  await load()
  loadingMore.value = false
}
</script>

<style scoped>
.my-content-page { max-width: 900px; margin: 0 auto; padding: 24px 0 60px; }
.page-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.page-header h2 { margin: 0; }
.load-more { text-align: center; padding: 20px; }
</style>
