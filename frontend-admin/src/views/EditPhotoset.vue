<template>
  <div>
    <div style="margin-bottom: 20px; display: flex; align-items: center; gap: 12px;">
      <el-button :icon="ArrowLeft" @click="$router.back()">返回</el-button>
      <h2 style="margin: 0;">编辑套图 #{{ photosetId }}</h2>
    </div>

    <el-card v-if="pageLoading">
      <el-skeleton :rows="8" animated />
    </el-card>

    <el-card v-else>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" maxlength="200" show-word-limit />
        </el-form-item>

        <el-form-item label="封面图" prop="cover">
          <div style="display: flex; flex-direction: column; gap: 8px;">
            <el-image v-if="form.cover" :src="form.cover" style="width: 200px; height: 133px; border-radius: 8px;" fit="cover" />
            <el-upload action="" :http-request="handleCoverUpload" :show-file-list="false" accept="image/*">
              <el-button plain :icon="UploadFilled">上传封面</el-button>
            </el-upload>
            <el-input v-model="form.cover" placeholder="或手动输入封面URL" />
          </div>
        </el-form-item>

        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="4" maxlength="1000" show-word-limit />
        </el-form-item>

        <el-form-item label="标签">
          <el-select v-model="form.tags" multiple filterable allow-create default-first-option placeholder="选择或输入标签" style="width: 100%">
            <el-option v-for="tag in availableTags" :key="tag.id" :label="tag.name" :value="tag.name" />
          </el-select>
        </el-form-item>

        <el-form-item label="分类">
          <el-select v-model="form.category" placeholder="选择分类（可选）" clearable style="width: 100%">
            <el-option v-for="cat in availableCategories" :key="cat.slug" :label="cat.name" :value="cat.slug" />
          </el-select>
        </el-form-item>

        <el-form-item label="收费方式">
          <el-radio-group v-model="form.is_free">
            <el-radio :label="1">免费</el-radio>
            <el-radio :label="0">付费</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="form.is_free === 0" label="价格">
          <el-input-number v-model="form.price" :min="0.01" :max="9999" :precision="2" />
        </el-form-item>

        <el-form-item label="图片列表">
          <div style="display: flex; flex-direction: column; gap: 10px; width: 100%;">
            <div v-for="(photo, index) in form.photos" :key="index" style="display: flex; align-items: center; gap: 10px;">
              <span style="width: 24px; text-align: center; color: #909399; font-size: 12px;">{{ index + 1 }}</span>
              <el-input v-model="photo.url" placeholder="图片URL" style="flex: 1;" />
              <el-input-number v-model="photo.sort_order" :min="0" controls-position="right" style="width: 90px;" />
              <el-button type="danger" :icon="Delete" circle size="small" @click="form.photos.splice(index, 1)" />
            </div>
            <div style="display: flex; gap: 10px; margin-top: 4px;">
              <el-upload action="" :http-request="handlePhotoUpload" :show-file-list="false" accept="image/*" multiple>
                <el-button plain :icon="UploadFilled" size="small">上传图片</el-button>
              </el-upload>
              <el-button plain :icon="Plus" size="small" @click="form.photos.push({ url: '', sort_order: form.photos.length })">添加URL</el-button>
            </div>
          </div>
        </el-form-item>

        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio label="published">已发布</el-radio>
            <el-radio label="pending">待审核</el-radio>
            <el-radio label="draft">草稿</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">保存修改</el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getPhotosetDetail, updatePhotoset, getTags, uploadImage, getPublicCategories } from '@/api'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Delete, Plus, UploadFilled } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const photosetId = Number(route.params.id)

const formRef = ref(null)
const loading = ref(false)
const pageLoading = ref(true)
const availableTags = ref([])
const availableCategories = ref([])

const form = reactive({
  title: '', cover: '', description: '',
  tags: [], is_free: 1, price: 0,
  category: '',
  photos: [], status: 'published'
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  cover: [{ required: true, message: '请输入封面URL', trigger: 'blur' }],
}

const loadData = async () => {
  try {
    const [detailRes, tagsRes, catRes] = await Promise.all([
      getPhotosetDetail(photosetId),
      getTags(),
      getPublicCategories()
    ])
    const ps = detailRes.data
    form.title = ps.title
    form.cover = ps.cover
    form.description = ps.description || ''
    form.is_free = ps.is_free
    form.price = ps.price || 0
    form.status = ps.status
    form.category = ps.category || ''
    form.tags = (ps.tags || []).map(t => t.name)
    form.photos = (ps.photos || []).map(p => ({ url: p.url, sort_order: p.sort_order }))
    availableTags.value = tagsRes.data || []
    availableCategories.value = catRes.data || []
  } catch {
    ElMessage.error('加载失败')
    router.back()
  } finally {
    pageLoading.value = false
  }
}

const handleCoverUpload = async ({ file }) => {
  try {
    const res = await uploadImage(file)
    form.cover = res.data.url
    ElMessage.success('封面上传成功')
  } catch { ElMessage.error('上传失败') }
}

const handlePhotoUpload = async ({ file }) => {
  try {
    const res = await uploadImage(file)
    form.photos.push({ url: res.data.url, sort_order: form.photos.length })
    ElMessage.success('图片上传成功')
  } catch { ElMessage.error('上传失败') }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try { await formRef.value.validate() } catch { return }
  loading.value = true
  try {
    await updatePhotoset(photosetId, {
      title: form.title, cover: form.cover, description: form.description || '',
      tags: form.tags, is_free: form.is_free,
      price: form.is_free === 1 ? 0 : (form.price || 0),
      photos: form.photos.filter(p => p.url.trim()),
      status: form.status,
      category: form.category || ''
    })
    ElMessage.success('保存成功')
    router.back()
  } catch { /* 错误由 request.js 全局处理 */ }
  finally { loading.value = false }
}

onMounted(loadData)
</script>
