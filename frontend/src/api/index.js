/**
 * API 接口聚合入口
 * 基于后端已实现的接口
 * 原有具名导出保持向后兼容，新模块从独立文件 re-export
 */

// 新增业务模块统一从此处转出，现有页面可继续 import { xxx } from '@/api'
export * from './community'
export * from './comment'
export * from './collection'
export * from './notification'
export * from './user'
export * from './explore'
export * from './creator'
export * from './review'
export * from './feedback'
export * from './pages'

import request from '@/utils/request'

// ============ 认证模块 ============

/**
 * 获取验证码
 * @returns { captcha_id, captcha_image }
 */
export function getCaptcha() {
  return request.get('/auth/captcha')
}

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
 * 获取套图列表（基础版 - 向后兼容）
 * @param {Object} params { page, page_size, tag, keyword, mine }
 */
export function getPhotosetList(params) {
  return request.get('/photosets', { params })
}

/**
 * 获取套图列表（高级版 - 支持所有筛选参数）
 * @param {Object} params {
 *   page, page_size, tag, keyword,
 *   category, price_min, price_max, is_free,
 *   sort_by, time_range, user_id, only_mine
 * }
 */
export function getPhotosetListAdvanced(params) {
  return request.get('/photosets/advanced', { params })
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

/**
 * 更新套图
 * 需要 creator 或 admin 权限，且为套图所有者
 * @param {Number} id 套图 ID
 * @param {Object} data 套图数据
 */
export function updatePhotoset(id, data) {
  return request.put(`/photosets/${id}`, data)
}

/**
 * 删除套图
 * @param {Number} id 套图ID
 */
export function deletePhotoset(id) {
  return request.delete(`/photosets/${id}`)
}

// ============ 标签模块 ============

/**
 * 获取所有标签
 */
export function getTags() {
  return request.get('/tags')
}

/**
 * 获取所有分类（公开接口，供高级搜索和上传表单使用）
 */
export function getCategories() {
  return request.get('/categories')
}

// ============ 站点设置模块 ============

/**
 * 获取站点设置（公开接口）
 */
export function getSiteSettings() {
  return request.get('/settings').then(res => res.data)
}

/**
 * 获取完整站点设置（包含导航菜单等所有配置）
 */
export function getSettings() {
  return request.get('/settings').then(res => res.data)
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

// ============ 会员模块 ============

/**
 * 获取会员套餐列表（公开接口）
 */
export function getMemberships() {
  return request.get('/memberships')
}

// ============ 订单模块 ============

/**
 * 创建订单（会员订阅或单套图购买）
 * @param {Object} data { type: 'membership'|'single', membership_id?: number, photoset_id?: number }
 */
export function createOrder(data) {
  return request.post('/orders', data)
}

/**
 * 模拟支付订单（开发环境用）
 * @param {Number} id 订单 ID
 * @returns { token, order } 返回新 Token（角色可能变更）
 */
export function payOrder(id) {
  return request.post(`/orders/${id}/pay`)
}

/**
 * 获取我的订单列表
 * @param {Object} params { page, page_size }
 */
export function getOrders(params) {
  return request.get('/orders', { params })
}

/**
 * 申请订单退款
 * @param {Number} id 订单ID
 */
export function refundOrder(id) {
  return request.post(`/orders/${id}/refund`)
}

// ============ 密码模块 ============

/**
 * 修改密码（需登录，需验证旧密码）
 * @param {Object} data { old_password, new_password }
 */
export function changePassword(data) {
  return request.put('/auth/password', data)
}

// ============ 密码重置模块 ============

/**
 * 请求密码重置（发送重置邮件）
 * @param {Object} data { email, captcha_id, captcha_code }
 */
export function forgotPassword(data) {
  return request.post('/auth/forgot-password', data)
}

/**
 * 通过 token 重置密码
 * @param {Object} data { token, new_password }
 */
export function resetPasswordByToken(data) {
  return request.post('/auth/reset-password', data)
}

/**
 * 检查邮件配置是否可用
 */
export function checkEmailConfig() {
  return request.get('/auth/email-config')
}

// ============ OAuth2 模块 ============

/**
 * 获取授权页面信息
 * @param {Object} params { client_id, redirect_uri, scope, state, response_type }
 */
export function getOAuthAuthorize(params) {
  return request.get('/oauth/authorize', { params })
}

/**
 * 用户确认授权
 * @param {Object} data { client_id, redirect_uri, scope, state, approved }
 */
export function confirmOAuthAuthorize(data) {
  return request.post('/oauth/authorize', data)
}

// ============ 套图回收站 / 分享 / 下载 ============

/** 回收站列表 */
export function getTrashPhotosets(params) {
  return request.get('/photosets/trash', { params })
}

/** 从回收站恢复 */
export function restorePhotoset(id) {
  return request.post(`/photosets/${id}/restore`)
}

/** 永久删除 */
export function permanentDeletePhotoset(id) {
  return request.delete(`/photosets/${id}/permanent`)
}

/** 下载套图 */
export function downloadPhotoset(id) {
  return request.get(`/photosets/${id}/download`)
}

/** 生成分享链接 */
export function getShareLink(id, source = 'web') {
  return request.get(`/photosets/${id}/share-link`, { params: { source } })
}

/** 验证分享链接（需携带 s/t/sig 签名参数） */
export function verifyShareLink(id, params) {
  return request.get(`/share/${id}`, { params })
}

// ============ 管理员套图审核 ============

/** 按状态查询套图（admin） */
export function getAdminPhotosets(params) {
  return request.get('/admin/photosets', { params })
}

/** 审核通过（admin） */
export function approvePhotoset(id) {
  return request.post(`/admin/photosets/${id}/approve`)
}

/** 审核拒绝（admin） */
export function rejectPhotoset(id, data) {
  return request.post(`/admin/photosets/${id}/reject`, data)
}
