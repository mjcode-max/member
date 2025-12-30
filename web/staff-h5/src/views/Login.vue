<template>
  <div class="login-page">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="decoration-circle circle-1"></div>
      <div class="decoration-circle circle-2"></div>
      <div class="decoration-circle circle-3"></div>
    </div>
    
    <!-- 登录容器 -->
    <div class="login-container">
      <!-- 头部logo区域 -->
      <div class="login-header">
        <div class="logo-container">
          <div class="logo-icon">
            <i class="logo-emoji">💅</i>
          </div>
          <div class="logo-rings">
            <div class="ring ring-1"></div>
            <div class="ring ring-2"></div>
          </div>
        </div>
        <h1 class="app-title">美甲师工作台</h1>
        <p class="app-subtitle">专业服务，精致体验</p>
      </div>
      
      <!-- 登录表单 -->
      <div class="login-form-container">
        <div class="form-card">
          <div class="form-header">
            <h3>欢迎回来</h3>
            <p>请登录您的账户开始工作</p>
          </div>
          
          <van-form @submit="handleLogin" class="login-form">
            <div class="form-group">
              <div class="input-container">
                <i class="input-icon">👤</i>
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
                <i class="input-icon">🔒</i>
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
                <i class="input-icon">👔</i>
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
                  <i class="login-icon">🚀</i>
                  立即登录
                </span>
              </button>
            </div>
          </van-form>
          
          <!-- 底部信息 -->
          <div class="form-footer">
            <div class="help-links">
              <span @click="handleForgotPassword">忘记密码？</span>
              <span @click="handleContact">联系管理员</span>
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
  username: 'staff_1_13800138000',
  password: '123456',
  user_type: 'staff'
})

const userTypeOptions = [
  { text: '美甲师', value: 'staff' }
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
  showToast('请联系管理员重置密码')
}

// 联系管理员
const handleContact = () => {
  showToast('管理员电话：400-123-4567')
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

// 背景装饰
.bg-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 0;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  
  &.circle-1 {
    width: 200px;
    height: 200px;
    top: -100px;
    right: -100px;
    animation: float 6s ease-in-out infinite;
  }
  
  &.circle-2 {
    width: 150px;
    height: 150px;
    bottom: -75px;
    left: -75px;
    animation: float 8s ease-in-out infinite reverse;
  }
  
  &.circle-3 {
    width: 100px;
    height: 100px;
    top: 50%;
    left: -50px;
    animation: float 10s ease-in-out infinite;
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

.logo-rings {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.ring {
  position: absolute;
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  
  &.ring-1 {
    width: 100px;
    height: 100px;
    animation: pulse 2s ease-in-out infinite;
  }
  
  &.ring-2 {
    width: 120px;
    height: 120px;
    animation: pulse 2s ease-in-out infinite 0.5s;
  }
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
    border-color: #667eea;
    box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
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
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  border-radius: 26px;
  color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
  
  &:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 12px 32px rgba(102, 126, 234, 0.4);
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
    color: #667eea;
    cursor: pointer;
    transition: color 0.2s ease;
    
    &:hover {
      color: #764ba2;
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
@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-20px); }
}

@keyframes pulse {
  0%, 100% { opacity: 0.3; transform: translate(-50%, -50%) scale(1); }
  50% { opacity: 0.1; transform: translate(-50%, -50%) scale(1.1); }
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
  
  .ring {
    &.ring-1 {
      width: 80px;
      height: 80px;
    }
    
    &.ring-2 {
      width: 100px;
      height: 100px;
    }
  }
}
</style>
