/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { ICustomReq } from '../@types'
import { AxiosInstance } from 'axios'

type Props = {
  httpAdapter: AxiosInstance
  preUploadUrl?: string
}

export class BaseController {
  httpAdapter: AxiosInstance
  req: ICustomReq
  preUploadUrl: string

  constructor({ httpAdapter, preUploadUrl }: Props) {
    this.httpAdapter = httpAdapter
    this.preUploadUrl = preUploadUrl
  }

  name: string = ''

  setReq(req) {
    this.req = req
  }

  // 上传
  upload() {}

  // 取消
  abort() {}

  // 重试
  retry() {}

  // 暂停
  pause() {
    throw new Error('not implemented by current controller')
  }

  // 恢复
  resume() {
    throw new Error('not implemented by current controller')
  }
}
