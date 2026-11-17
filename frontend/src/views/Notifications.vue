<template>
  <div class="notifications-page">
    <div class="page-header">
      <h2>通知中心</h2>
      <div>
        <el-button size="small" :icon="Check" @click="markAllRead">全部已读</el-button>
        <el-button size="small" type="danger" plain :icon="Delete" @click="clearAll" v-if="list.length">清空</el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab" @tab-change="reload">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="未读" :name="'unread'">
        <template #label>未读 <span v-if="store.unreadCount" class="unread-badge">{{ store.unreadCount }}</span></template>
      </el-tab-pane>
    </el-tabs>

    <div v-loading="store.loading">
      <el-empty v-if="!store.loading && list.length === 0" description="暂无通知" :image-size="100" />
      <div v-else class="notification-list">
        <div
          v-for="item in list"
          :key="item.id"
          :class="['notification-item', { unread: !item.is_read }]"
          @click="handleClick(item)"
        >
          <el-avatar :size="40" :src="item.sender_avatar">
            {{ getIconText(item.type) }}
          </el-avatar>
          <div class="notif-main">
            <div class="notif-head">
              <span class="notif-title">{{ item.title }}</span>
              <span class="notif-time">{{ formatTime(item.created_at) }}</span>
            </div>
            <div class="notif-content">{{ item.content }}</div>
            <div v-if="item.sender_name" class="notif-sender">来自 {{ item.sender_name }}</div>
          </div>
          <div class="notif-actions">
            <span v-if="!item.is_read" class="unread-dot" />
            <el-button link size="small" :icon="Delete" @click.stop="removeItem(item)" />
          </div>
        </div>
      </div>

      <div class="load-more" v-if="list.length < store.total">
        <el-button @click="loadMore">加载更多</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Check, Delete } from '@element-plus/icons-vue'
import { useNotificationStore } from '@/stores/notification'
import { onMessage } from '@/utils/websocket'

const router = useRouter()
const store = useNotificationStore()
const activeTab = ref('all')
const page = ref(1)
let unsub = null

const list = computed(() => store.notifications)

onMounted(() => {
  reload()
  // WebSocket 实时推送新通知
  unsub = onMessage('notification', (payload) => {
    store.pushNewNotification(payload?.notification || payload || {})
    ElMessage.info(payload?.title || '收到一条新通知')
  })
})

onUnmounted(() => {
  if (unsub) unsub()
})

function reload() {
  page.value = 1
  loadNotifications()
}

async function loadNotifications() {
  try {
    await store.fetchNotifications({
      page: page.value,
      page_size: 20,
      unread_only: activeTab.value === 'unread' ? 'true' : undefined
    })
  } catch (e) { /* store 内已处理 */ }
}

function loadMore() {
  page.value++
  loadNotifications()
}

async function markAllRead() {
  await store.markAllRead()
  ElMessage.success('已全部标记为已读')
}

async function clearAll() {
  try {
    await ElMessageBox.confirm('确定清空所有通知？', '提示', { type: 'warning' })
    await store.removeAll()
    ElMessage.success('已清空')
  } catch (e) { /* 取消 */ }
}

async function removeItem(item) {
  await store.removeNotification(item.id)
}

async function handleClick(item) {
  if (!item.is_read) await store.markRead(item.id)
  // 根据通知类型跳转
  if (item.target_type === 'post' && item.target_id) {
    router.push(`/community/post/${item.target_id}`)
  } else if (item.target_type === 'photoset' && item.target_id) {
    router.push(`/photoset/${item.target_id}`)
  }
}

function getIconText(type) {
  const map = { like: '赞', reply: '回', follow: '关', mention: '@', system: '系' }
  return map[type] || '知'
}

function formatTime(str) {
  if (!str) return ''
  const d = new Date(str)
  if (isNaN(d.getTime())) return str
  const diff = Date.now() - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins}分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}小时前`
  return d.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.notifications-page { max-width: 800px; margin: 0 auto; padding: 24px 0 60px; }
.page-header { display: flex; align-items: center; justify-content: space-between; }
.page-header h2 { margin: 0; }
.page-header > div { display: flex; gap: 8px; }
.unread-badge { background: #f56c6c; color: #fff; border-radius: 10px; padding: 1px 6px; font-size: 11px; }
.notification-list { background: #fff; border-radius: 12px; overflow: hidden; }
.notification-item { display: flex; align-items: center; gap: 14px; padding: 16px; border-bottom: 1px solid #f0f0f0; cursor: pointer; transition: background 0.2s; }
.notification-item:last-child { border-bottom: none; }
.notification-item:hover { background: #f8f8f8; }
.notification-item.unread { background: #f0f9ff; }
.notif-main { flex: 1; min-width: 0; }
.notif-head { display: flex; align-items: center; justify-content: space-between; }
.notif-title { font-weight: 600; font-size: 14px; }
.notif-time { font-size: 12px; color: #999; }
.notif-content { color: #666; font-size: 13px; margin-top: 4px; }
.notif-sender { color: #aaa; font-size: 12px; margin-top: 2px; }
.notif-actions { display: flex; align-items: center; gap: 10px; }
.unread-dot { width: 8px; height: 8px; border-radius: 50%; background: #f56c6c; }
.load-more { text-align: center; padding: 20px; }
</style>
