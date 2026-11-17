<template>
  <div class="feedback-page">
    <el-card>
      <template #header>
        <div class="header-row">
          <h3>意见反馈</h3>
          <el-button size="small" @click="$router.push('/user-level')">我的等级</el-button>
        </div>
      </template>

      <el-form label-width="90px" style="max-width: 640px">
        <el-form-item label="反馈类型" required>
          <el-radio-group v-model="form.category">
            <el-radio-button value="bug">问题反馈</el-radio-button>
            <el-radio-button value="suggestion">功能建议</el-radio-button>
            <el-radio-button value="complaint">投诉举报</el-radio-button>
            <el-radio-button value="other">其他</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="反馈内容" required>
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="6"
            minlength="10"
            maxlength="2000"
            show-word-limit
            placeholder="请详细描述你的问题或建议（至少 10 个字）"
          />
        </el-form-item>

        <el-form-item label="图片">
          <div class="image-toolbar">
            <el-button :icon="Plus" :loading="uploading" :disabled="form.images.length >= 5" @click="fileInput?.click()">
              添加图片（{{ form.images.length }}/5）
            </el-button>
            <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/webp" multiple style="display:none" @change="onFilesPicked" />
            <span class="upload-tip">支持 JPG/PNG/WebP，单张不超过 5MB</span>
          </div>
          <div v-if="form.images.length" class="uploaded-images">
            <div v-for="(url, i) in form.images" :key="i" class="uploaded-img">
              <el-image :src="url" fit="cover" />
              <el-icon class="img-remove" @click="removeImage(i)"><Close /></el-icon>
            </div>
          </div>
        </el-form-item>

        <el-form-item label="联系方式">
          <el-input v-model="form.contact" placeholder="邮箱或手机号（选填，便于我们回复你）" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="submitFeedback">提交反馈</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>

      <el-divider />
      <div class="tips">
        <p><strong>温馨提示：</strong></p>
        <ul>
          <li>你的反馈会帮助我们不断改进产品，感谢你的支持</li>
          <li>请尽量详细描述问题发生的场景，便于我们快速定位</li>
          <li>恶意或重复提交无效反馈可能被限制使用反馈功能</li>
        </ul>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Close } from '@element-plus/icons-vue'
import { submitFeedback as apiSubmitFeedback, uploadFeedbackImage } from '@/api'

const fileInput = ref(null)
const uploading = ref(false)
const submitting = ref(false)
const form = ref({ category: 'bug', content: '', images: [], contact: '' })

async function onFilesPicked(e) {
  const files = Array.from(e.target.files || [])
  if (!files.length) return
  uploading.value = true
  try {
    for (const file of files) {
      if (form.value.images.length >= 5) {
        ElMessage.warning('最多上传 5 张图片')
        break
      }
      if (file.size > 5 * 1024 * 1024) {
        ElMessage.warning(`${file.name} 超过 5MB，已跳过`)
        continue
      }
      const res = await uploadFeedbackImage(file)
      const url = res?.data?.url
      if (url) form.value.images.push(url)
    }
    ElMessage.success('图片上传成功')
  } catch (e) {
    ElMessage.error('图片上传失败')
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}

function removeImage(i) {
  form.value.images.splice(i, 1)
}

async function submitFeedback() {
  const content = form.value.content.trim()
  if (content.length < 10) {
    ElMessage.warning('反馈内容至少 10 个字')
    return
  }
  submitting.value = true
  try {
    await apiSubmitFeedback({
      category: form.value.category,
      content,
      images: form.value.images,
      contact: form.value.contact || undefined
    })
    ElMessage.success('感谢你的反馈，我们会尽快处理')
    resetForm()
  } catch (e) {
    ElMessage.error('提交失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

function resetForm() {
  form.value = { category: 'bug', content: '', images: [], contact: '' }
}
</script>

<style scoped>
.feedback-page { max-width: 800px; margin: 0 auto; padding: 24px 0 60px; }
.header-row { display: flex; align-items: center; justify-content: space-between; }
.header-row h3 { margin: 0; font-size: 18px; }
.image-toolbar { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.upload-tip { color: #999; font-size: 12px; }
.uploaded-images { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 12px; }
.uploaded-img { position: relative; width: 90px; height: 90px; }
.uploaded-img .el-image { width: 100%; height: 100%; border-radius: 6px; }
.img-remove { position: absolute; top: -6px; right: -6px; cursor: pointer; color: #f56c6c; background: #fff; border-radius: 50%; }
.tips { color: #666; font-size: 13px; line-height: 1.8; }
.tips p { margin: 0 0 6px; }
.tips ul { margin: 0; padding-left: 20px; }
</style>
