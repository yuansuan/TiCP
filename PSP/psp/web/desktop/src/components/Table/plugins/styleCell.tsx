/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

/**
 * @module RowEventPlugin
 * @description use this plugin to style default cell
 */

import * as React from 'react'
import { PluginContext, PluginProps } from '../@types'
import styled from 'styled-components'

const StyledCell = styled.span`
  display: inline-block;
  width: calc(100% - 20px);
  padding: 0 4px;
  vertical-align: middle;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: ${props => props['data-line-height']}px;
`
const StyledHeaderCell = styled.span`
  display: inline-block;
  width: calc(100% - 20px);
  padding: 0 4px;
  vertical-align: middle;
  font-family: 'PingFangSC-Regular';
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: ${props => props['data-line-height']}px;
`

const PLUGIN_ID = 'styleCell-plugin'

export function styleCell({ hooks }: PluginProps) {
  hooks.beforeRender.tap(PLUGIN_ID, (context: PluginContext) => {
    const { columns, props } = context
    const rowExpandedHeight = props?.rowExpandedHeight || 100

    context.columns = columns.map(column => {
      // style header
      const { header } = column
      column.header =
        typeof header === 'function' ? (
          <div
            style={{ lineHeight: `${props.headerHeight}px`, padding: '0 4px' }}>
            {header()}
          </div>
        ) : (
          <StyledHeaderCell data-line-height={props.headerHeight}>
            {header}
          </StyledHeaderCell>
        )

      // style cell
      const { render: originRender } = column.cell
      column.cell.render = props => {
        const renderRes = originRender && originRender(props)
        const expanded = !!props?.expanded
        const lineHeight =
          expanded && !props.hasChildren
            ? props.height - rowExpandedHeight
            : props.height

        return typeof renderRes === 'object' ? (
          <div style={{ lineHeight: `${lineHeight}px`, padding: '0 4px' }}>
            {renderRes}
          </div>
        ) : (
          <StyledCell
            data-line-height={lineHeight}
            title={props.rowData[props.dataKey]}>
            {renderRes === undefined ? props.rowData[props.dataKey] : renderRes}
          </StyledCell>
        )
      }

      return column
    })

    return context
  })
}
