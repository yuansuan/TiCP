import * as React from 'react'
import { computed } from 'mobx'
import { observer } from 'mobx-react'
import { message, Switch, Tooltip } from 'antd'

import { Modal } from '@/components'
import { UserList } from '@/domain/UserMG'
import { currentUser } from '@/domain'
import { StateWrapper } from './style'
import { sysConfig } from '@/domain'

interface IProps {
  rowData: any
}

@observer
export default class State extends React.Component<IProps> {
  @computed
  get user() {
    const { rowData } = this.props
    return UserList.get(rowData.id)
  }

  @computed
  get disableOpt() {
    return (
      currentUser.id === this.props.rowData.id || this.props.rowData.isInternal
    )
  }

  @computed
  get isActive() {
    return this.user?.enabled
  }

  private active = () => {
    Modal.showConfirm({
      content: sysConfig.enableThreeMembers
        ? `确认发起启用用户${this.user.name}申请吗？`
        : `确认启用用户${this.user.name}吗？`
    }).then(() =>
      UserList.active(this.user.id).then(res => {
        if (res.data) {
          res.success
            ? message.success(res.message)
            : message.error(res.message)
        } else {
          message.success('启用用户成功')
        }
      })
    )
  }

  private inactive = () => {
    Modal.showConfirm({
      content: sysConfig.enableThreeMembers
        ? `确认发起禁用用户${this.user.name}申请吗？`
        : `确认禁用用户${this.user.name}吗？`
    }).then(() =>
      UserList.inactive(this.user.id).then(res => {
        if (res.data) {
          res.success
            ? message.success(res.message)
            : message.error(res.message)
        } else {
          message.success('禁用用户成功')
        }
      })
    )
  }

  onChange = checked => {
    if (checked) {
      this.active()
    } else {
      this.inactive()
    }
  }

  render() {
    return (
      <StateWrapper>
        <Tooltip
          title={
            this.props.rowData.isInternal ? '内置系统账号不能被禁用或启用' : ''
          }>
          <Switch
            disabled={this.disableOpt}
            checked={this.isActive}
            onChange={this.onChange}
            size='small'
          />
        </Tooltip>
        {this.isActive ? '  启用' : '  禁用'}
      </StateWrapper>
    )
  }
}
