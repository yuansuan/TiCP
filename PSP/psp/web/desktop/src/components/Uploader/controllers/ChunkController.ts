/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import axios, { CancelTokenSource, AxiosInstance } from 'axios'
import * as qs from 'querystring'
import { message } from 'antd'
import { BaseController } from './BaseController'
import { UPLOAD_CHUNK_SIZE } from '../constant'

interface ChunkUploadQuery {
  upload_id: string
  path: string
  file_size: number
  offset: number
  slice_size: number
  finish: boolean
  project_id: string
  bucket: string
}

type Props = {
  httpAdapter: AxiosInstance
  chunkSize?: number
  preUploadUrl?: string
}

export const formatPath = (path: string): string => {
  if (!path?.length) {
    return './'
  } else {
    if (/^\.\/.*$/.test(path)) {
      return path
    } else if (/^\/.*$/.test(path)) {
      return `.${path}`
    } else if (/^[^\/].*$/.test(path)) {
      return `./${path}`
    } else {
      return './'
    }
  }
}

// 分片上传
export class ChunkController extends BaseController {
  name = 'chunk'

  index: number = 0
  source: CancelTokenSource
  sliceCount: number
  query: ChunkUploadQuery
  chunkSize = UPLOAD_CHUNK_SIZE

  async setUploadId(): Promise<void> {
    const {
      data: { upload_id },
    } = await this.httpAdapter.post(
      `${this.preUploadUrl}?${qs.stringify({ ...this.query })}`
    )
    this.query.upload_id = upload_id
  }

  isPaused = false

  constructor({ chunkSize = UPLOAD_CHUNK_SIZE, ...props }: Props) {
    super(props)

    this.chunkSize = chunkSize
  }

  setReq(req) {
    super.setReq(req)
    this.sliceCount = Math.ceil(this.req.file.size / this.chunkSize)
  }

  async upload() {
    const { file, data } = this.req

    this.query = {
      ...data,
    }

    const filePath = (file as any).webkitRelativePath || file.name
    this.query.path = [...data.dir.split('/'), ...filePath.split('/')]
      .filter(item => !!item)
      .join('/')

    this.query.path = formatPath(this.query.path)
    this.query.file_size = file.size

    await this.setUploadId()

    this.uploadCurrentChunk()
  }

  private uploadCurrentChunk() {
    const {
      action,
      file,
      headers,
      withCredentials,
      onSuccess,
      onError,
      onProgress,
    } = this.req
    const currentFileSlice = this.req.file.slice(
      this.chunkSize * this.index,
      this.chunkSize * this.index + this.chunkSize
    )
    const isLastChunk = !this.sliceCount || this.index === this.sliceCount - 1
    const query = {
      ...this.query,
      offset: this.chunkSize * this.index,
      slice_size: isLastChunk
        ? file.size - this.chunkSize * this.index
        : this.chunkSize,
      finish: isLastChunk,
    }
    const url = `${action}?${qs.stringify(query)}`
    const formData = new FormData()
    formData.append('slice', currentFileSlice)
    this.source = axios.CancelToken.source()
    this.httpAdapter
      .post(url, formData, {
        headers,
        withCredentials,
        onUploadProgress: (e: any) => {
          const percent = Math.min(
            ((this.index * this.chunkSize + e.loaded) / file.size) * 100,
            100
          )
          onProgress({ percent }, file)
        },
        cancelToken: this.source.token,
      })
      .then(res => {
        this.index += 1
        if (!this.sliceCount || this.index === this.sliceCount) {
          onSuccess(res, file)
        } else {
          if (!this.isPaused) {
            this.uploadCurrentChunk()
          }
        }
      })
      .catch(e => {
        if (!(e instanceof axios.Cancel)) {
          onError(e)
          message.error(`${file.name}上传失败`)
        }
      })
  }

  pause() {
    this.isPaused = true
    this.abort()
  }

  resume() {
    this.isPaused = false
    this.retry()
  }

  abort() {
    this.source.cancel('canceled by the user')
  }

  retry() {
    this.uploadCurrentChunk()
  }
}
