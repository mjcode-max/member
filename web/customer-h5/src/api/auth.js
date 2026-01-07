import request from './request'

// 微信登录：通过code换取openid和手机号并保存
export const wechatLoginByCode = (code, phone = '') => {
  return request({
    url: '/public/customer/wechat/login',
    method: 'post',
    data: { code, phone }
  })
}
