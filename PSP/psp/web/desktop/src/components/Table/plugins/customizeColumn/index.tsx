/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

/**
 * @module ComposablePlugin
 * @description use this plugin to custom compose column
 * when table.props.tableId is defined, the plugin will be enabled
 */

import React from 'react'
import { PluginProps, PluginContext } from '../../@types'
import { observable, action } from 'mobx'
import { Toolbar } from './Toolbar'
import { Observer } from 'mobx-react-lite'

const STORAGE_KEY = 'TABLE_CUSTOM_SETTINGS'
const PLUGIN_ID = 'customizable-plugin'

class Config {
  @observable config: Array<{ key: string; active: boolean }> = []
  @action
  setConfig = (config: Array<{ key: string; active: boolean }>) => {
    this.config = config
  }
}

export function customizeColumn({ hooks, forceUpdate }: PluginProps, { Icon }) {
  // use observable state to proxy Column
  const state = new Config()

  function getConfig(context) {
    if (!context) {
      return []
    }

    return context.defaultConfig.map(item => {
      if (typeof item === 'string') {
        return {
          key: item,
          active: true,
        }
      } else if (Array.isArray(item)) {
        return {
          key: item[0],
          active: item[1],
        }
      }
      return {
        ...item,
        active: item.active === undefined ? true : !!item.active,
      }
    })
  }

  function setConfig(tableId, config) {
    state.setConfig(config)

    // persist the settings
    let rootStorage = JSON.parse(localStorage.getItem(STORAGE_KEY))
    rootStorage[tableId].config = config
    localStorage.setItem(STORAGE_KEY, JSON.stringify(rootStorage))

    forceUpdate()
  }

  hooks.init.tap(PLUGIN_ID, (context: PluginContext) => {
    const { tableId } = context

    if (!tableId) {
      return context
    } else if (!context.defaultConfig) {
      console.error(`Table:${PLUGIN_ID} need the defaultConfig prop`)
      return context
    }

    // define root storage
    let rootStorage = null
    try {
      rootStorage = JSON.parse(localStorage.getItem(STORAGE_KEY))
      if (!rootStorage) {
        throw new Error(`${STORAGE_KEY} is not defined`)
      }
    } catch (err) {
      rootStorage = {}
      localStorage.setItem(STORAGE_KEY, JSON.stringify(rootStorage))
    }

    // use defaultConfig to define tableId storage
    // use columns.dataKeys to generate hash
    const hash = getConfig(context)
      .map(item => item.key)
      .sort()
      .join('')
    if (!rootStorage[tableId]) {
      rootStorage[tableId] = {
        config: getConfig(context),
        hash,
      }
    } else if (rootStorage[tableId].hash !== hash) {
      // if dataKeys changed, freshen the config
      rootStorage[tableId] = {
        config: getConfig(context),
        hash,
      }
    }
    localStorage.setItem(STORAGE_KEY, JSON.stringify(rootStorage))

    // update state
    setConfig(context.tableId, rootStorage[tableId].config)

    return context
  })

  hooks.beforeRender.tap(PLUGIN_ID, (context: PluginContext) => {
    const { columns } = context

    if (!context.tableId || !context.defaultConfig) {
      return context
    }

    // flat the Array<Object> to Array<string>
    const flatConfig = state.config
      .filter(item => item.active)
      .map(item => item.key)

    // filter and order the column which is active
    let finalColumns = columns
      .filter(item => !item.dataKey || flatConfig.includes(item.dataKey))
      .sort(
        (prev, next) =>
          flatConfig.findIndex(item => item === prev.dataKey) -
          flatConfig.findIndex(item => item === next.dataKey)
      )
    // slot toolbar to last column's header
    const lastColumn = finalColumns[finalColumns.length - 1]
    if (lastColumn) {
      const Header =
        typeof lastColumn.header === 'function' ? (
          lastColumn.header()
        ) : (
          <span>{lastColumn.header}</span>
        )

      lastColumn.header = (
        <Observer>
          {() => (
            <>
              {Header}
              <Toolbar
                Icon={Icon}
                setConfig={config => setConfig(context.tableId, config)}
                config={state.config}
                columns={columns.filter(item => item.dataKey)}
              />
            </>
          )}
        </Observer>
      )
    }

    context.columns = finalColumns
    return context
  })
}
