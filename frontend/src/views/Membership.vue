<template>
  <div class="membership-page">
    <div class="page-header">
      <h1>开通会员</h1>
      <p class="subtitle">成为会员，解锁全部付费套图</p>
    </div>

    <!-- 会员状态展示 -->
    <div v-if="userStore.isLoggedIn" class="membership-status">
      <div v-if="userStore.isMembershipActive" class="status-card active">
        <el-icon :size="32"><Crown /></el-icon>
        <div class="status-info">
          <h3>您已是尊贵的会员</h3>
          <p>会员有效期至 {{ formatDate(userStore.membershipExpires) }}</p>
        </div>
        <el-button type="warning" @click="handleSubscribe">续费会员</el-button>
      </div>
      <div v-else-if="userStore.membershipExpires" class="status-card expired">
        <el-icon :size="32"><Crown /></el-icon>
        <div class="status-info">
          <h3>会员已过期</h3>
          <p>于 {{ formatDate(userStore.membershipExpires) }} 到期</p>
        </div>
        <el-button type="primary" @click="handleSubscribe">重新开通</el-button>
      </div>
    </div>

    <!-- 套餐列表 -->
    <div v-loading="loading" class="plans-grid">
      <div v-if="!loading && plans.length === 0" class="empty-state">
        <el-empty description="暂无可用套餐" />
      </div>
      <div
        v-for="plan in plans"
        :key="plan.id"
        class="plan-card"
        :class="{ recommended: plan.recommended }"
        @click="handleSelectPlan(plan)"
      >
        <div v-if="plan.recommended" class="recommended-badge">推荐</div>
        <h3 class="plan-name">{{ plan.name }}</h3>
        <div class="plan-price">
          <span class="currency">¥</span>
          <span class="amount">{{ plan.price }}</span>
        </div>
        <p class="plan-duration">{{ plan.duration_days }} 天</p>
        <p class="plan-description">{{ plan.description }}</p>
        <el-button
          type="primary"
          :disabled="!userStore.isLoggedIn"
          class="plan-btn"
        >
          {{ userStore.isLoggedIn ? '立即开通' : '登录后开通' }}
        </el-button>
      </div>
    </div>

    <!-- 未登录提示 -->
    <div v-if="!userStore.isLoggedIn" class="login-prompt">
      <el-alert type="info" :closable="false" show-icon>
        <template #title>
          <span>登录后即可开通会员</span>
          <el-button type="primary" size="small" @click="$router.push('/login')">
            去登录
          </el-button>
        </template>
      </el-alert>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getMemberships, createOrder, payOrder } from '@/api'
import { Crown } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const plans = ref([])

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const loadPlans = async () => {
  loading.value = true
  try {
    const res = await getMemberships()
    // 标记推荐套餐
    plans.value = (res.data.list || res.data || []).map((plan, index) => ({
      ...plan,
      recommended: index === 0 // 默认第一个为推荐
    }))
  } catch (e) {
    console.error('加载套餐列表失败', e)
    ElMessage.error('加载套餐列表失败')
  } finally {
    loading.value = false
  }
}

const handleSelectPlan = async (plan) => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
    return
  }

  try {
    // 创建订单
    const res = await createOrder({
      type: 'membership',
      membership_id: plan.id
    })
    const order = res.data.order

    // 确认支付
    await ElMessageBox.confirm(
      `确认支付 ¥${order.amount} 开通 ${plan.name}？`,
      '确认支付',
      {
        confirmButtonText: '确认支付',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // 模拟支付
    const payRes = await payOrder(order.id)
    const { token, user } = payRes.data

    // 更新状态
    userStore.afterPayment(token, user)

    ElMessage.success('支付成功，欢迎成为会员！')

    // 刷新页面以更新状态
    setTimeout(() => {
      window.location.reload()
    }, 1500)
  } catch (e) {
    if (e !== 'cancel') {
      console.error('支付失败', e)
      ElMessage.error('支付失败，请重试')
    }
  }
}

const handleSubscribe = () => {
  // 如果已有套餐，默认选择第一个
  if (plans.value.length > 0) {
    handleSelectPlan(plans.value[0])
  }
}

onMounted(() => {
  loadPlans()
})
</script>

<style scoped>
.membership-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px 60px;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
}

.page-header h1 {
  font-size: 32px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.subtitle {
  color: #666;
  font-size: 16px;
}

.membership-status {
  margin-bottom: 40px;
}

.status-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.status-card.active {
  background: linear-gradient(135deg, #fff9e6 0%, #fff 100%);
  border: 1px solid #ffe7b3;
}

.status-card.active .el-icon {
  color: #ffa500;
}

.status-card.expired {
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
}

.status-info {
  flex: 1;
}

.status-info h3 {
  font-size: 18px;
  margin: 0 0 4px;
}

.status-info p {
  color: #666;
  margin: 0;
  font-size: 14px;
}

.plans-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
}

.plan-card {
  background: #fff;
  border-radius: 16px;
  padding: 32px 24px;
  text-align: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  border: 2px solid transparent;
}

.plan-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.plan-card.recommended {
  border-color: #409eff;
}

.recommended-badge {
  position: absolute;
  top: -12px;
  left: 50%;
  transform: translateX(-50%);
  background: #409eff;
  color: #fff;
  padding: 4px 16px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.plan-name {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 16px;
}

.plan-price {
  color: #f56c6c;
  margin-bottom: 8px;
}

.currency {
  font-size: 20px;
  vertical-align: top;
}

.amount {
  font-size: 48px;
  font-weight: 700;
  line-height: 1;
}

.plan-duration {
  color: #999;
  font-size: 14px;
  margin: 0 0 16px;
}

.plan-description {
  color: #666;
  font-size: 14px;
  line-height: 1.6;
  margin: 0 0 24px;
  min-height: 44px;
}

.plan-btn {
  width: 100%;
}

.empty-state {
  grid-column: 1 / -1;
  padding: 60px 0;
}

.login-prompt {
  margin-top: 32px;
}

.login-prompt .el-alert {
  border-radius: 12px;
}

.login-prompt :deep(.el-alert__title) {
  display: flex;
  align-items: center;
  gap: 16px;
}

@media (max-width: 768px) {
  .membership-page {
    padding: 24px 16px 40px;
  }

  .page-header h1 {
    font-size: 24px;
  }

  .plans-grid {
    grid-template-columns: 1fr;
  }

  .status-card {
    flex-direction: column;
    text-align: center;
  }
}
</style>
