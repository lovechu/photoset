/**
 * 用户状态管理
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getCurrentUser, login as apiLogin, register as apiRegister } from '@/api'

export const useUserStore = defineStore('user', () => {
  // 用户信息
  const user = ref(null)
  // Token
  const token = ref(localStorage.getItem('photoset_token') || '')

  // 初始化时尝试恢复登录态
  const init = async () => {
    const storedUser = localStorage.getItem('photoset_user')
    if (storedUser) {
      user.value = JSON.parse(storedUser)
    }
    if (token.value) {
      try {
        const res = await getCurrentUser()
        user.value = res.data.user
        localStorage.setItem('photoset_user', JSON.stringify(user.value))
      } catch (e) {
        // Token 失效，清除登录态
        logout()
      }
    }
  }

  // 计算属性：是否已登录
  const isLoggedIn = computed(() => !!token.value && !!user.value)

  // 计算属性：是否为创作者或管理员
  const isCreatorOrAdmin = computed(() => {
    if (!user.value) return false
    return ['creator', 'admin'].includes(user.value.role)
  })

  // 计算属性：是否为管理员
  const isAdmin = computed(() => user.value?.role === 'admin')

  // 计算属性：是否为会员
  const isMember = computed(() => {
    return ['member', 'creator', 'admin'].includes(user.value?.role)
  })

  // 登录
  const login = async (data) => {
    const res = await apiLogin(data)
    token.value = res.data.token
    user.value = res.data.user
    // 持久化
    localStorage.setItem('photoset_token', token.value)
    localStorage.setItem('photoset_user', JSON.stringify(user.value))
    return res
  }

  // 注册
  const register = async (data) => {
    const res = await apiRegister(data)
    user.value = res.data.user
    return res
  }

  // 登出
  const logout = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem('photoset_token')
    localStorage.removeItem('photoset_user')
  }

  // 更新用户信息
  const updateUser = (newUser) => {
    user.value = newUser
    localStorage.setItem('photoset_user', JSON.stringify(user.value))
  }

  return {
    user,
    token,
    isLoggedIn,
    isCreatorOrAdmin,
    isAdmin,
    isMember,
    init,
    login,
    register,
    logout,
    updateUser
  }
})
