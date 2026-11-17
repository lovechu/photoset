<template>
  <div class="community-search-page">
    <div class="search-bar">
      <el-input
        v-model="keyword"
        size="large"
        placeholder="搜索帖子、用户..."
        :prefix-icon="Search"
        clearable
        @keyup.enter="doSearch"
      >
        <template #append>
          <el-button :icon="Search" @click="doSearch">搜索</el-button>
        </template>
      </el-input>
    </div>

    <div v-if="searched" class="search-results">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="帖子" name="posts">
          <div v-loading="loadingPosts">
            <el-empty v-if="!loadingPosts && posts.length === 0" description="没有找到相关帖子" :image-size="80" />
            <PostCard v-for="post in posts" :key="post.id" :post="post" />
            <div class="load-more" v-if="posts.length < postsTotal">
              <el-button @click="loadMorePosts">加载更多</el-button>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="用户" name="users">
          <div v-loading="loadingUsers">
            <el-empty v-if="!loadingUsers && users.length === 0" description="没有找到相关用户" :image-size="80" />
            <div v-for="u in users" :key="u.id" class="user-item" @click="$router.push(`/community/user/${u.id}`)">
              <el-avatar :size="44" :src="u.avatar">{{ u.nickname?.charAt(0) }}</el-avatar>
              <div class="user-info">
                <div class="user-name">{{ u.nickname }}</div>
                <div class="user-bio">{{ u.bio || '这个人很懒，没有签名' }}</div>
              </div>
              <el-button size="small" @click.stop="$router.push(`/community/user/${u.id}`)">查看</el-button>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
    <el-empty v-else description="输入关键词开始搜索" :image-size="100" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Search } from '@element-plus/icons-vue'
import PostCard from '@/components/PostCard.vue'
import { searchPosts, searchUsersForMention } from '@/api'

const keyword = ref('')
const searched = ref(false)
const activeTab = ref('posts')
const posts = ref([])
const users = ref([])
const postsTotal = ref(0)
const postPage = ref(1)
const loadingPosts = ref(false)
const loadingUsers = ref(false)

function doSearch() {
  if (!keyword.value.trim()) return
  searched.value = true
  postPage.value = 1
  posts.value = []
  loadPosts()
  loadUsers()
}

async function loadPosts() {
  loadingPosts.value = true
  try {
    const res = await searchPosts({ keyword: keyword.value, page: postPage.value, page_size: 20 })
    const data = res?.data || {}
    const items = data.posts || []
    posts.value = postPage.value === 1 ? items : [...posts.value, ...items]
    postsTotal.value = data.pagination?.total || 0
  } catch (e) { console.warn('搜索帖子失败:', e.message) }
  finally { loadingPosts.value = false }
}

function loadMorePosts() { postPage.value++; loadPosts() }

async function loadUsers() {
  loadingUsers.value = true
  try {
    const res = await searchUsersForMention({ keyword: keyword.value, limit: 20 })
    const data = res?.data || {}
    users.value = data.users || (Array.isArray(data) ? data : [])
  } catch (e) { console.warn('搜索用户失败:', e.message) }
  finally { loadingUsers.value = false }
}
</script>

<style scoped>
.community-search-page {
  max-width: 1000px; margin: 0 auto; padding: 24px 20px 60px;
  color: #33373d; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}
.search-bar { margin-bottom: 20px; }
.user-item { display: flex; align-items: center; gap: 12px; padding: 14px 16px; background: #fff; border: 1px solid #e6e8eb; border-radius: 5px; margin-bottom: 10px; cursor: pointer; }
.user-item:hover { border-color: #cdd5e0; background: #f8fafc; }
.user-info { flex: 1; }
.user-name { font-weight: 600; color: #2b3038; }
.user-bio { color: #9aa1ac; font-size: 13px; margin-top: 2px; }
.load-more { text-align: center; padding: 20px; }
</style>
