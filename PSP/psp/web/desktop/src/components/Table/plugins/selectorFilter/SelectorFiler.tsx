/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Popover, Tooltip } from 'antd'
import { useLocalStore, useObserver } from 'mobx-react-lite'
import styled from 'styled-components'

const StyledLayout = styled.div`
  display: flex;
  line-height: inherit;

  > .selector {
    margin-left: 4px;
    display: flex;
    align-items: center;

    > .anticon {
      cursor: pointer;
      color: rgba(0, 0, 0, 0.25);

      &.active,
      &:hover {
        color: ${({ theme }) => theme.primaryColor};
      }
    }
  }
`

interface IProps {
  style?: React.CSSProperties
  Selector: any
  Icon: any
  header: any
  items: any[]
  onChange?: (
    selectedKeys: string[],
    info?: { node?: any; checked: boolean }
  ) => void
  selectedKeys?: string[]
  updateSelectedKeys?: (keys: string[]) => void
  searchable?: boolean
}

export function SelectorFilter(props: IProps) {
  const state = useLocalStore(() => ({
    _selectedKeys: [],
    visible: false,
    updateVisible(visible) {
      this.visible = visible
    },
    updateSelectedKeys(keys) {
      this._selectedKeys = keys
    },
  }))

  function updateSelectedKeys(keys) {
    const { selectedKeys, updateSelectedKeys } = props

    if (selectedKeys) {
      selectedKeys && updateSelectedKeys(keys)
    } else {
      state.updateSelectedKeys(keys)
    }
  }

  function onChange() {
    const { onChange } = props
    onChange && onChange(props.selectedKeys || state._selectedKeys)
  }

  const { header, items, Selector, Icon, style, searchable } = props
  const Header = typeof header === 'function' ? header() : <span>{header}</span>

  return useObserver(() => {
    const selectedKeys = props.selectedKeys || state._selectedKeys
    const active = selectedKeys.length > 0
    const selectedNames = selectedKeys
      .map(selectedKey =>
        props.items
          .filter(({ key }) => selectedKey === key)
          .map(item => item.name)
      )
      .join(', ')

    return (
      <StyledLayout style={style}>
        <div>{Header}</div>
        <div className='selector'>
          <Popover
            placement='bottom'
            content={
              <Selector
                items={items}
                selectedKeys={selectedKeys}
                updateSelectedKeys={updateSelectedKeys}
                onChange={onChange}
                searchable={
                  searchable === undefined
                    ? items && items.length >= 20
                    : searchable
                }
              />
            }
            trigger='click'
            visible={state.visible}
            onVisibleChange={state.updateVisible}>
            {active ? (
              <Tooltip title={selectedNames} placement='top'>
                <Icon className='active' type='table_filter' />
              </Tooltip>
            ) : (
              <Icon type='table_filter' />
            )}
          </Popover>
        </div>
      </StyledLayout>
    )
  })
}
