/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, computed } from 'mobx'
import { formatByte } from '@/utils/Validator'

export declare type UploadFileStatus =
  | 'error'
  | 'success'
  | 'done'
  | 'uploading'
  | 'removed'
  | 'paused'

export class UploaderFileProps {
  @observable origin: string
  @observable uid: string
  @observable size: number
  @observable name: string
  @observable fileName?: string
  @observable lastModified?: number
  @observable lastModifiedDate?: Date
  @observable url?: string
  @observable status?: UploadFileStatus
  @observable percent?: number
  @observable thumbUrl?: string
  @observable originFileObj?: File | Blob
  @observable response?: any
  @observable error?: any
  @observable linkProps?: any
  @observable type: string
  @observable speed?: number
  @observable by?: string
}

export class UploaderFile extends UploaderFileProps {
  private lastFilePercent = 0

  private speedArr = []
  private speedArrCount = 5 // 计算五次平均速度

  constructor(props: UploaderFileProps) {
    super()
    Object.assign(this, props)

    setInterval(() => {
      const speed = (this.size * (this.percent - this.lastFilePercent)) / 100
      this.lastFilePercent = this.percent
      this.speedArr = this.speedArr.concat(speed).slice(-this.speedArrCount)
      this.speed =
        this.speedArr.reduce((prev, curr) => prev + curr, 0) /
        this.speedArr.length
    }, 1000)
  }

  refresh(props: Omit<UploaderFileProps, 'origin'>) {
    Object.assign(this, props)
  }

  get statusText() {
    const map = {
      error: '上传失败',
      success: '上传成功',
      done: '上传完成',
      uploading: '上传中',
      removed: '已删除',
      paused: '已暂停',
    }
    return map[this.status] || '未知状态'
  }

  @computed
  get speedText() {
    return `${formatByte(this.speed)}/s`
  }

  @computed
  get loaded() {
    return (this.size * this.percent) / 100
  }

  @computed
  get totalSizeText() {
    return `${formatByte(this.size)}`
  }

  @computed
  get loadedSizeText() {
    return `${formatByte(this.loaded)}`
  }
}
