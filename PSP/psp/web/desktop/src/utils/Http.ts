/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import {
  createHttp,
  interceptError,
  interceptResponse
} from './HttpClient2'
import qs from 'qs'

function interceptRequest(Http) {
  Http.interceptors.request.use(config => {
    return config
  })
}

export const RawHttp = createHttp({
  baseURL: '/api/v1',
  timeout: 20 * 60000,
  withCredentials: true,
  paramsSerializer: function (params) {
    return qs.stringify(params)
  }
})
interceptRequest(RawHttp)

export const Http = createHttp({
  paramsSerializer: function (params) {
    return qs.stringify(params)
  }
})

interceptRequest(Http)
interceptError(Http)
interceptResponse(Http)
Http.interceptors.response.use(
  response => {
    const { data } = response

    return data
  },
  error => {
    const { response } = error
    if (response) {
      switch (response.status) {
        case 401:
          localStorage.setItem('needLogin', 'true')
          break
        default:
          break
      }
    }

    return Promise.reject(error)
  }
)
