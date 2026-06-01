<template>
  <div class="ip-geo-manage">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>IP地理位置管理</span>
          <el-button type="primary" @click="handleSave">保存配置</el-button>
        </div>
      </template>

      <!-- 状态卡片 -->
      <el-row :gutter="20" class="status-row">
        <el-col :span="6">
          <el-card shadow="hover" class="status-card">
            <div class="status-item">
              <el-icon :size="24" :color="config.enabled ? '#67c23a' : '#f56c6c'">
                <component :is="config.enabled ? 'CircleCheck' : 'CircleClose'" />
              </el-icon>
              <div class="status-info">
                <div class="status-label">服务状态</div>
                <div class="status-value">{{ config.enabled ? '已启用' : '已禁用' }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="status-card">
            <div class="status-item">
              <el-icon :size="24" :color="database.ipv4_loaded ? '#67c23a' : '#f56c6c'">
                <component :is="database.ipv4_loaded ? 'FolderChecked' : 'FolderRemove'" />
              </el-icon>
              <div class="status-info">
                <div class="status-label">IPv4数据库</div>
                <div class="status-value">{{ database.ipv4_loaded ? '已加载' : '未加载' }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="status-card">
            <div class="status-item">
              <el-icon :size="24" :color="database.ipv6_loaded ? '#67c23a' : '#e6a23c'">
                <component :is="database.ipv6_loaded ? 'FolderChecked' : 'FolderRemove'" />
              </el-icon>
              <div class="status-info">
                <div class="status-label">IPv6数据库</div>
                <div class="status-value">{{ database.ipv6_loaded ? '已加载' : '未加载' }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="status-card">
            <div class="status-item">
              <el-icon :size="24" color="#409eff"><Timer /></el-icon>
              <div class="status-info">
                <div class="status-label">上次更新</div>
                <div class="status-value">{{ config.last_update || '从未更新' }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 配置表单 -->
      <el-form :model="config" label-width="140px" class="config-form">
        <el-form-item label="启用服务">
          <el-switch v-model="config.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>

        <el-form-item label="下载镜像源">
          <el-radio-group v-model="config.mirror" @change="handleMirrorChange">
            <el-radio-button value="gitee">国内 (Gitee)</el-radio-button>
            <el-radio-button value="github">国外 (GitHub)</el-radio-button>
          </el-radio-group>
          <div class="form-tip">切换后将自动更新下方下载地址</div>
        </el-form-item>

        <el-divider content-position="left">IPv4 配置</el-divider>

        <el-form-item label="IPv4下载地址">
          <el-input v-model="config.download_url_v4" placeholder="请输入IPv4数据库下载地址">
            <template #append>
              <el-button @click="handleTestDownload(config.download_url_v4)">测试</el-button>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="IPv4数据库路径">
          <el-input v-model="config.database_path_v4" placeholder="请输入IPv4数据库文件存储路径" />
        </el-form-item>

        <el-divider content-position="left">IPv6 配置</el-divider>

        <el-form-item label="IPv6下载地址">
          <el-input v-model="config.download_url_v6" placeholder="请输入IPv6数据库下载地址">
            <template #append>
              <el-button @click="handleTestDownload(config.download_url_v6)">测试</el-button>
            </template>
          </el-input>
          <div class="form-tip">IPv6数据库为可选配置</div>
        </el-form-item>

        <el-form-item label="IPv6数据库路径">
          <el-input v-model="config.database_path_v6" placeholder="请输入IPv6数据库文件存储路径" />
        </el-form-item>

        <el-divider content-position="left">更新设置</el-divider>

        <el-form-item label="更新周期（天）">
          <el-input-number v-model="config.update_days" :min="1" :max="30" />
          <div class="form-tip">建议7-14天更新一次</div>
        </el-form-item>

        <el-form-item label="IPv4数据库信息" v-if="database.file_exists_v4">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="文件大小">{{ formatFileSize(database.file_size_v4) }}</el-descriptions-item>
            <el-descriptions-item label="文件状态">
              <el-tag type="success">存在</el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </el-form-item>

        <el-form-item label="IPv6数据库信息" v-if="database.file_exists_v6">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="文件大小">{{ formatFileSize(database.file_size_v6) }}</el-descriptions-item>
            <el-descriptions-item label="文件状态">
              <el-tag type="success">存在</el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </el-form-item>
      </el-form>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button type="warning" :loading="updating" @click="handleUpdate">
          <el-icon><Refresh /></el-icon>
          手动更新数据库
        </el-button>
        <el-button @click="handleRefresh">
          <el-icon><RefreshRight /></el-icon>
          刷新状态
        </el-button>
      </div>
    </el-card>

    <!-- IP测试卡片 -->
    <el-card class="box-card" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>IP查询测试</span>
        </div>
      </template>

      <el-form :inline="true">
        <el-form-item label="IP地址">
          <el-input v-model="testIp" placeholder="请输入IP地址（如：8.8.8.8 或 2001:db8::1）" style="width: 350px;" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="testing" @click="handleTestIP">查询</el-button>
        </el-form-item>
      </el-form>

      <el-descriptions v-if="testResult" :column="2" border style="margin-top: 20px;">
        <el-descriptions-item label="IP地址">{{ testResult.ip }}</el-descriptions-item>
        <el-descriptions-item label="简单位置">{{ testResult.location || '未知' }}</el-descriptions-item>
        <el-descriptions-item label="国家">{{ testResult.full_location?.country || '未知' }}</el-descriptions-item>
        <el-descriptions-item label="省份">{{ testResult.full_location?.province || '未知' }}</el-descriptions-item>
        <el-descriptions-item label="城市">{{ testResult.full_location?.city || '未知' }}</el-descriptions-item>
        <el-descriptions-item label="ISP">{{ testResult.full_location?.isp || '未知' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 更新日志卡片 -->
    <el-card class="box-card" style="margin-top: 20px;">
      <template #header>
        <div class="card-header">
          <span>更新日志</span>
        </div>
      </template>

      <el-timeline>
        <el-timeline-item
          v-for="(log, index) in updateLogs"
          :key="index"
          :timestamp="log.time"
          placement="top"
        >
          <el-card>
            <h4>{{ log.action }}</h4>
            <p>{{ log.detail }}</p>
          </el-card>
        </el-timeline-item>
      </el-timeline>

      <el-empty v-if="updateLogs.length === 0" description="暂无更新日志" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  CircleCheck,
  CircleClose,
  FolderChecked,
  FolderRemove,
  Timer,
  Refresh,
  RefreshRight
} from '@element-plus/icons-vue'
import request from '@/utils/request'

// 镜像源地址映射
const mirrorURLs = {
  gitee: {
    v4: 'https://gitee.com/lionsoul/ip2region/raw/master/data/ip2region_v4.xdb',
    v6: 'https://gitee.com/lionsoul/ip2region/raw/master/data/ip2region_v6.xdb'
  },
  github: {
    v4: 'https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v4.xdb',
    v6: 'https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region_v6.xdb'
  }
}

// 配置数据
const config = ref({
  enabled: true,
  mirror: 'gitee',
  download_url_v4: mirrorURLs.gitee.v4,
  download_url_v6: mirrorURLs.gitee.v6,
  update_days: 10,
  last_update: '',
  database_path_v4: 'data/ip2region_v4.xdb',
  database_path_v6: 'data/ip2region_v6.xdb'
})

// 数据库信息
const database = ref({
  file_exists_v4: false,
  file_exists_v6: false,
  file_size_v4: 0,
  file_size_v6: 0,
  ipv4_loaded: false,
  ipv6_loaded: false
})

// 状态
const loading = ref(false)
const updating = ref(false)
const testing = ref(false)

// 测试数据
const testIp = ref('')
const testResult = ref(null)

// 更新日志
const updateLogs = ref([])

// 加载配置
const loadConfig = async () => {
  loading.value = true
  try {
    const res = await request.get('/admin/ip-geo/config')
    if (res.data) {
      config.value = { ...config.value, ...res.data.config }
      database.value = { ...database.value, ...res.data.database }
    }
  } catch (error) {
    ElMessage.error('加载配置失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 保存配置
const handleSave = async () => {
  try {
    await request.put('/admin/ip-geo/config', config.value)
    ElMessage.success('配置保存成功')
  } catch (error) {
    ElMessage.error('保存配置失败: ' + error.message)
  }
}

// 更新数据库
const handleUpdate = async () => {
  try {
    await ElMessageBox.confirm('确定要更新IP地理位置数据库吗？将同时更新IPv4和IPv6数据库。', '确认更新', {
      type: 'warning'
    })
    
    updating.value = true
    await request.post('/admin/ip-geo/update')
    ElMessage.success('数据库更新成功')
    await loadConfig()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('更新失败: ' + error.message)
    }
  } finally {
    updating.value = false
  }
}

// 刷新状态
const handleRefresh = async () => {
  await loadConfig()
  ElMessage.success('状态已刷新')
}

// 镜像源切换
const handleMirrorChange = (val) => {
  const urls = mirrorURLs[val] || mirrorURLs.gitee
  config.value.download_url_v4 = urls.v4
  config.value.download_url_v6 = urls.v6
}

// 测试下载地址
const handleTestDownload = (url) => {
  if (url) {
    window.open(url, '_blank')
  } else {
    ElMessage.warning('请输入下载地址')
  }
}

// 测试IP查询
const handleTestIP = async () => {
  if (!testIp.value) {
    ElMessage.warning('请输入IP地址')
    return
  }

  testing.value = true
  try {
    const res = await request.get('/admin/ip-geo/test', {
      params: { ip: testIp.value }
    })
    testResult.value = res.data
  } catch (error) {
    ElMessage.error('查询失败: ' + error.message)
    testResult.value = null
  } finally {
    testing.value = false
  }
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes) return '-'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

// 加载更新日志
const loadLogs = async () => {
  try {
    const res = await request.get('/admin/ip-geo/logs')
    if (res.data?.list) {
      updateLogs.value = res.data.list.map(item => ({
        time: item.config?.last_update || '未知',
        action: '数据库更新',
        detail: `IPv4: ${item.config?.database_path_v4}, IPv6: ${item.config?.database_path_v6}`
      }))
    }
  } catch (error) {
    console.error('加载日志失败:', error)
  }
}

onMounted(() => {
  loadConfig()
  loadLogs()
})
</script>

<style scoped>
.ip-geo-manage {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-row {
  margin-bottom: 20px;
}

.status-card {
  cursor: default;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-info {
  flex: 1;
}

.status-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.status-value {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.config-form {
  max-width: 700px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.action-buttons {
  margin-top: 20px;
  display: flex;
  gap: 12px;
}
</style>
