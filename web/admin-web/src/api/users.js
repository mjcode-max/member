import request from './request'

// 获取用户列表
export const getUsers = (params) => {
  return request({
    url: '/users',
    method: 'get',
    params
  })
}

// 获取用户详情
export const getUserById = (id) => {
  return request({
    url: `/users/${id}`,
    method: 'get'
  })
}

// 创建用户
export const createUser = (data) => {
  return request({
    url: '/users',
    method: 'post',
    data
  })
}

// 更新用户
export const updateUser = (id, data) => {
  return request({
    url: `/users/${id}`,
    method: 'put',
    data
  })
}

// 更新用户状态
export const updateUserStatus = (id, status) => {
  return request({
    url: `/users/${id}/status`,
    method: 'put',
    data: { status }
  })
}

// 更新美甲师工作状态
export const updateWorkStatus = (id, workStatus) => {
  return request({
    url: `/users/${id}/work-status`,
    method: 'put',
    data: { work_status: workStatus }
  })
}

