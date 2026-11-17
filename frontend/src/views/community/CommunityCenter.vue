<template>
  <div class="community-center">
    <div class="page-header"><h2>社区中心</h2></div>
    <el-row :gutter="16">
      <el-col :xs="24" :sm="12" :md="8" v-for="item in entries" :key="item.path">
        <el-card class="entry-card" shadow="hover" @click="$router.push(item.path)">
          <div class="entry-icon"><el-icon :size="28"><component :is="item.icon" /></el-icon></div>
          <div class="entry-text">
            <h3>{{ item.title }}</h3>
            <p>{{ item.desc }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import {
  Document, ChatDotRound, Star, Folder, EditPen, Collection,
  Message, StarFilled, Timer, Warning, Tickets, User
} from '@element-plus/icons-vue'

const entries = [
  { title: '我的帖子', desc: '查看我发布的帖子', path: '/community/my-posts', icon: Document },
  { title: '我的回复', desc: '查看我的回复记录', path: '/community/my-replies', icon: ChatDotRound },
  { title: '我的点赞', desc: '我点赞过的帖子', path: '/community/my-likes', icon: StarFilled },
  { title: '社区收藏', desc: '收藏的帖子', path: '/community/favorites', icon: Star },
  { title: '草稿箱', desc: '未发布的草稿', path: '/community/drafts', icon: Tickets },
  { title: '私信消息', desc: '查看私信会话', path: '/community/conversations', icon: Message },
  { title: '我的等级', desc: '等级与成就', path: '/user-level', icon: User },
  { title: '积分商城', desc: '积分兑换好礼', path: '/points-mall', icon: Folder },
  { title: '浏览历史', desc: '看过的套图', path: '/history', icon: Timer },
  { title: '社区公约', desc: '社区行为规范', path: '/community/guidelines', icon: Warning },
  { title: '意见反馈', desc: '提交反馈建议', path: '/feedback', icon: EditPen },
  { title: '屏蔽列表', desc: '管理屏蔽用户', path: '/block-list', icon: Collection }
]
</script>

<style scoped>
.community-center { max-width: 1000px; margin: 0 auto; padding: 24px 0 60px; }
.page-header { margin-bottom: 20px; }
.page-header h2 { margin: 0; }
.entry-card { cursor: pointer; margin-bottom: 16px; display: flex; align-items: center; gap: 14px; }
.entry-card:hover { border-color: #409eff; }
.entry-icon { color: #409eff; flex-shrink: 0; }
.entry-text h3 { margin: 0 0 4px; font-size: 15px; }
.entry-text p { margin: 0; color: #999; font-size: 12px; }
</style>
