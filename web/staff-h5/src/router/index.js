import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '工作台', icon: 'home-o' }
      },
      {
        path: 'bookings',
        name: 'Bookings',
        component: () => import('@/views/Bookings/index.vue'),
        meta: { title: '我的日程', icon: 'calendar-o' }
      },
      {
        path: 'scanner',
        name: 'Scanner',
        component: () => import('@/views/Scanner.vue'),
        meta: { title: '扫码核销', icon: 'scan' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '个人中心', icon: 'user-o' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/404.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 美甲师端`
  }
  
  // 检查登录状态
  if (to.path !== '/login') {
    if (!userStore.token) {
      next('/login')
      return
    }
  } else {
    if (userStore.token) {
      next('/')
      return
    }
  }
  
  next()
})

export default router
