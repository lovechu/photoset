import request from '@/utils/request'

/**
 * 获取页面详情（通过 slug）
 * @param {string} slug - 页面标识
 * @returns {Promise}
 */
export function getPageBySlug(slug) {
  return request.get(`/pages/${slug}`)
}

/**
 * 获取所有已发布的页面（用于站点地图）
 * @returns {Promise}
 */
export function getAllPages() {
  return request.get('/pages')
}