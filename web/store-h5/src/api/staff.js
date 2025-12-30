import request from './request'

// 获取员工列表
export const getStaffList = () => {
  return request({
    url: '/store/staff',
    method: 'get'
  })
}

// 创建员工
export const createStaff = (data) => {
  return request({
    url: '/store/staff',
    method: 'post',
    data
  })
}

// 更新员工状态
export const updateStaffStatus = (id, data) => {
  return request({
    url: `/store/staff/${id}/status`,
    method: 'put',
    data
  })
}
