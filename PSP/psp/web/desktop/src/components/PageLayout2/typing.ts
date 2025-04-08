/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { RouteProps } from 'react-router-dom'
import { MenuProps as AntMenuProps } from 'antd/es/menu'

export type RouterType = {
  visible?: boolean | (() => boolean)
  component?: () => React.ReactNode
  icon?: React.ReactNode
  selectedIcon?: React.ReactNode
  path?: string
  search?: { [key: string]: any }
  name?: string | (() => string)
  children?: RouterType[]
  isMenu?: boolean
  key?: string
  id?: string
} & Omit<RouteProps, 'component'>

type HistoryType = 'browser' | 'hash'

export type MenuProps = {
  defaultExpandedKeys?: string[]
  routers: RouterType[]
  menuProps?: AntMenuProps
  Menu?: React.ElementType
  historyType?: HistoryType
  Link?: React.ElementType
  history?:
    | HistoryType
    | {
        push?: (path: any) => void
        listen?: (params: any) => void
        location?: any
      }
}

export type Props = {
  routers?: RouterType[]
  title?: string | (() => React.ReactNode)
  logo?: [any, any]
  SiderHeader?: React.ReactNode
  SiderFooter?: React.ReactNode
  HeaderToolbar?: React.ElementType
  showBreadcrumb?: boolean
  children?
} & MenuProps
