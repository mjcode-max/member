import request from './request'

// 创建支付
export const createPayment = (bookingId, data) => {
  return request({
    url: `/customer/bookings/${bookingId}/payment`,
    method: 'post',
    data
  })
}

// 获取支付状态
export const getPaymentStatus = (paymentId) => {
  return request({
    url: `/customer/payments/${paymentId}/status`,
    method: 'get'
  })
}
