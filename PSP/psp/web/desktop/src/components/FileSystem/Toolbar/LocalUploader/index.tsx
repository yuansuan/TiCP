import * as React from 'react'
import { Upload, Dropdown, message } from 'antd'
import { observer } from 'mobx-react'
import { filter } from 'rxjs/operators'
import { computed } from 'mobx'

import { Button, Modal } from '@/components'
import { Point } from '@/domain/FileSystem'
import { uploader } from '@/domain'
import { Validator } from '@/utils'
import { untilDestroyed } from '@/utils/operators'
import { StyledUploadMenu } from './style'

const MAX_UPLOAD_FILE_NUMBER = 100
interface IProps {
  path: string
  point: Point
  disabled: boolean
  showLoading?: any
  hideLoading?: any
}

@observer
export default class UploadButton extends React.Component<IProps> {
  @computed
  get parentNode() {
    const { point, path } = this.props
    return point.filterFirstNode(item => item.path === path)
  }

  private uploadFile = props => {
    this.upload(props, false)
  }

  private upload = (props, isDir) => {
    const task = uploader.upload(props)
    const { dirPath } = props.data

    task.status$
      .pipe(
        untilDestroyed(this),
        filter(status => status === 'done')
      )
      .subscribe(() => {
        this.props.point.service.fetch(dirPath)
      })
  }

  // 将 已存在的文件 存储 dir 中
  existFilesInDir = []
  // 收集 dir 中，准备上传的文件数据
  uploadFileDatasInDir = []
  // 是否可以开始长传文件夹
  startUploadDir = false
  // 拒绝上传文件
  rejectUploadDir = false

  private uploadDir = props => {
    this.uploadFileDatasInDir.push(props)

    if (this.startUploadDir) {
      if (this.existFilesInDir.length !== 0) {
        Modal.showConfirm({
          content: `${this.existFilesInDir.join(',')} 已存在，是否覆盖？`
        })
          .then(() => {
            this.uploadFileDatasInDir.forEach(p => this.upload(p, true))
          })
          .catch(() => {
            const { path: basePath } = this.props
            this.uploadFileDatasInDir
              .filter(
                p =>
                  !this.existFilesInDir.includes(
                    `${basePath}/${p.file.webkitRelativePath || p.file.name}`
                  )
              )
              .forEach(p => this.upload(p, true))
          })
          .finally(() => {
            // clear
            this.startUploadDir = false
            this.existFilesInDir = []
            this.uploadFileDatasInDir = []
          })
      } else {
        this.uploadFileDatasInDir.forEach(p => this.upload(p, true))
        // clear
        this.startUploadDir = false
        this.existFilesInDir = []
        this.uploadFileDatasInDir = []
      }
    }
  }

  // ignore hidden file
  private isHiddenFile = fileName => fileName.startsWith('.')

  // check File Name
  private checkFileName = fileName => {
    const { error } = Validator.filename(fileName)
    if (error) {
      message.error(error.message)
      return false
    } else {
      return true
    }
  }

  // check file exist
  isExitedFile = async filePath => {
    const { point } = this.props

    const existArr = await point.service.exist([filePath])
    return existArr[0]
  }

  private beforeUploadDir = async (file, fileList) => {
    const { path: basePath, showLoading, hideLoading } = this.props

    showLoading && showLoading('正在检查目录中的文件，请稍等 ...')

    if (this.isHiddenFile(file.name)) {
      return Promise.reject()
    }

    if (!this.checkFileName(file.name)) {
      return Promise.reject()
    }

    let fileListTmp = fileList.filter(file => !file.name.startsWith('.'))
    if (this.rejectUploadDir) {
      if (fileListTmp[fileListTmp.length - 1].name === file.name) {
        hideLoading && hideLoading()
        this.rejectUploadDir = false
      }
      return Promise.reject()
    }
    if (fileListTmp.length > MAX_UPLOAD_FILE_NUMBER) {
      Modal.showConfirm({
        title: '确认弹窗',
        content: '当选择目录中文件数大于100时，请将文件打包上传',
        CancelButton: null
      })
      this.rejectUploadDir = true
      return Promise.reject()
    }

    if (fileListTmp[fileListTmp.length - 1].name === file.name) {
      // beforeUpload last file
      // check exist file
      const pathes = fileListTmp.map(
        file => `${basePath}/${file.webkitRelativePath || file.name}`
      )

      await Promise.all(
        pathes.map(async path => {
          // TODO 批量检测文件是否存在，会一致pending，需要调研 GRPC 接口，以下是临时方案
          const isExist = await this.isExitedFile(path)
          isExist && this.existFilesInDir.push(path)
        })
      )

      hideLoading && hideLoading()
      this.startUploadDir = true
    }
    return Promise.resolve()
  }

  private beforeUpload = async file => {
    const { path: basePath } = this.props

    if (this.isHiddenFile(file.name)) {
      return Promise.reject()
    }

    if (!this.checkFileName(file.name)) {
      return Promise.reject()
    }

    // check exist file
    const path = `${basePath}/${file.webkitRelativePath || file.name}`
    const isExist = await this.isExitedFile(path)

    if (isExist) {
      return Modal.showConfirm({
        content: `${path} 已存在，是否覆盖？`
      })
    } else {
      return Promise.resolve()
    }
  }

  render() {
    const { path, disabled } = this.props

    return (
      <Dropdown
        disabled={disabled}
        overlay={
          <StyledUploadMenu>
            <Upload
              action='/api/v3/file/upload'
              multiple={true}
              showUploadList={false}
              data={{ dirPath: path }}
              beforeUpload={this.beforeUpload}
              customRequest={this.uploadFile}>
              <span className='uploadItem first'>上传文件</span>
            </Upload>

            <Upload
              action='/api/v3/file/upload'
              multiple={true}
              directory={true}
              showUploadList={false}
              data={{ dirPath: path }}
              beforeUpload={this.beforeUploadDir}
              customRequest={this.uploadDir}>
              <span className='uploadItem'>上传目录</span>
            </Upload>
          </StyledUploadMenu>
        }>
        <Button type='primary' icon='upload' ghost>
          上传
        </Button>
      </Dropdown>
    )
  }
}
