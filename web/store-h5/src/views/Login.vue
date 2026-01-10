<template>
  <div class="login-page">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="decoration-shape shape-1"></div>
      <div class="decoration-shape shape-2"></div>
      <div class="decoration-shape shape-3"></div>
    </div>
    
    <!-- 登录容器 -->
    <div class="login-container">
      <!-- 头部logo区域 -->
      <div class="login-header">
        <div class="logo-container">
          <div class="logo-icon">
            <i class="logo-emoji">🏪</i>
          </div>
          <div class="logo-glow"></div>
        </div>
        <h1 class="app-title">店长管理台</h1>
        <p class="app-subtitle">门店管理，一手掌控</p>
      </div>
      
      <!-- 登录表单 -->
      <div class="login-form-container">
        <div class="form-card">
          <div class="form-header">
            <h3>管理员登录</h3>
            <p>请使用您的管理员账户登录</p>
          </div>
          
          <van-form @submit="handleLogin" class="login-form">
            <div class="form-group">
              <div class="input-container">
                <i class="input-icon">👨‍💼</i>
                <input
                  v-model="loginForm.username"
                  type="text"
                  placeholder="请输入用户名"
                  class="form-input"
                  required
                />
              </div>
            </div>
            
            <div class="form-group">
              <div class="input-container">
                <i class="input-icon">🔐</i>
                <input
                  v-model="loginForm.password"
                  type="password"
                  placeholder="请输入密码"
                  class="form-input"
                  required
                />
              </div>
            </div>
            
            <div class="form-group">
              <div class="input-container" @click="showUserTypePicker = true">
                <i class="input-icon">🎯</i>
                <input
                  v-model="userTypeText"
                  type="text"
                  placeholder="请选择用户类型"
                  class="form-input"
                  readonly
                />
                <i class="arrow-icon">›</i>
              </div>
            </div>
            
            <div class="login-actions">
              <button
                type="submit"
                class="login-button"
                :disabled="loading"
              >
                <span v-if="loading" class="loading-text">
                  <i class="loading-icon">⏳</i>
                  登录中...
                </span>
                <span v-else>
                  <i class="login-icon">🔑</i>
                  管理员登录
                </span>
              </button>
            </div>
          </van-form>
          
          <!-- 底部信息 -->
          <div class="form-footer">
            <div class="help-links">
              <span @click="handleForgotPassword">忘记密码？</span>
              <span @click="handleContact">联系技术支持</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 用户类型选择器 -->
    <van-popup 
      v-model:show="showUserTypePicker" 
      position="bottom" 
      round
      :style="{ height: '40%' }"
    >
      <div class="popup-header">
        <div class="popup-title">选择用户类型</div>
        <div class="popup-close" @click="showUserTypePicker = false">×</div>
      </div>
      <van-picker
        :columns="userTypeOptions"
        @confirm="onUserTypeConfirm"
        @cancel="showUserTypePicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const showUserTypePicker = ref(false)

const loginForm = reactive({
  username: 'liukang',
  password: '123456',
  user_type: 'store_manager'
})

const userTypeOptions = [
  { text: '店长', value: 'store_manager' }
]

const userTypeText = computed(() => {
  const option = userTypeOptions.find(opt => opt.value === loginForm.user_type)
  return option ? option.text : ''
})

const onUserTypeConfirm = ({ selectedOptions }) => {
  loginForm.user_type = selectedOptions[0].value
  showUserTypePicker.value = false
}

// 忘记密码
const handleForgotPassword = () => {
  showToast('请联系技术支持重置密码')
}

// 联系技术支持
const handleContact = () => {
  showToast('技术支持：400-123-4567')
}

const handleLogin = async () => {
  loading.value = true
  try {
    const success = await userStore.loginAction(loginForm)
    if (success) {
      router.push('/')
    }
  } catch (error) {
    console.error('登录失败:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #ff6b6b 0%, #ffa726 100%);
  position: relative;
  overflow: hidden;
}

// 背景装饰（店长端使用不同的装饰）
.bg-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 0;
}

.decoration-shape {
  position: absolute;
  background: rgba(255, 255, 255, 0.1);
  
  &.shape-1 {
    width: 120px;
    height: 120px;
    border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
    top: 10%;
    right: -60px;
    animation: morph 8s ease-in-out infinite;
  }
  
  &.shape-2 {
    width: 160px;
    height: 160px;
    border-radius: 60% 40% 30% 70% / 60% 30% 70% 40%;
    bottom: 20%;
    left: -80px;
    animation: morph 10s ease-in-out infinite reverse;
  }
  
  &.shape-3 {
    width: 80px;
    height: 80px;
    border-radius: 40% 60% 60% 40% / 60% 40% 40% 60%;
    top: 60%;
    right: 10%;
    animation: morph 6s ease-in-out infinite;
  }
}

// 登录容器
.login-container {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 100vh;
  padding: 40px 20px;
}

// 头部区域
.login-header {
  text-align: center;
  color: white;
  margin-bottom: 40px;
}

.logo-container {
  position: relative;
  display: inline-block;
  margin-bottom: 24px;
}

.logo-icon {
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.3);
  position: relative;
  z-index: 2;
}

