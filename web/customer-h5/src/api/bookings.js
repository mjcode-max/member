import request from './request'

/**
 * 顾客端预约相关API
 */

// 获取可用时间段
export const getAvailableSlots = (params) => {
  return request({
    url: '/slots/available',
    method: 'get',
    params
  })
}

/**
 * 预约管理
 */

// 创建预约
export const createAppointment = (data) => {
  return request({
    url: '/appointments',
    method: 'post',
    data
  })
}

// 获取我的预约列表
export const getMyAppointments = () => {
  return request({
    url: '/appointments/my',
    method: 'get'
  })
}

// 获取我即将到来的预约
export const getUpcomingAppointments = () => {
  return request({
    url: '/appointments/my/upcoming',
    method: 'get'
  })
}

// 获取预约详情
export const getAppointmentById = (id) => {
  return request({
    url: `/appointments/${id}`,
    method: 'get'
  })
}

// 顾客取消预约（需提前3小时）
export const cancelAppointmentByCustomer = (id, reason = '') => {
  return request({
    url: '/appointments/cancel/customer',
    method: 'post',
    data: {
      appointment_id: id,
      reason
    }
  })
}

/**
 * 支付相关
 */

// 支付押金
export const payDeposit = (appointmentId, paymentMethod = 'wechat') => {
  return request({
    url: '/appointments/pay-deposit',
    method: 'post',
    data: {
      appointment_id: appointmentId,
      payment_method: paymentMethod
    }
  })
}

// 获取门店列表
export const getStores = (params) => {
  return request({
    url: '/public/stores',
    method: 'get',
    params
  })
}

// 获取门店详情
export const getStoreById = (id) => {
  return request({
    url: `/stores/${id}`,
    method: 'get'
  })
}

// 获取门店可用时段
export const getStoreAvailableSlots = (storeId, date) => {
  return request({
    url: '/slots/available',
    method: 'get',
    params: {
      store_id: storeId,
      date
    }
  })
}
