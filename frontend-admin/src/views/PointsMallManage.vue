<template>
  <div class="points-mall-manage">
    <div class="header-bar">
      <el-button type="primary" @click="handleCreate" :icon="Plus">
        新增商品
      </el-button>
      <el-input
        v-model="filterKeyword"
        placeholder="搜索商品名称"
        clearable
        @clear="fetchItems"
        @keyup.enter="fetchItems"
        style="width: 240px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <el-table :data="filteredItems" v-loading="loading" stripe style="width: 100%" border>
      <el-table-column prop="id" label="ID" width="70" align="center" />
      <el-table-column prop="name" label="商品名称" min-width="150" />
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column label="分类" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getCategoryTagType(row.category)" size="small">
            {{ getCategoryName(row.category) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getItemTypeTagType(row.item_type)" size="small">
            {{ getItemTypeName(row.item_type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="points_cost" label="积分价格" width="100" align="center">
        <template #default="{ row }">
          <span style="color: #E6A23C; font-weight: bold;">{{ row.points_cost }}</span>
        </template>
      </el-table-column>
      <el-table-column label="库存" width="120" align="center">
        <template #default="{ row }">
          <div v-if="row.is_unlimited">
            <el-tag type="success" size="small">无限</el-tag>
          </div>
          <div v-else>
            <span>{{ row.total_stock - row.used_stock }}</span>
            <span class="stock-hint">/ {{ row.total_stock }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="最低等级" width="100" align="center">
        <template #default="{ row }">
          <span>Lv.{{ row.min_level }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="sort_order" label="排序" width="80" align="center" />
      <el-table-column label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'" size="small">
            {{ row.is_active ? '上架' : '下架' }}
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
              title="确定要删除这个商品吗？删除后用户将无法兑换此商品。"
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

    <!-- 新建/编辑商品对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="resetForm"
    >
      <el-form ref="formRef" :model="editForm" :rules="formRules" label-width="100px">
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="editForm.name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="editForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入商品描述"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="商品分类" prop="category">
          <el-select v-model="editForm.category" placeholder="请选择商品分类">
            <el-option label="徽章" value="badge" />
            <el-option label="称号" value="title" />
            <el-option label="特权" value="privilege" />
            <el-option label="虚拟礼物" value="virtual_gift" />
          </el-select>
        </el-form-item>
        <el-form-item label="商品类型" prop="item_type">
          <el-select v-model="editForm.item_type" placeholder="请选择商品类型">
            <el-option label="徽章" value="badge" />
            <el-option label="称号" value="title" />
            <el-option label="VIP天数" value="vip_days" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="商品值" prop="item_value">
          <el-input v-model="editForm.item_value" placeholder="请输入商品值（如徽章ID、称号名称、天数等）" />
        </el-form-item>
        <el-form-item label="积分价格" prop="points_cost">
          <el-input-number v-model="editForm.points_cost" :min="0" />
        </el-form-item>
        
        <el-divider content-position="left">库存设置</el-divider>
        
        <el-form-item label="无限库存" prop="is_unlimited">
          <el-switch v-model="editForm.is_unlimited" />
          <span class="limit-hint">开启后库存无限制</span>
        </el-form-item>
        <el-form-item v-if="!editForm.is_unlimited" label="总库存" prop="total_stock">
          <el-input-number v-model="editForm.total_stock" :min="0" />
        </el-form-item>
        
        <el-divider content-position="left">其他设置</el-divider>
        
        <el-form-item label="最低等级" prop="min_level">
          <el-input-number v-model="editForm.min_level" :min="1" :max="10" />
          <span class="limit-hint">用户达到此等级才可兑换</span>
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="editForm.sort_order" :min="0" />
        </el-form-item>
        <el-form-item label="上架" prop="is_active">
          <el-switch v-model="editForm.is_active" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ isEditing ? '保存修改' : '创建商品' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, reactive, computed } from 'vue'
import { getPointsMallItems, createPointsMallItem, updatePointsMallItem, deletePointsMallItem } from '@/api/level'
import { ElMessage } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'

const loading = ref(false)
const itemList = ref([])
const filterKeyword = ref('')
const submitting = ref(false)
const dialogVisible = ref(false)
const isEditing = ref(false)

const editForm = reactive({
  id: null,
  name: '',
  description: '',
  category: 'badge',
  item_type: 'badge',
  item_value: '',
  points_cost: 0,
  is_unlimited: true,
  total_stock: 0,
  min_level: 1,
  sort_order: 0,
  is_active: true
})

const formRules = {
  name: [
    { required: true, message: '请输入商品名称', trigger: 'blur' },
    { min: 2, max: 100, message: '名称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ],
  item_type: [
    { required: true, message: '请选择商品类型', trigger: 'change' }
  ],
  item_value: [
    { required: true, message: '请输入商品值', trigger: 'blur' }
  ],
  points_cost: [
    { required: true, message: '请输入积分价格', trigger: 'blur' }
  ],
  total_stock: [
    { required: true, message: '请输入总库存', trigger: 'blur' }
  ]
}

const formRef = ref()

const filteredItems = computed(() => {
  if (!filterKeyword.value) return itemList.value
  const keyword = filterKeyword.value.toLowerCase()
  return itemList.value.filter(item => 
    item.name.toLowerCase().includes(keyword) ||
    (item.description && item.description.toLowerCase().includes(keyword))
  )
})

function getCategoryName(category) {
  const categoryMap = {
    'badge': '徽章',
    'title': '称号',
    'privilege': '特权',
    'virtual_gift': '虚拟礼物'
  }
  return categoryMap[category] || category
}

function getCategoryTagType(category) {
  const categoryMap = {
    'badge': 'primary',
    'title': 'success',
    'privilege': 'warning',
    'virtual_gift': 'info'
  }
  return categoryMap[category] || ''
}

function getItemTypeName(itemType) {
  const typeMap = {
    'badge': '徽章',
    'title': '称号',
    'vip_days': 'VIP天数',
    'custom': '自定义'
  }
  return typeMap[itemType] || itemType
}

function getItemTypeTagType(itemType) {
  const typeMap = {
    'badge': 'primary',
    'title': 'success',
    'vip_days': 'warning',
    'custom': 'info'
  }
  return typeMap[itemType] || ''
}

async function fetchItems() {
  loading.value = true
  try {
    const res = await getPointsMallItems()
    itemList.value = res.data || []
  } catch (error) {
    console.error('获取商品列表失败:', error)
    ElMessage.error('获取商品列表失败')
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  isEditing.value = false
  resetForm()
  dialogVisible.value = true
}

function handleEdit(item) {
  isEditing.value = true
  Object.assign(editForm, {
    id: item.id,
    name: item.name,
    description: item.description || '',
    category: item.category || 'badge',
    item_type: item.item_type || 'badge',
    item_value: item.item_value || '',
    points_cost: item.points_cost || 0,
    is_unlimited: item.is_unlimited || false,
    total_stock: item.total_stock || 0,
    min_level: item.min_level || 1,
    sort_order: item.sort_order || 0,
    is_active: item.is_active || false
  })
  dialogVisible.value = true
}

async function handleDelete(item) {
  try {
    await deletePointsMallItem(item.id)
    ElMessage.success('删除成功')
    fetchItems()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

function resetForm() {
  Object.assign(editForm, {
    id: null,
    name: '',
    description: '',
    category: 'badge',
    item_type: 'badge',
    item_value: '',
    points_cost: 0,
    is_unlimited: true,
    total_stock: 0,
    min_level: 1,
    sort_order: 0,
    is_active: true
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
      description: editForm.description.trim() || null,
      category: editForm.category,
      item_type: editForm.item_type,
      item_value: editForm.item_value.trim(),
      points_cost: editForm.points_cost,
      is_unlimited: editForm.is_unlimited,
      total_stock: editForm.is_unlimited ? -1 : editForm.total_stock,
      min_level: editForm.min_level,
      sort_order: editForm.sort_order,
      is_active: editForm.is_active
    }

    if (isEditing.value) {
      await updatePointsMallItem(editForm.id, formData)
      ElMessage.success('更新成功')
    } else {
      await createPointsMallItem(formData)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchItems()
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

onMounted(fetchItems)

const dialogTitle = computed(() => {
  return isEditing.value ? '编辑商品' : '新增商品'
})
</script>

<style scoped>
.header-bar {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stock-hint {
  font-size: 12px;
  color: #909399;
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