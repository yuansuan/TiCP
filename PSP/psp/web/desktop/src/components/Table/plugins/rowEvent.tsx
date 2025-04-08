/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

/**
 * @module RowEventPlugin
 * @description use this plugin to support listen row event
 */

import React from 'react'
import { PluginProps, PluginContext } from '../@types'

const PLUGIN_ID = 'rowEvent-plugin'

export function rowEvent({ hooks }: PluginProps) {
  hooks.beforeRender.tap(PLUGIN_ID, (context: PluginContext) => {
    const { columns, onRow } = context

    if (!onRow) {
      return context
    }

    context.columns = columns.map(column => {
      let originRender =
        column.cell.render || (props => <>{props.rowData[props.dataKey]}</>)

      column.cell.render = (cellProps: any) => {
        let element = originRender(cellProps)
        if (typeof element === 'string') {
          element = <>{element}</>
        }
        return React.cloneElement(
          element,
          onRow(cellProps.rowData, cellProps.rowIndex)
        )
      }

      return column
    })

    return context
  })
}
