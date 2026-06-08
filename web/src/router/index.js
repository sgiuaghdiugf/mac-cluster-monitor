import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('../views/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '监控面板', icon: 'Monitor' }
      },
      {
        path: '/devices',
        name: 'Devices',
        component: () => import('../views/Devices.vue'),
        meta: { title: '设备管理', icon: 'Cpu' }
      },
      {
        path: '/users',
        name: 'Users',
        component: () => import('../views/Users.vue'),
        meta: { title: '用户管理', icon: 'User', admin: true }
      }
    ]
  },
  {
    path: '/datascreen',
    name: 'DataScreen',
    component: () => import('../views/DataScreen.vue'),
    meta: { title: '数据大屏', admin: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (!to.meta.public && !userStore.token) {
    next('/login')
  } else if (to.path === '/login' && userStore.token) {
    next('/')
  } else if (to.meta.admin && !userStore.isAdmin) {
    next('/')
  } else {
    next()
  }
})

export default router
