/**
 * Vue Router 路由配置
 */
import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' }
  },
  {
    path: '/photoset/:id',
    name: 'PhotosetDetail',
    component: () => import('@/views/PhotosetDetail.vue'),
    meta: { title: '套图详情' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录', guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册', guest: true }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/ForgotPassword.vue'),
    meta: { title: '找回密码', guest: true }
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('@/views/ResetPassword.vue'),
    meta: { title: '重置密码', guest: true }
  },
  {
    path: '/create',
    name: 'CreatePhotoset',
    component: () => import('@/views/CreatePhotoset.vue'),
    meta: { title: '创建套图', requiresCreator: true }
  },
  {
    path: '/photoset/:id/edit',
    name: 'EditPhotoset',
    component: () => import('@/views/EditPhotoset.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/favorites',
    name: 'Favorites',
    component: () => import('@/views/Favorites.vue'),
    meta: { title: '我的收藏', requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { title: '个人中心', requiresAuth: true }
  },
  {
    path: '/membership',
    name: 'Membership',
    component: () => import('@/views/Membership.vue'),
    meta: { title: '会员订阅' }
  },
  {
    path: '/orders',
    name: 'Orders',
    component: () => import('@/views/Orders.vue'),
    meta: { title: '我的订单', requiresAuth: true }
  },
  {
    path: '/tags',
    name: 'Tags',
    component: () => import('@/views/Tags.vue'),
    meta: { title: '标签' }
  },
  {
    path: '/page/:pageType',
    name: 'StaticPage',
    component: () => import('@/views/StaticPage.vue'),
    meta: { title: '页面详情' }
  },
  {
    path: '/oauth/authorize',
    name: 'OAuthAuthorize',
    component: () => import('@/views/OAuthAuthorize.vue'),
    meta: { title: '授权确认' }
  },

  // ============ Phase 1 核心链路 ============
  { path: '/collections', name: 'CollectionList', component: () => import('@/views/CollectionList.vue'), meta: { title: '我的合集', requiresAuth: true } },
  { path: '/collection/:id', name: 'CollectionDetail', component: () => import('@/views/CollectionDetail.vue'), meta: { title: '合集详情', requiresAuth: true } },
  { path: '/history', name: 'History', component: () => import('@/views/History.vue'), meta: { title: '浏览历史', requiresAuth: true } },
  { path: '/notifications', name: 'Notifications', component: () => import('@/views/Notifications.vue'), meta: { title: '通知中心', requiresAuth: true } },
  { path: '/edit-profile', name: 'EditProfile', component: () => import('@/views/EditProfile.vue'), meta: { title: '编辑资料', requiresAuth: true } },
  { path: '/bind-email', name: 'BindEmail', component: () => import('@/views/BindEmail.vue'), meta: { title: '绑定邮箱', requiresAuth: true } },
  { path: '/login-history', name: 'LoginHistory', component: () => import('@/views/LoginHistory.vue'), meta: { title: '登录历史', requiresAuth: true } },
  { path: '/account-security', name: 'AccountSecurity', component: () => import('@/views/AccountSecurity.vue'), meta: { title: '账号安全', requiresAuth: true } },

  // ============ Phase 2 社区模块 ============
  { path: '/community', name: 'CommunityHome', component: () => import('@/views/community/CommunityHome.vue'), meta: { title: '社区' } },
  { path: '/community/post/:id', name: 'PostDetail', component: () => import('@/views/community/PostDetail.vue'), meta: { title: '帖子详情' } },
  { path: '/community/create', name: 'CreatePost', component: () => import('@/views/community/CreatePost.vue'), meta: { title: '发帖', requiresAuth: true } },
  { path: '/community/topic/:name', name: 'TopicPosts', component: () => import('@/views/community/TopicPosts.vue'), meta: { title: '话题' } },
  { path: '/community/tag/:name', name: 'TagPosts', component: () => import('@/views/community/TagPosts.vue'), meta: { title: '标签' } },
  { path: '/community/search', name: 'CommunitySearch', component: () => import('@/views/community/CommunitySearch.vue'), meta: { title: '社区搜索' } },
  { path: '/community/hot', name: 'HotPosts', component: () => import('@/views/community/HotPosts.vue'), meta: { title: '热门帖子' } },
  { path: '/community/center', name: 'CommunityCenter', component: () => import('@/views/community/CommunityCenter.vue'), meta: { title: '社区中心', requiresAuth: true } },
  { path: '/community/user/:id', name: 'CommunityUserProfile', component: () => import('@/views/community/UserProfile.vue'), meta: { title: '用户主页' } },
  { path: '/community/follow-list/:id', name: 'FollowList', component: () => import('@/views/community/FollowList.vue'), meta: { title: '关注列表', requiresAuth: true } },
  { path: '/community/my-posts', name: 'MyPosts', component: () => import('@/views/community/MyPosts.vue'), meta: { title: '我的帖子', requiresAuth: true } },
  { path: '/community/my-replies', name: 'MyReplies', component: () => import('@/views/community/MyReplies.vue'), meta: { title: '我的回复', requiresAuth: true } },
  { path: '/community/my-likes', name: 'MyLikes', component: () => import('@/views/community/MyLikes.vue'), meta: { title: '我的点赞', requiresAuth: true } },
  { path: '/community/favorites', name: 'CommunityFavorites', component: () => import('@/views/community/CommunityFavorites.vue'), meta: { title: '社区收藏', requiresAuth: true } },
  { path: '/community/drafts', name: 'Drafts', component: () => import('@/views/community/Drafts.vue'), meta: { title: '草稿箱', requiresAuth: true } },
  { path: '/community/conversations', name: 'Conversations', component: () => import('@/views/community/Conversations.vue'), meta: { title: '私信', requiresAuth: true } },
  { path: '/community/chat/:id', name: 'Chat', component: () => import('@/views/community/Chat.vue'), meta: { title: '聊天', requiresAuth: true } },
  { path: '/community/guidelines', name: 'CommunityGuidelines', component: () => import('@/views/community/CommunityGuidelines.vue'), meta: { title: '社区公约' } },

  // ============ Phase 3 账号与创作者 ============
  { path: '/user-level', name: 'UserLevel', component: () => import('@/views/UserLevel.vue'), meta: { title: '我的等级', requiresAuth: true } },
  { path: '/points-mall', name: 'PointsMall', component: () => import('@/views/PointsMall.vue'), meta: { title: '积分商城', requiresAuth: true } },
  { path: '/points-history', name: 'PointsHistory', component: () => import('@/views/PointsHistory.vue'), meta: { title: '积分记录', requiresAuth: true } },
  { path: '/leaderboard', name: 'Leaderboard', component: () => import('@/views/Leaderboard.vue'), meta: { title: '积分排行榜' } },
  { path: '/privacy-settings', name: 'PrivacySettings', component: () => import('@/views/PrivacySettings.vue'), meta: { title: '隐私设置', requiresAuth: true } },
  { path: '/block-list', name: 'BlockList', component: () => import('@/views/BlockList.vue'), meta: { title: '屏蔽列表', requiresAuth: true } },
  { path: '/deletion-status', name: 'DeletionStatus', component: () => import('@/views/DeletionStatus.vue'), meta: { title: '账号注销', requiresAuth: true } },
  { path: '/creator-dashboard', name: 'CreatorDashboard', component: () => import('@/views/CreatorDashboard.vue'), meta: { title: '创作者中心', requiresAuth: true } },
  { path: '/my-photosets', name: 'MyPhotosets', component: () => import('@/views/MyPhotosets.vue'), meta: { title: '我的套图', requiresAuth: true } },
  { path: '/trash', name: 'Trash', component: () => import('@/views/Trash.vue'), meta: { title: '回收站', requiresAuth: true } },
  { path: '/review-submissions', name: 'ReviewSubmissions', component: () => import('@/views/ReviewSubmissions.vue'), meta: { title: '投稿审核', requiresAuth: true } },

  // ============ Phase 4 长尾体验 ============
  { path: '/explore', name: 'Explore', component: () => import('@/views/Explore.vue'), meta: { title: '探索发现' } },
  { path: '/feedback', name: 'Feedback', component: () => import('@/views/Feedback.vue'), meta: { title: '意见反馈' } },
  { path: '/share/:id', name: 'ShareLanding', component: () => import('@/views/ShareLanding.vue'), meta: { title: '分享' } },

  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: { title: '页面不存在' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition
    return { top: 0 }
  }
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 更新页面标题
  document.title = to.meta.title ? `${to.meta.title}` : ''

  const userStore = useUserStore()

  // 需要登录的页面
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  // 需要 creator 或 admin 权限的页面
  if (to.meta.requiresCreator && !userStore.isCreatorOrAdmin) {
    next({ name: 'Home' })
    return
  }

  // 仅限未登录用户访问的页面（如登录、注册）
  if (to.meta.guest && userStore.isLoggedIn) {
    next({ name: 'Home' })
    return
  }

  next()
})

export default router
