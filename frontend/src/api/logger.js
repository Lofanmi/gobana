import request from '@/utils/request'

const basePath = '/api/gobana/v1'

export function loggerSearch(params) {
  return request({
    url: basePath + '/logger/search',
    method: 'post',
    data: params
  })
}
