/**
 * 用户状态管理
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getCurrentUser, login as apiLogin, register as apiRegister, updateProfile as apiUpdateProfile } from '@/api'
import { useNotificationStore } from './notification'

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
        if (res.data.user) {
          user.value = res.data.user
          localStorage.setItem('photoset_user', JSON.stringify(user.value))
        } else {
          // 没有用户信息，可能token已失效
          logout()
        }
      } catch (e) {
        // 网络或其他错误，保留当前状态，不自动登出
        console.warn('获取用户信息失败:', e.message)
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

  // 计算属性：会员过期时间
  const membershipExpires = computed(() => user.value?.membership_expires)

  // 计算属性：会员是否有效
  const isMembershipActive = computed(() => {
    if (!membershipExpires.value) return false
    return new Date(membershipExpires.value) > new Date()
  })

  // 计算属性：用户头像
  const avatar = computed(() => user.value?.avatar || '')

  // 计算属性：用户等级（从社区资料读取，缺省为 1）
  const level = computed(() => user.value?.level || 1)

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

  // 登出（同时重置通知状态）
  const logout = () => {
    token.value = ''
    user.value = null
    localStorage.removeItem('photoset_token')
    localStorage.removeItem('photoset_user')
    try {
      useNotificationStore().reset()
    } catch (e) { /* 通知 store 未初始化时忽略 */ }
  }

  // 更新用户信息
  const updateUser = (newUser) => {
    user.value = newUser
    localStorage.setItem('photoset_user', JSON.stringify(user.value))
  }

  // 更新用户资料（提交到后端后同步本地）
  const updateUserProfile = async (data) => {
    const res = await apiUpdateProfile(data)
    const updatedUser = res.data?.user || res.data
    if (updatedUser) updateUser(updatedUser)
    return res
  }

  // 支付订单后更新用户状态（接收后端返回的新 token + 用户信息）
  const afterPayment = (newToken, userData) => {
    token.value = newToken
    user.value = userData
    localStorage.setItem('photoset_token', newToken)
    localStorage.setItem('photoset_user', JSON.stringify(userData))
  }

  return {
    user,
    token,
    isLoggedIn,
    isCreatorOrAdmin,
    isAdmin,
    isMember,
    membershipExpires,
    isMembershipActive,
    avatar,
    level,
    init,
    login,
    register,
    logout,
    updateUser,
    updateUserProfile,
    afterPayment
  }
})
