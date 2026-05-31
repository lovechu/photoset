<template>
  <div class="achievement-manage">
    <div class="header-bar">
      <el-button type="primary" @click="handleCreate" :icon="Plus">
        新建成就
      </el-button>
      <el-input
        v-model="filterKeyword"
        placeholder="搜索成就名称"
        clearable
        @clear="fetchAchievements"
        @keyup.enter="fetchAchievements"
        style="width: 240px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <el-table :data="filteredAchievements" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column label="图标" width="80" align="center">
        <template #default="{ row }">
          <span style="font-size: 20px;">{{ row.icon }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="标识名" min-width="120" />
      <el-table-column prop="title" label="成就名称" min-width="150" />
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getTypeTagType(row.type)" size="small">{{ getTypeName(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="解锁条件" min-width="200">
        <template #default="{ row }">
          <span>{{ getConditionName(row.condition_type) }} {{ row.condition_value }}</span>
        </template>
      </el-table-column>
      <el-table-column label="奖励" min-width="150">
        <template #default="{ row }">
          <div v-if="row.reward_points || row.reward_title">
            <span v-if="row.reward_points">+{{ row.reward_points }}积分</span>
            <span v-if="row.reward_title">+{{ row.reward_title }}称号</span>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="sort_order" label="排序" width="80" align="center" />
      <el-table-column label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_hidden ? 'info' : 'success'" size="small">
            {{ row.is_hidden ? '隐藏' : '显示' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" align="center" fixed="right">
        <template #default="{ row }">
          <div class="action-buttons">
            <el-button
              type="primary"
              size="small"
              plain
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-popconfirm
              title="确定要删除这个成就吗？删除后用户将无法解锁此成就。"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button
                  type="danger"
                  size="small"
                  plain
                >
                  删除
                </el-button>
              </template>
            </el-popconfirm>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新建/编辑成就对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="resetForm"
    >
      <el-form ref="formRef" :model="editForm" :rules="formRules" label-width="100px">
        <el-form-item label="成就名称" prop="title">
          <el-input v-model="editForm.title" placeholder="请输入成就名称" />
        </el-form-item>
        <el-form-item label="标识名" prop="name">
          <el-input v-model="editForm.name" placeholder="请输入标识名（英文小写字母和下划线）" />
          <div class="slug-hint">
            <span class="slug-preview">用于程序内部标识，如: first_post</span>
          </div>
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-input v-model="editForm.icon" placeholder="请输入表情符号图标" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="editForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入成就描述"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="成就类型" prop="type">
          <el-select v-model="editForm.type" placeholder="请选择成就类型">
            <el-option label="发帖相关" value="post" />
            <el-option label="回复相关" value="reply" />
            <el-option label="点赞相关" value="like" />
            <el-option label="关注相关" value="follow" />
            <el-option label="等级相关" value="level" />
            <el-option label="特殊成就" value="special" />
          </el-select>
        </el-form-item>
        <el-form-item label="条件类型" prop="condition_type">
          <el-select v-model="editForm.condition_type" placeholder="请选择条件类型">
            <el-option label="发帖数量" value="post_count" />
            <el-option label="回复数量" value="reply_count" />
            <el-option label="获赞数量" value="like_received" />
            <el-option label="关注数量" value="following_count" />
            <el-option label="粉丝数量" value="follower_count" />
            <el-option label="达到等级" value="level_reached" />
            <el-option label="特殊条件" value="special" />
          </el-select>
        </el-form-item>
        <el-form-item label="条件值" prop="condition_value">
          <el-input-number v-model="editForm.condition_value" :min="0" />
        </el-form-item>
        
        <el-divider content-position="left">奖励设置</el-divider>
        
        <el-form-item label="奖励积分" prop="reward_points">
          <el-input-number v-model="editForm.reward_points" :min="0" />
        </el-form-item>
        <el-form-item label="奖励称号" prop="reward_title">
          <el-input v-model="editForm.reward_title" placeholder="请输入奖励称号" />
        </el-form-item>
        
        <el-divider content-position="left">显示设置</el-divider>
        
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="editForm.sort_order" :min="0" />
        </el-form-item>
        <el-form-item label="隐藏" prop="is_hidden">
          <el-switch v-model="editForm.is_hidden" />
          <span class="limit-hint">隐藏后用户在成就列表中不可见</span>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ isEditing ? '保存修改' : '创建成就' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, reactive, computed } from 'vue'
import { getAchievements, createAchievement, updateAchievement, deleteAchievement } from '@/api/level'
import { ElMessage } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'

const loading = ref(false)
const achievementList = ref([])
const filterKeyword = ref('')
const submitting = ref(false)
const dialogVisible = ref(false)
const isEditing = ref(false)

const editForm = reactive({
  id: null,
  name: '',
  title: '',
  description: '',
  icon: '',
  type: 'post',
  condition_type: 'post_count',
  condition_value: 0,
  reward_points: 0,
  reward_title: '',
  sort_order: 0,
  is_hidden: false
})

const formRules = {
  title: [
    { required: true, message: '请输入成就名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入标识名', trigger: 'blur' },
    { pattern: /^[a-z0-9_]+$/, message: '只能包含小写字母、数字和下划线', trigger: 'blur' },
    { min: 2, max: 50, message: '标识名长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  icon: [
    { required: true, message: '请输入图标', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择成就类型', trigger: 'change' }
  ],
  condition_type: [
    { required: true, message: '请选择条件类型', trigger: 'change' }
  ],
  condition_value: [
    { required: true, message: '请输入条件值', trigger: 'blur' }
  ]
}

const formRef = ref()

const filteredAchievements = computed(() => {
  if (!filterKeyword.value) return achievementList.value
  const keyword = filterKeyword.value.toLowerCase()
  return achievementList.value.filter(item => 
    item.name.toLowerCase().includes(keyword) ||
    item.title.toLowerCase().includes(keyword) ||
    (item.description && item.description.toLowerCase().includes(keyword))
  )
})

function getTypeName(type) {
  const typeMap = {
    'post': '发帖',
    'reply': '回复',
    'like': '点赞',
    'follow': '关注',
    'level': '等级',
    'special': '特殊'
  }
  return typeMap[type] || type
}

function getTypeTagType(type) {
  const typeMap = {
    'post': 'primary',
    'reply': 'success',
    'like': 'warning',
    'follow': 'info',
    'level': 'danger',
    'special': ''
  }
  return typeMap[type] || ''
}

function getConditionName(conditionType) {
  const conditionMap = {
    'post_count': '发帖数量',
    'reply_count': '回复数量',
    'like_received': '获赞数量',
    'following_count': '关注数量',
    'follower_count': '粉丝数量',
    'level_reached': '达到等级',
    'special': '特殊条件'
  }
  return conditionMap[conditionType] || conditionType
}

async function fetchAchievements() {
  loading.value = true
  try {
    const res = await getAchievements()
    achievementList.value = res.data || []
  } catch (error) {
    console.error('获取成就列表失败:', error)
    ElMessage.error('获取成就列表失败')
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  isEditing.value = false
  resetForm()
  dialogVisible.value = true
}

function handleEdit(achievement) {
  isEditing.value = true
  Object.assign(editForm, {
    id: achievement.id,
    name: achievement.name,
    title: achievement.title,
    description: achievement.description || '',
    icon: achievement.icon || '',
    type: achievement.type || 'post',
    condition_type: achievement.condition_type || 'post_count',
    condition_value: achievement.condition_value || 0,
    reward_points: achievement.reward_points || 0,
    reward_title: achievement.reward_title || '',
    sort_order: achievement.sort_order || 0,
    is_hidden: achievement.is_hidden || false
  })
  dialogVisible.value = true
}

async function handleDelete(achievement) {
  try {
    await deleteAchievement(achievement.id)
    ElMessage.success('删除成功')
    fetchAchievements()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

function resetForm() {
  Object.assign(editForm, {
    id: null,
    name: '',
    title: '',
    description: '',
    icon: '',
    type: 'post',
    condition_type: 'post_count',
    condition_value: 0,
    reward_points: 0,
    reward_title: '',
    sort_order: 0,
    is_hidden: false
  })
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
      title: editForm.title.trim(),
      description: editForm.description.trim() || null,
      icon: editForm.icon.trim(),
      type: editForm.type,
      condition_type: editForm.condition_type,
      condition_value: editForm.condition_value,
      reward_points: editForm.reward_points,
      reward_title: editForm.reward_title.trim() || null,
      sort_order: editForm.sort_order,
      is_hidden: editForm.is_hidden
    }

    if (isEditing.value) {
      await updateAchievement(editForm.id, formData)
      ElMessage.success('更新成功')
    } else {
      await createAchievement(formData)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchAchievements()
  } catch (error) {
    if (error.errorFields) {
      // 验证失败，不处理
    } else {
      ElMessage.error(isEditing.value ? '更新失败' : '创建失败')
    }
  } finally {
    submitting.value = false
  }
}

// 监听筛选条件变化
watch([filterKeyword], () => {
  // 由于使用计算属性，这里不需要额外操作
})

onMounted(fetchAchievements)

const dialogTitle = computed(() => {
  return isEditing.value ? '编辑成就' : '新建成就'
})
</script>

<style scoped>
.header-bar {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.slug-hint {
  margin-top: 6px;
  font-size: 12px;
}

.slug-preview {
  color: #909399;
  background: #f2f3f5;
  padding: 4px 8px;
  border-radius: 4px;
}

.limit-hint {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 8px;
}
</style>