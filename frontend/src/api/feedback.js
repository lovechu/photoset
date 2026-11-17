/**
 * 意见反馈 API
 */
import request from '@/utils/request'

export function submitFeedback(data) {
  return request.post('/feedback', data)
}

export function uploadFeedbackImage(file) {
  const formData = new FormData()
  formData.append('image', file)
  return request.post('/feedback/image', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
