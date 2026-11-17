/**
 * 社区模块 API
 * 帖子/分类/热度/关注/拉黑/用户主页/私信/等级积分
 */
import request from '@/utils/request'

// ============ 帖子（公开） ============

export function getPosts(params) {
  return request.get('/community/posts', { params })
}

export function getPostDetail(id, params) {
  return request.get(`/community/posts/${id}`, { params })
}

export function getReplies(id, params) {
  return request.get(`/community/posts/${id}/replies`, { params })
}

export function searchPosts(params) {
  return request.get('/community/posts/search', { params })
}

export function getCommunityCategories() {
  return request.get('/community/categories')
}

export function getHotPosts(params) {
  return request.get('/community/hot', { params })
}

export function getCommunityUserProfile(id) {
  return request.get(`/community/users/${id}`)
}

export function getCommunityUserPosts(id, params) {
  return request.get(`/community/users/${id}/posts`, { params })
}

export function searchUsersForMention(params) {
  return request.get('/community/users/search', { params })
}

export function getShareStats(id) {
  return request.get(`/community/posts/${id}/shares`)
}

// ============ 标签 / 话题（公开） ============

export function getPopularTags(params) {
  return request.get('/community/tags/popular', { params })
}

export function searchTags(params) {
  return request.get('/community/tags/search', { params })
}

export function getPostsByTagName(name, params) {
  return request.get(`/community/tags/${encodeURIComponent(name)}/posts`, { params })
}

export function getHotTopics(params) {
  return request.get('/community/topics/hot', { params })
}

export function searchTopics(params) {
  return request.get('/community/topics/search', { params })
}

export function getPostsByTopicName(name, params) {
  return request.get(`/community/topics/${encodeURIComponent(name)}/posts`, { params })
}

// ============ 帖子（需登录） ============

export function createPost(data) {
  return request.post('/community/posts', data)
}

export function updatePost(id, data) {
  return request.put(`/community/posts/${id}`, data)
}

export function deletePost(id) {
  return request.delete(`/community/posts/${id}`)
}

export function createReply(id, data) {
  return request.post(`/community/posts/${id}/reply`, data)
}

export function updateReply(id, data) {
  return request.put(`/community/replies/${id}`, data)
}

export function deleteReply(id) {
  return request.delete(`/community/replies/${id}`)
}

export function togglePostLike(id) {
  return request.post(`/community/posts/${id}/like`)
}

export function toggleReplyLike(id) {
  return request.post(`/community/replies/${id}/like`)
}

export function recordShare(id, data) {
  return request.post(`/community/posts/${id}/share`, data)
}

export function togglePostFavorite(id) {
  return request.post(`/community/posts/${id}/favorite`)
}

export function checkPostFavorite(id) {
  return request.get(`/community/posts/${id}/favorite/check`)
}

export function getMyCommunityFavorites(params) {
  return request.get('/community/me/favorites', { params })
}

export function reportPost(id, data) {
  return request.post(`/community/posts/${id}/report`, data)
}

export function reportReply(id, data) {
  return request.post(`/community/replies/${id}/report`, data)
}

// ============ 我的内容 / 草稿 ============

export function getMyPosts(params) {
  return request.get('/community/me/posts', { params })
}

export function getMyReplies(params) {
  return request.get('/community/me/replies', { params })
}

export function getMyPoints() {
  return request.get('/community/me/points')
}

export function saveDraft(data) {
  return request.post('/community/drafts', data)
}

export function getMyDrafts(params) {
  return request.get('/community/me/drafts', { params })
}

export function deleteDraft(id) {
  return request.delete(`/community/drafts/${id}`)
}

// ============ 推荐 ============

export function getRecommendations(params) {
  return request.get('/community/recommendations', { params })
}

export function getUserInterests() {
  return request.get('/community/interests')
}

// ============ 关注 ============

export function followUser(id) {
  return request.post(`/community/users/${id}/follow`)
}

export function unfollowUser(id) {
  return request.delete(`/community/users/${id}/follow`)
}

export function checkFollowing(id) {
  return request.get(`/community/users/${id}/follow/check`)
}

export function batchCheckFollowing(data) {
  return request.post('/community/users/batch-follow-check', data)
}

export function getFollowingList(id, params) {
  return request.get(`/community/users/${id}/following`, { params })
}

export function getFollowerList(id, params) {
  return request.get(`/community/users/${id}/followers`, { params })
}

// ============ 拉黑 / 静音 ============

export function blockUser(id) {
  return request.post(`/community/users/${id}/block`)
}

export function muteUser(id, data) {
  return request.post(`/community/users/${id}/mute`, data)
}

export function unblockUser(id) {
  return request.delete(`/community/users/${id}/block`)
}

export function getBlockStatus(id) {
  return request.get(`/community/users/${id}/block/status`)
}

export function getBlockedUsers(params) {
  return request.get('/community/me/blocks', { params })
}

// ============ 私信 ============

export function getConversations(params) {
  return request.get('/community/messages/conversations', { params })
}

export function getConversation(userId, params) {
  return request.get(`/community/messages/conversations/${userId}`, { params })
}

export function sendMessage(data) {
  return request.post('/community/messages', data)
}

export function getMessageUnreadCount() {
  return request.get('/community/messages/unread-count')
}

export function markMessageRead(id) {
  return request.put(`/community/messages/${id}/read`)
}

export function markConversationRead(userId) {
  return request.put(`/community/messages/conversations/${userId}/read`)
}

export function deleteMessage(id) {
  return request.delete(`/community/messages/${id}`)
}

// ============ 等级 / 成就 / 积分（公开） ============

export function getLevelConfigs() {
  return request.get('/community/levels')
}

export function getUserLevelById(id) {
  return request.get(`/community/users/${id}/level`)
}

export function getUserAchievementsById(id) {
  return request.get(`/community/users/${id}/achievements`)
}

export function getPointsMallItems(params) {
  return request.get('/community/points-mall/items', { params })
}

export function getPointsMallCategories() {
  return request.get('/community/points-mall/categories')
}

// ============ 等级 / 成就 / 积分（需登录） ============

export function getMyLevel() {
  return request.get('/community/user/level')
}

export function getMyAchievements() {
  return request.get('/community/user/achievements')
}

export function exchangePointsItem(data) {
  return request.post('/community/points-mall/exchange', data)
}

export function getPointsExchangeHistory(params) {
  return request.get('/community/points-mall/history', { params })
}

export function getPointsLeaderboard(params) {
  return request.get('/community/points/leaderboard', { params })
}

// ============ 社区上传 ============

export function uploadCommunityImage(file) {
  const formData = new FormData()
  formData.append('image', file)
  return request.post('/community/upload/image', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export function uploadCommunityVideo(file) {
  const formData = new FormData()
  formData.append('video', file)
  return request.post('/community/upload/video', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
