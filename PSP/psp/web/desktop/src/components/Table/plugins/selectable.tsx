/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

/**
 * @module SelectablePlugin
 * @description use this plugin to select row
 * when table.props.rowSelection is defined, the plugin will be enabled
 */

import React from 'react'
import { Checkbox } from 'antd'
import { PluginProps, PluginContext } from '../@types'
import { observable, action } from 'mobx'
import { Observer } from 'mobx-react-lite'

const PLUGIN_ID = 'selectable-plugin'

class State {
  @observable selectedKeys = []
  @action
  selectAll = keys => {
    this.selectedKeys = keys
  }
  @action
  selectInvert = () => {
    this.selectedKeys = []
  }
  @action
  select = (rowKey, checked) => {
    let selectedKeys = this.selectedKeys
    if (checked) {
      selectedKeys = [...selectedKeys, rowKey]
    } else {
      const index = selectedKeys.findIndex(item => item === rowKey)
      selectedKeys.splice(index, 1)
    }

    this.selectedKeys = selectedKeys
  }
}

export function selectable({ hooks }: PluginProps) {
  const state = new State()

  hooks.init.tap(PLUGIN_ID, (context: PluginContext) => {
    if (!context.rowSelection) {
      return context
    }

    const {
      rowSelection: { defaultSelectedKeys: keys = [] },
    } = context
    state.selectedKeys = [...keys]

    return context
  })

  hooks.useLayoutEffect.tap(PLUGIN_ID, (context: PluginContext) => {
    if (!context.rowSelection) {
      return
    }

    const {
      rowSelection: { selectedRowKeys, selectedKeys },
    } = context
    // reset selectedKeys by props
    const keys = selectedKeys || selectedRowKeys
    if (keys) {
      state.selectedKeys = [...keys]
    }
  })

  hooks.beforeRender.tap(PLUGIN_ID, (context: PluginContext) => {
    if (!context.rowSelection) {
      return context
    }

    const {
      rowSelection,
      columns,
      props: { data: dataSource, rowKey },
    } = context

    // set default width
    const props = { width: 50, ...rowSelection.props }

    const onChange = () =>
      rowSelection.onChange && rowSelection.onChange(state.selectedKeys)

    const onSelectAll = keys => {
      state.selectAll(keys)
      onChange()
      rowSelection.onSelectAll && rowSelection.onSelectAll(keys)
    }
    const onSelectInvert = () => {
      state.selectInvert()
      onChange()
      rowSelection.onSelectInvert && rowSelection.onSelectInvert()
    }
    const onSelect = (rowKey, checked) => {
      state.select(rowKey, checked)
      onChange()
      rowSelection.onSelect && rowSelection.onSelect(rowKey, checked)
    }

    columns.unshift({
      props: { ...props, width: 50 },
      header: () => (
        <Observer>
          {() => {
            const { selectedKeys } = state
            return (
              <Checkbox
                checked={
                  dataSource.length > 0 &&
                  dataSource.length === selectedKeys.length
                }
                indeterminate={
                  selectedKeys.length > 0 &&
                  selectedKeys.length < dataSource.length
                }
                onChange={e => {
                  if (e.target.checked) {
                    onSelectAll &&
                      onSelectAll(dataSource.map(item => item[rowKey]))
                  } else {
                    onSelectInvert && onSelectInvert()
                  }
                }}
              />
            )
          }}
        </Observer>
      ),
      cell: {
        render: props => (
          <Observer>
            {() => {
              const { selectedKeys } = state

              return (
                <Checkbox
                  checked={selectedKeys.includes(props.rowKey)}
                  onChange={e => {
                    onSelect(props.rowKey, e.target.checked)
                  }}
                />
              )
            }}
          </Observer>
        ),
      },
    })

    const originRowClassName = context.props.rowClassName
    context.props.rowClassName = rowData => {
      const classNames = []
      if (typeof originRowClassName === 'function') {
        classNames.push(originRowClassName(rowData))
      } else if (typeof originRowClassName === 'string') {
        classNames.push(originRowClassName)
      }

      const { selectedKeys = [] } = rowSelection
      if (selectedKeys.includes(rowData && rowData[rowKey])) {
        classNames.push('rs-table-row-selected')
      }

      return classNames.filter(item => !!item).join(' ')
    }

    return context
  })
}
