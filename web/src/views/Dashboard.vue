<template>
  <div class="dashboard">
    <!-- 管理员视图 -->
    <template v-if="userStore.isAdmin">
      <!-- 顶部统计卡片 -->
      <div class="stats-grid">
        <div class="stat-card glass-card" v-for="stat in adminStatsCards" :key="stat.key">
          <div class="stat-icon" :style="{ background: stat.gradient }">
            <el-icon :size="24"><component :is="stat.icon" /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">
              <span class="number">{{ animatedStats[stat.key] }}</span>
              <span v-if="stat.suffix" class="suffix">{{ stat.suffix }}</span>
            </div>
            <div class="stat-title">{{ stat.title }}</div>
          </div>
        </div>
      </div>

      <!-- 设备列表预览 -->
      <div class="main-section">
        <div class="devices-preview glass-card">
          <div class="section-header">
            <h3>最近设备</h3>
            <el-button type="primary" text @click="$router.push('/devices')">
              查看全部
              <el-icon><arrow-right /></el-icon>
            </el-button>
          </div>
          <div class="device-list">
            <div
              v-for="device in recentDevices"
              :key="device.id"
              class="device-item"
            >
              <div class="device-status" :class="device.is_online ? 'online' : 'offline'"></div>
              <div class="device-info">
                <div class="device-name">{{ device.hostname || device.ip_address }}</div>
                <div class="device-meta">{{ device.ip_address }} · {{ device.username || '未绑定' }}</div>
              </div>
              <div class="device-cpu">
                <el-progress
                  :percentage="Math.round(device.cpu_usage)"
                  :color="getCpuColor(device.cpu_usage)"
                  :stroke-width="6"
                  style="width: 100px"
                />
              </div>
            </div>
            <div v-if="recentDevices.length === 0" class="empty-state">
              <el-empty description="暂无设备数据" :image-size="80" />
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 普通用户视图 -->
    <template v-else>
      <!-- 用户信息头部 -->
      <div class="user-hero glass-card">
        <div class="hero-content">
          <div class="hero-avatar">
            <el-icon :size="32"><User /></el-icon>
          </div>
          <div class="hero-info">
            <h2>{{ userStore.username }}</h2>
            <div class="hero-license">
              <el-icon :size="14"><Key /></el-icon>
              <span>卡密: {{ userStore.licenseKey }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 用户统计卡片 -->
      <div class="stats-grid user-stats">
        <div class="stat-card glass-card profit-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #f59e0b, #ef4444)">
            <el-icon :size="24"><TrendCharts /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">
              <span class="number profit-number">¥{{ userDashboard.today_profit?.toFixed(2) || '0.00' }}</span>
            </div>
            <div class="stat-title">今日盈利</div>
          </div>
          <div class="profit-badge">实时更新</div>
        </div>

        <div class="stat-card glass-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #10b981, #06b6d4)">
            <el-icon :size="24"><Wallet /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">
              <span class="number">¥{{ userDashboard.total_profit?.toFixed(2) || '0.00' }}</span>
            </div>
            <div class="stat-title">累计盈利</div>
          </div>
        </div>

        <div class="stat-card glass-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #3b82f6, #8b5cf6)">
            <el-icon :size="24"><Monitor /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">
              <span class="number">{{ userDashboard.device_count || 0 }}</span>
              <span class="suffix">台</span>
            </div>
            <div class="stat-title">我的设备</div>
          </div>
        </div>

        <div class="stat-card glass-card">
          <div class="stat-icon" style="background: linear-gradient(135deg, #10b981, #38ef7d)">
            <el-icon :size="24"><CircleCheck /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">
              <span class="number">{{ userDashboard.online_count || 0 }}</span>
              <span class="suffix">台</span>
            </div>
            <div class="stat-title">在线设备</div>
          </div>
        </div>
      </div>

      <!-- 用户设备列表 -->
      <div class="user-devices glass-card">
        <div class="section-header">
          <h3>我的设备</h3>
          <div class="device-summary">
            <span class="online-text">{{ userDashboard.online_count || 0 }} 在线</span>
            <span class="divider">/</span>
            <span class="offline-text">{{ userDashboard.offline_count || 0 }} 离线</span>
          </div>
        </div>
        <div class="device-grid">
          <div
            v-for="device in userDashboard.devices"
            :key="device.id"
            class="device-card"
            :class="{ offline: !device.is_online }"
          >
            <div class="device-card-header">
              <div class="device-status-dot" :class="device.is_online ? 'online' : 'offline'"></div>
              <span class="device-card-status">{{ device.is_online ? '在线' : '离线' }}</span>
            </div>
            <div class="device-card-body">
              <div class="device-card-name">{{ device.hostname }}</div>
              <div class="device-card-ip">{{ device.ip_address }}</div>
            </div>
            <div class="device-card-footer">
              <div class="device-metric">
                <span class="metric-label">CPU</span>
                <el-progress
                  :percentage="Math.round(device.cpu_usage)"
                  :color="getCpuColor(device.cpu_usage)"
                  :stroke-width="4"
                  style="flex: 1"
                />
              </div>
              <div class="device-metric" v-if="device.memory_total > 0">
                <span class="metric-label">MEM</span>
                <el-progress
                  :percentage="Math.round((device.memory_used / device.memory_total) * 100)"
                  :stroke-width="4"
                  style="flex: 1"
                />
              </div>
            </div>
          </div>
          <div v-if="!userDashboard.devices || userDashboard.devices.length === 0" class="empty-state">
            <el-empty description="暂无绑定设备" :image-size="80" />
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Monitor, Cpu, Warning, CircleCheck, ArrowRight,
  Timer, DataLine, TrendCharts, UserFilled, User, Key, Wallet
} from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'
import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || ''
const userStore = useUserStore()

