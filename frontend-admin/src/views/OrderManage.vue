<template>
  <div class="order-manage">
    <div class="filter-bar">
      <el-button type="success" plain @click="handleExport" :loading="exporting">
        导出 CSV
      </el-button>
      <el-input
        v-model="filterKeyword"
        placeholder="搜索订单ID/用户/套图标题"
        clearable
        @clear="fetchOrders"
        @keyup.enter="fetchOrders"
        style="width: 260px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      
      <el-select v-model="filterStatus" placeholder="状态筛选" clearable @change="fetchOrders" style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="待支付" value="pending" />
        <el-option label="已支付" value="paid" />
        <el-option label="已退款" value="refunded" />
        <el-option label="已取消" value="cancelled" />
      </el-select>
      
      <el-select v-model="filterType" placeholder="类型筛选" clearable @change="fetchOrders" style="width: 140px">
        <el-option label="全部" value="" />
        <el-option label="套图购买" value="photoset" />
        <el-option label="会员订阅" value="membership" />
      </el-select>
      
      <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        value-format="YYYY-MM-DD"
        @change="fetchOrders"
        style="width: 320px"
      />
    </div>

    <el-table :data="orderList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="订单ID" width="80" align="center" />
      <el-table-column label="用户" min-width="160">
        <template #default="{ row }">
          <div v-if="row.user">
            <div class="user-info">
              <div>{{ row.user.nickname }}</div>
              <div class="user-email">{{ row.user.email }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="订单内容" min-width="200">
        <template #default="{ row }">
          <div v-if="row.target_type === 'photoset' && row.photoset">
            <div class="order-content">
              <strong>套图:</strong> {{ row.photoset.title }}
            </div>
          </div>
          <div v-else-if="row.target_type === 'membership'">
            <div class="order-content">
              <strong>会员类型:</strong> {{ formatMembershipType(row.target_id) }}
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="金额" width="100" align="right">
        <template #default="{ row }">
          ¥{{ (row.amount / 100).toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.target_type === 'photoset' ? 'warning' : 'success'" size="small">
            {{ row.target_type === 'photoset' ? '套图' : '会员' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170" align="center">
        <template #default="{ row }">
          {{ formatTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="支付/退款时间" width="170" align="center">
        <template #default="{ row }">
          <div v-if="row.paid_at">{{ formatTime(row.paid_at) }}</div>
          <div v-else-if="row.refunded_at">{{ formatTime(row.refunded_at) }}</div>
          <div v-else>-</div>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" align="center" fixed="right">
        <template #default="{ row }">
          <div class="action-buttons">
            <el-button
              v-if="row.status === 'paid'"
              type="danger"
              size="small"
              plain
              @click="handleRefund(row)"
              :loading="refundingId === row.id"
            >
              强制退款
            </el-button>
            <el-text v-else-if="row.status === 'refunded'" type="danger" size="small">
              已退款
            </el-text>
            <el-text v-else-if="row.status === 'cancelled'" type="info" size="small">
              已取消
            </el-text>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { getOrderList, adminRefundOrder, exportOrders } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'

const loading = ref(false)
const orderList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 筛选条件
const filterKeyword = ref('')
const filterStatus = ref('')
const filterType = ref('')
const dateRange = ref([])

// 导出状态
const exporting = ref(false)

// 退款状态
const refundingId = ref(null)

// 订单状态映射
const statusMap = {
  pending: '待支付',
  paid: '已支付',
  refunded: '已退款',
  cancelled: '已取消'
}

const statusTagMap = {
  pending: 'warning',
  paid: 'success',
  refunded: 'danger',
  cancelled: 'info'
}

function statusLabel(s) {
  return statusMap[s] || s
}

function statusTagType(s) {
  return statusTagMap[s] || 'info'
}

function formatTime(t) {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}

function formatMembershipType(typeId) {
  const map = {
    1: '月度会员',
    2: '季度会员',
    3: '年度会员'
  }
  return map[typeId] || `会员类型 ${typeId}`
}

async function fetchOrders() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value
    }

    if (filterKeyword.value) {
      params.keyword = filterKeyword.value
    }
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    if (filterType.value) {
      params.type = filterType.value
    }
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }

    const res = await getOrderList(params)
    orderList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('获取订单列表失败:', error)
    ElMessage.error('获取订单列表失败')
  } finally {
    loading.value = false
  }
}

function handleSizeChange(val) {
  pageSize.value = val
  currentPage.value = 1
  fetchOrders()
}

function handleCurrentChange(val) {
  currentPage.value = val
  fetchOrders()
}

async function handleRefund(order) {
  try {
    await ElMessageBox.confirm(
      `确定要强制退款订单 ${order.id} 吗？金额 ¥${(order.amount / 100).toFixed(2)} 将退还给用户。`,
      '强制退款确认',
      {
        confirmButtonText: '确认退款',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    refundingId.value = order.id
    await adminRefundOrder(order.id)
    ElMessage.success('退款成功')
    fetchOrders()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('退款失败')
    }
  } finally {
    refundingId.value = null
  }
}

// 监听筛选条件变化
watch([filterKeyword, filterStatus, filterType, dateRange], () => {
  currentPage.value = 1
  fetchOrders()
})

async function handleExport() {
  exporting.value = true
  try {
    const params = {}
    if (filterKeyword.value) params.keyword = filterKeyword.value
    if (filterStatus.value) params.status = filterStatus.value
    if (filterType.value) params.type = filterType.value
    if (dateRange.value && dateRange.value.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    const res = await exportOrders(params)
    const blob = new Blob([res], { type: 'text/csv;charset=utf-8' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'orders.csv'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

onMounted(fetchOrders)
</script>

<style scoped>
.filter-bar {
  margin-bottom: 20px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
}

.user-info .user-email {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.order-content {
  font-size: 13px;
  line-height: 1.4;
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 8px;
}
</style>