import * as React from 'react'
import { observer } from 'mobx-react'
import { Input, Switch, InputNumber, Button, message, Tooltip } from 'antd'
import { Label } from '@/components'
import { ConfigWrapper } from './style'
import sysConfig from '@/domain/SysConfig'
import { observable, computed } from 'mobx'
import { Validator } from '@/utils'
import { INSTALL_TYPE } from '@/utils/const'

@observer
export default class ADDomainConfig extends React.Component<any> {
  @observable server = {
    from: null,
    host: null,
    password: null,
    port: 25,
    use_tls: false,
    user_name: null,
    admin_addr: null,
  }

  @observable loading = false

  @observable message = {
    from: '',
    host: '',
    admin_addr: '',
    user_name: '',
    password: '',
  }

  @computed
  get enabled() {
    let noMessages =
      this.message.host === '' &&
      this.message.from === '' &&
      this.message.admin_addr === '' &&
      this.message.user_name === '' &&
      this.message.password === ''
    return (
      this.server.host &&
      this.server.from &&
      this.server.admin_addr &&
      this.server.user_name &&
      this.server.password &&
      noMessages
    )
  }

  componentDidMount() {
    const { email } = sysConfig.mailConfig
    this.server = { ...email?.server }
  }

  testConnect = async () => {
    this.loading = true
    try {
      await sysConfig.sendMail(this.server)
      message.success('测试发送成功')
      this.loading = false
    } catch (e: any) {
      if (e?.message?.includes('timeout')) {
        sysConfig.installType === INSTALL_TYPE.aio
          ? message.error(
              `无法连接，请检查防火墙设置策略或是否可访问邮件服务器！`
            )
          : message.error(`无法连接，请检查是否可以连接到邮件服务器！`)
      }

      this.loading = false
    }
  }

  updateConfig = () => {
    sysConfig.updateMailServer(this.server)
  }

  onChange = (type: string, value: string | boolean | number) => {
    this.server[type] = value
    if (type === 'use_tls') {
      this.updateConfig()
    }
  }

  validate(type: string, value: string | number) {
    let tmp = (value as string).trim()

    if (type === 'host') {
      if (tmp === '') return '服务器地址不能为空'
      if (!(Validator.isValidDomainName(tmp) || Validator.isValidIp(tmp)))
        return '服务器地址格式不对'
    }

    if (type === 'user_name') {
      if (tmp === '') return '邮箱服务器认证账号不能为空'
    }

    if (type === 'password') {
      if (tmp === '') return '邮箱服务器认证密码不能为空'
    }

    
    return ''
  }

  onBlur = (type: string, value: string | number) => {
    // 校验
    const msg = this.validate(type, value)
    if (msg !== '') {
      this.message[type] = msg
      return
    } else {
      // clear
      this.message[type] = ''
    }

    this.updateConfig()
  }

  render() {
    return (
      <ConfigWrapper>
        <div className='item'>
          <Label required starBefore={false} align={'left'}>
            AD服务器地址
          </Label>
          <Input
            className='field'
            value={this.server.host}
            onChange={e => {
              this.onChange('host', e.target.value)
            }}
            onBlur={e => {
              this.onBlur('host', e.target.value)
            }}
          />
          <p className={'msg'}>{this.message.host}</p>
        </div>
        <div className='item'>
          <Label align={'left'}>端口</Label>
          <InputNumber
            style={{ marginLeft: 30 }}
            size='small'
            min={0}
            max={65535}
            step={1}
            value={this.server.port}
            precision={0}
            parser={value => (isNaN(parseInt(value)) ? 0 : parseInt(value))}
            onChange={value => {
              this.onChange('port', value)
            }}
            onBlur={e => {
              this.onBlur('port', e.target.value)
            }}
          />
        </div>
        <div className='item'>
          <Label align={'left'}>是否加密(TLS)</Label>
          <Switch
            style={{ marginLeft: 30 }}
            checked={this.server.use_tls}
            onChange={value => {
              this.onChange('use_tls', value)
            }}
          />
        </div>
      
        <div className='item'>
          <div style={{ width: 440 }}>
            <Button
              onClick={this.testConnect}
              disabled={!this.enabled}
              loading={this.loading}
              style={{ float: 'right' }}>
              发送测试
            </Button>
          </div>
        </div>
      </ConfigWrapper>
    )
  }
}
