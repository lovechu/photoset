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

        <!-- 图片列表（WordPress 网格风格） -->
        <el-form-item label="图片列表">
          <div class="photos-grid">
            <!-- 已上传图片 -->
            <div
              v-for="(photo, index) in form.photos"
              :key="photo.tempId || index"
              class="photo-card"
              :class="{ selected: selectedPhotos.includes(photo.tempId || index) }"
              @click="togglePhotoSelect(photo.tempId || index)"
            >
              <div class="photo-thumb">
                <el-image :src="photo.url" fit="cover" />
                <div class="photo-check">
                  <el-icon><Check /></el-icon>
                </div>
                <div class="photo-delete" @click.stop="removePhoto(index)">
                  <el-icon><Close /></el-icon>
                </div>
              </div>
              <div class="photo-meta">
                <span class="photo-index">{{ index + 1 }}</span>
                <span class="photo-name">{{ getFileName(photo.url) }}</span>
              </div>
            </div>

            <!-- 上传中的图片 -->
            <div
              v-for="item in uploadQueue"
              :key="item.tempId"
              class="photo-card uploading"
            >
              <div class="photo-thumb">
                <img :src="item.preview" class="preview-img" />
                <div class="upload-overlay">
                  <el-progress
                    type="circle"
                    :percentage="item.progress"
                    :width="40"
                    :stroke-width="3"
                  />
                </div>
              </div>
              <div class="photo-meta">
                <span class="photo-index">--</span>
                <span class="photo-name uploading-text">{{ item.name }}</span>
              </div>
            </div>

            <!-- 上传按钮 -->
            <el-upload
              action=""
              :http-request="handlePhotoUpload"
              :show-file-list="false"
              accept="image/*"
              multiple
            >
              <div class="photo-card add-btn">
                <el-icon :size="32"><Plus /></el-icon>
                <span>添加图片</span>
              </div>
            </el-upload>
          </div>

          <!-- 批量操作栏 -->
          <div v-if="selectedPhotos.length > 0" class="batch-bar">
            <span>已选 {{ selectedPhotos.length }} 张</span>
            <el-button size="small" @click="selectedPhotos = []">取消选择</el-button>
            <el-button size="small" type="danger" @click="batchDelete">删除选中</el-button>
          </div>

          <!-- 手动添加 URL -->
          <div class="manual-add" style="margin-top: 12px">
            <el-input v-model="manualUrl" placeholder="或手动输入图片URL" style="width: 300px; margin-right: 8px" />
            <el-button @click="addManualUrl">添加</el-button>
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
import { Delete, Plus, UploadFilled, Close, Check } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref(null)
const loading = ref(false)
const previewCover = ref(false)
const availableTags = ref([])
const availableCategories = ref([])
const coverInputRef = ref(null)
const isUploading = ref(false)
const selectedPhotos = ref([])
const manualUrl = ref('')

// uploadQueue 中每个 item: { tempId, name, preview, progress, file }
const uploadQueue = ref([])

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

// 生成唯一 ID
const genId = () => Date.now() + '-' + Math.random().toString(36).slice(2)

// 生成预览图
const genPreview = (file) => {
  return new Promise((resolve) => {
    const reader = new FileReader()
    reader.onload = (e) => resolve(e.target.result)
    reader.readAsDataURL(file)
  })
}

// 处理图片上传（串行队列 + 进度条）
const handlePhotoUpload = async ({ file, fileList }) => {
  const files = fileList && fileList.length > 1 ? fileList : [file]

  // 先加入队列显示预览
  for (const f of files) {
    const preview = await genPreview(f)
    const tempId = genId()
    uploadQueue.value.push({ tempId, name: f.name, preview, progress: 0, file: f })
  }

  if (isUploading.value) return
  isUploading.value = true

  while (uploadQueue.value.length > 0) {
    const item = uploadQueue.value[0]
    try {
      const res = await uploadImage(item.file)
      // 移除队列，加入已上传列表
      uploadQueue.value.shift()
      form.photos.push({
        tempId: item.tempId,
        url: res.data.url,
        sort_order: form.photos.length
      })
      ElMessage.success(`${item.name} 上传成功`)
    } catch {
      ElMessage.error(`${item.name} 上传失败`)
      uploadQueue.value.shift() // 失败也移出，避免卡住
    }
  }

  isUploading.value = false
}

// 切换选中
const togglePhotoSelect = (id) => {
  const idx = selectedPhotos.value.indexOf(id)
  if (idx > -1) selectedPhotos.value.splice(idx, 1)
  else selectedPhotos.value.push(id)
}

// 批量删除
const batchDelete = () => {
  const toDelete = new Set(selectedPhotos.value)
  form.photos = form.photos.filter(p => !toDelete.has(p.tempId || p))
  selectedPhotos.value = []
  ElMessage.success('已删除')
}

// 从文件名提取
const getFileName = (url) => {
  if (!url) return ''
  const parts = url.split('/')
  return parts[parts.length - 1]
}

// 手动添加 URL
const addManualUrl = () => {
  if (!manualUrl.value.trim()) return
  form.photos.push({ tempId: genId(), url: manualUrl.value.trim(), sort_order: form.photos.length })
  manualUrl.value = ''
}

// 移除单张
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

/* WordPress 风格图片网格 */
.photos-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
  margin-bottom: 12px;
}

.photo-card {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
  background: #f5f5f5;
}

.photo-card:hover {
  border-color: #409eff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.photo-card.selected {
  border-color: #409eff;
}

.photo-thumb {
  position: relative;
  width: 100%;
  aspect-ratio: 1;
  overflow: hidden;
}

.photo-thumb .el-image {
  width: 100%;
  height: 100%;
}

.photo-check {
  position: absolute;
  top: 6px;
  left: 6px;
  width: 22px;
  height: 22px;
  background: #409eff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.2s;
}

.photo-card.selected .photo-check {
  opacity: 1;
}

.photo-delete {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 22px;
  height: 22px;
  background: rgba(0,0,0,0.5);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.2s;
}

.photo-card:hover .photo-delete {
  opacity: 1;
}

.photo-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  font-size: 12px;
  color: #666;
}

.photo-index {
  width: 18px;
  height: 18px;
  background: #409eff;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  flex-shrink: 0;
}

.photo-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.uploading-text {
  color: #999;
  font-style: italic;
}

/* 上传中状态 */
.photo-card.uploading .photo-thumb {
  background: #f0f0f0;
}

.preview-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.upload-overlay {
  position: absolute;
  inset: 0;
  background: rgba(255,255,255,0.85);
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 添加按钮 */
.add-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 100%;
  min-height: 150px;
  border: 2px dashed #dcdfe6;
  color: #909399;
  transition: all 0.2s;
}

.add-btn:hover {
  border-color: #409eff;
  color: #409eff;
}

.add-btn span {
  font-size: 12px;
}

/* 批量操作栏 */
.batch-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  background: #f0f9ff;
  border-radius: 8px;
  color: #409eff;
  font-size: 13px;
  margin-bottom: 8px;
}

.batch-bar span {
  flex: 1;
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
