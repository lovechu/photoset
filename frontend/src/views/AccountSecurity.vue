<template>
  <div class="account-security-page">
    <el-card>
      <template #header><h3>账号安全</h3></template>

      <div class="security-section">
        <h4>修改密码</h4>
        <el-form label-width="90px" style="max-width: 460px">
          <el-form-item label="当前密码">
            <el-input v-model="pwdForm.old_password" type="password" show-password placeholder="请输入当前密码" />
          </el-form-item>
          <el-form-item label="新密码">
            <el-input v-model="pwdForm.new_password" type="password" show-password placeholder="请输入新密码" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :loading="changingPwd" @click="changePassword">修改密码</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-divider />

      <div class="security-section">
        <div class="section-head">
          <h4>登录设备</h4>
          <el-button size="small" type="danger" plain :disabled="!devices.length" @click="deactivateOthers">退出其他设备</el-button>
        </div>
        <el-table :data="devices" v-loading="loadingDevices" stripe>
          <el-table-column prop="device_name" label="设备名称" min-width="160" />
          <el-table-column prop="device_type" label="类型" width="100" />
          <el-table-column prop="os" label="系统" width="120" />
          <el-table-column prop="browser" label="浏览器" width="120" />
          <el-table-column prop="last_active_at" label="最后活跃" width="180" />
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button link type="danger" size="small" @click="deactivateDevice(row)">退出</el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-if="!loadingDevices && !devices.length" description="暂无其他登录设备" :image-size="50" />
      </div>

      <el-divider />

      <div class="security-section">
        <h4>快捷入口</h4>
        <div class="quick-links">
          <el-button :icon="Message" @click="$router.push('/bind-email')">绑定邮箱</el-button>
          <el-button :icon="Clock" @click="$router.push('/login-history')">登录历史</el-button>
          <el-button :icon="Lock" @click="$router.push('/privacy-settings')">隐私设置</el-button>
          <el-button type="danger" plain :icon="Warning" @click="$router.push('/deletion-status')">账号注销</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Message, Clock, Lock, Warning } from '@element-plus/icons-vue'
import { changePassword as apiChangePassword, getUserDevices, deactivateDevice as apiDeactivateDevice, deactivateOtherDevices } from '@/api'

const pwdForm = ref({ old_password: '', new_password: '' })
const changingPwd = ref(false)
const devices = ref([])
const loadingDevices = ref(false)

onMounted(() => loadDevices())

async function changePassword() {
  if (!pwdForm.value.old_password || !pwdForm.value.new_password) {
    ElMessage.warning('请填写完整')
    return
  }
  changingPwd.value = true
  try {
    await apiChangePassword({ ...pwdForm.value })
    ElMessage.success('密码修改成功，请重新登录')
    pwdForm.value = { old_password: '', new_password: '' }
  } catch (e) {
    ElMessage.error('修改失败')
  } finally {
    changingPwd.value = false
  }
}

async function loadDevices() {
  loadingDevices.value = true
  try {
    const res = await getUserDevices()
    const data = res?.data || {}
    devices.value = data.list || data.devices || (Array.isArray(data) ? data : [])
  } catch (e) {
    devices.value = []
  } finally {
    loadingDevices.value = false
  }
}

async function deactivateDevice(row) {
  try {
    await ElMessageBox.confirm(`确定退出设备「${row.device_name || '该设备'}」？`, '提示', { type: 'warning' })
    await apiDeactivateDevice(row.device_id || row.id)
    devices.value = devices.value.filter(d => (d.device_id || d.id) !== (row.device_id || row.id))
    ElMessage.success('已退出该设备')
  } catch (e) { /* 取消 */ }
}

async function deactivateOthers() {
  try {
    await ElMessageBox.confirm('确定退出其他所有设备？', '提示', { type: 'warning' })
    await deactivateOtherDevices()
    ElMessage.success('已退出其他设备')
    loadDevices()
  } catch (e) { /* 取消 */ }
}
</script>

<style scoped>
.account-security-page { max-width: 1000px; margin: 0 auto; padding: 24px 0 60px; }
.el-card__header h3 { margin: 0; font-size: 18px; }
.security-section h4 { margin: 0 0 16px; font-size: 16px; }
.section-head { display: flex; align-items: center; justify-content: space-between; }
.quick-links { display: flex; flex-wrap: wrap; gap: 12px; }
</style>
