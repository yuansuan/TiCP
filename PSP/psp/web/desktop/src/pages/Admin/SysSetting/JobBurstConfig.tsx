import * as React from 'react'
import { observer } from 'mobx-react'
import { Switch, InputNumber } from 'antd'
import { Label } from '@/components'
import { ConfigWrapper } from './style'
import sysConfig from '@/domain/SysConfig'
import { observable, computed } from 'mobx'

@observer
export default class JobConfig extends React.Component<any> {
  @observable job = {
    threshold: 1, //爆发阈值, 默认一小时
    unit: 'hour', //阈值数值单位 minute、hour
    enable: false // 是否开启
  }

  @observable loading = false

  @observable message = {
    queue: ''
  }

  async componentDidMount() {
    const { data } = await sysConfig.fetchJobBurstConfig()
    this.job = data
  }
  updateConfig = () => {
    sysConfig.updateJobBurstConfig(this.job)
  }

  onChange = (type: string, value: string | boolean | number) => {
    this.job[type] = value
    if (type === 'enable' || type === 'queue') {
      this.updateConfig()
    }
  }

  validate(type: string, value: string | number) {
    let tmp = (value as string).trim()
    if (type === 'queue' && tmp === '') {
      return '默认队列不能为空'
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
          <Label align={'left'}>自动爆发</Label>
          <Switch
            checkedChildren='开启'
            unCheckedChildren='关闭'
            style={{ marginLeft: 30 }}
            checked={this.job.enable}
            onChange={value => {
              this.onChange('enable', value)
            }}
          />
        </div>
        {this.job?.enable && (
          <>
            <div className='item'>
              <div className='left'>
                <Label align={'left'}>作业最大等待时长</Label>
                <InputNumber
                  className='field'
                  min={0}
                  step={1}
                  value={this.job.threshold}
                  precision={0}
                  parser={text =>
                    text && Math.round(Number(text.replace(/[^0-9.]+/g, '')))
                  }
                  onChange={value => {
                    this.onChange('threshold', value)
                  }}
                  onBlur={e => {
                    this.onBlur('threshold', e.target.value)
                  }}
                />
                <div className='unit'>分钟</div>
              </div>
            </div>
          </>
        )}
      </ConfigWrapper>
    )
  }
}
