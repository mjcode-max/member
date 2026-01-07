import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, logout } from '@/api/auth'
import { ElMessage } from 'element-plus'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

  const isLoggedIn = computed(() => !!token.value)

  // 登录
  const loginAction = async (loginForm) => {
    try {
      const response = await login(loginForm)
      if (response.code === 200) {
        // 验证用户角色，后台只允许管理员登录
        if (response.data.user && response.data.user.role !== 'admin') {
          ElMessage.error('只有管理员可以登录后台系统')
          return false
        }
        
        token.value = response.data.token
        userInfo.value = response.data.user
        
        // 保存到本地存储
        localStorage.setItem('token', response.data.token)
        localStorage.setItem('userInfo', JSON.stringify(response.data.user))
        
        ElMessage.success('登录成功')
        return true
      } else {
        ElMessage.error(response.message || '登录失败')
        return false
      }
    } catch (error) {
      ElMessage.error(error.message || '登录失败，请检查网络连接')
      return false
    }
  }

  // 登出
  const logoutAction = async () => {
    try {
      await logout()
    } catch (error) {
      console.error('登出失败:', error)
    } finally {
      // 清除本地存储
      token.value = ''
      userInfo.value = {}
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    }
  }

  // 初始化用户信息
  const initUser = () => {
    const savedToken = localStorage.getItem('token')
    const savedUserInfo = localStorage.getItem('userInfo')
    
    if (savedToken) {
      token.value = savedToken
    }
    
    if (savedUserInfo) {
      try {
        userInfo.value = JSON.parse(savedUserInfo)
      } catch (error) {
        console.error('解析用户信息失败:', error)
        localStorage.removeItem('userInfo')
      }
    }
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    loginAction,
    logoutAction,
    initUser
  }
})
