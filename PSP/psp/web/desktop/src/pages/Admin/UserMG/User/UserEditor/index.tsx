import * as React from 'react'
import { observable, action, computed } from 'mobx'
import { observer } from 'mobx-react'
import { message } from 'antd'

import { GroupList, RoleList } from '@/domain/UserMG'
import { RoleType, isRoleConflict } from '@/domain/UserMG/Role'
import BasicInfo from './BasicInfo'
import { SelectEditor, PendingFooter } from '../../components'
import { checkPerm } from '../../utils'
import { UserEditorWrapper } from './style'

interface IProps {
  user: any
  onCancel: () => void
  onOk: () => void
}

enum MessageType {
  EMAIL = 'email',
  MOBILE = 'mobile'
}

const SelectEditorTitle = { leftTab: '角色' }

@observer
export default class UserEditor extends React.Component<IProps> {
  @observable errMessage = {
    email: '',
    mobile: ''
  }
  @observable adding = false
  @action
  updateEmail = email => (this.props.user.email = email)
  @action
  updateMobile = mobile => (this.props.user.mobile = mobile)
  @action
  updateEnableOenapi = enable_openapi => (this.props.user.enable_openapi = enable_openapi)

  @action
  updateSelectedGroupKeys = (keys, items) => {
    this.props.user.groups = items.map(i => i.id)
    this.props.user.groupNames = items.map(i => i.name)
  }

  @action
  updateSelectedRoleKeys = (keys, items) => {
    if (keys.length > 1 && isRoleConflict(keys)) {
      message.error(
        '系统管理员，安全管理员和审计管理员，这三种内置管理员角色只能选择其中之一。'
      )
      return
    }
    this.props.user.roles = items.map(i => i.id)
    this.props.user.roleNames = items.map(i => i.name)
  }

  @action
  updateAdding = flag => (this.adding = flag)
  @action
  updateErrorMessage = (type: MessageType, message) =>
    (this.errMessage[type] = message)

  @computed
  get email() {
    return this.props.user.email
  }

  @computed
  get mobile() {
    return this.props.user.mobile
  }

  @computed
  get selectedGroupKeys() {
    return this.props.user.groups
  }

  @computed
  get selectedRoleKeys() {
    return this.props.user.roles
  }

  @computed
  get isInternal() {
    return this.props.user.isInternal
  }

  private onOk = () => {
    if (this.selectedRoleKeys.length === 0) {
      message.error('用户角色不能为空')
      return
    }
    if (this.errMessage.email) {
      message.error(this.errMessage.email)
      return
    }

    if (this.errMessage.mobile) {
      message.error(this.errMessage.mobile)
      return
    }

    this.updateAdding(true)
    this.props.user
      .update()
      .then(() => {
        message.success('用户修改成功')
        this.props.onOk()
      })
      .catch(e => {
        if (e.fake) {
          if (e.success) {
            // 关闭对话框
            this.props.onCancel()
          }
        }
      })
      .finally(() => {
        this.updateAdding(false)
        checkPerm()
      })
  }

  @computed
  get leftList() {
    return {
      list: RoleList.roleList,
      fetch: () => RoleList.fetch().then(res => res.data.roles)
    }
  }

  @computed
  get rightList() {
    return {
      list: GroupList.groupList,
      fetch: () => GroupList.fetch().then(res => res.data.groups)
    }
  }

  render() {
    return (
      <UserEditorWrapper>
        <BasicInfo
          name={this.props.user.name}
          email={this.email}
          mobile={this.mobile}
          enable_openapi={this.props.user.enable_openapi}
          updateEmail={this.updateEmail}
          updateMobile={this.updateMobile}
          updateEnableOpenapi={this.updateEnableOenapi}
          updateError={this.updateErrorMessage
            
          }
        />
        <SelectEditor
          selectedLeftKeys={this.selectedRoleKeys}
          // selectedRightKeys={this.selectedGroupKeys}
          updateSelectedLeftKeys={this.updateSelectedRoleKeys}
          // updateSelectedRightKeys={this.updateSelectedGroupKeys}
          leftDisabledCondition={item =>
            item.type === RoleType.ROLE_NORMAL_USER ||
            (this.isInternal && item.type === RoleType.ROLE_ADMIN)
          }
          LeftList={this.leftList}
          // RightList={this.rightList}
          title={SelectEditorTitle}
        />

        <PendingFooter
          onCancel={this.props.onCancel}
          onOk={this.onOk}
          processing={this.adding}
        />
      </UserEditorWrapper>
    )
  }
}
