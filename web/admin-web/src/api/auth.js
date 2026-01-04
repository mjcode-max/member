import request from './request'

// 用户登录
// 请求参数: { username: string, password?: string }
// username: 用户名（员工）或手机号（顾客）
// password: 密码（员工必填，顾客不需要）
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
