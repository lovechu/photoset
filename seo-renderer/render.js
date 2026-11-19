/**
 * 无头浏览器渲染
 * - 复用单个 Browser 实例，每个任务独立 Page 用完关闭
 * - 拦截图片/字体/媒体二进制下载以提速（DOM 里的 <img src>/alt、OG 依然保留在 HTML 中）
 * - 访问本地前端容器（走内网，不经公网 CDN/防盗链）
 * - 读取页面内的 prerender-status-code meta，决定真实 HTTP 状态
 */
import puppeteer from 'puppeteer-core'

const FRONTEND_ORIGIN = process.env.FRONTEND_ORIGIN || 'http://frontend'
const CHROME_PATH = process.env.CHROME_PATH || '/usr/bin/chromium'
const RENDER_TIMEOUT = Number(process.env.RENDER_TIMEOUT || 30000)

let browser = null
let launching = null

async function getBrowser() {
  if (browser && browser.connected) return browser
  if (launching) return launching
  launching = puppeteer.launch({
    executablePath: CHROME_PATH,
    headless: 'new',
    args: [
      '--no-sandbox',
      '--disable-setuid-sandbox',
      '--disable-dev-shm-usage',
      '--disable-gpu',
      '--no-zygote',
      '--single-process'
    ]
  }).then(b => {
    browser = b
    launching = null
    b.on('disconnected', () => { browser = null })
    return b
  }).catch(e => {
    launching = null
    throw e
  })
  return launching
}

/**
 * 渲染一个站内路径，返回 { html, status, canonical }
 */
export async function renderPath(p) {
  const b = await getBrowser()
  const page = await b.newPage()
  try {
    await page.setViewport({ width: 1280, height: 800 })
    // 拦截媒体/字体下载提速；保留 HTML/DOM 中的真实 URL 与 alt
    await page.setRequestInterception(true)
    page.on('request', req => {
      const type = req.resourceType()
      if (['image', 'media', 'font'].includes(type)) {
        req.abort().catch(() => {})
      } else {
        req.continue().catch(() => {})
      }
    })

    const url = `${FRONTEND_ORIGIN}${p}`
    const resp = await page.goto(url, {
      waitUntil: 'networkidle2',
      timeout: RENDER_TIMEOUT
    })
    // 等 Vue 异步渲染完成：轮询 #app 内是否有内容，最多 8 秒
    try {
      await page.waitForFunction(
        () => {
          const app = document.querySelector('#app')
          return app && app.textContent.trim().length > 20
        },
        { timeout: 8000 }
      )
    } catch {
      // 空内容也返回，交给调用方按状态码判断
    }

    const html = await page.content()
    const meta = await page.evaluate(() => {
      const statusEl = document.querySelector('meta[name="prerender-status-code"]')
      const canonicalEl = document.querySelector('link[rel="canonical"]')
      return {
        status: statusEl ? Number(statusEl.getAttribute('content')) : 200,
        canonical: canonicalEl ? canonicalEl.getAttribute('href') : ''
      }
    })

    let status = meta.status
    if (!Number.isInteger(status) || status < 100 || status > 599) status = 200
    // 前端 SPA 对不存在资源没有明确状态时，用 HTTP 响应兜底判断
    if (status === 200 && resp && resp.status() >= 400) status = resp.status()

    return { html, status, canonical: meta.canonical }
  } finally {
    await page.close().catch(() => {})
  }
}

export async function closeBrowser() {
  if (browser) {
    await browser.close().catch(() => {})
    browser = null
  }
}
