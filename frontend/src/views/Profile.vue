<template>
  <div class="profile-page">
    <div class="page-header">
      <h1>个人中心</h1>
    </div>

    <!-- 用户信息卡片 -->
    <div class="user-card">
      <el-avatar :size="64">{{ userStore.user?.nickname?.charAt(0) }}</el-avatar>
      <div class="user-info">
        <h2>{{ userStore.user?.nickname }}</h2>
        <el-tag size="small" :type="roleTagType">{{ roleLabel }}</el-tag>
        <div v-if="userStore.membershipExpires" class="membership-info">
          <el-tag type="success" v-if="userStore.isMembershipActive">
            会员有效至 {{ formatDate(userStore.membershipExpires) }}
          </el-tag>
          <el-tag type="info" v-else>
            会员已过期
          </el-tag>
        </div>
        <p>{{ userStore.user?.email }}</p>
        <p class="join-date">注册时间：{{ formatDate(userStore.user?.created_at) }}</p>
      </div>
    </div>

    <!-- Tabs -->
    <el-tabs v-model="activeTab" class="profile-tabs">
      <el-tab-pane label="个人资料" name="info">
        <div class="info-card">
          <div class="info-item">
            <span class="label">昵称</span>
            <span class="value">{{ userStore.user?.nickname }}</span>
            <el-button size="small" text disabled>编辑</el-button>
          </div>
          <div class="info-item">
            <span class="label">邮箱</span>
            <span class="value">{{ userStore.user?.email }}</span>
          </div>
          <div class="info-item">
            <span class="label">角色</span>
            <el-tag :type="roleTagType">{{ roleLabel }}</el-tag>
          </div>
        </div>
      </el-tab-pane>

      <el-tab-pane label="我的收藏" name="favorites">
        <div class="tab-action">
          <el-button type="primary" @click="$router.push('/favorites')">
            查看全部收藏
          </el-button>
        </div>
      </el-tab-pane>

      <el-tab-pane label="我的订单" name="orders">
        <div class="tab-action">
          <el-button type="primary" @click="$router.push('/orders')">
            查看全部订单
          </el-button>
        </div>
      </el-tab-pane>

      <el-tab-pane v-if="userStore.isCreatorOrAdmin" label="我的发布" name="published">
        <div class="photoset-grid" v-loading="loading">
          <div v-if="!loading && photosets.length === 0" class="empty-state">
            <el-empty description="还没有发布套图">
              <el-button type="primary" @click="$router.push('/create')">去发布</el-button>
            </el-empty>
          </div>
          <PhotosetCard v-for="item in photosets" :key="item.id" :data="item" @deleted="handleDeleted" />
        </div>
        <div v-if="total > 0" class="pagination-wrapper">
          <el-pagination
            v-model:current-page="currentPage"
            :total="total"
            :page-size="12"
            layout="prev, pager, next"
            @current-change="loadMyPhotosets"
          />
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getPhotosetList } from '@/api'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import PhotosetCard from '@/components/PhotosetCard.vue'

const userStore = useUserStore()
const activeTab = ref('info')
const loading = ref(false)
const photosets = ref([])
const total = ref(0)
const currentPage = ref(1)

const roleLabel = computed(() => {
  const labels = { admin: '管理员', creator: '创作者', member: '会员', user: '普通用户' }
  return labels[userStore.user?.role] || '普通用户'
})

const roleTagType = computed(() => {
  const types = { admin: 'danger', creator: 'warning', member: 'success', user: 'info' }
  return types[userStore.user?.role] || 'info'
})

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

const loadMyPhotosets = async () => {
  loading.value = true
  try {
    const res = await getPhotosetList({ mine: true, page: currentPage.value, page_size: 12 })
    photosets.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error('加载我的发布失败', e)
  } finally {
    loading.value = false
  }
}

// 处理套图删除
const handleDeleted = (deletedId) => {
  photosets.value = photosets.value.filter(item => item.id !== deletedId)
  total.value -= 1

  if (photosets.value.length === 0 && currentPage.value > 1) {
    currentPage.value -= 1
    loadMyPhotosets()
  }

  ElMessage.success('删除成功')
}

onMounted(() => {
  if (userStore.isCreatorOrAdmin) {
    loadMyPhotosets()
  }
})
</script>

<style scoped>
.profile-page { padding: 0 0 60px; max-width: 1000px; margin: 0 auto; }
.page-header { text-align: center; padding: 40px 0; }
.page-header h1 { font-size: 28px; font-weight: 600; color: #1a1a1a; }
.user-card {
  background: #fff; border-radius: 12px; padding: 24px;
  display: flex; align-items: center; gap: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04); margin-bottom: 24px;
}
.user-info h2 { font-size: 20px; margin: 4px 0; }
.user-info p { color: #666; font-size: 14px; margin: 4px 0; }
.join-date { color: #999 !important; }
.membership-info { margin: 6px 0; }
.profile-tabs { background: #fff; border-radius: 12px; padding: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.04); }
.info-card { max-width: 500px; }
.info-item { display: flex; align-items: center; padding: 16px 0; border-bottom: 1px solid #f0f0f0; }
.info-item .label { width: 80px; color: #999; flex-shrink: 0; }
.info-item .value { flex: 1; }
.tab-action { padding: 40px 0; text-align: center; }
.photoset-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 20px; }
.empty-state { grid-column: 1 / -1; padding: 60px 0; }
.pagination-wrapper { display: flex; justify-content: center; margin-top: 24px; }
</style>
