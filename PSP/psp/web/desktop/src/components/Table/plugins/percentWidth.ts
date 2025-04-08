/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { PluginContext, PluginProps, ColumnProps } from '../@types'

const PLUGIN_ID = 'percentWidth-plugin'

export function percentWidth({ hooks }: PluginProps) {
  function formatColumns(columns: ColumnProps[]) {
    columns.forEach(item => {
      if (item.props && typeof item.props.width === 'string') {
        item.props.width = undefined
      }
    })

    return columns
  }

  hooks.beforeRender.tap(PLUGIN_ID, (context: PluginContext) => {
    const { columns, props } = context
    let widthArr = columns.map(item => item.props && item.props.width)

    // format width
    const percentReg = /^[1-9][0-9]?%$|^100%$/
    widthArr = widthArr.map((item, index) => {
      if (typeof item === 'string' && !percentReg.test(item)) {
        const column = columns[index]
        console.warn(
          `PercentWidthPlugin: The ${column.name} column's width(${item}) is invalid.Please use interge or '1-100%'.`
        )
        return undefined
      }

      return item
    })

    const isAdaptive =
      widthArr.findIndex(
        item => typeof item === 'string' || item === undefined
      ) > -1

    if (!isAdaptive) {
      return context
    }

    if (props.width === undefined) {
      context.columns = formatColumns(columns)
      return context
    }

    const exactWidth = (widthArr.filter(
      item => typeof item === 'number'
    ) as number[]).reduce((total, item) => total + item, 0)
    const totalWidth = props.width - exactWidth
    const totalPercent = (widthArr.filter(
      item => typeof item === 'string'
    ) as string[]).reduce((total, item) => total + parseInt(item), 0) as number
    const totalUndefined = widthArr.filter(item => item === undefined).length

    if (totalPercent > 100) {
      console.warn(
        'PercentWidthPlugin: The total percent assigned to column is more than 100.'
      )
    }

    context.columns = columns.map(column => {
      const props = column.props || {}
      const { width } = props
      if (typeof width === 'string') {
        column.props.width = totalWidth * (parseInt(width) / 100)
      } else if (width === undefined) {
        if (totalPercent < 100) {
          column.props = column.props || {}
          column.props.width =
            (totalWidth * ((100 - totalPercent) / 100)) / totalUndefined
        }
      }
      return column
    })

    return context
  })
}
