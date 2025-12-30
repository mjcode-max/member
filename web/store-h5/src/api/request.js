import axios from 'axios'
import { showSuccessToast, showFailToast } from 'vant'
import { useUserStore } from '@/stores/user'

// 创建axios实例
const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    
    // 添加认证token
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    
    return config
  },
  (error) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const { data } = response
    
    // 如果是文件下载等特殊响应，直接返回
    if (response.config.responseType === 'blob') {
      return response
    }
    
    // 统一处理响应数据
    if (data.code === 200) {
      return data
    } else if (data.code === 401) {
      // token过期，清除用户信息并跳转到登录页
      const userStore = useUserStore()
      userStore.logoutAction()
      window.location.href = '/login'
      return Promise.reject(new Error(data.message || '认证失败'))
    } else {
      showFailToast(data.message || '请求失败')
      return Promise.reject(new Error(data.message || '请求失败'))
    }
  },
  (error) => {
    console.error('响应错误:', error)
    
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          showFailToast('认证失败，请重新登录')
          const userStore = useUserStore()
          userStore.logoutAction()
          window.location.href = '/login'
          break
        case 403:
          showFailToast('权限不足')
          break
        case 404:
          showFailToast('请求的资源不存在')
          break
        case 500:
          showFailToast('服务器内部错误')
          break
        default:
          showFailToast(data?.message || '请求失败')
      }
    } else if (error.code === 'ECONNABORTED') {
      showFailToast('请求超时')
    } else {
      showFailToast('网络错误')
    }
    
    return Promise.reject(error)
  }
)

export default request
