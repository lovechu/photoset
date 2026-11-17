/**
 * 通知状态管理
 * 集中管理未读计数，避免多组件重复请求
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getNotifications, getUnreadCount, markNotificationRead, markAllNotificationsRead, deleteNotification, deleteAllNotifications } from '@/api'

export const useNotificationStore = defineStore('notification', () => {
  const unreadCount = ref(0)
  const notifications = ref([])
  const total = ref(0)
  const loading = ref(false)

  const hasUnread = computed(() => unreadCount.value > 0)

  async function fetchUnreadCount() {
    try {
      const res = await getUnreadCount()
      const data = extractData(res)
      // 后端返回 { unread_count: N }
      if (data && typeof data.unread_count === 'number') {
        unreadCount.value = data.unread_count
      } else if (data && typeof data.count === 'number') {
        unreadCount.value = data.count
      } else if (typeof data === 'number') {
        unreadCount.value = data
      }
    } catch (e) {
      console.warn('获取未读通知数失败:', e.message)
    }
  }

  async function fetchNotifications(params = {}) {
    loading.value = true
    try {
      const res = await getNotifications(params)
      const data = extractData(res)
      // 后端返回 { notifications: [], pagination: { total, page, page_size } }
      if (data && Array.isArray(data.notifications)) {
        notifications.value = data.notifications
        total.value = data.pagination?.total || data.notifications.length
      } else if (data && Array.isArray(data.list)) {
        notifications.value = data.list
        total.value = data.total || data.list.length
      } else if (Array.isArray(data)) {
        notifications.value = data
        total.value = data.length
      }
      return notifications.value
    } catch (e) {
      console.warn('获取通知列表失败:', e.message)
      return []
    } finally {
      loading.value = false
    }
  }

  async function markRead(id) {
    try {
      await markNotificationRead(id)
      const item = notifications.value.find(n => n.id === id)
      if (item && !item.is_read) {
        item.is_read = true
        if (unreadCount.value > 0) unreadCount.value--
      }
    } catch (e) {
      console.warn('标记通知已读失败:', e.message)
    }
  }

  async function markAllRead() {
    try {
      await markAllNotificationsRead()
      notifications.value.forEach(n => { n.is_read = true })
      unreadCount.value = 0
    } catch (e) {
      console.warn('全部标记已读失败:', e.message)
    }
  }

  async function removeNotification(id) {
    try {
      await deleteNotification(id)
      const item = notifications.value.find(n => n.id === id)
      notifications.value = notifications.value.filter(n => n.id !== id)
      if (item && !item.is_read && unreadCount.value > 0) unreadCount.value--
    } catch (e) {
      console.warn('删除通知失败:', e.message)
    }
  }

  async function removeAll() {
    try {
      await deleteAllNotifications()
      notifications.value = []
      unreadCount.value = 0
    } catch (e) {
      console.warn('清空通知失败:', e.message)
    }
  }

  // WebSocket 推送新通知时调用
  function pushNewNotification(notification) {
    notifications.value.unshift(notification)
    unreadCount.value++
  }

  function reset() {
    unreadCount.value = 0
    notifications.value = []
    total.value = 0
  }

  return {
    unreadCount,
    notifications,
    total,
    loading,
    hasUnread,
    fetchUnreadCount,
    fetchNotifications,
    markRead,
    markAllRead,
    removeNotification,
    removeAll,
    pushNewNotification,
    reset
  }
})

function extractData(res) {
  if (res && res.data && typeof res.data === 'object') return res.data
  if (res && typeof res === 'object' && !res.code) return res
  return null
}
