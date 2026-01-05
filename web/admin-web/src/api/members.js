import request from './request'

// 获取会员列表
export const getMembers = (params) => {
  return request({
    url: '/members',
    method: 'get',
    params
  })
}

// 获取会员详情
export const getMemberById = (id) => {
  return request({
    url: `/members/${id}`,
    method: 'get'
  })
}

// 创建会员
export const createMember = (data) => {
  return request({
    url: '/members',
    method: 'post',
    data
  })
}

// 更新会员
export const updateMember = (id, data) => {
  return request({
    url: `/members/${id}`,
    method: 'put',
    data
  })
}

// 更新会员状态
export const updateMemberStatus = (id, status) => {
  return request({
    url: `/members/${id}/status`,
    method: 'put',
    data: { status }
  })
}

// 停用会员
export const disableMember = (id) => {
  return updateMemberStatus(id, 'inactive')
}

// 启用会员
export const enableMember = (id) => {
  return updateMemberStatus(id, 'active')
}

// 获取会员使用记录（消费记录）
export const getMemberUsages = (id) => {
  return request({
    url: `/members/${id}/usages`,
    method: 'get'
  })
}

// 创建使用记录
export const createUsage = (memberId, data) => {
  return request({
    url: `/members/${memberId}/usages`,
    method: 'post',
    data
  })
}

// 删除使用记录
export const deleteUsage = (usageId) => {
  return request({
    url: `/usages/${usageId}`,
    method: 'delete'
  })
}

// 获取所有使用记录
export const getUsages = (params) => {
  return request({
    url: '/usages',
    method: 'get',
    params
  })
}
