<template>
  <div class="login-page">
    <div class="login-container">
      <!-- Logo区域 -->
      <div class="logo-section">
        <div class="logo">
          <div class="logo-icon">💅</div>
          <div class="logo-text">美甲美睫预约</div>
        </div>
        <div class="welcome-text">欢迎使用美甲美睫预约系统</div>
      </div>

      <!-- 登录按钮区域 -->
      <div class="login-section">
        <button 
          class="wechat-login-btn"
          @click="handleWechatLogin"
          :disabled="loading"
        >
          <div class="btn-content">
            <div class="btn-icon">📱</div>
            <div class="btn-text">
              {{ loading ? '登录中...' : '微信一键登录' }}
            </div>
          </div>
        </button>
        
        <div class="login-tips">
          <div class="tip-item">• 使用微信授权登录，安全便捷</div>
          <div class="tip-item">• 自动获取微信头像和昵称</div>
          <div class="tip-item">• 支持快速预约和会员服务</div>
        </div>
      </div>

      <!-- 底部信息 -->
      <div class="footer-section">
        <div class="footer-text">
          登录即表示同意
          <span class="link-text" @click="showPrivacyPolicy">《隐私政策》</span>
          和
          <span class="link-text" @click="showUserAgreement">《用户协议》</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { wechatLogin, getWechatConfig } from '@/api/auth'
import { showToast, showDialog } from 'vant'
import { isWechatBrowser, getWechatCode, initWechatSDK } from '@/utils/wechat'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)

// 微信登录
const handleWechatLogin = async () => {
  if (loading.value) return
  
  loading.value = true
  
  try {
    // 检查是否在微信环境中
    if (!isWechatBrowser()) {
      showToast('请在微信中打开此页面')
      return
    }

    // 获取微信授权码
    const code = await getWechatCode()
    if (!code) {
      showToast('获取微信授权失败')
      return
    }

    // 调用后端登录接口
    const response = await wechatLogin(code)
    if (response.code === 200) {
      // 保存用户信息
      userStore.setUserInfo(response.data.user)
      
      showToast.success('登录成功')
      
      // 跳转到首页
      router.push('/')
    } else {
      showToast(response.message || '登录失败')
    }
  } catch (error) {
    console.error('微信登录失败:', error)
    showToast('登录失败，请重试')
  } finally {
    loading.value = false
  }
}


// 初始化微信JS-SDK
const initWechatSDKConfig = async () => {
  try {
    const url = window.location.href.split('#')[0]
    const response = await getWechatConfig(url)
    
    if (response.code === 200) {
      const config = response.data
      await initWechatSDK(config)
    }
  } catch (error) {
    console.error('获取微信配置失败:', error)
  }
}

// 显示隐私政策
const showPrivacyPolicy = () => {
  showDialog({
    title: '隐私政策',
    message: '我们重视您的隐私保护，详细内容请查看完整版隐私政策。',
    confirmButtonText: '我知道了'
  })
}

// 显示用户协议
const showUserAgreement = () => {
  showDialog({
    title: '用户协议',
    message: '使用本服务即表示您同意遵守相关用户协议条款。',
    confirmButtonText: '我知道了'
  })
}

onMounted(() => {
  // 初始化微信JS-SDK
  initWechatSDKConfig()
})
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 20px;
  padding: 40px 30px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.logo-section {
  text-align: center;
  margin-bottom: 40px;
  
  .logo {
    margin-bottom: 20px;
    
    .logo-icon {
      font-size: 60px;
      margin-bottom: 10px;
    }
    
    .logo-text {
      font-size: 24px;
      font-weight: bold;
      color: #333;
    }
  }
  
  .welcome-text {
    font-size: 16px;
    color: #666;
    line-height: 1.5;
  }
}

.login-section {
  margin-bottom: 40px;
  
  .wechat-login-btn {
    width: 100%;
    height: 60px;
    background: linear-gradient(135deg, #07c160 0%, #00d4aa 100%);
    border: none;
    border-radius: 30px;
    color: white;
    font-size: 18px;
    font-weight: bold;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 8px 20px rgba(7, 193, 96, 0.3);
    
    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 12px 25px rgba(7, 193, 96, 0.4);
    }
    
    &:active {
      transform: translateY(0);
    }
    
    &:disabled {
      opacity: 0.7;
      cursor: not-allowed;
      transform: none;
    }
    
    .btn-content {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 10px;
      
      .btn-icon {
        font-size: 24px;
      }
      
      .btn-text {
        font-size: 18px;
      }
    }
  }
  
  .login-tips {
    margin-top: 30px;
    
    .tip-item {
      font-size: 14px;
      color: #666;
      margin-bottom: 8px;
      line-height: 1.4;
    }
  }
}

.footer-section {
  text-align: center;
  
  .footer-text {
    font-size: 12px;
    color: #999;
    line-height: 1.5;
    
    .link-text {
      color: #667eea;
      cursor: pointer;
      
      &:hover {
        text-decoration: underline;
      }
    }
  }
}
</style>
