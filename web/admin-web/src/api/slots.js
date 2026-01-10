import request from './request'

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

