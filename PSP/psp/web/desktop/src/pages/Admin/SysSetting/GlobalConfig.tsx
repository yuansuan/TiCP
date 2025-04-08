import * as React from 'react'
import { observer } from 'mobx-react'
import { message, Select } from 'antd'
import { Label } from '@/components'
import { ConfigWrapper } from './style'
import sysConfig from '@/domain/SysConfig'
import { observable } from 'mobx'
import { EditableText } from '@/components'
import { Validator } from '@/utils'

const { Option } = Select

@observer
export default class GlobalConfig extends React.Component<any> {
  @observable homedir = ''
  @observable show_verify_code = 3

  componentDidMount() {
    const { homedir, show_verify_code } = sysConfig.userConfig
    this.homedir = homedir
    this.show_verify_code =
      typeof show_verify_code !== 'number' ? 3 : show_verify_code
  }

  changeVCode = value => {
    sysConfig.updateUserVCode(value)
    this.show_verify_code = value
  }

  onConfirm = async value => {
    // 校验 homedir 文件路径是否存在
    let res = await sysConfig.checkHomedir(value)
    if (res.success) {
      this.homedir = value.trim()
      sysConfig.updateUserHomeDir(this.homedir)
    }
  }

  render() {
    return (
      <ConfigWrapper>
        <div className='item'>
          <Label align={'left'}>home 默认路径</Label>
          <EditableText
            style={{ width: 340, marginLeft: 30 }}
            beforeConfirm={value => {
              if (!value) {
                message.error('home 默认路径不能为空')
                return false
              }

              let flag = Validator.isValidPath(value)

              if (!flag) {
                message.error('home 默认路径格式不对')
                return false
              }

              return flag
            }}
            help='请确保新的 home 默认路径满足权限要求，用户和组是 root，权限是 755。'
            defaultValue={this.homedir}
            onConfirm={this.onConfirm}
          />
        </div>
        <div className='item'>
          <Label align={'left'}>登录验证码</Label>
          <Select
            value={this.show_verify_code}
            onChange={this.changeVCode}
            size={'small'}
            style={{ width: 240, marginLeft: 30 }}>
            <Option value={0}>不显示</Option>
            <Option value={-1}>永久显示</Option>
            <Option value={1}>登录出错1次后显示</Option>
            <Option value={3}>登录出错3次后显示</Option>
            <Option value={5}>登录出错5次后显示</Option>
          </Select>
        </div>
      </ConfigWrapper>
    )
  }
}
