import * as React from 'react'
import { Upload, Dropdown, message } from 'antd'
import { Validator } from '@/utils'

import { UploadMenu } from './style'

interface IProps {
  beforeUpload?: (params?: any) => Promise<any>
  upload: (params: any, isDir: boolean) => void
  data?: any
  disabled?: boolean
}

export default class LocalUploader extends React.Component<IProps> {
  static uploadFileParamsForDir = []
  static curruntUploadDirLastFile = null
  private beforeUpload = file => {
    // ignore hidden file
    if (file.name.startsWith('.')) {
      return Promise.reject()
    }

    // check filename
    const { error } = Validator.filename(file.name)
    if (error) {
      message.error(error.message)
      return Promise.reject(error)
    }

    // props.beforeUpload
    const { beforeUpload } = this.props
    if (beforeUpload) {
      return beforeUpload({
        file,
        data: this.props.data,
      })
    }

    return Promise.resolve()
  }

  private beforeUploadDir = (file, fileList) => {
    // props.beforeUpload
    const { beforeUpload } = this.props

    let lastFile = fileList.at(-1)
    LocalUploader.curruntUploadDirLastFile = fileList.at(-1)

    if (beforeUpload) {
      if (file.uid === lastFile.uid) {
        return beforeUpload({
          file,
          data: this.props.data,
          isDir: true,
          isLast: true,
        })
      }
    }

    // ignore hidden file
    if (file.name.startsWith('.')) {
      return Promise.reject()
    }

    // check filename
    const { error } = Validator.filename(file.name)
    if (error) {
      message.error(error.message)
      return Promise.reject(error)
    }

    return Promise.resolve()
  }

  uploadFile = props => {
    this.props.upload(props, false)
  }

  uploadDir = props => {
    // 收集要上传的参数，等待信号统一上传
    // this.props.upload(props, true)
    LocalUploader.uploadFileParamsForDir.push(props)

    if (LocalUploader.curruntUploadDirLastFile?.uid === props.file?.uid) {
      LocalUploader.uploadFileParamsForDir.forEach(params => {
        this.props.upload(params, true)
      })
      LocalUploader.uploadFileParamsForDir = []
      LocalUploader.curruntUploadDirLastFile = null
    } 
  }

  render() {
    const { data, children, disabled } = this.props

    return (
      <Dropdown
        placement='bottomCenter'
        disabled={disabled}
        overlay={
          <UploadMenu>
            <Upload
              className='upload'
              action='/api/v3/file/upload'
              multiple={true}
              // hack: ignore set state in unmounted component warning
              fileList={[]}
              data={data}
              showUploadList={false}
              beforeUpload={this.beforeUpload}
              customRequest={this.uploadFile}>
              <span className='uploadItem'>上传文件</span>
            </Upload>
            <Upload
              className='upload'
              action='/api/v3/file/upload'
              multiple={true}
              // hack: ignore set state in unmounted component warning
              fileList={[]}
              data={data}
              directory={true}
              showUploadList={false}
              beforeUpload={this.beforeUploadDir}
              customRequest={this.uploadDir}>
              <span className='uploadItem'>上传目录</span>
            </Upload>
          </UploadMenu>
        }>
        {children}
      </Dropdown>
    )
  }
}
