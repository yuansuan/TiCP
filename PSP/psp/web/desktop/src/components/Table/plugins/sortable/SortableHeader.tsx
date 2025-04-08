/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

/**
 * @module SortablePlugin
 * @description use this plugin to sort column
 * when column.sorter is defined, the plugin will be enabled
 */

import React from 'react'
import styled from 'styled-components'
import { ColumnProps } from '../../@types'
import { observer } from 'mobx-react-lite'
import { Model } from '.'

const StyledLayout = styled.div`
  display: flex;
  line-height: inherit;

  > .sorter {
    margin-left: 5px;
    cursor: pointer;
    display: flex;
    align-items: center;
    color: rgba(0, 0, 0, 0.25);

    &.active,
    &:hover {
      color: ${({ theme }) => theme.primaryColor};
    }
  }
`

enum SortType {
  default = '',
  asc = 'asc',
  desc = 'desc',
}

type Props = {
  Header: React.ReactNode
  Icon: React.ElementType
  column: ColumnProps
  model: Model
}

export const SortableHeader = observer(function SortableHeader({
  Header,
  Icon,
  column,
  model,
}: Props) {
  const { setSortKey, setSortType } = model

  function onSort() {
    if (model.sortKey !== column.dataKey) {
      setSortKey(column.dataKey)
      setSortType(SortType.default)
    }

    let nextSortType
    switch (model.sortType) {
      case SortType.asc: {
        nextSortType = SortType.desc
        break
      }
      case SortType.desc: {
        nextSortType = SortType.default
        break
      }
      case SortType.default: {
        nextSortType = SortType.asc
        break
      }
      default:
        break
    }

    model.setSortType(nextSortType)
    column.sorter({ sortType: nextSortType, sortKey: column.dataKey })
  }

  let icon = <Icon type='table_sort' />
  let active = false
  if (model.sortKey === column.dataKey) {
    if (model.sortType === SortType.asc) {
      icon = <Icon type='table_sort_up' />
      active = true
    } else if (model.sortType === SortType.desc) {
      icon = <Icon type='table_sort_down' />
      active = true
    } else {
      icon = <Icon type='table_sort' />
    }
  }

  return (
    <StyledLayout>
      <div>{Header}</div>
      <span className={`sorter ${active ? 'active' : ''}`} onClick={onSort}>
        {icon}
      </span>
    </StyledLayout>
  )
})
