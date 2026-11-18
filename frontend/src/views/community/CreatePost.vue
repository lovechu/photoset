<template>
  <div class="composer-page">
    <div class="composer">
      <!-- 顶部标题栏 -->
      <div class="composer-topbar">
        <span class="composer-title">{{ isEdit ? '编辑话题' : '新建话题' }}</span>
        <div class="composer-top-meta">
          <span class="draft-hint" v-if="form.content">已自动保留内容</span>
          <button class="close-btn" @click="$router.back()" title="关闭">×</button>
        </div>
      </div>

      <!-- 分类 / 类型 / 可见范围 -->
      <div class="composer-options">
        <label class="opt">
          <span class="opt-label">分类</span>
          <el-select v-model="form.category" placeholder="选择分类" size="small" class="cat-select">
            <template #prefix>
              <span class="cat-dot" :style="{ background: currentCategoryColor }"></span>
            </template>
            <el-option v-for="cat in categories" :key="cat.slug" :label="cat.name" :value="cat.slug">
              <span class="cat-dot" :style="{ background: categoryColor(cat.slug) }"></span>
              <span>{{ cat.name }}</span>
            </el-option>
          </el-select>
        </label>
        <label class="opt">
          <span class="opt-label">类型</span>
          <el-select v-model="form.post_type" size="small" class="type-select">
            <el-option label="动态" value="dynamic" />
            <el-option label="文章" value="article" />
            <el-option label="提问" value="question" />
            <el-option label="随手拍" value="quick" />
          </el-select>
        </label>
        <label class="opt">
          <span class="opt-label">可见</span>
          <el-select v-model="form.visibility" size="small" class="vis-select">
            <el-option label="公开" value="public" />
            <el-option label="会员可见" value="member" />
            <el-option label="VIP 可见" value="vip" />
          </el-select>
        </label>
        <label class="opt opt-original">
          <el-switch v-model="form.is_original" size="small" />
          <span>原创</span>
        </label>
      </div>

      <!-- 标题输入 -->
      <div class="title-row">
        <input
          v-model="form.title"
          class="title-input"
          maxlength="100"
          :placeholder="form.post_type === 'article' ? '请输入文章标题（必填）' : '话题标题（选填）'"
        />
      </div>

      <!-- Markdown 工具栏 -->
      <div class="md-toolbar">
        <button class="tool-btn" title="加粗" @click.prevent="wrapSelection('**', '**', '加粗文本')"><b>B</b></button>
        <button class="tool-btn" title="斜体" @click.prevent="wrapSelection('*', '*', '斜体文本')"><i>I</i></button>
        <button class="tool-btn" title="删除线" @click.prevent="wrapSelection('~~', '~~', '删除线')"><s>S</s></button>
        <span class="tool-sep"></span>
        <button class="tool-btn wide" title="标题" @click.prevent="prefixLine('## ', '标题')">H</button>
        <button class="tool-btn" title="引用" @click.prevent="prefixLine('> ', '引用内容')">❝</button>
        <button class="tool-btn" title="链接" @click.prevent="insertLink()">🔗</button>
        <span class="tool-sep"></span>
        <button class="tool-btn wide" title="无序列表" @click.prevent="prefixLine('- ', '列表项')">• 列表</button>
        <button class="tool-btn wide" title="有序列表" @click.prevent="prefixLine('1. ', '列表项')">1. 列表</button>
        <button class="tool-btn wide" title="代码块" @click.prevent="wrapSelection('\n```\n', '\n```\n', '代码')">代码块</button>
        <span class="tool-sep"></span>
        <button class="tool-btn wide" :disabled="uploading" title="插入图片" @click.prevent="imageInput?.click()">▣ 图片</button>
        <button class="tool-btn wide" :disabled="uploadingVideo" title="插入视频" @click.prevent="videoInput?.click()">▶ 视频</button>
        <input ref="imageInput" type="file" accept="image/*" multiple style="display:none" @change="onImagesPicked" />
        <input ref="videoInput" type="file" accept="video/*" style="display:none" @change="onVideoPicked" />
        <span class="tool-insert-tip">{{ uploading ? '上传图片中…' : uploadingVideo ? '上传视频中…' : '支持 Markdown' }}</span>
      </div>

      <!-- 正文 -->
      <div class="editor-body">
        <textarea
          ref="contentArea"
          v-model="form.content"
          class="content-input"
          maxlength="5000"
          placeholder="写下你的话题内容…"
        ></textarea>
      </div>

      <!-- 已插入图片 -->
      <div v-if="uploadedImages.length" class="uploaded-images">
        <div v-for="(url, i) in uploadedImages" :key="i" class="uploaded-img">
          <el-image :src="url" fit="cover" />
          <span class="img-remove" @click="removeUploaded(i)">×</span>
        </div>
      </div>

      <!-- 标签 / 话题 -->
      <div class="composer-fields">
        <div class="field">
          <label class="field-label">标签</label>
          <el-select
            v-model="form.tags"
            multiple filterable allow-create default-first-option
            class="field-select" size="small"
            placeholder="输入标签后回车（可选）"
          />
        </div>
        <div class="field">
          <label class="field-label">话题</label>
          <el-select
            v-model="form.topics"
            multiple filterable allow-create default-first-option
            class="field-select" size="small"
            placeholder="输入话题后回车（可选）"
          />
        </div>
      </div>

      <!-- 底部操作栏 -->
      <div class="composer-footer">
        <div class="footer-left">
          <button class="text-btn" @click="saveDraft">存草稿</button>
          <button class="text-btn" @click="$router.back()">取消</button>
        </div>
        <div class="footer-right">
          <span class="char-count">{{ form.content.length }} / 5000</span>
          <button class="submit-btn" :disabled="submitting" @click="submitPost">
            {{ submitting ? '发表中…' : (isEdit ? '保存修改' : '创建话题') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useCommunityStore } from '@/stores/community'
import {
  createPost, updatePost, uploadCommunityImage, uploadCommunityVideo,
  saveDraft as apiSaveDraft, getPostDetail, getMyDrafts
} from '@/api'

const route = useRoute()
const router = useRouter()
const communityStore = useCommunityStore()

const categories = ref([])
const uploading = ref(false)
const uploadingVideo = ref(false)
const submitting = ref(false)
const imageInput = ref(null)
const videoInput = ref(null)
const contentArea = ref(null)
const uploadedImages = ref([])
const editId = ref(null)
const draftId = ref(null)

const isEdit = ref(false)
const form = ref({
  content: '', title: '', category: '', post_type: 'dynamic',
  visibility: 'public', is_original: true, tags: [], topics: []
})

const catPalette = ['#e07a5f', '#4f6ef7', '#81b29a', '#9b5de5']
function categoryColor(slug) {
  const i = categories.value.findIndex(c => c.slug === slug)
  return catPalette[(i < 0 ? 0 : i) % catPalette.length]
}
const currentCategoryColor = computed(() => categoryColor(form.value.category))

onMounted(async () => {
  categories.value = await communityStore.fetchCategories()
  if (route.query.edit) {
    isEdit.value = true
    editId.value = route.query.edit
    const res = await getPostDetail(editId.value)
    const p = res?.data?.post
    if (p) {
      form.value = {
        content: p.content || '', title: p.title || '', category: p.category || '',
        post_type: p.post_type || 'dynamic', visibility: p.visibility || 'public',
        is_original: !!p.is_original,
        tags: (p.tags || []).map(t => t.name),
        topics: (p.topics || []).map(t => t.name)
      }
    }
  } else if (route.query.draft) {
    // 草稿继续编辑：优先读路由 state 里携带的草稿对象
    const draft = window.history.state?.draft
    if (draft) {
      form.value = {
        content: draft.content || '', title: draft.title || '', category: draft.category || '',
        post_type: draft.post_type || 'dynamic', visibility: draft.visibility || 'public',
        is_original: draft.is_original !== false,
        tags: draft.tags || [], topics: draft.topics || []
      }
      draftId.value = route.query.draft
    } else {
      // 刷新后 state 丢失，从草稿列表回退查找
      try {
        const res = await getMyDrafts({ page: 1, page_size: 100 })
        const found = (res?.data?.drafts || []).find(d => String(d.id) === String(route.query.draft))
        if (found) {
          form.value = {
            content: found.content || '', title: found.title || '', category: found.category || '',
            post_type: found.post_type || 'dynamic', visibility: found.visibility || 'public',
            is_original: found.is_original !== false,
            tags: found.tags || [], topics: found.topics || []
          }
          draftId.value = found.id
        }
      } catch (e) { /* 忽略，使用空白编辑器 */ }
    }
  } else if (categories.value.length && !form.value.category) {
    form.value.category = categories.value[0].slug
  }
})

// ===== Markdown 工具 =====
function selection() {
  const el = contentArea.value
  if (!el) return { start: 0, end: 0 }
  return { start: el.selectionStart, end: el.selectionEnd }
}
function applyEdit(start, end, text, cursorPos) {
  form.value.content = form.value.content.slice(0, start) + text + form.value.content.slice(end)
  nextTick(() => {
    const el = contentArea.value
    if (!el) return
    el.focus()
    const pos = cursorPos != null ? cursorPos : start + text.length
    el.setSelectionRange(pos, pos)
  })
}
function wrapSelection(before, after, placeholder) {
  const { start, end } = selection()
  const selected = form.value.content.slice(start, end) || placeholder
  const text = before + selected + after
  applyEdit(start, end, text, start + before.length + selected.length)
}
function prefixLine(prefix, placeholder) {
  const { start } = selection()
  const lineStart = form.value.content.lastIndexOf('\n', start - 1) + 1
  const rest = form.value.content.slice(lineStart)
  const lineEndIdx = rest.indexOf('\n')
  const lineEnd = lineEndIdx < 0 ? form.value.content.length : lineStart + lineEndIdx
  const line = form.value.content.slice(lineStart, lineEnd) || placeholder
  applyEdit(lineStart, lineEnd, prefix + line, lineStart + prefix.length + line.length)
}
function insertLink() {
  const { start, end } = selection()
  const selected = form.value.content.slice(start, end) || '链接文字'
  const text = `[${selected}](https://)`
  applyEdit(start, end, text, start + selected.length + 3)
}

async function onImagesPicked(e) {
  const files = Array.from(e.target.files || [])
  if (!files.length) return
  uploading.value = true
  try {
    for (const file of files) {
      const res = await uploadCommunityImage(file)
      const url = res?.data?.url || res?.data?.image_url || res?.url
      if (url) {
        uploadedImages.value.push(url)
        form.value.content += `\n![图片](${url})\n`
      }
    }
    ElMessage.success('图片已插入')
  } catch (err) {
    ElMessage.error('图片上传失败')
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}

async function onVideoPicked(e) {
  const file = e.target.files?.[0]
  if (!file) return
  uploadingVideo.value = true
  try {
    const res = await uploadCommunityVideo(file)
    const url = res?.data?.url || res?.data?.video_url || res?.url
    if (url) {
      form.value.content += `\n<video src="${url}" controls></video>\n`
      ElMessage.success('视频已插入')
    }
  } catch (err) {
    ElMessage.error('视频上传失败')
  } finally {
    uploadingVideo.value = false
    e.target.value = ''
  }
}

function removeUploaded(i) {
  const url = uploadedImages.value[i]
  uploadedImages.value.splice(i, 1)
  form.value.content = form.value.content.replace(`![图片](${url})`, '')
}

function buildPayload() {
  return {
    content: form.value.content,
    title: form.value.title || undefined,
    category: form.value.category || undefined,
    post_type: form.value.post_type,
    visibility: form.value.visibility,
    is_original: form.value.is_original,
    tags: form.value.tags,
    topics: form.value.topics
  }
}

async function submitPost() {
  if (!form.value.content.trim()) {
    ElMessage.warning('请输入内容')
    return
  }
  if (!form.value.category) {
    ElMessage.warning('请选择分类')
    return
  }
  if (form.value.post_type === 'article' && !form.value.title.trim()) {
    ElMessage.warning('文章需要标题')
    return
  }
  submitting.value = true
  try {
    const payload = buildPayload()
    if (isEdit.value) {
      await updatePost(editId.value, payload)
      ElMessage.success('已保存')
      router.replace(`/community/post/${editId.value}`)
    } else {
      const res = await createPost(payload)
      const id = res?.data?.id
      ElMessage.success('发表成功')
      router.replace(id ? `/community/post/${id}` : '/community')
    }
  } catch (e) {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

async function saveDraft() {
  try {
    await apiSaveDraft({
      title: form.value.title || '',
      content: form.value.content,
      category: form.value.category || '',
      post_type: form.value.post_type,
      visibility: form.value.visibility
    })
    ElMessage.success('已保存到草稿箱')
  } catch (e) {
    ElMessage.error('保存草稿失败')
  }
}
</script>

<style scoped>
.composer-page {
  max-width: 860px; margin: 0 auto; padding: 24px 20px 60px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  color: #33373d; font-size: 14px; line-height: 1.5;
}
.composer {
  background: #fff; border: 1px solid #e6e8eb; border-radius: 6px; overflow: hidden;
  box-shadow: 0 1px 3px rgba(0,0,0,0.04);
}

/* 顶部标题栏 */
.composer-topbar {
  display: flex; align-items: center; justify-content: space-between;
  background: #2b3038; color: #fff; padding: 13px 20px;
}
.composer-title { font-size: 15px; font-weight: 700; }
.composer-top-meta { display: flex; align-items: center; gap: 14px; }
.draft-hint { font-size: 11px; color: #a9b0b9; }
.close-btn {
  background: none; border: none; color: #c4cad2; font-size: 22px; line-height: 1;
  cursor: pointer; padding: 0 2px;
}
.close-btn:hover { color: #fff; }

/* 选项行 */
.composer-options {
  display: flex; align-items: center; gap: 18px; flex-wrap: wrap;
  padding: 12px 20px; background: #fafbfc; border-bottom: 1px solid #eef0f2;
}
.opt { display: flex; align-items: center; gap: 7px; }
.opt-label { font-size: 12px; color: #8d95a0; }
.cat-dot { display: inline-block; width: 8px; height: 8px; border-radius: 50%; margin-right: 7px; vertical-align: middle; }
.cat-select :deep(.el-select__wrapper) { min-width: 120px; }
.type-select { width: 96px; }
.vis-select { width: 116px; }
.opt-original { gap: 6px; font-size: 12px; color: #6b7280; }
:deep(.el-option) { display: flex; align-items: center; }

/* 标题 */
.title-row { padding: 18px 22px 6px; }
.title-input {
  width: 100%; border: none; outline: none; background: transparent;
  font-size: 24px; font-weight: 700; color: #2b3038; padding: 4px 0;
}
.title-input::placeholder { color: #c2c8d0; }

/* Markdown 工具栏 */
.md-toolbar {
  display: flex; align-items: center; gap: 2px; flex-wrap: wrap;
  margin: 8px 20px 0; padding: 6px 8px;
  background: #f7f8f9; border: 1px solid #e6e8eb; border-bottom: none;
  border-radius: 5px 5px 0 0;
}
.tool-btn {
  min-width: 28px; height: 28px; padding: 0 6px;
  background: none; border: none; border-radius: 3px; cursor: pointer;
  font-size: 13px; color: #5a606a;
}
.tool-btn.wide { padding: 0 9px; font-size: 12px; }
.tool-btn:hover:not(:disabled) { background: #e8ecf1; color: #2b3038; }
.tool-btn:disabled { opacity: 0.5; cursor: default; }
.tool-sep { width: 1px; height: 18px; background: #dfe3e8; margin: 0 5px; }
.tool-insert-tip { margin-left: auto; font-size: 11px; color: #9aa1ac; }

/* 正文 */
.editor-body { margin: 0 20px; }
.content-input {
  display: block; width: 100%; min-height: 300px; resize: vertical;
  border: 1px solid #e6e8eb; border-radius: 0 0 5px 5px;
  padding: 16px; font-size: 14px; line-height: 1.8; color: #33373d;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  outline: none;
}
.content-input:focus { border-color: #4f6ef7; }
.content-input::placeholder { color: #b6bcc5; }

.uploaded-images { display: flex; flex-wrap: wrap; gap: 10px; margin: 14px 22px 0; }
.uploaded-img { position: relative; width: 80px; height: 80px; }
.uploaded-img :deep(.el-image) { width: 100%; height: 100%; border-radius: 4px; }
.img-remove {
  position: absolute; top: -7px; right: -7px; width: 18px; height: 18px;
  background: #2b3038; color: #fff; border-radius: 50%; cursor: pointer;
  font-size: 13px; line-height: 18px; text-align: center;
}

/* 标签 / 话题 */
.composer-fields { display: flex; gap: 16px; padding: 18px 22px; flex-wrap: wrap; }
.field { flex: 1; min-width: 200px; }
.field-label { display: block; font-size: 12px; color: #8d95a0; margin-bottom: 6px; }
.field-select { width: 100%; }

/* 底部操作栏 */
.composer-footer {
  display: flex; align-items: center; justify-content: space-between;
  padding: 14px 20px; background: #fafbfc; border-top: 1px solid #eef0f2;
}
.footer-left { display: flex; gap: 16px; }
.text-btn { background: none; border: none; font-size: 12px; color: #8d95a0; cursor: pointer; padding: 4px 0; }
.text-btn:hover { color: #2b3038; }
.footer-right { display: flex; align-items: center; gap: 14px; }
.char-count { font-size: 11px; color: #b6bcc5; }
.submit-btn {
  padding: 9px 22px; background: #4f6ef7; color: #fff; border: none; border-radius: 4px;
  font-size: 13px; font-weight: 600; cursor: pointer;
}
.submit-btn:hover:not(:disabled) { background: #3d5ae0; }
.submit-btn:disabled { opacity: 0.6; cursor: default; }

@media (max-width: 640px) {
  .composer-page { padding: 12px 8px 40px; }
  .title-input { font-size: 20px; }
  .composer-options { gap: 12px; }
  .content-input { min-height: 240px; }
}
</style>
