import request from './request'

// 更新我的工作状态（美甲师端：更新自己的工作状态）
export const updateMyWorkStatus = (data, userId) => {
  // userId可以从外部传入，或者从userStore获取
  if (!userId) {
    // 尝试从localStorage获取
    try {
      const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
      userId = userInfo.id
    } catch (e) {
      console.error('获取用户ID失败:', e)
    }
  }
  
  if (!userId) {
    return Promise.reject(new Error('无法获取用户ID'))
  }
  
  return request({
    url: `/users/${userId}/work-status`,
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
