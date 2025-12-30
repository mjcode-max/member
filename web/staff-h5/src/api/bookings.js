import request from './request'

// 获取美甲师的工作安排
export const getStaffSchedule = (date) => {
  return request({
    url: '/staff/schedule',
    method: 'get',
    params: { date }
  })
}

// 获取美甲师的预约列表
export const getStaffBookings = (params) => {
  return request({
    url: '/staff/bookings',
    method: 'get',
    params
  })
}

// 确认预约
export const confirmBooking = (bookingId) => {
  return request({
    url: `/staff/bookings/${bookingId}/confirm`,
    method: 'put'
  })
}

// 完成预约
export const completeBooking = (bookingId, data) => {
  return request({
    url: `/staff/bookings/${bookingId}/complete`,
    method: 'put',
    data
  })
}

// 取消预约
export const cancelBooking = (bookingId, reason) => {
  return request({
    url: `/staff/bookings/${bookingId}/cancel`,
    method: 'put',
    data: { reason }
  })
}
