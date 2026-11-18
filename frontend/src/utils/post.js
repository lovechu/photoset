/**
 * 社区帖子内容解析工具
 * 后端 Post 没有独立图片字段，图片/视频内嵌在 content 中（Markdown 或 HTML）
 */

// 图片文件扩展名（含常见格式）
const IMG_EXT_RE = /\.(?:jpg|jpeg|png|gif|webp|avif|bmp)(?:[?#][^\s)"'<>]*)?$/i

// 裸图片 URL 正则：匹配 http(s) 或相对 /uploads 路径的图片地址
const BARE_IMG_RE = /(?:https?:\/\/[^\s)"'<>]+|\/uploads\/[^\s)"'<>]+?\.(?:jpg|jpeg|png|gif|webp|avif|bmp)(?:[?#][^\s)"'<>]*)?)/gi

// 提取内容中的图片地址
export function extractImages(content) {
  if (!content || typeof content !== 'string') return []
  const urls = []
  const add = (u) => { if (u && !urls.includes(u)) urls.push(u) }
  // Markdown: ![alt](url)
  const mdRe = /!\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)/g
  let m
  while ((m = mdRe.exec(content)) !== null) add(m[2])
  // HTML: <img src="url">
  const htmlRe = /<img[^>]+src=["']([^"']+)["'][^>]*>/gi
  while ((m = htmlRe.exec(content)) !== null) add(m[1])
  // App 发帖：视频封面图 [video_thumb]url[/video_thumb]
  const vtRe = /\[video_thumb\]\s*([^\s\[\]]+)\s*\[\/video_thumb\]/gi
  while ((m = vtRe.exec(content)) !== null) add(m[1])
  // App 发帖：图片裸 URL 直接写在正文里（如 https://域名/uploads/photos/xxx.jpg）
  // 用更严格的正则：必须是 /uploads/ 路径且以图片扩展名结尾，避免误抓普通链接
  const bareRe = /(?:https?:\/\/[^\s)"'<>]+\/uploads\/[^\s)"'<>]+|\/uploads\/[^\s)"'<>]+)/gi
  while ((m = bareRe.exec(content)) !== null) {
    if (IMG_EXT_RE.test(m[0])) add(m[0].replace(/[,，。]$/, ''))
  }
  return urls
}

