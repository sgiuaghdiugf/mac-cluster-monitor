<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside width="240px" class="sidebar glass-card">
      <div class="logo">
        <div class="logo-icon">
          <el-icon :size="28"><Monitor /></el-icon>
        </div>
        <span>Mac Monitor</span>
      </div>

      <el-menu
        :default-active="$route.path"
        router
        class="sidebar-menu"
      >
        <el-menu-item index="/dashboard" class="menu-item">
          <el-icon><Monitor /></el-icon>
          <span>监控面板</span>
          <el-badge v-if="userStore.isAdmin && deviceStats.offline > 0" :value="deviceStats.offline" :max="99" type="danger" class="menu-badge" />
        </el-menu-item>

        <el-menu-item index="/devices" class="menu-item">
          <el-icon><Cpu /></el-icon>
          <span>设备管理</span>
        </el-menu-item>

        <el-menu-item index="/users" v-if="userStore.isAdmin" class="menu-item">
          <el-icon><UserFilled /></el-icon>
          <span>用户管理</span>
        </el-menu-item>

        <el-menu-item index="/datascreen" v-if="userStore.isAdmin" class="menu-item">
          <el-icon><DataBoard /></el-icon>
          <span>数据大屏</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container class="main-container">
      <!-- 顶部导航栏 -->
      <el-header class="header glass-card">
        <div class="header-left">
          <h3 class="page-title">{{ pageTitle }}</h3>
        </div>
        <div class="header-right">
          <div class="status-indicator" v-if="userStore.isAdmin">
            <span class="status-dot online"></span>
            <span>{{ deviceStats.online }} 在线</span>
            <span class="status-divider">|</span>
            <span class="status-dot offline"></span>
            <span>{{ deviceStats.offline }} 离线</span>
          </div>
          <div class="status-indicator" v-else>
            <span class="status-dot online"></span>
            <span>我的设备</span>
          </div>

          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <div class="user-avatar">
                <el-icon><User /></el-icon>
              </div>
              <span class="username">{{ userStore.username }}</span>
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu class="glass-dropdown">
                <el-dropdown-item command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade-slide" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Monitor, Cpu, User, UserFilled, ArrowDown, SwitchButton, DataBoard } from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const deviceStats = ref({
  online: 0,
  offline: 0,
  total: 0
})

const pageTitle = computed(() => {
  const titles = {
    '/dashboard': userStore.isAdmin ? '监控面板' : '我的面板',
    '/devices': userStore.isAdmin ? '设备管理' : '我的设备',
    '/users': '用户管理'
  }
  return titles[route.path] || 'Mac Cluster Monitor'
})

// 获取设备统计
const fetchDeviceStats = async () => {
  try {
    const response = await axios.get(`${import.meta.env.VITE_API_URL || ''}/api/devices/stats`)
    deviceStats.value = response.data
  } catch (error) {
    console.error('获取设备统计失败', error)
  }
}

const handleCommand = (command) => {
  if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      userStore.logout()
      ElMessage.success('已退出登录')
      router.push('/login')
    })
  }
}

onMounted(() => {
  fetchDeviceStats()
  setInterval(fetchDeviceStats, 30000)
})

watch(() => route.path, () => {
  fetchDeviceStats()
})
</script>

<style scoped>
.layout-container {
  min-height: 100vh;
  position: relative;
  z-index: 1;
}

/* 侧边栏样式 */
.sidebar {
  background: rgba(15, 23, 42, 0.8) !important;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  position: fixed;
  left: 20px;
  top: 20px;
  bottom: 20px;
  border-radius: 16px;
  z-index: 100;
  width: 240px !important;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.logo {
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #f1f5f9;
  font-size: 18px;
  font-weight: bold;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  gap: 12px;
}

.logo-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #3b82f6, #06b6d4);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 4px 15px rgba(59, 130, 246, 0.4);
}

.sidebar-menu {
  border-right: none !important;
  padding: 10px;
  background: transparent !important;
}

.menu-item {
  margin: 5px 0 !important;
  border-radius: 10px !important;
  height: 48px !important;
  line-height: 48px !important;
  color: #94a3b8 !important;
  position: relative;
  transition: all 0.3s ease !important;
}

.menu-item:hover {
  background: rgba(59, 130, 246, 0.1) !important;
  color: #e2e8f0 !important;
}

.menu-item.is-active {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.2), rgba(6, 182, 212, 0.2)) !important;
  color: #60a5fa !important;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 2px 8px rgba(59, 130, 246, 0.2);
}

.menu-item.is-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 24px;
  background: linear-gradient(180deg, #3b82f6, #06b6d4);
  border-radius: 0 3px 3px 0;
}

.menu-badge {
  margin-left: auto;
  margin-right: 8px;
}

/* 主容器 */
.main-container {
  margin-left: 280px;
}

/* 顶部导航 */
.header {
  background: rgba(15, 23, 42, 0.7) !important;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: sticky;
  top: 20px;
  right: 20px;
  left: 280px;
  z-index: 99;
  margin: 20px 20px 0 0;
  border-radius: 16px;
  height: 64px !important;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.page-title {
  margin: 0;
  color: #f1f5f9;
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #94a3b8;
  padding: 6px 14px;
  background: rgba(15, 23, 42, 0.5);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.status-dot.online {
  background: #10b981;
  box-shadow: 0 0 10px rgba(16, 185, 129, 0.5);
}

.status-dot.offline {
  background: #ef4444;
  box-shadow: 0 0 10px rgba(239, 68, 68, 0.5);
}

.status-divider {
  color: rgba(148, 163, 184, 0.3);
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.user-info {
  cursor: pointer;
  display: flex;
  align-items: center;
  color: #e2e8f0;
  gap: 8px;
  transition: all 0.3s ease;
  padding: 6px 12px;
  border-radius: 10px;
}

.user-info:hover {
  background: rgba(59, 130, 246, 0.1);
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #3b82f6, #06b6d4);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.username {
  font-weight: 500;
}

/* 主内容区 */
.main-content {
  padding: 20px;
  min-height: calc(100vh - 104px);
  position: relative;
  z-index: 1;
}

/* 过渡动画 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>

<style>
/* 全局下拉菜单毛玻璃效果 */
.glass-dropdown {
  background: rgba(15, 23, 42, 0.95) !important;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 12px !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4) !important;
  padding: 8px !important;
}

.glass-dropdown .el-dropdown-menu__item {
  color: #e2e8f0 !important;
  border-radius: 8px !important;
  margin: 4px 0 !important;
  padding: 8px 16px !important;
  transition: all 0.2s ease !important;
}

.glass-dropdown .el-dropdown-menu__item:hover {
  background: rgba(59, 130, 246, 0.2) !important;
  color: #60a5fa !important;
}
</style>
