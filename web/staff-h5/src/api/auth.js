import request from './request'

// 用户登录
export const login = (data) => {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

// 用户登出
export const logout = () => {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
}

// 刷新token
export const refreshToken = () => {
  return request({
    url: '/auth/refresh',
    method: 'post'
  })
}
