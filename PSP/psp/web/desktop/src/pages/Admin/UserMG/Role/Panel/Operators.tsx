import * as React from 'react'
import { observer } from 'mobx-react'
import { computed } from 'mobx'
import { message } from 'antd'
import { Modal } from '@/components'
import { Role, RoleList } from '@/domain/UserMG'
import { RoleType, RoleTypeName } from '@/domain/UserMG/Role'
import RoleEditor from '../RoleEditor'
import { OperatorsWrapper } from './style'

interface IProps {
  rowData: any
}

@observer
export default class Operators extends React.Component<IProps> {
  @computed
  get role() {
    const { rowData } = this.props
    return RoleList.get(rowData.id)
  }

  @computed
  get isNormalRole() {
    const { rowData } = this.props
    return rowData.type === RoleType.ROLE_NORMAL_USER
  }

  @computed
  get isInternal() {
    const { rowData } = this.props
    return rowData.isInternal === true
  }
  // 管理员不可编辑
  @computed
  get noEditRole() {
    return this.isInternal
  }

  @computed
  get noDeleteRole() {
    return this.isInternal
  }

  public render() {
    return (
      <OperatorsWrapper>
        <div
          className={`action ${this.noEditRole ? 'disabled' : ''}`}
          onClick={() => this.editRole(this.props.rowData)}>
          <span>编辑</span>
        </div>
        <div
          className={`action ${this.noDeleteRole ? 'disabled' : ''}`}
          onClick={this.delete}>
          <span>删除</span>
        </div>
      </OperatorsWrapper>
    )
  }

  private editRole = rowData => {
    if (this.noEditRole) {
      return
    }

    Modal.show({
      title: '编辑角色',
      bodyStyle: { padding: 0, height: 710, background: '#F0F5FD' },
      width: 1130,
      footer: null,
      content: ({ onCancel, onOk }) => (
        <RoleEditor
          role={new Role(this.role.toRequest())}
          onCancel={onCancel}
          onOk={onOk}
        />
      )
    })
  }

  private delete = () => {
    if (this.noDeleteRole) {
      return
    }
    Modal.showConfirm({
      content: `确认删除角色${this.role.name}吗？`
    }).then(() =>
      RoleList.delete(this.role.id, this.role.name).then(() =>
        message.success('角色删除成功')
      )
    )
  }
}
