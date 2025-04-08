import * as React from 'react'
import { computed } from 'mobx'
import { observer } from 'mobx-react'
import { message } from 'antd'

import { Modal } from '@/components'
import { Group, GroupList } from '@/domain/UserMG'
import GroupEditor from '../GroupEditor'
import { OperatorsWrapper } from './style'
import { sysConfig } from '@/domain'

interface IProps {
  rowData: any
}

@observer
export default class Operators extends React.Component<IProps> {
  @computed
  get group() {
    const { rowData } = this.props
    return GroupList.get(rowData.id)
  }

  private editGroup = () => {
    if (this.group.approve_status === 0 && sysConfig.enableThreeMembers) {
      message.warn(`用户组${this.group.name}有未完成的审批，请等待审批结束`)
      return
    }

    Modal.show({
      title: '编辑用户组',
      bodyStyle: { height: 710, background: '#F0F5FD' },
      width: 1130,
      footer: null,
      content: ({ onCancel, onOk }) => (
        <GroupEditor
          group={new Group(this.group.toRequest())}
          onCancel={onCancel}
          onOk={onOk}
        />
      ),
    })
  }

  private delete = () => {
    if (this.group.approve_status === 0 && sysConfig.enableThreeMembers) {
      message.warn(`用户组${this.group.name}有未完成的审批，请等待审批结束`)
      return
    }

    Modal.showConfirm({
      content: sysConfig.enableThreeMembers
        ? `确认发起删除用户组${this.group.name}的申请吗？`
        : `确认删除用户组${this.group.name}吗？`,
    }).then(() =>
      GroupList.delete(this.group.id, this.group.name, {
        roles: this.group.roles,
        roleNames: this.props.rowData.roles,
        users: this.group.users,
        userNames: this.props.rowData.users,
      }).then(res => {
        if (res.data?.isAskRequest) {
          res.success
            ? message.success(res.message)
            : message.error(res.message)
        } else {
          message.success('用户组删除成功')
        }
      })
    )
  }

  render() {
    return (
      <OperatorsWrapper>
        <div className='action' onClick={this.editGroup}>
          <span>编辑</span>
        </div>
        <div className='action' onClick={this.delete}>
          <span>删除</span>
        </div>
      </OperatorsWrapper>
    )
  }
}
