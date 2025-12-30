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
      >
        {{ item.title }}
      </van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const tabbarItems = [
  {
    name: 'Dashboard',
    path: '/dashboard',
    title: '工作台',
    icon: 'home-o'
  },
  {
    name: 'Bookings',
    path: '/bookings',
    title: '我的日程',
    icon: 'calendar-o'
  },
  {
    name: 'Scanner',
    path: '/scanner',
    title: '扫码核销',
    icon: 'scan'
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
}
</style>
