<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '200px'" class="sidebar">
      <div class="logo">
        <div class="logo-icon" v-if="!isCollapse">💅</div>
        <span v-if="!isCollapse">美甲美睫管理系统</span>
        <span v-else>💅</span>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :unique-opened="true"
        router
        class="sidebar-menu"
      >
        <el-menu-item 
          v-for="route in menuRoutes"
          :key="route.path"
          :index="route.path"
        >
          <el-icon>
            <component :is="getIconComponent(route.meta.icon)" />
          </el-icon>
          <template #title>{{ route.meta.title }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 主内容区 -->
    <el-container>
      <!-- 顶部导航 -->
      <el-header class="header">
        <div class="header-left">
          <el-button
            type="text"
            @click="toggleCollapse"
            class="collapse-btn"
          >
            <el-icon>
              <Expand v-if="isCollapse" />
              <Fold v-else />
            </el-icon>
          </el-button>
          
          <el-breadcrumb separator="/">
            <el-breadcrumb-item
              v-for="item in breadcrumbList"
              :key="item.path"
              :to="item.path"
            >
              {{ item.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" :src="userStore.userInfo.avatar">
                {{ userStore.userInfo.username?.charAt(0) || 'A' }}
              </el-avatar>
              <span class="username">{{ userStore.userInfo.username || '管理员' }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人资料</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主内容 -->
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessageBox } from 'element-plus'
import { 
  DataBoard, 
  Shop, 
  Calendar, 
  User, 
  Avatar, 
  Money, 
  DataAnalysis, 
  Setting 
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)

// 图标组件映射
const iconMap = {
  DataBoard,
  Shop,
  Calendar,
  User,
  Avatar,
  Money,
  DataAnalysis,
  Setting
}

// 获取图标组件
const getIconComponent = (iconName) => {
  return iconMap[iconName] || DataBoard
}

// 菜单路由
const menuRoutes = computed(() => {
  const routes = router.getRoutes()
  const mainRoute = routes.find(route => route.path === '/')
  if (mainRoute && mainRoute.children) {
    return mainRoute.children.filter(child => 
      child.meta && 
      child.meta.title && 
      !child.meta.hidden
    )
  }
  return []
})

// 当前激活的菜单
const activeMenu = computed(() => route.path)

// 面包屑导航
const breadcrumbList = computed(() => {
  const matched = route.matched.filter(item => item.meta && item.meta.title)
  return matched.map(item => ({
    path: item.path,
    title: item.meta.title
  }))
})

// 切换侧边栏折叠状态
const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

// 处理用户下拉菜单命令
const handleCommand = async (command) => {
  switch (command) {
    case 'profile':
      // 跳转到个人资料页面
      router.push('/profile')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        await userStore.logoutAction()
        router.push('/login')
      } catch (error) {
        // 用户取消
      }
      break
  }
}
</script>

<style lang="scss" scoped>
.layout-container {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  transition: width 0.3s;
  
  .logo {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 16px;
    font-weight: bold;
    border-bottom: 1px solid #434a50;
    
    .logo-icon {
      font-size: 24px;
      margin-right: 8px;
    }
  }
  
  .sidebar-menu {
    border: none;
    background-color: #304156;
    
    :deep(.el-menu-item) {
      color: #bfcbd9;
      
      &:hover {
        background-color: #263445;
        color: white;
      }
      
      &.is-active {
        background-color: #409eff;
        color: white;
      }
    }
  }
}

.header {
  background-color: white;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  
  .header-left {
    display: flex;
    align-items: center;
    
    .collapse-btn {
      margin-right: 20px;
      font-size: 18px;
    }
  }
  
  .header-right {
    .user-info {
      display: flex;
      align-items: center;
      cursor: pointer;
      
      .username {
        margin: 0 8px;
        font-size: 14px;
      }
    }
  }
}

.main-content {
  background-color: #f5f5f5;
  padding: 20px;
  overflow-y: auto;
}
</style>
