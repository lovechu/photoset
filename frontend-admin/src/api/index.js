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

export function batchApprovePhotoSets(ids) {
  return request.post('/admin/photosets/batch/approve', { ids })
}

export function batchRejectPhotoSets(ids, reason) {
  return request.post('/admin/photosets/batch/reject', { ids, reason })
}

export function batchDeletePhotoSets(ids) {
  return request.post('/admin/photosets/batch/delete', { ids })
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

export function exportUsers(params) {
  return request.get('/admin/users/export', { params, responseType: 'blob' })
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

export function createUser(data) {
  return request.post('/admin/users', data)
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

// 导出套图列表
export function exportPhotosets(params) {
  return request.get('/admin/photosets/export', { params, responseType: 'blob' })
}

// 获取订单列表（带分页和筛选）
export function getOrderList(params) {
  return request.get('/admin/orders', { params })
}

export function exportOrders(params) {
  return request.get('/admin/orders/export', { params, responseType: 'blob' })
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

// ============ 支付配置 ============

export function testAlipayConnection(data) {
  return request.post('/admin/payment/alipay/test', data)
}

export function testWechatPayConnection(data) {
  return request.post('/admin/payment/wechat/test', data)
}

// ============ 会员套餐管理 APIs ============

export function getMembershipList(params) {
  return request.get('/admin/memberships', { params })
}

export function createMembership(data) {
  return request.post('/admin/memberships', data)
}

export function updateMembership(id, data) {
  return request.put(`/admin/memberships/${id}`, data)
}

export function deleteMembership(id) {
  return request.delete(`/admin/memberships/${id}`)
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

// ============ 系统监控 ============

export function getSystemStatus() {
  return request.get('/admin/system/status')
}

export function getSystemHealth() {
  return request.get('/admin/system/health')
}

/**
 * 重启后端服务（仅管理员）
 */
export function restartServer() {
  return request.post('/admin/system/restart')
}

/**
 * 健康检查（公开接口，用于重启后轮询）
 * baseURL 已是 /api，直接请求 /health 即可
 */
export function healthCheck() {
  return request.get('/health', { timeout: 5000 })
}

// ============ 数据备份 ============

export function createBackup() {
  return request.post('/admin/backups')
}

export function getBackupList() {
  return request.get('/admin/backups')
}

export function downloadBackup(filename) {
  return request.get(`/admin/backups/${filename}/download`, { responseType: 'blob' })
}

export function deleteBackup(filename) {
  return request.delete(`/admin/backups/${filename}`)
}

// ============ OAuth2 应用管理 ============

export function getOAuthClients() {
  return request.get('/admin/oauth/clients')
}

export function createOAuthClient(data) {
  return request.post('/admin/oauth/clients', data)
}

export function updateOAuthClient(id, data) {
  return request.put(`/admin/oauth/clients/${id}`, data)
}

export function deleteOAuthClient(id) {
  return request.delete(`/admin/oauth/clients/${id}`)
}

export function getOAuthClient(id) {
  return request.get(`/admin/oauth/clients/${id}`)
}

// ============ 回收站（软删除管理） ============

/**
 * 获取回收站（软删除）套图列表
 */
export function getTrashList(params) {
  return request.get('/photosets/trash', { params })
}

/**
 * 恢复软删除的套图
 */
export function restorePhotoset(id) {
  return request.post(`/photosets/${id}/restore`)
}

/**
 * 永久删除套图
 */
export function permanentDeletePhotoset(id) {
  return request.delete(`/photosets/${id}/permanent`)
}

// ============ 用户管理增强 APIs ============

/**
 * 获取用户登录历史（管理员）
 * @param {number} userId 用户ID
 * @param {object} params { page, page_size }
 */
export function getUserLoginHistory(userId, params) {
  return request.get(`/admin/users/${userId}/login-history`, { params })
}

/**
 * 获取用户设备列表（管理员）
 * @param {number} userId 用户ID
 */
export function getUserDevices(userId) {
  return request.get(`/admin/users/${userId}/devices`)
}

/**
 * 停用用户设备（管理员）
 * @param {number} userId 用户ID
 * @param {string} deviceId 设备ID
 */
export function deactivateUserDevice(userId, deviceId) {
  return request.delete(`/admin/users/${userId}/devices/${deviceId}`)
}

/**
 * 获取用户隐私设置（管理员）
 * @param {number} userId 用户ID
 */
export function getUserPrivacySettings(userId) {
  return request.get(`/admin/users/${userId}/privacy-settings`)
}

/**
 * 更新用户隐私设置（管理员）
 * @param {number} userId 用户ID
 * @param {object} data 隐私设置数据
 */
export function updateUserPrivacySettings(userId, data) {
  return request.put(`/admin/users/${userId}/privacy-settings`, data)
}
