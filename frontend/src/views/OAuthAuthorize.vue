<template>
  <div class="oauth-authorize">
    <div class="authorize-card">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="error-state">
        <div class="error-icon">!</div>
        <h2>授权失败</h2>
        <p>{{ error }}</p>
        <button class="btn btn-primary" @click="goHome">返回首页</button>
      </div>

      <!-- 授权页面 -->
      <div v-else-if="clientInfo" class="authorize-content">
        <div class="app-info">
          <img v-if="clientInfo.client?.logo_url" :src="clientInfo.client.logo_url" alt="应用Logo" class="app-logo" />
          <div v-else class="app-logo-placeholder">A</div>
          <h1 class="app-name">{{ clientInfo.client?.name }}</h1>
          <p v-if="clientInfo.client?.description" class="app-desc">{{ clientInfo.client.description }}</p>
        </div>

        <div class="scopes-section">
          <p class="scopes-title">该应用请求以下权限：</p>
          <ul class="scopes-list">
            <li v-for="scope in clientInfo.scopes" :key="scope" class="scope-item">
              <span class="scope-icon">✓</span>
              <span class="scope-name">{{ getScopeName(scope) }}</span>
            </li>
          </ul>
        </div>

        <div v-if="!isLoggedIn" class="login-prompt">
          <p>请先登录后再进行授权</p>
          <button class="btn btn-primary" @click="goLogin">去登录</button>
        </div>

        <div v-else class="authorize-actions">
          <button class="btn btn-secondary" @click="handleDeny">拒绝</button>
          <button class="btn btn-primary" @click="handleApprove" :disabled="submitting">
            {{ submitting ? '处理中...' : '授权' }}
          </button>
        </div>

        <p class="cancel-tip">
          <a href="javascript:void(0)" @click="handleDeny">取消并返回</a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { getOAuthAuthorize, confirmOAuthAuthorize } from '@/api/index'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(true)
const error = ref('')
const clientInfo = ref(null)
const submitting = ref(false)

const isLoggedIn = computed(() => userStore.isLoggedIn)

// 权限范围名称映射
const scopeNames = {
  'userinfo': '获取基本信息（昵称、头像）',
  'userinfo:email': '获取邮箱地址',
  'photosets:read': '读取套图列表',
  'favorites:read': '读取收藏列表',
  'favorites:write': '管理收藏',
  'community:read': '读取社区内容',
  'community:write': '发布社区内容'
}

function getScopeName(scope) {
  return scopeNames[scope] || scope
}

function goHome() {
  router.push('/')
}

function goLogin() {
  router.push({
    name: 'Login',
    query: { redirect: route.fullPath }
  })
}

async function loadAuthorizeInfo() {
  loading.value = true
  error.value = ''

  const { client_id, redirect_uri, scope, state, response_type } = route.query

  if (!client_id || !redirect_uri) {
    error.value = '缺少必要的授权参数'
    loading.value = false
    return
  }

  try {
    const res = await getOAuthAuthorize({
      client_id,
      redirect_uri,
      scope: scope || '',
      state: state || '',
      response_type: response_type || 'code'
    })

    if (res.data?.code === 0) {
      clientInfo.value = res.data.data
    } else {
      error.value = res.data?.message || '获取应用信息失败'
    }
  } catch (e) {
    if (e.response?.data?.message) {
      error.value = e.response.data.message
    } else {
      error.value = '网络错误，请稍后重试'
    }
  } finally {
    loading.value = false
  }
}

async function handleApprove() {
  submitting.value = true

  const { client_id, redirect_uri, scope, state } = route.query

  try {
    const res = await confirmOAuthAuthorize({
      client_id,
      redirect_uri,
      scope: scope || '',
      state: state || '',
      approved: true
    })

    if (res.data?.code === 0) {
      // 重定向到第三方应用
      const redirectUrl = res.data.data?.redirect_url
      if (redirectUrl) {
        window.location.href = redirectUrl
      }
    } else {
      alert(res.data?.message || '授权失败，请重试')
    }
  } catch (e) {
    alert(e.response?.data?.message || '授权失败，请重试')
  } finally {
    submitting.value = false
  }
}

async function handleDeny() {
  const { client_id, redirect_uri, state } = route.query

  try {
    const res = await confirmOAuthAuthorize({
      client_id,
      redirect_uri,
      scope: '',
      state: state || '',
      approved: false
    })

    if (res.data?.code === 0) {
      const redirectUrl = res.data.data?.redirect_url
      if (redirectUrl) {
        window.location.href = redirectUrl
      }
    } else {
      router.push('/')
    }
  } catch (e) {
    router.push('/')
  }
}

onMounted(() => {
  loadAuthorizeInfo()
})
</script>

<style scoped>
.oauth-authorize {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.authorize-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  padding: 40px;
  max-width: 440px;
  width: 100%;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 20px 0;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-icon {
  width: 60px;
  height: 60px;
  background: #fee2e2;
  color: #ef4444;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: bold;
  margin: 0 auto 16px;
}

.error-state h2 {
  margin: 0 0 8px;
  color: #1f2937;
  font-size: 20px;
}

.error-state p {
  color: #6b7280;
  margin: 0 0 24px;
}

.authorize-content {
  text-align: center;
}

.app-info {
  margin-bottom: 32px;
}

.app-logo {
  width: 72px;
  height: 72px;
  border-radius: 16px;
  object-fit: cover;
  margin-bottom: 16px;
}

.app-logo-placeholder {
  width: 72px;
  height: 72px;
  border-radius: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 32px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.app-name {
  margin: 0 0 8px;
  color: #1f2937;
  font-size: 24px;
  font-weight: 600;
}

.app-desc {
  margin: 0;
  color: #6b7280;
  font-size: 14px;
}

.scopes-section {
  background: #f9fafb;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 24px;
  text-align: left;
}

.scopes-title {
  margin: 0 0 12px;
  color: #374151;
  font-size: 14px;
  font-weight: 500;
}

.scopes-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.scope-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 0;
  color: #4b5563;
  font-size: 14px;
}

.scope-icon {
  color: #10b981;
  font-weight: bold;
}

.login-prompt {
  margin-bottom: 24px;
}

.login-prompt p {
  color: #6b7280;
  margin: 0 0 16px;
}

.authorize-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-bottom: 16px;
}

.btn {
  padding: 12px 32px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-secondary {
  background: #f3f4f6;
  color: #374151;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.cancel-tip {
  margin: 0;
  font-size: 14px;
}

.cancel-tip a {
  color: #6b7280;
  text-decoration: none;
}

.cancel-tip a:hover {
  color: #374151;
  text-decoration: underline;
}
</style>