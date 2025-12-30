import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/home',
    children: [
      {
        path: 'home',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
        meta: { title: '首页', icon: 'home-o' }
      },
      {
        path: 'booking',
        name: 'Booking',
        component: () => import('@/views/Booking/index.vue'),
        meta: { title: '预约', icon: 'calendar-o' }
      },
      {
        path: 'booking/create',
        name: 'BookingCreate',
        component: () => import('@/views/Booking/Create.vue'),
        meta: { title: '创建预约', hidden: true }
      },
      {
        path: 'member',
        name: 'Member',
        component: () => import('@/views/Member/index.vue'),
        meta: { title: '会员', icon: 'vip-card-o' }
      },
      {
        path: 'member/code/:id',
        name: 'MemberCode',
        component: () => import('@/views/Member/Code.vue'),
        meta: { title: '会员码', hidden: true }
      },
      {
        path: 'booking/:id/payment',
        name: 'BookingPayment',
        component: () => import('@/views/Booking/Payment.vue'),
        meta: { title: '支付押金', hidden: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: '我的', icon: 'user-o' }
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录', hidden: true }
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
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 美甲美睫预约`
  }
  
  next()
})

export default router
