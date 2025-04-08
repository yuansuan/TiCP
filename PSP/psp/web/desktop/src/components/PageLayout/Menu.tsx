/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { stringify } from 'querystring'
import { Menu as AntMenu } from 'antd'
import { Link as RRDLink } from 'react-router-dom'
import { RouterType, MenuProps } from './typing'
import { COLLAPSED_WIDTH } from './constant'
import { getSelectedKey } from '@/utils'
const { SubMenu, Item } = AntMenu

const StyledLayout = styled.div`
  .ant-menu-item {
    display: flex;
    align-items: center;

    .anticon.ysicon {
      font-size: 24px;
    }
  }
  .ant-menu-submenu-title .anticon.ysicon {
    font-size: 24px;
  }

  .ant-menu-inline-collapsed {
    width: ${COLLAPSED_WIDTH}px;
  }

  .ant-menu {
    height: 100%;
    overflow-y: auto;
    overflow-x: hidden;

    &::-webkit-scrollbar-track {
      box-shadow: inset 0 0 6px rgba(255, 255, 255, 0.3);
      border-radius: 10px;
      background-color: white;
    }

    &::-webkit-scrollbar {
      width: 2px;
      background-color: white;
    }

    &::-webkit-scrollbar-thumb {
      border-radius: 10px;
      box-shadow: inset 0 0 6px rgba(255, 255, 255, 0.3);
      background-color: ${({ theme }) => theme.primaryColor};
    }

    &.ant-menu-dark {
      background-color: #001529;

      .ant-menu-inline.ant-menu-sub {
        background-color: #001940;
      }

      .ant-menu-item {
        &:not(.ant-menu-item-selected) {
          &:hover {
            background-color: #001b43;
          }
        }

        &:active {
          background-color: #0034b4;
        }
      }

      .ant-menu-item,
      .ant-menu-submenu {
        a {
          color: hsla(0, 0%, 100%, 0.65);

          &:hover {
            color: white;
          }
        }

        &.ant-menu-item-selected {
          a {
            color: white;
          }
        }
      }
    }

    &.ant-menu-light {
      background-color: #f1f3f6;

      .ant-menu-item {
        &:not(.ant-menu-item-selected) {
          &:hover {
            background-color: rgba(0, 49, 169, 0.11);
          }
        }
      }

      .ant-menu-item,
      .ant-menu-submenu {
        > .ant-menu {
          background-color: #f1f3f6;
        }

        a {
          color: #333;
        }

        &.ant-menu-item-selected {
          background-color: #0034b4;

          a {
            color: white;
          }
        }
      }
    }
  }

  .ant-menu-item,
  .ant-menu-submenu {
    overflow: hidden;

    a {
      &:hover,
      &:active,
      &:link,
      &:visited {
        text-decoration: none;
      }

      &::before {
        position: absolute;
        top: 0;
        right: 0;
        bottom: 0;
        left: 0;
        background-color: transparent;
        content: '';
      }
    }
  }

  .ant-menu-inline-collapsed > .ant-menu-item,
  .ant-menu-inline-collapsed
    > .ant-menu-item-group
    > .ant-menu-item-group-list
    > .ant-menu-item,
  .ant-menu-inline-collapsed
    > .ant-menu-item-group
    > .ant-menu-item-group-list
    > .ant-menu-submenu
    > .ant-menu-submenu-title,
  .ant-menu-inline-collapsed > .ant-menu-submenu > .ant-menu-submenu-title {
    padding: 0 20px !important;
  }
`

const routerFilter = (item: RouterType) => {
  let name = item.name
  if (typeof item.name === 'function') {
    name = item.name()
  }

  if (name === undefined) {
    return false
  }

  if (item.isMenu === false) {
    return false
  }

  if (typeof item.visible === 'function') {
    return item.visible()
  }

  if (item.visible !== undefined) {
    return !!item.visible
  }

  return true
}

const getKey = (item: RouterType) => {
  let name: string = ''
  if (typeof item.name === 'function') {
    name = item.name()
  } else if (typeof item.name === 'string') {
    name = item.name
  }

  return item.key || item.path || name
}

type Props = MenuProps & {
  style?: React.CSSProperties
}

export const Menu = observer(function Menu({
  style,
  Link = RRDLink,
  historyType = 'hash',
  menuProps,
  Menu,
  routers,
  history
}: Props) {
  const selectedKey = getSelectedKey('hash')
  const [selectedKeys, setSelectedKeys] = useState([])
  useEffect(() => {
    setSelectedKeys([selectedKey])
  }, [selectedKey])

  const getName = (item: RouterType) => {
    if (item.customName) {
      return item.customName()
    }

    if (typeof item.name === 'function') {
      return item.name()
    } else {
      return item.name
    }
  }

  const getIcon = (item: RouterType) => {
    const icon =
      item.path === selectedKey ? item.selectedIcon || item.icon : item.icon
    return icon
  }

  return (
    <StyledLayout style={style}>
      <AntMenu
        mode='inline'
        theme='dark'
        className='menu'
        {...menuProps}
        selectedKeys={selectedKeys}>
        {Menu ? (
          <Menu />
        ) : (
          routers.filter(routerFilter).map(item =>
            item.children && item.children.filter(routerFilter).length > 0 ? (
              <SubMenu
                key={getKey(item)}
                title={
                  <span style={{ display: 'flex', alignItems: 'center' }}>
                    {getIcon(item)}
                    <span
                      style={{
                        width: '100%',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                        overflow: 'hidden'
                      }}>
                      {getName(item)}
                    </span>
                  </span>
                }>
                {item.children.filter(routerFilter).map(sub => (
                  <Item key={getKey(sub)} id={sub.id}>
                    <Link
                      style={{ paddingLeft: 4 }}
                      to={{
                        pathname: sub.path,
                        search: sub.search && `?${stringify(sub.search)}`,
                        state:
                          typeof history === 'object' &&
                          history?.location?.state
                      }}>
                      {getIcon(sub)}
                      {getName(sub)}
                    </Link>
                  </Item>
                ))}
              </SubMenu>
            ) : (
              <Item key={getKey(item)} id={item.id}>
                {getIcon(item)}
                <span>
                  <Link
                    to={{
                      pathname: item.path,
                      search: item.search && `?${stringify(item.search)}`,
                      state:
                        typeof history === 'object' && history?.location?.state
                    }}>
                    {getName(item)}
                  </Link>
                </span>
              </Item>
            )
          )
        )}
      </AntMenu>
    </StyledLayout>
  )
})
