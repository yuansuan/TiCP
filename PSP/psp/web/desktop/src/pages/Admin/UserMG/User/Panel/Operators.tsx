import * as React from 'react'
import { computed } from 'mobx'
import { observer } from 'mobx-react'
import { message, Tooltip } from 'antd'
import {
  EditOutlined,
  DeleteOutlined,
  KeyOutlined,
  CopyOutlined
} from '@ant-design/icons'
import { Modal, Icon } from '@/components'
import { UserList, User } from '@/domain/UserMG'
import { currentUser } from '@/domain'
import UserEditor from '../UserEditor'
import { OperatorsWrapper } from './style'
import { sysConfig } from '@/domain'
import { copy2clipboard } from '@/utils'

interface IProps {
  rowData: any
}

@observer
export default class Operators extends React.Component<IProps> {
  @computed
  get user() {
    const { rowData } = this.props
    return UserList.get(rowData.id)
  }

  @computed
  get isInternalOrSelf() {
    return (
      currentUser.id === this.props.rowData.id || this.props.rowData.isInternal
    )
  }

  private editUser = () => {
    if (this.isInternalOrSelf) {
      return
    }

    if (this.user.approve_status === 0 && sysConfig.enableThreeMembers) {
      message.warn(`用户${this.user.name}有未完成的审批，请等待审批结束`)
      return
    }

    Modal.show({
      title: '编辑用户',
      bodyStyle: {
        height: 710,
        background: '#F0F5FD'
      },
      width: currentUser.authType === 'local' ? 1200 : 900,
      footer: null,
      content: ({ onCancel, onOk }) => (
        <UserEditor
          onCancel={onCancel}
          onOk={onOk}
          user={new User(this.user.toRequest())}
        />
      )
    })
  }

  private copyPasswd = pwd => {
    copy2clipboard(pwd)
    message.success('新密码复制成功')
  }

  private resetPwd = () => {
    Modal.showConfirm({
      content: `确定重置用户 ${this.user.name} 的密码 ?`
    }).then(() =>
      UserList.resetPwd(this.user.id).then(res => {
        Modal.show({
          title: '密码重置成功',
          content: (
            <div
              style={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center'
              }}>
              <p>用户 {this.user.name} 的密码重置成功</p>
              <p>
                新密码：{res.data}{' '}
                <CopyOutlined onClick={() => this.copyPasswd(res.data)} />
              </p>
            </div>
          ),
          footer: null
        })
      })
    )
  }

  private delete = () => {
    if (this.isInternalOrSelf) {
      return
    }

    if (this.user.approve_status === 0 && sysConfig.enableThreeMembers) {
      message.warn(`用户${this.user.name}有未完成的审批，请等待审批结束`)
      return
    }

    Modal.showConfirm({
      content: sysConfig.enableThreeMembers
        ? `确认发起删除用户${this.user.name}申请吗？`
        : `确认删除用户${this.user.name}吗？`
    }).then(() =>
      UserList.delete(this.user.id).then(res => {
        message.success('用户删除成功')
      })
    )
  }

  render() {
    const disabledClassName = this.isInternalOrSelf ? 'disabled' : ''

    return (
      <OperatorsWrapper>
        <div className={`action ${disabledClassName}`} onClick={this.editUser}>
          <Tooltip title='编辑'>
            <EditOutlined />
          </Tooltip>
        </div>
        {currentUser?.isLdapEnabled ? (
          <div className={`action ${disabledClassName}`} onClick={this.delete}>
            <Tooltip title='删除'>
              <DeleteOutlined />
            </Tooltip>
          </div>
        ) : null}
        {currentUser?.isSuperAdmin ? (
          <div
            className={`action ${
              currentUser.id === this.props.rowData.id ? 'disabled' : ''
            }`}
            onClick={this.resetPwd}>
            <Tooltip title='重置密码'>
              <KeyOutlined />
            </Tooltip>
          </div>
        ) : null}
      </OperatorsWrapper>
    )
  }
}
