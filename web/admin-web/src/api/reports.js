import request from './request'

// 获取管理员数据看板
export const getAdminDashboard = () => {
  return request({
    url: '/admin/dashboard',
    method: 'get'
  })
}

// 获取报表数据
export const getReports = (params) => {
  return request({
    url: '/admin/reports',
    method: 'get',
    params
  })
}
