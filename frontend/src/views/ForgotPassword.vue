<template>
  <div class="auth-page">
    <div class="auth-container">
      <div class="auth-card">
        <div class="auth-header">
          <router-link to="/" class="logo">
            <el-icon :size="32"><Camera /></el-icon>
          </router-link>
          <h1>找回密码</h1>
          <p>输入注册邮箱，我们将发送重置链接</p>
        </div>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          class="auth-form"
          @submit.prevent="handleSubmit"
        >
          <el-form-item prop="email">
            <el-input
              v-model="form.email"
              placeholder="请输入注册邮箱"
              size="large"
              :prefix-icon="Message"
            />
          </el-form-item>

          <el-form-item prop="captcha_code">
            <div class="captcha-container">
              <el-input
                v-model="form.captcha_code"
                placeholder="请输入验证码"
                size="large"
                class="captcha-input"
              />
              <div class="captcha-image-wrapper">
                <img
                  v-if="form.captcha_image"
                  :src="form.captcha_image"
                  alt="验证码"
                  class="captcha-image"
                  @click="refreshCaptcha"
                />
                <el-button
                  v-else
                  type="default"
                  size="large"
                  :icon="Refresh"
                  @click="refreshCaptcha"
                >
                  获取验证码
                </el-button>
              </div>
            </div>
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              native-type="submit"
              class="submit-btn"
            >
              {{ loading ? '发送中...' : '发送重置链接' }}
            </el-button>
          </el-form-item>
        </el-form>

        <div class="auth-footer">
          <span>想起密码了？</span>
          <router-link to="/login">立即登录</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCaptcha, forgotPassword } from '@/api'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import { Camera, Message, Refresh } from '@element-plus/icons-vue'

const router = useRouter()

const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  email: '',
  captcha_code: '',
  captcha_id: '',
  captcha_image: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  captcha_code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { min: 4, message: '验证码至少 4 个字符', trigger: 'blur' }
  ]
}

const refreshCaptcha = async () => {
  try {
    const res = await request.get('/auth/captcha?action=forgot')
    form.captcha_id = res.data.captcha_id
    form.captcha_image = res.data.captcha_image
  } catch (e) {
    console.error('获取验证码失败:', e)
  }
}

onMounted(() => {
  refreshCaptcha()
})

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch (e) {
    return
  }

  loading.value = true

  try {
    await forgotPassword({
      email: form.email,
      captcha_id: form.captcha_id,
      captcha_code: form.captcha_code
    })
    ElMessage.success('重置链接已发送到您的邮箱，请查收（30分钟内有效）')
    router.push('/login')
  } catch (e) {
    // 刷新验证码
    refreshCaptcha()
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
  margin-left: 4px;
}

.auth-footer a:hover {
  text-decoration: underline;
}

/* 验证码样式 */
.captcha-container {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.captcha-input {
  flex: 1;
}

.captcha-image-wrapper {
  width: 120px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.captcha-image {
  width: 100%;
  height: 100%;
  border-radius: 4px;
  cursor: pointer;
  border: 1px solid #dcdfe6;
}

.captcha-image:hover {
  opacity: 0.9;
}

@media (max-width: 768px) {
  .auth-card {
    padding: 32px 24px;
  }
}
</style>
