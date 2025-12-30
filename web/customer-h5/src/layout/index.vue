<template>
  <div class="layout">
    <!-- 主内容区 -->
    <div class="main-content">
      <router-view />
    </div>
    
    <!-- 底部导航 -->
    <van-tabbar v-model="activeTab" route>
      <van-tabbar-item
        v-for="item in tabbarItems"
        :key="item.name"
        :to="item.path"
        :icon="item.icon"
        @click="handleTabClick(item)"
      >
        {{ item.title }}
      </van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const tabbarItems = [
  {
    name: 'Home',
    path: '/home',
    title: '首页',
    icon: 'home-o'
  },
  {
    name: 'Booking',
    path: '/booking',
    title: '预约',
    icon: 'calendar-o'
  },
  {
    name: 'Member',
    path: '/member',
    title: '会员',
    icon: 'vip-card-o'
  },
  {
    name: 'Profile',
    path: '/profile',
    title: '我的',
    icon: 'user-o'
  }
]

const activeTab = computed(() => {
  const currentRoute = route.name
  const tabbarItem = tabbarItems.find(item => item.name === currentRoute)
  return tabbarItem ? tabbarItems.indexOf(tabbarItem) : 0
})

// 处理tab点击事件
const handleTabClick = (item) => {
  console.log('Tab clicked:', item)
  router.push(item.path)
}
</script>

<style lang="scss" scoped>
.layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 50px; // 为底部导航留出空间
  box-sizing: border-box;
}

// 让tabbar使用flex布局而不是fixed定位
:deep(.van-tabbar) {
  flex-shrink: 0; // 防止tabbar被压缩
}
</style>
