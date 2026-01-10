import request from './request'

/**
 * 后台管理预约相关API
 */

/**
 * 预约管理
 */

// 获取预约列表
export const getAppointments = (params) => {
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

// 根据日期范围获取预约
export const getAppointmentsByDateRange = (startDate, endDate, storeId = null) => {
  const params = {
    start_date: startDate,
    end_date: endDate
  }
  if (storeId) {
    params.store_id = storeId
  }
  return request({
    url: '/appointments/store',
    method: 'get',
    params
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

// 完成预约
export const completeAppointment = (appointmentId) => {
  return request({
    url: '/appointments/complete',
    method: 'post',
    data: {
      appointment_id: appointmentId
    }
  })
}

// 后台取消预约
export const cancelAppointment = (appointmentId, reason = '') => {
  return request({
    url: '/appointments/cancel/technician',
    method: 'post',
    data: {
      appointment_id: appointmentId,
      reason
    }
  })
}

/**
 * 时段管理
 */

// 获取时段列表
export const getSlots = (params) => {
  return request({
    url: '/slots',
    method: 'get',
    params
  })
}

// 获取可用时段列表
export const getAvailableSlots = (params) => {
  return request({
    url: '/slots/available',
    method: 'get',
    params
  })
}

// 获取时段详情
export const getSlotById = (id) => {
  return request({
    url: `/slots/${id}`,
    method: 'get'
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

// 锁定时段
export const lockSlot = (slotId, count = 1) => {
  return request({
    url: '/slots/lock',
    method: 'post',
    data: {
      slot_id: slotId,
      count
    }
  })
}

// 解锁时段
export const unlockSlot = (slotId, count = 1) => {
  return request({
    url: '/slots/unlock',
    method: 'post',
    data: {
      slot_id: slotId,
      count
    }
  })
}

// 预约时段
export const bookSlot = (slotId, count = 1) => {
  return request({
    url: '/slots/book',
    method: 'post',
    data: {
      slot_id: slotId,
      count
    }
  })
}

// 释放时段
export const releaseSlot = (slotId, count = 1) => {
  return request({
    url: '/slots/release',
    method: 'post',
    data: {
      slot_id: slotId,
      count
    }
  })
}

// 重新计算时段容量
export const recalculateSlotCapacity = (data) => {
  return request({
    url: '/slots/recalculate-capacity',
    method: 'post',
    data
  })
}

/**
 * 时段模板管理
 */

// 获取时段模板列表
export const getTemplates = (params) => {
  return request({
    url: '/slot-templates',
    method: 'get',
    params
  })
}

// 获取时段模板详情
export const getTemplateById = (id) => {
  return request({
    url: `/slot-templates/${id}`,
    method: 'get'
  })
}

// 创建时段模板
export const createTemplate = (data) => {
  return request({
    url: '/slot-templates',
    method: 'post',
    data
  })
}

// 更新时段模板
export const updateTemplate = (id, data) => {
  return request({
    url: `/slot-templates/${id}`,
    method: 'put',
    data
  })
}

// 删除时段模板
export const deleteTemplate = (id) => {
  return request({
    url: `/slot-templates/${id}`,
    method: 'delete'
  })
}
