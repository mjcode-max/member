import request from './request'

// 获取支付列表
export const getPayments = (params) => {
  return request({
    url: '/admin/payments',
    method: 'get',
    params
  })
}

// 获取支付详情
export const getPaymentById = (id) => {
  return request({
    url: `/admin/payments/${id}`,
    method: 'get'
  })
}

// 申请退款
export const refundPayment = (id, data) => {
  return request({
    url: `/admin/payments/${id}/refund`,
    method: 'post',
    data
  })
}

// 重试支付
export const retryPayment = (id) => {
  return request({
    url: `/admin/payments/${id}/retry`,
    method: 'post'
  })
}

// 导出支付记录
export const exportPayments = (params) => {
  return request({
    url: '/admin/payments/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}
