/* Copyright (C) 2016-present, Yuansuan.cn */
import qs from 'qs'
import { UPLOAD_CHUNK_SIZE } from '@/constant'
import Uploader from '@/components/Uploader'
import { UploaderFile } from '@/components/Uploader/UploaderFile'
import { UploadChangeParam } from 'antd/es/upload'
import axios, { AxiosInstance, CancelTokenSource } from 'axios'
import { ICustomReq } from '@/components/Uploader/@types'
import { BaseController } from '@/components/Uploader/controllers/BaseController'
import { formatPath } from '@/components/Uploader/controllers/ChunkController'
import { message } from 'antd'
import AsyncPool from '../asyncpool/AsyncPool'
import { Http } from '@/utils'
import { uploader } from '@/domain'
import { v4 as uuid } from 'uuid'

type ControllerProps = {
  httpAdapter: AxiosInstance
  chunkSize?: number
  preUploadUrl?: string
}

interface ChunkUploadQuery {
  upload_id: string
  path: string
  file_size: number
  offset: number
  slice_size: number
  finish: boolean
  bucket: string
  tempDirPath: string
  tempUpload: string
  cross: boolean
  is_cloud: boolean
}

/**
 * 区块的上传进度计算器
 */
class ChunkPercentCalculator {
  /**
   * 已上传的区块大小
   */
  private uploaded: number = 0

  /**
   * 记录每个区块的进度
   */
  private chunkProgress: number[]

  constructor(private total: number, chunkCount: number) {
    this.chunkProgress = Array(chunkCount).fill(0)
  }

  /**
   * 更新区块的大小
   */
  updateChunkProgress(index: number, loaded: number) {
    this.uploaded += Math.max(loaded - this.chunkProgress[index], 0)
    this.chunkProgress[index] = Math.max(loaded, this.chunkProgress[index])
  }

  /**
   * 计算区块上传进度
   */
  get percent() {
    return Math.min(Math.ceil((this.uploaded / this.total) * 100), 100)
  }
}

/**
 * 全局上传任务协程池, 用于限制全局上传文件的并发数量
 *
 * 默认数量为 10
 */
let chunkAsyncPool: AsyncPool

/**
 * 全局上传前置任务协程池, 用于限制同时执行 preupload 的数量
 *
 * 默认数量为 10
 */
const prepareAsyncPool: AsyncPool = new AsyncPool()

/**
 * 异步上传控制器
 */
class AsyncChunkController extends BaseController {
  public name: string = 'AsyncChunk'

  /**
   * 分片大小
   */
  private chunkSize: number = UPLOAD_CHUNK_SIZE

  /**
   * 分片数量
   */
  private chunkCount: number = 0

  /**
   * 用于控制上传暂停或恢复
   */
  private tokenSources: { [key: string]: CancelTokenSource }

  /**
   * 当前上传是否暂停
   */
  private isPaused: boolean = false

  /**
   * 区块的上传参数
   */
  private query: ChunkUploadQuery

  /**
   * 已完成上传的区块
   */
  private uploadedParts: { [id: number]: boolean }

  /**
   * 区块上传进度计算器
   */
  private progress: ChunkPercentCalculator

  constructor(props: ControllerProps) {
    super(props)

    if (props.chunkSize) {
      this.chunkSize = props.chunkSize
    }

    this.tokenSources = {}
    this.uploadedParts = {}
  }

  /**
   * 保存请求对象并计算区块数量
   */
  setReq(req: ICustomReq) {
    super.setReq(req)
    this.chunkCount = Math.ceil(this.req.file.size / this.chunkSize)
    this.progress = new ChunkPercentCalculator(req.file.size, this.chunkCount)
  }

  /**
   * 获取单个上传任务的ID
   */
  async setUploadID(): Promise<void> {
    const {
      data: { upload_id }
    } = await this.httpAdapter.post(
      `${this.preUploadUrl}?${qs.stringify({ ...this.query })}`
    )

    this.query.upload_id = upload_id
  }

  /**
   * 执行上传任务
   *  1. 初始化上传参数
   *  2. 开启并发队列并执行上传
   *  3. 所有区块完成后标记文件上传完成
   */
  async upload() {
    const { file, data, onError } = this.req
    // 检测是初次上传还是暂停之后再次上传，需要保持相同的upload_id
    if (!this.query || !this.query.upload_id) {
      this.query = { ...data }
      const filePath = (file as any).webkitRelativePath || file.name
      this.query.path = [...data.dir?.split('/'), ...filePath?.split('/')]
      .filter(item => !!item)
      .join('/')
      
      this.query.path = formatPath(this.query.path)
      this.query.file_size = file.size

      if (this.query.tempDirPath) {
        let tempPath = formatPath(this.query.path)
        this.query.path = tempPath.replace(/^\./, this.query.tempDirPath) // 匹配以 ./ 开头的字符串并替换为tempDirPath
      }

      // 获取每次上传需要的 uploadID
      // 限制准备上传接口的并发数量
      await prepareAsyncPool.once(async () => {
        return await this.setUploadID()
      })
    }

    try {
      // 由于服务端是在第一个请求时创建的文件
      // 为了保证后续并发没问题，所以这里先发送一个
      // 空分片让服务端先把文件创建出来
      // await this.prepareAsyncQueue()
      await this.queueChunks()
      await this.finishChunk()

      // 防止出现协程泄漏，在完成上传之后，取消所有的任务
      this.abort()
    } catch (error: any) {
      if (!(error instanceof axios.Cancel)) {
        onError(error)
        // message.error(`${file.name}上传失败`)
      }
    }
  }

  /**
   * 暂停上传任务
   */
  pause() {
    this.isPaused = true
    this.abort()
  }

