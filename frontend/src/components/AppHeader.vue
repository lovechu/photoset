<template>
  <header class="app-header">
    <div class="header-container">
      <!-- Logo -->
      <router-link to="/" class="logo">
        <el-icon :size="28"><Camera /></el-icon>
        <span>{{ siteStore.siteName }}</span>
      </router-link>

      <!-- 导航 -->
      <nav class="nav-links">
        <router-link to="/" class="nav-item">首页</router-link>
        <router-link to="/community" class="nav-item">社区</router-link>
        <router-link to="/explore" class="nav-item">发现</router-link>
        <!-- 分类平铺导航 -->
        <router-link
          v-for="cat in navCategories"
          :key="cat.slug"
          :to="{ path: '/', query: { category: cat.slug } }"
          class="nav-item"
          :class="{ 'router-link-active': route.query.category === cat.slug }"
        >
          {{ getCategoryName(cat.slug) }}
        </router-link>
        <!-- 标签页面入口 -->
        <router-link to="/tags" class="nav-item">标签</router-link>
        <router-link to="/create" v-if="userStore.isCreatorOrAdmin" class="nav-item">
          <el-icon><Plus /></el-icon>
          创建套图
        </router-link>
        <router-link v-if="userStore.isLoggedIn && !userStore.isMember" to="/membership" class="nav-item vip">
          <el-icon><Medal /></el-icon>
          开通会员
        </router-link>
      </nav>

      <!-- 用户状态 -->
      <div class="user-area">
        <template v-if="userStore.isLoggedIn">
          <!-- 通知入口（未读红点） -->
          <el-badge :value="notificationStore.unreadCount" :hidden="!notificationStore.hasUnread" class="header-badge">
            <el-icon class="header-icon" :size="22" @click="router.push('/notifications')"><Bell /></el-icon>
          </el-badge>
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" :src="userStore.avatar || defaultAvatar">
                {{ userStore.user?.nickname?.charAt(0) }}
              </el-avatar>
              <span class="username">{{ userStore.user?.nickname }}</span>
              <el-tag v-if="userStore.user?.role !== 'user'" size="small" type="warning">
                {{ roleLabel }}
              </el-tag>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>个人中心
                </el-dropdown-item>
                <el-dropdown-item command="collections">
                  <el-icon><Folder /></el-icon>我的合集
                </el-dropdown-item>
                <el-dropdown-item command="history">
                  <el-icon><Clock /></el-icon>浏览历史
                </el-dropdown-item>
                <el-dropdown-item command="user-level">
                  <el-icon><Star /></el-icon>我的等级
                </el-dropdown-item>
                <el-dropdown-item command="creator-dashboard" divided>
                  <el-icon><DataAnalysis /></el-icon>创作者中心
                </el-dropdown-item>
                <el-dropdown-item command="membership">
                  <el-icon><Medal /></el-icon>会员中心
                </el-dropdown-item>
                <el-dropdown-item command="orders">
                  <el-icon><Document /></el-icon>我的订单
                </el-dropdown-item>
                <el-dropdown-item command="settings">
                  <el-icon><Setting /></el-icon>账号设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <template v-else>
          <router-link to="/login">
            <el-button type="primary" plain>登录</el-button>
          </router-link>
          <router-link to="/register">
            <el-button type="primary">注册</el-button>
          </router-link>
        </template>
      </div>

      <!-- 移动端菜单按钮 -->
      <el-icon class="mobile-menu" :size="24" @click="mobileMenuVisible = !mobileMenuVisible">
        <Menu v-if="!mobileMenuVisible" />
        <Close v-else />
      </el-icon>
    </div>

    <!-- 移动端菜单 -->
    <transition name="slide-down">
      <div v-if="mobileMenuVisible" class="mobile-menu-panel">
        <router-link to="/" @click="mobileMenuVisible = false">首页</router-link>
        <router-link to="/community" @click="mobileMenuVisible = false">社区</router-link>
        <router-link to="/explore" @click="mobileMenuVisible = false">发现</router-link>
        <!-- 移动端分类导航 -->
        <router-link
          v-for="cat in navCategories"
          :key="cat.slug"
          :to="{ path: '/', query: { category: cat.slug } }"
          @click="mobileMenuVisible = false"
        >
          {{ getCategoryName(cat.slug) }}
        </router-link>
        <router-link to="/tags" @click="mobileMenuVisible = false">标签</router-link>
        <router-link v-if="userStore.isCreatorOrAdmin" to="/create" @click="mobileMenuVisible = false">
          创建套图
        </router-link>
        <template v-if="userStore.isLoggedIn">
          <a @click="router.push('/profile'); mobileMenuVisible = false">个人中心</a>
          <a @click="router.push('/notifications'); mobileMenuVisible = false">通知中心</a>
          <a @click="router.push('/collections'); mobileMenuVisible = false">我的合集</a>
          <a @click="router.push('/history'); mobileMenuVisible = false">浏览历史</a>
          <a @click="router.push('/user-level'); mobileMenuVisible = false">我的等级</a>
          <a @click="router.push('/creator-dashboard'); mobileMenuVisible = false">创作者中心</a>
          <a v-if="!userStore.isMember" @click="router.push('/membership'); mobileMenuVisible = false">开通会员</a>
          <a @click="router.push('/orders'); mobileMenuVisible = false">我的订单</a>
          <a @click="handleCommand('logout'); mobileMenuVisible = false">退出登录</a>
        </template>
        <template v-else>
          <router-link to="/login" @click="mobileMenuVisible = false">登录</router-link>
          <router-link to="/register" @click="mobileMenuVisible = false">注册</router-link>
        </template>
      </div>
    </transition>
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useSiteStore } from '@/stores/site'
import { useNotificationStore } from '@/stores/notification'
import {
  Camera, Plus, User, SwitchButton, Menu, Close, Medal, Document,
  Bell, Folder, Clock, Star, DataAnalysis, Setting
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getSettings, getCategories } from '@/api'
import {
  connectWebSocket, disconnectWebSocket, onMessage, getWebSocketState
} from '@/utils/websocket'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const siteStore = useSiteStore()
const notificationStore = useNotificationStore()
let unsubscribes = []
const mobileMenuVisible = ref(false)
const defaultAvatar = ''
const navCategories = ref([])
const allCategories = ref([])

