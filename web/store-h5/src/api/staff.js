import request from './request'

// 获取员工列表（店长端：获取自己门店的美甲师）
export const getStaffList = () => {
  return request({
    url: '/users',
    method: 'get',
    params: {
      role: 'technician'
    }
  })
}

// 创建员工（店长端：创建美甲师）
export const createStaff = (data) => {
  return request({
    url: '/users',
    method: 'post',
    data: {
      ...data,
      role: 'technician' // 固定为美甲师
    }
  })
}

// 更新员工工作状态
export const updateStaffStatus = (id, data) => {
  return request({
    url: `/users/${id}/work-status`,
    method: 'put',
    data
  })
}
