import request from './request'

// 获取会员列表
export const getMembers = (params) => {
  return request({
    url: '/admin/members',
    method: 'get',
    params
  })
}

// 获取会员详情
export const getMemberById = (id) => {
  return request({
    url: `/admin/members/${id}`,
    method: 'get'
  })
}

// 创建会员
export const createMember = (data) => {
  return request({
    url: '/admin/members',
    method: 'post',
    data
  })
}

// 更新会员
export const updateMember = (id, data) => {
  return request({
    url: `/admin/members/${id}`,
    method: 'put',
    data
  })
}

// 停用会员
export const disableMember = (id) => {
  return request({
    url: `/admin/members/${id}/disable`,
    method: 'put'
  })
}

// 启用会员
export const enableMember = (id) => {
  return request({
    url: `/admin/members/${id}/enable`,
    method: 'put'
  })
}

// 获取会员消费记录
export const getMemberConsumptions = (id, params) => {
  return request({
    url: `/admin/members/${id}/consumptions`,
    method: 'get',
    params
  })
}

// 上传会员人脸
export const uploadMemberFace = (id, data) => {
  return request({
    url: `/admin/members/${id}/face`,
    method: 'post',
    data
  })
}
