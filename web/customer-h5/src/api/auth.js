import request from './request'

// 微信登录
export const wechatLogin = (code) => {
  return request({
    url: '/customer/auth/wechat/login',
    method: 'post',
    data: { code }
  })
}

// 获取微信JS-SDK配置
export const getWechatConfig = (url) => {
  return request({
    url: '/customer/auth/wechat/config',
    method: 'get',
    params: { url }
  })
}
