<template>
  <div class="user-profile-page" v-loading="loading">
    <template v-if="profile">
      <div class="profile-card">
        <div class="cover-banner" :style="coverStyle">
          <div class="cover-overlay"></div>
        </div>
        <div class="profile-body">
          <div class="profile-head">
            <el-avatar :size="84" :src="profile.avatar" class="profile-avatar">
              {{ profile.nickname?.charAt(0) || 'U' }}
            </el-avatar>
            <div class="profile-info">
              <div class="name-row">
                <h2>{{ profile.nickname }}</h2>
                <el-tag size="small" type="warning" effect="dark" class="level-tag">Lv.{{ profile.level || 1 }}</el-tag>
                <el-tag v-if="isSelf" size="small" type="info" effect="plain">我</el-tag>
              </div>
              <p class="bio">{{ profile.bio || '这个人很懒，没有签名' }}</p>
              <div class="profile-meta">
                <span v-if="profile.ip_location" class="meta-item">IP {{ profile.ip_location }}</span>
                <span class="meta-item">{{ formatDate(profile.created_at) }} 加入</span>
              </div>
            </div>
            <div class="profile-actions">
              <el-button
                v-if="isSelf"
                type="primary"
                plain
                @click="$router.push('/edit-profile')"
              >编辑资料</el-button>
              <el-button
                v-else-if="userStore.isLoggedIn"
                type="primary"
                @click="toggleFollow"
              >{{ isFollowing ? '已关注' : '+ 关注' }}</el-button>
            </div>
          </div>

          <div class="profile-stats">
            <div @click="switchTab('posts')"><strong>{{ profile.post_count || 0 }}</strong><span>帖子</span></div>
            <div><strong>{{ profile.reply_count || 0 }}</strong><span>回复</span></div>
            <div @click="$router.push(`/community/follow-list/${profile.id}?type=following`)"><strong>{{ profile.following_count || 0 }}</strong><span>关注</span></div>
            <div @click="$router.push(`/community/follow-list/${profile.id}?type=followers`)"><strong>{{ profile.follower_count || 0 }}</strong><span>粉丝</span></div>
            <div><strong>{{ profile.like_count || 0 }}</strong><span>获赞</span></div>
          </div>

          <div v-if="unlockedAchievements.length" class="achievement-strip">
            <span class="strip-label">成就</span>
            <el-tooltip
              v-for="a in unlockedAchievements.slice(0, 8)"
              :key="a.id"
              :content="a.description || a.title || a.name"
              placement="top"
            >
              <span class="badge">{{ a.icon || '🏅' }} {{ a.title || a.name }}</span>
            </el-tooltip>
            <span v-if="unlockedAchievements.length > 8" class="badge-more">+{{ unlockedAchievements.length - 8 }}</span>
          </div>
        </div>
      </div>

      <div class="profile-tabs-card">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="TA 的帖子" name="posts">
            <el-empty v-if="!posts.length && !postsLoading" description="暂无帖子" :image-size="60" />
            <PostCard v-for="post in posts" :key="post.id" :post="post" />
            <div class="load-more" v-if="posts.length < postsTotal">
              <el-button :loading="postsLoading" @click="loadMorePosts">加载更多</el-button>
            </div>
          </el-tab-pane>
          <el-tab-pane :label="`成就 (${unlockedAchievements.length})`" name="achievements">
            <div class="achievement-grid" v-loading="achievementsLoading">
              <el-empty v-if="!achievements.length && !achievementsLoading" description="暂无成就数据" :image-size="60" />
              <div
                v-for="a in achievements"
                :key="a.id"
                :class="['achievement-item', { locked: !a.is_unlocked }]"
              >
                <div class="achievement-icon">{{ a.icon || '🏅' }}</div>
                <div class="achievement-name">{{ a.title || a.name }}</div>
                <div class="achievement-desc">{{ a.description }}</div>
                <div class="achievement-status">
                  <el-tag v-if="a.is_unlocked" size="small" type="success">已解锁</el-tag>
                  <el-tag v-else size="small" type="info">未解锁</el-tag>
                </div>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </template>
    <el-result v-else-if="!loading" icon="warning" title="用户不存在">
      <template #extra><el-button type="primary" @click="$router.back()">返回</el-button></template>
    </el-result>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import PostCard from '@/components/PostCard.vue'
