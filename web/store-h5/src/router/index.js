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
        meta: { title: '数据看板', icon: 'chart-trending-o' }
      },
      {
        path: 'staff',
        name: 'Staff',
        component: () => import('@/views/Staff/index.vue'),
        meta: { title: '员工管理', icon: 'friends-o' }
      },
      {
        path: 'staff/create',
        name: 'StaffCreate',
        component: () => import('@/views/Staff/Create.vue'),
        meta: { title: '添加员工', hidden: true }
      },
      {
        path: 'members',
        name: 'Members',
        component: () => import('@/views/Members/index.vue'),
        meta: { title: '会员管理', icon: 'vip-card-o' }
      },
      {
        path: 'members/create',
        name: 'MemberCreate',
        component: () => import('@/views/Members/Create.vue'),
        meta: { title: '创建会员', hidden: true }
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
    document.title = `${to.meta.title} - 店长端`
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
