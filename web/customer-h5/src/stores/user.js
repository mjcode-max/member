import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref({})
  const memberInfo = ref({})

  const isLoggedIn = computed(() => !!userInfo.value.openid || !!userInfo.value.phone)
  const isMember = computed(() => !!memberInfo.value.id)

  // 设置用户信息
  const setUserInfo = (info) => {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  // 设置会员信息
  const setMemberInfo = (info) => {
    memberInfo.value = info
    localStorage.setItem('memberInfo', JSON.stringify(info))
  }

  // 清除用户信息
  const clearUserInfo = () => {
    userInfo.value = {}
    memberInfo.value = {}
    localStorage.removeItem('userInfo')
    localStorage.removeItem('memberInfo')
  }

  // 初始化用户信息
  const initUser = () => {
    const savedUserInfo = localStorage.getItem('userInfo')
    const savedMemberInfo = localStorage.getItem('memberInfo')
    
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
    userInfo,
    memberInfo,
    isLoggedIn,
    isMember,
    setUserInfo,
    setMemberInfo,
    clearUserInfo,
    initUser,
    getUserPhone
  }
})
