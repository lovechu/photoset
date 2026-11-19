/**
 * 统一 SEO head 管理
 * 基于 @unhead/vue，所有公开页面通过本封装设置 title / description / canonical / OG / JSON-LD。
 * 站点信息从 useSiteStore 读取（站点名/logo），未加载时用默认值兜底，避免爬虫首屏为空。
 */
import { computed, unref } from 'vue'
import { useHead } from '@unhead/vue'
import { useSiteStore } from '@/stores/site'

export const SITE_URL = 'https://tt.cy.mk'
export const SITE_NAME = 'PhotoSet'
export const DEFAULT_DESC = 'PhotoSet 是一个摄影套图分享社区，发现高清优质套图、图集与原创摄影作品，浏览热门话题与用户分享。'
const DEFAULT_LOGO = `${SITE_URL}/favicon.svg`

/**
 * 设置页面 SEO 信息
 * @param {object|Ref|() => object} options 响应式对象/ref/getter
 *   - title: 页面标题（自动追加站点名）
 *   - description / path / image / type / noindex / jsonLd
 *   - statusCode: 内部软状态（404/410），写入 meta 供渲染服务读取，不改变浏览器行为
 */
export function useSeo(options) {
  const siteStore = useSiteStore()

  const head = computed(() => {
    const o = unref(typeof options === 'function' ? options() : options) || {}
    const siteName = (siteStore.siteName || '').trim() || SITE_NAME
    const title = o.title ? `${o.title} - ${siteName}` : siteName
    const description = o.description || DEFAULT_DESC
    const path = o.path || (typeof window !== 'undefined' ? window.location.pathname : '/')
    const canonical = path.startsWith('http') ? path : `${SITE_URL}${path}`
    const siteLogo = siteStore.siteLogo || DEFAULT_LOGO
    const image = o.image ? (o.image.startsWith('http') ? o.image : `${SITE_URL}${o.image}`) : siteLogo
    const robots = o.noindex ? 'noindex, nofollow' : 'index, follow'

    const meta = [
      { name: 'description', content: description },
      { name: 'robots', content: robots },
      { property: 'og:site_name', content: siteName },
      { property: 'og:type', content: o.type || 'website' },
      { property: 'og:title', content: title },
      { property: 'og:description', content: description },
      { property: 'og:url', content: canonical },
      { property: 'og:image', content: image },
      { name: 'twitter:card', content: 'summary_large_image' },
      { name: 'twitter:title', content: title },
      { name: 'twitter:description', content: description },
      { name: 'twitter:image', content: image }
    ]
    if (o.statusCode) {
      meta.push({ name: 'prerender-status-code', content: String(o.statusCode) })
    }

    const obj = {
      title,
      meta,
      link: [{ rel: 'canonical', href: canonical }]
    }
    if (o.jsonLd) {
      obj.script = (Array.isArray(o.jsonLd) ? o.jsonLd : [o.jsonLd]).map(data => ({
        type: 'application/ld+json',
        children: JSON.stringify(data)
      }))
    }
    return obj
  })

  useHead(head)
}

/** 站点级默认 head（在 App.vue 调用一次，补默认描述） */
export function useSiteSeo() {
  useHead({
    meta: [{ name: 'description', content: DEFAULT_DESC }]
  })
}
