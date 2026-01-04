<template>
  <div class="profile-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-info">
          <h1 class="page-title">个人中心</h1>
        </div>
      </div>
    </div>
    
    <!-- 用户信息 -->
    <div class="user-info">
      <div class="user-card">
        <div class="user-avatar">
          <div class="avatar-circle">
            {{ (userStore.userInfo.username || userStore.userInfo.name || 'U').charAt(0) }}
          </div>
          <div class="user-details">
            <div class="user-name">{{ userStore.userInfo.username || userStore.userInfo.name || '用户' }}</div>
            <div class="user-role">店长</div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 功能菜单 -->
    <div class="menu-section">
      <div class="menu-item" @click="handleMenuClick('store')">
        <div class="menu-icon">🏪</div>
        <div class="menu-content">
          <div class="menu-title">门店信息</div>
          <div class="menu-desc">查看和编辑门店信息</div>
        </div>
        <div class="menu-arrow">›</div>
      </div>
      
      <div class="menu-item" @click="handleMenuClick('password')">
        <div class="menu-icon">🔒</div>
        <div class="menu-content">
          <div class="menu-title">修改密码</div>
          <div class="menu-desc">更改登录密码</div>
        </div>
        <div class="menu-arrow">›</div>
      </div>
      
      <div class="menu-item" @click="handleMenuClick('about')">
        <div class="menu-icon">ℹ️</div>
        <div class="menu-content">
          <div class="menu-title">关于我们</div>
          <div class="menu-desc">应用版本和帮助信息</div>
        </div>
        <div class="menu-arrow">›</div>
      </div>
    </div>
    
    <!-- 退出登录 -->
    <div class="logout-section">
      <button class="logout-btn" @click="handleLogout">
        退出登录
      </button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { showToast, showConfirmDialog } from 'vant'

const router = useRouter()
const userStore = useUserStore()

// 处理菜单点击
const handleMenuClick = (type) => {
  switch (type) {
    case 'store':
      showToast('门店信息功能开发中...')
      break
    case 'password':
      showToast('修改密码功能开发中...')
      break
    case 'about':
      showToast('关于我们功能开发中...')
      break
  }
}

// 处理退出登录
const handleLogout = async () => {
  try {
    await showConfirmDialog({
      title: '确认退出',
      message: '确定要退出登录吗？'
    })
    
    await userStore.logoutAction()
    router.push('/login')
  } catch (error) {
    // 用户取消
  }
}
</script>

<style lang="scss" scoped>
.profile-page {
  background: linear-gradient(180deg, #f8f9ff 0%, #f0f2ff 100%);
  min-height: 100vh;
  padding-bottom: 20px;
}

// 页面头部
.page-header {
  background: linear-gradient(135deg, #ff6b6b 0%, #ffa726 100%);
  padding: 20px;
  color: white;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

// 用户信息
.user-info {
  padding: 20px 16px;
}

.user-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.user-avatar {
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar-circle {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 600;
  color: white;
}

.user-details {
  flex: 1;
}

.user-name {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.user-role {
  font-size: 14px;
  color: #666;
  background: #f0f2ff;
  padding: 4px 12px;
  border-radius: 12px;
  display: inline-block;
}

// 菜单区域
.menu-section {
  padding: 0 16px;
  margin-top: 16px;
}

.menu-item {
  background: white;
  border-radius: 12px;
  padding: 16px 20px;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  
  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  }
  
  &:active {
    transform: translateY(0);
  }
}

.menu-icon {
  font-size: 24px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9ff;
  border-radius: 10px;
}

.menu-content {
  flex: 1;
}

.menu-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 2px;
}

.menu-desc {
  font-size: 13px;
  color: #999;
}

.menu-arrow {
  font-size: 20px;
  color: #ccc;
  font-weight: 300;
}

// 退出登录
.logout-section {
  padding: 20px 16px;
  margin-top: 32px;
}

.logout-btn {
  width: 100%;
  height: 48px;
  background: linear-gradient(135deg, #ff4d4f, #ff7875);
  color: white;
  border: none;
  border-radius: 24px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(255, 77, 79, 0.3);
  }
  
  &:active {
    transform: translateY(0);
  }
}

// 响应式设计
@media (max-width: 375px) {
  .user-card {
    padding: 20px;
  }
  
  .avatar-circle {
    width: 50px;
    height: 50px;
    font-size: 20px;
  }
  
  .user-name {
    font-size: 18px;
  }
  
  .menu-item {
    padding: 14px 16px;
  }
  
  .menu-icon {
    font-size: 20px;
    width: 36px;
    height: 36px;
  }
}
</style>