import * as React from 'react'
import { observer } from 'mobx-react'
import { observable } from 'mobx'
import { Label } from '@/components'
import { ConfigWrapper } from './style'
import { Radio, message } from 'antd'
import sysConfig from '@/domain/SysConfig'
import { Http } from '@/utils'

const radioStyle = {
  display: 'block',
  height: 60,
}

const descStyle = {
  fontSize: 12,
  margin: '5px 0',
}

@observer
export default class FireWallConfig extends React.Component<any> {
  @observable level = sysConfig.firewallConfig.level
  @observable updating = false

  onChange = async e => {
    const level = e.target.value
    this.updating = true

    try {
      const res = await Http.get('/sysconfig/optinfo')
      if (res.data?.id) {
        const status = res.data.status === 0 ? false : true
        if (status) {
          message.error('远程运维已开启，防火墙策略不能修改，请关闭远程运维后操作。')
          return
        }
      }
      await sysConfig.updateFireWallConfig(level)
      message.info('防火墙策略修改成功，策略生效需要大约3分钟，请稍后...')
    } finally {
      this.updating = false
    }

    this.level = level
  }

  render() {
    return (
      <ConfigWrapper>
        <div className='item'>
          <Label align={'left'}>对外访问控制</Label>
          <Radio.Group
            disabled={this.updating}
            className='field'
            value={this.level}
            onChange={this.onChange}>
            <Radio style={radioStyle} value={'all'}>
              不限制
              <p style={descStyle}>对外访问不做任何限制</p>
            </Radio>
            <Radio style={radioStyle} value={'assign'}>
              仅访问远算或者渠道商网络
              <p style={descStyle}>对外仅能访问远算或者渠道商网络</p>
            </Radio>
            <Radio style={radioStyle} value={'none'}>
              完全限制
              <p style={descStyle}>对外不能做任何访问</p>
            </Radio>
          </Radio.Group>
        </div>
      </ConfigWrapper>
    )
  }
}
