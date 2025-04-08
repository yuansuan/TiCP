/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Dropdown, Menu } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { currentUser, env } from '@/domain'
import { Http } from '@/utils'
import { Modal, Icon } from '@/components'
import { CaretDownFilled, CaretUpFilled } from '@ant-design/icons'
import PersonalSetting from '@/pages/PersonalSetting'
import { Hover } from '@/components'
import { buryPoint } from '@/utils'
const sessionkey = 'LOGIN_PASSWD_EXPIRED_NOTIFICATION'

const StyledMenu = styled(Menu)`
  .ant-dropdown-menu-item,
  .ant-dropdown-menu-submenu-title {
    padding: 8px 20px;
    color: #666;

    a {
      color: #666;
    }
  }
  .ant-dropdown-menu-item {
    > * {
      display: flex;
      align-items: center;
    }

    .anticon.ysicon {
      margin-right: 8px;
      font-size: 16px;
      color: #666;

      &.hovered {
        color: ${props => props.theme.primaryColor};

        > span {
          background: white;
          display: inline-block;
        }
      }
    }
  }
`

const StyledUserInfo = styled.div`
  display: flex;
  align-items: center;
  padding: 0 10px;
  /* position: absolute;
  left: 20px;
  top: 12px; */

  > .username {
    margin-left: 4px;
    max-width: 150px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
`
type Props = {
  type?: 'portal' | 'inside'
}
export const UserInfo = observer(function UserInfo({ type = 'inside' }: Props) {
  const state = useLocalStore(() => ({
    visible: false,
    setVisible(flag) {
      this.visible = flag
    },
    workOrderCount: 0,
    setWorkOrderCount(count) {
      this.workOrderCount = count
    }
  }))
  const { visible } = state

  async function logout() {
    await Modal.showConfirm({
      title: '退出登录',
      content: '确认要退出登录吗？'
    })
    await Http.post('/auth/logout')
    // don't use history.push which will be intercepted
    location.reload()
    localStorage.removeItem('userId')
    localStorage.removeItem('SystemPerm')
    localStorage.removeItem('CURRENTROUTERPATH')
    localStorage.removeItem('GlobalConfig')
    sessionStorage.removeItem(sessionkey)
    document.cookie = 'access_token=; expires=Thu, 01 Jan 1970 00:00:01 GMT;'
    document.cookie = 'refresh_token=; expires=Thu, 01 Jan 1970 00:00:01 GMT;'
  }

  if (!currentUser.id) return null

  function onMenuClick(props) {
    const { key } = props
    buryPoint({
      category: '个人信息',
      action: {
        person: '个人设置',
        account: '账户管理',
        handbook: '用户手册',
        logout: '退出登录'
      }[key]
    })
    state.setVisible(false)
  }

  const personSetting = () => {
    Modal.show({
      title: '个人设置',
      content: <PersonalSetting />
    })
  }
  return (
    <Dropdown
      visible={state.visible}
      onVisibleChange={visible => state.setVisible(visible)}
      overlay={
        <StyledMenu onClick={onMenuClick}>
          {currentUser?.isLdapEnabled && (
            <Menu.Item key='person'>
              <Hover
                render={hovered => (
                  <div onClick={personSetting}>
                    <Icon
                      className={hovered ? 'hovered' : ''}
                      type='personal_setting_default'
                    />
                    <span>个人设置</span>
                  </div>
                )}
              />
            </Menu.Item>
          )}
          <Menu.Item key='logout' onClick={logout}>
            <Hover
              render={hovered => (
                <div>
                  <Icon
                    className={hovered ? 'hovered' : ''}
                    type='logout_default'
                  />
                  <span>退出登录</span>
                </div>
              )}
            />
          </Menu.Item>
        </StyledMenu>
      }
      placement='bottomCenter'>
      <StyledUserInfo id='ys_header_user_menu'>
        <Icon
          type={
            state.visible ? 'global_setting_active' : 'global_setting_default'
          }
          className={`${state.visible ? 'active' : ''}`}>
          <span></span>
        </Icon>
        <div className='username' title={currentUser.name}>
          {currentUser.name || '--'}
        </div>
        {visible && (
          <CaretUpFilled style={{ marginLeft: 4, fontSize: '14px' }} />
        )}
        {!visible && (
          <CaretDownFilled style={{ marginLeft: 4, fontSize: '14px' }} />
        )}
      </StyledUserInfo>
    </Dropdown>
  )
})
