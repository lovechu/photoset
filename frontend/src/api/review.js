/**
 * 套图评价 API
 */
import request from '@/utils/request'

export function listReviews(photosetId, params) {
  return request.get(`/photosets/${photosetId}/reviews`, { params })
}

export function getReviewSummary(photosetId) {
  return request.get(`/photosets/${photosetId}/reviews/summary`)
}

export function getMyReview(photosetId) {
  return request.get(`/photosets/${photosetId}/reviews/mine`)
}

export function createReview(photosetId, data) {
  return request.post(`/photosets/${photosetId}/reviews`, data)
}

export function updateReview(photosetId, reviewId, data) {
  return request.put(`/photosets/${photosetId}/reviews/${reviewId}`, data)
}

export function deleteReview(photosetId, reviewId) {
  return request.delete(`/photosets/${photosetId}/reviews/${reviewId}`)
}
