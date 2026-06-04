<template>
  <el-container class="admin-layout">
    <el-aside width="220px" class="admin-aside">
      <div class="logo">
        <h2>PhotoSet</h2>
        <p>管理后台</p>
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <span>数据看板</span>
        </el-menu-item>
        <el-menu-item index="/review">
          <el-icon><PictureFilled /></el-icon>
          <span>内容审核</span>
        </el-menu-item>
        <el-menu-item index="/users">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item index="/orders">
          <el-icon><ShoppingBag /></el-icon>
          <span>订单管理</span>
        </el-menu-item>
        <el-menu-item index="/memberships">
          <el-icon><Coin /></el-icon>
          <span>会员管理</span>
        </el-menu-item>
        <el-menu-item index="/tags">
          <el-icon><PriceTag /></el-icon>
          <span>标签管理</span>
        </el-menu-item>
        <el-menu-item index="/categories">
          <el-icon><FolderOpened /></el-icon>
          <span>分类管理</span>
        </el-menu-item>
        <el-menu-item index="/pages">
          <el-icon><Document /></el-icon>
          <span>页面管理</span>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>站点设置</span>
        </el-menu-item>
        <el-menu-item index="/developer">
          <el-icon><Connection /></el-icon>
          <span>开发者中心</span>
        </el-menu-item>
        <el-menu-item index="/oauth">
          <el-icon><Key /></el-icon>
          <span>OAuth应用管理</span>
        </el-menu-item>
        <el-menu-item index="/logs">
          <el-icon><Notebook /></el-icon>
          <span>操作日志</span>
        </el-menu-item>
        <el-sub-menu index="/system">
          <template #title>
            <el-icon><Monitor /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/system/monitor">
            <span>系统监控</span>
          </el-menu-item>
          <el-menu-item index="/system/performance">
            <span>性能分析</span>
          </el-menu-item>
          <el-menu-item index="/system/backup">
            <span>数据备份</span>
          </el-menu-item>
          <el-menu-item index="/system/trash">
            <span>回收站管理</span>
          </el-menu-item>
          <el-menu-item index="/system/ip-geo">
            <span>IP地理位置</span>
          </el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="/community">
          <template #title>
            <el-icon><ChatDotRound /></el-icon>
            <span>社区管理</span>
          </template>
          <el-menu-item index="/community/posts">
            <span>帖子管理</span>
          </el-menu-item>
          <el-menu-item index="/community/replies">
            <span>回帖管理</span>
          </el-menu-item>
          <el-menu-item index="/community/keywords">
            <span>敏感词管理</span>
          </el-menu-item>
          <el-menu-item index="/community/categories">
            <span>分类管理</span>
          </el-menu-item>
          <el-menu-item index="/community/reports">
            <span>举报处理</span>
          </el-menu-item>
          <el-menu-item index="/community/users">
            <span>用户积分管理</span>
          </el-menu-item>
          <el-menu-item index="/community/levels">
            <span>等级管理</span>
          </el-menu-item>
          <el-menu-item index="/community/achievements">
            <span>成就管理</span>
          </el-menu-item>
          <el-menu-item index="/community/points-mall">
            <span>积分商城</span>
          </el-menu-item>
          <el-menu-item index="/community/notifications">
            <span>通知管理</span>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="admin-header">
        <div class="header-left">
          <span class="page-title">{{ currentPageTitle }}</span>
        </div>
        <div class="header-right">

          <el-dropdown trigger="click">
            <span class="user-info">
              <el-icon><UserFilled /></el-icon>
              {{ adminStore.user?.nickname || '管理员' }}
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="admin-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { Connection, FolderOpened, ChatDotRound, Monitor, Key, Coin } from '@element-plus/icons-vue'


const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()


const activeMenu = computed(() => route.path)

const currentPageTitle = computed(() => {
  return route.meta?.title || '管理后台'
})

function handleLogout() {
  adminStore.logout()
  router.push('/login')
}


</script>

<style scoped>
.admin-layout {
  height: 100vh;
  overflow: hidden;
}

.admin-aside {
  background-color: #304156;
  overflow-y: auto;
  overflow-x: hidden;
}

.logo {
  height: 60px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 2px;
}

.logo p {
  margin: 2px 0 0;
  font-size: 12px;
  color: #bfcbd9;
}

.el-menu {
  border-right: none;
}

.admin-header {
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  padding: 0 24px;
  height: 60px;
}

.header-left .page-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.header-right .user-info {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  color: #606266;
  font-size: 14px;
}

.admin-main {
  background: #f0f2f5;
  overflow-y: auto;
}
</style>
