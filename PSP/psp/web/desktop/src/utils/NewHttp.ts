/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

/**
 * Fetch 模块用于处理原生 http 请求
 * 注：Http 模块用于处理经过 apiFramework format 的请求
 */

import axios from 'axios'
import qs from 'qs'
import {
  AxiosInstance,
  interceptError,
  interceptResponse,
  createHttp
} from '@/utils'
import { message } from 'antd'
import { env } from '@/domain'
import { Modal } from '@/components'
import { single } from '@/utils'

const setHeader = (config, key: string, value: string) =>
  value && (config.headers[key] = value)

const interceptRequest = Http => {
  Http.interceptors.request.use(config => {
    return config
  })
}
export const NewRawHttp = (url: string) => {
  const NewRawHttp = createHttp()
  NewRawHttp.interceptors.request.use(config => {
    // TODO: 根据分区地址设置baseURL
    config.baseURL = `${url}/api`
    return config
  })
  interceptRequest(NewRawHttp)
  return NewRawHttp
}

export const NewHttp = (url: string): AxiosInstance => {
  const NewHttp: AxiosInstance = <any>axios.create({
    timeout: 20 * 60000,
    withCredentials: true,
    paramsSerializer: function (params) {
      const o = Object.keys(params).reduce((obj, key) => {
        const item = params[key]
        if (Object.prototype.toString.call(item) === '[object Object]') {
          obj[key] = JSON.stringify(item)
        } else {
          obj[key] = item
        }
        return obj
      }, {})

      return qs.stringify(o)
    }
  })

  NewHttp.interceptors.request.use(config => {
    return config
  })
  interceptRequest(NewHttp)
  interceptError(NewHttp)
  interceptResponse(NewHttp)
  NewHttp.interceptors.response.use(
    response => {
      const { data } = response
      if (!data.success) {
        const { disableErrorMessage, formatErrorMessage } = response.config

        if (data.code === 120006) {
          single('user-not-exist-modal', () =>
            Modal.showConfirm({
              title: '消息提示',
              content: '账号异常，点击确认重新登录',
              closable: false,
              CancelButton: null
            }).then(() => {
              // env.logout()
            })
          )
        }
        let msg
        if (!disableErrorMessage) {
          msg = data.message
          if (formatErrorMessage) {
            msg = formatErrorMessage(msg)
          }
          msg && message.error(msg)
        }

        return Promise.reject({
          ...data,
          message: msg
        })
      }
      return data
    },
    error => {
      const { response } = error
      if (response) {
        switch (response.status) {
          case 409: {
            return single('login-conflict-modal', () =>
              Modal.showConfirm({
                cancelButtonProps: { style: { display: 'none' } },
                closable: false,
                content: '用户登录冲突，请重新登录'
              }).then(() => {
                // env.logout()
              })
            )
          }
          case 401:
            // window.location.replace(`/api/sso/login${window.location.hash}`)
            break
          default:
            break
        }
      }

      return Promise.reject(error)
    }
  )
  return NewHttp
}
