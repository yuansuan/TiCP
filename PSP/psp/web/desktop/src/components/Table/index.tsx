/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, {
  useMemo,
  useState,
  useCallback,
  useLayoutEffect,
  useEffect,
} from 'react'
import { Cell, HeaderCell, Column } from 'rsuite-table'
import { AsyncParallelHook, SyncWaterfallHook } from 'tapable'
import 'rsuite-table/dist/css/rsuite-table.min.css'
import Icon from '../Icon'
import Chooser from '../Chooser'
import { StyledTable } from './style'
import { TableProps } from './@types'
import { harmonyProps, harmonyColumn, createPluginContext } from './utils'
import { ColumnProps } from 'rsuite-table/lib/Column'
import { CellProps } from 'rsuite-table/lib/Cell'
import plugins from './plugins'

const YSCell: any = ({ render, ...props }) => (
  <Cell key={props?.rowKey} {...(props as CellProps)}>
    {render && render({ ...props })}
  </Cell>
)

const useForceUpdate = () => {
  const [, updateState] = React.useState<any>()
  const forceUpdate = useCallback(() => updateState({}), [])

  return forceUpdate
}

export * from './@types'

export default function YSTable(tableProps: TableProps) {
  const forceUpdate = useForceUpdate()
  const hooks = useMemo(
    () => ({
      init: new AsyncParallelHook(['context']),
      beforeRender: new SyncWaterfallHook(['context']),
      useLayoutEffect: new SyncWaterfallHook(['context']),
      useEffect: new SyncWaterfallHook(['context']),
    }),
    []
  )

  const mergedProps = harmonyProps(tableProps)
  const pluginContext = createPluginContext(mergedProps)

  const [] = useState(() => {
    // apply plugins
    const finalPlugins = [
      ...plugins,
      ...(tableProps.plugins || []).filter(
        p => typeof p === 'object' || typeof p === 'function'
      ),
    ]
    finalPlugins.forEach(plugin => {
      if (typeof plugin === 'function') {
        plugin({ hooks, forceUpdate }, { Icon, Selector: Chooser })
      } else if (plugin.apply) {
        plugin.apply({ hooks, forceUpdate }, { Icon, Selector: Chooser })
      }
    })

    // trigger init hooks
    hooks.init.callAsync(pluginContext, () => {})

    return finalPlugins
  })

  let finalprops: TableProps = hooks.beforeRender.call(pluginContext)
  // plugin may change the columns property, harmony the columns
  if (finalprops) {
    finalprops.columns = finalprops.columns.map(harmonyColumn)
  } else {
    // default
    finalprops = mergedProps
  }

  // mount original data to props
  finalprops.props.data = tableProps.props.data

  // useLayoutEffect
  useLayoutEffect(() => {
    hooks.useLayoutEffect.call(pluginContext)
  })

  // useEffect
  useEffect(() => {
    hooks.useEffect.call(pluginContext)
  })

  return (
    <StyledTable
      {...finalprops.props}
      headerHeight={finalprops.props.headerHeight + 1}>
      {finalprops.columns.map((column, index) => {
        const { props: columnProps, headerClassName, header, cell } = column
        const { props: cellProps, render } = cell

        return (
          <Column key={index} {...(columnProps as ColumnProps)}>
            <HeaderCell {...({} as CellProps)}>
              <div className={headerClassName}>
                {typeof header === 'function' ? header() : <>{header}</>}
              </div>
            </HeaderCell>
            {render ? (
              <YSCell {...cellProps} render={render} />
            ) : (
              <Cell {...(cellProps as CellProps)} />
            )}
          </Column>
        )
      })}
    </StyledTable>
  )
}
