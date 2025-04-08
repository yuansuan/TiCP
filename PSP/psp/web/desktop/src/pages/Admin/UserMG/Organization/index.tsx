import { observer } from 'mobx-react'
import React, { useEffect, useState } from 'react'

import { Icon } from '@/components'
import { Context, useModel, useStore } from './store'
import { Resizable } from 're-resizable'
import { useLayoutRect } from '@/utils/hooks'
import { UserList } from './UserList'
import { sysConfig, useResize } from '@/domain'
import styled from 'styled-components'
import { Toolbar } from './Operators'
import { Menu } from './Menu'
import { Input, message } from 'antd'
import Button from 'antd/es/button'
import { Modal } from '@/components'
import organization from '@/domain/UserMG/UserOfOrgList'

export const StyledWrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;

  .header {
    margin-bottom: 10px;
  }

  .main {
    display: flex;
    width: 100%;
    margin-top: 4px;
    border: 1px solid ${props => props.theme.borderColor};
    background: #ffffff;

    .menu {
      .resizeBar {
        z-index: 5;
      }
    }
    .mockbar {
      position: relative;
      width: 2px;
      background-color: ${props => props.theme.borderColor};
      z-index: 1;

      > .wrapper {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 10px;
        height: 26px;
        color: #c9c9c9;
        background-color: #eee;
        border-radius: 5px;
        > .icon {
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
        }
      }
    }
    .panel {
      flex: 1;
      display: flex;
      flex-direction: column;
      padding: 20px 20px 0px;

      > .toolbar {
        padding-bottom: 10px;
        display: flex;
        justify-content: space-between;
      }

      > .list {
        flex: 1;
      }
    }
  }
`
export const UserOrg = observer(function OrganizationManagement() {
  const store = useStore()
  const { selectedKeys, userList } = store
  const [width, setWidth] = useState(300)
  const [headerRect, headerRef, headerResize] = useLayoutRect()
  const [rect, ref, resize] = useResize()
  const [fetch, loading] = store.getUserList()

  useEffect(() => {
    headerResize()
    setTimeout(resize, 100)
  }, [])

  function search(e) {
    const { value } = e.target
    store.setSearchKey(value)
  }

  const selectedList = selectedKeys.map(keys => {
    const user = userList?.filter(user => keys === user.id)[0]
    return user
  })

  const haveApprove = selectedList.filter(list => list?.approve_status === 0)

  const userNames = selectedList.map(list => list?.name)

  function active() {
    Modal.showConfirm({
      content: sysConfig.enableThreeMembers
        ? `确认发起启用用户${userNames}申请吗？`
        : `确认启用用户${userNames}吗？`
    }).then(async () => {
      await selectedList.map(user =>
        organization
          .active(user.id, user.name, {
            roles: user.roles,
            roleNames: user.roleNames,
            groups: [],
            groupNames: []
          })
          .then(res => {
            if (res.data) {
              res.success
                ? message.success(res.message)
                : message.error(res.message)
            } else {
              message.success(`启用用户${user.name}成功`)
            }
          })
      )
      fetch()
    })
  }

  function inactive() {
    if (haveApprove.length !== 0 && sysConfig.enableThreeMembers) {
      message.warn(`用户${userNames}中有未完成的审批，请等待审批结束`)
      return
    }

    Modal.showConfirm({
      content: sysConfig.enableThreeMembers
        ? `确认发起禁用用户${userNames}申请吗？`
        : `确认禁用用户${userNames}吗？`
    }).then(async () => {
      await Promise.all(
        selectedList.map(user =>
          organization
            .inactive(user.id, user.name, {
              roles: user.roles,
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
                message.success(`禁用用户${user.name}成功`)
              }
            })
        )
      )
      fetch()
    })
  }

  return (
    <StyledWrapper style={{ height: 'calc(100vh - 228px)' }}>
      <div className='header' ref={headerRef}>
        <Toolbar />
      </div>
      <div
        className='main'
        ref={ref}
        style={{ height: `calc(100% - ${headerRect.height}px)` }}>
        <div className='menu'>
          <Resizable
            handleClasses={{ right: 'resizeBar' }}
            minWidth={120}
            enable={{ right: true }}
            size={{ width, height: '100%' }}
            onResizeStop={(e, direction, ref, d) => {
              setWidth(width + d.width)
              resize()
            }}>
            <Menu />
          </Resizable>
        </div>
        <div className='mockbar'>
          <div className='wrapper'>
            <Icon className='icon' type='drag' />
          </div>
        </div>
        <div className='panel'>
          <div className='toolbar'>
            <div>
              <Button
                style={{ marginRight: 10 }}
                disabled={selectedKeys.length === 0}
                onClick={active}>
                启用
              </Button>
              <Button disabled={selectedKeys.length === 0} onClick={inactive}>
                禁用
              </Button>
            </div>
            <Input.Search
              style={{ width: 200 }}
              maxLength={64}
              placeholder='按登录名称搜索'
              value={store.searchKey}
              onChange={search}
              allowClear
            />
          </div>
          <div className='list'>
            <UserList
              width={rect.width - width - 43}
              height={rect.height - 120}
            />
          </div>
        </div>
      </div>
    </StyledWrapper>
  )
})

export function Organization() {
  const defaultModel = useModel()
  const finalModel = defaultModel

  return (
    <Context.Provider value={finalModel}>
      <UserOrg />
    </Context.Provider>
  )
}
