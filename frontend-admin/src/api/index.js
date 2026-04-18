import request from '@/utils/request'

// ============ 认证模块 ============

/**
 * 获取验证码
 * @returns { captcha_id, captcha_image }
 */
export function getCaptcha() {
  return request.get('/auth/captcha')
}

export function login(data) {
  return request.post('/auth/login', data)
}

export function getCurrentUser() {
  return request.get('/auth/me')
}

export function getStats() {
  return request.get('/admin/stats')
}

export function getStatsTrend(days = 7) {
  return request.get('/admin/stats/trend', { params: { days } })
}

// ============ 操作日志 ============

export function getAdminLogs(params) {
  return request.get('/admin/logs', { params })
}

export function getPhotoSetsByStatus(params) {
  return request.get('/admin/photosets', { params })
}

export function approvePhotoSet(id) {
  return request.post(`/admin/photosets/${id}/approve`)
}

export function rejectPhotoSet(id, reason = '') {
  return request.post(`/admin/photosets/${id}/reject`, { reason })
}

export function getUserList(params) {
  return request.get('/admin/users', { params })
}

export function banUser(id, status) {
  return request.put(`/admin/users/${id}/ban`, { status: Number(status) })
}

export function updateUserRole(id, role) {
  return request.put(`/admin/users/${id}/role`, { role })
}

export function getUserDetail(id) {
  return request.get(`/admin/users/${id}`)
}

export function resetUserPassword(id, newPassword) {
  return request.put(`/admin/users/${id}/password`, { new_password: newPassword })
}

// ============ 套图编辑模块 ============

export function getPhotosetDetail(id) {
  return request.get(`/photosets/${id}`)
}

export function updatePhotoset(id, data) {
  return request.put(`/photosets/${id}`, data)
}

export function getTags() {
  return request.get('/tags')
}

export function uploadImage(file) {
  const formData = new FormData()
  formData.append('image', file)
  return request.post('/upload/image', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// ============ Phase 5 新增接口 ============

// 管理员删除套图
export function deletePhotoset(id) {
  return request.delete(`/photosets/${id}`)
}

// 获取订单列表（带分页和筛选）
export function getOrderList(params) {
  return request.get('/admin/orders', { params })
}

// 强制退款（管理员无时间限制）
export function adminRefundOrder(id) {
  return request.post(`/admin/orders/${id}/refund`)
}

// 标签管理 APIs
export function getTagList(params) {
  return request.get('/admin/tags', { params })
}

export function createTag(data) {
  return request.post('/admin/tags', data)
}

export function updateTag(id, data) {
  return request.put(`/admin/tags/${id}`, data)
}

export function deleteTag(id) {
  return request.delete(`/admin/tags/${id}`)
}

// ============ 分类管理 APIs ============

export function getCategoryList(params) {
  return request.get('/admin/categories', { params })
}

export function createCategory(data) {
  return request.post('/admin/categories', data)
}

export function updateCategory(id, data) {
  return request.put(`/admin/categories/${id}`, data)
}

export function deleteCategory(id) {
  return request.delete(`/admin/categories/${id}`)
}

export function getPublicCategories() {
  return request.get('/categories')
}

// ============ 站点设置 ============

export function getSettings() {
  return request.get('/admin/settings').then(res => res.data?.data || res.data)
}

export function updateSettings(data) {
  return request.put('/admin/settings', data)
}

// ============ 存储配置 ============

export function getStorageStatus() {
  return request.get('/admin/storage/status')
}

export function testStorageConnection(data) {
  return request.post('/admin/storage/test', data)
}

// ============ 邮件配置 ============

export function testMailConnection() {
  return request.post('/admin/mail/test-connection')
}

export function getMailConfig() {
  return request.get('/admin/mail/config')
}

export function sendMailTest(data) {
  return request.post('/admin/mail/send-test', data)
}

// ============ 水印配置 ============

export function getWatermarkInfo() {
  return request.get('/admin/watermark/info')
}

// ============ 开发者中心 ============

export function getApiKeys() {
  return request.get('/admin/dev/api-keys')
}

export function createApiKey(name) {
  return request.post('/admin/dev/api-keys', { name })
}

export function deleteApiKey(id) {
  return request.delete(`/admin/dev/api-keys/${id}`)
}

export function getApiDocs() {
  return request.get('/admin/dev/api-docs')
}

export function getSignUrlDocs() {
  return request.get('/admin/dev/sign-url-docs')
}
