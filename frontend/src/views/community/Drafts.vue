<template>
  <MyContentLayout title="草稿箱" :fetcher="fetcher" empty-text="还没有草稿">
    <template #default="{ items }">
      <div v-for="draft in items" :key="draft.id" class="draft-card">
        <div class="draft-title" @click="editDraft(draft)">{{ draft.title || '无标题草稿' }}</div>
        <div class="draft-content">{{ draft.content?.slice(0, 120) }}</div>
        <div class="draft-meta">
          <span>{{ formatTime(draft.updated_at || draft.created_at) }}</span>
          <el-button link size="small" type="primary" @click="editDraft(draft)">继续编辑</el-button>
          <el-button link size="small" type="danger" :icon="Delete" @click="removeDraft(draft)">删除</el-button>
        </div>
      </div>
    </template>
  </MyContentLayout>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Delete } from '@element-plus/icons-vue'
import MyContentLayout from './MyContentLayout.vue'
import { getMyDrafts, deleteDraft } from '@/api'
import { formatTime } from '@/utils/post'

const router = useRouter()
const refreshKey = ref(0)

const fetcher = (params) => getMyDrafts(params).then(res => {
  const data = res?.data || {}
  return { items: data.drafts || [], total: data.pagination?.total || 0 }
})

function editDraft(draft) {
  router.push({ path: '/community/create', query: { draft: draft.id } })
}

async function removeDraft(draft) {
  try {
    await ElMessageBox.confirm('确定删除这条草稿？', '提示', { type: 'warning' })
    await deleteDraft(draft.id)
    ElMessage.success('已删除')
    refreshKey.value++
  } catch (e) { /* 取消 */ }
}
</script>

<style scoped>
.draft-card { background: #fff; border-radius: 10px; padding: 16px 20px; margin-bottom: 12px; }
.draft-title { font-size: 16px; font-weight: 600; cursor: pointer; margin-bottom: 6px; }
.draft-title:hover { color: #409eff; }
.draft-content { color: #666; font-size: 13px; line-height: 1.6; }
.draft-meta { display: flex; align-items: center; gap: 16px; margin-top: 10px; color: #999; font-size: 12px; }
.draft-meta .el-button { margin-left: auto; }
.draft-meta .el-button + .el-button { margin-left: 0; }
</style>
