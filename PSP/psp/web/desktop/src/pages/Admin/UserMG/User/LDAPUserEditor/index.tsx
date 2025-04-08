import * as React from 'react'
import styled from 'styled-components'
import { RoleList, User, GroupList, UserList } from '@/domain/UserMG'
import { isRoleConflict } from '@/domain/UserMG/Role'
import { Validator, copyText } from '@/utils'
import {
  UserFormSteps,
  UserFormStepsFooter,
  UserForm,
  SelectEditor,
  Loading
} from '../../components'
import LDAPUserPreview from '../LDAPUserPreview'
import { observer } from 'mobx-react'
import { observable, action, computed } from 'mobx'
import { message } from 'antd'
import { sysConfig } from '@/domain'

const Wrapper = styled.div`
  display: flex;
  align-items: center;
  flex-direction: column;
`

interface IProps {
  user: User
  onCancel?: () => void
  onOk?: () => void
}

@observer
export default class LDAPUserEditor extends React.Component<IProps> {
  // 0 第一步，1，第二步 2，完成显示成功结果（这一步不可以倒退）
  state = {
    currentStep: 0,
    resUser: new User(),
    noFetch: false
  }

  @observable user = this.props.user

  @observable messages = {
    name: '',
    display_name: '',
    mobile: '',
    email: ''
  }

  @observable loading = false

  @action
  onChange = (type, value, isLDAP) => {
    this.user[type] = value
  }

  onBlur = (type, value) => {
    if (type === 'name') {
      if (value === '') return '登录名称不能为空'
      if (value.length > 64) return '登录名称不能超过 64 个字符'
      if (!Validator.isValidName(value)) return '登录名称只能包含字母数字下划线'
      return ''
    }
    if (type === 'email') {
      if (value === '') return ''
      if (value.length > 64) return '邮箱长度不能超过 64 个字符'
      if (!Validator.isValidEmail(value)) return '邮箱格式不对'
      return ''
    }
    if (type === 'mobile') {
      if (value === '') return ''
      if (!Validator.isValidPhoneNumber(value)) return '手机格式不对'
      return ''
    }
    if (type === 'display_name') {
      if (value === '') return ''
      if (!Validator.isValidInputName(value))
        return '用户名只能包含字母数字下划线和中文'
      return ''
    }
    return ''
  }

  onNext = () => {
    if (
      [
        this.onBlur('name', this.user.name),
        this.onBlur('email', this.user.email),
        this.onBlur('mobile', this.user.mobile)
      ].every(str => !str)
    ) {
      this.setState({ currentStep: 1 })
    }
  }

  onPre = () => {
    this.setState({ currentStep: 0 })
  }

  onOk = () => {
    if (this.user.roles.length === 0) {
      message.error('用户角色不能为空')
      return
    }
    this.loading = true

    UserList.updateLDAP(this.user.id, {
      enabled: this.user.enabled,
      name: this.user.name,
      email: this.user.email,
      mobile: this.user.mobile,
      roles: this.user.roles,
      roleNames: this.user.roleNames
    })
      .then(res => {
        if (res.data?.isAskRequest) {
          res.success
            ? message.success(res.message)
            : message.error(res.message)

          this.setState({
            currentStep: 2,
            resUser: new User(this.user),
            noFetch: true
          })

          UserList.fetch()
        } else {
          message.success('编辑用户成功')

          const user = res.data.users.filter(u => u.id === this.user.id)

          this.setState({
            currentStep: 2,
            resUser: new User(user[0])
          })
        }
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
        this.loading = false
      })
  }

  onCopy = () => {
    copyText(
      JSON.stringify(this.state.resUser),
      () => {
        message.success('复制成功')
      },
      () => {
        message.error('复制失败')
      }
    )
  }

  @action
  updateSelectedRoleKeys = (keys, items) => {
    if (
      keys.length > 1 &&
      isRoleConflict(keys) &&
      sysConfig.enableThreeMembers
    ) {
      message.error(
        '系统管理员，安全管理员和审计管理员，这三种内置管理员角色只能选择其中之一。'
      )
      return
    }
    this.user.roles = items.map(i => i.id)
    this.user.roleNames = items.map(i => i.name)
  }

  @computed
  get selectedGroupKeys() {
    return this.user.groups
  }

  @computed
  get selectedRoleKeys() {
    return this.user.roles
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
    const { onCancel } = this.props
    const { currentStep } = this.state
    return (
      <Wrapper>
        <Loading
          loading={this.loading}
          message={'编辑用户中...'}
          style={{ height: 700 }}
        />
        {currentStep !== 2 ? (
          <>
            <UserFormSteps current={currentStep} />
            {currentStep === 0 ? (
              <UserForm
                isEdit
                user={this.user}
                messages={this.messages}
                onChange={this.onChange}
                onBlur={this.onBlur}
              />
            ) : (
              <SelectEditor
                selectedLeftKeys={this.selectedRoleKeys}
                updateSelectedLeftKeys={this.updateSelectedRoleKeys}
                LeftList={this.leftList}
                RightList={this.rightList}
                title={{ leftTab: '角色' }}
              />
            )}
          </>
        ) : (
          <div style={{ marginBottom: 50 }}>
            <LDAPUserPreview
              user={this.state.resUser}
              noFetch={this.state.noFetch}
            />
          </div>
        )}
        <UserFormStepsFooter
          isEdit
          current={currentStep}
          onCancel={onCancel}
          onOk={this.onOk}
          onPre={this.onPre}
          onNext={this.onNext}
          onCopy={this.onCopy}
          onClose={this.props.onOk}
        />
      </Wrapper>
    )
  }
}
