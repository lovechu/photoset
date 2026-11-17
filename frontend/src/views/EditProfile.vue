<template>
  <div class="edit-profile-page">
    <el-card>
      <template #header><h3>编辑资料</h3></template>
      <el-form label-width="90px" style="max-width: 520px">
        <el-form-item label="头像">
          <div class="avatar-edit">
            <el-avatar :size="72" :src="previewAvatar || userStore.avatar">
              {{ userStore.user?.nickname?.charAt(0) }}
            </el-avatar>
            <input ref="avatarInput" type="file" accept="image/*" style="display:none" @change="onAvatarPicked" />
            <el-button size="small" :icon="Camera" :loading="uploadingAvatar" @click="avatarInput?.click()">更换头像</el-button>
          </div>
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nickname" maxlength="30" placeholder="你的昵称" />
        </el-form-item>
        <el-form-item label="个性签名">
          <el-input v-model="form.bio" type="textarea" :rows="3" maxlength="200" show-word-limit placeholder="介绍一下自己吧" />
        </el-form-item>
        <el-form-item label="封面图">
          <el-input v-model="form.cover_image" placeholder="封面图地址（选填）" />
        </el-form-item>
        <el-form-item label="IP属地">
          <el-input v-model="form.ip_location" placeholder="选填" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving" @click="saveProfile">保存修改</el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { Camera } from '@element-plus/icons-vue'
import { uploadAvatar } from '@/api'

const userStore = useUserStore()
const avatarInput = ref(null)
const form = ref({ nickname: '', bio: '', cover_image: '', ip_location: '' })
const previewAvatar = ref('')
const uploadingAvatar = ref(false)
const saving = ref(false)

onMounted(() => {
  const u = userStore.user || {}
  form.value = {
    nickname: u.nickname || '',
    bio: u.bio || '',
    cover_image: u.cover_image || '',
    ip_location: u.ip_location || ''
  }
  previewAvatar.value = u.avatar || ''
})

async function onAvatarPicked(e) {
  const file = e.target.files?.[0]
  if (!file) return
  uploadingAvatar.value = true
  try {
    const res = await uploadAvatar(file)
    const url = res?.data?.url || res?.data?.avatar_url || res?.url
    if (url) {
      previewAvatar.value = url
      form.value.avatar = url
      ElMessage.success('头像已上传')
    }
  } catch (err) {
    ElMessage.error('头像上传失败')
  } finally {
    uploadingAvatar.value = false
    e.target.value = ''
  }
}

async function saveProfile() {
  if (!form.value.nickname.trim()) {
    ElMessage.warning('昵称不能为空')
    return
  }
  saving.value = true
  try {
    await userStore.updateUserProfile({ ...form.value })
    ElMessage.success('资料已保存')
    routerBack()
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

function routerBack() {
  window.history.back()
}
</script>

<style scoped>
.edit-profile-page { max-width: 760px; margin: 0 auto; padding: 24px 0 60px; }
.el-card__header h3 { margin: 0; font-size: 18px; }
.avatar-edit { display: flex; align-items: center; gap: 16px; }
</style>
