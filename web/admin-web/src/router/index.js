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
        meta: { title: '数据看板', icon: 'DataBoard' }
      },
      {
        path: 'stores',
        name: 'Stores',
        component: () => import('@/views/Stores/index.vue'),
        meta: { title: '门店管理', icon: 'Shop' }
      },
      {
        path: 'stores/create',
        name: 'StoreCreate',
        component: () => import('@/views/Stores/Create.vue'),
        meta: { title: '创建门店', hidden: true }
      },
      {
        path: 'stores/:id/edit',
        name: 'StoreEdit',
        component: () => import('@/views/Stores/Edit.vue'),
        meta: { title: '编辑门店', hidden: true }
      },
      {
        path: 'bookings',
        name: 'Bookings',
        component: () => import('@/views/Bookings/index.vue'),
        meta: { title: '预约管理', icon: 'Calendar' }
      },
      {
        path: 'members',
        name: 'Members',
        component: () => import('@/views/Members/index.vue'),
        meta: { title: '会员管理', icon: 'User' }
      },
      {
        path: 'members/create',
        name: 'MemberCreate',
        component: () => import('@/views/Members/Form.vue'),
        meta: { title: '新增会员', hidden: true }
      },
      {
        path: 'members/:id',
        name: 'MemberDetail',
        component: () => import('@/views/Members/Detail.vue'),
        meta: { title: '会员详情', hidden: true }
      },
      {
        path: 'members/:id/edit',
        name: 'MemberEdit',
        component: () => import('@/views/Members/Form.vue'),
        meta: { title: '编辑会员', hidden: true }
      },
      {
        path: 'members/:id/usages',
        name: 'MemberUsages',
        component: () => import('@/views/Members/Usages.vue'),
        meta: { title: '使用记录', hidden: true }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/Users/index.vue'),
        meta: { title: '员工管理', icon: 'UserFilled' }
      },
      {
        path: 'users/create',
        name: 'UserCreate',
        component: () => import('@/views/Users/Create.vue'),
        meta: { title: '新建员工', hidden: true }
      },
      {
        path: 'users/:id/edit',
        name: 'UserEdit',
        component: () => import('@/views/Users/Edit.vue'),
        meta: { title: '编辑员工', hidden: true }
      },
      {
        path: 'payments',
        name: 'Payments',
        component: () => import('@/views/Payments/index.vue'),
        meta: { title: '支付管理', icon: 'Money' }
      },
      {
        path: 'reports',
        name: 'Reports',
        component: () => import('@/views/Reports/index.vue'),
        meta: { title: '报表中心', icon: 'DataAnalysis' }
      },
      {
        path: 'templates',
        name: 'Templates',
        component: () => import('@/views/Templates/index.vue'),
        meta: { title: '时段模板', icon: 'Clock' }
      },
      {
        path: 'templates/create',
        name: 'TemplateCreate',
        component: () => import('@/views/Templates/Form.vue'),
        meta: { title: '创建时段模板', hidden: true }
      },
      {
        path: 'templates/:id/edit',
        name: 'TemplateEdit',
        component: () => import('@/views/Templates/Form.vue'),
        meta: { title: '编辑时段模板', hidden: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '个人资料', icon: 'User', hidden: true }
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
    document.title = `${to.meta.title} - 美甲美睫管理系统`
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