// ============ 管理员数据 ============
const adminStatsCards = ref([
  { key: 'total', title: '设备总数', icon: Monitor, gradient: 'linear-gradient(135deg, #667eea, #764ba2)', suffix: '台' },
  { key: 'online', title: '在线设备', icon: CircleCheck, gradient: 'linear-gradient(135deg, #11998e, #38ef7d)', suffix: '台' },
  { key: 'offline', title: '离线设备', icon: Warning, gradient: 'linear-gradient(135deg, #eb3349, #f45c43)', suffix: '台' },
  { key: 'onlineRate', title: '在线率', icon: Cpu, gradient: 'linear-gradient(135deg, #3b82f6, #06b6d4)', suffix: '%' }
])

const animatedStats = ref({ total: 0, online: 0, offline: 0, onlineRate: 0 })
const deviceStats = ref({ total: 0, online: 0, offline: 0, online_rate: 0 })
const recentDevices = ref([])

// ============ 普通用户数据 ============
const userDashboard = ref({
  device_count: 0,
  online_count: 0,
  offline_count: 0,
  today_profit: 0,
  total_profit: 0,
  devices: []
})

let refreshTimer = null

// ============ 管理员方法 ============
const fetchAdminStats = async () => {
  try {
    const response = await axios.get(`${API_URL}/api/devices/stats`)
    deviceStats.value = response.data
    animateNumber('total', response.data.total)
    animateNumber('online', response.data.online)
    animateNumber('offline', response.data.offline)
    animateNumber('onlineRate', parseFloat(response.data.online_rate).toFixed(1))
  } catch (error) {
    console.error('获取统计失败', error)
  }
}

const fetchRecentDevices = async () => {
  try {
    const response = await axios.get(`${API_URL}/api/devices`)
    recentDevices.value = (response.data.devices || []).slice(0, 6)
  } catch (error) {
    console.error('获取设备列表失败', error)
  }
}

const animateNumber = (key, target) => {
  const duration = 1000
  const steps = 60
  const stepValue = (target - animatedStats.value[key]) / steps
  let currentStep = 0
  const timer = setInterval(() => {
    currentStep++
    if (currentStep >= steps) {
      animatedStats.value[key] = target
      clearInterval(timer)
    } else {
      animatedStats.value[key] = Math.round((animatedStats.value[key] + stepValue) * 10) / 10
    }
  }, duration / steps)
}

// ============ 普通用户方法 ============
const fetchUserDashboard = async () => {
  try {
    const response = await axios.get(`${API_URL}/api/user/dashboard`)
    userDashboard.value = response.data
  } catch (error) {
    console.error('获取用户面板失败', error)
  }
}

// ============ 通用方法 ============
const getCpuColor = (usage) => {
  if (usage > 90) return '#ef4444'
  if (usage > 70) return '#f59e0b'
  return '#10b981'
}

onMounted(async () => {
  if (userStore.isAdmin) {
    await Promise.all([fetchAdminStats(), fetchRecentDevices()])
  } else {
    await fetchUserDashboard()
  }

  // 定时刷新
  refreshTimer = setInterval(async () => {
    if (userStore.isAdmin) {
      await Promise.all([fetchAdminStats(), fetchRecentDevices()])
    } else {
      await fetchUserDashboard()
    }
  }, 5000)
})

onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer)
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

/* 统计卡片网格 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 24px;
  position: relative;
  overflow: hidden;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-right: 16px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #f1f5f9;
  margin-bottom: 4px;
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.number {
  background: linear-gradient(135deg, #60a5fa, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.profit-number {
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-size: 32px;
}

.suffix {
  font-size: 16px;
  color: #94a3b8;
  font-weight: 500;
}

.stat-title {
  font-size: 13px;
  color: #94a3b8;
  font-weight: 500;
}

/* 盈利卡片特效 */
.profit-card {
  border: 1px solid rgba(251, 191, 36, 0.2) !important;
}

.profit-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 10px;
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
  font-weight: 600;
  animation: glow-pulse 2s ease-in-out infinite;
}

@keyframes glow-pulse {
  0%, 100% { box-shadow: 0 0 5px rgba(251, 191, 36, 0.2); }
  50% { box-shadow: 0 0 15px rgba(251, 191, 36, 0.4); }
}

/* 用户英雄区域 */
.user-hero {
  padding: 28px 32px;
  margin-bottom: 20px;
  border: 1px solid rgba(59, 130, 246, 0.15) !important;
}

.hero-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.hero-avatar {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 8px 24px rgba(59, 130, 246, 0.3);
}

.hero-info h2 {
  margin: 0 0 6px;
  color: #f1f5f9;
  font-size: 24px;
  font-weight: 700;
}

.hero-license {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #64748b;
  font-size: 13px;
  font-family: 'Courier New', monospace;
  background: rgba(15, 23, 42, 0.6);
  padding: 4px 12px;
  border-radius: 8px;
  border: 1px solid rgba(100, 116, 139, 0.2);
}

/* 主内容区 */
.main-section {
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.devices-preview {
  padding: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h3 {
  margin: 0;
  font-size: 18px;
  color: #f1f5f9;
  font-weight: 600;
}

/* 设备列表 */
.device-list {
  max-height: 400px;
  overflow-y: auto;
}

.device-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border-radius: 10px;
  transition: all 0.2s ease;
  cursor: pointer;
  margin-bottom: 8px;
}

.device-item:hover {
  background: rgba(59, 130, 246, 0.1);
}

.device-status {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-right: 12px;
  animation: pulse 2s infinite;
}

.device-status.online {
  background: #10b981;
  box-shadow: 0 0 10px rgba(16, 185, 129, 0.5);
}

.device-status.offline {
  background: #ef4444;
  box-shadow: 0 0 10px rgba(239, 68, 68, 0.5);
}

.device-info {
  flex: 1;
}

.device-name {
  font-size: 14px;
  color: #e2e8f0;
  font-weight: 500;
  margin-bottom: 2px;
}

.device-meta {
  font-size: 12px;
  color: #94a3b8;
}

.device-cpu {
  margin-left: 12px;
}

/* 用户设备网格 */
.user-devices {
  padding: 24px;
}

.device-summary {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.online-text { color: #10b981; }
.offline-text { color: #ef4444; }
.divider { color: #475569; }

.device-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.device-card {
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(51, 65, 85, 0.5);
  border-radius: 14px;
  padding: 18px;
  transition: all 0.3s ease;
}

.device-card:hover {
  border-color: rgba(59, 130, 246, 0.3);
  box-shadow: 0 4px 20px rgba(59, 130, 246, 0.1);
  transform: translateY(-2px);
}

.device-card.offline {
  opacity: 0.7;
}

.device-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
}

.device-status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.device-status-dot.online {
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
}

.device-status-dot.offline {
  background: #ef4444;
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.5);
}

.device-card-status {
  font-size: 12px;
  font-weight: 600;
  color: #94a3b8;
}

.device-card-body {
  margin-bottom: 14px;
}

.device-card-name {
  font-size: 16px;
  font-weight: 600;
  color: #f1f5f9;
  margin-bottom: 4px;
}

.device-card-ip {
  font-size: 12px;
  color: #64748b;
  font-family: 'Courier New', monospace;
}

.device-card-footer {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.device-metric {
  display: flex;
  align-items: center;
  gap: 8px;
}

.metric-label {
  font-size: 11px;
  color: #64748b;
  font-weight: 600;
  min-width: 28px;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  grid-column: 1 / -1;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>
