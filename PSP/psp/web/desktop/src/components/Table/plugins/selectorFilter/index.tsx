/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

/**
 * @module SelectorFilterPlugin
 * @description use this plugin to filter column with Selector
 * when column.selector is defined, the plugin will be enabled
 */

import React from 'react'
import { SelectorFilter } from './SelectorFiler'
import { PluginProps, PluginContext } from '../../@types'

const PLUGIN_ID = 'selectorFilter-plugin'

export function selectorFilter({ hooks }: PluginProps, { Icon, Selector }) {
  hooks.beforeRender.tap(PLUGIN_ID, (context: PluginContext) => {
    const { columns } = context
    columns
      .filter(
        column =>
          !!column.filter &&
          column.filter.items &&
          column.filter.items.length > 0
      )
      .map(column => {
        const originHeader = column.header
        column.header = () => (
          <SelectorFilter
            style={{ lineHeight: `${context.props.headerHeight}px` }}
            header={originHeader}
            {...column.filter}
            Icon={Icon}
            Selector={Selector}
          />
        )
        return column
      })

    return context
  })
}