// 提取内容中的视频地址
export function extractVideos(content) {
  if (!content || typeof content !== 'string') return []
  const urls = []
  const add = (u) => { if (u && !urls.includes(u)) urls.push(u) }
  const videoRe = /<video[^>]+src=["']([^"']+)["'][^>]*>/gi
  let m
  while ((m = videoRe.exec(content)) !== null) add(m[1])
  const mp4Re = /(https?:\/\/[^\s)"'<>]+?\.(?:mp4|webm|ogg))(?:[?#][^\s)"'<>]*)?/gi
  while ((m = mp4Re.exec(content)) !== null) add(m[1])
  // App 发帖：[video]url[/video] 或 [video]url
  const tagRe = /\[video\]\s*([^\s\[\]]+)\s*(?:\[\/video\])?/gi
  while ((m = tagRe.exec(content)) !== null) add(m[1])
  return urls
}

// 去掉内容里的图片/视频标记和 HTML 标签，得到纯文本摘要
export function toPlainText(content, maxLen = 200) {
  if (!content || typeof content !== 'string') return ''
  let text = content
    .replace(/!\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)/g, '')
    .replace(/\[video\][\s\S]*?(?:\[\/video\]|$)/gi, '')
    .replace(/\[video_thumb\][\s\S]*?(?:\[\/video_thumb\]|$)/gi, '')
    .replace(/<img[^>]*>/gi, '')
    .replace(/<video[^>]*>.*?<\/video>/gi, '')
    .replace(/<[^>]+>/g, '')
    // 去掉帖子正文里的裸图片/视频 URL（App 发帖格式），避免摘要被长链接刷屏
    .replace(/(?:https?:\/\/[^\s)"'<>]+\/uploads\/[^\s)"'<>]+|(?:\.\/|\/)?uploads\/[^\s)"'<>]+)/gi, '')
    // 去掉 App 表情标记 [微笑]/[笑哭] 及话题标记残留，避免正文漏出方括号文字
    .replace(/\[[\u4e00-\u9fa5A-Za-z0-9]{1,8}\]/g, '')
    .replace(/^#{1,6}\s+/gm, '')
    .replace(/[*_`~]/g, '')
    .replace(/\n{3,}/g, '\n\n')
    .trim()
  if (maxLen && text.length > maxLen) {
    text = text.slice(0, maxLen) + '...'
  }
  return text
}

// 转义 HTML 属性中的特殊字符，防止 Markdown alt/url 注入事件处理器
function escapeAttr(s) {
  return String(s)
    .replace(/&/g, '&amp;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
}

// 仅允许安全的图片协议（http/https 或相对路径），阻断 javascript: 等伪协议
function safeImageUrl(url) {
  const u = String(url).trim()
  if (/^(https?:\/\/|\/|\.\/|\.\.\/)/i.test(u)) return u
  return ''
}

// 把内容里的图片/视频标记渲染为 HTML（用于详情页）
// 兼容 Markdown ![alt](url)、App 裸图片 URL、[video]url[/video]、#话题#、@提及
export function renderContent(content) {
  if (!content || typeof content !== 'string') return ''
  // 先整体转义，杜绝任何标签/事件注入
  let html = content
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')

  // 用占位符暂存已生成的标签，避免后续正则再次匹配到标签内部的 URL
  const slots = []
  const stash = (tag) => {
    const token = `\u0000SLOT${slots.length}\u0000`
    slots.push(tag)
    return token
  }

  // Markdown 图片 ![alt](url)
  html = html.replace(/!\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)/g, (_, alt, url) => {
    const safeUrl = safeImageUrl(url)
    if (!safeUrl) return ''
    return stash(`<img src="${safeUrl}" alt="${attr(alt)}" class="post-content-img" />`)
  })

  // App 格式：[video_thumb]url[/video_thumb] -> 封面图（可能出现在 video 之前或之后）
  html = html.replace(/\[video_thumb\]\s*([^\s\[\]]+)\s*\[\/video_thumb\]/gi, (_, url) => {
    const u = String(url).replace(/&amp;/g, '&')
    if (!/^(https?:\/\/|\/)/i.test(u)) return ''
    return stash(`<img src="${u}" class="post-content-img" />`)
  })

  // App 格式：[video]url[/video]（或没有闭合标签）-> video 元素
  html = html.replace(/\[video\]\s*([^\s\[\]]+)\s*(?:\[\/video\])?/gi, (_, url) => {
    const u = String(url).replace(/&amp;/g, '&')
    if (!/^(https?:\/\/|\/)/i.test(u)) return ''
    return stash(`<video src="${u}" controls class="post-content-video" style="max-width:100%;border-radius:4px;margin:12px 0;"></video>`)
  })

  // App 格式：正文中的裸图片 URL（/uploads/... 且图片扩展名）-> img
  html = html.replace(/(?:https?:\/\/[^\s)"'<>]+\/uploads\/[^\s)"'<>]+|\/uploads\/[^\s)"'<>]+)/gi, (url) => {
    const u = url.replace(/[,，。]$/, '').replace(/&amp;/g, '&')
    if (!IMG_EXT_RE.test(u)) return url
    return stash(`<img src="${u}" class="post-content-img" />`)
  })

  // #话题# -> 高亮 span
  html = html.replace(/#([^#\s]{1,30})#/g, (_, topic) =>
    `<span class="post-topic">#${attr(topic)}#</span>`)

  // App 表情标记 [微笑] -> emoji，未收录的常见表情做通用映射
  html = html.replace(/\[([\u4e00-\u9fa5]{1,6})\]/g, (full, name) => {
    const map = {
      微笑: '😊', 笑: '😄', 大笑: '😆', 可爱: '🥰', 色: '😍', 害羞: '😊',
      哭: '😢', 笑哭: '😂', 大哭: '😭', 生气: '😠', 愤怒: '😡', 惊讶: '😲',
      赞: '👍', 厉害: '👍', OK: '👌', 加油: '💪', 爱心: '❤️', 心: '❤️',
      花: '🌸', 玫瑰: '🌹', 太阳: '☀️', 月亮: '🌙', 星星: '⭐', 火: '🔥',
      思考: '🤔', 困: '😴', 睡: '😴', 吐: '🤢', 晕: '😵', 酷: '😎'
    }
    return map[name] || full
  })

  // 换行 -> br
  html = html.replace(/\n/g, '<br/>')

  // 回填占位符
  html = html.replace(/\u0000SLOT(\d+)\u0000/g, (_, i) => slots[Number(i)] || '')
  return html
}

// 属性值已被整体转义，这里只还原 &amp;，避免 URL 查询参数被二次转义
function attr(s) {
  return String(s).replace(/&amp;/g, '&')
}

export function formatTime(str) {
  if (!str) return ''
  const d = new Date(str)
  if (isNaN(d.getTime())) return str
  const diff = Date.now() - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins}分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}天前`
  return d.toLocaleDateString('zh-CN')
}
