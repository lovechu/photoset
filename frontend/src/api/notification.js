/**
 * 通知中心 API
 */
import request from '@/utils/request'

export function getNotifications(params) {
  return request.get('/community/notifications', { params })
}

export function getUnreadCount() {
  return request.get('/community/notifications/unread-count')
}

export function markNotificationRead(id) {
  return request.put(`/community/notifications/${id}/read`)
}

export function markAllNotificationsRead() {
  return request.put('/community/notifications/read-all')
}

export function deleteNotification(id) {
  return request.delete(`/community/notifications/${id}`)
}

export function deleteAllNotifications() {
  return request.delete('/community/notifications')
}
