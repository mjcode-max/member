import request from './request'

// 获取时段模板列表
export const getTemplates = (params) => {
  return request({
    url: '/slot-templates',
    method: 'get',
    params
  })
}

// 获取时段模板详情
export const getTemplateById = (id) => {
  return request({
    url: `/slot-templates/${id}`,
    method: 'get'
  })
}

// 创建时段模板
export const createTemplate = (data) => {
  return request({
    url: '/slot-templates',
    method: 'post',
    data
  })
}

// 更新时段模板
export const updateTemplate = (id, data) => {
  return request({
    url: `/slot-templates/${id}`,
    method: 'put',
    data
  })
}

// 删除时段模板
export const deleteTemplate = (id) => {
  return request({
    url: `/slot-templates/${id}`,
    method: 'delete'
  })
}

