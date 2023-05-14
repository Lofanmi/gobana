import request from '@/utils/request'

const basePath = '/api/gobana/v1'

export function configGetBackendList(params) {
  return request({
    url: basePath + '/config/backend_list',
    method: 'get',
    params
  })
}

export function configGetStorageList(params) {
  return request({
    url: basePath + '/config/storage_list',
    method: 'get',
    params
  })
}
