<template>
  <div class="create-page">
    <div class="page-header">
      <h1>创建套图</h1>
      <p>发布您的摄影作品到平台</p>
    </div>

    <div class="create-form-wrapper">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        class="create-form"
      >
        <!-- 标题 -->
        <el-form-item label="套图标题" prop="title">
          <el-input
            v-model="form.title"
            placeholder="请输入套图标题"
            maxlength="200"
            show-word-limit
          />
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
            <el-input
              v-model="form.cover"
              placeholder="或手动输入封面URL"
              style="margin-top: 8px"
            >
              <template #append>
                <el-button @click="previewCover = true">预览</el-button>
              </template>
            </el-input>
          </div>
          <div v-if="form.cover" class="cover-preview">
            <el-image :src="form.cover" fit="cover" />
          </div>
        </el-form-item>

        <!-- 描述 -->
        <el-form-item label="套图描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请输入套图描述（选填）"
            maxlength="1000"
            show-word-limit
          />
        </el-form-item>

        <!-- 标签 -->
        <el-form-item label="标签">
          <el-select
            v-model="form.tags"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="选择或输入标签"
            style="width: 100%"
          >
            <el-option
              v-for="tag in availableTags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.name"
            />
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

        <!-- 价格（付费时显示） -->
        <el-form-item v-if="form.is_free === 0" label="价格" prop="price">
          <el-input-number
            v-model="form.price"
            :min="0.01"
            :max="9999"
            :precision="2"
            :step="1"
            placeholder="请输入价格"
          >
            <template #append>元</template>
          </el-input-number>
        </el-form-item>

        <!-- 图片列表 -->
        <el-form-item label="图片列表">
          <div class="photos-editor">
            <div
              v-for="(photo, index) in form.photos"
              :key="index"
              class="photo-item"
            >
              <div class="photo-number">{{ index + 1 }}</div>
              <el-input
                v-model="photo.url"
                placeholder="图片URL"
                style="flex: 1"
              />
              <el-input-number
                v-model="photo.sort_order"
                :min="0"
                :max="999"
                placeholder="排序"
                controls-position="right"
                style="width: 100px"
              />
              <el-button
                type="danger"
                :icon="Delete"
                circle
                @click="removePhoto(index)"
              />
            </div>
            <div class="photo-actions">
              <el-upload
                action=""
                :http-request="handlePhotoUpload"
                :show-file-list="false"
                accept="image/*"
                multiple
              >
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

        <!-- 提交按钮 -->
        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" @click="handleSubmit">
            {{ loading ? '提交中...' : '提交' }}
          </el-button>
          <el-button size="large" @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 封面预览弹窗 -->
    <el-dialog v-model="previewCover" title="封面预览" width="600px">
      <el-image v-if="form.cover" :src="form.cover" fit="contain" style="width: 100%" />
      <div v-else>暂无封面</div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { createPhotoset, getTags, uploadImage, getCategories } from '@/api'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { Delete, Plus, UploadFilled } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref(null)
const loading = ref(false)
const previewCover = ref(false)
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
  is_free: [
    { required: true, message: '请选择收费方式', trigger: 'change' }
  ],
  price: [
    {
      validator: (rule, value, callback) => {
        if (form.is_free === 0 && (!value || value <= 0)) {
          callback(new Error('付费套图请设置价格'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  status: [
    { required: true, message: '请选择发布状态', trigger: 'change' }
  ]
}

// 加载可用标签
const loadTags = async () => {
  try {
    const res = await getTags()
    availableTags.value = res.data || []
  } catch (e) {
    console.error('加载标签失败', e)
  }
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const res = await getCategories()
    availableCategories.value = res.data || []
  } catch (e) {
    console.error('加载分类失败', e)
  }
}

// 添加图片
const addPhoto = () => {
  form.photos.push({
    url: '',
    sort_order: form.photos.length
  })
}

// 触发封面上传
const triggerCoverUpload = () => coverInputRef.value?.click()

// 处理封面上传
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
  } catch (err) {
    ElMessage.error('封面上传失败')
  }
  e.target.value = '' // 重置 input
}

// 处理图片上传
const handlePhotoUpload = async ({ file }) => {
  try {
    const res = await uploadImage(file)
    form.photos.push({
      url: res.data.url,
      sort_order: form.photos.length
    })
    ElMessage.success('图片上传成功')
  } catch (err) {
    ElMessage.error('图片上传失败')
  }
}

// 移除图片
const removePhoto = (index) => {
  form.photos.splice(index, 1)
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch (e) {
    return
  }

  // 过滤掉空 URL 的图片
  const validPhotos = form.photos.filter(p => p.url.trim())

  loading.value = true

  try {
    const data = {
      title: form.title,
      cover: form.cover,
      description: form.description || '',
      tags: form.tags,
      is_free: form.is_free,
      price: form.is_free === 1 ? 0 : (form.price || 0),
      photos: validPhotos,
      status: form.status,
      category: form.category || ''
    }

    const res = await createPhotoset(data)

    ElMessage.success('创建成功')

    // 跳转到详情页
    router.push(`/photoset/${res.data.id}`)
  } catch (e) {
    console.error('创建失败', e)
  } finally {
    loading.value = false
  }
}

// 权限检查
onMounted(() => {
  if (!userStore.isCreatorOrAdmin) {
    ElMessage.error('您没有权限访问此页面')
    router.push('/')
    return
  }
  loadTags()
  loadCategories()
})
</script>

<style scoped>
.create-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 0 60px;
}

.page-header {
  text-align: center;
  padding: 40px 0;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.page-header p {
  color: #666;
  font-size: 15px;
}

.create-form-wrapper {
  background: #fff;
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.create-form :deep(.el-form-item__label) {
  font-weight: 500;
}

.cover-preview {
  margin-top: 12px;
  width: 200px;
  height: 133px;
  border-radius: 8px;
  overflow: hidden;
}

.cover-preview :deep(.el-image) {
  width: 100%;
  height: 100%;
}

.cover-upload-area {
  display: flex;
  flex-direction: column;
}

.cover-preview-box {
  width: 200px;
  height: 133px;
  border: 2px dashed #dcdfe6;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 0.2s;
}

.cover-preview-box:hover {
  border-color: #409eff;
}

.cover-preview-box .el-image {
  width: 100%;
  height: 100%;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #909399;
  gap: 8px;
}

.upload-placeholder span {
  font-size: 12px;
}

.photo-actions {
  display: flex;
  gap: 12px;
  margin-top: 8px;
}

.photos-editor {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.photo-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.photo-number {
  width: 28px;
  height: 28px;
  background: #409eff;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .create-form-wrapper {
    padding: 20px;
  }

  .create-form :deep(.el-form-item__label) {
    float: none;
    text-align: left;
    display: block;
  }

  .create-form :deep(.el-form-item__content) {
    margin-left: 0 !important;
    display: block;
  }

  .photo-item {
    flex-wrap: wrap;
  }
}
</style>
