import request from './request'

// 微信登录：通过code换取openid并保存
export const wechatLoginByCode = (code) => {
  return request({
    url: '/public/customer/wechat/login',
    method: 'post',
    data: { code }
  })
}
