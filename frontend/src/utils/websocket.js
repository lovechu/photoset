/**
 * WebSocket 单例连接
 * 集中处理 token 鉴权、心跳、指数退避重连，通知与私信共用同一连接
 * 后端消息格式: { type: string, payload: any }
 */
import { ElMessage } from 'element-plus'

let socket = null
let reconnectTimer = null
let heartbeatTimer = null
let reconnectAttempts = 0
let manualClosed = false

const MAX_RECONNECT_DELAY = 30000
const HEARTBEAT_INTERVAL = 25000

// 消息订阅者：Map<type, Set<handler>>
const subscribers = new Map()
// 全局消息监听（接收所有消息）
const globalHandlers = new Set()
// 连接状态回调
const statusHandlers = new Set()

function currentToken() {
  return localStorage.getItem('photoset_token') || ''
}

function wsBaseURL() {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  // 开发环境下 Vite 代理 /api 到后端；生产与页面同源
  return `${protocol}//${window.location.host}/api/community/ws`
}

function notifyStatus(status) {
  statusHandlers.forEach(handler => {
    try { handler(status) } catch (e) { console.error('WS 状态回调异常:', e) }
  })
}

function dispatch(type, payload, raw) {
  const handlers = subscribers.get(type)
  if (handlers) {
    handlers.forEach(handler => {
      try { handler(payload, raw) } catch (e) { console.error(`WS [${type}] 处理器异常:`, e) }
    })
  }
  globalHandlers.forEach(handler => {
    try { handler(type, payload, raw) } catch (e) { console.error('WS 全局处理器异常:', e) }
  })
}

function scheduleReconnect() {
  if (manualClosed || reconnectTimer) return
  reconnectAttempts++
  const delay = Math.min(1000 * Math.pow(2, reconnectAttempts - 1), MAX_RECONNECT_DELAY)
  console.warn(`[WS] 连接断开，${delay}ms 后重连（第 ${reconnectAttempts} 次）`)
  reconnectTimer = setTimeout(() => {
    reconnectTimer = null
    connectWebSocket()
  }, delay)
}

function startHeartbeat() {
  stopHeartbeat()
  heartbeatTimer = setInterval(() => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ type: 'ping', payload: {} }))
    }
  }, HEARTBEAT_INTERVAL)
}

function stopHeartbeat() {
  if (heartbeatTimer) {
    clearInterval(heartbeatTimer)
    heartbeatTimer = null
  }
}

/** 建立 WebSocket 连接（未登录时不建立） */
export function connectWebSocket() {
  const token = currentToken()
  if (!token) {
    console.info('[WS] 未登录，跳过连接')
    return
  }
  if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
    return
  }

  manualClosed = false
  const url = `${wsBaseURL()}?token=${encodeURIComponent(token)}`
  try {
    socket = new WebSocket(url)
  } catch (e) {
    console.error('[WS] 创建连接失败:', e)
    scheduleReconnect()
    return
  }

  socket.onopen = () => {
    reconnectAttempts = 0
    notifyStatus('open')
    startHeartbeat()
  }

  socket.onmessage = (event) => {
    let msg
    try {
      msg = JSON.parse(event.data)
    } catch (e) {
      console.warn('[WS] 收到非 JSON 消息:', event.data)
      return
    }
    if (!msg || !msg.type) return
    dispatch(msg.type, msg.payload, msg)
  }

  socket.onerror = (error) => {
    console.error('[WS] 连接错误:', error)
  }

  socket.onclose = (event) => {
    stopHeartbeat()
    notifyStatus('close')
    socket = null
    if (!manualClosed && event.code !== 1000) {
      scheduleReconnect()
    }
  }
}

/** 主动断开连接（登出时调用，不触发重连） */
export function disconnectWebSocket() {
  manualClosed = true
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  reconnectAttempts = 0
  stopHeartbeat()
  if (socket) {
    try {
      socket.close(1000, 'client logout')
    } catch (e) { /* 忽略关闭异常 */ }
    socket = null
  }
  notifyStatus('close')
}

/** 订阅某种类型的消息，返回取消订阅函数 */
export function onMessage(type, handler) {
  if (!subscribers.has(type)) subscribers.set(type, new Set())
  subscribers.get(type).add(handler)
  return () => {
    const handlers = subscribers.get(type)
    if (handlers) handlers.delete(handler)
  }
}

/** 订阅所有消息，返回取消订阅函数 */
export function onAnyMessage(handler) {
  globalHandlers.add(handler)
  return () => globalHandlers.delete(handler)
}

/** 订阅连接状态变化，返回取消订阅函数 */
export function onStatusChange(handler) {
  statusHandlers.add(handler)
  return () => statusHandlers.delete(handler)
}

/** 获取当前连接状态 */
export function getWebSocketState() {
  if (!socket) return 'closed'
  switch (socket.readyState) {
    case WebSocket.CONNECTING: return 'connecting'
    case WebSocket.OPEN: return 'open'
    case WebSocket.CLOSING: return 'closing'
    default: return 'closed'
  }
}

/** 发送消息到后端 */
export function sendWebSocketMessage(type, payload) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({ type, payload }))
    return true
  }
  console.warn('[WS] 连接未就绪，消息未发送')
  return false
}
