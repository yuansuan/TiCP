/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { TableProps, ColumnProps, PluginContext } from './@types'
import cloneDeep from 'lodash/cloneDeep'

export function harmonyColumn(item: ColumnProps) {
  const dataKey = item.dataKey || item.cell?.props?.dataKey

  return {
    ...item,
    name: item.name || (typeof item.header === 'string' ? item.header : ''),
    dataKey,
    cell: item.cell || {
      props: {
        dataKey,
      },
      render: undefined,
    },
  }
}

export function harmonyProps(tableProps: TableProps) {
  const { data: originalData = [] } = tableProps.props
  // delete the data property to avoid inefficient clone
  Reflect.deleteProperty(tableProps.props, 'data')
  // deep clone props
  const mergedProps = cloneDeep(tableProps)
  // recover the data property
  tableProps.props.data = originalData

  // just transmit minimal subset of data as read-only information to avoid arbitrary mutation
  const { rowKey } = mergedProps.props
  mergedProps.props.data = originalData.map(item => ({
    ...(rowKey ? { [rowKey]: item[rowKey] } : null),
  }))

  // define the table row's default height
  mergedProps.props.rowHeight = mergedProps.props.rowHeight || 54
  mergedProps.props.headerHeight = mergedProps.props.headerHeight || 54
  mergedProps.props.height = mergedProps.props.height || 250

  // define default empty
  mergedProps.props.renderEmpty =
    mergedProps.props.renderEmpty ||
    (() => (
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}>
        <img
          style={{ height: '160px', margin: '10px 0' }}
          src={require('../assets/images/nodata.png')}
        />
      </div>
    ))

  // harmony columns
  mergedProps.columns = mergedProps.columns.map(harmonyColumn)

  return mergedProps
}

export function createPluginContext(props: TableProps): PluginContext {
  return cloneDeep(props)
}
