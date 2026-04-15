<template>
  <div class="site-settings">
    <el-card>
      <template #header>
        <span>站点设置</span>
      </template>

      <el-tabs v-model="activeTab">
        <!-- 基本信息 -->
        <el-tab-pane label="基本信息" name="general">
          <el-form :model="form" label-width="140px" style="max-width: 600px; margin-top: 16px;">
            <el-form-item label="网站名称">
              <el-input v-model="form.site_title" placeholder="例：PhotoSet" />
            </el-form-item>
            <el-form-item label="网站描述">
              <el-input v-model="form.site_description" type="textarea" :rows="3" placeholder="网站简介" />
            </el-form-item>
            <el-form-item label="ICP 备案号">
              <el-input v-model="form.site_icp" placeholder="例：京ICP备XXXXXXXX号" />
            </el-form-item>
            <el-form-item label="开放注册">
              <el-switch v-model="form.register_enabled" active-value="true" inactive-value="false" />
            </el-form-item>
            <el-form-item label="要求邮箱验证">
              <el-switch v-model="form.email_verify_required" active-value="true" inactive-value="false" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="save('general')">保存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- SEO 设置 -->
        <el-tab-pane label="SEO 设置" name="seo">
          <el-form :model="form" label-width="140px" style="max-width: 600px; margin-top: 16px;">
            <el-form-item label="关键词">
              <el-input v-model="form.site_keywords" type="textarea" :rows="3"
                placeholder="多个关键词用英文逗号分隔，例：摄影,套图,写真" />
              <div class="form-tip">用于 meta keywords，有助于搜索引擎收录</div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="save('seo')">保存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 关于我 -->
        <el-tab-pane label="关于我" name="about">
          <el-form :model="form" label-width="140px" style="max-width: 700px; margin-top: 16px;">
            <el-form-item label="页面内容">
              <el-input
                v-model="form.about_content"
                type="textarea"
                :rows="12"
                placeholder="支持 Markdown 格式，例：# 关于我&#10;我是一名摄影爱好者..."
              />
              <div class="form-tip">支持 Markdown 语法</div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="save('about')">保存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 邮件设置 -->
        <el-tab-pane label="邮件设置" name="mail">
          <el-form :model="form" label-width="140px" style="max-width: 600px; margin-top: 16px;">
            <el-form-item label="SMTP 服务器">
              <el-input v-model="form.smtp_host" placeholder="例：smtp.qq.com" />
            </el-form-item>
            <el-form-item label="SMTP 端口">
              <el-input v-model="form.smtp_port" placeholder="例：465 或 587" />
            </el-form-item>
            <el-form-item label="SMTP 用户名">
              <el-input v-model="form.smtp_user" placeholder="发件邮箱地址" />
            </el-form-item>
            <el-form-item label="SMTP 密码">
              <el-input v-model="form.smtp_password" type="password" show-password placeholder="邮箱授权码" />
            </el-form-item>
            <el-form-item label="发件人名称">
              <el-input v-model="form.smtp_from_name" placeholder="例：PhotoSet 平台" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="save('mail')">保存</el-button>
              <el-button @click="testMail" :loading="testingMail" style="margin-left: 12px;">发送测试邮件</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 水印设置 -->
        <el-tab-pane label="水印设置" name="watermark">
          <el-form :model="form" label-width="140px" style="max-width: 600px; margin-top: 16px;">
            <el-form-item label="启用水印">
              <el-switch v-model="form.watermark_enabled" active-value="true" inactive-value="false" />
            </el-form-item>
            <el-form-item label="水印文字">
              <el-input v-model="form.watermark_text" placeholder="例：© PhotoSet" />
            </el-form-item>
            <el-form-item label="水印透明度">
              <el-slider v-model="watermarkOpacity" :min="0" :max="100" show-input style="width: 400px;" />
              <div class="form-tip">0 = 全透明，100 = 不透明，建议 20-40</div>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="save('watermark')">保存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 存储与 CDN -->
        <el-tab-pane label="存储与 CDN" name="storage">
          <div style="max-width: 680px; margin-top: 16px;">
            <!-- 当前状态 -->
            <div v-if="storageStatus" class="storage-status-card">
              <div class="status-header">
                <h4>当前存储</h4>
                <el-tag :type="storageStatus.type === 'local' ? 'info' : 'success'" size="small">
                  {{ storageStatus.label }}
                </el-tag>
              </div>
              <div class="status-body">
                <template v-if="storageStatus.type === 'local'">
                  <p>📁 路径：{{ storageStatus.path }}</p>
                </template>
                <template v-else>
                  <p>☁️ 提供商：{{ storageStatus.provider || 'S3 兼容' }}</p>
                  <p v-if="storageStatus.cdn_domain">🌐 CDN 域名：{{ storageStatus.cdn_domain }}</p>
                  <p>🔑 AK 已配置：{{ storageStatus.s3_access_key_set ? '是' : '否' }}</p>
                  <p>🔑 SK 已配置：{{ storageStatus.s3_secret_key_set ? '是' : '否' }}</p>
                  <p>📦 Bucket 已配置：{{ storageStatus.s3_bucket_set ? '是' : '否' }}</p>
                </template>
              </div>
            </div>

            <el-divider />

            <el-alert
              title="修改存储配置后需重启后端服务才能生效"
              type="warning"
              :closable="false"
              show-icon
              style="margin-bottom: 20px;"
            />

            <el-form :model="form" label-width="140px">
              <el-form-item label="存储类型">
                <el-radio-group v-model="form.storage_type">
                  <el-radio value="local">本地存储</el-radio>
                  <el-radio value="s3">S3 兼容存储</el-radio>
                  <el-radio value="r2">Cloudflare R2</el-radio>
                </el-radio-group>
              </el-form-item>

              <!-- S3 通用配置 -->
              <template v-if="form.storage_type === 's3'">
                <el-form-item label="S3 Endpoint">
                  <el-input v-model="form.s3_endpoint" placeholder="例：https://oss-cn-hangzhou.aliyuncs.com" />
                  <div class="form-tip">阿里 OSS、AWS S3、MinIO 等的 Endpoint</div>
                </el-form-item>
                <el-form-item label="Region">
                  <el-input v-model="form.s3_region" placeholder="例：cn-hangzhou 或 us-east-1" />
                </el-form-item>
                <el-form-item label="Access Key ID">
                  <el-input v-model="form.s3_access_key" placeholder="Access Key ID" />
                </el-form-item>
                <el-form-item label="Access Key Secret">
                  <el-input v-model="form.s3_secret_key" type="password" show-password placeholder="Access Key Secret" />
                </el-form-item>
                <el-form-item label="Bucket">
                  <el-input v-model="form.s3_bucket" placeholder="存储桶名称" />
                </el-form-item>
              </template>

              <!-- R2 配置 -->
              <template v-if="form.storage_type === 'r2'">
                <el-form-item label="Account ID">
                  <el-input v-model="form.r2_account_id" placeholder="Cloudflare Account ID" />
                </el-form-item>
                <el-form-item label="Access Key ID">
                  <el-input v-model="form.s3_access_key" placeholder="R2 API Token Access Key ID" />
                </el-form-item>
                <el-form-item label="Access Key Secret">
                  <el-input v-model="form.s3_secret_key" type="password" show-password placeholder="R2 API Token Secret" />
                </el-form-item>
                <el-form-item label="Bucket">
                  <el-input v-model="form.s3_bucket" placeholder="R2 存储桶名称" />
                </el-form-item>
              </template>

              <!-- CDN 域名（S3/R2 通用） -->
              <el-form-item v-if="form.storage_type !== 'local'" label="CDN / 公开域名">
                <el-input v-model="form.cdn_domain" placeholder="例：https://cdn.example.com" />
                <div class="form-tip">用于生成图片访问 URL，建议配置自定义域名或 CDN 加速域名</div>
              </el-form-item>

              <el-form-item>
                <el-button type="primary" :loading="saving" @click="save('storage')">保存配置</el-button>
                <el-button
                  v-if="form.storage_type !== 'local'"
                  :loading="testingStorage"
                  @click="testStorage"
                  style="margin-left: 12px;"
                >
                  测试连接
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getSettings, updateSettings, getStorageStatus, testStorageConnection } from '@/api/index'

