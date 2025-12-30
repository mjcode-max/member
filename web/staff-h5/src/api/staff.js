import request from './request'

// 更新我的工作状态
export const updateMyWorkStatus = (data) => {
  return request({
    url: '/staff/work-status',
    method: 'put',
    data
  })
}

// 获取今日预约
export const getTodayBookings = () => {
  return request({
    url: '/staff/bookings/today',
    method: 'get'
  })
}

// 获取我的预约列表
export const getMyBookings = (params) => {
  return request({
    url: '/staff/bookings',
    method: 'get',
    params
  })
}

// 完成预约服务
export const completeBooking = (bookingId, data) => {
  return request({
    url: `/staff/bookings/${bookingId}/complete`,
    method: 'put',
    data
  })
}

// 取消预约
export const cancelBooking = (bookingId, data) => {
  return request({
    url: `/staff/bookings/${bookingId}/cancel`,
    method: 'put',
    data
  })
}
