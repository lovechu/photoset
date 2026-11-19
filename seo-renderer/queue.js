/**
 * 渲染任务队列
 * - 并发限制（默认 1，避免 Chromium 打满内存）
 * - 同一 URL 任务合并（防缓存击穿：并发请求只渲染一次）
 * - 失败重试（最多 2 次）
 */
import { renderPath } from './render.js'
import { writeCache } from './cache.js'

const CONCURRENCY = Math.max(1, Number(process.env.CONCURRENCY || 1))
const MAX_RETRY = 2

let active = 0
const waiters = []
// 同 URL 合并：path -> Promise
const inflight = new Map()

function acquire() {
  return new Promise(resolve => {
    if (active < CONCURRENCY) {
      active++
      resolve()
    } else {
      waiters.push(resolve)
    }
  })
}

function release() {
  if (waiters.length) {
    const next = waiters.shift()
    next() // active 数量不变，交给下一个任务
  } else {
    active--
  }
}

async function doRender(p) {
  let lastErr
  for (let attempt = 0; attempt <= MAX_RETRY; attempt++) {
    try {
      const result = await renderPath(p)
      // 5xx 与 404/410 都写缓存（404/410 短 TTL），避免反复渲染
      await writeCache(p, {
        html: result.html,
        status: result.status,
        canonical: result.canonical,
        version: '1'
      })
      return result
    } catch (e) {
      lastErr = e
      if (attempt < MAX_RETRY) await new Promise(r => setTimeout(r, 500 * (attempt + 1)))
    }
  }
  throw lastErr
}

/** 渲染并写缓存（带并发控制 + 同 URL 合并） */
export function renderAndCache(p) {
  if (inflight.has(p)) return inflight.get(p)
  const task = (async () => {
    await acquire()
    try {
      return await doRender(p)
    } finally {
      release()
      inflight.delete(p)
    }
  })()
  inflight.set(p, task)
  return task
}

export function queueStats() {
  return { active, waiting: waiters.length, inflight: inflight.size }
}
