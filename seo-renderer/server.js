/**
 * PhotoSet SEO 渲染服务入口
 * - HTTP 服务：Nginx 爬虫分流命中/未命中统一走这里；/render /invalidate 供内部调用
 * - Redis Stream 消费者：消费 Go 后端发布的 render/invalidate 事件（可靠、不丢消息）
 * - 全部接口仅内部 Docker 网络可达 + Token 认证
 */
import http from 'node:http'
import Redis from 'ioredis'
import { normalizePath, readCache, deleteCache, clearAllCache } from './cache.js'
import { renderAndCache, queueStats } from './queue.js'
import { closeBrowser } from './render.js'

const PORT = Number(process.env.PORT || 3000)
const INTERNAL_TOKEN = process.env.SEO_INTERNAL_TOKEN || ''
const REDIS_URL = process.env.REDIS_URL || 'redis://redis:6379'
const STREAM = 'photoset:seo:render'
const GROUP = 'photoset-seo-group'
const CONSUMER = `renderer-${process.pid}`
const VERSION = '1'

const redis = new Redis(REDIS_URL, { maxRetriesPerRequest: null, lazyConnect: false })

function json(res, status, data) {
  res.writeHead(status, { 'Content-Type': 'application/json; charset=utf-8' })
  res.end(JSON.stringify(data))
}

function authorized(req) {
  if (!INTERNAL_TOKEN) return true // 未配置 token 时仅靠内网隔离
  const token = req.headers['x-seo-token'] || ''
  return token === INTERNAL_TOKEN
}

/**
 * 核心：按路径返回渲染 HTML（带真实状态码）
 * 命中缓存直接返回；未命中则渲染+写缓存后返回；渲染失败返回 503 且不写长缓存
 */
async function servePath(p, res) {
  const cached = await readCache(p)
  if (cached) {
    res.writeHead(cached.status || 200, {
      'Content-Type': 'text/html; charset=utf-8',
      'X-SEO-Cache': 'HIT',
      'X-SEO-Status': String(cached.status || 200)
    })
    res.end(cached.html)
    return
  }

  try {
    const result = await renderAndCache(p)
    res.writeHead(result.status || 200, {
      'Content-Type': 'text/html; charset=utf-8',
      'X-SEO-Cache': 'MISS',
      'X-SEO-Status': String(result.status || 200)
    })
    res.end(result.html)
  } catch (e) {
    console.error('[SEO] 渲染失败:', p, e.message)
    // 渲染服务不可用：返回 503，不把错误写成长缓存
    res.writeHead(503, { 'Content-Type': 'text/plain; charset=utf-8', 'Retry-After': '30' })
    res.end('Service Unavailable')
  }
}

const server = http.createServer(async (req, res) => {
  try {
    const url = new URL(req.url, `http://${req.headers.host || 'localhost'}`)

    // 健康检查（内网）
    if (url.pathname === '/health') {
      return json(res, 200, { ok: true, queue: queueStats() })
    }

    // 所有业务接口需内部 token
    if (!authorized(req)) return json(res, 403, { error: 'forbidden' })

    // GET /seo?path=xxx —— Nginx 爬虫分流入口
    if (req.method === 'GET' && url.pathname === '/seo') {
      const p = normalizePath(url.searchParams.get('path') || '/')
      if (!p) return json(res, 400, { error: 'invalid path' })
      return await servePath(p, res)
    }

    // POST /render —— 主动渲染（内容发布事件/手动触发）
    if (req.method === 'POST' && url.pathname === '/render') {
      const body = await readBody(req)
      const p = normalizePath(body.path)
      if (!p) return json(res, 400, { error: 'invalid path' })
      // 主动渲染时先删旧缓存，强制刷新
      await deleteCache(p)
      try {
        const result = await renderAndCache(p)
        return json(res, 200, { ok: true, path: p, status: result.status })
      } catch (e) {
        return json(res, 503, { ok: false, path: p, error: e.message })
      }
    }

    // POST /invalidate —— 清除缓存（内容删除事件/手动触发）
    if (req.method === 'POST' && url.pathname === '/invalidate') {
      const body = await readBody(req)
      if (body.all) {
        await clearAllCache()
        return json(res, 200, { ok: true, all: true })
      }
      const p = normalizePath(body.path)
      if (!p) return json(res, 400, { error: 'invalid path' })
      await deleteCache(p)
      return json(res, 200, { ok: true, path: p })
    }

    return json(res, 404, { error: 'not found' })
  } catch (e) {
    console.error('[SEO] 请求处理异常:', e)
    return json(res, 500, { error: 'internal error' })
  }
})

function readBody(req) {
  return new Promise((resolve, reject) => {
    let data = ''
    req.on('data', c => {
      data += c
      if (data.length > 10240) req.destroy() // 限制请求体 10KB
    })
    req.on('end', () => {
      try {
        resolve(data ? JSON.parse(data) : {})
      } catch {
        reject(new Error('invalid json'))
      }
    })
    req.on('error', reject)
  })
}

// ---- Redis Stream 消费者（可靠事件，不丢消息）----

async function initConsumer() {
  try {
    await redis.xgroup('CREATE', STREAM, GROUP, '0', 'MKSTREAM')
  } catch (e) {
    // BUSYGROUP 表示组已存在，正常
    if (!String(e.message).includes('BUSYGROUP')) {
      console.error('[SEO] 创建消费者组失败:', e.message)
    }
  }
  console.log('[SEO] Stream 消费者已启动')
  while (true) {
    try {
      const res = await redis.xreadgroup(
        'GROUP', GROUP, CONSUMER,
        'COUNT', 1, 'BLOCK', 5000,
        'STREAMS', STREAM, '>'
      )
      if (!res) continue
      for (const [, messages] of res) {
        for (const [id, fields] of messages) {
          try {
            const action = fields.action
            const rawPath = fields.path || '/'
            const p = normalizePath(rawPath)
            if (p) {
              if (action === 'invalidate') {
                await deleteCache(p)
              } else {
                await deleteCache(p)
                await renderAndCache(p)
              }
              console.log(`[SEO] 处理事件 ${id} ${action} ${p}`)
            }
            await redis.xack(STREAM, GROUP, id)
          } catch (e) {
            console.error(`[SEO] 处理事件 ${id} 失败:`, e.message)
            // 不 ACK，下次重投递（实现至少一次投递）
          }
        }
      }
    } catch (e) {
      console.error('[SEO] Stream 读取异常:', e.message)
      await new Promise(r => setTimeout(r, 2000))
    }
  }
}

// 启动时先清空旧消费者组中未确认的pending（简单重启策略：本服务无长任务，直接接管）
async function bootstrap() {
  await initConsumer()
}

server.listen(PORT, '0.0.0.0', () => {
  console.log(`[SEO] 渲染服务监听 0.0.0.0:${PORT}`)
  bootstrap()
})

// 优雅退出
async function shutdown() {
  await closeBrowser()
  await redis.quit().catch(() => {})
  server.close(() => process.exit(0))
}
process.on('SIGTERM', shutdown)
process.on('SIGINT', shutdown)
