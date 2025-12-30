import request from './request'

// 获取会员列表
export const getMembers = (params) => {
  return request({
    url: '/store/members',
    method: 'get',
    params
  })
}

// 创建会员
export const createMember = (data) => {
  return request({
    url: '/store/members',
    method: 'post',
    data
  })
}

// 上传会员人脸照片
export const uploadFaceImage = (memberId, formData) => {
  return request({
    url: `/store/members/${memberId}/face`,
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