.logo-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100px;
  height: 100px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.3) 0%, transparent 70%);
  border-radius: 50%;
  animation: glow 3s ease-in-out infinite;
}

.app-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 8px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

.app-subtitle {
  font-size: 16px;
  opacity: 0.9;
  margin: 0;
}

// 表单容器
.login-form-container {
  width: 100%;
  max-width: 400px;
  margin: 0 auto;
}

.form-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 32px 24px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.form-header {
  text-align: center;
  margin-bottom: 32px;
  
  h3 {
    font-size: 22px;
    font-weight: 700;
    color: #333;
    margin-bottom: 8px;
  }
  
  p {
    font-size: 14px;
    color: #666;
    margin: 0;
  }
}

// 表单样式
.form-group {
  margin-bottom: 20px;
}

.input-container {
  position: relative;
  display: flex;
  align-items: center;
  background: white;
  border-radius: 16px;
  padding: 0 16px;
  border: 2px solid #f0f0f0;
  transition: all 0.3s ease;
  
  &:focus-within {
    border-color: #ff6b6b;
    box-shadow: 0 0 0 4px rgba(255, 107, 107, 0.1);
  }
}

.input-icon {
  font-size: 18px;
  margin-right: 12px;
  color: #999;
}

.form-input {
  flex: 1;
  border: none;
  outline: none;
  padding: 16px 0;
  font-size: 16px;
  background: transparent;
  color: #333;
  
  &::placeholder {
    color: #999;
  }
  
  &[readonly] {
    cursor: pointer;
  }
}

.arrow-icon {
  font-size: 16px;
  color: #ccc;
  margin-left: 8px;
}

// 登录按钮
.login-actions {
  margin-top: 32px;
}

.login-button {
  width: 100%;
  height: 52px;
  background: linear-gradient(135deg, #ff6b6b, #ffa726);
  border: none;
  border-radius: 26px;
  color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 8px 24px rgba(255, 107, 107, 0.3);
  
  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 12px 32px rgba(255, 107, 107, 0.4);
  }
  
  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }
}

.loading-text,
.login-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.loading-icon {
  animation: spin 1s linear infinite;
}

// 表单底部
.form-footer {
  margin-top: 24px;
  text-align: center;
}

.help-links {
  display: flex;
  justify-content: space-between;
  
  span {
    font-size: 14px;
    color: #ff6b6b;
    cursor: pointer;
    transition: color 0.2s ease;
    
    &:hover {
      color: #ffa726;
    }
  }
}

// 弹窗样式
.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.popup-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.popup-close {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: #e9ecef;
  }
}

// 动画效果
@keyframes morph {
  0%, 100% {
    border-radius: 30% 70% 70% 30% / 30% 30% 70% 70%;
    transform: rotate(0deg);
  }
  25% {
    border-radius: 58% 42% 75% 25% / 76% 46% 54% 24%;
    transform: rotate(90deg);
  }
  50% {
    border-radius: 50% 50% 33% 67% / 55% 27% 73% 45%;
    transform: rotate(180deg);
  }
  75% {
    border-radius: 33% 67% 58% 42% / 63% 68% 32% 37%;
    transform: rotate(270deg);
  }
}

@keyframes glow {
  0%, 100% { opacity: 0.5; transform: translate(-50%, -50%) scale(1); }
  50% { opacity: 0.8; transform: translate(-50%, -50%) scale(1.2); }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.form-card {
  animation: slideInUp 0.6s ease-out;
}

.login-header {
  animation: slideInUp 0.6s ease-out 0.2s both;
}

// 响应式设计
@media (max-width: 375px) {
  .login-container {
    padding: 20px 16px;
  }
  
  .app-title {
    font-size: 24px;
  }
  
  .form-card {
    padding: 24px 20px;
  }
  
  .logo-icon {
    width: 60px;
    height: 60px;
    font-size: 28px;
  }
  
  .logo-glow {
    width: 80px;
    height: 80px;
  }
}
</style>
