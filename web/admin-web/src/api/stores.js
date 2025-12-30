import request from './request'

// 获取门店列表
export const getStores = (params) => {
  return request({
    url: '/admin/stores',
    method: 'get',
    params
  })
}

// 获取门店详情
export const getStoreById = (id) => {
  return request({
    url: `/admin/stores/${id}`,
    method: 'get'
  })
}

// 创建门店
export const createStore = (data) => {
  return request({
    url: '/admin/stores',
    method: 'post',
    data
  })
}

// 更新门店
export const updateStore = (id, data) => {
  return request({
    url: `/admin/stores/${id}`,
    method: 'put',
    data
  })
}

// 删除门店
export const deleteStore = (id) => {
  return request({
    url: `/admin/stores/${id}`,
    method: 'delete'
  })
}
