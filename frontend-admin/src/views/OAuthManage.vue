<template>
  <div class="oauth-manage">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>OAuth2 应用管理</span>
          <el-button type="primary" @click="showCreateDialog = true">
            <el-icon><Plus /></el-icon> 创建应用
          </el-button>
        </div>
      </template>

      <p class="tip">管理第三方 OAuth2 应用，用于授权第三方应用访问用户数据。</p>

      <el-table :data="oauthClients" v-loading="loading" stripe>
        <el-table-column prop="name" label="应用名称" min-width="150" />
        <el-table-column prop="client_id" label="Client ID" min-width="280">
          <template #default="{ row }">
            <div class="key-cell">
              <code>{{ row.client_id }}</code>
              <el-button size="small" text @click="copyText(row.client_id)">
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
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" text @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" text @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑应用对话框 -->
    <el-dialog v-model="showCreateDialog" :title="editingClient ? '编辑应用' : '创建应用'" width="600px">
      <el-form :model="createForm" label-width="120px">
        <el-form-item label="应用名称" required>
          <el-input v-model="createForm.name" placeholder="例如：我的第三方应用" maxlength="100" />
        </el-form-item>
        <el-form-item label="应用描述">
          <el-input v-model="createForm.description" type="textarea" :rows="3" placeholder="应用功能描述" />
        </el-form-item>
        <el-form-item label="重定向URI" required>
          <el-input v-model="createForm.redirect_uris" placeholder="例如：https://example.com/callback,http://localhost:3000/callback" />
          <div class="form-tip">多个URI用逗号分隔</div>
        </el-form-item>
        <el-form-item label="权限范围" required>
          <el-select v-model="createForm.scopes" multiple placeholder="选择权限范围">
            <el-option label="用户基本信息" value="userinfo" />
            <el-option label="用户邮箱" value="userinfo:email" />
            <el-option label="读取套图" value="photosets:read" />
            <el-option label="读取收藏" value="favorites:read" />
            <el-option label="管理收藏" value="favorites:write" />
            <el-option label="读取社区" value="community:read" />
            <el-option label="发布社区" value="community:write" />
          </el-select>
        </el-form-item>
        <el-form-item label="应用Logo">
          <el-input v-model="createForm.logo_url" placeholder="应用Logo URL" />
        </el-form-item>
        <el-form-item label="状态" v-if="editingClient">
          <el-switch v-model="createForm.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleSave">
          {{ editingClient ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 应用创建成功对话框 -->
    <el-dialog v-model="showClientResult" title="应用创建成功" width="600px">
      <el-alert type="warning" :closable="false" class="mb-16">
        <strong>请立即复制并妥善保存！</strong> Client Secret 只会显示这一次，之后无法找回。
      </el-alert>
      <div class="client-result">
        <p><strong>应用名称：</strong>{{ newClient.name }}</p>
        <p><strong>Client ID：</strong><code class="key-code">{{ newClient.client_id }}</code></p>
        <p><strong>Client Secret：</strong><code class="key-code">{{ newClient.client_secret }}</code></p>
      </div>
      <template #footer>
        <el-button type="primary" @click="copyNewClient">复制全部信息</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, CopyDocument } from '@element-plus/icons-vue'
import { getOAuthClients, createOAuthClient, updateOAuthClient, deleteOAuthClient } from '@/api/index'

const loading = ref(false)
const creating = ref(false)
const oauthClients = ref([])
const showCreateDialog = ref(false)
const showClientResult = ref(false)
const editingClient = ref(null)
const createForm = ref({
  name: '',
  description: '',
  redirect_uris: '',
  scopes: [],
  logo_url: '',
  status: 1
})
const newClient = ref({})

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString('zh-CN')
}

async function copyText(text) {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

async function loadOAuthClients() {
  loading.value = true
  try {
    const res = await getOAuthClients()
    oauthClients.value = res.data || []
  } catch (e) {
    // 错误已在拦截器处理
  } finally {
    loading.value = false
  }
}

function handleEdit(row) {
  editingClient.value = row
  createForm.value = {
    name: row.name,
    description: row.description || '',
    redirect_uris: row.redirect_uris,
    scopes: row.scopes ? row.scopes.split(',') : [],
    logo_url: row.logo_url || '',
    status: row.status
  }
  showCreateDialog.value = true
}

async function handleSave() {
  if (!createForm.value.name || createForm.value.name.length < 2) {
    ElMessage.warning('请输入应用名称（至少2个字符）')
    return
  }
  if (!createForm.value.redirect_uris) {
    ElMessage.warning('请输入重定向URI')
    return
  }
  if (createForm.value.scopes.length === 0) {
    ElMessage.warning('请选择至少一个权限范围')
    return
  }

  creating.value = true
  try {
    const data = {
      ...createForm.value,
      redirect_uris: createForm.value.redirect_uris.split(',').map(uri => uri.trim()).filter(uri => uri),
      scopes: createForm.value.scopes
    }
    
    if (editingClient.value) {
      await updateOAuthClient(editingClient.value.id, data)
      ElMessage.success('应用已更新')
    } else {
      const res = await createOAuthClient(data)
      newClient.value = res.data?.client || {}
      showCreateDialog.value = false
      showClientResult.value = true
      ElMessage.success('应用创建成功')
    }
    
    createForm.value = {
      name: '',
      description: '',
      redirect_uris: '',
      scopes: [],
      logo_url: '',
      status: 1
    }
    editingClient.value = null
    loadOAuthClients()
  } catch (e) {
    // 错误已在拦截器处理
  } finally {
    creating.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(
      `确定要删除应用「${row.name}」吗？此操作不可恢复，所有已授权的令牌将被撤销。`,
      '删除确认',
      { type: 'warning' }
    )
    await deleteOAuthClient(row.id)
    ElMessage.success('应用已删除')
    loadOAuthClients()
  } catch (e) {
    if (e !== 'cancel') {
      // 错误已在拦截器处理
    }
  }
}

function copyNewClient() {
  const text = `应用名称: ${newClient.value.name}\nClient ID: ${newClient.value.client_id}\nClient Secret: ${newClient.value.client_secret}`
  copyText(text)
  showClientResult.value = false
}

onMounted(() => {
  loadOAuthClients()
})
</script>

<style scoped>
.oauth-manage {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 18px;
  font-weight: 600;
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

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.client-result {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 8px;
}

.client-result p {
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