const activeTab = ref('general')
const saving = ref(false)
const testingMail = ref(false)
const testingStorage = ref(false)
const loaded = ref(false)
const storageStatus = ref(null)

const form = ref({
  site_title: '',
  site_description: '',
  site_icp: '',
  register_enabled: 'true',
  email_verify_required: 'false',
  site_keywords: '',
  about_content: '',
  smtp_host: '',
  smtp_port: '',
  smtp_user: '',
  smtp_password: '',
  smtp_from_name: '',
  watermark_enabled: 'false',
  watermark_text: '',
  watermark_opacity: '30',
  // 存储配置
  storage_type: 'local',
  s3_endpoint: '',
  s3_region: '',
  s3_access_key: '',
  s3_secret_key: '',
  s3_bucket: '',
  r2_account_id: '',
  cdn_domain: '',
})

// 水印透明度用数字绑定 slider
const watermarkOpacity = computed({
  get: () => parseInt(form.value.watermark_opacity) || 30,
  set: (v) => { form.value.watermark_opacity = String(v) }
})

onMounted(async () => {
  try {
    const data = await getSettings()
    if (data) {
      Object.keys(form.value).forEach(key => {
        if (data[key] !== undefined) {
          form.value[key] = data[key]
        }
      })
    }
    loaded.value = true
    // 加载存储状态
    loadStorageStatus()
  } catch (e) {
    ElMessage.error('加载配置失败')
  }
})

