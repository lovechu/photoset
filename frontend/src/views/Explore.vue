<template>
  <div class="explore-page">
    <div class="explore-hero">
      <h2>探索发现</h2>
      <p>发现你可能感兴趣的精彩套图</p>
    </div>

    <el-tabs v-model="activeTab" @tab-change="reload">
      <el-tab-pane label="推荐" name="feed" />
      <el-tab-pane label="热门" name="hot" />
    </el-tabs>

    <div v-loading="loading">
      <el-empty v-if="!loading && items.length === 0" description="暂无推荐内容" :image-size="100" />
      <div v-else class="explore-grid">
        <el-card
          v-for="item in items"
          :key="item.id"
          class="explore-card"
          shadow="hover"
          @click="goDetail(item)"
        >
          <div class="cover-wrap">
            <el-image :src="item.cover_image" fit="cover" class="cover" :lazy="true" />
            <span v-if="item.is_vip" class="vip-tag">VIP</span>
            <span v-if="item.source === 'hot'" class="hot-tag">热门</span>
            <span class="photo-count">{{ item.photo_count || 0 }} 图</span>
          </div>
          <div class="card-body">
            <div class="card-title">{{ item.title || '未命名套图' }}</div>
            <div class="card-meta">
              <span>{{ item.view_count || 0 }} 浏览</span>
              <span>{{ item.fav_count || 0 }} 收藏</span>
              <span class="price" :class="{ free: !item.price }">
                {{ item.price ? `¥${Number(item.price).toFixed(2)}` : '免费' }}
              </span>
            </div>
          </div>
        </el-card>
      </div>

      <div class="load-more" v-if="hasMore">
        <el-button @click="loadMore" :loading="loadingMore">加载更多</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getExploreFeed, getHotPhotosets } from '@/api'

const router = useRouter()
const activeTab = ref('feed')
const items = ref([])
const page = ref(1)
const hasMore = ref(false)
const loading = ref(false)
const loadingMore = ref(false)

onMounted(() => load())

function reload() {
  page.value = 1
  items.value = []
  load()
}

async function load() {
  loading.value = true
  try {
    if (activeTab.value === 'hot') {
      const res = await getHotPhotosets({ page: page.value, page_size: 20 })
      const data = res?.data || {}
      items.value = data.items || []
      hasMore.value = false // 热门接口当前为占位实现
    } else {
      const res = await getExploreFeed({ page: page.value, page_size: 20 })
      const data = res?.data || {}
      const list = data.items || []
      items.value = page.value === 1 ? list : [...items.value, ...list]
      hasMore.value = !!data.has_more
    }
  } catch (e) {
    console.warn('加载探索内容失败:', e.message)
  } finally {
    loading.value = false
  }
}

function loadMore() {
  page.value++
  load()
}

function goDetail(item) {
  router.push(`/photoset/${item.id}`)
}
</script>

<style scoped>
.explore-page { max-width: 1100px; margin: 0 auto; padding: 24px 0 60px; }
.explore-hero { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); border-radius: 16px; padding: 32px; color: #fff; margin-bottom: 20px; }
.explore-hero h2 { margin: 0 0 6px; font-size: 28px; }
.explore-hero p { margin: 0; opacity: 0.85; }
.explore-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 16px; }
.explore-card { overflow: hidden; cursor: pointer; }
.cover-wrap { position: relative; height: 160px; }
.cover { width: 100%; height: 100%; }
.vip-tag, .hot-tag { position: absolute; top: 8px; left: 8px; padding: 2px 8px; border-radius: 4px; font-size: 11px; color: #fff; background: #f5a623; }
.hot-tag { left: auto; right: 8px; background: #f56c6c; }
.photo-count { position: absolute; bottom: 8px; right: 8px; padding: 2px 8px; border-radius: 4px; font-size: 11px; color: #fff; background: rgba(0,0,0,0.5); }
.card-title { font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.card-meta { display: flex; gap: 12px; color: #999; font-size: 12px; margin-top: 8px; }
.card-meta .price { margin-left: auto; color: #f56c6c; font-weight: 600; }
.card-meta .price.free { color: #67c23a; }
.load-more { text-align: center; padding: 24px; }
</style>
