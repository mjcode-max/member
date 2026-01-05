import request from './request'

// 获取门店列表（公开接口，无需登录）
// 支持筛选：status (operating/closed/shutdown), name (门店名称模糊搜索)
// 支持分页：page, page_size
// 默认只返回营业中的门店（status=operating）
export const getStores = (params) => {
  return request({
    url: '/public/stores',
    method: 'get',
    params
  })
}

// 获取门店详情（如果需要，可以后续添加公开接口）
export const getStoreById = (id) => {
  return request({
    url: `/stores/${id}`,
    method: 'get'
  })
}
