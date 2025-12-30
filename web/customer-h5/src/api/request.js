import axios from 'axios'
import { showToast } from 'vant'

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
    } else {
      showToast.fail(data.message || '请求失败')
      return Promise.reject(new Error(data.message || '请求失败'))
    }
  },
  (error) => {
    console.error('响应错误:', error)
    
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 400:
          showToast.fail(data?.message || '请求参数错误')
          break
        case 404:
          showToast.fail('请求的资源不存在')
          break
        case 500:
          showToast.fail('服务器内部错误')
          break
        default:
          showToast.fail(data?.message || '请求失败')
      }
    } else if (error.code === 'ECONNABORTED') {
      showToast.fail('请求超时')
    } else {
      showToast.fail('网络错误')
    }
    
    return Promise.reject(error)
  }
)

export default request
