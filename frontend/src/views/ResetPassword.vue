<template>
  <div class="auth-page">
    <div class="auth-container">
      <div class="auth-card">
        <div class="auth-header">
          <router-link to="/" class="logo">
            <el-icon :size="32"><Camera /></el-icon>
          </router-link>
          <h1>重置密码</h1>
          <p>请输入您的新密码</p>
        </div>

        <div v-if="errorMsg" class="error-banner">
          <el-alert :title="errorMsg" type="error" show-icon :closable="false" />
        </div>

        <el-form
          v-if="!errorMsg"
          ref="formRef"
          :model="form"
          :rules="rules"
          class="auth-form"
          @submit.prevent="handleSubmit"
        >
          <el-form-item prop="new_password">
            <el-input
              v-model="form.new_password"
              type="password"
              placeholder="请输入新密码（至少6位）"
              size="large"
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>

          <el-form-item prop="confirm_password">
            <el-input
              v-model="form.confirm_password"
              type="password"
              placeholder="请确认新密码"
              size="large"
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              native-type="submit"
              class="submit-btn"
            >
              {{ loading ? '重置中...' : '重置密码' }}
            </el-button>
          </el-form-item>
        </el-form>

        <div v-if="successMsg" class="success-banner">
          <el-alert :title="successMsg" type="success" show-icon :closable="false" />
          <div style="text-align: center; margin-top: 16px;">
            <router-link to="/login" class="login-link">立即登录</router-link>
          </div>
        </div>

        <div class="auth-footer">
          <router-link to="/login">返回登录</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { resetPasswordByToken } from '@/api'
import { ElMessage } from 'element-plus'
import { Camera, Lock } from '@element-plus/icons-vue'

const route = useRoute()

const formRef = ref(null)
const loading = ref(false)
const errorMsg = ref('')
const successMsg = ref('')

const form = reactive({
  new_password: '',
  confirm_password: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.new_password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少 6 个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

onMounted(() => {
  const token = route.query.token
  if (!token) {
    errorMsg.value = '重置链接无效，缺少必要参数'
  }
})

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch (e) {
    return
  }

  const token = route.query.token
  if (!token) {
    errorMsg.value = '重置链接无效'
    return
  }

  loading.value = true

  try {
    await resetPasswordByToken({
      token: token,
      new_password: form.new_password
    })
    successMsg.value = '密码重置成功！请使用新密码登录'
  } catch (e) {
    console.error('重置密码失败:', e)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  min-height: calc(100vh - 180px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.auth-container {
  width: 100%;
  max-width: 400px;
}

.auth-card {
  background: #fff;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  color: #fff;
  margin-bottom: 16px;
}

.auth-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.auth-header p {
  color: #666;
  font-size: 14px;
}

.error-banner {
  margin-bottom: 24px;
}

.success-banner {
  margin-bottom: 24px;
}

.login-link {
  color: #409eff;
  text-decoration: none;
  font-size: 16px;
  font-weight: 500;
}

.login-link:hover {
  text-decoration: underline;
}

.auth-form {
  margin-bottom: 24px;
}

.auth-form :deep(.el-form-item) {
  margin-bottom: 20px;
}

.auth-form :deep(.el-input__wrapper) {
  padding: 4px 12px;
}

.submit-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
}

.auth-footer {
  text-align: center;
  color: #666;
  font-size: 14px;
}

.auth-footer a {
  color: #409eff;
  text-decoration: none;
}

.auth-footer a:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .auth-card {
    padding: 32px 24px;
  }
}
</style>
