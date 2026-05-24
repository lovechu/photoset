import request from '@/utils/request'

// ============ 帖子管理 ============

/**
 * 获取社区帖子列表
 * @param {Object} params - { page, page_size, status }
 */
export function getCommunityPosts(params) {
  return request.get('/admin/community/posts', { params })
}

/**
 * 置顶/取消置顶帖子
 * @param {number} id - 帖子ID
 */
export function togglePostPin(id) {
  return request.put(`/admin/community/posts/${id}/pin`)
}

/**
 * 更新帖子状态（审核通过/拒绝）
 * @param {number} id - 帖子ID
 * @param {Object} data - { status }
 */
export function updatePostStatus(id, data) {
  return request.put(`/admin/community/posts/${id}/status`, data)
}

/**
 * 删除帖子
 * @param {number} id - 帖子ID
 */
export function deletePost(id) {
  return request.delete(`/admin/community/posts/${id}`)
}

// ============ 回帖管理 ============

/**
 * 获取回帖列表
 * @param {Object} params - { page, page_size, post_id }
 */
export function getCommunityReplies(params) {
  return request.get('/admin/community/replies', { params })
}

/**
 * 删除回帖
 * @param {number} id - 回帖ID
 */
export function deleteReply(id) {
  return request.delete(`/admin/community/replies/${id}`)
}

// ============ 敏感词管理 ============

/**
 * 获取敏感词列表
 * @param {Object} params - { page, page_size, is_active }
 */
export function getKeywords(params) {
  return request.get('/admin/community/keywords', { params })
}

/**
 * 新增敏感词
 * @param {Object} data - { word, replacement }
 */
export function addKeyword(data) {
  return request.post('/admin/community/keywords', data)
}

/**
 * 更新敏感词
 * @param {number} id - 敏感词ID
 * @param {Object} data - { word, replacement, is_active }
 */
export function updateKeyword(id, data) {
  return request.put(`/admin/community/keywords/${id}`, data)
}

/**
 * 删除敏感词
 * @param {number} id - 敏感词ID
 */
export function deleteKeyword(id) {
  return request.delete(`/admin/community/keywords/${id}`)
}

/**
 * 重新加载敏感词到内存
 */
export function reloadKeywords() {
  return request.put('/admin/community/keywords/reload')
}

// ============ 举报处理 ============

/**
 * 获取举报列表
 * @param {Object} params - { page, page_size, status }
 */
export function getReports(params) {
  return request.get('/admin/community/reports', { params })
}

/**
 * 处理举报
 * @param {number} id - 举报ID
 * @param {Object} data - { status, note }
 */
export function resolveReport(id, data) {
  return request.put(`/admin/community/reports/${id}/resolve`, data)
}

// ============ 用户积分管理 ============

/**
 * 获取社区用户列表
 * @param {Object} params - { page, page_size, level }
 */
export function getCommunityUsers(params) {
  return request.get('/admin/community/users', { params })
}

/**
 * 调整用户积分
 * @param {number} id - 用户ID
 * @param {Object} data - { points, reason }
 */
export function adjustUserPoints(id, data) {
  return request.put(`/admin/community/users/${id}/points`, data)
}

// ============ 统计 ============

/**
 * 获取社区统计数据
 */
export function getCommunityStats() {
  return request.get('/admin/community/stats')
}
