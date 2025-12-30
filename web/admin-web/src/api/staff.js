import request from './request'

// 获取员工列表
export const getStaff = (params) => {
  return request({
    url: '/admin/staff',
    method: 'get',
    params
  })
}

// 获取员工详情
export const getStaffById = (id) => {
  return request({
    url: `/admin/staff/${id}`,
    method: 'get'
  })
}

// 创建员工
export const createStaff = (data) => {
  return request({
    url: '/admin/staff',
    method: 'post',
    data
  })
}

// 更新员工
export const updateStaff = (id, data) => {
  return request({
    url: `/admin/staff/${id}`,
    method: 'put',
    data
  })
}

// 员工离职
export const resignStaff = (id) => {
  return request({
    url: `/admin/staff/${id}/resign`,
    method: 'put'
  })
}

// 员工复职
export const rehireStaff = (id) => {
  return request({
    url: `/admin/staff/${id}/rehire`,
    method: 'put'
  })
}

// 获取员工排班
export const getStaffSchedule = (id, params) => {
  return request({
    url: `/admin/staff/${id}/schedule`,
    method: 'get',
    params
  })
}

// 更新员工排班
export const updateStaffSchedule = (id, data) => {
  return request({
    url: `/admin/staff/${id}/schedule`,
    method: 'put',
    data
  })
}
