/**
 * 合集 API
 */
import request from '@/utils/request'

export function createCollection(data) {
  return request.post('/collections', data)
}

export function listCollections(params) {
  return request.get('/collections', { params })
}

export function getCollection(id) {
  return request.get(`/collections/${id}`)
}

export function updateCollection(id, data) {
  return request.put(`/collections/${id}`, data)
}

export function deleteCollection(id) {
  return request.delete(`/collections/${id}`)
}

export function addCollectionItem(id, data) {
  return request.post(`/collections/${id}/items`, data)
}

export function removeCollectionItem(id, photosetId) {
  return request.delete(`/collections/${id}/items/${photosetId}`)
}

export function batchAddCollectionItems(id, data) {
  return request.post(`/collections/${id}/items/batch`, data)
}

export function getCollectionsForPhotoset(photosetId) {
  return request.get(`/collections/by-photoset/${photosetId}`)
}
