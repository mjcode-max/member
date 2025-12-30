import request from './request'

// 获取店长端仪表板数据
export const getStoreDashboard = () => {
  return request({
    url: '/store/dashboard',
    method: 'get'
  })
}

// 获取门店统计信息
export const getStoreStats = (params) => {
  return request({
    url: '/store/stats',
    method: 'get',
    params
  })
}

// 获取今日数据概览
export const getTodayOverview = () => {
  return request({
    url: '/store/today-overview',
    method: 'get'
  })
}

// 获取员工状态统计
export const getStaffStats = () => {
  return request({
    url: '/store/staff-stats',
    method: 'get'
  })
}

// 获取最新会员列表
export const getRecentMembers = (params) => {
  return request({
    url: '/store/recent-members',
    method: 'get',
    params
  })
}

