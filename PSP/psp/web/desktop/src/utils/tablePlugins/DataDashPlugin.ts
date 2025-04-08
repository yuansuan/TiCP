/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

const replaceByDash = value => ([undefined, ''].includes(value) ? '--' : value)

import * as React from 'react'

export class DataDashPlugin extends React.Component {
  // plugin id
  id = 'data-dash-plugin'
  options = null

  constructor(props?) {
    super(props)

    this.options = props
  }

  apply = table => {
    table.hooks.beforeRender.tap(this.id, context => {
      const { columns, props } = context.props
      if (columns?.length === 0) return
      context.props.columns = columns?.map(column => {
        const { render } = column.cell

        // render 方法不存在
        if (!render) {
          column.cell.render = props => {
            const value = props.rowData[props.dataKey]
            return replaceByDash(value)
          }
        } else {
          // render 方法存在，包装 render 方法
          let originRender = render
          column.cell.render = props => {
            props.rowData[props.dataKey] = replaceByDash(
              props.rowData[props.dataKey]
            )
            return originRender.call(column.cell, props)
          }
        }

        return column
      })

      return context
    })
  }
}
