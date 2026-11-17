<template>
  <div class="deletion-status-page">
    <el-card>
      <template #header><h3>账号注销</h3></template>

      <template v-if="status?.has_pending_deletion">
        <el-result icon="warning" title="注销申请处理中" :sub-title="`申请时间：${formatDate(status.requested_at)}`">
          <template #extra>
            <div class="pending-info">
              <p>冷静期剩余 <strong>{{ status.days_remaining }}</strong> 天</p>
              <p class="deadline">截止日期：{{ formatDate(status.deadline) }}</p>
              <p class="tip">冷静期内登录账号，注销申请将自动取消。如需保留账号，也可点击下方按钮主动撤销。</p>
              <el-button type="primary" @click="cancelDeletion">撤销注销申请</el-button>
            </div>
          </template>
        </el-result>
      </template>

      <template v-else>
        <el-alert
          type="warning"
          :closable="false"
          show-icon
          title="注销账号是不可逆操作，请谨慎考虑"
          description="注销后将删除你的账号数据、帖子、收藏等全部信息，30 天冷静期内登录可撤销申请。"
        />
        <div class="deletion-form">
          <el-form label-width="90px" style="max-width: 500px">
            <el-form-item label="登录密码">
              <el-input v-model="form.password" type="password" show-password placeholder="请输入当前密码以确认身份" />
            </el-form-item>
            <el-form-item label="注销原因">
              <el-input v-model="form.reason" type="textarea" :rows="3" placeholder="请告诉我们注销原因（选填）" />
            </el-form-item>
            <el-form-item>
              <el-button type="danger" :loading="submitting" @click="requestDeletion">申请注销</el-button>
            </el-form-item>
          </el-form>
        </div>
      </template>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDeletionStatus, requestDeletion as apiRequestDeletion, cancelDeletion as apiCancelDeletion } from '@/api'

const status = ref(null)
const form = ref({ password: '', reason: '' })
const submitting = ref(false)

onMounted(() => load())
async function load() {
  try {
    const res = await getDeletionStatus()
    status.value = res?.data || { has_pending_deletion: false }
  } catch (e) {
    status.value = { has_pending_deletion: false }
  }
}

async function requestDeletion() {
  if (!form.value.password) return ElMessage.warning('请输入密码')
  try {
    await ElMessageBox.confirm('确定申请注销账号？此操作在 30 天冷静期内可撤销。', '危险操作', {
      type: 'warning', confirmButtonText: '确认申请', cancelButtonText: '再想想'
    })
    submitting.value = true
    await apiRequestDeletion({ ...form.value })
    ElMessage.success('注销申请已提交')
    form.value = { password: '', reason: '' }
    load()
  } catch (e) { /* 取消 */ }
  finally { submitting.value = false }
}

async function cancelDeletion() {
  try {
    await ElMessageBox.confirm('确定撤销注销申请？', '提示', { type: 'info' })
    await apiCancelDeletion()
    ElMessage.success('已撤销注销申请')
    load()
  } catch (e) { /* 取消 */ }
}

function formatDate(str) {
  if (!str) return ''
  return new Date(str).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.deletion-status-page { max-width: 760px; margin: 0 auto; padding: 24px 0 60px; }
.el-card__header h3 { margin: 0; font-size: 18px; }
.deletion-form { margin-top: 24px; }
.pending-info p { margin: 6px 0; }
.pending-info strong { color: #f56c6c; font-size: 20px; }
.pending-info .deadline { color: #999; }
.pending-info .tip { color: #e6a23c; font-size: 13px; margin: 16px 0; }
</style>
