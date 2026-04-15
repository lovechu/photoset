<template>
  <div class="login-page">
    <el-card class="login-card" shadow="always">
      <template #header>
        <div class="login-header">
          <h2>PhotoSet 管理后台</h2>
          <p>请使用管理员账号登录</p>
        </div>
      </template>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="0"
        size="large"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="email">
          <el-input
            v-model="form.email"
            placeholder="管理员邮箱"
            :prefix-icon="Message"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        
        <!-- 验证码 -->
        <el-form-item prop="captcha_code">
          <div class="captcha-container">
            <el-input
              v-model="form.captcha_code"
              placeholder="请输入验证码"
              style="flex: 1"
            />
            <div class="captcha-image-wrapper">
              <img 
                v-if="form.captcha_image" 
                :src="form.captcha_image" 
                alt="验证码" 
                class="captcha-image"
                @click="refreshCaptcha"
                style="cursor: pointer; height: 40px; border-radius: 4px;"
              />
              <el-button
                v-else
                type="default"
                :icon="Refresh"
                @click="refreshCaptcha"
                size="default"
              >
                获取验证码
              </el-button>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            style="width: 100%"
            @click="handleLogin"
          >
            登 录
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { ElMessage } from 'element-plus'
import { Message, Lock, Refresh } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()
const adminStore = useAdminStore()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  email: '',
  password: '',
  captcha_code: '',
  captcha_id: '',
  captcha_image: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少 6 位', trigger: 'blur' }
  ],
  captcha_code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { min: 4, message: '验证码至少 4 位', trigger: 'blur' }
  ]
}

// 刷新验证码
const refreshCaptcha = async () => {
  try {
    const res = await request.get('/auth/captcha', {
      params: { action: 'login' }
    })
    form.captcha_id = res.data.captcha_id
    form.captcha_image = res.data.captcha_image
  } catch (e) {
    console.error('获取验证码失败:', e)
    ElMessage.error('获取验证码失败，请重试')
  }
}

// 组件挂载时获取验证码
onMounted(() => {
  refreshCaptcha()
})

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const data = await adminStore.login({ 
      email: form.email, 
      password: form.password,
      captcha_id: form.captcha_id,
      captcha_code: form.captcha_code
    })
    if (data.user.role !== 'admin') {
      ElMessage.error('非管理员账号，无法进入后台')
      adminStore.logout()
      return
    }
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch (err) {
    // 登录失败时刷新验证码
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 420px;
  border-radius: 12px;
}

.login-card :deep(.el-card__header) {
  padding: 30px 30px 10px;
}

.login-card :deep(.el-card__body) {
  padding: 20px 30px 30px;
}

.login-header {
  text-align: center;
}

.login-header h2 {
  margin: 0 0 8px;
  font-size: 24px;
  color: #303133;
}

.login-header p {
  margin: 0;
  font-size: 14px;
  color: #909399;
}

/* 验证码样式 */
.captcha-container {
  display: flex;
  align-items: center;
  width: 100%;
}

.captcha-image-wrapper {
  cursor: pointer;
  position: relative;
  height: 40px;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid #dcdfe6;
  transition: all 0.3s ease;
  background-color: #f5f7fa;
}

.captcha-image-wrapper:hover {
  border-color: #409eff;
  box-shadow: 0 0 8px rgba(64, 158, 255, 0.2);
}

.captcha-image-container {
  position: relative;
  width: 120px;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.captcha-image {
  height: 100%;
  object-fit: contain;
}

.captcha-refresh-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.captcha-image-wrapper:hover .captcha-refresh-overlay {
  opacity: 1;
}

.captcha-refresh-overlay .el-icon {
  color: white;
  font-size: 20px;
}

.captcha-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 120px;
  height: 100%;
  color: #909399;
  font-size: 12px;
}

.captcha-loading .el-icon {
  margin-right: 4px;
  animation: rotating 2s linear infinite;
}

@keyframes rotating {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
