/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import ReactDOM from 'react-dom'
import { Upload, message } from 'antd'
import {
  UploadProps as AntUploadProps,
  UploadChangeParam,
  RcFile
} from 'antd/es/upload'
import { UploadFile } from 'antd/es/upload/interface'
import { Observer } from 'mobx-react-lite'
import { action, observable, runInAction } from 'mobx'
import { UploaderFile } from './UploaderFile'
import { ICustomReq } from './@types'
import { FullController } from './controllers/FullController'
import { ChunkController } from './controllers/ChunkController'
import { BaseController } from './controllers/BaseController'
import axios, { AxiosInstance } from 'axios'
import { nextTick, single } from '@/utils'
import { UPLOAD_CHUNK_SIZE } from './constant'

export * from './UploaderFile'

export type UploadProps = Omit<AntUploadProps, 'beforeUpload' | 'onChange'> & {
  origin: string
  // 上传方式：分片上传/全量上传
  by?: 'full' | 'chunk'
  // 分片大小（仅在 by 属性为 chunk 时生效）
  chunkSize?: number
  httpAdapter?: AxiosInstance
  beforeUpload?: (fileList: RcFile[]) => string[] | Promise<string[]>
  onChange?: (
    info: UploadChangeParam<UploadFile & { origin: string }> & {
      origin: string
    }
  ) => void
}

export type UploaderProps = {
  // 分片大小（仅在 by 属性为 chunk 时生效）
  chunkSize?: number
  preUploadUrl?: string
}

export default class Uploader {
  ref
  uploader
  controllerMap = new Map<string, BaseController>()

  @observable config: UploaderProps = {
    chunkSize: UPLOAD_CHUNK_SIZE,
    // '/filemanager/pre-upload' 盒子统一接口 允许自定义
    preUploadUrl: '/filemanager/pre-upload'
  }
  @observable props: Omit<UploadProps, 'origin'> = {}
  @observable fileList: UploaderFile[] = []
  @observable antFileList: any[]

  private thisTickUploadFiles = []
  private thisTickBeforeloadRes = null

  constructor(config?: Partial<UploaderProps>) {
    if (config) {
      this.updateConfig(config)
    }

    const el = document.createElement('div')
    el.style.display = 'none'
    document.querySelector('body').appendChild(el)
    ReactDOM.render(
      <Observer>
        {() => (
          <Upload
            {...(this.props as any)}
            beforeUpload={this.beforeUpload}
            showUploadList={false}
            customRequest={this.customRequest}
            fileList={this.antFileList}
            ref={uploader => (this.uploader = uploader)}>
            <div ref={ref => (this.ref = ref)} />
          </Upload>
        )}
      </Observer>,
      el
    )
  }

  @action
  updateConfig = (config: Partial<UploaderProps>) => {
    Object.assign(this.config, config)
  }

  beforeUpload = (file: RcFile) => {
    this.thisTickUploadFiles.push(file)
    const { beforeUpload } = this.props

    return new Promise<void>((resolve, reject) => {
      const handler = (data: string[]) => {
        if (data.includes(file.uid)) {
          resolve()
        } else {
          reject()
        }
        nextTick(() => {
          this.thisTickUploadFiles = []
          this.thisTickBeforeloadRes = null
        })
      }
      nextTick(() => {
        // 保证只执行一次beforeUpload, 防止多次弹窗
        if (!this.thisTickBeforeloadRes) {
          this.thisTickBeforeloadRes = beforeUpload
            ? beforeUpload(this.thisTickUploadFiles)
            : (this.thisTickUploadFiles || []).map(item => item.uid)
        }
        if (this.thisTickBeforeloadRes instanceof Promise) {
          this.thisTickBeforeloadRes
            .then(data => {
              handler(data)
            })
            // let the promise be fullfilled when user have canceled the upload
            .catch(() => handler([]))
        } else {
          handler(this.thisTickBeforeloadRes)
        }
      })
    })
  }

  upload({ origin, ...props }: UploadProps) {
    this.props = props
    // manage upload data
    this.props.onChange = data => {
      this.globalOnChangeHandler({ ...data, origin })
      const uploaderFile = this.fileList.find(
        item => item.uid === data.file.uid
      )
      props.onChange({ ...data, origin: uploaderFile?.origin })

      if (data.file.status === 'done') {
        single(
          'upload-success-message',
          () => message.success('文件上传成功').promise
        )
      }
    }
    // trigger Uploader click event to select upload file
    setTimeout(() => {
      this.ref.click()
    }, 0)
  }

  remove(id: string) {
    const index = this.fileList.findIndex(f => f.uid === id)
    if (index < 0) {
      return
    }
    const file = this.fileList[index]
    switch (file.status) {
      case 'uploading':
        const controller = this.controllerMap.get(id)
        if (controller) {
          controller.abort()
          this.controllerMap.delete(id)
        }
        break
      default:
        break
    }

    file.status = 'removed'
    this.props.onRemove && this.props.onRemove(file as any)
    this.fileList.splice(index, 1)
  }

  retry(id: string) {
    // reset ant upload file status
    // will trigger onChange event
    const index = this.antFileList.findIndex(f => f.uid === id)
    if (index < 0) {
      return
    }

    const file = this.antFileList[index]
    if (file.status !== 'error') {
      return
    }

    file.status = 'uploading'
    file.error = null
    file.response = null
    const controller = this.controllerMap.get(id)
    controller.retry()
  }

  async pause(id: string) {
    const index = this.fileList.findIndex(f => f.uid === id)
    if (index < 0) return
    const file = this.fileList[index]
    const controller = this.controllerMap.get(id)
    await controller.pause()
    file.status = 'paused'
  }

  async resume(id: string) {
    const index = this.fileList.findIndex(f => f.uid === id)
    if (index < 0) return
    const file = this.fileList[index]

    const controller = this.controllerMap.get(id)
    await controller.resume()
    file.status = 'uploading'
  }

  private customRequest = async (req: ICustomReq) => {
    const controller = this.controllerMap.get((req.file as any).uid)
    controller.setReq(req)
    await controller.upload()
  }

  private globalOnChangeHandler = (
    data: UploadChangeParam & { origin: string }
  ) => {
    runInAction(() => {
      this.antFileList = data.fileList
      const uploadFile = data.file
      if (!uploadFile) return
      const index = this.fileList.findIndex(item => item.uid === uploadFile.uid)
      if (index < 0) {
        const ctrl = this.setControllerMap(data.file.uid)
        this.fileList.push(
          new UploaderFile({ origin: data.origin, ...data.file, by: ctrl.name })
        )
      } else {
        this.fileList[index].refresh(data.file)
      }
    })
  }

  private setControllerMap(uid) {
    const defaultHttpAdapter = axios.create()
    const { by, httpAdapter = defaultHttpAdapter, chunkSize } = this.props
    const controller =
      by === 'chunk'
        ? new ChunkController({
            httpAdapter,
            chunkSize: chunkSize || this.config.chunkSize,
            preUploadUrl: this.config.preUploadUrl
          })
        : new FullController({ httpAdapter })
    this.controllerMap.set(uid, controller)
    return controller
  }
}
