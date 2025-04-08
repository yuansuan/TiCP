/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useReducer } from 'react'
import { Layout } from 'antd'
import { MenuFoldOutlined, MenuUnfoldOutlined } from '@ant-design/icons'
import { StyledLayout, StyledBody, StyledContent, StyledHeader } from './style'
import { createStore } from '@/utils/reducer'
import { usePrevious } from '@/utils/reducer'
import Headroom from 'react-headroom'
import Env from './Env'
import { COLLAPSED_WIDTH, SIDER_WIDTH } from './constant'
import { Breadcrumb } from './Breadcrumb'
import { Menu as YSMenu } from './Menu'
import { Props } from './typing'

const store = createStore(function useStore() {
  const menuExpanded = useReducer(Env.reducer, Env.init())

  return {
    menuExpanded,
  }
})
const { Provider, useStore } = store

const { Sider } = Layout

function SelfLayout({
  defaultExpandedKeys = [],
  routers = [],
  children,
  HeaderToolbar,
  title,
  logo = [
    <img src={require('../../assets/images/components/logo_text.svg')} />,
    <img src={require('../../assets/images/components/logo.svg')} />,
  ],
  history = 'browser',
  showBreadcrumb = false,
  SiderFooter,
  SiderHeader,
  ...restProps
}: Props) {
  const {
    menuExpanded: [menuExpanded, dispatch],
  } = useStore()

  // openKeys
  const [openKeys, setOpenKeys] = useState(defaultExpandedKeys)
  const previousOpenedKeys = usePrevious(openKeys)

  function toggleMenu() {
    if (!menuExpanded) {
      setOpenKeys(previousOpenedKeys || [])
    } else {
      setOpenKeys([])
    }

    dispatch({
      type: 'TOGGLE_MENU',
      payload: !menuExpanded,
    })

    // trigger window resize
    setTimeout(() => {
      window.dispatchEvent(new Event('resize'))
    }, 300)
  }

  return (
    <StyledLayout>
      <Sider
        trigger={null}
        collapsible
        width={SIDER_WIDTH}
        collapsedWidth={COLLAPSED_WIDTH}
        collapsed={!menuExpanded}>
        <div className='ys-pageLayout-title'>
          {menuExpanded ? logo[0] : logo[1]}
          {menuExpanded && (
            <div className='text'>
              {typeof title === 'function' ? title() : title}
            </div>
          )}
        </div>
        <div className='ys-pageLayout-menu-container'>
          {SiderHeader && (
            <div className='ys-pageLayout-sider-header'>{SiderHeader}</div>
          )}
          <YSMenu
            style={{ flex: 1 }}
            history={history}
            routers={routers}
            {...{
              ...restProps,
              menuProps: {
                ...restProps.menuProps,
                openKeys,
                onOpenChange: keys => setOpenKeys(keys as string[]),
              },
            }}
          />
          {SiderFooter && (
            <div className='ys-pageLayout-sider-footer'>{SiderFooter}</div>
          )}
        </div>
      </Sider>
      <StyledBody>
        {/* <Headroom
          style={{
            ...(menuExpanded
              ? {
                  paddingLeft: SIDER_WIDTH,
                }
              : {
                  paddingLeft: COLLAPSED_WIDTH,
                }),
          }}> */}
          <StyledHeader
          style={{
            ...(menuExpanded
              ? {
                  paddingLeft: SIDER_WIDTH,
                }
              : {
                  paddingLeft: COLLAPSED_WIDTH,
                })
          }}>
            <div className='ys-pageLayout-toggle'>
              {menuExpanded && <MenuFoldOutlined onClick={toggleMenu} />}
              {!menuExpanded && <MenuUnfoldOutlined onClick={toggleMenu} />}
            </div>
            {showBreadcrumb && (
              <div className='ys-pageLayout-breadcrumb'>
                <Breadcrumb routers={routers} history={history} />
              </div>
            )}
            <div className='ys-pageLayout-toolbar'>
              {HeaderToolbar && <HeaderToolbar />}
            </div>
          </StyledHeader>
        {/* </Headroom> */}
        <StyledContent className={menuExpanded ? 'expanded' : ''}>
          {children}
        </StyledContent>
      </StyledBody>
    </StyledLayout>
  )
}

type PageLayoutType = React.SFC<Props> & {
  useStore: typeof store.useStore
}
const PageLayout: PageLayoutType = (props => {
  return (
    <Provider>
      <SelfLayout {...props} />
    </Provider>
  )
}) as PageLayoutType
PageLayout.useStore = useStore

export default PageLayout
