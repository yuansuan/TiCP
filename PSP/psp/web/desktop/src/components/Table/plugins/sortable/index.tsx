/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

/**
 * @module SortablePlugin
 * @description use this plugin to sort column
 * when column.sorter is defined, the plugin will be enabled
 */

import * as React from 'react'
import { PluginProps, PluginContext } from '../../@types'
import { SortableHeader } from './SortableHeader'
import { observable, action } from 'mobx'

const PLUGIN_ID = 'sortable-plugin'

export enum SortType {
  default = '',
  asc = 'asc',
  desc = 'desc',
}

export class Model {
  @observable sortKey = ''
  @observable sortType = SortType.default

  @action
  setSortKey = key => {
    this.sortKey = key
  }

  @action
  setSortType = type => {
    this.sortType = type
  }
}

const model = new Model()

export function sortable({ hooks }: PluginProps, { Icon }) {
  hooks.beforeRender.tap(PLUGIN_ID, (context: PluginContext) => {
    const { columns } = context
    columns
      .filter(column => !!column.sorter)
      .map(column => {
        const Header =
          typeof column.header === 'function' ? (
            column.header()
          ) : (
            <span>{column.header}</span>
          )

        column.header = (
          <SortableHeader
            model={model}
            Header={Header}
            Icon={Icon}
            column={column}
          />
        )

        return column
      })

    return context
  })
}
