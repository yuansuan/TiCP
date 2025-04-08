import axios from 'axios'
import { message } from 'antd'
import qs from 'qs'

import { ignoredErrorCodes, overrideErrorCodes } from '@/utils/ErrorCodes'
import wrapClient from './wrapClient'
import IAxiosInstance from './IAxiosInstance'
import cacheAdapter from './CacheAdapter'

const Http: IAxiosInstance = <any>axios.create({
  baseURL: '/api/v1',
  timeout: 20 * 60000,
  withCredentials: true,
  adapter: cacheAdapter(),
  paramsSerializer: function (params) {
    return qs.stringify(params, { arrayFormat: 'repeat'})
  }
})

wrapClient(Http)
Http.interceptors.response.use(response => {
  const { data } = response
  
  if (!data.success) {
    const { disableErrorMessage, formatErrorMessage } = response.config
    if (!disableErrorMessage && !ignoredErrorCodes.includes(data.code)) {
      let msg = overrideErrorCodes.get(data.code) || data.message
      if (formatErrorMessage) {
        msg = formatErrorMessage(msg)
      }
      message.error(msg)
    }
    if(data.code ===11013){
      return ''
    }else {
      return Promise.reject(data.code)
    } 
    
  }
  return data
})

export default Http
