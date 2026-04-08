/**
 * API 接口定义
 * 基于后端已实现的接口
 */
import request from '@/utils/request'

// ============ 认证模块 ============

/**
 * 用户注册
 * @param {Object} data { nickname, email, password }
 */
export function register(data) {
  return request.post('/auth/register', data)
}

/**
 * 用户登录
 * @param {Object} data { email, password }
 * @returns { token, user }
 */
export function login(data) {
  return request.post('/auth/login', data)
}

/**
 * 获取当前用户信息
 * 需要登录态
 */
export function getCurrentUser() {
  return request.get('/auth/me')
}

// ============ 套图模块 ============

/**
 * 获取套图列表
 * @param {Object} params { page, page_size, tag }
 */
export function getPhotosetList(params) {
  return request.get('/photosets', { params })
}

/**
 * 获取套图详情
 * 可选携带 token，后端根据身份决定返回内容
 * @param {Number} id 套图 ID
 */
export function getPhotosetDetail(id) {
  return request.get(`/photosets/${id}`)
}

/**
 * 创建套图
 * 需要 creator 或 admin 权限
 * @param {Object} data 套图数据
 */
export function createPhotoset(data) {
  return request.post('/photosets', data)
}

// ============ 标签模块 ============

/**
 * 获取所有标签
 */
export function getTags() {
  return request.get('/tags')
}

// ============ 健康检查 ============

/**
 * 健康检查
 */
export function healthCheck() {
  return request.get('/health')
}

// ============ 收藏模块 ============

/**
 * 收藏套图（幂等，重复收藏返回成功）
 * @param {Number} photosetId 套图 ID
 */
export function addFavorite(photosetId) {
  return request.post(`/favorites/${photosetId}`)
}

/**
 * 取消收藏
 * @param {Number} photosetId 套图 ID
 */
export function removeFavorite(photosetId) {
  return request.delete(`/favorites/${photosetId}`)
}

/**
 * 获取我的收藏列表
 * @param {Object} params { page, page_size }
 */
export function getFavorites(params) {
  return request.get('/favorites', { params })
}

// ============ 上传模块 ============

/**
 * 上传图片
 * @param {File} file 文件对象
 * @returns { url } 可访问的图片 URL
 */
export function uploadImage(file) {
  const formData = new FormData()
  formData.append('image', file)
  return request.post('/upload/image', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
