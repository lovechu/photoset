<template>
  <div class="guidelines-page">
    <el-card>
      <template #header><h2>社区公约</h2></template>
      <div class="guidelines-content">
        <p class="intro">欢迎来到 PhotoSet 社区！为了营造友好、健康的交流环境，请共同遵守以下约定：</p>

        <h3>一、文明交流</h3>
        <ul>
          <li>尊重每一位成员，不进行人身攻击、谩骂或歧视</li>
          <li>理性讨论，对不同观点保持包容</li>
          <li>不使用侮辱性、攻击性或煽动性语言</li>
        </ul>

        <h3>二、内容规范</h3>
        <ul>
          <li>不发布违法违规、色情暴力、恐怖主义等内容</li>
          <li>不发布垃圾广告、营销推广、刷屏信息</li>
          <li>不发布虚假信息、谣言或侵犯他人隐私的内容</li>
          <li>转载内容请注明来源，尊重原创</li>
        </ul>

        <h3>三、版权尊重</h3>
        <ul>
          <li>发布作品需为自己原创或已获授权</li>
          <li>转载他人作品须获得许可并标注出处</li>
          <li>不得盗用他人图片、文字冒充自己的作品</li>
        </ul>

        <h3>四、账号安全</h3>
        <ul>
          <li>妥善保管账号密码，不共享账号</li>
          <li>发现账号异常请及时修改密码</li>
          <li>不利用技术手段干扰平台正常运行</li>
        </ul>

        <h3>五、违规处理</h3>
        <p>违反以上公约的内容可能被删除，违规账号可能被警告、禁言或封禁。情节严重的将移交相关机构处理。</p>
        <p>如发现违规行为，请使用帖子的「举报」功能，或前往<a @click="$router.push('/feedback')">意见反馈</a>提交。</p>

        <div class="footer-note">本公约可能不定期更新，请留意最新版本。感谢你的理解与配合，让我们一起维护良好的社区氛围。</div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
</script>

<style scoped>
.guidelines-page { max-width: 800px; margin: 0 auto; padding: 24px 0 60px; }
.el-card__header h2 { margin: 0; }
.guidelines-content { line-height: 1.9; color: #333; }
.guidelines-content .intro { color: #606266; background: #f4f4f5; padding: 14px 18px; border-radius: 8px; }
.guidelines-content h3 { margin: 24px 0 8px; font-size: 17px; color: #409eff; }
.guidelines-content ul { padding-left: 22px; }
.guidelines-content li { margin: 4px 0; }
.guidelines-content a { color: #409eff; cursor: pointer; }
.footer-note { margin-top: 28px; padding-top: 16px; border-top: 1px solid #eee; color: #999; font-size: 13px; }
</style>
