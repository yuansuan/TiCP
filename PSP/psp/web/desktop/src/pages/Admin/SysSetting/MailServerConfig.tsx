import * as React from 'react'
import { observer } from 'mobx-react'
import { Input, Switch, InputNumber, Button, message, Tooltip } from 'antd'
import { Label } from '@/components'
import { Icon } from '@/components'
import { ConfigWrapper } from './style'
import sysConfig from '@/domain/SysConfig'
import { observable, computed } from 'mobx'
import { Validator } from '@/utils'
import { INSTALL_TYPE } from '@/utils/const'
import debounce from 'lodash/debounce'


@observer
export default class MailServerConfig extends React.Component<any> {
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
    const { email_config } = sysConfig.mailServerConfig
    this.server = email_config || this.server
  }

  sendMail = async () => {
    this.loading = true
    try {
      await sysConfig.sendMail()
      message.success('测试邮件发送成功')
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
   debouncedUpdateConfig = debounce(()=> {
    this.updateConfig();
  }, 1000);

  updateConfig = () => {
    sysConfig.updateMailServer(this.server)
  }

  onChange = (type: string, value: string | number) => {
    
    // 校验
    if( type == 'use_tls' || type == 'port'){
      // @ts-ignore
      this.server[type] = value
      this.updateConfig()
    }else {
      if (this.server[type] === (value as string)?.trim()) {
        return;
      } 
      this.server[type] = value
      const msg = this.validate(type, value)
      if (msg !== '') {
        this.message[type] = msg
        return
      } else {
        // clear
        this.message[type] = ''
      }
    
      this.debouncedUpdateConfig(); 
    }
    
  }

  validate(type: string, value: string | number) {
    let tmp = (value as string)?.trim()

    if (type === 'host') {
      if (tmp === '') return '邮件服务器地址不能为空'
      if (!(Validator.isValidDomainName(tmp) || Validator.isValidIp(tmp)))
        return '邮件服务器地址格式不对'
    }

    if (type === 'user_name') {
      if (tmp === '') return '邮箱服务器认证账号不能为空'
    }

    if (type === 'password') {
      if (tmp === '') return '邮箱服务器认证密码不能为空'
    }

    if (type === 'from') {
      if (tmp === '') return '邮件发送人地址不能为空'
      if (!Validator.isValidEmail(tmp)) return '邮件发送人地址格式不对'
    }

    if (type === 'admin_addr') {
      if (tmp === '') return '管理员邮件地址不能为空'
      let addrs = tmp.split(',')
      if (addrs.length > 10) return '管理员邮件地址最多支持10个'
      if (addrs.some(addr => !Validator.isValidEmail(addr)))
        return '管理员邮件地址格式不对'
    }

    return ''
  }


  render() {
    return (
      <ConfigWrapper>
        <div className='item'>
          <Label required starBefore={false} align={'left'}>
            邮件服务器地址
          </Label>
          <Input
            className='field'
            value={this.server.host}
            onChange={e => {
              this.onChange('host', e.target.value)
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
          <Label align={'left'} required starBefore={false}>
            邮箱服务器认证账号
          </Label>
          <Input
            className='field'
            value={this.server.user_name}
            onChange={e => {
              this.onChange('user_name', e.target.value)
            }}
          />
          <p className={'msg'}>{this.message.user_name}</p>
        </div>
        <div className='item'>
          <Label align={'left'} required starBefore={false}>
            邮箱服务器认证密码
          </Label>
          <Input.Password
            className='field'
            value={this.server.password}
            onChange={e => {
              this.onChange('password', e.target.value)
            }}
          />
          <p className={'msg'}>{this.message.password}</p>
        </div>
        <div className='item'>
          <Label required starBefore={false} align={'left'}>
            邮件发送人地址
          </Label>
          <Input
            className='field'
            value={this.server.from}
            onChange={e => {
              this.onChange('from', e.target.value)
            }}
          />
          <p className={'msg'}>{this.message.from}</p>
        </div>
        <div className='item'>
          <Label align={'left'} required starBefore={false}>
            管理员邮件地址
          </Label>
          <Input
            className='field'
            value={this.server.admin_addr}
            onChange={e => {
              this.onChange('admin_addr', e.target.value)
            }}
          />
          <Tooltip title='管理员邮件地址最多支持10个地址，以逗号(,)隔开'>
            <Icon style={{ padding: '0 5px' }} type={'help-circle'} />
          </Tooltip>
          <p className={'msg'}>{this.message.admin_addr}</p>
        </div>
        <div className='item'>
          <div style={{ width: 510 }}>
            <Button
              onClick={this.sendMail}
              disabled={!this.enabled}
              loading={this.loading}
              style={{ float: 'right' }}>
              邮件发送测试
            </Button>
          </div>
        </div>
      </ConfigWrapper>
    )
  }
}
