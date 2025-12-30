import request from './request'

// 获取门店列表
export const getStores = (params) => {
  return request({
    url: '/customer/stores',
    method: 'get',
    params
  })
}

// 获取门店详情
export const getStoreById = (id) => {
  return request({
    url: `/customer/stores/${id}`,
    method: 'get'
  })
}
