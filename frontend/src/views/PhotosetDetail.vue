<template>
  <div class="photoset-detail" v-loading="loading">
    <!-- 加载失败 -->
    <div v-if="!loading && error" class="error-state">
      <el-result
        icon="error"
        title="加载失败"
        :sub-title="error"
      >
        <template #extra>
          <el-button type="primary" @click="loadDetail">重试</el-button>
          <el-button @click="$router.back()">返回</el-button>
        </template>
      </el-result>
    </div>

    <!-- 详情内容 -->
    <template v-if="!loading && !error && detail">
      <!-- 返回按钮 -->
      <div class="back-bar">
        <el-button @click="$router.back()" :icon="ArrowLeft" text>
          返回列表
        </el-button>
      </div>

      <!-- 套图头部信息 -->
      <div class="detail-header">
        <div class="header-info">
          <h1 class="detail-title">{{ detail.title }}</h1>
          <div class="detail-meta">
            <span :class="['badge', detail.is_free ? 'badge-free' : 'badge-paid']">
              {{ detail.is_free ? '免费' : '付费' }}
            </span>
            <span v-if="!detail.is_free && detail.price" class="price">
              ¥{{ detail.price.toFixed(2) }}
            </span>
            <span class="author">
              <el-avatar :size="20">
                {{ detail.user?.nickname?.charAt(0) || 'U' }}
              </el-avatar>
              {{ detail.user?.nickname || '匿名用户' }}
            </span>
            <span class="date">
              {{ formatDate(detail.created_at) }}
            </span>
            <span class="favorite-btn" @click="handleToggleFavorite">
              <el-icon :size="20">
                <StarFilled v-if="isFavorited" style="color: #f5a623" />
                <Star v-else />
              </el-icon>
            </span>
          </div>
          <p class="detail-desc" v-if="detail.description">
            {{ detail.description }}
          </p>
          <!-- 标签 -->
          <div class="detail-tags" v-if="detail.tags?.length">
            <el-tag
              v-for="tag in detail.tags"
              :key="tag.id"
              size="small"
            >
              {{ tag.name }}
            </el-tag>
          </div>
        </div>
      </div>

      <!-- 图片画廊 -->
      <div class="photo-gallery">
        <!-- 免费套图或已付费：展示全部图片 -->
        <template v-if="detail.is_free || photos.length > 0">
          <div class="gallery-grid">
            <div
              v-for="(photo, index) in photos"
              :key="photo.id"
              class="gallery-item"
              @click="openViewer(index)"
            >
              <el-image
                :src="photo.url"
                :alt="`图片 ${index + 1}`"
                fit="cover"
                loading="lazy"
              >
                <template #error>
                  <div class="image-error">
                    <el-icon :size="32"><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
            </div>
          </div>
          <div class="gallery-tip">
            共 {{ photos.length }} 张图片，点击可放大查看
          </div>
        </template>

        <!-- 付费套图无权限：显示付费墙 -->
        <template v-else>
          <div class="paywall">
            <div class="paywall-cover">
              <el-image
                :src="detail.cover"
                fit="cover"
              />
              <div class="paywall-overlay">
                <el-icon :size="48"><Lock /></el-icon>
                <h3>此套图为付费内容</h3>
                <p>成为会员或单独购买即可解锁全部 {{ photoCount }} 张图片</p>
                <div class="paywall-actions">
                  <el-button type="primary" size="large">
                    开通会员 ¥{{ memberPrice }}/月
                  </el-button>
                  <el-button size="large">
                    单独购买 ¥{{ detail.price?.toFixed(2) }}
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>

      <!-- 图片预览器 -->
      <el-image-viewer
        v-if="showViewer"
        :url-list="viewerImages"
        :initial-index="viewerIndex"
        @close="showViewer = false"
      />
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getPhotosetDetail, addFavorite, removeFavorite } from '@/api'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Picture, Lock, Star, StarFilled } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const error = ref('')
const detail = ref(null)
const showViewer = ref(false)
const viewerIndex = ref(0)
const isFavorited = ref(false)

// 图片列表（可能为空）
const photos = computed(() => detail.value?.photos || [])

// 可预览的图片 URL 列表
const viewerImages = computed(() => photos.value.map(p => p.url))

