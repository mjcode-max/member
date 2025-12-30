import request from './request'

// 获取预约列表
export const getBookings = (params) => {
  return request({
    url: '/admin/bookings',
    method: 'get',
    params
  })
}

// 获取预约详情
export const getBookingById = (id) => {
  return request({
    url: `/admin/bookings/${id}`,
    method: 'get'
  })
}

// 确认预约
export const confirmBooking = (id) => {
  return request({
    url: `/admin/bookings/${id}/confirm`,
    method: 'put'
  })
}

// 完成预约
export const completeBooking = (id) => {
  return request({
    url: `/admin/bookings/${id}/complete`,
    method: 'put'
  })
}

// 取消预约
export const cancelBooking = (id) => {
  return request({
    url: `/admin/bookings/${id}/cancel`,
    method: 'put'
  })
}
