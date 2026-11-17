<template>
  <div class="bind-email-page">
    <el-card>
      <template #header><h3>绑定邮箱</h3></template>
      <div v-if="userStore.user?.email" class="current-email">
        <span>当前邮箱：</span>
        <el-tag type="success">{{ userStore.user.email }}</el-tag>
      </div>
      <el-form v-else label-width="90px" style="max-width: 480px">
        <el-form-item label="邮箱地址">
          <el-input v-model="form.email" placeholder="请输入邮箱地址" />
        </el-form-item>
        <el-form-item label="验证码">
          <el-input v-model="form.code" placeholder="6 位验证码">
            <template #append>
              <el-button :disabled="countdown > 0" :loading="sending" @click="sendCode">
                {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="submitBind">确认绑定</el-button>
        </el-form-item>
      </el-form>
      <el-alert type="info" :closable="false" show-icon style="max-width:480px;margin-top:12px">
        绑定邮箱后可用于找回密码和接收系统通知。
      </el-alert>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { sendEmailCode, bindEmail } from '@/api'

const userStore = useUserStore()
const form = ref({ email: '', code: '' })
const sending = ref(false)
const submitting = ref(false)
const countdown = ref(0)
let timer = null

onUnmounted(() => { if (timer) clearInterval(timer) })

async function sendCode() {
  if (!form.value.email || !/^\S+@\S+\.\S+$/.test(form.value.email)) {
    ElMessage.warning('请输入正确的邮箱地址')
    return
  }
  sending.value = true
  try {
    await sendEmailCode({ email: form.value.email, purpose: 'bind' })
    ElMessage.success('验证码已发送')
    countdown.value = 60
    timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0 && timer) { clearInterval(timer); timer = null }
    }, 1000)
  } catch (e) {
    ElMessage.error('发送失败')
  } finally {
    sending.value = false
  }
}

async function submitBind() {
  if (!form.value.code) {
    ElMessage.warning('请输入验证码')
    return
  }
  submitting.value = true
  try {
    const res = await bindEmail({ email: form.value.email, code: form.value.code })
    const updatedUser = res?.data?.user
    if (updatedUser) userStore.updateUser({ ...userStore.user, ...updatedUser })
    ElMessage.success('邮箱绑定成功')
  } catch (e) {
    ElMessage.error('绑定失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.bind-email-page { max-width: 760px; margin: 0 auto; padding: 24px 0 60px; }
.el-card__header h3 { margin: 0; font-size: 18px; }
.current-email { display: flex; align-items: center; gap: 12px; padding: 12px 0; }
</style>
