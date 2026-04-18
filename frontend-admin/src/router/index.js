import { createRouter, createWebHistory } from 'vue-router'
import { useAdminStore } from '@/stores/admin'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    component: () => import('@/layout/AdminLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '数据看板', icon: 'DataAnalysis' }
      },
      {
        path: 'review',
        name: 'ContentReview',
        component: () => import('@/views/ContentReview.vue'),
        meta: { title: '内容审核', icon: 'PictureFilled' }
      },
      {
        path: 'users',
        name: 'UserManage',
        component: () => import('@/views/UserManage.vue'),
        meta: { title: '用户管理', icon: 'User' }
      },
      {
        path: 'orders',
        name: 'OrderManage',
        component: () => import('@/views/OrderManage.vue'),
        meta: { title: '订单管理', icon: 'ShoppingBag' }
      },
      {
        path: 'tags',
        name: 'TagManage',
        component: () => import('@/views/TagManage.vue'),
        meta: { title: '标签管理', icon: 'PriceTag' }
      },
      {
        path: 'categories',
        name: 'CategoryManage',
        component: () => import('@/views/CategoryManage.vue'),
        meta: { title: '分类管理', icon: 'FolderOpened' }
      },
      {
        path: 'pages',
        name: 'Pages',
        component: () => import('@/views/Pages.vue'),
        meta: { title: '页面管理', icon: 'Document' }
      },
      {
        path: 'settings',
        name: 'SiteSettings',
        component: () => import('@/views/SiteSettings.vue'),
        meta: { title: '站点设置', icon: 'Setting' }
      },
      {
        path: 'developer',
        name: 'DeveloperCenter',
        component: () => import('@/views/DeveloperCenter.vue'),
        meta: { title: '开发者中心', icon: 'Code' }
      },
      {
        path: 'photoset/:id/edit',
        name: 'AdminEditPhotoset',
        component: () => import('@/views/EditPhotoset.vue')
      },
      {
        path: 'logs',
        name: 'AdminLogs',
        component: () => import('@/views/AdminLogs.vue'),
        meta: { title: '操作日志', icon: 'Notebook' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  const store = useAdminStore()

  if (store.token && !store.user) {
    try {
      await store.fetchUser()
    } catch {
      store.logout()
      return next('/login')
    }
  }

  if (to.meta.public) {
    if (store.isLoggedIn) return next('/dashboard')
    return next()
  }

  if (!store.isLoggedIn) return next('/login')
  if (!store.isAdmin) {
    alert('权限不足，仅管理员可访问后台')
    store.logout()
    return next('/login')
  }

  next()
})

export default router
