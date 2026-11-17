/**
 * 套图评论 API
 */
import request from '@/utils/request'

export function listComments(photosetId, params) {
  return request.get(`/photosets/${photosetId}/comments`, { params })
}

export function createComment(photosetId, data) {
  return request.post(`/photosets/${photosetId}/comments`, data)
}

export function getCommentReplies(photosetId, commentId, params) {
  return request.get(`/photosets/${photosetId}/comments/${commentId}/replies`, { params })
}

export function deleteComment(photosetId, commentId) {
  return request.delete(`/photosets/${photosetId}/comments/${commentId}`)
}

export function toggleCommentLike(photosetId, commentId) {
  return request.post(`/photosets/${photosetId}/comments/${commentId}/like`)
}
