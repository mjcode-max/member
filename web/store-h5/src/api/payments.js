import request from './request'

/**
 * 门店端支付相关API
 */

// 创建支付
export const createPayment = (data) => {
  return request({
    url: '/payments',
    method: 'post',
    data
  })
}

// 获取支付记录列表
export const getPayments = (params) => {
  return request({
    url: '/payments',
    method: 'get',
    params
  })
}

// 获取支付详情
export const getPaymentById = (id) => {
  return request({
    url: `/payments/${id}`,
    method: 'get'
  })
}

// 获取支付状态
export const getPaymentStatus = (paymentId) => {
  return request({
    url: `/payments/${paymentId}/status`,
    method: 'get'
  })
}

// 获取门店收入统计
export const getStoreRevenue = (params) => {
  return request({
    url: '/reports/revenue',
    method: 'get',
    params
  })
}

// 获取预约统计
export const getAppointmentStats = (params) => {
  return request({
    url: '/reports/appointments',
    method: 'get',
    params
  })
}

