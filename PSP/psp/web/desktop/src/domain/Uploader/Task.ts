/**
 * @module Task
 * @description task progress/task operators
 */
import nanoid from 'nanoid'
import { observable, action, computed } from 'mobx'
import { Subject, merge } from 'rxjs'
import { takeUntil, filter, throttleTime } from 'rxjs/operators'
import SparkMD5 from 'spark-md5'

import axios from 'axios'

import { createMobxStore, Fetch } from '@/utils'
import { currentUser } from '@/domain'
import UploadError, { MergeAndMd5Error } from './UploadError'
import { eventEmitter, IEventData } from '@/utils'
import Worker from '@/worker/md5.worker'

export enum TaskStatus {
  inited = 'inited',
  uploading = 'uploading',
  paused = 'paused',
  pausing = 'pausing',
  error = 'error',
  done = 'done',
  aborted = 'aborted',
  mergeAndMd5 = 'mergeAndMd5'
}

export const statusMap = {
  inited: '排队中',
  uploading: '上传中',
  paused: '已暂停',
  pausing: '暂停中',
  error: '上传失败',
  done: '已完成',
  aborted: '已取消',
  mergeAndMd5: '合并与校验中'
}

interface ITask {
  id: string
  speed: number
  loaded: number
  chunkTotal: number
  status: TaskStatus
  target: any
  cursor: number
  result: any
}

const SIZE_POINTS = {
  '1G': 1 * 1024 * 1024 * 1024,
  '3G': 3 * 1024 * 1024 * 1024,
  '6G': 6 * 1024 * 1024 * 1024,
  '10G': 10 * 1024 * 1024 * 1024
}

let CHUNK_SIZE = 3 * 1024 * 1024
//  1G 合并分片同步与异步方式分割点
let SYNC_SIZE = 1024 * 1024 * 1024

const computeChuckSize = total => {
  if (total <= SIZE_POINTS['1G']) {
    return 2 * 1024 * 1024 // 2M
  } else if (total <= SIZE_POINTS['3G']) {
    return 3 * 1024 * 1024 // 3M
  } else if (total <= SIZE_POINTS['6G']) {
    return 5 * 1024 * 1024 // 5M
  } else if (total <= SIZE_POINTS['10G']) {
    return 6 * 1024 * 1024 // 6M
  } else {
    return 10 * 1024 * 1024 // 10M
  }
}

const BIG_FILE_SIZE_CHUNK_HUMBER = 10

export default class Task implements ITask {
  readonly id = nanoid()
  chunkTotal = 0
  action = ''
  params = {}
  lastLoaded = 0
  prevTime = 0
  axiosSource = null
  dirPath = ''

  // 首次
  progressStartMark = true

  abort$ = new Subject()
  pause$ = new Subject()
  status$ = null
  updateSpeed$ = new Subject<number>()

  @observable isDir = false
  @observable speed = 0
  @observable loaded = 0
  @observable status = TaskStatus.inited
  @observable cursor = 0
  @observable target = null
  @observable result = null
  @observable error = null

  concurrent = 3
  subLoaded = []
  subResults = []
  subErrors = []
  subRetryTimes = []
  subCursor = []
  subPaused = []

  md5 = null

  // The retryTimes identify the remain times to retry
  retryTimes = 3

  // 收集需要计算md5值的大文件
  waitingCheckMD5Files = new Map()

  @computed
  get percent() {
    const { target } = this
    if (target && !target.size) {
      return 100
    }
    return target ? (this.loaded / target.size) * 100 : 0
  }