async function loadStorageStatus() {
  try {
    const res = await getStorageStatus()
    storageStatus.value = res.data || res
  } catch {
    // 静默处理
  }
}

async function save(group) {
  saving.value = true
  try {
    // 按 tab 保存对应的字段
    const groupKeys = {
      general: ['site_title', 'site_description', 'site_icp', 'register_enabled', 'email_verify_required'],
      seo: ['site_keywords'],
      about: ['about_content'],
      mail: ['smtp_host', 'smtp_port', 'smtp_user', 'smtp_password', 'smtp_from_name'],
      watermark: ['watermark_enabled', 'watermark_text', 'watermark_opacity'],
      storage: ['storage_type', 's3_endpoint', 's3_region', 's3_access_key', 's3_secret_key', 's3_bucket', 'r2_account_id', 'cdn_domain'],
    }
    const keys = groupKeys[group] || Object.keys(form.value)
    const payload = {}
    keys.forEach(k => { payload[k] = form.value[k] })

    await updateSettings(payload)
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function testMail() {
  testingMail.value = true
  try {
    // 先保存邮件配置
    await save('mail')
    ElMessage.info('测试邮件功能待后端实现后生效')
  } finally {
    testingMail.value = false
  }
}

async function testStorage() {
  testingStorage.value = true
  try {
    await testStorageConnection({
      storage_type: form.value.storage_type,
      s3_endpoint: form.value.s3_endpoint,
      s3_region: form.value.s3_region,
      s3_access_key: form.value.s3_access_key,
      s3_secret_key: form.value.s3_secret_key,
      s3_bucket: form.value.s3_bucket,
      r2_account_id: form.value.r2_account_id,
      cdn_domain: form.value.cdn_domain,
    })
    ElMessage.success('连接成功！存储配置可用')
  } catch {
    // 拦截器已处理
  } finally {
    testingStorage.value = false
  }
}
</script>

<style scoped>
.site-settings {
  padding: 0;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.5;
}

.storage-status-card {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 16px 20px;
  background: #fafafa;
}

.status-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.status-header h4 {
  margin: 0;
  font-size: 15px;
  color: #303133;
}

.status-body p {
  margin: 6px 0;
  font-size: 13px;
  color: #606266;
}
</style>
