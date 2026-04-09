/**
 * axios 请求封装
 * 统一处理 token、响应拦截、错误提示
 */
import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const request = axios.create({
  baseURL: '/api',
  timeout: 15000
})

// 请求拦截器：自动携带 Token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('photoset_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器：统一处理错误
request.interceptors.response.use(
  (response) => {
    const res = response.data
    // 根据 code 判断业务是否成功
    if (res.code === 0) {
      return res
    }
    // 业务错误（code !== 0）
    ElMessage.error(res.message || '请求失败')
    return Promise.reject(new Error(res.message || '请求失败'))
  },
  (error) => {
    // HTTP 状态码错误处理
    const status = error.response?.status
    const message = error.response?.data?.message || '网络请求失败'

    switch (status) {
      case 401:
        // Token 无效或过期，清除本地 token，跳转登录
        localStorage.removeItem('photoset_token')
        localStorage.removeItem('photoset_user')
        ElMessage.error('登录已过期，请重新登录')
        router.push('/login')
        break
      case 403:
        ElMessage.error('权限不足')
        break
      case 404:
        ElMessage.error('资源不存在')
        break
      case 500:
        ElMessage.error('服务器错误，请稍后重试')
        break
      default:
        ElMessage.error(message)
    }

    return Promise.reject(error)
  }
)

export default request
