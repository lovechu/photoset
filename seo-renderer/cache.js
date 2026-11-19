/**
 * HTML 文件缓存管理
 * - 路径安全映射（防目录穿越）
 * - 每个 URL 存 .html + .meta（HTTP 状态、canonical、生成时间、TTL、版本）
 * - 原子写（临时文件 + rename），避免读到半写文件
 * - 404/410 状态也持久化，缓存命中时返回真实状态
 */
import fs from 'node:fs/promises'
import path from 'node:path'
import crypto from 'node:crypto'

const CACHE_DIR = process.env.CACHE_DIR || '/var/cache/seo'
const DEFAULT_TTL = Number(process.env.CACHE_TTL || 3600) * 1000 // 1 小时
const ERROR_TTL = 5 * 60 * 1000 // 404/410 短缓存 5 分钟

// 允许渲染的公开路由白名单（规范化后精确匹配）
const ROUTE_PATTERNS = [
  /^\/$/,
  /^\/community$/,
  /^\/explore$/,
  /^\/tags$/,
  /^\/leaderboard$/,
  /^\/community\/hot$/,
  /^\/community\/post\/\d+$/,
  /^\/community\/user\/\d+$/,
  /^\/community\/topic\/[^/]+$/,
  /^\/community\/tag\/[^/]+$/,
  /^\/photoset\/\d+$/
]

/** 规范化路径：去 query/hash、解码、压缩斜杠、校验白名单。返回 null 表示非法。 */
export function normalizePath(raw) {
  if (typeof raw !== 'string' || !raw.startsWith('/')) return null
  let p
  try {
    p = decodeURIComponent(raw.split('?')[0].split('#')[0])
  } catch {
    return null
  }
  // 禁止目录穿越与协议
  if (p.includes('..') || p.includes('\\') || p.includes('://')) return null
  // 压缩重复斜杠
  p = p.replace(/\/{2,}/g, '/')
  if (p.length > 1 && p.endsWith('/')) p = p.slice(0, -1)
  if (!ROUTE_PATTERNS.some(re => re.test(p))) return null
  return p
}

function keyFor(p) {
  return crypto.createHash('sha256').update(p).digest('hex')
}

function filePaths(p) {
  const key = keyFor(p)
  const dir = path.join(CACHE_DIR, key.slice(0, 2))
  return {
    dir,
    html: path.join(dir, `${key}.html`),
    meta: path.join(dir, `${key}.meta`)
  }
}

/** 读缓存；未命中/过期/损坏返回 null */
export async function readCache(p) {
  const { html, meta } = filePaths(p)
  try {
    const [metaRaw, htmlRaw] = await Promise.all([
      fs.readFile(meta, 'utf8'),
      fs.readFile(html, 'utf8')
    ])
    const m = JSON.parse(metaRaw)
    const age = Date.now() - m.generatedAt
    const ttl = m.status === 404 || m.status === 410 ? ERROR_TTL : DEFAULT_TTL
    if (age < 0 || age > ttl) return null
    return { html: htmlRaw, ...m }
  } catch {
    return null
  }
}

/** 原子写 HTML + meta */
export async function writeCache(p, { html, status, canonical, version }) {
  const { dir, html: htmlPath, meta: metaPath } = filePaths(p)
  await fs.mkdir(dir, { recursive: true })
  const meta = JSON.stringify({
    status: status || 200,
    canonical: canonical || '',
    generatedAt: Date.now(),
    version: version || '1',
    path: p
  })
  const tmpHtml = `${htmlPath}.${process.pid}.tmp`
  const tmpMeta = `${metaPath}.${process.pid}.tmp`
  await fs.writeFile(tmpHtml, html, 'utf8')
  await fs.writeFile(tmpMeta, meta, 'utf8')
  await fs.rename(tmpHtml, htmlPath)
  await fs.rename(tmpMeta, metaPath)
}

/** 删除单个 URL 缓存 */
export async function deleteCache(p) {
  const { html, meta } = filePaths(p)
  await Promise.all([
    fs.unlink(html).catch(() => {}),
    fs.unlink(meta).catch(() => {})
  ])
}

/** 清空全部缓存 */
export async function clearAllCache() {
  await fs.rm(CACHE_DIR, { recursive: true, force: true }).catch(() => {})
  await fs.mkdir(CACHE_DIR, { recursive: true }).catch(() => {})
}
