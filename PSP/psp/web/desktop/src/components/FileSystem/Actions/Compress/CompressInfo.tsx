import * as React from 'react'
import { Input, Select, message } from 'antd'
import { observer } from 'mobx-react'
import { observable, action } from 'mobx'

import { Button } from '@/components'
import { StyledCompressInfo } from './style'
import ChooseFileAction from '../ChooseFile'
import { RootPoint } from '@/domain/FileSystem'
import { ZipType } from '@/utils/const'

const Option = Select.Option

interface IProps {
  zipName: string
  onCancel: any
  onOk: any
  points: RootPoint[]
  parentPath: string
}

@observer
export default class CompressInfo extends React.Component<IProps> {
  @observable zipName = ''
  @observable compressType = 'Zip'
  @observable destPath = this.props.parentPath
  @action
  updateZipName = zipName => (this.zipName = zipName)
  @action
  updateCompressType = type => (this.compressType = type)

  zipNameRef = null

  constructor(props) {
    super(props)

    this.updateZipName(props.zipName)
  }

  componentDidMount() {
    this.zipNameRef && this.zipNameRef.select()
  }

  onCancel = () => {
    this.props.onCancel()
  }

  onOk = () => {
    const { zipName, compressType, destPath } = this

    if (!zipName) {
      message.error('请输入压缩包名称')
      return
    }

    if (/[\\/,;'`"]/.test(zipName)) {
      message.error(
        '压缩包名称中包含非法字符：斜杆、反斜杆、单引号、双引号、反引号、逗号或分号'
      )
      return
    }

    this.props.onOk({ zipName, compressType, destPath })
  }

  render() {
    const {
      zipName,
      compressType,
      updateZipName,
      updateCompressType,
      onCancel,
    } = this

    const { points } = this.props

    return (
      <StyledCompressInfo>
        <div className='body'>
          <div className='module'>
            <span className='name'>文件名：</span>
            <Input
              ref={ref => (this.zipNameRef = ref)}
              autoFocus
              value={zipName}
              maxLength={120}
              className='widget'
              onChange={e => {
                updateZipName(e.target.value)
              }}
            />
          </div>
          <div className='module'>
            <span className='name'>类型：</span>
            <Select
              className='widget'
              value={compressType}
              onChange={updateCompressType}>
              {Object.keys(ZipType).map(key => (
                <Option key={key} value={key}>
                  {ZipType[key]}
                </Option>
              ))}
            </Select>
          </div>
          <div className='module'>
            <span className='name'>压缩到：</span>
            <ChooseFileAction
              points={points}
              disabledKeys={['__unknown__dir__']}
              path={this.destPath}
              onOk={(keys: string[], next, currPoint) => {
                console.log(keys)
                this.destPath = keys[0]
                next()
              }}
              hasPerm={true}
              onCancel={(keys, next) => next()}>
              <Button type='link'>{this.destPath}</Button>
            </ChooseFileAction>
          </div>
        </div>
        <div className='footer'>
          <div className='footerMain'>
            <Button onClick={onCancel}>取消</Button>
            <Button type='primary' onClick={this.onOk}>
              确认
            </Button>
          </div>
        </div>
      </StyledCompressInfo>
    )
  }
}
