import { observer } from 'mobx-react'
import styled from 'styled-components'
import { message, Switch, Tooltip } from 'antd'
import React from 'react'
import { currentUser, sysConfig } from '@/domain'
import { Modal } from '@/components'
import organization from '@/domain/UserMG/UserOfOrgList'
import { useStore } from '../store'

export const StateWrapper = styled.div`
  margin: 0 10px;
`

type Props = {
  user: any
}

export const State = observer(function State({ user }: Props) {
  const store = useStore()
  const [fetch, loading] = store.getUserList()
  const disableOpt = currentUser.id === user?.id || user?.isInternal
  const isActive = user.enabled

  function active() {
    if (user.approve_status === 0 && sysConfig.enableThreeMembers) {
      message.warn(`用户${user.name}有未完成的审批，请等待审批结束`)
      return
    }

    Modal.showConfirm({
      content: sysConfig.enableThreeMembers
        ? `确认发起启用用户${user.name}申请吗？`
        : `确认启用用户${user.name}吗？`,
    }).then(() =>
      organization
        .active(user.id, user.name, {
          roles: user.roles,
          roleNames: user.roleNames,
          groups: [],
          groupNames: [],
        })
        .then(res => {
          fetch()
          if (res.data?.isAskRequest) {
            res.success
              ? message.success(res.message)
              : message.error(res.message)
          } else {
            message.success('启用用户成功')
          }
        })
    )
  }
  function inactive() {
    if (user.approve_status === 0 && sysConfig.enableThreeMembers) {
      message.warn(`用户${user.name}有未完成的审批，请等待审批结束`)
      return
    }

    Modal.showConfirm({
      content: sysConfig.enableThreeMembers
        ? `确认发起禁用用户${user.name}申请吗？`
        : `确认禁用用户${user.name}吗？`,
    }).then(() =>
      organization
        .inactive(user.id, user.name, {
          roles: user.roles,
          roleNames: user.roleNames,
          groups: [],
          groupNames: [],
        })
        .then(res => {
          fetch()
          if (res.data?.isAskRequest) {
            res.success
              ? message.success(res.message)
              : message.error(res.message)
          } else {
            message.success('禁用用户成功')
          }
        })
    )
  }
  function onChange(checked) {
    if (checked) {
      active()
    } else {
      inactive()
    }
  }
  return (
    <StateWrapper>
      <Tooltip title={user?.isInternal ? '内置系统账号不能被禁用或启用' : ''}>
        <Switch
          disabled={disableOpt}
          checked={isActive}
          onChange={onChange}
          size='small'
        />
      </Tooltip>
      {user.enabled ? '  启用' : '  禁用'}
    </StateWrapper>
  )
})
