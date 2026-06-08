import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || ''

export const useUserStore = defineStore('user', () => {
  // State
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

  // Getters
  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.username || '')
  const role = computed(() => userInfo.value?.role || 'user')
  const licenseKey = computed(() => userInfo.value?.license_key || '')
  const isAdmin = computed(() => role.value === 'admin')

  // Actions
  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
    axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
  }

  const setUserInfo = (info) => {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  const clearUser = () => {
    token.value = ''
    userInfo.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    delete axios.defaults.headers.common['Authorization']
  }

  const login = async (credentials) => {
    const response = await axios.post(`${API_URL}/api/login`, credentials)
    const { token: newToken, user } = response.data
    setToken(newToken)
    setUserInfo(user)
    return user
  }

  const logout = () => {
    clearUser()
  }

  // 初始化axios默认header
  if (token.value) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    username,
    role,
    licenseKey,
    isAdmin,
    login,
    logout,
    setToken,
    setUserInfo,
    clearUser
  }
})
