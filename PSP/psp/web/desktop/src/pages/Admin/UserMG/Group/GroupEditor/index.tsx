import * as React from 'react'
import { observable, action, computed } from 'mobx'
import { observer } from 'mobx-react'
import { Input, message } from 'antd'

import { Validator } from '@/utils'
import { RoleList, UserList, GroupList } from '@/domain/UserMG'
import { isRoleConflict, RoleType } from '@/domain/UserMG/Role'

import { GroupEditorWrapper } from './style'
import { SelectEditor, PendingFooter } from '../../components'
import { checkPerm } from '../../utils'
import { sysConfig } from '@/domain'

interface IProps {
  group: any
  onCancel: () => void
  onOk: () => void
}

const SelectEditorTitle = { leftTab: '用户', rightTab: '角色' }
const MAX_USER_NUM = 300

@observer
export default class GroupEditor extends React.Component<IProps> {
  @observable errorMessage = ''
  @observable adding = false
  @action
  updateName = name => (this.props.group.name = name)

  @action
  updateSelectedRoleKeys = (keys, items) => {
    if (keys.includes(RoleType.ROLE_NORMAL_USER)) {
      message.error('用户组不允许添加普通用户角色')
      return
    }

    if (
      keys.length > 1 &&
      isRoleConflict(keys) &&
      sysConfig.enableThreeMembers
    ) {
      message.error(
        '系统管理员，安全管理员和审计管理员，这三种内置角色只能选择其中之一。'
      )
      return
    }

    this.props.group.roles = items.map(i => i.id)
    this.props.group.roleNames = items.map(i => i.name)
  }

  @action
  updateSelectedUserKeys = (keys, items) => {
    if (keys.length > 0 && !this.editMode && sysConfig.enableThreeMembers) {
      message.error(`注意: 申请创建新用户组时, 不能选择用户`)
      return
    }

    if (keys.length > MAX_USER_NUM) {
      message.error(`用户组不能超过 ${MAX_USER_NUM} 人`)
      return
    }
    this.props.group.users = items.map(i => i.id)
    this.props.group.userNames = items.map(i => i.name)
  }
  @action
  updateAdding = flag => (this.adding = flag)

  @computed
  get name() {
    return this.props.group.name
  }

  @computed
  get selectedUserNames() {
    return this.props.group.userNames
  }

  @computed
  get selectedRoleNames() {
    return this.props.group.roleNames
  }

  @computed
  get selectedUserKeys() {
    return this.props.group.users
  }

  @computed
  get selectedRoleKeys() {
    return this.props.group.roles
  }

  @computed
  get editMode() {
    return this.props.group.id > 0
  }

  public componentDidMount() {
    if (this.editMode) {
      this.props.group.fetch()
    }
  }

  private onOk = () => {
    if (!this.name) {
      message.error('请输入用户组名称')
      return
    }

    if (this.errorMessage) {
      // 校验不通过
      message.error(this.errorMessage)
      return
    }

    if (this.selectedUserKeys.length > MAX_USER_NUM) {
      message.error(`用户组不能超过 ${MAX_USER_NUM} 人`)
      return
    }

    this.updateAdding(true)

    if (this.editMode) {
      this.props.group
        .update({
          name: this.name,
          roles: this.selectedRoleKeys,
          users: this.selectedUserKeys,
          userNames: this.selectedUserNames,
          roleNames: this.selectedRoleNames
        })
        .then(res => {
          if (res.data?.isAskRequest) {
            res.success
              ? message.success(res.message)
              : message.error(res.message)
          } else {
            message.success('用户组修改成功')
          }

          this.props.onOk()
        })
        .finally(() => {
          this.updateAdding(false)
          checkPerm()
        })
    } else {
      GroupList.add({
        name: this.name,
        roles: this.selectedRoleKeys,
        users: this.selectedUserKeys,
        userNames: this.selectedUserNames,
        roleNames: this.selectedRoleNames
      })
        .then(res => {
          if (res.data?.isAskRequest) {
            res.success
              ? message.success(res.message)
              : message.error(res.message)
          } else {
            message.success('用户组添加成功')
          }
          this.props.onOk()
        })
        .finally(() => {
          this.updateAdding(false)
        })
    }
  }

  @computed
  get leftList() {
    return {
      list: UserList.enabledUsers,
      fetch: () => UserList.fetch().then(res => res.data.users)
    }
  }

  @computed
  get rightList() {
    return {
      list: RoleList.roleList,
      fetch: () => RoleList.fetch().then(res => res.data.roles)
    }
  }

  private onChangeName = e => this.updateName(e.target.value)
  private onBlurName = e => {
    const name = e.target.value
    if (!Validator.isValidInputName(name)) {
      this.errorMessage = '用户组名称只能包含字母,汉字,数字和下划线'
    } else if (name.length > 32) {
      this.errorMessage = '用户组名称长度不能大于 32 字符'
    } else {
      this.errorMessage = ''
    }
  }

  render() {
    return (
      <GroupEditorWrapper>
        <div className='groupName'>
          <span>*</span>
          用户组名称：
          <Input
            size='small'
            value={this.name}
            onChange={this.onChangeName}
            autoFocus
            onBlur={this.onBlurName}
            onFocus={e => e.target.select()}
          />
        </div>
        <div style={{ fontSize: 12, color: 'red' }}>
          注意：一个用户组最多包含 {MAX_USER_NUM} 个用户
        </div>
        <SelectEditor
          selectedLeftKeys={this.selectedUserKeys}
          selectedRightKeys={this.selectedRoleKeys}
          updateSelectedLeftKeys={this.updateSelectedUserKeys}
          updateSelectedRightKeys={this.updateSelectedRoleKeys}
          LeftList={this.leftList}
          RightList={this.rightList}
          title={SelectEditorTitle}
        />

        <PendingFooter
          onCancel={this.props.onCancel}
          onOk={this.onOk}
          processing={this.adding}
        />
      </GroupEditorWrapper>
    )
  }
}
