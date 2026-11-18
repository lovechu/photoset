/**
 * 用户资料与安全 API
 */
import request from '@/utils/request'

export function updateProfile(data) {
  return request.put('/auth/profile', data)
}

export function uploadAvatar(file) {
  const formData = new FormData()
  formData.append('image', file)
  return request.post('/user/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export function sendEmailCode(data) {
  return request.post('/auth/send-email-code', data)
}

export function verifyEmailCode(data) {
  return request.post('/auth/verify-email-code', data)
}

export function bindEmail(data) {
  return request.put('/auth/bind-email', data)
}

export function getLoginHistory(params) {
  return request.get('/user/login-history', { params })
}

export function getUserDevices() {
  return request.get('/user/devices')
}

export function deactivateDevice(deviceId) {
  return request.delete(`/user/devices/${deviceId}`)
}

export function deactivateOtherDevices() {
  return request.delete('/user/devices')
}

export function getPrivacySettings() {
  return request.get('/user/privacy-settings')
}

export function updatePrivacySettings(data) {
  return request.put('/user/privacy-settings', data)
}

export function requestDeletion() {
  return request.post('/auth/request-deletion')
}

export function cancelDeletion() {
  return request.post('/auth/cancel-deletion')
}

export function getDeletionStatus() {
  return request.get('/auth/deletion-status')
}

// ============ 浏览历史 ============

export function recordViewHistory(photosetId) {
  return request.post('/user/view-history', { photoset_id: photosetId })
}

export function listViewHistory(params) {
  return request.get('/user/view-history', { params })
}

export function deleteViewHistory(id) {
  return request.delete(`/user/view-history/${id}`)
}

export function clearViewHistory() {
  return request.delete('/user/view-history')
}

export function batchDeleteViewHistory(data) {
  return request.post('/user/view-history/batch-delete', data)
}

// ============ 推送 Token ============

export function registerPushToken(data) {
  return request.post('/user/push-tokens', data)
}

export function unregisterPushToken(data) {
  return request.delete('/user/push-tokens', { data })
}