import { useUserStore } from '@/stores/user'
import {
  getCommunityUserProfile, getCommunityUserPosts,
  getUserAchievementsById,
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
const achievements = ref([])
const achievementsLoading = ref(false)
const isFollowing = ref(false)
const activeTab = ref('posts')

const isSelf = computed(() => userStore.isLoggedIn && Number(userId) === userStore.user?.id)
const unlockedAchievements = computed(() => achievements.value.filter(a => a.is_unlocked))
const coverStyle = computed(() => {
  const url = profile.value?.cover_image
  return url
    ? { backgroundImage: `url(${url})` }
    : { background: 'linear-gradient(120deg, #4f6ef7 0%, #9b5de5 50%, #e07a5f 100%)' }
})

onMounted(() => {
  loadProfile()
  loadPosts()
  loadAchievements()
  if (userStore.isLoggedIn && !isSelf.value) checkFollow()
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

async function loadAchievements() {
  achievementsLoading.value = true
  try {
    const res = await getUserAchievementsById(userId)
    const data = res?.data
    achievements.value = Array.isArray(data) ? data : (data?.achievements || [])
  } catch (e) { achievements.value = [] }
  finally { achievementsLoading.value = false }
}

function switchTab(name) { activeTab.value = name }
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
.user-profile-page { max-width: 960px; margin: 0 auto; padding: 24px 0 60px; }
.profile-card { background: #fff; border-radius: 12px; overflow: hidden; margin-bottom: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.06); }
.cover-banner { height: 160px; background-size: cover; background-position: center; position: relative; }
.cover-overlay { position: absolute; inset: 0; background: linear-gradient(to bottom, rgba(0,0,0,0) 40%, rgba(0,0,0,0.15) 100%); }
.profile-body { padding: 0 28px 24px; }
.profile-head { display: flex; align-items: flex-start; gap: 20px; margin-top: -42px; position: relative; }
.profile-avatar { flex-shrink: 0; border: 4px solid #fff; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.profile-info { flex: 1; padding-top: 44px; min-width: 0; }
.name-row { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
.name-row h2 { margin: 0; font-size: 22px; }
.level-tag { font-weight: 600; }
.bio { color: #666; margin: 8px 0; line-height: 1.5; }
.profile-meta { display: flex; flex-wrap: wrap; gap: 16px; color: #999; font-size: 12px; }
.profile-actions { padding-top: 48px; flex-shrink: 0; }
.profile-stats { display: flex; gap: 40px; margin-top: 24px; padding: 18px 0; border-top: 1px solid #f0f0f0; border-bottom: 1px solid #f0f0f0; }
.profile-stats > div { display: flex; flex-direction: column; cursor: pointer; align-items: center; }
.profile-stats > div:hover strong { color: #4f6ef7; }
.profile-stats strong { font-size: 20px; color: #3d444e; }
.profile-stats span { color: #999; font-size: 12px; margin-top: 2px; }
.achievement-strip { display: flex; flex-wrap: wrap; align-items: center; gap: 8px; margin-top: 18px; }
.strip-label { color: #999; font-size: 12px; margin-right: 4px; }
.badge { background: #fff7e6; border: 1px solid #ffd591; color: #ad6800; padding: 3px 10px; border-radius: 14px; font-size: 12px; cursor: default; }
.badge-more { color: #999; font-size: 12px; }
.profile-tabs-card { background: #fff; border-radius: 12px; padding: 8px 20px 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.load-more { text-align: center; padding: 20px; }
.achievement-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 14px; padding: 8px 0; }
.achievement-item { border: 1px solid #e6e8eb; border-radius: 8px; padding: 16px; text-align: center; background: #fff; }
.achievement-item.locked { opacity: .55; filter: grayscale(.6); }
.achievement-icon { font-size: 30px; margin-bottom: 8px; }
.achievement-name { font-weight: 600; color: #3d444e; margin-bottom: 4px; }
.achievement-desc { font-size: 12px; color: #999; min-height: 32px; line-height: 1.4; }
.achievement-status { margin-top: 8px; }
@media (max-width: 768px) {
  .profile-head { flex-direction: column; gap: 0; }
  .profile-info { padding-top: 12px; }
  .profile-actions { padding-top: 16px; }
  .profile-stats { gap: 20px; justify-content: space-between; }
}
</style>
