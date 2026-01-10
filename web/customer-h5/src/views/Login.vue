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
import { wechatLoginByCode } from '@/api/auth'
import { showToast, showDialog, showInputDialog } from 'vant'
import { isWechatBrowser, redirectToWechatAuth, getPhoneNumberBySDK } from '@/utils/wechat'

const router = useRouter()
const loading = ref(false)
const phoneInput = ref('')

// 微信登录完整流程
const handleWechatLogin = async () => {
  if (loading.value) return
  
  // 检查是否在微信环境中
  if (!isWechatBrowser()) {
    showToast('请在微信中打开此页面')
    return
  }

  // 检查URL中是否有code（微信授权回调）
  const urlParams = new URLSearchParams(window.location.search)
  const code = urlParams.get('code')
  
  if (code) {
    // 有code，尝试获取手机号
    let phone = ''
    
    // 尝试通过JS-SDK获取手机号
    try {
      const phoneResult = await getPhoneNumberBySDK()
      phone = phoneResult.phone || ''
    } catch (error) {
      console.warn('通过JS-SDK获取手机号失败:', error)
    }
    
    // 如果无法获取手机号，提示用户输入
    if (!phone) {
      try {
        const result = await showInputDialog({
          title: '请输入手机号',
          message: '为了更好的服务体验，请输入您的手机号',
          placeholder: '请输入11位手机号',
          validator: (value) => {
            const phoneRegex = /^1[3-9]\d{9}$/
            if (!value) {
              return '请输入手机号'
            }
            if (!phoneRegex.test(value)) {
              return '请输入正确的手机号'
            }
            return true
          }
        })
        phone = result.value
      } catch (error) {
        // 用户取消输入，仍然可以登录（手机号为空）
        console.log('用户取消输入手机号')
      }
    }
    
    // 调用后端接口通过code换取openid和手机号并保存
    loading.value = true
    showToast.loading({
      message: '正在登录...',
      forbidClick: true,
      duration: 0
    })
    
    try {
      const response = await wechatLoginByCode(code, phone)
      if (response.code === 200) {
        showToast.success('登录成功')
        // 清除URL参数，跳转到首页
        window.history.replaceState({}, '', window.location.pathname)
        router.push('/')
      } else {
        showToast(response.message || '登录失败')
      }
    } catch (error) {
      console.error('登录失败:', error)
      showToast('登录失败，请重试')
    } finally {
      loading.value = false
    }
    return
  }

  // 没有code，跳转到微信授权页面
  try {
    redirectToWechatAuth()
  } catch (error) {
    console.error('跳转微信授权失败:', error)
    // 如果是开发环境的内网地址错误，显示更详细的提示
    if (error.message && error.message.includes('内网地址')) {
      showDialog({
        title: '开发环境配置提示',
        message: error.message + '\n\n详细说明请查看控制台',
        confirmButtonText: '我知道了'
      })
    } else {
      showToast(error.message || '跳转微信授权失败')
    }
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

// 页面加载时，如果URL中有code，自动触发登录
onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search)
  const code = urlParams.get('code')
  if (code && isWechatBrowser()) {
    handleWechatLogin()
  }
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
