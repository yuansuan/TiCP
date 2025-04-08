/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import message from 'antd/lib/message'
import 'antd/lib/message/style/index.css'
import { single } from '@/utils'

/**
 * 不需要全局的错误提示消息时，请求需传入disableErrorMessage: true
 * 需要根据不同错误码做不同处理时，自行进行try catch
 * example:
 *  try {
 *    await Http.post('/api/test/url', data, { disableErrorMessage: true })
 *  } catch(e) {
 *    // custom error handle
 * }
 */

export const interceptError = instance => {
  instance.interceptors.response.use(
    response => response,
    error => {
      if (error.message === 'Network Error') {
        single(
          'http-client-network-error-message',
          () => message.error('网络异常').promise
        )
      }

      const { response } = error
      if (response) {
        const { formatErrorMessage, disableErrorMessage } = response.config
        if (
          !disableErrorMessage &&
          response &&
          response.data &&
          response.data.message
        ) {
          let msg = response.data.message
          if (formatErrorMessage) {
            msg = formatErrorMessage(msg)
          }
          message.error(msg)
        }
      }

      return Promise.reject(error)
    }
  )
}

export const interceptResponse = instance => {
  instance.interceptors.response.use(response => {
    const { data } = response
    if (!data.success) {
      const { disableErrorMessage, formatErrorMessage } = response.config
      let msg
      if (!disableErrorMessage) {
        msg = data.message
        if (formatErrorMessage) {
          msg = formatErrorMessage(msg)
        }
        message.error(msg)
      }
      return Promise.reject({
        ...data,
        message: msg,
      })
    }

    return response
  })
}
