import request from './request'

// 注册会员
export const registerMember = (data) => {
  return request({
    url: '/customer/members/register',
    method: 'post',
    data
  })
}

// 获取会员信息
export const getMemberInfo = (phone) => {
  return request({
    url: `/customer/members/phone/${phone}`,
    method: 'get'
  })
}

// 获取会员码
export const getMemberCode = (memberId) => {
  return request({
    url: `/customer/members/${memberId}/code`,
    method: 'get'
  })
}

// 生成会员码
export const generateMemberCode = (data) => {
  return request({
    url: '/customer/members/verify',
    method: 'post',
    data
  })
}