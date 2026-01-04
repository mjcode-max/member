import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))
  const memberInfo = ref(JSON.parse(localStorage.getItem('memberInfo') || '{}'))

  const isLoggedIn = computed(() => !!token.value || !!userInfo.value.phone)
  const isMember = computed(() => !!memberInfo.value.id)

  // 设置token
  const setToken = (newToken) => {
    token.value = newToken
    if (newToken) {
      localStorage.setItem('token', newToken)
    } else {
      localStorage.removeItem('token')
    }
  }

  // 设置用户信息
  const setUserInfo = (info) => {
    userInfo.value = info
    if (info && Object.keys(info).length > 0) {
      localStorage.setItem('userInfo', JSON.stringify(info))
    } else {
      localStorage.removeItem('userInfo')
    }
  }

  // 设置会员信息
  const setMemberInfo = (info) => {
    memberInfo.value = info
    if (info && Object.keys(info).length > 0) {
      localStorage.setItem('memberInfo', JSON.stringify(info))
    } else {
      localStorage.removeItem('memberInfo')
    }
  }

  // 清除用户信息
  const clearUserInfo = () => {
    token.value = ''
    userInfo.value = {}
    memberInfo.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    localStorage.removeItem('memberInfo')
  }

  // 初始化用户信息
  const initUser = () => {
    const savedToken = localStorage.getItem('token')
    const savedUserInfo = localStorage.getItem('userInfo')
    const savedMemberInfo = localStorage.getItem('memberInfo')
    
    if (savedToken) {
      token.value = savedToken
    }
    
    if (savedUserInfo) {
      try {
        userInfo.value = JSON.parse(savedUserInfo)
      } catch (error) {
        console.error('解析用户信息失败:', error)
        localStorage.removeItem('userInfo')
        userInfo.value = {}
      }
    }
    
    if (savedMemberInfo) {
      try {
        memberInfo.value = JSON.parse(savedMemberInfo)
      } catch (error) {
        console.error('解析会员信息失败:', error)
        localStorage.removeItem('memberInfo')
        memberInfo.value = {}
      }
    }
  }

  // 获取用户手机号
  const getUserPhone = () => {
    return userInfo.value.phone || ''
  }

  return {
    token,
    userInfo,
    memberInfo,
    isLoggedIn,
    isMember,
    setToken,
    setUserInfo,
    setMemberInfo,
    clearUserInfo,
    initUser,
    getUserPhone
  }
})
