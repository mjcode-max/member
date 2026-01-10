import request from './request'

/**
 * 支付相关API
 */

// 创建支付（通用）
export const createPayment = (data) => {
  return request({
    url: '/payments',
    method: 'post',
    data
  })
}

// 获取支付状态
export const getPaymentStatus = (paymentId) => {
  return request({
    url: `/payments/${paymentId}/status`,
    method: 'get'
  })
}

// 微信登录（获取用户信息）
export const wechatLogin = (code) => {
  return request({
    url: '/public/customer/wechat/login',
    method: 'post',
    data: { code }
  })
}

// 获取会员信息
export const getMemberInfo = () => {
  return request({
    url: '/members/my',
    method: 'get'
  })
}

// 根据手机号查询会员
export const getMemberByPhone = (phone) => {
  return request({
    url: `/members/phone/${phone}`,
    method: 'get'
  })
}
