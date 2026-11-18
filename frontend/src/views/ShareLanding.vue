<template>
  <div class="share-landing" v-loading="loading">
    <template v-if="valid && photoset">
      <div class="share-header">
        <el-icon :size="22" color="#409eff"><Picture /></el-icon>
        <span>分享的套图</span>
      </div>
      <el-card class="share-card">
        <el-image
          v-if="photoset.cover || photoset.cover_image"
          :src="photoset.cover || photoset.cover_image"
          fit="cover"
          class="share-cover"
        />
        <div class="share-info">
          <h2>{{ photoset.title || '未命名套图' }}</h2>
          <div class="share-meta">
            <span>{{ photoset.photo_count || 0 }} 张图片</span>
            <span>{{ photoset.view_count || 0 }} 浏览</span>
            <el-tag v-if="photoset.is_vip || photoset.is_free === 0" size="small" type="warning">VIP</el-tag>
          </div>
          <p class="share-desc">{{ photoset.description || '暂无描述' }}</p>
          <div class="share-actions">
            <el-button type="primary" :icon="View" @click="goDetail">查看完整套图</el-button>
            <el-button :icon="Star" @click="goDetail">收藏</el-button>
          </div>
        </div>
      </el-card>
      <div class="share-footer">
        <p>来自 PhotoSet · 发现精彩套图</p>
        <el-button link type="primary" @click="$router.push('/')">访问首页</el-button>
      </div>
    </template>

    <el-result
      v-else-if="!loading"
      icon="warning"
      :title="errorTitle"
      :sub-title="errorSubtitle"
    >
      <template #extra>
        <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Picture, View, Star } from '@element-plus/icons-vue'
import { verifyShareLink, getPhotosetDetail } from '@/api'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const valid = ref(false)
const photoset = ref(null)
const errorTitle = ref('链接无效')
const errorSubtitle = ref('该分享链接可能已过期或被修改')

onMounted(() => {
  verifyAndLoad()
})

async function verifyAndLoad() {
  loading.value = true
  try {
    const { id } = route.params
    const { s, t, sig } = route.query
    if (!s || !t || !sig) {
      errorTitle.value = '链接无效'
      errorSubtitle.value = '分享链接缺少必要参数'
      return
    }
    // 验签
    const res = await verifyShareLink(id, { s, t, sig })
    const data = res?.data || {}
    if (!data.valid) {
      errorTitle.value = '链接已失效'
      errorSubtitle.value = '该分享链接签名无效，请联系分享者重新获取'
      return
    }
    valid.value = true
    // 验签通过后拉取套图详情
    const detailRes = await getPhotosetDetail(data.photoset_id || id)
    photoset.value = detailRes?.data || null
  } catch (e) {
    errorTitle.value = '链接验证失败'
    errorSubtitle.value = e?.message || '请检查链接是否完整'
  } finally {
    loading.value = false
  }
}

function goDetail() {
  router.push(`/photoset/${photoset.value?.id || route.params.id}`)
}
</script>

<style scoped>
.share-landing { max-width: 640px; margin: 0 auto; padding: 40px 20px 60px; }
.share-header { display: flex; align-items: center; gap: 8px; margin-bottom: 16px; color: #606266; font-weight: 600; }
.share-card { overflow: hidden; }
.share-cover { width: 100%; height: 260px; }
.share-info { padding: 20px; }
.share-info h2 { margin: 0 0 10px; font-size: 22px; }
.share-meta { display: flex; align-items: center; gap: 16px; color: #999; font-size: 13px; }
.share-desc { color: #666; line-height: 1.7; margin: 12px 0; }
.share-actions { display: flex; gap: 12px; margin-top: 16px; }
.share-footer { text-align: center; margin-top: 24px; color: #bbb; font-size: 13px; }
.share-footer p { margin: 0 0 6px; }
</style>
