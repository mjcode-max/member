import request from './request'

// 获取门店列表
export const getStores = (params) => {
  return request({
    url: '/stores',
    method: 'get',
    params
  })
}

// 获取门店详情
export const getStoreById = (id) => {
  return request({
    url: `/stores/${id}`,
    method: 'get'
  })
}

// 创建门店
export const createStore = (data) => {
  return request({
    url: '/stores',
    method: 'post',
    data
  })
}

// 更新门店
export const updateStore = (id, data) => {
  return request({
    url: `/stores/${id}`,
    method: 'put',
    data
  })
}

// 删除门店
export const deleteStore = (id) => {
  return request({
    url: `/stores/${id}`,
    method: 'delete'
  })
}

// 更新门店状态
export const updateStoreStatus = (id, status) => {
  return request({
    url: `/stores/${id}/status`,
    method: 'put',
    data: { status }
  })
}
