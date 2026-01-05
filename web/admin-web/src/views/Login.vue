<template>
  <div class="login-container">
    <div class="login-form">
      <div class="login-header">
        <h2>美甲美睫管理系统</h2>
        <p>总公司后台管理</p>
      </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form-content"
        @submit.prevent="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            size="large"
            prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item prop="user_type">
          <el-select
            v-model="loginForm.user_type"
            placeholder="请选择用户类型"
            size="large"
            style="width: 100%"
          >
            <el-option label="管理员" value="admin" />
            <el-option label="店长" value="store_manager" />
            <el-option label="美甲师" value="staff" />
          </el-select>
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            style="width: 100%"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref()
const loading = ref(false)

const loginForm = reactive({
  username: 'admin',
  password: 'password',
  user_type: 'admin'
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  user_type: [
    { required: true, message: '请选择用户类型', trigger: 'change' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    const success = await userStore.loginAction(loginForm)
    debugger
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
.login-container {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-form {
  width: 400px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  padding: 40px;
  
  .login-header {
    text-align: center;
    margin-bottom: 30px;
    
    h2 {
      color: #333;
      margin-bottom: 8px;
      font-size: 24px;
    }
    
    p {
      color: #666;
      font-size: 14px;
    }
  }
  
  .login-form-content {
    .el-form-item {
      margin-bottom: 20px;
    }
  }
}
</style>