// 图片数量（用于付费墙提示）
const photoCount = computed(() => {
  // 如果后端返回了照片数量就用返回的，否则从 photos 数组长度推断
  return detail.value?.photo_count || photos.value.length || '多张'
})

// 会员价格（模拟）
const memberPrice = ref(29.9)

// 加载详情
const loadDetail = async () => {
  const id = route.params.id
  if (!id) {
    error.value = '缺少套图 ID'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const res = await getPhotosetDetail(Number(id))
    detail.value = res.data
    isFavorited.value = res.data.is_favorited || false

    // 根据返回数据判断是否可以看到完整图片
    if (!detail.value.is_free && (!photos.value || photos.value.length === 0)) {
      // 说明是付费套图且当前用户无权查看
      if (!userStore.isLoggedIn) {
        ElMessage.info('登录后可查看付费套图全部内容')
      } else if (userStore.user?.role === 'user') {
        ElMessage.info('此套图为付费内容，开通会员或购买后可查看')
      }
    }
  } catch (e) {
    console.error('加载套图详情失败', e)
    error.value = '加载套图详情失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

// 打开图片预览器
const openViewer = (index) => {
  viewerIndex.value = index
  showViewer.value = true
}

// 切换收藏
const handleToggleFavorite = () => {
  if (!userStore.isLoggedIn) {
    ElMessage.info('请先登录')
    router.push({ name: 'Login', query: { redirect: route.fullPath } })
    return
  }

  if (isFavorited.value) {
    removeFavorite(detail.value.id).then(() => {
      isFavorited.value = false
      ElMessage.success('已取消收藏')
    }).catch(() => {
      ElMessage.error('操作失败')
    })
  } else {
    addFavorite(detail.value.id).then(() => {
      isFavorited.value = true
      ElMessage.success('收藏成功')
    }).catch(() => {
      ElMessage.error('操作失败')
    })
  }
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.photoset-detail {
  padding: 0 0 60px;
}

.error-state {
  padding: 80px 0;
}

/* 返回按钮 */
.back-bar {
  padding: 16px 0;
}

/* 头部信息 */
.detail-header {
  background: #fff;
  padding: 24px;
  border-radius: 12px;
  margin-bottom: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.detail-title {
  font-size: 28px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 12px;
}

.detail-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 16px;
}

.badge {
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.badge-free {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
  color: #fff;
}

.badge-paid {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.price {
  font-size: 18px;
  font-weight: 600;
  color: #f56c6c;
}

.author {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #666;
}

.date {
  font-size: 14px;
  color: #999;
}

.favorite-btn {
  cursor: pointer;
  display: flex;
  align-items: center;
  transition: transform 0.2s;
}

.favorite-btn:hover {
  transform: scale(1.2);
}

.detail-desc {
  font-size: 15px;
  color: #666;
  line-height: 1.8;
  margin-bottom: 16px;
}

.detail-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

/* 画廊 */
.gallery-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
}

.gallery-item {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  aspect-ratio: 3/2;
}

.gallery-item :deep(.el-image) {
  width: 100%;
  height: 100%;
}

.gallery-item:hover {
  transform: scale(1.02);
  transition: transform 0.3s;
}

.image-error {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  color: #ccc;
}

.gallery-tip {
  text-align: center;
  padding: 16px;
  color: #999;
  font-size: 14px;
}

/* 付费墙 */
.paywall {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.paywall-cover {
  position: relative;
  max-height: 400px;
  overflow: hidden;
}

.paywall-cover :deep(.el-image) {
  width: 100%;
  filter: blur(20px);
  transform: scale(1.1);
}

.paywall-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #fff;
  padding: 40px;
  text-align: center;
}

.paywall-overlay h3 {
  font-size: 24px;
  margin: 16px 0 8px;
}

.paywall-overlay p {
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 24px;
}

.paywall-actions {
  display: flex;
  gap: 16px;
}

@media (max-width: 768px) {
  .detail-header {
    padding: 16px;
  }

  .detail-title {
    font-size: 22px;
  }

  .gallery-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
  }

  .paywall-actions {
    flex-direction: column;
    width: 100%;
  }

  .paywall-actions .el-button {
    width: 100%;
  }
}
</style>
