import request from './request'

/**
 * 后台管理支付相关API
 */

/**
 * 支付管理
 */

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

// 获取预约相关的支付
export const getAppointmentPayments = (appointmentId) => {
  return request({
    url: `/admin/appointments/${appointmentId}/payments`,
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

/**
 * 退款管理
 */

// 获取退款列表
export const getRefunds = (params) => {
  return request({
    url: '/admin/refunds',
    method: 'get',
    params
  })
}

// 获取退款详情
export const getRefundById = (id) => {
  return request({
    url: `/admin/refunds/${id}`,
    method: 'get'
  })
}

// 获取预约相关的退款
export const getAppointmentRefunds = (appointmentId) => {
  return request({
    url: `/admin/appointments/${appointmentId}/refunds`,
    method: 'get'
  })
}

/**
 * 预约相关支付
 */

// 获取预约押金状态
export const getAppointmentDepositStatus = (appointmentId) => {
  return request({
    url: `/admin/appointments/${appointmentId}/deposit`,
    method: 'get'
  })
}

// 手动退还押金
export const refundAppointmentDeposit = (appointmentId, reason = '') => {
  return request({
    url: `/admin/appointments/${appointmentId}/deposit/refund`,
    method: 'post',
    data: { reason }
  })
}
