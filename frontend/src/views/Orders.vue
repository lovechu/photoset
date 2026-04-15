<template>
  <div class="orders-page">
    <div class="page-header">
      <h1>我的订单</h1>
    </div>

    <!-- 订单列表 -->
    <div v-loading="loading" class="orders-container">
      <div v-if="!loading && orders.length === 0" class="empty-state">
        <el-empty description="暂无订单">
          <el-button type="primary" @click="$router.push('/')">去逛逛</el-button>
        </el-empty>
      </div>

      <div v-else class="orders-list">
        <div v-for="order in orders" :key="order.id" class="order-card">
          <div class="order-header">
            <span class="order-no">订单号：{{ order.order_no }}</span>
            <el-tag :type="statusType(order.status)" size="small">
              {{ statusLabel(order.status) }}
            </el-tag>
          </div>

          <div class="order-body">
            <div class="order-type">
              <el-icon v-if="order.type === 'membership'"><Medal /></el-icon>
              <el-icon v-else><Picture /></el-icon>
              <span>{{ order.type === 'membership' ? '会员订阅' : '套图购买' }}</span>
            </div>

            <div class="order-detail">
              <span v-if="order.type === 'membership' && order.membership">
                {{ order.membership.name }}
              </span>
              <span v-else-if="order.type === 'single' && order.photoset">
                {{ order.photoset.title }}
              </span>
            </div>

            <div class="order-amount">
              <span class="currency">¥</span>
              <span class="amount">{{ order.amount }}</span>
            </div>
          </div>

          <div class="order-footer">
            <span class="order-time">{{ formatDateTime(order.created_at) }}</span>
            <el-button
              v-if="order.status === 'pending'"
              type="primary"
              size="small"
              @click="handlePay(order)"
            >
              去支付
            </el-button>
          </div>

          <!-- 退款区域 -->
          <div v-if="order.status === 'paid'" class="order-refund">
            <div v-if="canRefund(order)" class="refund-countdown">
              <el-text type="info" size="small">
                可退款：剩余 {{ formatRefundTimeLeft(order) }}
              </el-text>
            </div>
            <el-button
              v-if="canRefund(order)"
              type="danger"
              size="small"
              :icon="RefreshLeft"
              @click="handleRefund(order)"
            >
              申请退款
            </el-button>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="total > 0" class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next, jumper"
          @current-change="loadOrders"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { getOrders, payOrder, refundOrder } from '@/api'
import { Medal, Picture, RefreshLeft } from '@element-plus/icons-vue'

const userStore = useUserStore()
const loading = ref(false)
const orders = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10

const formatDateTime = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const statusLabel = (status) => {
  const labels = {
    pending: '待支付',
    paid: '已支付',
    cancelled: '已取消',
    refunded: '已退款'
  }
  return labels[status] || status
}

const statusType = (status) => {
  const types = {
    pending: 'warning',
    paid: 'success',
    cancelled: 'info',
    refunded: 'danger'
  }
  return types[status] || 'info'
}

const loadOrders = async () => {
  loading.value = true
  try {
    const res = await getOrders({
      page: currentPage.value,
      page_size: pageSize
    })
    orders.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error('加载订单失败', e)
    ElMessage.error('加载订单失败')
  } finally {
    loading.value = false
  }
}

const handlePay = async (order) => {
  try {
    await ElMessageBox.confirm(
      `确认支付 ¥${order.amount}？`,
      '确认支付',
      {
        confirmButtonText: '确认支付',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await payOrder(order.id)
    const { token, user } = res.data

    // 更新状态
    userStore.afterPayment(token, user)

    ElMessage.success('支付成功！')

    // 刷新订单列表
    loadOrders()
  } catch (e) {
    if (e !== 'cancel') {
      console.error('支付失败', e)
      ElMessage.error('支付失败，请重试')
    }
  }
}

// 计算是否可以退款（48小时内）
const canRefund = (order) => {
  if (order.status !== 'paid') return false

  const orderTime = new Date(order.created_at)
  const now = new Date()
  const hoursDiff = (now - orderTime) / (1000 * 60 * 60)

  return hoursDiff <= 48
}

// 格式化退款剩余时间
const formatRefundTimeLeft = (order) => {
  const orderTime = new Date(order.created_at)
  const now = new Date()
  const totalHours = 48
  const usedHours = (now - orderTime) / (1000 * 60 * 60)
  const leftHours = totalHours - usedHours

  if (leftHours < 1) {
    const leftMinutes = Math.floor(leftHours * 60)
    return `${leftMinutes}分钟`
  } else if (leftHours < 24) {
    return `${Math.floor(leftHours)}小时`
  } else {
    const days = Math.floor(leftHours / 24)
    const hours = Math.floor(leftHours % 24)
    return `${days}天${hours}小时`
  }
}

// 处理退款申请
const handleRefund = async (order) => {
  try {
    await ElMessageBox.confirm(
      `确定要申请退款吗？退款金额：¥${order.amount}。退款申请提交后不可撤销。`,
      '确认退款',
      {
        confirmButtonText: '确认退款',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    loading.value = true
    await refundOrder(order.id)

    // 刷新订单列表
    loadOrders()
    ElMessage.success('退款申请已提交')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error?.response?.data?.error || '退款申请失败')
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.orders-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 40px 20px 60px;
}

.page-header {
  text-align: center;
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #1a1a1a;
}

.orders-container {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  min-height: 400px;
}

.empty-state {
  padding: 60px 0;
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.order-card {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 16px;
  transition: box-shadow 0.2s;
}

.order-card:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.order-no {
  color: #999;
  font-size: 13px;
}

.order-body {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
}

.order-type {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #333;
  font-size: 15px;
  min-width: 100px;
}

.order-type .el-icon {
  color: #409eff;
}

.order-detail {
  flex: 1;
  color: #666;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.order-amount {
  color: #f56c6c;
  font-size: 18px;
  font-weight: 600;
}

.currency {
  font-size: 14px;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-time {
  color: #999;
  font-size: 13px;
}

.order-refund {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--el-border-color-lighter);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.refund-countdown {
  color: var(--el-color-info);
  font-size: 13px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

@media (max-width: 768px) {
  .orders-page {
    padding: 24px 16px 40px;
  }

  .page-header h1 {
    font-size: 24px;
  }

  .orders-container {
    padding: 16px;
  }

  .order-body {
    flex-wrap: wrap;
  }

  .order-detail {
    order: 3;
    flex-basis: 100%;
    margin-top: 8px;
  }

  .order-amount {
    margin-left: auto;
  }
}
</style>
