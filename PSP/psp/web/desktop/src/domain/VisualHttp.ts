/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import axios from 'axios'
import qs from 'qs'

import {
  ignoredErrorCodes,
  overrideErrorCodes
} from './Visualization/ErrorCodes'
import { AxiosInstance } from '@/utils'
import { message } from 'antd'
import { Http } from '@/utils'

let token: string = ''
export function setToken(t: string) {
  token = t
}
export function getToken() {
  return token
}
export async function init() {
  const res = await Http.get('/visual/token')
  setToken(res.data)
}

export const VisualHttp: AxiosInstance = <any>axios.create({
  timeout: 20 * 60000,
  paramsSerializer: function (params) {
    return qs.stringify(params)
  }
})

VisualHttp.interceptors.request.use(config => {
  config.baseURL = '/visual'
  config.headers['Authorization'] = `Bearer ${getToken()}`
  return config
})
VisualHttp.interceptors.response.use(response => {
  const { data } = response
  const { disableErrorMessage, formatErrorMessage } = response.config
  if (!data.success) {
    if (!disableErrorMessage && !ignoredErrorCodes.includes(data.errorCode)) {
      let msg = overrideErrorCodes.get(data.errorCode) || data.message
      if (formatErrorMessage) {
        msg = formatErrorMessage(msg)
      }
      message.error(msg)
    }
    return Promise.reject(data.errorCode)
  }
  return data
})
