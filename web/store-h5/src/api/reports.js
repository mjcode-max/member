import request from './request'

// 获取店长数据看板
export const getStoreDashboard = () => {
  return request({
    url: '/store/dashboard',
    method: 'get'
  })
}
