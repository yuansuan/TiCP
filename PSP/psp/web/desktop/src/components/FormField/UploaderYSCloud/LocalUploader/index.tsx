import * as React from 'react'
import { Upload, Dropdown, message } from 'antd'
import { Validator } from '@/utils'

import { UploadMenu } from './style'

interface IProps {
  beforeUpload?: (params?: any) => Promise<any>
  upload: (params: any, isDir: boolean) => void
  data?: any
}

export default class LocalUploader extends React.Component<IProps> {
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

  uploadFile = props => {
    this.props.upload(props, false)
  }

  uploadDir = props => {
    this.props.upload(props, true)
  }

  render() {
    const { data, children } = this.props

    return (
      <Dropdown
        placement='bottomCenter'
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
              beforeUpload={this.beforeUpload}
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
