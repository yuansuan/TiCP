import * as React from 'react'
import { observer } from 'mobx-react'
import { Checkbox,InputNumber,Tooltip } from 'antd'
import { ConfigWrapper } from './style'
import { Icon } from '@/components'
import sysConfig from '@/domain/SysConfig'
import { observable } from 'mobx'

@observer
export default class MailNotificationConfig extends React.Component<any> {
  @observable notification = {
    node_breakdown: false,
    disk_usage: false,
    agent_breakdown: false,
    job_fail_num: false,
  }


  componentDidMount() {
    const { notification } = sysConfig.mailInfoConfig
    this.notification = notification || this.notification
  }

  updateConfig = () => {
    sysConfig.updateMailNotifation(this.notification)
  }

  onChange = (type: string, value: boolean) => {
    this.notification[type] = value
    this.updateConfig()
  }

  onBlur = () => {
    this.updateConfig()
  }


  render() {
    return (
      <ConfigWrapper>
        <div className='item'>
          <Checkbox
            checked={this.notification.node_breakdown}
            onChange={e => {
              this.onChange('node_breakdown', e.target.checked)
            }}>
            调度器节点下线提醒
          </Checkbox>
          <Tooltip title='调度器节点下线且超过60s未上线，发送邮件通知！'>
            <Icon style={{ padding: '0 5px' }} type={'help-circle'} />
          </Tooltip>
        </div>
        <div className='item'>
          <Checkbox
            checked={this.notification.disk_usage}
            onChange={e => {
              this.onChange('disk_usage', e.target.checked)
            }}>
            存储使用率超阈值提醒
          </Checkbox>
          <Tooltip title='磁盘存储使用率超过阈值，发送邮件通知！'>
            <Icon style={{ padding: '0 5px' }} type={'help-circle'} />
          </Tooltip>
        </div>
        <div className='item'>
          <Checkbox
            checked={this.notification.agent_breakdown}
            onChange={e => {
              this.onChange('agent_breakdown', e.target.checked)
            }}>
            监控采集服务下线提醒
          </Checkbox>
          <Tooltip title='节点的监控采集服务下线且超过60s未上线，发送邮件通知！'>
            <Icon style={{ padding: '0 5px' }} type={'help-circle'} />
          </Tooltip>
        </div>
        <div className='item'>
          <Checkbox
            checked={this.notification.job_fail_num}
            onChange={e => {
              this.onChange('job_fail_num', e.target.checked)
            }}>
            求解作业失败超阈值提醒
          </Checkbox>
          <Tooltip title='求解的作业在24小时内失败数超过阈值，发送邮件通知！'>
            <Icon style={{ padding: '0 5px' }} type={'help-circle'} />
          </Tooltip>
        </div>
      </ConfigWrapper>
    )
  }
}
