import { isRoleConflict } from '@/domain/UserMG/Role'
import { message } from 'antd'
import { observer, useLocalStore } from 'mobx-react'
import * as React from 'react'
import { PendingFooter } from '../../components'
import BasicInfo from '../../User/UserEditor/BasicInfo'
import { UserEditorWrapper } from '../../User/UserEditor/style'
import { checkPerm } from '../../utils'
import SelectEditor from './SelectEditor'
import organization from '@/domain/UserMG/UserOfOrgList'
import { sysConfig } from '@/domain'

type Props = {
  user: any
  onCancel: () => void
  onOk: () => void
}

export const Edit = observer(function Edit({ user, onCancel, onOk }: Props) {
  const state = useLocalStore(() => ({
    isInternal: user.isInternal,
    adding: false,
    updateAdding(flag) {
      this.adding = flag
    },
    name: user.name,
    email: user.email,
    updateEmail(email) {
      this.email = email
    },
    mobile: user.mobile,
    updateMobile(mobile) {
      this.mobile = mobile
    },
    errMessage: {
      email: '',
      mobile: ''
    },
    updateErrorMessage(type, message) {
      this.errMessage[type] = message
    },
    selectedRoleKeys: user.roles,
    updateSelectedRoleKeys(keys, items) {
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

      this.selectedRoleKeys = items.map(i => i.id)
      user.roleNames = items.map(i => i.name)
    }
  }))

  const {
    adding,
    name,
    email,
    errMessage,
    mobile,
    updateEmail,
    updateErrorMessage,
    updateMobile,
    selectedRoleKeys,
    updateSelectedRoleKeys,
    isInternal,
    updateAdding
  } = state

  async function ok() {
    if (user.roleNames.length === 0) {
      message.error('用户角色不能为空')
      return
    }

    if (errMessage.email) {
      message.error(errMessage.email)
      return
    }

    if (errMessage.mobile) {
      message.error(errMessage.mobile)
      return
    }

    updateAdding(true)
    organization
      .update(user.id, {
        mobile,
        email,
        name,
        roles: selectedRoleKeys,
        roleNames: user.roleNames,
        groups: [],
        groupNames: []
      })
      .then(res => {
        if (res.data?.isAskRequest) {
          res.success
            ? message.success(res.message)
            : message.error(res.message)
        } else {
          message.success('用户修改成功')
        }
        onOk()
      })
      .finally(() => {
        updateAdding(false)
        checkPerm()
      })
  }
  return (
    <UserEditorWrapper>
      <BasicInfo
        name={name}
        email={email}
        mobile={mobile}
        updateEmail={updateEmail}
        updateMobile={updateMobile}
        updateError={updateErrorMessage}
      />
      <SelectEditor
        selectedKeys={selectedRoleKeys}
        updateSelectedKeys={updateSelectedRoleKeys}
        title={'角色'}
      />

      <PendingFooter onCancel={onCancel} onOk={ok} processing={adding} />
    </UserEditorWrapper>
  )
})
