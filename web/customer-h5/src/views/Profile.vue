<template>
  <div class="profile-page">
    <!-- 用户信息 -->
    <div class="user-info">
      <div class="user-card">
        <div class="user-avatar">
          <div class="avatar" :style="{ backgroundImage: userStore.userInfo.avatar ? `url(${userStore.userInfo.avatar})` : 'none' }">
            {{ userStore.userInfo.name?.charAt(0) || '用' }}
          </div>
          <div class="user-details">
            <div class="user-name">{{ userStore.userInfo.nickname || userStore.userInfo.name || '未登录' }}</div>
            <div class="user-phone">{{ userStore.userInfo.phone || '微信用户' }}</div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 功能菜单 -->
    <div class="menu-section">
      <div class="menu-item" @click="$router.push('/booking')">
        <div class="menu-icon">📅</div>
        <div class="menu-title">我的预约</div>
        <div class="menu-arrow">></div>
      </div>
      
      <div class="menu-item" @click="$router.push('/member')">
        <div class="menu-icon">💎</div>
        <div class="menu-title">会员中心</div>
        <div class="menu-arrow">></div>
      </div>
      
      <div class="menu-item" @click="handleMenuClick('stores')">
        <div class="menu-icon">🏪</div>
        <div class="menu-title">门店查询</div>
        <div class="menu-arrow">></div>
      </div>
    </div>
    
    <!-- 设置菜单 -->
    <div class="menu-section">
      <div class="menu-item" @click="handleMenuClick('about')">
        <div class="menu-icon">ℹ️</div>
        <div class="menu-title">关于我们</div>
        <div class="menu-arrow">></div>
      </div>
      
      <div class="menu-item" @click="handleMenuClick('contact')">
        <div class="menu-icon">📞</div>
        <div class="menu-title">客服电话</div>
        <div class="menu-arrow">></div>
      </div>
    </div>
    
    <!-- 登录/退出按钮 -->
    <div class="action-section">
      <button
        v-if="!userStore.isLoggedIn"
        class="action-btn login-btn"
        @click="handleLogin"
      >
        登录
      </button>
      
      <button
        v-else
        class="action-btn logout-btn"
        @click="handleLogout"
      >
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
    case 'stores':
      showToast('门店查询功能开发中...')
      break
    case 'about':
      showToast('关于我们功能开发中...')
      break
    case 'contact':
      showToast('客服电话：400-123-4567')
      break
  }
}

// 处理登录
const handleLogin = () => {
  router.push('/login')
}

// 处理退出登录
const handleLogout = async () => {
  try {
    await showConfirmDialog({
      title: '确认退出',
      message: '确定要退出登录吗？'
    })
    
    userStore.clearUserInfo()
    showToast.success('已退出登录')
  } catch (error) {
    // 用户取消
  }
}
</script>

<style lang="scss" scoped>
.profile-page {
  background-color: #ffffff;
  min-height: 100vh;
}

.user-info {
  padding: 16px;
  
  .user-card {
    background: white;
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  .user-avatar {
    display: flex;
    align-items: center;
    gap: 16px;
    
    .avatar {
      width: 60px;
      height: 60px;
      border-radius: 50%;
      background: linear-gradient(135deg, #667eea, #764ba2);
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 24px;
      font-weight: 600;
      background-size: cover;
      background-position: center;
    }
    
    .user-details {
      .user-name {
        font-size: 18px;
        font-weight: 600;
        color: #333;
        margin-bottom: 4px;
      }
      
      .user-phone {
        font-size: 14px;
        color: #666;
      }
    }
  }
}

.menu-section {
  margin: 16px;
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  
  .menu-item {
    display: flex;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid #f0f0f0;
    cursor: pointer;
    transition: background-color 0.2s;
    
    &:last-child {
      border-bottom: none;
    }
    
    &:hover {
      background-color: #f8f9fa;
    }
    
    .menu-icon {
      font-size: 20px;
      margin-right: 12px;
      width: 24px;
      text-align: center;
    }
    
    .menu-title {
      flex: 1;
      font-size: 16px;
      color: #333;
    }
    
    .menu-arrow {
      color: #999;
      font-size: 16px;
    }
  }
}

.action-section {
  padding: 16px;
  margin-top: 32px;
  
  .action-btn {
    width: 100%;
    padding: 14px;
    border: none;
    border-radius: 8px;
    font-size: 16px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    
    &.login-btn {
      background: linear-gradient(135deg, #667eea, #764ba2);
      color: white;
      
      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
      }
    }
    
    &.logout-btn {
      background: #ff4d4f;
      color: white;
      
      &:hover {
        background: #ff7875;
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(255, 77, 79, 0.3);
      }
    }
  }
}
</style>
