import request from './request'

/**
 * 门店端预约相关API
 */

/**
 * 预约管理
 */

// 获取门店预约列表
export const getStoreAppointments = (params) => {
  return request({
    url: '/appointments/store',
    method: 'get',
    params
  })
}

// 获取预约详情
export const getAppointmentById = (id) => {
  return request({
    url: `/appointments/${id}`,
    method: 'get'
  })
}

// 确认顾客到店（退还押金）
export const confirmArrival = (appointmentId) => {
  return request({
    url: '/appointments/confirm-arrival',
    method: 'post',
    data: {
      appointment_id: appointmentId
    }
  })
}

// 完成预约服务
export const completeAppointment = (appointmentId) => {
  return request({
    url: '/appointments/complete',
    method: 'post',
    data: {
      appointment_id: appointmentId
    }
  })
}

// 美甲师取消预约
export const cancelAppointmentByTechnician = (appointmentId, reason = '') => {
  return request({
    url: '/appointments/cancel/technician',
    method: 'post',
    data: {
      appointment_id: appointmentId,
      reason
    }
  })
}

// 根据日期范围获取预约
export const getAppointmentsByDateRange = (startDate, endDate) => {
  return request({
    url: '/appointments/store',
    method: 'get',
    params: {
      start_date: startDate,
      end_date: endDate
    }
  })
}

/**
 * 美甲师管理
 */

// 获取美甲师列表
export const getTechnicians = (params) => {
  return request({
    url: '/users',
    method: 'get',
    params: {
      ...params,
      role: 'technician'
    }
  })
}

// 获取美甲师详情
export const getTechnicianById = (id) => {
  return request({
    url: `/users/${id}`,
    method: 'get'
  })
}

// 更新美甲师工作状态
export const updateTechnicianWorkStatus = (id, workStatus) => {
  return request({
    url: `/users/${id}/work-status`,
    method: 'put',
    data: { work_status: workStatus }
  })
}

/**
 * 时段管理
 */

// 获取可用时间段
export const getAvailableSlots = (params) => {
  return request({
    url: '/slots/available',
    method: 'get',
    params
  })
}

// 获取所有时段列表
export const getSlots = (params) => {
  return request({
    url: '/slots',
    method: 'get',
    params
  })
}

// 生成时段
export const generateSlots = (data) => {
  return request({
    url: '/slots/generate',
    method: 'post',
    data
  })
}

/**
 * 门店管理
 */

// 获取当前门店信息
export const getCurrentStore = () => {
  return request({
    url: '/stores/my',
    method: 'get'
  })
}

// 更新门店信息
export const updateStore = (id, data) => {
  return request({
    url: `/stores/${id}`,
    method: 'put',
    data
  })
}

/**
 * 会员管理
 */

// 获取会员列表
export const getMembers = (params) => {
  return request({
    url: '/members',
    method: 'get',
    params
  })
}

// 获取会员详情
export const getMemberById = (id) => {
  return request({
    url: `/members/${id}`,
    method: 'get'
  })
}

// 根据手机号查询会员
export const getMemberByPhone = (phone) => {
  return request({
    url: `/members/phone/${phone}`,
    method: 'get'
  })
}

// 创建会员
export const createMember = (data) => {
  return request({
    url: '/members',
    method: 'post',
    data
  })
}

// 更新会员信息
export const updateMember = (id, data) => {
  return request({
    url: `/members/${id}`,
    method: 'put',
    data
  })
}

