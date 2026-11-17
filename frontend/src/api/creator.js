/**
 * 创作者数据统计 API
 */
import request from '@/utils/request'

export function getCreatorStats() {
  return request.get('/creator/stats')
}

export function getCreatorDailyStats(params) {
  return request.get('/creator/stats/daily', { params })
}

export function getCreatorPhotosetStats(params) {
  return request.get('/creator/stats/photosets', { params })
}
