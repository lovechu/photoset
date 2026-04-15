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
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" :src="defaultAvatar">
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
                <el-dropdown-item command="membership" v-if="userStore.isLoggedIn">
                  <el-icon><Medal /></el-icon>会员中心
                </el-dropdown-item>
                <el-dropdown-item command="orders" v-if="userStore.isLoggedIn">
                  <el-icon><Document /></el-icon>我的订单
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
        <router-link v-if="userStore.isCreatorOrAdmin" to="/create" @click="mobileMenuVisible = false">
          创建套图
        </router-link>
        <template v-if="userStore.isLoggedIn">
          <a @click="router.push('/profile'); mobileMenuVisible = false">个人中心</a>
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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useSiteStore } from '@/stores/site'
import { Camera, Plus, User, SwitchButton, Menu, Close, Medal, Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const siteStore = useSiteStore()
const mobileMenuVisible = ref(false)
const defaultAvatar = ''

onMounted(() => {
  siteStore.fetchSettings()
})

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
