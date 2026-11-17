<template>
  <div class="user-profile-page" v-loading="loading">
    <template v-if="profile">
      <div class="profile-card">
        <el-avatar :size="80" :src="profile.avatar" class="profile-avatar">
          {{ profile.nickname?.charAt(0) || 'U' }}
        </el-avatar>
        <div class="profile-info">
          <h2>{{ profile.nickname }}</h2>
          <p class="bio">{{ profile.bio || '这个人很懒，没有签名' }}</p>
          <div class="profile-meta">
            <span v-if="profile.ip_location">📍 {{ profile.ip_location }}</span>
            <span>等级 Lv.{{ profile.level || 1 }}</span>
            <span>{{ formatDate(profile.created_at) }} 加入</span>
          </div>
          <div class="profile-stats">
            <div><strong>{{ profile.post_count || 0 }}</strong><span>帖子</span></div>
            <div @click="$router.push(`/community/follow-list/${profile.id}?type=following`)"><strong>{{ profile.following_count || 0 }}</strong><span>关注</span></div>
            <div @click="$router.push(`/community/follow-list/${profile.id}?type=followers`)"><strong>{{ profile.follower_count || 0 }}</strong><span>粉丝</span></div>
            <div><strong>{{ profile.like_count || 0 }}</strong><span>获赞</span></div>
          </div>
        </div>
        <el-button
          v-if="userStore.isLoggedIn && profile.id !== userStore.user?.id"
          type="primary"
          @click="toggleFollow"
        >{{ isFollowing ? '已关注' : '+ 关注' }}</el-button>
      </div>

      <div class="user-posts">
        <div class="section-head"><h3>TA 的帖子</h3></div>
        <el-empty v-if="!posts.length && !postsLoading" description="暂无帖子" :image-size="60" />
        <PostCard v-for="post in posts" :key="post.id" :post="post" />
        <div class="load-more" v-if="posts.length < postsTotal">
          <el-button @click="loadMorePosts">加载更多</el-button>
        </div>
      </div>
    </template>
    <el-result v-else-if="!loading" icon="warning" title="用户不存在">
      <template #extra><el-button type="primary" @click="$router.back()">返回</el-button></template>
    </el-result>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import PostCard from '@/components/PostCard.vue'
import { useUserStore } from '@/stores/user'
import {
  getCommunityUserProfile, getCommunityUserPosts,
  followUser, unfollowUser, checkFollowing
} from '@/api'

const route = useRoute()
const userStore = useUserStore()
const userId = route.params.id

const loading = ref(false)
const profile = ref(null)
const posts = ref([])
const postsTotal = ref(0)
const postPage = ref(1)
const postsLoading = ref(false)
const isFollowing = ref(false)

onMounted(() => {
  loadProfile()
  loadPosts()
  if (userStore.isLoggedIn && Number(userId) !== userStore.user?.id) checkFollow()
})

async function loadProfile() {
  loading.value = true
  try {
    const res = await getCommunityUserProfile(userId)
    profile.value = res?.data || null
  } catch (e) { console.warn('加载用户资料失败:', e.message) }
  finally { loading.value = false }
}

async function loadPosts() {
  postsLoading.value = true
  try {
    const res = await getCommunityUserPosts(userId, { page: postPage.value, page_size: 20 })
    const data = res?.data || {}
    const items = data.posts || []
    posts.value = postPage.value === 1 ? items : [...posts.value, ...items]
    postsTotal.value = data.pagination?.total || data.total || 0
  } catch (e) { console.warn('加载用户帖子失败:', e.message) }
  finally { postsLoading.value = false }
}

function loadMorePosts() { postPage.value++; loadPosts() }

async function checkFollow() {
  try {
    const res = await checkFollowing(userId)
    isFollowing.value = !!(res?.data?.is_following ?? res?.data?.following)
  } catch (e) { /* 静默 */ }
}

async function toggleFollow() {
  if (!userStore.isLoggedIn) return ElMessage.info('请先登录')
  try {
    if (isFollowing.value) await unfollowUser(userId)
    else await followUser(userId)
    isFollowing.value = !isFollowing.value
    ElMessage.success(isFollowing.value ? '已关注' : '已取消关注')
  } catch (e) { ElMessage.error('操作失败') }
}

function formatDate(str) {
  if (!str) return ''
  return new Date(str).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.user-profile-page { max-width: 900px; margin: 0 auto; padding: 24px 0 60px; }
.profile-card { background: #fff; border-radius: 12px; padding: 28px; display: flex; gap: 20px; align-items: flex-start; margin-bottom: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.profile-avatar { flex-shrink: 0; }
.profile-info { flex: 1; }
.profile-info h2 { margin: 0 0 8px; }
.bio { color: #666; margin: 0 0 10px; }
.profile-meta { display: flex; flex-wrap: wrap; gap: 16px; color: #999; font-size: 12px; margin-bottom: 16px; }
.profile-stats { display: flex; gap: 32px; }
.profile-stats > div { display: flex; flex-direction: column; cursor: pointer; }
.profile-stats strong { font-size: 18px; }
.profile-stats span { color: #999; font-size: 12px; }
.section-head h3 { margin: 0 0 14px; }
.load-more { text-align: center; padding: 20px; }
</style>
