<template>
  <div class="create-page">
    <div class="page-header">
      <h1>编辑套图</h1>
      <p>修改您的摄影作品信息</p>
    </div>

    <div v-if="pageLoading" class="loading-wrapper">
      <el-skeleton :rows="8" animated />
    </div>

    <div v-else class="create-form-wrapper">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        class="create-form"
      >
        <!-- 标题 -->
        <el-form-item label="套图标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入套图标题" maxlength="200" show-word-limit />
        </el-form-item>

        <!-- 封面图 -->
        <el-form-item label="封面图" prop="cover">
          <div class="cover-upload-area">
            <div class="cover-preview-box" @click="triggerCoverUpload">
              <el-image v-if="form.cover" :src="form.cover" fit="cover" />
              <div v-else class="upload-placeholder">
                <el-icon :size="32"><Plus /></el-icon>
                <span>点击上传封面</span>
              </div>
            </div>
            <input ref="coverInputRef" type="file" accept="image/*" hidden @change="handleCoverUpload" />
            <el-input v-model="form.cover" placeholder="或手动输入封面URL" style="margin-top: 8px" />
          </div>
        </el-form-item>

        <!-- 描述 -->
        <el-form-item label="套图描述">
          <el-input v-model="form.description" type="textarea" :rows="4" placeholder="请输入套图描述（选填）" maxlength="1000" show-word-limit />
        </el-form-item>

        <!-- 标签 -->
        <el-form-item label="标签">
          <el-select v-model="form.tags" multiple filterable allow-create default-first-option placeholder="选择或输入标签" style="width: 100%">
            <el-option v-for="tag in availableTags" :key="tag.id" :label="tag.name" :value="tag.name" />
          </el-select>
        </el-form-item>

        <!-- 分类 -->
        <el-form-item label="分类">
          <el-select v-model="form.category" placeholder="选择分类（可选）" clearable style="width: 100%">
            <el-option v-for="cat in availableCategories" :key="cat.slug" :label="cat.name" :value="cat.slug" />
          </el-select>
        </el-form-item>

        <!-- 免费/付费 -->
        <el-form-item label="收费方式" prop="is_free">
          <el-radio-group v-model="form.is_free">
            <el-radio :value="1">免费</el-radio>
            <el-radio :value="0">付费</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 价格 -->
        <el-form-item v-if="form.is_free === 0" label="价格" prop="price">
          <el-input-number v-model="form.price" :min="0.01" :max="9999" :precision="2" :step="1" />
        </el-form-item>

        <!-- 图片列表 -->
        <el-form-item label="图片列表">
          <div class="photos-editor">
            <div v-for="(photo, index) in form.photos" :key="index" class="photo-item">
              <div class="photo-number">{{ index + 1 }}</div>
              <el-input v-model="photo.url" placeholder="图片URL" style="flex: 1" />
              <el-input-number v-model="photo.sort_order" :min="0" :max="999" controls-position="right" style="width: 100px" />
              <el-button type="danger" :icon="Delete" circle @click="removePhoto(index)" />
            </div>
            <div class="photo-actions">
              <el-upload action="" :http-request="handlePhotoUpload" :show-file-list="false" accept="image/*" multiple>
                <el-button type="primary" plain :icon="UploadFilled">上传图片</el-button>
              </el-upload>
              <el-button plain :icon="Plus" @click="addPhoto">手动添加URL</el-button>
            </div>
          </div>
        </el-form-item>

        <!-- 发布状态 -->
        <el-form-item label="发布状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio value="published">直接发布</el-radio>
            <el-radio value="pending">待审核</el-radio>
            <el-radio value="draft">存为草稿</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 按钮 -->
        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" @click="handleSubmit">
            {{ loading ? '保存中...' : '保存修改' }}
          </el-button>
          <el-button size="large" @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getPhotosetDetail, updatePhotoset, getTags, uploadImage, getCategories } from '@/api'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { Delete, Plus, UploadFilled } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const photosetId = Number(route.params.id)

const formRef = ref(null)
const loading = ref(false)
const pageLoading = ref(true)
const availableTags = ref([])
const availableCategories = ref([])
const coverInputRef = ref(null)

const form = reactive({
  title: '',
  cover: '',
  description: '',
  tags: [],
  is_free: 1,
  price: 0,
  category: '',
  photos: [],
  status: 'published'
})

