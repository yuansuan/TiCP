import * as React from 'react'
import { observer } from 'mobx-react'
import { InputNumber, Checkbox } from 'antd'
import { Label } from '@/components'
import { ConfigWrapper } from './style'
import sysConfig from '@/domain/SysConfig'
import { observable } from 'mobx'

@observer
export default class PasswdConfig extends React.Component<any> {
  @observable passwd = {
    max_len: 32,
    min_len: 8,
    min_age: 0,
    max_age: 7776000,
    in_history: 3,
    expire_warning: 604800,
    max_failure: 6,
    failure_count_interval: 300,
    lockout_duration: 1800,
    strength: [],
    show_verify_code: 3
  }

  componentDidMount() {
    const { password } = sysConfig.userConfig
    this.passwd = password
  }

  updateConfig = () => {
    sysConfig.updateUserPwdConf(this.passwd)
  }

  onChange = (type: string, value: number | string, checked?: boolean) => {
    if (type === 'strength') {
      if (checked) {
        this.passwd.strength.push(value)
      } else {
        ;(this.passwd.strength as any).remove(value)
      }
      this.updateConfig()
    } else {
      this.passwd[type] = value
    }
  }

  onBlur = () => {
    this.updateConfig()
  }

  render() {
    return (
      <ConfigWrapper>
        <div className='item'>
          <Label align={'left'}>密码最小长度</Label>
          <InputNumber
            style={{ marginLeft: 30, marginRight: 10 }}
            min={6}
            max={65535}
            step={1}
            value={this.passwd?.min_len}
            precision={0}
            parser={value => (isNaN(parseInt(value)) ? 6 : parseInt(value))}
            onChange={value => {
              this.onChange('min_len', value)
            }}
            onBlur={e => {
              this.onBlur()
            }}
          />
          位
        </div>
        <div className='item'>
          <Label align={'left'}>密码强度</Label>
          <Checkbox
            style={{ marginLeft: 30 }}
            value='Char_Upper'
            onChange={e => {
              this.onChange('strength', e.target.value, e.target.checked)
            }}
            checked={this.passwd?.strength?.includes('Char_Upper')}>
            大写字母
          </Checkbox>
          <Checkbox
            value='Char_Lower'
            onChange={e => {
              this.onChange('strength', e.target.value, e.target.checked)
            }}
            checked={this.passwd?.strength?.includes('Char_Lower')}>
            小写字母
          </Checkbox>
          <Checkbox
            value='Number'
            onChange={e => {
              this.onChange('strength', e.target.value, e.target.checked)
            }}
            checked={this.passwd?.strength?.includes('Number')}>
            数字
          </Checkbox>
          <Checkbox
            value='Char_Special'
            onChange={e => {
              this.onChange('strength', e.target.value, e.target.checked)
            }}
            checked={this.passwd?.strength?.includes('Char_Special')}>
            特殊符号
          </Checkbox>
        </div>
        <div className='item'>
          <Label align={'left'}>密码最短期限</Label>
          <InputNumber
            style={{ marginLeft: 30, marginRight: 10 }}
            min={0}
            max={65535}
            step={1}
            value={this.passwd?.min_age}
            precision={0}
            parser={value => (isNaN(parseInt(value)) ? 0 : parseInt(value))}
            onChange={value => {
              this.onChange('min_age', value)
            }}
            onBlur={e => {
              this.onBlur()
            }}
          />
          天
        </div>
        <div className='item'>
          <Label align={'left'}>密码最长期限</Label>
          <InputNumber
            style={{ marginLeft: 30, marginRight: 10 }}
            min={0}
            max={65535}
            step={1}
            value={this.passwd?.max_age}
            precision={0}
            parser={value => (isNaN(parseInt(value)) ? 0 : parseInt(value))}
            onChange={value => {
              this.onChange('max_age', value)
            }}
            onBlur={e => {
              this.onBlur()
            }}
          />
          天
        </div>
        <div className='item'>
          <Label align={'left'}>历史密码个数</Label>
          <InputNumber
            style={{ marginLeft: 30, marginRight: 10 }}
            min={0}
            max={65535}
            step={1}
            value={this.passwd?.in_history}
            precision={0}
            parser={value => (isNaN(parseInt(value)) ? 0 : parseInt(value))}
            onChange={value => {
              this.onChange('in_history', value)
            }}
            onBlur={e => {
              this.onBlur()
            }}
          />
          个
        </div>
        <div className='item'>
          <Label align={'left'}>密码过期警告</Label>
          <InputNumber
            style={{ marginLeft: 30, marginRight: 10 }}
            min={0}
            max={65535}
            step={1}
            value={this.passwd?.expire_warning}
            precision={0}
            parser={value => (isNaN(parseInt(value)) ? 0 : parseInt(value))}
            onChange={value => {
              this.onChange('expire_warning', value)
            }}
            onBlur={e => {
              this.onBlur()
            }}
          />
          天
        </div>
        <div className='item'>
          <Label align={'left'}>失败允许次数</Label>
          <InputNumber
            style={{ marginLeft: 30, marginRight: 10 }}
            min={0}
            max={65535}
            step={1}
            value={this.passwd?.max_failure}
            precision={0}
            parser={value => (isNaN(parseInt(value)) ? 0 : parseInt(value))}
            onChange={value => {
              this.onChange('max_failure', value)
            }}
            onBlur={e => {
              this.onBlur()
            }}
          />
          次
        </div>
        <div className='item'>
          <Label align={'left'}>失败复位时间</Label>
          <InputNumber
            style={{ marginLeft: 30, marginRight: 10 }}
            min={0}
            max={65535}
            step={1}
            value={this.passwd?.failure_count_interval}
            precision={0}
            parser={value => (isNaN(parseInt(value)) ? 0 : parseInt(value))}
            onChange={value => {
              this.onChange('failure_count_interval', value)
            }}
            onBlur={e => {
              this.onBlur()
            }}
          />
          秒
        </div>
        <div className='item'>
          <Label align={'left'}>自动解锁时间</Label>
          <InputNumber
            style={{ marginLeft: 30, marginRight: 10 }}
            min={0}
            max={65535}
            step={1}
            value={this.passwd?.lockout_duration}
            precision={0}
            parser={value => (isNaN(parseInt(value)) ? 0 : parseInt(value))}
            onChange={value => {
              this.onChange('lockout_duration', value)
            }}
            onBlur={e => {
              this.onBlur()
            }}
          />
          秒
        </div>
      </ConfigWrapper>
    )
  }
}
