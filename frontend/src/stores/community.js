/**
 * 社区状态管理
 * 缓存社区分类、当前圈子、发帖草稿
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCommunityCategories } from '@/api'

export const useCommunityStore = defineStore('community', () => {
  const categories = ref([])
  const currentCategory = ref('')
  const loaded = ref(false)

  async function fetchCategories() {
    if (loaded.value && categories.value.length) return categories.value
    try {
      const res = await getCommunityCategories()
      const data = extractData(res)
      // 后端返回 { data: { categories: [...] } }
      if (data && Array.isArray(data.categories)) {
        categories.value = data.categories
      } else if (data && Array.isArray(data.list)) {
        categories.value = data.list
      } else if (Array.isArray(data)) {
        categories.value = data
      } else if (Array.isArray(res)) {
        categories.value = res
      }
      loaded.value = true
      return categories.value
    } catch (e) {
      console.warn('获取社区分类失败:', e.message)
      return []
    }
  }

  function setCurrentCategory(slug) {
    currentCategory.value = slug || ''
  }

  function reset() {
    categories.value = []
    currentCategory.value = ''
    loaded.value = false
  }

  return {
    categories,
    currentCategory,
    loaded,
    fetchCategories,
    setCurrentCategory,
    reset
  }
})

function extractData(res) {
  if (res && res.data && typeof res.data === 'object') return res.data
  if (res && typeof res === 'object' && !res.code) return res
  return null
}
