<template>
  <div class="conversations-page">
    <el-card>
      <template #header>
        <div class="header-row">
          <h3>私信</h3>
          <el-badge :value="messageUnread" :hidden="!messageUnread">
            <el-button :icon="Refresh" circle @click="loadConversations" />
          </el-badge>
        </div>
      </template>
      <div v-loading="loading">
        <el-empty v-if="!loading && conversations.length === 0" description="还没有私信，去社区认识新朋友吧" :image-size="80" />
        <div v-for="c in conversations" :key="c.user_id" class="conversation-item" @click="openChat(c)">
          <el-badge :value="c.unread_count" :hidden="!c.unread_count">
            <el-avatar :size="46" :src="c.avatar">{{ c.username?.charAt(0) || 'U' }}</el-avatar>
          </el-badge>
          <div class="conv-info">
            <div class="conv-name">{{ c.username || `用户 ${c.user_id}` }}</div>
            <div class="conv-last">{{ c.last_message || '开始聊天吧' }}</div>
          </div>
          <span class="conv-time">{{ formatTime(c.last_msg_time) }}</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Refresh } from '@element-plus/icons-vue'
import { getConversations, getMessageUnreadCount } from '@/api'
import { formatTime } from '@/utils/post'
import { onMessage } from '@/utils/websocket'

const router = useRouter()
const conversations = ref([])
const loading = ref(false)
const messageUnread = ref(0)
let unsubscribe = null

onMounted(() => {
  loadConversations()
  loadUnread()
  unsubscribe = onMessage('message', () => {
    loadConversations()
    loadUnread()
  })
})

onUnmounted(() => {
  if (unsubscribe) {
    unsubscribe()
    unsubscribe = null
  }
})

async function loadConversations() {
  loading.value = true
  try {
    const res = await getConversations()
    const data = res?.data || {}
    conversations.value = data.conversations || []
  } catch (e) { console.warn('加载会话失败:', e.message) }
  finally { loading.value = false }
}

async function loadUnread() {
  try {
    const res = await getMessageUnreadCount()
    const data = res?.data || {}
    messageUnread.value = data.unread_count ?? data.count ?? 0
  } catch (e) { /* 静默 */ }
}

function openChat(c) {
  router.push({
    path: `/community/chat/${c.user_id}`,
    query: { username: c.username, avatar: c.avatar }
  })
}
</script>

<style scoped>
.conversations-page { max-width: 720px; margin: 0 auto; padding: 24px 0 60px; }
.header-row { display: flex; align-items: center; justify-content: space-between; }
.header-row h3 { margin: 0; }
.conversation-item { display: flex; align-items: center; gap: 14px; padding: 14px 0; border-bottom: 1px solid #f0f0f0; cursor: pointer; }
.conversation-item:hover { background: #f8f8f8; }
.conv-info { flex: 1; min-width: 0; }
.conv-name { font-weight: 600; }
.conv-last { color: #999; font-size: 13px; margin-top: 2px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.conv-time { color: #bbb; font-size: 12px; }
</style>
