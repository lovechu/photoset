import request from '@/utils/request'

// ============ 等级配置管理 ============

/**
 * 获取所有等级配置
 */
export function getLevelConfigs() {
  return request.get('/admin/community/levels')
}

/**
 * 更新等级配置
 * @param {number} id - 等级ID
 * @param {Object} data - 等级配置数据
 */
export function updateLevelConfig(id, data) {
  return request.put(`/admin/community/levels/${id}`, data)
}

// ============ 成就管理 ============

/**
 * 获取所有成就列表（管理员）
 */
export function getAchievements() {
  return request.get('/admin/community/achievements')
}

/**
 * 新建成就
 * @param {Object} data - 成就数据
 */
export function createAchievement(data) {
  return request.post('/admin/community/achievements', data)
}

/**
 * 更新成就
 * @param {number} id - 成就ID
 * @param {Object} data - 成就数据
 */
export function updateAchievement(id, data) {
  return request.put(`/admin/community/achievements/${id}`, data)
}

/**
 * 删除成就
 * @param {number} id - 成就ID
 */
export function deleteAchievement(id) {
  return request.delete(`/admin/community/achievements/${id}`)
}

// ============ 积分商城管理 ============

/**
 * 获取积分商城商品列表（管理员）
 */
export function getPointsMallItems() {
  return request.get('/admin/community/points-mall/items')
}

/**
 * 创建积分商城商品
 * @param {Object} data - 商品数据
 */
export function createPointsMallItem(data) {
  return request.post('/admin/community/points-mall/items', data)
}

/**
 * 更新积分商城商品
 * @param {number} id - 商品ID
 * @param {Object} data - 商品数据
 */
export function updatePointsMallItem(id, data) {
  return request.put(`/admin/community/points-mall/items/${id}`, data)
}

/**
 * 删除积分商城商品
 * @param {number} id - 商品ID
 */
export function deletePointsMallItem(id) {
  return request.delete(`/admin/community/points-mall/items/${id}`)
}

// ============ 兑换记录 ============

/**
 * 获取所有兑换记录（管理员）
 * @param {Object} params - { page, page_size }
 */
export function getAllExchanges(params) {
  return request.get('/admin/community/exchanges', { params })
}
