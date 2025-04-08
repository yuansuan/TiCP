/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

export interface ICustomReq {
  onProgress: (event: { percent: number }, file: File) => void
  onError: (event: Error, body?: Object) => void
  onSuccess: (body: Object, file: File) => void
  data: any
  filename: string
  file: File
  withCredentials: boolean
  action: string
  headers: Object
}
