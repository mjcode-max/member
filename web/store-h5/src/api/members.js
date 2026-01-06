import request from './request'

// 获取会员列表
export const getMembers = (params) => {
  return request({
    url: '/members',
    method: 'get',
    params
  })
}

// 创建会员（支持FormData）
export const createMember = (data) => {
  // 如果是FormData，不设置Content-Type，让浏览器自动设置
  const config = {
    url: '/members',
    method: 'post',
    data
  }
  
  // FormData需要特殊处理，不设置Content-Type
  if (data instanceof FormData) {
    config.headers = {
      'Content-Type': 'multipart/form-data'
    }
  }
  
  return request(config)
}

// 上传会员人脸照片
export const uploadFaceImage = (memberId, formData) => {
  return request({
    url: `/members/${memberId}/face`,
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取会员详情
export const getMemberById = (memberId) => {
  return request({
    url: `/members/${memberId}`,
    method: 'get'
  })
}

// 获取会员使用记录列表
export const getMemberUsages = (memberId) => {
  return request({
    url: `/members/${memberId}/usages`,
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