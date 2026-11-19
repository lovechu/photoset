<template>
  <div class="collection-detail">
    <div class="back-bar">
      <el-button @click="$router.back()" :icon="ArrowLeft" text>返回合集列表</el-button>
    </div>

    <div v-loading="loading">
      <template v-if="collection">
        <div class="collection-header">
          <div>
            <h2>{{ collection.name }}</h2>
            <p class="desc">{{ collection.description || '暂无描述' }}</p>
            <div class="meta">
              <el-tag size="small" :type="collection.is_public ? 'success' : 'info'">
                {{ collection.is_public ? '公开' : '私密' }}
              </el-tag>
              <span>共 {{ items.length }} 个套图</span>
            </div>
          </div>
          <el-button type="primary" :icon="Plus" @click="openAddDialog">添加套图</el-button>
        </div>

        <el-empty v-if="items.length === 0" description="合集还是空的，去添加套图吧" :image-size="100" />
        <div v-else class="photoset-grid">
          <div v-for="item in items" :key="item.id" class="photoset-item">
            <el-image :src="item.cover || item.cover_url" fit="cover" class="item-cover" @click="goPhotoset(item)" />
            <div class="item-title" @click="goPhotoset(item)">{{ item.title || `套图 #${item.photoset_id}` }}</div>
            <el-button link size="small" type="danger" :icon="Delete" @click="removeItem(item)">移除</el-button>
          </div>
        </div>
      </template>
      <el-result v-else-if="!loading" icon="warning" title="合集不存在" sub-title="可能已被删除">
        <template #extra><el-button type="primary" @click="$router.push('/collections')">返回列表</el-button></template>
      </el-result>
    </div>

    <!-- 添加套图弹窗：展示我的收藏供选择 -->
    <el-dialog v-model="addDialogVisible" title="添加套图" width="640px">
      <div v-loading="favLoading">
        <el-empty v-if="!favLoading && favorites.length === 0" description="收藏夹为空，先去收藏一些套图" :image-size="60" />
        <el-checkbox-group v-else v-model="selectedPhotosets" class="fav-grid">
          <el-checkbox v-for="fav in favorites" :key="fav.photoset_id || fav.id" :value="fav.photoset_id || fav.id" border>
            <div class="fav-option">
              <el-image :src="fav.cover || fav.cover_image" fit="cover" style="width:40px;height:40px;border-radius:4px" />
              <span>{{ fav.title || `套图 #${fav.photoset_id || fav.id}` }}</span>
            </div>
          </el-checkbox>
        </el-checkbox-group>
      </div>
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="adding" @click="confirmAdd">添加选中</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Plus, Delete } from '@element-plus/icons-vue'
import {
  getCollection, removeCollectionItem, batchAddCollectionItems, getFavorites
} from '@/api'
import { useSeo } from '@/composables/useSeo'

const route = useRoute()
const router = useRouter()
const collectionId = route.params.id

useSeo(() => ({
  title: collection.value?.name ? `${collection.value.name} - 合集` : '合集详情',
  description: collection.value?.description || '查看 PhotoSet 用户创建的摄影合集',
  path: route.fullPath,
  noindex: true
}))

const loading = ref(false)
const collection = ref(null)
const items = ref([])

const addDialogVisible = ref(false)
const favLoading = ref(false)
const favorites = ref([])
const selectedPhotosets = ref([])
const adding = ref(false)

onMounted(() => loadCollection())

async function loadCollection() {
  loading.value = true
  try {
    const res = await getCollection(collectionId)
    collection.value = res?.data || null
    items.value = collection.value?.items || []
  } catch (e) {
    console.warn('加载合集失败:', e.message)
  } finally {
    loading.value = false
  }
}

async function openAddDialog() {
  addDialogVisible.value = true
  selectedPhotosets.value = []
  favLoading.value = true
  try {
    const res = await getFavorites({ page: 1, page_size: 100 })
    const data = res?.data || {}
    favorites.value = data.list || []
  } catch (e) {
    favorites.value = []
  } finally {
    favLoading.value = false
  }
}

async function confirmAdd() {
  if (!selectedPhotosets.value.length) {
    ElMessage.warning('请选择要添加的套图')
    return
  }
  adding.value = true
  try {
    await batchAddCollectionItems(collectionId, { photoset_ids: selectedPhotosets.value })
    addDialogVisible.value = false
    await loadCollection()
    ElMessage.success('已添加')
  } catch (e) {
    ElMessage.error('添加失败')
  } finally {
    adding.value = false
  }
}

async function removeItem(item) {
  const pid = item.photoset_id || item.id
  try {
    await ElMessageBox.confirm('确定从合集移除？', '提示', { type: 'warning' })
    await removeCollectionItem(collectionId, pid)
    items.value = items.value.filter(i => (i.photoset_id || i.id) !== pid)
    ElMessage.success('已移除')
  } catch (e) { /* 取消 */ }
}

function goPhotoset(item) {
  router.push(`/photoset/${item.photoset_id || item.id}`)
}
</script>

<style scoped>
.collection-detail { max-width: 1100px; margin: 0 auto; padding: 0 0 60px; }
.back-bar { padding: 16px 0; }
.collection-header { display: flex; align-items: flex-start; justify-content: space-between; background: #fff; padding: 24px; border-radius: 12px; margin-bottom: 24px; }
.collection-header h2 { margin: 0 0 8px; }
.desc { color: #666; margin: 0 0 12px; }
.meta { display: flex; align-items: center; gap: 12px; color: #999; font-size: 13px; }
.photoset-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 16px; }
.photoset-item { background: #fff; border-radius: 8px; overflow: hidden; padding-bottom: 8px; }
.item-cover { width: 100%; height: 140px; cursor: pointer; }
.item-title { padding: 10px 12px 4px; font-size: 14px; cursor: pointer; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.photoset-item .el-button { margin-left: 12px; }
.fav-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 10px; max-height: 400px; overflow-y: auto; }
.fav-grid .el-checkbox { margin: 0; }
.fav-option { display: flex; align-items: center; gap: 8px; }
</style>
