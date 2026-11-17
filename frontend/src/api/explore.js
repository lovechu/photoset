/**
 * 探索/发现 API
 */
import request from '@/utils/request'

export function getExploreFeed(params) {
  return request.get('/explore/feed', { params })
}

export function getHotPhotosets(params) {
  return request.get('/explore/hot', { params })
}
