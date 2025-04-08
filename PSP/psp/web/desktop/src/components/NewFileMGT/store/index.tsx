/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { createStore } from '@/utils/reducer'
import { DirectoryTree } from './DirectoryTree'
import { serverFactory } from './common'
import { History } from './History'
import { useLocalStore } from 'mobx-react-lite'
import { BaseDirectory } from '@/utils/FileSystem'
import { useCallback } from 'react'
import { newBoxServer, FileServer } from '@/server'

const widgets = {
  download: null,
  upload: null,
  move: null,
  delete: null,
  refresh: null,
  testFile: null,
  newDir: null,
  history: null,
  edit: null,
  preview: null,
  rename: null,
  'cloud-app': null,
  'custom-column': null,
  record: null
}

export type Store = ReturnType<typeof useModel>

type WidgetType = keyof typeof widgets
type WidgetFactory = (store: Store) => React.ReactNode
type WidgetConfig = (WidgetType | [WidgetType, WidgetFactory])[]

export const defaultWidgets = Object.keys(widgets) as WidgetConfig

type DIR = {
  name: string
  path: string
}

export type ModelProps = {
  widgets: (props: WidgetConfig) => WidgetConfig
  server?: FileServer
  startDIRS?: DIR[]
}

export function useModel(props?: Partial<ModelProps>) {
  return useLocalStore(() => ({
    server: serverFactory(props?.server || newBoxServer),
    dirTree: new DirectoryTree(),
    dir: new BaseDirectory(),
    setDir(dir) {
      this.dir = dir
    },

    async refresh() {},
    dirLoading: false,
    setDirLoading(flag) {
      this.dirLoading = flag
    },
    history: new History(),
    nodeId: undefined,
    setNodeId(nodeId) {
      this.nodeId = nodeId
    },
    get currentNode() {
      const { dirTree, nodeId } = this
      return dirTree.filterFirstNode(item => item.id === nodeId)
    },
    selectedKeys: [],
    setSelectedKeys(keys) {
      this.selectedKeys = [...keys]
    },
    searchKey: undefined,
    setSearchKey(key) {
      this.searchKey = key
    },
    async initDirTree() {
      this.dirTree.setChildren(
        props?.startDIRS || [
          {
            name: 'Home',
            path: '.'
          }
        ]
      )

      const node = this.dirTree.children[0]
      this.setNodeId(node.id)
    },
    useRefresh: function useRefresh(): [() => Promise<void>, boolean] {
      const store = this

      const fetch = useCallback(
        async function fetch() {
          const { nodeId, dirTree, server } = store
          if (!nodeId) {
            return
          }

          const node = dirTree.filterFirstNode(item => item.id === nodeId)
          try {
            store.setDirLoading(true)
            const dir = await server.fetch(node.path)

            // fix race condition
            if (nodeId !== store.nodeId) {
              store.setDirLoading(false)
              return
            }

            store.setDir(dir)
            await server.sync(
              dirTree.filterFirstNode(item => item.path === dir.path)
            )

            store.setSelectedKeys([])
            store.setSearchKey('')
          } finally {
            store.setDirLoading(false)
          }
        },
        [store.nodeId]
      )

      return [fetch, store.dirLoading]
    },
    widgets: props?.widgets ? props.widgets(defaultWidgets) : defaultWidgets,
    getWidget(key: WidgetType) {
      const widget = this.widgets.find(item => {
        if ((Array.isArray(item) && item[0] === key) || item === key) {
          return true
        }
        return false
      })

      if (!widget) {
        return <></>
      }

      if (Array.isArray(widget)) {
        return widget[1](this)
      }

      return null
    },
    isWidgetVisible(key): boolean {
      const widget = this.widgets.find(item => {
        if ((Array.isArray(item) && item[0] === key) || item === key) {
          return true
        }
        return false
      })

      return !!widget
    }
  }))
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