const rules = {
  title: [
    { required: true, message: '请输入套图标题', trigger: 'blur' },
    { max: 200, message: '标题不能超过200个字符', trigger: 'blur' }
  ],
  cover: [
    { required: true, message: '请输入封面图URL', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (!value) {
          callback()
          return
        }
        
        // 允许相对路径（以 /uploads/ 开头的路径）
        if (value.startsWith('/uploads/')) {
          callback()
          return
        }
        
        // 也允许完整的URL
        try {
          const url = new URL(value)
          if (['http:', 'https:'].includes(url.protocol)) {
            callback()
            return
          }
        } catch {
          // 不是有效的完整URL
        }
        
        callback(new Error('请输入有效的URL地址或以/uploads/开头的相对路径'))
      },
      trigger: 'blur'
    }
  ],
  status: [{ required: true, message: '请选择发布状态', trigger: 'change' }]
}

// 加载现有数据回填
const loadData = async () => {
  try {
    const [detailRes, tagsRes, catRes] = await Promise.all([
      getPhotosetDetail(photosetId),
      getTags(),
      getCategories()
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
  } catch (e) {
    ElMessage.error('加载套图数据失败')
    router.back()
  } finally {
    pageLoading.value = false
  }
}

const triggerCoverUpload = () => coverInputRef.value?.click()

const handleCoverUpload = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  try {
    const res = await uploadImage(file)
    form.cover = res.data.url
    ElMessage.success('封面上传成功')
    
    // 手动触发验证，确保验证器立即运行
    if (formRef.value) {
      formRef.value.validateField('cover')
    }
  } catch {
    ElMessage.error('封面上传失败')
  }
  e.target.value = ''
}

const handlePhotoUpload = async ({ file }) => {
  try {
    const res = await uploadImage(file)
    form.photos.push({ url: res.data.url, sort_order: form.photos.length })
    ElMessage.success('图片上传成功')
  } catch {
    ElMessage.error('图片上传失败')
  }
}

const addPhoto = () => form.photos.push({ url: '', sort_order: form.photos.length })
const removePhoto = (index) => form.photos.splice(index, 1)

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch { return }

  const validPhotos = form.photos.filter(p => p.url.trim())
  loading.value = true
  try {
    await updatePhotoset(photosetId, {
      title: form.title,
      cover: form.cover,
      description: form.description || '',
      tags: form.tags,
      is_free: form.is_free,
      price: form.is_free === 1 ? 0 : (form.price || 0),
      photos: validPhotos,
      status: form.status,
      category: form.category || ''
    })
    ElMessage.success('保存成功')
    router.push(`/photoset/${photosetId}`)
  } catch (e) {
    console.error('保存失败', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (!userStore.isCreatorOrAdmin) {
    ElMessage.error('您没有权限访问此页面')
    router.push('/')
    return
  }
  loadData()
})
</script>

<style scoped>
/* 复用 CreatePhotoset.vue 的全部样式 */
.create-page { max-width: 800px; margin: 0 auto; padding: 0 0 60px; }
.page-header { text-align: center; padding: 40px 0; }
.page-header h1 { font-size: 28px; font-weight: 600; color: #1a1a1a; margin-bottom: 8px; }
.page-header p { color: #666; font-size: 15px; }
.create-form-wrapper { background: #fff; border-radius: 12px; padding: 32px; box-shadow: 0 2px 8px rgba(0,0,0,.04); }
.loading-wrapper { padding: 40px; }
.cover-upload-area { display: flex; flex-direction: column; }
.cover-preview-box { width: 200px; height: 133px; border: 2px dashed #dcdfe6; border-radius: 8px; overflow: hidden; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.cover-preview-box:hover { border-color: #409eff; }
.cover-preview-box .el-image { width: 100%; height: 100%; }
.upload-placeholder { display: flex; flex-direction: column; align-items: center; color: #909399; gap: 8px; }
.upload-placeholder span { font-size: 12px; }
.photos-editor { display: flex; flex-direction: column; gap: 12px; }
.photo-item { display: flex; align-items: center; gap: 12px; }
.photo-number { width: 28px; height: 28px; background: #409eff; color: #fff; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 12px; font-weight: 600; flex-shrink: 0; }
.photo-actions { display: flex; gap: 12px; margin-top: 8px; }
</style>
