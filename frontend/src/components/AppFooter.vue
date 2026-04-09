<template>
  <footer class="app-footer">
    <div class="footer-container">
      <div class="footer-info">
        <div class="footer-brand">
          <el-icon :size="20"><Camera /></el-icon>
          <span>PhotoSet</span>
        </div>
        <p class="footer-desc">专业摄影套图浏览平台，发现美好视觉作品</p>
      </div>
      <div class="footer-links">
        <div class="link-group">
          <h4>关于</h4>
          <a href="#">关于我们</a>
          <a href="#">使用协议</a>
          <a href="#">隐私政策</a>
        </div>
        <div class="link-group">
          <h4>帮助</h4>
          <a href="#">常见问题</a>
          <a href="#">联系我们</a>
        </div>
      </div>
    </div>
    <div class="footer-bottom">
      <p>© 2026 PhotoSet. All rights reserved.</p>
    </div>
  </footer>
</template>

<script setup>
import { Camera } from '@element-plus/icons-vue'
</script>

<style scoped>
.app-footer {
  background: #2c3e50;
  color: #ecf0f1;
  margin-top: auto;
}

.footer-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 48px 20px 32px;
  display: flex;
  justify-content: space-between;
  gap: 48px;
}

.footer-brand {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 12px;
}

.footer-desc {
  color: #95a5a6;
  font-size: 14px;
}

.footer-links {
  display: flex;
  gap: 64px;
}

.link-group h4 {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #fff;
}

.link-group a {
  display: block;
  color: #95a5a6;
  text-decoration: none;
  font-size: 14px;
  margin-bottom: 8px;
  transition: color 0.2s;
}

.link-group a:hover {
  color: #409eff;
}

.footer-bottom {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding: 20px;
  text-align: center;
  color: #7f8c8d;
  font-size: 13px;
}

@media (max-width: 768px) {
  .footer-container {
    flex-direction: column;
    gap: 32px;
    padding: 32px 20px 24px;
  }

  .footer-links {
    gap: 32px;
  }
}
</style>
