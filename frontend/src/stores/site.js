/**
 * 站点设置状态管理
 * 从后端 /api/settings 获取，带容错默认值
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getSiteSettings } from '@/api'

export const useSiteStore = defineStore('site', () => {
  // 站点设置原始数据
  const settings = ref({})
  const loaded = ref(false)

  // 计算属性：带默认值的快捷访问
  const siteName = computed(() => settings.value.site_title || 'PhotoSet')
  const siteDescription = computed(() => settings.value.site_description || '专业摄影套图浏览平台，发现美好视觉作品')
  const seoTitle = computed(() => settings.value.seo_title || '')
  const seoKeywords = computed(() => settings.value.seo_keywords || '')
  const seoDescription = computed(() => settings.value.seo_description || '')
  const copyrightYear = computed(() => settings.value.copyright_year || new Date().getFullYear().toString())
  const icpLicense = computed(() => settings.value.icp_license || '')

  // 获取所有设置
  const fetchSettings = async () => {
    try {
      const res = await getSiteSettings()
      // request 拦截器在 code===0 时返回 {code,message,data}
      // 兼容直接返回 data 对象的情况
      let data = null
      if (res && res.data && typeof res.data === 'object') {
        data = res.data
      } else if (res && typeof res === 'object' && !res.code) {
        data = res
      }
      if (data) {
        settings.value = data
      }
    } catch (e) {
      // 后端接口不可用时静默降级，使用默认值
      console.warn('获取站点设置失败，使用默认值:', e.message)
    } finally {
      loaded.value = true
    }
  }

  // 根据 key 获取单个设置值
  const get = (key, defaultValue = '') => {
    return settings.value[key] || defaultValue
  }

  return {
    settings,
    loaded,
    siteName,
    siteDescription,
    seoTitle,
    seoKeywords,
    seoDescription,
    copyrightYear,
    icpLicense,
    fetchSettings,
    get
  }
})
