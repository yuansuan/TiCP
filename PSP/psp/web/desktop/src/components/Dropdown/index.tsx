/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Dropdown as AntDropdown, Menu } from 'antd'
import { DownOutlined } from '@ant-design/icons'
import { StyledA } from './style'

interface IMenuContent {
  onClick?: () => void
  disabled?: boolean
  key?: string
  children: React.ReactNode
}

interface IProps {
  menuContentList: IMenuContent[]
}

export const Dropdown = (props: IProps) => {
  const { menuContentList } = props

  return (
    <AntDropdown
      placement='bottomLeft'
      overlay={
        <Menu>
          {menuContentList.map(content => (
            <Menu.Item
              disabled={content.disabled}
              onClick={content.onClick}
              key={content.key || Math.random()}>
              {content.children}
            </Menu.Item>
          ))}
        </Menu>
      }>
      <StyledA className='ant-dropdown-link' onClick={e => e.preventDefault()}>
        <DownOutlined />
      </StyledA>
    </AntDropdown>
  )
}
