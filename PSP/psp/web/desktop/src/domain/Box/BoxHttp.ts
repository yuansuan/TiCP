/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

/**
 * Fetch 模块用于处理原生 http 请求
 * 注：Http 模块用于处理经过 apiFramework format 的请求
 */

import axios from 'axios'
import qs from 'qs'
import { env } from '@/domain'
import { AxiosInstance } from '@/utils'
import { single } from '@/utils'
import { errorMapping } from './errorMapping'
import { message } from 'antd'
import { Modal } from '@/components'
import { Http } from '@/utils'
import { globalSizes } from './states'

const BoxHttp: AxiosInstance = <any>axios.create({
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

BoxHttp.interceptors.request.use(config => {
  return config
})

BoxHttp.interceptors.request.use(config => {
  // if (!env.isPersonal && config.url.includes('mkdir')) {
  //   try {
  //     Http.post('/filerecord/record', {
  //       type: 5,
  //       info: {
  //         storage_size: 0,
  //         file_name: config.data.path.split('/').pop(),
  //         file_type: 2
  //       }
  //     })
  //   } catch (error) {
  //     console.error(error)
  //   }
  // }
  return config
})

const uploadDirectories = {}

BoxHttp.interceptors.response.use(response => {
  if (response.config.url.includes('filemanager/upload')) {
    if (response.data.success) {
      const query = response.config.url.substring(
        response.config.url.indexOf('?') + 1
      )
      const key = qs.parse(query)

      if (!(key?.isEdit && key?.isEdit == 'true')) {
        if (key?.directory === 'true') {
          if (key._uid && !uploadDirectories[key._uid]) {
            uploadDirectories[key._uid] = true
            const subLength = key.dir.length === 1 ? 1 : 2
            // Http.post('/filerecord/record', {
            //   type: 1,
            //   info: {
            //     storage_size: globalSizes[key._uid],
            //     file_name: key.path
            //       .substring(key.dir.length + subLength)
            //       .split('/')[0],
            //     file_type: 2
            //   }
            // })
          }
        } else {
          if (key.finish === 'true') {
            // Http.post('/filerecord/record', {
            //   type: 1,
            //   info: {
            //     storage_size: key.file_size,
            //     file_name: key.path.split('/').pop(),
            //     file_type: 1
            //   }
            // })
          }
        }
      }
    }
    return response
  }

  if (response.config.url.includes('download')) {
    if (response.data.success) {
      const { types, paths: namePaths } = JSON.parse(response.config.data)

      if (!types) return response

      const names = namePaths.map(path => path.split('/').pop())
      const name =
        names.length === 1
          ? names[0]
          : `[批量下载]${
              names.length > 2
                ? names.slice(0, 2).join(',') + '等.zip'
                : names.join(',') + '.zip'
            }`
      const {
        data: { total_size }
      } = response.data
      // Http.post('/filerecord/record', {
      //   type: 2,
      //   info: {
      //     storage_size: total_size,
      //     file_name: name,
      //     file_type: types.length === 1 ? (types[0] ? 1 : 2) : 3 || 0
      //   }
      // })
    }
  }

  return response
})

BoxHttp.interceptors.response.use(response => {
  const {
    data: { code, ...data }
  } = response

  if (!data.success) {
    const { disableErrorMessage, formatErrorMessage } = response.config

    if (!disableErrorMessage) {
      let msg
      if (Object.keys(errorMapping).includes('' + code)) {
        msg = errorMapping[code]
      } else {
        msg = data.message
      }

      if (formatErrorMessage) {
        msg = formatErrorMessage(msg)
      }
      msg && message.error(msg)
    }

    return Promise.reject({ response })
  }
  return response.data
})

export default BoxHttp
