import request from './request'

// 用户登录（手机号登录，顾客）
// 请求参数: { username: string }
// username: 手机号（顾客登录不需要密码）
export const login = (data) => {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

// 获取当前用户信息
export const getCurrentUser = () => {
  return request({
    url: '/auth/me',
    method: 'get'
  })
}

// 用户登出
export const logout = () => {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
}

// 微信登录（如果后端支持）
export const wechatLogin = (code) => {
  return request({
    url: '/customer/auth/wechat/login',
    method: 'post',
    data: { code }
  })
}

// 获取微信JS-SDK配置（如果后端支持）
export const getWechatConfig = (url) => {
  return request({
    url: '/customer/auth/wechat/config',
    method: 'get',
    params: { url }
  })
}
