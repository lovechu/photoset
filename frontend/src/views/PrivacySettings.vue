<template>
  <div class="privacy-settings-page">
    <el-card v-loading="loading">
      <template #header><h3>隐私设置</h3></template>
      <div class="settings-list">
        <div v-for="item in settingItems" :key="item.key" class="setting-item">
          <div class="setting-info">
            <div class="setting-title">{{ item.title }}</div>
            <div class="setting-desc">{{ item.desc }}</div>
          </div>
          <el-switch v-model="form[item.key]" @change="markDirty" />
        </div>
      </div>
      <div class="settings-footer">
        <el-button type="primary" :loading="saving" :disabled="!dirty" @click="save">保存设置</el-button>
        <span v-if="dirty" class="dirty-tip">有未保存的修改</span>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getPrivacySettings, updatePrivacySettings } from '@/api'

const settingItems = [
  { key: 'show_profile', title: '公开个人资料', desc: '其他人可以查看我的个人主页资料' },
  { key: 'show_posts', title: '公开我的帖子', desc: '其他人可以查看我发布的帖子' },
  { key: 'show_favorites', title: '公开我的收藏', desc: '其他人可以查看我的收藏列表' },
  { key: 'allow_search', title: '允许被搜索', desc: '允许其他人通过搜索找到我' },
  { key: 'allow_message', title: '允许私信', desc: '允许其他人给我发送私信' }
]

const form = ref({
  show_profile: true, show_posts: true, show_favorites: false,
  allow_search: true, allow_message: true
})
const loading = ref(false)
const saving = ref(false)
const dirty = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const res = await getPrivacySettings()
    const data = res?.data || {}
    settingItems.forEach(item => {
      if (typeof data[item.key] === 'boolean') form.value[item.key] = data[item.key]
    })
  } catch (e) { console.warn('加载隐私设置失败:', e.message) }
  finally { loading.value = false }
})

function markDirty() { dirty.value = true }

async function save() {
  saving.value = true
  try {
    await updatePrivacySettings({ ...form.value })
    dirty.value = false
    ElMessage.success('设置已保存')
  } catch (e) { ElMessage.error('保存失败') }
  finally { saving.value = false }
}
</script>

<style scoped>
.privacy-settings-page { max-width: 760px; margin: 0 auto; padding: 24px 0 60px; }
.el-card__header h3 { margin: 0; font-size: 18px; }
.setting-item { display: flex; align-items: center; justify-content: space-between; padding: 18px 0; border-bottom: 1px solid #f0f0f0; }
.setting-title { font-weight: 500; }
.setting-desc { color: #999; font-size: 13px; margin-top: 4px; }
.settings-footer { display: flex; align-items: center; gap: 12px; padding-top: 20px; }
.dirty-tip { color: #e6a23c; font-size: 13px; }
</style>
