<template>
  <router-link :to="`/photoset/${data.id}`" class="photoset-card">
    <!-- 封面图 -->
    <div class="card-cover">
      <el-image
        :src="data.cover"
        :alt="data.title"
        fit="cover"
        loading="lazy"
        :preview-src-list="[data.cover]"
      >
        <template #error>
          <div class="image-error">
            <el-icon :size="40"><Picture /></el-icon>
          </div>
        </template>
      </el-image>
      <!-- 免费/付费标识 -->
      <span :class="['badge', data.is_free ? 'badge-free' : 'badge-paid']">
        {{ data.is_free ? '免费' : '付费' }}
      </span>
      <!-- 价格（付费套图显示） -->
      <span v-if="!data.is_free && data.price" class="price">
        ¥{{ data.price.toFixed(2) }}
      </span>
    </div>

    <!-- 卡片内容 -->
    <div class="card-content">
      <h3 class="card-title">{{ data.title }}</h3>
      <p class="card-desc">{{ data.description || '暂无描述' }}</p>

      <!-- 标签 -->
      <div class="card-tags" v-if="data.tags?.length">
        <el-tag
          v-for="tag in data.tags.slice(0, 3)"
          :key="tag.id"
          size="small"
          type="info"
        >
          {{ tag.name }}
        </el-tag>
      </div>

      <!-- 作者信息 -->
      <div class="card-author">
        <el-avatar :size="24">
          {{ data.user?.nickname?.charAt(0) || 'U' }}
        </el-avatar>
        <span class="author-name">{{ data.user?.nickname || '匿名用户' }}</span>
        <span class="photo-count">
          <el-icon><PictureFilled /></el-icon>
          {{ data.photo_count || '' }}
        </span>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import { Picture, PictureFilled } from '@element-plus/icons-vue'

defineProps({
  data: {
    type: Object,
    required: true
  }
})
</script>

<style scoped>
.photoset-card {
  display: block;
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  text-decoration: none;
  color: inherit;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.card-cover {
  position: relative;
  width: 100%;
  padding-top: 66.67%; /* 3:2 比例 */
  background: #f0f0f0;
  overflow: hidden;
}

.card-cover :deep(.el-image) {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.image-error {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  color: #ccc;
}

.badge {
  position: absolute;
  top: 12px;
  left: 12px;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  z-index: 1;
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
  position: absolute;
  top: 12px;
  right: 12px;
  background: rgba(0, 0, 0, 0.7);
  color: #fff;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 600;
}

.card-content {
  padding: 16px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-desc {
  font-size: 13px;
  color: #666;
  margin-bottom: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.card-tags :deep(.el-tag) {
  margin: 0;
}

.card-author {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.author-name {
  font-size: 13px;
  color: #666;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.photo-count {
  font-size: 12px;
  color: #999;
  display: flex;
  align-items: center;
  gap: 2px;
}

@media (max-width: 768px) {
  .card-content {
    padding: 12px;
  }

  .card-title {
    font-size: 14px;
  }

  .card-desc {
    font-size: 12px;
  }
}
</style>
