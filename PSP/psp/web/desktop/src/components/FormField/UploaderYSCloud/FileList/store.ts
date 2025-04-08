/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { useEffect } from 'react'
import { Http } from '@/utils'
import difference from 'lodash/difference'
import { message } from 'antd'
import { action, computed, observable } from 'mobx'
import { formatByte } from '@/utils/Validator'
import { statusMap } from '@/domain/Uploader/Task'

export type Props = {
  fileList: any[]
  deleteAction: (path: string) => void
  setMainFileKeysAction: (keys: string[]) => void
}

export function useModel({
  fileList,
  deleteAction,
  setMainFileKeysAction,
}: Props) {
  useEffect(() => {
    store.update({
      fileList: fileList.map(v => {
        const newFile = new FileTreeFile(v)
        const oldFile = store.fileList.find(f => f.path === v.path)
        // 有覆盖上传的情况则关闭展开
        oldFile &&
          store.update({
            expandedKeys: store.expandedKeys.filter(k => k !== oldFile.path),
          })
        return newFile
      }),
    })
  }, [fileList])

  const store = useLocalStore(() => ({
    loading: false,
    fileList,
    update(data: any) {
      Object.assign(this, data)
    },
    expandedKeys: [],
    onExpand: async keys => {
      const oldExpandedKeys = store.expandedKeys.slice()
      store.update({
        expandedKeys: [...new Set(keys)],
      })

      const diffList = difference(keys, oldExpandedKeys)
      if (diffList.length === 0) return

      const clickedKey = diffList[0]
      if (clickedKey[0] !== '/') {
        message.warn('请等待上传结束')
        return
      }

      const item = bfsSearch(
        store.fileList,
        item => item.path === clickedKey
      )?.[0]
      if (!item) {
        return
      }

      if (!item.children || item.children.length <= 0) {
        store.update({
          loading: true,
        })
      }

      const { data } = await Http.get('/file/list', {
        params: { path: clickedKey },
      }).catch(() => {
        store.update({
          loading: false,
        })
      })

      if (Array.isArray(data?.files)) {
        let mainKeysOfItem = []
        if (Array.isArray(item.children)) {
          mainKeysOfItem = item.children.map(v => (v.isMain ? v.path : null))
        }
        item.children = data.files.map(
          f =>
            new FileTreeFile({
              ...f,
              isDir: f.is_dir,
              parent: item,
              status: 'done',
              isMain: mainKeysOfItem.includes(f.path),
            })
        )
        store.update({
          loading: false,
        })
      }
    },

    fetch: async () => {},
    // create
    createFolder: id => {},

    // delete
    deleteNode: async (path, isServerDel = true) => {
      {
        store.update({
          loading: true,
        })

        if (isServerDel) {
          await Http.post(
            '/file/delete',
            { paths: [path] },
            { formatErrorMessage: msg => `删除失败：${msg}` }
          ).catch(() =>
            store.update({
              loading: false,
            })
          )
        }

        deleteAction(path)
        store.update({
          loading: false,
        })

        const item = bfsSearch(store.fileList, item => item.path === path)?.[0]
        if (!item) {
          return
        }
        message.success(`已删除 ${item.name}`)

        if (item.parent) {
          item.parent.children = item.parent.children.filter(
            v => v.path !== path
          )
        } else {
          store.fileList = store.fileList.filter(v => v.path !== path)
        }
      }
    },

    // upload
    upload: (id: string, isUploadingDirectory = false) => {},

    setMain: (path, bool) => {
      // psp作业输出路径的原因批量提交可能导致批量作业回传的文件覆盖，因此只设置一个主文件单次提交
      if (!bool) {
        bfsSearch(store.fileList, item => {
          item.isMain = false
          return false
        })
        setMainFileKeysAction([])
        return
      }

      let mainItem = null
      bfsSearch(store.fileList, item => {
        if (item.path !== path) {
          item.isMain = false
        }

        if (item.path === path) {
          mainItem = item
        }
        return false
      })
      let rootNode = mainItem.parent ?? mainItem
      while (rootNode.parent) {
        rootNode = rootNode.parent
      }
      const rootPath = rootNode.path.match(/\/([^/]*$)/)?.[1]
      const finalPath = rootPath + mainItem.path.split(rootPath)?.[1]

      setMainFileKeysAction([finalPath])
    },

    // get mainFileKeys(): string[] {
    //   const mainFiles = bfsSearch(store.fileList, item => item.isMain)
    //   return mainFiles.map(v => v.path)
    // },
  }))

  return store
}

function bfsSearch(
  list: Array<FileTreeFile>,
  pickFunc: (item: FileTreeFile) => boolean
): Array<FileTreeFile> {
  const queue = []
  queue.push(...list)
  const results = []
  while (queue.length > 0) {
    const item = queue.shift()
    if (pickFunc.apply(item, [item])) {
      results.push(item)
    } else {
      if (Array.isArray(item.children)) {
        queue.push(...item.children)
      }
    }
  }

  return results
}

class FileTreeFile {
  @observable name: string = undefined
  @observable parent: FileTreeFile = undefined
  @observable children: FileTreeFile[] = undefined
  @observable size: number = undefined
  @observable path: string = undefined
  @observable status: string = ''
  @observable from: string = ''
  @observable isMain: boolean = false
  @observable isDir: boolean = false

  @computed
  get displaySize() {
    return formatByte(this.size)
  }

  @computed
  get displayStatus() {
    return statusMap[this.status]
  }

  @computed
  get displayFrom() {
    let file = this
    while (file.parent) {
      file = file.parent as any
    }
    return file.from === 'local' ? '本地' : '服务器'
  }

  constructor(props: Partial<FileTreeFile>) {
    this.update(props)
  }

  @action
  update(props: Partial<FileTreeFile>) {
    Object.assign(this, props)
  }
}

const state = createStore(useModel)

export const Provider = state.Provider
export const Context = state.Context
export const useStore = state.useStore
