/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import axios, { AxiosRequestConfig } from 'axios'
import qs from 'query-string'
import { interceptError, interceptResponse } from './wrapClient'
import { AxiosInstance } from './IAxiosInstance'

export const Http = createHttp()
interceptError(Http)
interceptResponse(Http)

export function createHttp<T = AxiosInstance>(config?: AxiosRequestConfig): T {
  const Http: AxiosInstance = <any>axios.create({
    baseURL: '/api/',
    timeout: 60000,
    withCredentials: true,
    paramsSerializer: function(params) {
      return qs.stringify(params, { arrayFormat: 'bracket' })
    },
    ...config,
  })

  return Http as any
}
