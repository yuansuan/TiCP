/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

/**
 * @module RowEventPlugin
 * @description use this plugin to style default cell
 */

import { PluginContext, PluginProps } from '../@types'

const PLUGIN_ID = 'placeholder-plugin'

export function placeholder({ hooks }: PluginProps) {
  hooks.beforeRender.tap(PLUGIN_ID, (context: PluginContext) => {
    const { columns } = context

    context.columns = columns.map(column => {
      const { render } = column.cell
      if (!render) {
        column.cell.render = props => {
          const value = props.rowData[props.dataKey]
          return [undefined, ''].includes(value) ? '--' : value
        }
      }

      return column
    })

    return context
  })
}
