import * as React from 'react'
import { Form, Input, message, Tooltip } from 'antd'
import { observer } from 'mobx-react'
import { observable, action, computed } from 'mobx'

import { UserList, RoleList, User, GroupList } from '@/domain/UserMG'
import { RoleType, isRoleConflict } from '@/domain/UserMG/Role'
import { UserCreatorWrapper } from './style'
import { PendingFooter, SelectEditor } from '../../components'

interface IProps {
  user: User
  onCancel?: () => void
  onOk?: () => void
}

const SelectEditorTitle = { leftTab: '角色' }
@observer
export default class UserCreator extends React.Component<IProps> {
  @observable username: string
  @observable password: string
  @observable adding = false
  @action
  updateUsername = name => {
    if (!name) {
      message.error('请输入用户名！')
      return
    } else if (name.length > 64) {
      message.error('用户名长度不能大于64个字符')
      return
    }
    this.username = name
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

  @computed
  get selectedRoleKeys() {
    return this.props.user.roles
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

  public async componentDidMount() {
    await UserList.fetchDisabled()
  }

  private onCancel = () => {
    const { onCancel } = this.props
    onCancel && onCancel()
  }

  private onOk = () => {
    const { onOk, user: userData } = this.props
    if (this.username === undefined) {
      message.error('请选择用户')
      return
    }
    this.updateAdding(true)

    //   const list = this.userIds.map(id => {
    //     const user = UserList.disabledUsers.filter(n => n.id === id)[0]
    //     return {
    //       id: id,
    //       name: user.name,
    //       roles: userData.roles,
    //       roleNames: userData.roleNames
    //     }
    //   })

    //   if (this.username) {
    //     UserList.addAll(list)
    //       .then(() => {
    //         if (
    //           UserList.successUsername &&
    //           UserList.successUsername.length !== 0
    //         ) {
    //           message.success(`用户${UserList.successUsername}添加成功`, 3)
    //         }
    //         if (
    //           UserList.failureUsername &&
    //           UserList.failureUsername.length !== 0
    //         ) {
    //           message.error(`用户${UserList.failureUsername}添加失败`, 3)
    //         }
    //         onOk && onOk()
    //       })
    //       .finally(() => this.updateAdding(false))
    //   }
  }

  filterUser = (input, option) =>
    option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
  hidden = id => {
    const userName = UserList.disabledUsers
      .filter(user => id.includes(user.id))
      .map(user => user.name)
      .toString()

    return <Tooltip title={userName}>+{id.length}...</Tooltip>
  }

  render() {
    return (
      <UserCreatorWrapper>
        <div className='nameSelect'>
          <Form
            name='basic'
            labelCol={{ span: 8 }}
            wrapperCol={{ span: 16 }}
            style={{ maxWidth: 400 }}
            autoComplete='on'>
            <Form.Item
              label='用户名'
              name='username'
              rules={[{ required: true, message: '请输入用户名！' }]}>
              <Input maxLength={64} onChange={this.updateUsername} />
            </Form.Item>

            <Form.Item
              label='密码'
              name='password'
              rules={[{ required: true, message: '请输入密码！' }]}>
              <Input.Password />
            </Form.Item>
          </Form>
        </div>

        <SelectEditor
          selectedLeftKeys={this.selectedRoleKeys}
          updateSelectedLeftKeys={this.updateSelectedRoleKeys}
          leftDisabledCondition={item =>
            item.type === RoleType.ROLE_NORMAL_USER
          }
          LeftList={this.leftList}
          RightList={this.rightList}
          title={SelectEditorTitle}
        />

        <PendingFooter
          onCancel={this.onCancel}
          onOk={this.onOk}
          processing={this.adding}
        />
      </UserCreatorWrapper>
    )
  }
}
