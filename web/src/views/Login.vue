<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="login-header">
          <el-icon :size="40" color="#409EFF"><Monitor /></el-icon>
          <h2>Mac Cluster Monitor</h2>
          <p>集群设备监控系统</p>
        </div>
      </template>
      
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        @keyup.enter="handleSubmit"
      >
        <el-form-item label="手机号" prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入手机号"
            :prefix-icon="User"
            size="large"
          />
        </el-form-item>

        <el-form-item label="卡密" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入卡密"
            :prefix-icon="Lock"
            size="large"
            show-password
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleSubmit"
            style="width: 100%"
          >
            登录
          </el-button>
        </el-form-item>
        
        <div class="login-footer">
          <span class="footer-text">Mac Cluster Monitor v2.0</span>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Monitor } from '@element-plus/icons-vue'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref()
const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await userStore.login(form)
        ElMessage.success('登录成功')
        router.push('/')
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '登录失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%);
  position: relative;
}

.login-card {
  width: 420px;
  background: rgba(15, 23, 42, 0.7) !important;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  box-shadow: 
    0 8px 32px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.login-header {
  text-align: center;
  padding: 30px 0 25px;
}

.login-header h2 {
  margin: 18px 0 6px;
  color: #f1f5f9;
  font-size: 26px;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.login-header p {
  color: #94a3b8;
  font-size: 14px;
  margin: 0;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
}

.footer-text {
  color: #475569;
  font-size: 12px;
}
</style>
