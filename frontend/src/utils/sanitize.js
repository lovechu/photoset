// 轻量 HTML 净化：移除 script/style/iframe 等危险标签及所有 on* 事件属性和 javascript: 伪协议。
// 用于管理员编辑的 Markdown 页面内容，防止 v-html 触发存储型 XSS。
const FORBIDDEN_TAGS = ['SCRIPT', 'STYLE', 'IFRAME', 'OBJECT', 'EMBED', 'LINK', 'META', 'FORM', 'SVG', 'MATH']

export function sanitizeHtml(html) {
  if (!html || typeof html !== 'string') return ''
  const doc = new DOMParser().parseFromString(html, 'text/html')
  const walk = (node) => {
    for (const child of Array.from(node.children)) {
      const tag = child.tagName?.toUpperCase()
      if (FORBIDDEN_TAGS.includes(tag)) {
        child.remove()
        continue
      }
      for (const attr of Array.from(child.attributes || [])) {
        const name = attr.name.toLowerCase()
        const value = attr.value.trim().toLowerCase()
        if (name.startsWith('on') ||
            ((name === 'href' || name === 'src' || name === 'xlink:href') && value.startsWith('javascript:'))) {
          child.removeAttribute(attr.name)
        }
      }
      walk(child)
    }
  }
  walk(doc.body)
  return doc.body.innerHTML
}
