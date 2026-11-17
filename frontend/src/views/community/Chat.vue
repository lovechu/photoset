<template>
  <div class="chat-page">
    <el-card>
      <template #header>
        <div class="chat-header">
          <el-button :icon="ArrowLeft" text @click="$router.back()" />
          <el-avatar :size="32" :src="route.query.avatar">{{ username?.charAt(0) }}</el-avatar>
          <span class="chat-name">{{ username || userId }}</span>
        </div>
      </template>

      <div class="messages-box" ref="messagesBox">
        <el-empty v-if="!messages.length" description="发送消息开始聊天" :image-size="60" />
        <div v-for="msg in messages" :key="msg.id" :class="['message', { mine: msg.from_user_id === myId }]">
          <el-avatar :size="34" :src="msg.from_user_id === myId ? myAvatar : (msg.from_user?.avatar || route.query.avatar)">
            {{ msg.from_user_id === myId ? '我' : (username?.charAt(0) || 'U') }}
          </el-avatar>
          <div class="message-body">
            <div class="message-content">{{ msg.content }}</div>
            <div class="message-time">{{ formatTime(msg.created_at) }}</div>
          </div>
        </div>
      </div>

      <div class="chat-compose">
        <el-input
          v-model="input"
          :rows="2"
          type="textarea"
          placeholder="输入消息，Enter 发送"
          resize="none"
          @keydown.enter.exact.prevent="send"
        />
        <el-button type="primary" :disabled="!input.trim()" @click="send">发送</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { getConversation, sendMessage, markConversationRead } from '@/api'
import { formatTime } from '@/utils/post'
import { onMessage } from '@/utils/websocket'

const route = useRoute()
const userStore = useUserStore()
const userId = route.params.id
const username = computed(() => route.query.username)
const myId = computed(() => userStore.user?.id)
const myAvatar = computed(() => userStore.avatar)

const messages = ref([])
const input = ref('')
const messagesBox = ref(null)
const page = ref(1)
let unsub = null

onMounted(() => {
  loadMessages()
  markConversationRead(userId).catch(() => {})
  unsub = onMessage('message', (payload) => {
    const msg = payload?.message
    if (msg && (msg.from_user_id === Number(userId) || msg.to_user_id === Number(userId))) {
      messages.value.push(msg)
      scrollToBottom()
    }
  })
})

onUnmounted(() => { if (unsub) unsub() })

async function loadMessages() {
  try {
    const res = await getConversation(userId, { page: page.value, page_size: 50 })
    const data = res?.data || {}
    messages.value = data.messages || []
    await nextTick()
    scrollToBottom()
  } catch (e) { console.warn('加载聊天记录失败:', e.message) }
}

function scrollToBottom() {
  if (messagesBox.value) messagesBox.value.scrollTop = messagesBox.value.scrollHeight
}

async function send() {
  const content = input.value.trim()
  if (!content) return
  try {
    const res = await sendMessage({ to_user_id: Number(userId), content })
    const msg = res?.data?.message
    messages.value.push(msg || {
      id: Date.now(),
      from_user_id: myId.value,
      to_user_id: Number(userId),
      content,
      created_at: new Date().toISOString()
    })
    input.value = ''
    await nextTick()
    scrollToBottom()
  } catch (e) { /* request 拦截器已提示 */ }
}
</script>

<style scoped>
.chat-page { max-width: 720px; margin: 0 auto; padding: 24px 0 60px; }
.chat-header { display: flex; align-items: center; gap: 10px; }
.chat-name { font-weight: 600; }
.messages-box { height: 480px; overflow-y: auto; padding: 16px; background: #f5f7fa; border-radius: 8px; }
.message { display: flex; gap: 10px; margin-bottom: 16px; }
.message.mine { flex-direction: row-reverse; }
.message-body { max-width: 70%; }
.message-content { background: #fff; padding: 10px 14px; border-radius: 8px; line-height: 1.5; word-break: break-word; }
.message.mine .message-content { background: #409eff; color: #fff; }
.message-time { font-size: 11px; color: #bbb; margin-top: 4px; }
.message.mine .message-time { text-align: right; }
.chat-compose { display: flex; gap: 10px; align-items: flex-end; margin-top: 16px; }
.chat-compose .el-button { height: 40px; }
</style>
