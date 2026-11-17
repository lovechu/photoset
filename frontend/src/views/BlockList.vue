<template>
  <div class="block-list-page">
    <el-card>
      <template #header><h3>屏蔽列表</h3></template>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="全部" name="all" />
        <el-tab-pane label="拉黑" name="block" />
        <el-tab-pane label="静音" name="mute" />
      </el-tabs>
      <div v-loading="loading">
        <el-empty v-if="!loading && filtered.length === 0" description="没有屏蔽的用户" :image-size="60" />
        <div v-for="item in filtered" :key="item.id" class="block-item">
          <el-avatar :size="42" :src="item.avatar">{{ item.username?.charAt(0) || 'U' }}</el-avatar>
          <div class="block-info">
            <div class="block-name">{{ item.username || `用户 ${item.blocked_id}` }}</div>
            <div class="block-meta">
              <el-tag size="small" :type="item.block_type === 'block' ? 'danger' : 'warning'">
                {{ item.block_type === 'block' ? '已拉黑' : '已静音' }}
              </el-tag>
              <span>{{ formatTime(item.created_at) }}</span>
            </div>
          </div>
          <el-button size="small" @click="unblock(item)">解除</el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getBlockedUsers, unblockUser } from '@/api'
import { formatTime } from '@/utils/post'

const list = ref([])
const loading = ref(false)
const activeTab = ref('all')

const filtered = computed(() => {
  if (activeTab.value === 'all') return list.value
  return list.value.filter(item => item.block_type === activeTab.value)
})

onMounted(() => load())
async function load() {
  loading.value = true
  try {
    const res = await getBlockedUsers({ page: 1, page_size: 100 })
    const data = res?.data || {}
    list.value = data.list || []
  } catch (e) { console.warn('加载屏蔽列表失败:', e.message) }
  finally { loading.value = false }
}

async function unblock(item) {
  try {
    await ElMessageBox.confirm(`确定解除对「${item.username || '该用户'}」的屏蔽？`, '提示', { type: 'warning' })
    await unblockUser(item.blocked_id)
    list.value = list.value.filter(i => i.id !== item.id)
    ElMessage.success('已解除')
  } catch (e) { /* 取消 */ }
}
</script>

<style scoped>
.block-list-page { max-width: 760px; margin: 0 auto; padding: 24px 0 60px; }
.el-card__header h3 { margin: 0; font-size: 18px; }
.block-item { display: flex; align-items: center; gap: 14px; padding: 14px 0; border-bottom: 1px solid #f0f0f0; }
.block-info { flex: 1; }
.block-name { font-weight: 600; }
.block-meta { display: flex; align-items: center; gap: 10px; margin-top: 4px; font-size: 12px; color: #999; }
</style>
