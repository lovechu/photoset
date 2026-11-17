/**
 * 社区帖子内容解析工具
 * 后端 Post 没有独立图片字段，图片/视频内嵌在 content 中（Markdown 或 HTML）
 */

// 提取内容中的图片地址
export function extractImages(content) {
  if (!content || typeof content !== 'string') return []
  const urls = []
  // Markdown: ![alt](url)
  const mdRe = /!\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)/g
  let m
  while ((m = mdRe.exec(content)) !== null) {
    if (m[2] && !urls.includes(m[2])) urls.push(m[2])
  }
  // HTML: <img src="url">
  const htmlRe = /<img[^>]+src=["']([^"']+)["'][^>]*>/gi
  while ((m = htmlRe.exec(content)) !== null) {
    if (m[1] && !urls.includes(m[1])) urls.push(m[1])
  }
  return urls
}

// 提取内容中的视频地址
export function extractVideos(content) {
  if (!content || typeof content !== 'string') return []
  const urls = []
  const videoRe = /<video[^>]+src=["']([^"']+)["'][^>]*>/gi
  let m
  while ((m = videoRe.exec(content)) !== null) {
    if (m[1] && !urls.includes(m[1])) urls.push(m[1])
  }
  const mp4Re = /(https?:\/\/[^\s)"'<>]+?\.(?:mp4|webm|ogg))(?:[?#][^\s)"'<>]*)?/gi
  while ((m = mp4Re.exec(content)) !== null) {
    if (m[1] && !urls.includes(m[1])) urls.push(m[1])
  }
  return urls
}

// 去掉内容里的图片/视频标记和 HTML 标签，得到纯文本摘要
export function toPlainText(content, maxLen = 200) {
  if (!content || typeof content !== 'string') return ''
  let text = content
    .replace(/!\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)/g, '')
    .replace(/<img[^>]*>/gi, '')
    .replace(/<video[^>]*>.*?<\/video>/gi, '')
    .replace(/<[^>]+>/g, '')
    .replace(/^#{1,6}\s+/gm, '')
    .replace(/[*_`~]/g, '')
    .replace(/\n{3,}/g, '\n\n')
    .trim()
  if (maxLen && text.length > maxLen) {
    text = text.slice(0, maxLen) + '...'
  }
  return text
}

// 把内容里的图片/视频标记渲染为简单可点击的 HTML（用于详情页）
export function renderContent(content) {
  if (!content || typeof content !== 'string') return ''
  let html = content
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    // Markdown 图片 -> img
    .replace(/!\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)/g, (_, alt, url) =>
      `<img src="${url}" alt="${alt}" class="post-content-img" />`)
    // 换行 -> br
    .replace(/\n/g, '<br/>')
  return html
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