  /**
   * 恢复上传任务
   */
  resume() {
    this.isPaused = false
    this.retry()
  }

  /**
   * 发出取消信号
   */
  abort() {
    Object.keys(this.tokenSources).forEach(k => {
      if (this.tokenSources[k]) {
        this.tokenSources[k].cancel('canceled by user')
        this.tokenSources[k] = undefined
      }
    })
    this.tokenSources = {} // cleanup
  }

  /**
   * 重试当前上传任务
   */
  retry() {
    this.upload()
  }

  /**
   * 并发上传所有的区块
   */
  private async queueChunks() {
    let index = 0
    await chunkAsyncPool.run({
      next: () => (index === this.chunkCount ? null : index++),
      work: this.doUploadChunk.bind(this)
    })
  }

  /**
   * 准备一个空分片让服务端先创建文件
   */
  private async prepareAsyncQueue() {
    // 只有在未完成任意分片的情况下才可以让服务端创建相应的文件
    // 否则会导致已上传分片被服务端truncate成空数据
    if (Object.keys(this.uploadedParts).length === 0) {
      // empty chunk to truncate remote file
      const { action, formdata, tokenSource } = this.buildDataFromChunk(
        this.req.action,
        0,
        new Blob()
      )
      await this.httpAdapter.post(action, formdata, {
        headers: this.req.headers,
        withCredentials: this.req.withCredentials,
        cancelToken: tokenSource.token
      })
    }
  }

  /**
   * 表示文件以上传完成
   */
  private async finishChunk() {
    // empty chunk to finished upload
    const { action, formdata, tokenSource } = this.buildDataFromChunk(
      this.req.action,
      this.req.file.size,
      new Blob()
    )

    const resp = await this.httpAdapter.post(action, formdata, {
      headers: this.req.headers,
      withCredentials: this.req.withCredentials,
      cancelToken: tokenSource.token
    })

    this.req.onSuccess(resp, this.req.file)
  }

  /**
   * 并发上传每个区块
   */
  private async doUploadChunk(index: number) {
    if (!this.uploadedParts[index] && !this.isPaused) {
      const { file, onProgress } = this.req

      const { action, formdata, tokenSource } = this.buildChunkRequest(index)
      this.tokenSources[index] = tokenSource
      await this.httpAdapter.post(action, formdata, {
        headers: this.req.headers,
        withCredentials: this.req.withCredentials,
        onUploadProgress: (e: any) => {
          this.progress.updateChunkProgress(index, e.loaded)
          onProgress({ percent: this.progress.percent }, file)
        },
        cancelToken: tokenSource.token
      })

      delete this.tokenSources[index]
      this.uploadedParts[index] = true
    }
    return index
  }

  /**
   * 根据区块ID构建相对应的请求
   */
  private buildChunkRequest(index: number) {
    const { action, file } = this.req

    const offset = this.chunkSize * index
    const chunk = file.slice(offset, offset + this.chunkSize)
    return this.buildDataFromChunk(action, offset, chunk)
  }

  /**
   * 从区块构建一个请求参数
   */
  private buildDataFromChunk(action: string, offset: number, chunk: Blob) {
    const formdata = new FormData()
    formdata.append('slice', chunk)

    const query = {
      ...this.query,
      offset: offset,
      slice_size: chunk.size,
      finish: (chunk.size === 0 && offset !== 0) || this.query.file_size === 0
    }
    // set custom params
    Object.keys(query).forEach(key => {
      formdata.append(key, query[key])
    })

    return {
      action: `${action}`,
      formdata: formdata,
      tokenSource: axios.CancelToken.source()
    }
  }
}
/**
 * 劫持文件上传的控制器，增加异步上传特性
 */
export const hijackUploaderController = (
  uploader: Uploader,
  concurrency?: number
) => {
  if (uploader['__hijacked__'] == undefined) {
    uploader['__hijacked__'] = true

    chunkAsyncPool = new AsyncPool({ concurrency: concurrency })
    const original = uploader['setControllerMap'].bind(uploader)

    uploader['setControllerMap'] = function (this: Uploader, uid: string) {
      const { by, httpAdapter, chunkSize } = this.props
      if (by === 'chunk') {
        const controller = new AsyncChunkController({
          httpAdapter,
          chunkSize,
          preUploadUrl: this.config.preUploadUrl
        })
        this.controllerMap.set(uid, controller)
        return controller
      }

      return original(uid)
    }.bind(uploader)
    
    uploader['globalOnChangeHandler'] = function (
      data: UploadChangeParam & { origin: string }
      ) {
      this.antFileList = data.fileList
      const uploadFile = data.file
      if (!uploadFile) return
      const index = this.fileList.findIndex(item => item.uid === uploadFile.uid)
      if (index < 0) {
        console.log('上传文件初始化中...')
        const ctrl = this.setControllerMap(data.file.uid)

        this.fileList.push(
          new UploaderFile({
            origin: data.origin,
            ...data.file,
            by: ctrl.name
          })
        )
        // requestIdleCallback(
        //   () => {
        //     this.fileList.push(
        //       new UploaderFile({
        //         origin: data.origin,
        //         ...data.file,
        //         by: ctrl.name
        //       })
        //     )
        //   },
        //   { timeout: 500 }
        // )
      } else {
        console.log('上传文件状态更新中...')
        this.fileList[index].refresh(data.file)
        // requestIdleCallback(
        //   () => {
        //     this.fileList[index].refresh(data.file)
        //   },
        //   { timeout: 1000 }
        // )
      }
    }.bind(uploader)

    console.log(
      'Hijacking the upload controller for async/concurrent controller'
    )
  }
}