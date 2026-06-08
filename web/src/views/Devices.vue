<template>
  <div class="devices-page">
    <!-- 页面头部 -->
    <div class="page-header glass-card">
      <div class="header-left">
        <div class="header-icon">
          <el-icon :size="22"><Monitor /></el-icon>
        </div>
        <div class="header-text">
          <h2>{{ userStore.isAdmin ? '设备管理' : '我的设备' }}</h2>
          <span class="subtitle">{{ userStore.isAdmin ? '管理所有集群设备' : '查看已绑定的设备' }}</span>
        </div>
      </div>
      <div class="header-actions">
        <el-button @click="fetchDevices" :loading="loading" class="refresh-btn">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 统计概览 (仅管理员) -->
    <div class="stats-row" v-if="userStore.isAdmin">
      <div class="mini-stat glass-card">
        <div class="mini-stat-value">{{ stats.total }}</div>
        <div class="mini-stat-label">设备总数</div>
      </div>
      <div class="mini-stat glass-card online">
        <div class="mini-stat-value">{{ stats.online }}</div>
        <div class="mini-stat-label">在线</div>
      </div>
      <div class="mini-stat glass-card offline">
        <div class="mini-stat-value">{{ stats.offline }}</div>
        <div class="mini-stat-label">离线</div>
      </div>
      <div class="mini-stat glass-card">
        <div class="mini-stat-value">{{ stats.unbound }}</div>
        <div class="mini-stat-label">未绑定</div>
      </div>
    </div>

    <!-- 设备表格 -->
    <div class="table-section glass-card">
      <el-table :data="devices" v-loading="loading" class="custom-table">
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <div class="status-indicator" :class="row.is_online ? 'online' : 'offline'">
              <div class="status-dot"></div>
              <span>{{ row.is_online ? '在线' : '离线' }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="hostname" label="设备名称" min-width="130">
          <template #default="{ row }">
            <div class="device-name-cell">
              <span class="device-name">{{ row.hostname }}</span>
              <span class="device-mac">{{ row.mac_address }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="ip_address" label="IP地址" width="140">
          <template #default="{ row }">
            <code class="ip-code">{{ row.ip_address || '--' }}</code>
          </template>
        </el-table-column>

        <el-table-column label="绑定用户" min-width="120" v-if="userStore.isAdmin">
          <template #default="{ row }">
            <div class="user-bind-cell" v-if="row.username">
              <el-icon :size="14"><User /></el-icon>
              <span>{{ row.username }}</span>
            </div>
            <span v-else class="unbound-text">未绑定</span>
          </template>
        </el-table-column>

        <el-table-column label="CPU" width="140">
          <template #default="{ row }">
            <div class="progress-cell" v-if="row.is_online">
              <el-progress
                :percentage="Math.round(row.cpu_usage)"
                :color="getCpuColor(row.cpu_usage)"
                :stroke-width="6"
              />
            </div>
            <span v-else class="text-muted">--</span>
          </template>
        </el-table-column>

        <el-table-column label="内存" width="160">
          <template #default="{ row }">
            <div class="memory-cell" v-if="row.is_online && row.memory_total > 0">
              <el-progress
                :percentage="Math.round((row.memory_used / row.memory_total) * 100)"
                :color="getMemoryColor(row)"
                :stroke-width="6"
              />
              <span class="memory-text">{{ formatMemory(row.memory_used) }}/{{ formatMemory(row.memory_total) }}</span>
            </div>
            <span v-else class="text-muted">--</span>
          </template>
        </el-table-column>

        <el-table-column label="最后在线" width="160">
          <template #default="{ row }">
            <span class="time-text">{{ row.last_seen_at ? formatTime(row.last_seen_at) : '从未' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180" fixed="right" v-if="userStore.isAdmin">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button
                v-if="!row.username"
                type="primary"
                text
                size="small"
                @click="showBindDeviceDialog(row)"
              >
                绑定
              </el-button>
              <el-button
                v-if="row.username"
                type="warning"
                text
                size="small"
                @click="handleUnbindDevice(row)"
              >
                解绑
              </el-button>
              <el-button
                type="danger"
                text
                size="small"
                @click="handleDelete(row)"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="!loading && devices.length === 0" class="empty-state">
        <el-empty description="暂无设备数据" :image-size="100" />
      </div>
    </div>

    <!-- 绑定设备对话框 -->
    <el-dialog
      v-model="bindDialogVisible"
      title="绑定设备到用户"
      width="440px"
      destroy-on-close
      class="custom-dialog"
    >
      <div class="bind-device-info">
        <span>设备: <strong>{{ bindDevice?.hostname }}</strong></span>
      </div>
      <el-select
        v-model="bindUserId"
        placeholder="请选择用户"
        style="width: 100%; margin-top: 16px"
        filterable
      >
        <el-option
          v-for="user in userList"
          :key="user.id"
          :label="`${user.username} (${user.device_count}台设备)`"
          :value="user.id"
        />
      </el-select>
      <template #footer>
        <el-button @click="bindDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleBindDevice" :loading="bindLoading" :disabled="!bindUserId">绑定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Monitor, Refresh, User } from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'
import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || ''
const userStore = useUserStore()

const loading = ref(false)
const devices = ref([])
const stats = ref({ total: 0, online: 0, offline: 0, unbound: 0 })

// 绑定设备
const bindDialogVisible = ref(false)
const bindLoading = ref(false)
const bindDevice = ref(null)
const bindUserId = ref(null)
const userList = ref([])

// 获取设备列表
const fetchDevices = async () => {
  loading.value = true
  try {
    const response = await axios.get(`${API_URL}/api/devices`)
    devices.value = response.data.devices || []

    // 计算统计
    if (userStore.isAdmin) {
      const total = devices.value.length
      const online = devices.value.filter(d => d.is_online).length
      const unbound = devices.value.filter(d => !d.username).length
      stats.value = { total, online, offline: total - online, unbound }
    }
  } catch (error) {
    ElMessage.error('获取设备列表失败')
  } finally {
    loading.value = false
  }
}

// 显示绑定设备对话框
const showBindDeviceDialog = async (device) => {
  bindDevice.value = device
  bindUserId.value = null
  bindDialogVisible.value = true

  try {
    const response = await axios.get(`${API_URL}/api/admin/users`)
    userList.value = (response.data.users || []).filter(u => u.role === 'user')
  } catch (error) {
    ElMessage.error('获取用户列表失败')
  }
}

// 绑定设备
const handleBindDevice = async () => {
  if (!bindUserId.value || !bindDevice.value) return

  bindLoading.value = true
  try {
    await axios.put(`${API_URL}/api/admin/devices/${bindDevice.value.id}/bind`, {
      user_id: bindUserId.value
    })
    ElMessage.success('设备绑定成功')
    bindDialogVisible.value = false
    fetchDevices()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '绑定失败')
  } finally {
    bindLoading.value = false
  }
}

// 解绑设备
const handleUnbindDevice = async (device) => {
  try {
    await ElMessageBox.confirm(
      `确定要解绑设备 "${device.hostname}" 与用户 "${device.username}" 的绑定吗?`,
      '解绑确认',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await axios.put(`${API_URL}/api/admin/devices/${device.id}/unbind`)
    ElMessage.success('设备解绑成功')
    fetchDevices()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '解绑失败')
    }
  }
}

// 删除设备
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除设备 "${row.hostname}" (${row.mac_address}) 吗?`,
      '删除确认',
      { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
    )
    await axios.delete(`${API_URL}/api/admin/devices/${row.id}`)
    ElMessage.success('删除成功')
    fetchDevices()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

// 格式化内存
const formatMemory = (mb) => {
  if (!mb) return '0 MB'
  if (mb >= 1024) return (mb / 1024).toFixed(1) + ' GB'
  return mb + ' MB'
}

// 获取内存颜色
const getMemoryColor = (row) => {
  if (!row.memory_total) return '#10b981'
  const pct = (row.memory_used / row.memory_total) * 100
  if (pct > 90) return '#ef4444'
  if (pct > 70) return '#f59e0b'
  return '#10b981'
}

// 获取CPU颜色
const getCpuColor = (usage) => {
  if (usage > 90) return '#ef4444'
  if (usage > 70) return '#f59e0b'
  return '#10b981'
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchDevices()
})
</script>

<style scoped>
.devices-page {
  padding: 0;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 28px;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background: linear-gradient(135deg, #10b981, #06b6d4);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 6px 20px rgba(16, 185, 129, 0.35);
}

.header-text h2 {
  margin: 0;
  font-size: 22px;
  color: #f1f5f9;
  font-weight: 700;
}

.subtitle {
  font-size: 13px;
  color: #64748b;
}

.refresh-btn {
  background: rgba(15, 23, 42, 0.6) !important;
  border: 1px solid rgba(51, 65, 85, 0.5) !important;
  color: #94a3b8 !important;
  border-radius: 10px;
}

.refresh-btn:hover {
  border-color: rgba(59, 130, 246, 0.4) !important;
  color: #60a5fa !important;
}

/* 统计概览 */
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.mini-stat {
  padding: 20px;
  text-align: center;
}

.mini-stat-value {
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #60a5fa, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 4px;
}

.mini-stat.online .mini-stat-value {
  background: linear-gradient(135deg, #10b981, #38ef7d);
  -webkit-background-clip: text;
  background-clip: text;
}

.mini-stat.offline .mini-stat-value {
  background: linear-gradient(135deg, #ef4444, #f97316);
  -webkit-background-clip: text;
  background-clip: text;
}

.mini-stat-label {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 1px;
}

/* 表格区域 */
.table-section {
  padding: 24px;
  overflow: hidden;
}

/* 状态指示器 */
.status-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.status-indicator.online .status-dot {
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
}

.status-indicator.online {
  color: #10b981;
}

.status-indicator.offline .status-dot {
  background: #ef4444;
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.5);
}

.status-indicator.offline {
  color: #ef4444;
}

/* 设备名称 */
.device-name-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.device-name {
  font-size: 14px;
  color: #e2e8f0;
  font-weight: 500;
}

.device-mac {
  font-size: 11px;
  color: #64748b;
  font-family: 'Courier New', monospace;
}

.ip-code {
  background: rgba(15, 23, 42, 0.6);
  padding: 2px 8px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #60a5fa;
  border: 1px solid rgba(51, 65, 85, 0.4);
}

/* 绑定用户 */
.user-bind-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #e2e8f0;
  font-size: 13px;
}

.unbound-text {
  color: #475569;
  font-style: italic;
  font-size: 12px;
}

/* 进度条 */
.progress-cell {
  width: 100%;
}

.memory-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.memory-text {
  font-size: 11px;
  color: #64748b;
}

.time-text {
  font-size: 12px;
  color: #94a3b8;
}

.text-muted {
  color: #475569;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 4px;
}

/* 绑定对话框 */
.bind-device-info {
  font-size: 14px;
  color: #e2e8f0;
  padding: 12px;
  background: rgba(15, 23, 42, 0.6);
  border-radius: 8px;
  border: 1px solid rgba(51, 65, 85, 0.4);
}

.empty-state {
  display: flex;
  justify-content: center;
  padding: 60px 0;
}

/* 表格样式 */
:deep(.el-table) {
  background: transparent !important;
  --el-table-bg-color: transparent;
  --el-table-tr-bg-color: transparent;
  --el-table-header-bg-color: rgba(15, 23, 42, 0.5);
  color: #e2e8f0;
}

:deep(.el-table th) {
  background: rgba(15, 23, 42, 0.8) !important;
  color: #94a3b8 !important;
  font-weight: 600;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid rgba(59, 130, 246, 0.15) !important;
}

:deep(.el-table td) {
  background: transparent !important;
  border-bottom: 1px solid rgba(51, 65, 85, 0.3) !important;
  color: #e2e8f0;
}

:deep(.el-table tr:hover > td) {
  background: rgba(59, 130, 246, 0.08) !important;
}

:deep(.el-table__fixed-right),
:deep(.el-table__fixed-right .el-table__cell) {
  background: rgba(15, 23, 42, 0.9) !important;
}

:deep(.el-table__fixed-right tr:hover > td) {
  background: rgba(59, 130, 246, 0.12) !important;
}

:deep(.el-empty__description p) {
  color: #64748b;
}

/* 对话框样式 */
:deep(.el-dialog) {
  background: #1e293b !important;
  border: 1px solid rgba(51, 65, 85, 0.6);
  border-radius: 16px !important;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid rgba(51, 65, 85, 0.4);
}

:deep(.el-dialog__title) {
  color: #f1f5f9 !important;
  font-weight: 600;
}

:deep(.el-form-item__label) {
  color: #94a3b8 !important;
}

:deep(.el-select .el-input__wrapper) {
  background: #0f172a !important;
  box-shadow: 0 0 0 1px #334155 inset !important;
}

:deep(.el-select .el-input__inner) {
  color: #e2e8f0 !important;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>
