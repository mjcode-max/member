import request from './request'

// 获取预约可用时间段
export const getAvailableSlots = (params) => {
  return request({
    url: '/slots/available',
    method: 'get',
    params
  })
}

// 创建预约
export const createBooking = (data) => {
  return request({
    url: '/customer/bookings',
    method: 'post',
    data
  })
}

// 获取用户预约列表
export const getUserBookings = (params) => {
  return request({
    url: '/customer/bookings',
    method: 'get',
    params
  })
}

// 获取我的预约列表（别名）
export const getMyBookings = getUserBookings

// 获取预约详情
export const getBookingById = (bookingId) => {
  return request({
    url: `/customer/bookings/${bookingId}`,
    method: 'get'
  })
}

// 取消预约
export const cancelBooking = (bookingId, data) => {
  return request({
    url: `/customer/bookings/${bookingId}/cancel`,
    method: 'put',
    data
  })
}
