import request from './request'

// 验证会员身份
export const verifyMember = (data) => {
  return request({
    url: '/staff/members/verify',
    method: 'post',
    data
  })
}

// 获取会员信息
export const getMemberInfo = (memberId) => {
  return request({
    url: `/staff/members/${memberId}`,
    method: 'get'
  })
}

// 记录会员消费
export const recordConsumption = (data) => {
  return request({
    url: '/staff/members/consumption',
    method: 'post',
    data
  })
}
