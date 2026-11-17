<template>
  <div class="follow-list-page">
    <el-card>
      <template #header>
        <el-tabs v-model="activeTab">
          <el-tab-pane label="关注" name="following" />
          <el-tab-pane label="粉丝" name="followers" />
        </el-tabs>
      </template>
      <div v-loading="loading">
        <el-empty v-if="!loading && users.length === 0" description="暂无用户" :image-size="80" />
        <div v-for="u in users" :key="u.id" class="user-item">
          <el-avatar :size="44" :src="u.avatar" @click="$router.push(`/community/user/${u.id}`)">
            {{ u.nickname?.charAt(0) }}
          </el-avatar>
          <div class="user-info" @click="$router.push(`/community/user/${u.id}`)">
            <div class="user-name">{{ u.nickname }}</div>
            <div class="user-bio">{{ u.bio || '这个人很懒，没有签名' }}</div>
          </div>
          <el-button size="small" @click="$router.push(`/community/user/${u.id}`)">查看</el-button>
        </div>
        <div class="load-more" v-if="users.length < total">
          <el-button @click="loadMore">加载更多</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getFollowingList, getFollowerList } from '@/api'

const route = useRoute()
const userId = route.params.id
const activeTab = ref(route.query.type === 'followers' ? 'followers' : 'following')
const users = ref([])
const total = ref(0)
const page = ref(1)
const loading = ref(false)

onMounted(() => load())
watch(activeTab, () => { page.value = 1; users.value = []; load() })

async function load() {
  loading.value = true
  try {
    const api = activeTab.value === 'following' ? getFollowingList : getFollowerList
    const res = await api(userId, { page: page.value, page_size: 20 })
    const data = res?.data || {}
    const items = data.users || []
    users.value = page.value === 1 ? items : [...users.value, ...items]
    total.value = data.total || 0
  } catch (e) { console.warn('加载列表失败:', e.message) }
  finally { loading.value = false }
}
function loadMore() { page.value++; load() }
</script>

<style scoped>
.follow-list-page { max-width: 720px; margin: 0 auto; padding: 24px 0 60px; }
.user-item { display: flex; align-items: center; gap: 12px; padding: 14px 0; border-bottom: 1px solid #f0f0f0; cursor: pointer; }
.user-info { flex: 1; }
.user-name { font-weight: 600; }
.user-bio { color: #999; font-size: 13px; margin-top: 2px; }
.load-more { text-align: center; padding: 16px; }
</style>
