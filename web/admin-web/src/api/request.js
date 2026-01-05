import axios from 'axios'
// import { ElMessage } from 'element-plus'
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
    
    // 后台系统标识
    config.headers['X-Client-Type'] = 'admin-web'
    
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
    
    // 统一处理响应数据 - 后端返回格式: { code, message, data }
    // 注意：分页接口返回 code: 0，其他接口返回 code: 200
    // 成功状态码：0, 200, 201
    if (data.code === 200 || data.code === 201 || data.code === 0) {
      return data
    } else if (data.code === 401 || data.code === 403) {
      // token过期或权限不足，清除用户信息并跳转到登录页
      const userStore = useUserStore()
      userStore.logoutAction()
      window.location.href = '/login'
      return Promise.reject(new Error(data.message || '认证失败'))
    } else {
      // 其他错误，不显示 Message，只记录日志
      console.error('请求失败:', data.message || '请求失败')
      return Promise.reject(new Error(data.message || '请求失败'))
    }
  },
  (error) => {
    console.error('响应错误:', error)
    
    if (error.response) {
      const { status, data } = error.response
      
      // 如果响应体中有 message，优先使用
      const errorMessage = data?.message || data?.data?.message
      
      switch (status) {
        case 401:
          // 认证失败，清除用户信息并跳转到登录页，不显示 Message
          console.error('认证失败:', errorMessage || '认证失败，请重新登录')
          const userStore = useUserStore()
          userStore.logoutAction()
          window.location.href = '/login'
          break
        case 403:
          console.error('权限不足:', errorMessage || '权限不足')
          break
        case 404:
          console.error('请求的资源不存在:', errorMessage || '请求的资源不存在')
          break
        case 500:
          console.error('服务器内部错误:', errorMessage || '服务器内部错误')
          break
        default:
          console.error('请求失败:', errorMessage || '请求失败')
      }
    } else if (error.code === 'ECONNABORTED') {
      console.error('请求超时')
    } else {
      console.error('网络错误')
    }
    
    return Promise.reject(error)
  }
)

export default request
