/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { Layout } from 'antd'
import { history } from '@/utils'
import { Menu } from '@/components/PageLayout/Menu'
import { RouterType } from '@/components/PageLayout/typing'
import { Breadcrumb } from '@/components/PageLayout/Breadcrumb'
import { Header } from './Header'
import { MenuFoldOutlined, MenuUnfoldOutlined } from '@ant-design/icons'
import { useResize } from '@/domain'

const { Content, Sider } = Layout

interface IProps {
  height: string
}

const StyledLayout = styled.div<IProps>`
  > .layout {
    > .header {
      position: fixed;
      z-index: 1;
      width: calc(100% - 60px);
    }

    > .main {
      min-height: {props => props.height || "calc(100vh - 76px)"};

      aside {
        position: fixed;
        height: {props => props.height || "calc(100vh - 76px)"};

        .ant-layout-sider-children {
          padding-bottom: 60px;

          .sider-footer {
            position: absolute;
            bottom: 60px;
            left: 0;
            width: 100%;
            font-size: 14px;
            color: #333333;
            line-height: 22px;
            margin: 20px 0;
            display: flex;
            align-items: center;
            justify-content: center;
            cursor: pointer;

            &:hover {
              color: #001b43;
            }

            > .icon {
              margin-right: 12px;
            }
          }
        }
      }

      .ant-layout-sider-light {
        background-color: #f1f3f6;
        box-shadow: 0 0px 0px 0 rgba(213, 213, 213, 0.86);
      }

      > .body {
        // padding-left: 180px;

        > .toolbar {
          display: flex;
          margin-top: 10px;
          margin-bottom: 10px;
          align-items: center;
          font-size: 16px;
        }
      }
    }
  }
`

type Props = {
  routers: RouterType[]
  children?: React.ReactNode
}

const findSelectedKeys = (menuItems, pathname) => {
  for (let item of menuItems) {
    if (item.children) {
      const selectedKeys = findSelectedKeys(item.children, pathname)
      if (selectedKeys) {
        return [item.key, ...selectedKeys].filter(item => !!item)
      }
    } else if (item.path && pathname.includes(item.path)) {
      return [item.key]
    }
  }
  return null
}
export const CompanyLayout = observer(function CompanyLayout({
  routers,
  children
}: Props) {
  const [rect, ref] = useResize()
  const [collapsed, setCollapsed] = useState(false)
  const [defaultSelectedKeys, setDefaultSelectedKeys] = useState([])
  const currentPath = window.localStorage.getItem('CURRENTROUTERPATH')
  const SliderWidth = 150
  useEffect(() => {
    if (currentPath && currentPath.includes('sys')) {
      setDefaultSelectedKeys([currentPath])
    }
  }, [currentPath])

  function onMenuClick({ item, key, keyPath, selectedKeys, domEvent }) {
    window.localStorage.setItem('CURRENTROUTERPATH', key)
  }
  return (
    <StyledLayout ref={ref} height={rect.height + 'px'}>
      <Layout className='layout'>
        {/* <Header className='header' /> */}
        <Layout className='main'>
          <Sider
            theme='light'
            ref={ref}
            trigger={null}
            collapsible
            width={SliderWidth}
            collapsedWidth={60}
            collapsed={collapsed}
            onCollapse={value => setCollapsed(value)}>
            <Menu
              style={{ flex: 1 }}
              routers={routers}
              menuProps={{
                theme: 'light',
                onSelect: onMenuClick,
                // defaultOpenKeys: [''],
                defaultSelectedKeys: defaultSelectedKeys
              }}
            />
          </Sider>
          <Layout
            className='body'
            style={{ paddingLeft: !collapsed ? SliderWidth : 60 }}>
            <div className='toolbar'>
              {React.createElement(
                collapsed ? MenuUnfoldOutlined : MenuFoldOutlined,
                {
                  className: 'trigger',
                  onClick: () => setCollapsed(!collapsed)
                }
              )}
              <Breadcrumb routers={routers} history={history} />
            </div>
            <Content>{children}</Content>
          </Layout>
        </Layout>
      </Layout>
    </StyledLayout>
  )
})
