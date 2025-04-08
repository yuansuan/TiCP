/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { BaseController } from './BaseController'
import axios from 'axios'
import { message } from 'antd'

// 全量上传控制器
export class FullController extends BaseController {
  name = 'full'
  source = axios.CancelToken.source()

  constructor(props) {
    super(props)
  }

  upload() {
    const {
      action,
      filename,
      file,
      data,
      headers,
      withCredentials,
      onSuccess,
      onError,
      onProgress,
    } = this.req
    // 构造要上传的数据
    const formData = new FormData()
    formData.append('file', file)
    formData.append('filename', filename)
    Object.keys(data).forEach(key => formData.append(key, data[key]))
    this.httpAdapter
      .post(action, formData, {
        headers,
        withCredentials,
        onUploadProgress: (e: any) => {
          const percent = (e.loaded / e.total) * 100
          onProgress({ percent }, file)
        },
        cancelToken: this.source.token,
      })
      .then(res => {
        onSuccess(res, file)
      })
      .catch(e => {
        if (!(e instanceof axios.Cancel)) {
          onError(e)
          message.error(`${file.name}上传失败`)
        }
      })
  }

  abort() {
    this.source.cancel('canceled by the user')
  }

  retry() {
    this.upload()
  }
}