onMounted(() => {
  siteStore.fetchSettings()
  loadNavMenu()
  loadAllCategories()
  if (userStore.isLoggedIn) {
    initRealtime()
  }
})

onUnmounted(() => {
  unsubscribes.forEach(fn => fn())
  unsubscribes = []
})

// 登录态变化时连接/断开 WebSocket
watch(() => userStore.isLoggedIn, (loggedIn) => {
  if (loggedIn) {
    initRealtime()
  } else {
    disconnectWebSocket()
    unsubscribes.forEach(fn => fn())
    unsubscribes = []
  }
})

function initRealtime() {
  if (getWebSocketState() === 'open' || getWebSocketState() === 'connecting') return
  notificationStore.fetchUnreadCount()

  unsubscribes.push(onMessage('notification', (payload) => {
    notificationStore.pushNewNotification(payload?.notification || payload || {})
    ElMessage.info(payload?.title || '收到一条新通知')
  }))
  unsubscribes.push(onMessage('unread_count', (payload) => {
    const count = payload?.unread_count ?? payload?.count
    if (typeof count === 'number') {
      notificationStore.unreadCount = count
    }
  }))
  unsubscribes.push(onMessage('message', () => {
    // 私信消息由聊天页自行订阅处理，这里仅可提示
  }))
  connectWebSocket()
}

async function loadNavMenu() {
  try {
    const data = await getSettings()
    if (data?.nav_menu) {
      try {
        const menuItems = JSON.parse(data.nav_menu)
        navCategories.value = menuItems.filter(item => item.slug)
      } catch (parseError) {
        console.error('解析导航菜单JSON失败:', parseError)
        navCategories.value = []
      }
    } else {
      navCategories.value = []
    }
  } catch (e) {
    console.error('加载导航菜单失败', e)
    navCategories.value = []
  }
}

async function loadAllCategories() {
  try {
    const res = await getCategories()
    // 处理响应数据，兼容两种格式
    if (res?.data?.list) {
      allCategories.value = res.data.list
    } else if (Array.isArray(res?.data)) {
      allCategories.value = res.data
    } else if (Array.isArray(res)) {
      allCategories.value = res
    }
  } catch (e) {
    console.error('加载分类列表失败', e)
    allCategories.value = []
  }
}

function getCategoryName(slug) {
  const cat = allCategories.value.find(c => c.slug === slug)
  return cat ? cat.name : slug || '未知分类'
}

const roleLabel = computed(() => {
  const labels = {
    admin: '管理员',
    creator: '创作者',
    member: '会员'
  }
  return labels[userStore.user?.role] || ''
})

const handleCommand = (command) => {
  switch (command) {
    case 'logout':
      userStore.logout()
      ElMessage.success('已退出登录')
      router.push('/')
      break
    case 'profile':
      router.push('/profile')
      break
    case 'membership':
      router.push('/membership')
      break
    case 'orders':
      router.push('/orders')
      break
    case 'collections':
      router.push('/collections')
      break
    case 'history':
      router.push('/history')
      break
    case 'user-level':
      router.push('/user-level')
      break
    case 'creator-dashboard':
      router.push('/creator-dashboard')
      break
    case 'settings':
      router.push('/account-security')
      break
  }
}
</script>

<style scoped>
.app-header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  gap: 32px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  color: #333;
  font-size: 20px;
  font-weight: 600;
}

.logo:hover {
  color: #409eff;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 24px;
}

.nav-item {
  text-decoration: none;
  color: #666;
  font-size: 15px;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: color 0.2s;
}

.nav-item:hover,
.nav-item.router-link-active {
  color: #409eff;
}

.nav-dropdown {
  cursor: pointer;
  display: flex;
  align-items: center;
}

.user-area {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.header-badge {
  display: flex;
  align-items: center;
}

.header-icon {
  cursor: pointer;
  color: #606266;
  transition: color 0.2s;
}

.header-icon:hover {
  color: #409eff;
}

.mobile-menu {
  display: none;
  cursor: pointer;
  margin-left: auto;
}

.mobile-menu-panel {
  background: #fff;
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  border-top: 1px solid #eee;
}

.mobile-menu-panel a {
  text-decoration: none;
  color: #666;
  font-size: 16px;
  padding: 8px 0;
}

.mobile-menu-panel a:active {
  color: #409eff;
}

/* 动画 */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 768px) {
  .nav-links,
  .user-area {
    display: none;
  }

  .mobile-menu {
    display: flex;
  }

  .header-container {
    gap: 16px;
  }
}
</style>
