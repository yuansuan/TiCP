import * as React from 'react'
import { observer } from 'mobx-react'
import { Input, Switch, InputNumber, Select } from 'antd'
import { Label } from '@/components'
import { ConfigWrapper } from './style'
import sysConfig from '@/domain/SysConfig'
import { observable, computed } from 'mobx'
import { Validator } from '@/utils'
import { INSTALL_TYPE } from '@/utils/const'
import { Http } from '@/utils'
@observer
export default class JobConfig extends React.Component<any> {
  @observable job = {
    queue: '' //作业队列
  }

  @observable loading = false
  @observable queues = []

  @observable message = {
    queue: ''
  }

  @computed
  get defaultQueue() {
    return this.queues
      ?.filter(item => item.select)
      ?.map(item => item.queue_name === this.job.queue)
  }
  async componentDidMount() {
    const { data } = await sysConfig.fetchJobConfig()
    this.job = data

    const {
      data: { queues }
    } = await Http.get('/app/queue')

    this.queues = queues
  }
  updateConfig = () => {
    sysConfig.updateJobConfig(this.job)
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
          <div className='left'>
            <Label align={'left'}>默认作业队列名称</Label>
            <Select
              className='field'
              onChange={value => {
                this.onChange('queue', value)
              }}
              allowClear={true}
              value={this.job.queue}
              placeholder='请选择默认队列'>
              {this.queues?.map(item => {
                return (
                  <Select.Option
                    title={item.queue_name}
                    key={item.queue_name}
                    value={item.queue_name}>
                    {item.queue_name}
                  </Select.Option>
                )
              })}
            </Select>
            <p className={'msg'}>{this.message.queue}</p>
          </div>
        </div>
      </ConfigWrapper>
    )
  }
}
