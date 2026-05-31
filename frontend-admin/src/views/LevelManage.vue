<template>
  <div class="level-manage">
    <div class="header-bar">
      <h2>等级配置管理</h2>
      <span class="subtitle">管理用户等级系统配置（预置10个等级，仅支持编辑）</span>
    </div>

    <el-table :data="levelList" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="level" label="等级" width="70" align="center">
        <template #default="{ row }">
          <el-tag :color="row.color" effect="dark" size="small" style="color: white">
            {{ row.level }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="等级名称" min-width="120" />
      <el-table-column label="图标" width="100" align="center">
        <template #default="{ row }">
          <span v-if="row.icon" style="font-size: 20px;">{{ row.icon }}</span>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="颜色" width="80" align="center">
        <template #default="{ row }">
          <div class="color-preview" :style="{ backgroundColor: row.color }"></div>
        </template>
      </el-table-column>
      <el-table-column label="积分范围" width="150" align="center">
        <template #default="{ row }">
          <span>{{ row.min_points }} - {{ row.max_points }}</span>
        </template>
      </el-table-column>
      <el-table-column label="权限" min-width="300">
        <template #default="{ row }">
          <div class="permissions">
            <el-tag v-if="row.can_create_post" type="success" size="small">发帖</el-tag>
            <el-tag v-if="row.can_create_reply" type="success" size="small">回帖</el-tag>
            <el-tag v-if="row.can_upload_image" type="success" size="small">上传图片</el-tag>
            <el-tag v-if="row.can_upload_video" type="success" size="small">上传视频</el-tag>
            <el-tag v-if="row.can_create_topic" type="success" size="small">创建话题</el-tag>
            <el-tag v-if="row.can_pin_post" type="warning" size="small">置顶</el-tag>
            <el-tag v-if="row.can_delete_reply" type="danger" size="small">删除回帖</el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="每日限制" width="200" align="center">
        <template #default="{ row }">
          <div class="limits">
            <span>发帖: {{ row.max_post_per_day }}</span>
            <span>回帖: {{ row.max_reply_per_day }}</span>
            <span>图片: {{ row.max_image_per_post }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="升级奖励" min-width="150">
        <template #default="{ row }">
          <div v-if="row.reward_points || row.reward_badge">
            <span v-if="row.reward_points">+{{ row.reward_points }}积分</span>
            <span v-if="row.reward_badge">+{{ row.reward_badge }}徽章</span>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" align="center" fixed="right">
        <template #default="{ row }">
          <el-button
            type="primary"
            size="small"
            plain
            @click="handleEdit(row)"
          >
            编辑
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 编辑等级配置对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="'编辑等级配置 - ' + (editForm.name || '')"
      width="600px"
      @close="resetForm"
    >
      <el-form ref="formRef" :model="editForm" :rules="formRules" label-width="100px">
        <el-form-item label="等级名称" prop="name">
          <el-input v-model="editForm.name" placeholder="请输入等级名称" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-input v-model="editForm.icon" placeholder="请输入图标名称（如 Star）" />
          <div class="icon-hint">
            <span>可用图标: Star, Trophy, Medal, Crown, Diamond, Award, Gem, Heart, Like, Fire</span>
          </div>
        </el-form-item>
        <el-form-item label="颜色" prop="color">
          <el-color-picker v-model="editForm.color" />
          <span class="color-text">{{ editForm.color }}</span>
        </el-form-item>
        <el-form-item label="最小积分" prop="min_points">
          <el-input-number v-model="editForm.min_points" :min="0" :max="editForm.max_points" />
        </el-form-item>
        <el-form-item label="最大积分" prop="max_points">
          <el-input-number v-model="editForm.max_points" :min="editForm.min_points" />
        </el-form-item>
        
        <el-divider content-position="left">权限设置</el-divider>
        
        <el-form-item label="发帖权限" prop="can_create_post">
          <el-switch v-model="editForm.can_create_post" />
        </el-form-item>
        <el-form-item label="回帖权限" prop="can_create_reply">
          <el-switch v-model="editForm.can_create_reply" />
        </el-form-item>
        <el-form-item label="上传图片" prop="can_upload_image">
          <el-switch v-model="editForm.can_upload_image" />
        </el-form-item>
        <el-form-item label="上传视频" prop="can_upload_video">
          <el-switch v-model="editForm.can_upload_video" />
        </el-form-item>
        <el-form-item label="创建话题" prop="can_create_topic">
          <el-switch v-model="editForm.can_create_topic" />
        </el-form-item>
        <el-form-item label="置顶帖子" prop="can_pin_post">
          <el-switch v-model="editForm.can_pin_post" />
        </el-form-item>
        <el-form-item label="删除回帖" prop="can_delete_reply">
          <el-switch v-model="editForm.can_delete_reply" />
        </el-form-item>
        
        <el-divider content-position="left">每日限制</el-divider>
        
        <el-form-item label="每日发帖" prop="max_post_per_day">
          <el-input-number v-model="editForm.max_post_per_day" :min="0" />
          <span class="limit-hint">0表示无限制</span>
        </el-form-item>
        <el-form-item label="每日回帖" prop="max_reply_per_day">
          <el-input-number v-model="editForm.max_reply_per_day" :min="0" />
          <span class="limit-hint">0表示无限制</span>
        </el-form-item>
        <el-form-item label="每帖图片" prop="max_image_per_post">
          <el-input-number v-model="editForm.max_image_per_post" :min="0" />
          <span class="limit-hint">0表示无限制</span>
        </el-form-item>
        <el-form-item label="每帖视频" prop="max_video_per_post">
          <el-input-number v-model="editForm.max_video_per_post" :min="0" />
          <span class="limit-hint">0表示无限制</span>
        </el-form-item>
        <el-form-item label="帖子长度" prop="max_post_length">
          <el-input-number v-model="editForm.max_post_length" :min="0" />
          <span class="limit-hint">0表示无限制</span>
        </el-form-item>
        
        <el-divider content-position="left">升级奖励</el-divider>
        
        <el-form-item label="奖励积分" prop="reward_points">
          <el-input-number v-model="editForm.reward_points" :min="0" />
        </el-form-item>
        <el-form-item label="奖励徽章" prop="reward_badge">
          <el-input v-model="editForm.reward_badge" placeholder="请输入徽章名称" />
        </el-form-item>
        <el-form-item label="奖励称号" prop="reward_title">
          <el-input v-model="editForm.reward_title" placeholder="请输入称号名称" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          保存修改
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, computed } from 'vue'
import { getLevelConfigs, updateLevelConfig } from '@/api/level'
import { ElMessage } from 'element-plus'
const loading = ref(false)
const levelList = ref([])
const dialogVisible = ref(false)
const submitting = ref(false)

const editForm = reactive({
  id: null,
  level: 0,
  name: '',
  icon: '',
  color: '#409EFF',
  min_points: 0,
  max_points: 0,
  can_create_post: false,
  can_create_reply: false,
  can_upload_image: false,
  can_upload_video: false,
  can_create_topic: false,
  can_pin_post: false,
  can_delete_reply: false,
  max_post_per_day: 0,
  max_reply_per_day: 0,
  max_image_per_post: 0,
  max_video_per_post: 0,
  max_post_length: 0,
  reward_points: 0,
  reward_badge: '',
  reward_title: ''
})

const formRules = {
  name: [
    { required: true, message: '请输入等级名称', trigger: 'blur' },
    { min: 2, max: 20, message: '名称长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  icon: [
    { required: true, message: '请输入图标名称', trigger: 'blur' }
  ],
  color: [
    { required: true, message: '请选择颜色', trigger: 'change' }
  ],
  min_points: [
    { required: true, message: '请输入最小积分', trigger: 'blur' }
  ],
  max_points: [
    { required: true, message: '请输入最大积分', trigger: 'blur' }
  ]
}

const formRef = ref()

async function fetchLevels() {
  loading.value = true
  try {
    const res = await getLevelConfigs()
    levelList.value = res.data || []
  } catch (error) {
    console.error('获取等级配置失败:', error)
    ElMessage.error('获取等级配置失败')
  } finally {
    loading.value = false
  }
}

function handleEdit(level) {
  Object.assign(editForm, {
    id: level.id,
    level: level.level,
    name: level.name,
    icon: level.icon || '',
    color: level.color || '#409EFF',
    min_points: level.min_points || 0,
    max_points: level.max_points || 0,
    can_create_post: level.can_create_post || false,
    can_create_reply: level.can_create_reply || false,
    can_upload_image: level.can_upload_image || false,
    can_upload_video: level.can_upload_video || false,
    can_create_topic: level.can_create_topic || false,
    can_pin_post: level.can_pin_post || false,
    can_delete_reply: level.can_delete_reply || false,
    max_post_per_day: level.max_post_per_day || 0,
    max_reply_per_day: level.max_reply_per_day || 0,
    max_image_per_post: level.max_image_per_post || 0,
    max_video_per_post: level.max_video_per_post || 0,
    max_post_length: level.max_post_length || 0,
    reward_points: level.reward_points || 0,
    reward_badge: level.reward_badge || '',
    reward_title: level.reward_title || ''
  })
  dialogVisible.value = true
}

function resetForm() {
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

async function submitForm() {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true

    const formData = {
      name: editForm.name.trim(),
      icon: editForm.icon.trim(),
      color: editForm.color,
      min_points: editForm.min_points,
      max_points: editForm.max_points,
      can_create_post: editForm.can_create_post,
      can_create_reply: editForm.can_create_reply,
      can_upload_image: editForm.can_upload_image,
      can_upload_video: editForm.can_upload_video,
      can_create_topic: editForm.can_create_topic,
      can_pin_post: editForm.can_pin_post,
      can_delete_reply: editForm.can_delete_reply,
      max_post_per_day: editForm.max_post_per_day,
      max_reply_per_day: editForm.max_reply_per_day,
      max_image_per_post: editForm.max_image_per_post,
      max_video_per_post: editForm.max_video_per_post,
      max_post_length: editForm.max_post_length,
      reward_points: editForm.reward_points,
      reward_badge: editForm.reward_badge.trim() || null,
      reward_title: editForm.reward_title.trim() || null
    }

    await updateLevelConfig(editForm.id, formData)
    ElMessage.success('更新成功')
    dialogVisible.value = false
    fetchLevels()
  } catch (error) {
    if (error.errorFields) {
      // 验证失败，不处理
    } else {
      ElMessage.error('更新失败')
    }
  } finally {
    submitting.value = false
  }
}

onMounted(fetchLevels)
</script>

<style scoped>
.header-bar {
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.header-bar h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.header-bar .subtitle {
  font-size: 14px;
  color: #909399;
}

.color-preview {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  border: 1px solid #dcdfe6;
  margin: 0 auto;
}

.permissions {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.limits {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 12px;
  color: #606266;
}

.icon-hint {
  margin-top: 6px;
  font-size: 12px;
  color: #909399;
}

.color-text {
  margin-left: 12px;
  color: #606266;
}

.limit-hint {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}
</style>