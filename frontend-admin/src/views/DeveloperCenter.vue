<template>
  <div class="developer-center">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>开发者中心</span>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <!-- API 密钥管理 -->
        <el-tab-pane label="API 密钥" name="keys">
          <div class="section">
            <div class="section-header">
              <h3>API 密钥管理</h3>
              <el-button type="primary" @click="showCreateDialog = true">
                <el-icon><Plus /></el-icon> 创建新密钥
              </el-button>
            </div>
            <p class="tip">API 密钥用于调用平台开放接口，请妥善保管，不要泄露给他人。</p>

            <el-table :data="apiKeys" v-loading="loading" stripe>
              <el-table-column prop="name" label="名称" min-width="120" />
              <el-table-column prop="key" label="密钥" min-width="280">
                <template #default="{ row }">
                  <div class="key-cell">
                    <code>{{ row.key }}</code>
                    <el-button size="small" text @click="copyText(row.key)">
                      <el-icon><CopyDocument /></el-icon>
                    </el-button>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                    {{ row.status === 1 ? '启用' : '禁用' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="created_at" label="创建时间" width="180">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
              <el-table-column prop="last_used" label="最后使用" width="180">
                <template #default="{ row }">
                  {{ row.last_used ? formatDate(row.last_used) : '从未使用' }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="100" fixed="right">
                <template #default="{ row }">
                  <el-button type="danger" size="small" text @click="handleDelete(row)">
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <!-- API 文档 -->
        <el-tab-pane label="API 文档" name="docs">
          <div class="section">
            <div class="section-header">
              <h3>开放 API 文档</h3>
              <el-button @click="refreshDocs">
                <el-icon><Refresh /></el-icon> 刷新
              </el-button>
            </div>

            <el-alert type="info" :closable="false" class="mb-16">
              <template #title>
                <strong>认证方式</strong>：在请求 Header 中添加 <code>Authorization: Bearer &lt;token&gt;</code>
                <br>
                <strong>Content-Type</strong>：<code>application/json</code>
              </template>
            </el-alert>

            <div v-for="cat in apiDocs" :key="cat.category" class="api-category">
              <h4>{{ cat.category }}</h4>
              <el-table :data="cat.endpoints" stripe size="small">
                <el-table-column prop="method" label="方法" width="100">
                  <template #default="{ row }">
                    <el-tag :type="getMethodType(row.method)" size="small">
                      {{ row.method }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="path" label="路径" min-width="250">
                  <template #default="{ row }">
                    <code>{{ row.path }}</code>
                  </template>
                </el-table-column>
                <el-table-column prop="desc" label="说明" />
              </el-table>
            </div>
          </div>
        </el-tab-pane>

        <!-- 签名 URL 文档 -->
        <el-tab-pane label="签名 URL" name="sign">
          <div class="section">
            <div class="section-header">
              <h3>图片签名 URL 机制</h3>
            </div>

            <el-alert type="warning" :closable="false" class="mb-16">
              <template #title>
                <strong>付费图片防盗链</strong>：付费套图中的图片使用签名 URL 访问，需要正确签名才能查看。
              </template>
            </el-alert>

            <div class="sign-info">
              <h4>签名参数</h4>
              <el-table :data="signParams" size="small">
                <el-table-column prop="name" label="参数" width="120">
                  <template #default="{ row }">
                    <code>{{ row.name }}</code>
                  </template>
                </el-table-column>
                <el-table-column prop="desc" label="说明" />
              </el-table>
            </div>

            <div class="sign-info">
              <h4>签名算法</h4>
              <p>使用 HMAC-SHA256 算法对 URL 路径和过期时间进行签名。</p>
              <div class="code-block">
                <p><strong>签名公式：</strong></p>
                <pre>message = "path?expires=timestamp"
sign = HMAC-SHA256(message, secret_key)</pre>
              </div>
            </div>

            <div class="sign-info">
              <h4>Python 示例</h4>
              <pre class="code-block">{{ signCode.python }}</pre>
            </div>

            <div class="sign-info">
              <h4>JavaScript 示例</h4>
              <pre class="code-block">{{ signCode.javascript }}</pre>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 创建密钥对话框 -->
    <el-dialog v-model="showCreateDialog" title="创建 API 密钥" width="500px">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="密钥名称" required>
          <el-input v-model="createForm.name" placeholder="例如：我的应用密钥" maxlength="50" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreate">创建</el-button>
      </template>
    </el-dialog>

    <!-- 密钥创建成功对话框 -->
    <el-dialog v-model="showKeyResult" title="密钥创建成功" width="600px">
      <el-alert type="warning" :closable="false" class="mb-16">
        <strong>请立即复制并妥善保存！</strong> 此密钥只会显示这一次，之后无法找回。
      </el-alert>
      <div class="key-result">
        <p><strong>密钥名称：</strong>{{ newKey.name }}</p>
        <p><strong>Key：</strong><code class="key-code">{{ newKey.key }}</code></p>
        <p><strong>Secret：</strong><code class="key-code">{{ newKey.secret }}</code></p>
      </div>
      <template #footer>
        <el-button type="primary" @click="copyNewKey">复制全部信息</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, CopyDocument, Refresh } from '@element-plus/icons-vue'
import { getApiKeys, createApiKey, deleteApiKey, getApiDocs, getSignUrlDocs } from '@/api/index'

const activeTab = ref('keys')
const loading = ref(false)
const creating = ref(false)
const apiKeys = ref([])
const apiDocs = ref([])
const showCreateDialog = ref(false)
const showKeyResult = ref(false)
const createForm = ref({ name: '' })
const newKey = ref({})
const signParams = ref([
  { name: 'sign', desc: 'HMAC-SHA256 签名（十六进制）' },
  { name: 'expires', desc: '签名过期时间戳（Unix 秒级时间戳）' },
])
const signCode = ref({
  python: `import hmac
import hashlib
import time

def generate_sign_url(path, secret_key, expires=3600):
    expires_at = int(time.time()) + expires
    message = "%s?expires=%d" % (path, expires_at)
    sign = hmac.new(secret_key.encode(), message.encode(), hashlib.sha256).hexdigest()
    return "%s&sign=%s" % (message, sign)`,
  javascript: `// generateSignUrl(path, secretKey, expires=3600)
const expiresAt = Math.floor(Date.now() / 1000) + expires;
const message = path + "?expires=" + expiresAt;
const sign = crypto.createHmac('sha256', secretKey).update(message).digest('hex');
return message + "&sign=" + sign;`,
})

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString('zh-CN')
}

function getMethodType(method) {
  const types = {
    'GET': 'success',
    'POST': 'primary',
    'PUT': 'warning',
    'DELETE': 'danger',
  }
  return types[method] || 'info'
}

async function copyText(text) {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

async function loadApiKeys() {
  loading.value = true
  try {
    const res = await getApiKeys()
    apiKeys.value = res.data || []
  } catch (e) {
    // 错误已在拦截器处理
  } finally {
    loading.value = false
  }
}

async function refreshDocs() {
  try {
    const res = await getApiDocs()
    apiDocs.value = res.data?.docs || []
    ElMessage.success('文档已刷新')
  } catch (e) {
    // 错误已在拦截器处理
  }
}

async function handleCreate() {
  if (!createForm.value.name || createForm.value.name.length < 2) {
    ElMessage.warning('请输入密钥名称（至少2个字符）')
    return
  }
  creating.value = true
  try {
    const res = await createApiKey(createForm.value.name)
    newKey.value = res.data?.key || {}
    showCreateDialog.value = false
    showKeyResult.value = true
    createForm.value.name = ''
    loadApiKeys()
  } catch (e) {
    // 错误已在拦截器处理
  } finally {
    creating.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(
      `确定要删除密钥「${row.name}」吗？此操作不可恢复。`,
      '删除确认',
      { type: 'warning' }
    )
    await deleteApiKey(row.id)
    ElMessage.success('密钥已删除')
    loadApiKeys()
  } catch (e) {
    if (e !== 'cancel') {
      // 错误已在拦截器处理
    }
  }
}

function copyNewKey() {
  const text = `名称: ${newKey.value.name}\nKey: ${newKey.value.key}\nSecret: ${newKey.value.secret}`
  copyText(text)
  showKeyResult.value = false
}

onMounted(() => {
  loadApiKeys()
  refreshDocs()
})
</script>

<style scoped>
.developer-center {
  padding: 20px;
}

.card-header {
  font-size: 18px;
  font-weight: 600;
}

.section {
  padding: 10px 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h3 {
  margin: 0;
  font-size: 16px;
}

.tip {
  color: #909399;
  font-size: 13px;
  margin-bottom: 16px;
}

.mb-16 {
  margin-bottom: 16px;
}

.key-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.key-cell code {
  font-family: monospace;
  font-size: 12px;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
}

.api-category {
  margin-bottom: 24px;
}

.api-category h4 {
  margin-bottom: 12px;
  font-size: 15px;
  color: #303133;
}

.sign-info {
  margin-bottom: 24px;
}

.sign-info h4 {
  margin-bottom: 12px;
  font-size: 15px;
  color: #303133;
}

.code-block {
  background: #f5f7fa;
  padding: 12px 16px;
  border-radius: 8px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  line-height: 1.6;
  overflow-x: auto;
}

.code-block pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
}

.key-result {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 8px;
}

.key-result p {
  margin: 8px 0;
  word-break: break-all;
}

.key-code {
  display: block;
  background: #fff;
  padding: 8px;
  border-radius: 4px;
  margin-top: 4px;
  font-family: monospace;
}
</style>
