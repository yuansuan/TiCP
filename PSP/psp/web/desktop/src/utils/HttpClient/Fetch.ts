/**
 * Fetch 模块用于处理原生 http 请求
 * 注：Http 模块用于处理经过 apiFramework format 的请求
 */

import axios from 'axios'
import qs from 'qs'

import wrapClient from './wrapClient'
import IAxiosInstance from './IAxiosInstance'

const Fetch: IAxiosInstance = <any>axios.create({
  baseURL: '/api/v1',
  timeout: 20 * 60000,
  withCredentials: true,
  paramsSerializer: function (params) {
    return qs.stringify(params)
  }
})

wrapClient(Fetch)
Fetch.interceptors.response.use(response => response.data)

export default Fetch
