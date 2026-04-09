import request from '@/utils/request'

export function login(data) {
  return request.post('/auth/login', data)
}

export function getCurrentUser() {
  return request.get('/auth/me')
}

export function getStats() {
  return request.get('/admin/stats')
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
  return request.put(`/admin/users/${id}/ban`, { status })
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
