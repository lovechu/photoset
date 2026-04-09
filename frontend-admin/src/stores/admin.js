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
    user.value = res.data
    localStorage.setItem('photoset_admin_user', JSON.stringify(res.data))
    return res.data
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('photoset_admin_token')
    localStorage.removeItem('photoset_admin_user')
  }

  return { token, user, isLoggedIn, isAdmin, login, fetchUser, logout }
})