  constructor({ file: target, action, params = {}, dirPath, isDir = false }) {
    this.axiosSource = axios.CancelToken.source()
    this.action = action
    this.params = params
    this.isDir = isDir
    this.dirPath = dirPath

    // CHUNK_SIZE = computeChuckSize(target.size)
    CHUNK_SIZE = 50 * 1024 * 1024 // 50M

    const chunkTotal = Math.max(Math.ceil(target.size / CHUNK_SIZE), 1)

    this.update({
      target,
      chunkTotal
    })

    this.subResults = new Array(chunkTotal).fill(null)
    this.subErrors = new Array(chunkTotal).fill(null)
    this.subRetryTimes = new Array(chunkTotal).fill(this.retryTimes)
    this.subLoaded = new Array(chunkTotal).fill(0)
    this.subPaused = new Array(this.concurrent).fill(false)
    this.subCursor = Array.from(new Array(this.concurrent), (_, i) => i)

    // 对于进行分块上传的文件, 也就是大文件, 进行异步计算 md5
    if (chunkTotal > BIG_FILE_SIZE_CHUNK_HUMBER) {
      // 计算上传文件的 md5
      this.computeMd5(target)
    }

    this.status$ = createMobxStore(() => this.status)

    // update speed
    this.updateSpeed$
      .pipe(
        takeUntil(
          merge(
            this.abort$,
            this.status$.pipe(filter(status => status === TaskStatus.done))
          )
        ),
        throttleTime(500)
      )
      .subscribe((speed: number) => this.update({ speed }))
  }

  @action
  update = (props: Partial<ITask>) => {
    Object.assign(this, props)
  }

  @computed
  get displayState() {
    return statusMap[this.status]
  }

  syncComputeMd5(file) {
    return new Promise((resolve, reject) => {
      const spark = new SparkMD5.ArrayBuffer()

      let blobSlice = File.prototype.slice,
        currentChunk = 0
      const fileReader = new FileReader()

      fileReader.onload = e => {
        spark.append(e.target.result) // Append array buffer
        currentChunk++

        if (currentChunk < this.chunkTotal) {
          loadNext()
        } else {
          resolve(spark.end())
        }
      }

      fileReader.onerror = function () {
        console.warn('md5 计算失败')
        reject('md5 计算失败')
      }

      function loadNext() {
        let start = currentChunk * CHUNK_SIZE,
          end = start + CHUNK_SIZE >= file.size ? file.size : start + CHUNK_SIZE

        fileReader.readAsArrayBuffer(blobSlice.call(file, start, end))
      }

      loadNext()
    })
  }

  // 校验md5
  checkMD5 = async formData => {
    try {
      await Fetch.post(this.action, formData, { timeout: 0 })
      eventEmitter.once(`FILE_UPLOAD_${this.id}`, (data: IEventData) => {
        const msg = data.message
        console.debug(`FILE_UPLOAD_${this.id}:`, msg)

        if (msg.success) {
          // 计算完之后给删除
          this.waitingCheckMD5Files.delete(this.id)
          this.update({
            status: TaskStatus.done
          })
        } else {
          this.error = new MergeAndMd5Error(msg.errCode)
          this.update({
            status: TaskStatus.error
          })
        }
      })

      // 超时处理
      setTimeout(() => {
        this.error = this.error || {
          message: '合并与校验超时，联系管理员'
        }
        this.update({
          status: TaskStatus.error
        })
        eventEmitter.off(`FILE_UPLOAD_${this.id}`)
      }, 1000 * 60 * 30)
    } catch (err) {
      this.error = err

      if (err.message === 'Network Error') {
        this.error.message = '网络异常'
      }

      if (err?.response?.status >= 500) {
        this.error.message = '合并与校验失败，请联系管理员'
      }

      this.update({
        status: TaskStatus.error
      })

      return
    }
  }

  computeMd5(file) {
    let _eventName = 'md5'

    const worker = new Worker()

    console.debug('send message: start compute md5')

    worker.postMessage({
      eventName: _eventName,
      eventData: {
        file,
        chunkTotal: this.chunkTotal,
        CHUNK_SIZE,
        uploadId: this.id
      }
    })

    worker.addEventListener('message', async event => {
      const { eventName, eventData } = event.data
      const { uploadId, md5 } = eventData

      if (eventName === _eventName) {
        if (uploadId === this.id) {
          this.md5 = md5
          console.debug('receive message: md5', md5)

          if ([this.waitingCheckMD5Files.values()].length === 0) return

          //  如果有对应的md5 算出来，就发起校验请求，确保完整性正确
          const formData = this.waitingCheckMD5Files.get(uploadId)
          formData.append('md5', md5)

          await this.checkMD5(formData)
        }

        worker.terminate()
      }
    })
  }

