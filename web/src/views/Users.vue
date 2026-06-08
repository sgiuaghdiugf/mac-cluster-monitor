<template>
  <div class="users-page">
    <!-- 页面头部 -->
    <div class="page-header glass-card">
      <div class="header-left">
        <div class="header-icon">
          <el-icon :size="22"><UserFilled /></el-icon>
        </div>
        <div class="header-text">
          <h2>用户管理</h2>
          <span class="subtitle">管理系统用户与卡密绑定</span>
        </div>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreateDialog" class="create-btn">
          <el-icon><Plus /></el-icon>
          创建用户
        </el-button>
        <div class="view-toggle">
          <button :class="['toggle-btn', { active: viewMode === 'card' }]" @click="viewMode = 'card'">
            <el-icon><Grid /></el-icon>
            <span>卡片</span>
          </button>
          <button :class="['toggle-btn', { active: viewMode === 'table' }]" @click="viewMode = 'table'">
            <el-icon><List /></el-icon>
            <span>表格</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 统计概览 -->
    <div class="stats-row">
      <div class="mini-stat glass-card">
        <div class="mini-stat-value">{{ users.length }}</div>
        <div class="mini-stat-label">总用户数</div>
      </div>
      <div class="mini-stat glass-card">
        <div class="mini-stat-value">{{ users.filter(u => u.role === 'admin').length }}</div>
        <div class="mini-stat-label">管理员</div>
      </div>
      <div class="mini-stat glass-card">
        <div class="mini-stat-value">{{ users.filter(u => u.role === 'user').length }}</div>
        <div class="mini-stat-label">普通用户</div>
      </div>
      <div class="mini-stat glass-card">
        <div class="mini-stat-value">{{ totalDevices }}</div>
        <div class="mini-stat-label">已绑定设备</div>
      </div>
    </div>

    <!-- 用户卡片列表 -->
    <div v-if="viewMode === 'card'" class="users-grid">
      <div
        v-for="user in users"
        :key="user.id"
        class="user-card glass-card"
        :class="{ 'admin-card': user.role === 'admin' }"
      >
        <!-- 卡片头部 -->
        <div class="card-top">
          <div class="user-avatar" :class="user.role">
            <el-icon :size="20"><User /></el-icon>
          </div>
          <div class="user-brief">
            <div class="user-name-row">
              <span class="user-name">{{ user.username }}</span>
              <el-tag
                :type="user.role === 'admin' ? 'danger' : 'info'"
                size="small"
                effect="dark"
                class="role-tag"
              >
                {{ user.role === 'admin' ? '管理员' : '普通用户' }}
              </el-tag>
            </div>
            <div class="user-contact">
            <span>{{ user.created_at }}</span>
          </div>
          </div>
        </div>

        <!-- 卡密信息 -->
        <div class="license-section">
          <div class="license-label">卡密 (License Key)</div>
          <div class="license-value">
            <template v-if="visibleLicenses.has(user.id)">
              <code>{{ user.license_key }}</code>
              <el-button
                type="primary"
                text
                size="small"
                @click="copyLicense(user.license_key)"
                class="copy-btn"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
              <el-button type="primary" text size="small" @click="visibleLicenses.delete(user.id)" class="copy-btn">隐藏</el-button>
            </template>
            <el-button v-else type="primary" text size="small" @click="visibleLicenses.add(user.id)">点击显示卡密</el-button>
          </div>
        </div>

        <!-- 设备信息 -->
        <div class="device-section">
          <div class="device-stat">
            <div class="device-stat-value">{{ user.device_count }}</div>
            <div class="device-stat-label">绑定设备</div>
          </div>
          <div class="device-actions">
            <el-button
              type="primary"
              text
              size="small"
              @click="showBindDialog(user)"
              v-if="user.role === 'user'"
            >
              <el-icon><Link /></el-icon>
              绑定设备
            </el-button>
          </div>
        </div>

        <!-- 底部信息 -->
        <div class="card-footer">
          <span class="footer-time">创建于 {{ user.created_at }}</span>
          <div class="footer-actions">
            <el-button
              type="primary"
              text
              size="small"
              @click="showEditDialog(user)"
            >
              编辑
            </el-button>
            <el-button
              v-if="user.role !== 'admin'"
              type="danger"
              text
              size="small"
              @click="handleDelete(user)"
            >
              删除
            </el-button>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && users.length === 0" class="empty-state">
        <el-empty description="暂无用户数据" :image-size="120" />
      </div>
    </div>

    <!-- 表格视图 -->
    <div v-if="viewMode === 'table'" class="table-view glass-card">
      <el-table
        :data="users"
        v-loading="loading"
        border
        style="width: 100%"
        :header-cell-style="{ background: 'rgba(30, 41, 59, 0.8)', color: '#94a3b8', borderColor: 'rgba(51, 65, 85, 0.6)' }"
        :row-style="{ background: 'rgba(15, 23, 42, 0.4)', color: '#cbd5e1', borderColor: 'rgba(51, 65, 85, 0.3)' }"
        :cell-style="{ borderColor: 'rgba(51, 65, 85, 0.3)' }"
      >
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="phone" label="手机号" min-width="130">
          <template #default="{ row }">
            <span>{{ row.phone || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="卡密" min-width="140">
          <template #default="{ row }">
            <template v-if="visibleLicenses.has(row.id)">
              <code class="table-license">{{ row.license_key }}</code>
              <el-button type="primary" text size="small" @click="visibleLicenses.delete(row.id)" style="margin-left:4px">隐藏</el-button>
            </template>
            <el-button v-else type="primary" text size="small" @click="visibleLicenses.add(row.id)">点击显示</el-button>
          </template>
        </el-table-column>
        <el-table-column label="角色" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : 'info'" size="small" effect="dark">
              {{ row.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="device_count" label="绑定设备" width="90" align="center">
          <template #default="{ row }">
            <span>{{ row.device_count || 0 }}台</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" min-width="160" />
        <el-table-column label="操作" width="240" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text size="small" @click="copyLicense(row.license_key)">
              复制卡密
            </el-button>
            <el-button type="primary" text size="small" @click="showBindDialog(row)" v-if="row.role === 'user'">
              绑定
            </el-button>
            <el-button type="danger" text size="small" @click="handleDelete(row)" v-if="row.role !== 'admin'">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 创建用户对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="创建用户"
      width="480px"
      destroy-on-close
      class="custom-dialog"
    >
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="70px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="createForm.username" placeholder="请输入用户名称" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="createForm.phone" placeholder="请输入手机号（用于登录）" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-radio-group v-model="createForm.role">
            <el-radio-button label="user">普通用户</el-radio-button>
            <el-radio-button label="admin">管理员</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <div class="form-tip">
          卡密将由系统根据手机号自动生成。用户使用手机号和卡密登录。
        </div>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="createLoading">创建</el-button>
      </template>
    </el-dialog>

    <!-- 编辑用户对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑用户"
      width="480px"
      destroy-on-close
      class="custom-dialog"
    >
      <el-form label-width="80px">
        <el-form-item label="用户名">
          <span class="form-text">{{ currentUser?.username }}</span>
        </el-form-item>
        <el-form-item label="卡密">
          <code class="form-license">{{ currentUser?.license_key }}</code>
        </el-form-item>
        <el-form-item label="角色">
          <el-radio-group v-model="editForm.role">
            <el-radio-button label="user">普通用户</el-radio-button>
            <el-radio-button label="admin">管理员</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleUpdateUser" :loading="editLoading">保存</el-button>
      </template>
    </el-dialog>

    <!-- 绑定设备对话框 -->
    <el-dialog
      v-model="bindDialogVisible"
      title="绑定设备"
      width="560px"
      destroy-on-close
      class="custom-dialog"
    >
      <div class="bind-header">
        <span>为用户 <strong>{{ bindUser?.username }}</strong> 绑定设备</span>
        <span class="bind-count">当前已绑定 {{ bindUser?.device_count || 0 }} 台</span>
      </div>
      <div class="bind-list" v-loading="bindLoading">
        <div
          v-for="device in unboundDevices"
          :key="device.id"
          class="bind-item"
          @click="handleBindDevice(device.id)"
        >
          <div class="bind-item-left">
            <div class="bind-device-status" :class="device.is_online ? 'online' : 'offline'"></div>
            <div class="bind-device-info">
              <div class="bind-device-name">{{ device.hostname }}</div>
              <div class="bind-device-meta">{{ device.mac_address }} · {{ device.ip_address || '无IP' }}</div>
            </div>
          </div>
          <el-button type="primary" size="small">绑定</el-button>
        </div>
        <div v-if="!bindLoading && unboundDevices.length === 0" class="bind-empty">
          <el-empty description="暂无可绑定的设备" :image-size="60" />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UserFilled, User, Plus, CopyDocument, Link, Grid, List } from '@element-plus/icons-vue'
import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || ''

// 数据
const users = ref([])
const loading = ref(false)
const viewMode = ref('card') // card | table
const visibleLicenses = reactive(new Set()) // 表格中已展开显示卡密的行ID
const totalDevices = computed(() => users.value.reduce((sum, u) => sum + (u.device_count || 0), 0))

// 创建用户
const createDialogVisible = ref(false)
const createLoading = ref(false)
const createFormRef = ref(null)
const createForm = ref({
  username: '',
  phone: '',
  role: 'user'
})

const createRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 编辑用户
const editDialogVisible = ref(false)
const editLoading = ref(false)
const currentUser = ref(null)
const editForm = ref({ role: 'user' })

// 绑定设备
const bindDialogVisible = ref(false)
const bindLoading = ref(false)
const bindUser = ref(null)
const unboundDevices = ref([])

// 获取用户列表
const fetchUsers = async () => {
  loading.value = true
  try {
    const response = await axios.get(`${API_URL}/api/admin/users`)
    users.value = response.data.users || []
  } catch (error) {
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 复制卡密
const copyLicense = async (key) => {
  try {
    await navigator.clipboard.writeText(key)
    ElMessage.success('卡密已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

// 显示创建对话框
const showCreateDialog = () => {
  createForm.value = { username: '', phone: '', role: 'user' }
  createDialogVisible.value = true
}

// 创建用户（卡密由后端根据手机号MD5生成）
const handleCreate = async () => {
  const valid = await createFormRef.value?.validate().catch(() => false)
  if (!valid) return

  createLoading.value = true
  try {
    const res = await axios.post(`${API_URL}/api/admin/users`, createForm.value)
    ElMessage.success('用户创建成功，卡密：' + (res.data.user?.license_key || ''))
    createDialogVisible.value = false
    fetchUsers()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '创建用户失败')
  } finally {
    createLoading.value = false
  }
}

// 显示编辑对话框
const showEditDialog = (user) => {
  currentUser.value = user
  editForm.value = { role: user.role || 'user' }
  editDialogVisible.value = true
}

// 更新用户
const handleUpdateUser = async () => {
  if (!currentUser.value) return

  editLoading.value = true
  try {
    await axios.put(`${API_URL}/api/admin/users/${currentUser.value.id}`, {
      role: editForm.value.role
    })
    ElMessage.success('用户更新成功')
    editDialogVisible.value = false
    fetchUsers()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '更新用户失败')
  } finally {
    editLoading.value = false
  }
}

// 删除用户
const handleDelete = (user) => {
  ElMessageBox.confirm(
    `确定要删除用户 "${user.username}" 吗? 该用户绑定的设备将自动解绑。`,
    '警告',
    { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
  ).then(async () => {
    try {
      await axios.delete(`${API_URL}/api/admin/users/${user.id}`)
      ElMessage.success('用户删除成功')
      fetchUsers()
    } catch (error) {
      ElMessage.error(error.response?.data?.error || '删除用户失败')
    }
  })
}

// 显示绑定设备对话框
const showBindDialog = async (user) => {
  bindUser.value = user
  bindDialogVisible.value = true
  bindLoading.value = true
  try {
    const response = await axios.get(`${API_URL}/api/admin/devices/unbound`)
    unboundDevices.value = response.data.devices || []
  } catch (error) {
    ElMessage.error('获取未绑定设备失败')
  } finally {
    bindLoading.value = false
  }
}

// 绑定设备
const handleBindDevice = async (deviceId) => {
  try {
    await axios.put(`${API_URL}/api/admin/devices/${deviceId}/bind`, {
      user_id: bindUser.value.id
    })
    ElMessage.success('设备绑定成功')
    bindDialogVisible.value = false
    fetchUsers()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '绑定设备失败')
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.users-page {
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

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 6px 20px rgba(59, 130, 246, 0.35);
}

.header-text h2 {
  margin: 0;
  font-size: 22px;
  color: #f1f5f9;
  font-weight: 700;
  letter-spacing: 0.3px;
}

.subtitle {
  font-size: 13px;
  color: #64748b;
}

.create-btn {
  background: linear-gradient(135deg, #3b82f6, #8b5cf6) !important;
  border: none !important;
  height: 40px;
  padding: 0 24px;
  border-radius: 10px;
  font-weight: 600;
  box-shadow: 0 4px 15px rgba(59, 130, 246, 0.3);
  transition: all 0.3s ease;
}

.create-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(59, 130, 246, 0.4);
}

/* 视图切换 */
.view-toggle {
  display: flex;
  gap: 0;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid rgba(51, 65, 85, 0.5);
}

.toggle-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 7px 14px;
  border: none;
  background: rgba(15, 23, 42, 0.6);
  color: #94a3b8;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s ease;
}

.toggle-btn:hover {
  color: #e2e8f0;
  background: rgba(30, 41, 59, 0.8);
}

.toggle-btn.active {
  background: rgba(59, 130, 246, 0.25);
  color: #60a5fa;
}

/* 统计概览 */
.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
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

.mini-stat-label {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 1px;
}

/* 用户卡片网格 */
.users-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 20px;
}

.user-card {
  padding: 24px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.user-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #3b82f6, #8b5cf6);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.user-card:hover {
  transform: translateY(-3px);
  border-color: rgba(59, 130, 246, 0.3) !important;
  box-shadow: 0 8px 30px rgba(59, 130, 246, 0.15);
}

.user-card:hover::before {
  opacity: 1;
}

.admin-card::before {
  background: linear-gradient(90deg, #ef4444, #f97316);
  opacity: 1;
}

/* 卡片顶部 */
.card-top {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 18px;
}

.user-avatar {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.user-avatar.admin {
  background: linear-gradient(135deg, #ef4444, #f97316);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.user-avatar.user {
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.user-brief {
  flex: 1;
  min-width: 0;
}

.user-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.user-name {
  font-size: 16px;
  font-weight: 600;
  color: #f1f5f9;
}

.role-tag {
  border-radius: 6px;
}

.user-contact {
  font-size: 12px;
  color: #64748b;
}

.no-contact {
  color: #475569;
  font-style: italic;
}

/* 卡密区域 */
.license-section {
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(51, 65, 85, 0.5);
  border-radius: 10px;
  padding: 12px 14px;
  margin-bottom: 16px;
}

.license-label {
  font-size: 11px;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  margin-bottom: 6px;
  font-weight: 600;
}

.license-value {
  display: flex;
  align-items: center;
  gap: 8px;
}

.license-value code {
  flex: 1;
  font-family: 'JetBrains Mono', 'Fira Code', 'Courier New', monospace;
  font-size: 12px;
  color: #60a5fa;
  background: none;
  word-break: break-all;
}

.copy-btn {
  padding: 4px !important;
}

/* 设备区域 */
.device-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-top: 1px solid rgba(51, 65, 85, 0.4);
  border-bottom: 1px solid rgba(51, 65, 85, 0.4);
  margin-bottom: 14px;
}

.device-stat {
  display: flex;
  align-items: baseline;
  gap: 6px;
}

.device-stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #f1f5f9;
}

.device-stat-label {
  font-size: 12px;
  color: #64748b;
}

/* 卡片底部 */
.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-time {
  font-size: 11px;
  color: #475569;
}

.footer-actions {
  display: flex;
  gap: 4px;
}

/* 表单 */
.form-tip {
  font-size: 12px;
  color: #64748b;
  padding: 8px 12px;
  background: rgba(15, 23, 42, 0.6);
  border-radius: 8px;
  border: 1px solid rgba(51, 65, 85, 0.4);
  line-height: 1.6;
}

.form-text {
  color: #e2e8f0;
  font-weight: 500;
}

.form-license {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 12px;
  color: #60a5fa;
  background: rgba(15, 23, 42, 0.6);
  padding: 4px 8px;
  border-radius: 6px;
  border: 1px solid rgba(51, 65, 85, 0.5);
}

/* 绑定对话框 */
.bind-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  color: #e2e8f0;
  font-size: 14px;
}

.bind-count {
  font-size: 12px;
  color: #64748b;
  background: rgba(15, 23, 42, 0.6);
  padding: 4px 10px;
  border-radius: 6px;
}

.bind-list {
  max-height: 400px;
  overflow-y: auto;
}

.bind-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  border-radius: 10px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid rgba(51, 65, 85, 0.4);
}

.bind-item:hover {
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.3);
}

.bind-item-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.bind-device-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.bind-device-status.online {
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
}

.bind-device-status.offline {
  background: #ef4444;
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.5);
}

.bind-device-name {
  font-size: 14px;
  color: #e2e8f0;
  font-weight: 500;
}

.bind-device-meta {
  font-size: 11px;
  color: #64748b;
  font-family: 'Courier New', monospace;
}

.bind-empty {
  padding: 40px 0;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
  grid-column: 1 / -1;
}

/* 表格视图 */
.table-view {
  padding: 0;
  overflow-x: auto;
}

.table-view :deep(.el-table) {
  --el-table-border-color: rgba(51, 65, 85, 0.5);
  --el-table-header-bg-color: rgba(30, 41, 59, 0.8);
  --el-table-row-hover-bg-color: rgba(51, 65, 85, 0.4);
  --el-table-current-row-bg-color: rgba(59, 130, 246, 0.1);
  background: transparent;
  color: #cbd5e1;
}

.table-view :deep(.el-table th.el-table__cell) {
  background: rgba(30, 41, 59, 0.9);
  color: #94a3b8;
  border-bottom: 1px solid rgba(51, 65, 85, 0.6);
}

.table-view :deep(.el-table td.el-table__cell) {
  background: rgba(15, 23, 42, 0.3);
  border-bottom: 1px solid rgba(51, 65, 85, 0.3);
}

.table-view :deep(.el-table--striped .el-table__body tr.el-table__row--striped td.el-table__cell) {
  background: rgba(15, 23, 42, 0.5);
}

.table-view :deep(.el-table .el-table__row:hover > td.el-table__cell) {
  background: rgba(51, 65, 85, 0.4) !important;
}

.table-license {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 12px;
  color: #60a5fa;
  background: none;
}

/* 对话框样式 */
:deep(.el-dialog) {
  background: #1e293b !important;
  border: 1px solid rgba(51, 65, 85, 0.6);
  border-radius: 16px !important;
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid rgba(51, 65, 85, 0.4);
  padding: 20px 24px;
}

:deep(.el-dialog__title) {
  color: #f1f5f9 !important;
  font-weight: 600;
}

:deep(.el-dialog__body) {
  padding: 24px !important;
}

:deep(.el-form-item__label) {
  color: #94a3b8 !important;
}

:deep(.el-input__wrapper) {
  background: #0f172a !important;
  box-shadow: 0 0 0 1px #334155 inset !important;
  border-radius: 8px;
}

:deep(.el-input__inner) {
  color: #e2e8f0 !important;
}

:deep(.el-radio-button__inner) {
  background: #0f172a;
  border-color: #334155;
  color: #94a3b8;
}

:deep(.el-radio-button__original-radio:checked + .el-radio-button__inner) {
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  border-color: #3b82f6;
  color: white;
}
</style>
