import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, getCurrentUser } from '@/api'

export const useAdminStore = defineStore('admin', () => {
  const token = ref(localStorage.getItem('photoset_admin_token') || '')
  const user = ref(JSON.parse(localStorage.getItem('photoset_admin_user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(credentials) {
    const res = await loginApi(credentials)
    const data = res.data
    token.value = data.token
    user.value = data.user
    localStorage.setItem('photoset_admin_token', data.token)
    localStorage.setItem('photoset_admin_user', JSON.stringify(data.user))
    return data
  }

  async function fetchUser() {
    const res = await getCurrentUser()
    // /auth/me 使用 OptionalAuth：token 失效时返回 200 + user:null，不会抛 401
    const profile = res.data?.user
    if (!profile) {
      logout()
      throw new Error('登录已失效')
    }
    user.value = profile
    localStorage.setItem('photoset_admin_user', JSON.stringify(profile))
    return profile
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('photoset_admin_token')
    localStorage.removeItem('photoset_admin_user')
  }

  return { token, user, isLoggedIn, isAdmin, login, fetchUser, logout }
})
