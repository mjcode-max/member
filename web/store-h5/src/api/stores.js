import request from './request'

// 获取门店详情
export const getStoreById = (id) => {
  return request({
    url: `/stores/${id}`,
    method: 'get'
  })
}

// 获取门店列表
export const getStores = (params) => {
  return request({
    url: '/stores',
    method: 'get',
    params
  })
}

