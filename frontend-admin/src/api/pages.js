import request from '@/utils/request'

/**
 * 获取页面列表（管理员）
 * @param {object} params - 分页和筛选参数
 * @returns {Promise}
 */
export function getPageList(params) {
  return request.get('/admin/pages', { params })
}

/**
 * 创建页面
 * @param {object} data - 页面数据
 * @returns {Promise}
 */
export function createPage(data) {
  return request.post('/admin/pages', data)
}

/**
 * 获取页面详情（管理员编辑用）
 * @param {number} id - 页面 ID
 * @returns {Promise}
 */
export function getPage(id) {
  return request.get(`/admin/pages/${id}`)
}

/**
 * 更新页面
 * @param {number} id - 页面 ID
 * @param {object} data - 更新数据
 * @returns {Promise}
 */
export function updatePage(id, data) {
  return request.put(`/admin/pages/${id}`, data)
}

/**
 * 删除页面
 * @param {number} id - 页面 ID
 * @returns {Promise}
 */
export function deletePage(id) {
  return request.delete(`/admin/pages/${id}`)
}