  start = async () => {
    const { status } = this
    const canStartStates = [
      TaskStatus.paused,
      TaskStatus.inited,
      TaskStatus.error
    ]
    // you can only start a paused or inited or error task
    if (!canStartStates.includes(status)) {
      return
    }

    // 恢复各个子任务暂停状态
    this.subPaused = new Array(this.concurrent).fill(false)
    // 恢复 CancelToken
    this.axiosSource = axios.CancelToken.source()

    // 没有获取md5值，并且是小文件
    if (!this.md5 && this.chunkTotal <= BIG_FILE_SIZE_CHUNK_HUMBER) {
      this.md5 = await this.syncComputeMd5(this.target)
    }

    // begin uploading
    this._upload()
  }

  private _upload = () => {
    // avoid duplicate upload
    if (this.status === TaskStatus.uploading) {
      return
    }

    this.update({
      status: TaskStatus.uploading
    })

    this.prevTime = Date.now()

    const uploadPatch = async cursor => {
      if (cursor > this.chunkTotal - 1) {
        return
      }

      if (
        this.status === TaskStatus.aborted ||
        this.status === TaskStatus.paused ||
        this.subPaused[cursor % this.concurrent]
      ) {
        return
      }

      // get chunk by cursor
      const chunk = this.target.slice(
        CHUNK_SIZE * cursor,
        CHUNK_SIZE * (cursor + 1)
      )

      try {
        const { target } = this
        const formData = new FormData()
        formData.append('user_id', currentUser.id)
        formData.append('upload_id', this.id)
        formData.append('path', target.uploadPath)
        formData.append('file_size', target.size)
        formData.append('slice_size', chunk.size)
        formData.append('offset', CHUNK_SIZE * cursor + '')
        formData.append('finish', '0')
        formData.append('slice', chunk, 'slice')
        formData.append('seq', cursor)
        // set custom params
        Object.keys(this.params).forEach(key => {
          formData.append(key, this.params[key])
        })

        let result

        try {
          const res = await Fetch.post(this.action, formData, {
            timeout: 0, // 永不超时
            cancelToken: this.axiosSource.token,
            disableErrorMessage: true,
            onUploadProgress: progressEvent => {
              // 如果出现错误, 不更新上传进度相关数据
              if (this.status === TaskStatus.error) {
                return
              }

              // omit error progress
              if (progressEvent.total <= 0) {
                return
              }
              // 处理已经上传了多少，上传速度
              const chunkPercent = Math.round(
                (progressEvent.loaded * 100) / progressEvent.total
              )

              this.subLoaded[cursor] = (chunk.size * chunkPercent) / 100

              const now = Date.now()
              // calculate speed per second
              if (
                now - this.prevTime >= 1000 ||
                this.progressStartMark === true
              ) {
                const speed =
                  ((this.loaded - this.lastLoaded) /
                    (Date.now() - this.prevTime)) *
                  1000
                this.updateSpeed$.next(speed >= 0 ? speed : 0) // 上传速度为大致估算
                this.prevTime = now
                this.lastLoaded = this.loaded
                this.progressStartMark = false
              }

              let totalLoaded = this.subLoaded.reduce((p, c) => p + c, 0)

              this.update({
                loaded: this.loaded >= totalLoaded ? this.loaded : totalLoaded
              })
            }
          })

          if (res && !res.success) {
            throw new UploadError(res.code)
          } else {
            result = true
            console.debug('分片上传成功', cursor)
          }
        } catch (err) {
          result = false

          if (err.message !== '用户暂停当前分片上传任务') {
            // 其它错误抛出
            throw err
          } else {
            // 人为暂停
            return
          }
        }

        if (this.status === TaskStatus.pausing) {
          this.subPaused[cursor % this.concurrent] = true
          if (this.subPaused.every(p => p)) {
            this.update({ status: TaskStatus.paused })
          }
        }

        this.subResults[cursor] = result

        // upload complete
        if (this.subResults.every(res => res)) {
          this.update({
            status: TaskStatus.mergeAndMd5
          })
          // 等待计算出来后，再进行发起md5校验请求
          if (!this.md5) {
            const tempMergeData = new FormData()
            tempMergeData.append('user_id', currentUser.id)
            tempMergeData.append('upload_id', this.id)
            tempMergeData.append('path', target.uploadPath)
            tempMergeData.append('file_size', target.size)
            tempMergeData.append('slice_size', chunk.size)
            tempMergeData.append('offset', CHUNK_SIZE * cursor + '')
            tempMergeData.append('finish', '1')
            tempMergeData.append('is_dir', this.isDir ? '1' : '0')
            tempMergeData.append('slice', chunk, 'slice')
            tempMergeData.append('seq', -1 + '')
            tempMergeData.append('total_slice', this.chunkTotal + '')
            this.waitingCheckMD5Files.set(this.id, tempMergeData)
          } else {
            const mergeData = new FormData()

            mergeData.append('user_id', currentUser.id)
            mergeData.append('upload_id', this.id)
            mergeData.append('path', target.uploadPath)
            mergeData.append('file_size', target.size)
            mergeData.append('slice_size', chunk.size)
            mergeData.append('offset', CHUNK_SIZE * cursor + '')
            mergeData.append('finish', '1')
            mergeData.append('is_dir', this.isDir ? '1' : '0')
            mergeData.append('slice', chunk, 'slice')
            mergeData.append('seq', -1 + '')
            mergeData.append('md5', this.md5)
            mergeData.append('total_slice', this.chunkTotal + '')
            // 小于等于 SYNC_SIZE 同步, 大于 SYNC_SIZE 异步
            const isSync = target.size <= SYNC_SIZE
            mergeData.append('sync', isSync ? '1' : '0')

            console.debug(target.uploadPath, this.md5)

            let res = null

            try {
              res = await Fetch.post(this.action, mergeData, { timeout: 0 })
              // 异步
              if (!isSync) {
                // 事件监听，文件合并成功或失败消息
                eventEmitter.once(
                  `FILE_UPLOAD_${this.id}`,
                  (data: IEventData) => {
                    const msg = data.message
                    console.debug(`FILE_UPLOAD_${this.id}:`, msg)

                    if (msg.success) {
                      this.update({
                        status: TaskStatus.done
                      })
                    } else {
                      this.error = new MergeAndMd5Error(msg.errCode)
                      this.update({
                        status: TaskStatus.error
                      })
                    }
                  }
                )

                // 超时处理
                setTimeout(() => {
                  this.error = this.error || {
                    message: '合并与校验超时，联系管理员'
                  }
                  this.update({
                    status: TaskStatus.error
                  })
                  eventEmitter.off(`FILE_UPLOAD_${this.id}`)
                }, 1000 * 60 * 30)
              } else {
                // 同步
                if (!res || res.success) {
                  this.update({
                    status: TaskStatus.done
                  })
                } else {
                  this.update({
                    status: TaskStatus.error
                  })
                  throw new MergeAndMd5Error(res.code)
                }
              }
              return
            } catch (err) {
              this.error = err

              if (err.message === 'Network Error') {
                this.error.message = '网络异常'
              }

              if (err?.response?.status >= 500) {
                this.error.message = '合并与校验失败，请联系管理员'
              }

              this.update({
                status: TaskStatus.error
              })

              return
            }
          }
        }

        this.subCursor[cursor % this.concurrent] = cursor + this.concurrent

        uploadPatch(cursor + this.concurrent)
      } catch (err) {
        // 处理网络异常，以及其它错误，重试机制
        this.subErrors[cursor] = err
        this.error = err

        if (err.message === 'Network Error') {
          this.error.message = '网络异常'
        }

        if (err?.response?.status >= 500) {
          this.error.message = '上传失败，请联系管理员'
        }

        if (this.subRetryTimes[cursor] <= 0 || err instanceof UploadError) {
          this.update({
            status: TaskStatus.error
          })
        } else {
          // minus one time
          this.subRetryTimes[cursor] -= 1

          // start upload patch
          uploadPatch(cursor)
        }
      }
    }

    // start upload patch
    Promise.all(this.subCursor.map(cursor => uploadPatch(cursor)))
  }

  abort = () => {
    this.update({
      status: TaskStatus.aborted
    })
    this.abort$.next()
    this.axiosSource.cancel('用户删除上传任务')
  }

  pause = () => {
    this.axiosSource.cancel('用户暂停当前分片上传任务')
    this.update({
      status: TaskStatus.paused // 去掉暂停中的设定, 直接暂停
    })
    this.pause$.next()
  }
}
