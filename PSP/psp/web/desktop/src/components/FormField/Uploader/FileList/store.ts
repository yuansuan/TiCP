/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { transaction } from 'mobx'
import { useEffect } from 'react'
import { Http } from '@/utils'
import { message } from 'antd'
import { toJS } from 'mobx'
import { currentFileList, FileTreeFile } from './Files'
import eventEmitter from '@/utils/EventEmitter'

export type Props = {
  expandedKeys: string[]
  onExpandedKeysChange: (keys: string[]) => void
  fileList: any[]
  deleteAction: (path: string) => void
  onFileListChange: (files: any) => void
  setMainFileKeysAction: (keys: any[]) => void
  mainFilePath: string
  jobDirPath: string
  beforeUploadLocalFile?: (params: any) => Promise<any>
  uploadLocalFile?: (params: any, isDir: boolean) => void
  uploadServerFile?: (
    files: any[],
    master?: string,
    targetPath?: string
  ) => void
  onScroll?: (e) => void
}

export function useModel({
  fileList,
  deleteAction,
  expandedKeys,
  onExpandedKeysChange,
  setMainFileKeysAction,
  mainFilePath,
  jobDirPath,
  onFileListChange,
  beforeUploadLocalFile,
  uploadLocalFile,
  uploadServerFile,
  onScroll
}: Props) {
  useEffect(() => {
    // 采用批量更新，全部完成之后再通知
    transaction(() => {
      let oldFilePathMap = new Map(currentFileList.files.map(f => [f.path, f]))

      let newFiles = fileList.map(v => {
        const newFile = new FileTreeFile(v)
        // // 有覆盖上传的情况则关闭展开
        // oldFile &&
        //   store.update({
        //     expandedKeys: store.expandedKeys.filter(k => k !== oldFile.path),
        //   })
        return oldFilePathMap.has(v.path) ? oldFilePathMap.get(v.path) : newFile
      })
      currentFileList.setFiles(newFiles)
      onFileListChange(currentFileList.files)
    })
  }, [fileList])

  const store = useLocalStore(() => ({
    loading: false,
    update(data: any) {
      Object.assign(this, data)
    },
    expandedKeys: toJS(expandedKeys) || [],
    freshDir: async path => {
      if (!path) return

      const item = bfsSearch(
        currentFileList.files,
        item => item.path === path
      )?.[0]

      if (!item) {
        return
      }

      store.update({
        loading: true
      })

      const { data } = await Http.get('/file/ls', {
        params: { path: path }
      }).catch(() => {
        store.update({
          loading: false
        })
      })

      if (!item.children || item.children.length <= 0) {
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
                from: item.from,
                isMain: mainKeysOfItem.includes(f.path)
              })
          )
        }
      } else {
        if (Array.isArray(data?.files)) {
          let mainKeysOfItem = []

          if (Array.isArray(item.children)) {
            mainKeysOfItem = item.children.map(v => (v.isMain ? v.path : null))
          }

          let childPaths = item.children.map(c => c.path)

          for (let i = 0; i < data.files.length; i++) {
            let currentFile = data.files[i]

            if (!childPaths.includes(currentFile.path)) {
              item.children.push(
                new FileTreeFile({
                  ...currentFile,
                  isDir: currentFile.is_dir,
                  parent: item,
                  status: 'done',
                  from: item.from,
                  isMain: mainKeysOfItem.includes(currentFile.path)
                })
              )
            }
          }
        }
      }

      store.update({
        loading: false
      })
      onFileListChange(currentFileList.files)
    },
    onExpand: async (keys, records) => {
      if (records.status !== 'done') {
        message.warn('请等待文件上传完成')
        return false
      }

      const oldExpandedKeys = store.expandedKeys.slice()

      const isNeedExpand = !oldExpandedKeys.includes(keys.at(-1))

      store.update({
        expandedKeys: [...new Set(keys)]
      })

      // 记录展开的
      onExpandedKeysChange &&
        onExpandedKeysChange([...new Set(keys)] as string[])

      // 没有要展开的row
      if (keys.length === 0) return false

      if (!isNeedExpand) {
        return false
      }

      let clickedKey = keys.at(-1)
      let item = records

      store.update({
        loading: true
      })
      const { data } = await Http.get('/file/ls', {
        params: { path: clickedKey }
      }).catch(() => {
        store.update({
          loading: false
        })
        return false
      })

      if (!item.children || item.children.length <= 0) {
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
                from: item.from,
                isMain: mainKeysOfItem.includes(f.path)
              })
          )
        }
      } else {
        // 已经有数据，部分更新 children, 找不同
        if (Array.isArray(data?.files)) {
          let mainKeysOfItem = []

          if (Array.isArray(item.children)) {
            mainKeysOfItem = item.children.map(v => (v.isMain ? v.path : null))
          }

          let childPaths = item.children.map(c => c.path)

          for (let i = 0; i < data.files.length; i++) {
            let currentFile = data.files[i]

            if (!childPaths.includes(currentFile.path)) {
              item.children.push(
                new FileTreeFile({
                  ...currentFile,
                  isDir: currentFile.is_dir,
                  parent: item,
                  status: 'done',
                  from: item.from,
                  isMain: mainKeysOfItem.includes(currentFile.path)
                })
              )
            }
          }
        }
      }
      store.update({
        loading: false
      })
      onFileListChange(currentFileList.files)

      return true
    },

    fetch: async () => {},
    // create
    createFolder: async record => {
      store.update({
        loading: true
      })

      try {
        const { data } = await Http.get('/file/ls', {
          params: { path: record.path }
        })

        let childrenFileNames = data?.files.map(f => f.name)
        let newFolderNamePrefix = '新建文件夹'
        let newFolderName = '新建文件夹'
        let i = 1

        while (childrenFileNames.includes(newFolderName)) {
          newFolderName = newFolderNamePrefix + i
          i++
        }

        let newFolder = new FileTreeFile({
          ...record,
          path: `${record.path}/${newFolderName}`,
          name: newFolderName,
          children: [],
          isDir: true,
          parent: record,
          status: 'done',
          from: 'local',
          isMain: false,
          isRenaming: true
        })

        await Http.post(
          '/file/create_dir',
          {
            path: newFolder.path
          },
          { formatErrorMessage: msg => `创建新文件夹失败：${msg}` }
        )
        record?.children?.unshift(newFolder)
      } catch (e) {
        store.update({
          loading: false
        })
      }
      store.update({
        loading: false
      })
    },
    // rename
    rename: async (newName, record) => {
      if (newName === record.name) {
        return
      }

      store.update({
        loading: true
      })

      await Http.put(
        '/file/rename',
        {
          path: record.path,
          new_name: newName
        },
        { formatErrorMessage: msg => `重命名失败：${msg}` }
      ).catch(() =>
        store.update({
          loading: false
        })
      )

      record.name = newName

      let parentPath = record?.parent?.path

      if (!parentPath) {
        let tmpArr = record.path.split('/')
        tmpArr.pop()
        parentPath = tmpArr.join('/')
      }

      record.path = `${parentPath}/${newName}`

      store.update({
        loading: false
      })
    },
    // delete
    deleteNode: async (path, isServerDel = true) => {
      {
        store.update({
          loading: true
        })

        if (isServerDel) {
          await Http.post(
            '/file/delete',
            { paths: [path] },
            { formatErrorMessage: msg => `删除失败：${msg}` }
          ).catch(() =>
            store.update({
              loading: false
            })
          )
        }

        deleteAction(path)
        store.update({
          loading: false
        })

        const item = bfsSearch(
          currentFileList.files,
          item => item.path === path
        )?.[0]
        if (!item) {
          return
        }
        message.success(`已删除 ${item.name}`)

        if (item.parent) {
          item.parent.children = item.parent.children.filter(
            v => v.path !== path
          )
        } else {
          currentFileList.files = currentFileList.files.filter(
            v => v.path !== path
          )
        }

        onFileListChange(currentFileList.files)
      }
    },

    beforeUpload: (params, isLocal) => {
      if (isLocal) {
        return beforeUploadLocalFile
          ? beforeUploadLocalFile(params)
          : Promise.resolve()
      } else {
        return Promise.resolve()
      }
    },

    // upload
    upload: (params, isUploadingDirectory = false, isLocal, targetPath?) => {
      if (isLocal) {
        uploadLocalFile(params, isUploadingDirectory)
        eventEmitter.once('AFTER_UPLOAD_REFRESH', () => {
          const master = params.data ? params.data.master : targetPath
          store.freshDir(master)
        })
      } else {
        uploadServerFile(params, null, targetPath)
        eventEmitter.once('AFTER_SERVER_UPLOAD_REFRESH', () => {
          const master = params.data ? params.data.master : targetPath
          store.freshDir(master)
        })
      }
    },

    setMain: (path, bool) => {
      // psp作业输出路径的原因批量提交可能导致批量作业回传的文件覆盖，因此只设置一个主文件单次提交
      if (!bool) {
        bfsSearch(currentFileList.files, item => {
          item.isMain = false
          return false
        })
        setMainFileKeysAction([])
        return
      }

      let mainItem = null
      bfsSearch(currentFileList.files, item => {
        if (item.path !== path) {
          item.isMain = false
        }

        if (item.path === path) {
          mainItem = item
        }
        return false
      })

      setMainFileKeysAction([mainItem])
    },

    // get mainFileKeys(): string[] {
    //   const mainFiles = bfsSearch(store.fileList, item => item.isMain)
    //   return mainFiles.map(v => v.path)
    // },
    onScroll: e => {
      onScroll && onScroll(e)
    }
  }))

  useEffect(() => {
    if (mainFilePath) {
      let pathes = mainFilePath.split(jobDirPath)

      let subFiles = pathes[1].split('/')

      let mainFileName = subFiles.pop()

      let expandedKeys = []

      let prefix = jobDirPath

      while (subFiles.length) {
        let dirName = subFiles.shift()
        expandedKeys.push(`${prefix}${dirName}`)
        prefix = `${prefix}${dirName}/`
      }

      let promises = []
      let keys = []

      expandedKeys.forEach(key => {
        promises.push(() => {
          const currentFileData = bfsSearch(
            currentFileList.files,
            item => item.path === key
          )?.[0]

          if (currentFileData) {
            keys.push(key)
            return new Promise((solve, reject) =>
              store
                .onExpand(keys, currentFileData)
                .then(() => {
                  solve(true)
                })
                .catch(e => {
                  reject(false)
                })
            )
          } else {
            return Promise.resolve()
          }
        })
      })


      promises.push(() => {
        return new Promise((resolve, reject) => {
          const currentFileData = bfsSearch(
            currentFileList.files,
            item => item.path === mainFilePath
          )?.[0]
          if (currentFileData) {
            currentFileData?.update({
              isMain: true
            })
          }
          store.setMain(mainFilePath, true)
          resolve(true)
        })
      })

      promises.reduce(
        (previousPromise, nextPromise) =>
          previousPromise.then(() => nextPromise()),
        Promise.resolve()
      )
      // store.setMain(store.mainFilePath, true)
    }
  }, [mainFilePath])

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

const state = createStore(useModel)

export const Provider = state.Provider
export const Context = state.Context
export const useStore = state.useStore
