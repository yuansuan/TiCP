import * as React from 'react'
import { message, Spin } from 'antd'
import { Button } from '@/components'
import { observer } from 'mobx-react'
import { action, observable } from 'mobx'
import { StepWrapper } from './style'
import { CopyToClipboard } from 'react-copy-to-clipboard'
import license from '@/domain/License'

interface IProps {
  successCallback: () => void
  successMsg?: string
}

const initLicenseAttrs = {
  name: '',
  version: '',
  expiry: '',
  machine: '',
  key: ''
}
@observer
export default class LicenseForm extends React.Component<IProps> {
  @observable loading = false
  @observable machineId = ''
  @observable updating = false
  @observable content = ''
  @observable licenseAttribute = initLicenseAttrs

  @action
  updateContent = content => (this.content = content)

  private inputRef

  constructor(props) {
    super(props)
    this.inputRef = React.createRef()
  }

  async componentDidMount() {
    this.loading = true
    try {
      const res = await license.getMachineId()
      this.machineId = res.data?.id
    } finally {
      this.loading = false
    }
  }

  submitFile = async e => {
    const fileContent = await this.readFile(e.target.files[0])
    this.updateContent(fileContent)

    try {
      const yaml = require('js-yaml')
      const doc = yaml.load(this.content)
      if (this.validateFileContent(doc)) {
        ;['expiry', 'key', 'machine', 'name', 'version'].forEach(
          key => (this.licenseAttribute[key] = doc.license[key])
        )

        this.updating = true

        license
          .updateLicense(this.licenseAttribute)
          .then(res => {
            if (res.success) {
              message.success(this.props.successMsg || '许可证激活成功')
              this.props.successCallback()
            }
          })
          .finally(() => {
            this.updating = false
          })
      } else {
        message.error('请上传正确的许可证文件')
      }
    } catch (error) {
      message.error('请上传正确的许可证文件')
    }
  }

  changeFile = e => {
    e.target.value = null
  }

  validateFileContent = doc => {
    if (!doc.license) {
      return false
    } else {
      return ['expiry', 'key', 'machineid', 'name', 'version'].every(
        key => doc.license[key] && typeof doc.license[key] === 'string'
      )
    }
  }

  readFile(file) {
    if (file) {
      return new Promise((resolve, reject) => {
        var reader = new FileReader()
        reader.readAsText(file, 'UTF-8')
        reader.onload = evt => {
          resolve(evt.target.result)
        }
        reader.onerror = () => {
          reject('failed')
        }
      })
    } else {
      return ' '
    }
  }

  render() {
    return (
      <StepWrapper>
        <Spin
          tip={this.updating ? '许可证激活中, 请稍后...' : '加载中, 请稍后...'}
          spinning={this.loading || this.updating}>
          <div>
            <div className='text'>
              1.点击‘复制’按钮，复制下面的许可证申请码后发送给远算的工作人员。
            </div>
            <div className='tips license'>
              许可证申请码: <p className='machineId'>{this.machineId}</p>
              <CopyToClipboard
                text={this.machineId}
                onCopy={() => {
                  message.success('复制许可证申请码成功')
                }}>
                <Button className='btn'>复制</Button>
              </CopyToClipboard>
            </div>
          </div>
          <div>
            2.远算的工作人员会发送给您新的许可证文件，点击下面按钮，
            选中许可证文件进行上传并激活。
            <div className='upload'>
              <input
                style={{ display: 'none' }}
                type='file'
                ref={this.inputRef}
                accept='.yml, .yaml'
                onChange={this.submitFile}
                onClick={this.changeFile}
                className='file'
              />
              <Button
                icon='upload'
                onClick={() => {
                  this.inputRef.current.click()
                }}
                className='file'>
                上传许可证并激活
              </Button>
            </div>
          </div>
        </Spin>
      </StepWrapper>
    )
  }
}
