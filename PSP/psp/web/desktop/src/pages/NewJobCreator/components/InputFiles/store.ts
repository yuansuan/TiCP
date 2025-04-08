/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { env, uploader, NewBoxHttp } from '@/domain'
import { Fetch, Http } from '@/utils'
import { currentUser } from '@/domain'
import { JobDirectory } from '@/domain/JobBuilder/JobDirectory'
import { Modal } from '@/components'
import { message } from 'antd'
import { showFailure, showFileSelector } from '@/components/NewFileMGT'
import {
  showFailure as newShowFailure,
  showFileSelector as newShowFileSelector
} from '@/components/NewFileMGT'
import { v4 as uuid } from 'uuid'
import { useStore as useJobStore } from '../../store'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'

const UPLOAD_ID = uuid()

function useModel() {
  const createdMap = React.useRef({})

  const jobStore = useJobStore()

  const store = useLocalStore(() => ({
    expandedKeysSet: new Set(
      jobStore.fileTree.filterNodes(node => !node.isFile).map(item => item.id)
    ),
    get expandedKeys() {
      return [...this.expandedKeysSet]
    },
    onExpand: keys => {
      store.expandedKeysSet = new Set(keys)
    },

    // create
    createFolder: id => {
      store.expandedKeysSet.add(id)
      const parentNode = jobStore.fileTree.filterFirstNode(
        item => item.id === id
      )
      parentNode.unshift(new JobDirectory({ name: '' }))
    },

    // delete
    deleteNode: async id => {
      const node = jobStore.fileTree.filterFirstNode(item => item.id === id)
      if (node && node.name) {
        await Modal.showConfirm({
          title: '删除文件',
          content: `确认要删除当前${node.isFile ? '文件' : '目录'}吗？`
        })
        const { uid, path } = node
        await jobStore.draft.deleteFiles([path])
        if (uid) {
          uploader.remove(uid)
        }
        message.success('文件删除成功')
      }
      jobStore.fileTree.removeFirstNode(item => item.id === id)
    },

    // upload
    upload: async (id: string, isUploadingDirectory = false) => {
      const currentPath = jobStore.fileTree.filterFirstNode(
        node => node.id === id
      ).path
      const uniqueID = uuid()

      if (!jobStore.tempDirPath) await jobStore.fetchTempDirPath()

      uploader.upload({
        origin: UPLOAD_ID,
        by: 'chunk',
        multiple: true,
        action: '/storage/upload',
        httpAdapter: Fetch,
        data: {
          directory: isUploadingDirectory,
          _uid: uniqueID,
          dir: currentPath,
          tempDirPath: jobStore.tempDirPath,
          user_name: currentUser.name,
          cross: jobStore.isCross,
          is_cloud:jobStore.isCloud,
        },
        directory: isUploadingDirectory,
        // 遇到同名文件（夹），让用户挑选覆盖哪些原有文件（夹）
        beforeUpload: async uploadingFiles => {
          /*
           * @return: 上传后的目标位置
           * 示例：
           *   - 在 parent 目录上传 childFile.png 文件，返回 'parent/childFile.png'
           *   - 在 parent 目录上传 childDirectory 文件夹，返回'parent/childDirectory'
           */
          const getTargetPath = file => {
            const relativePathOfFile = file.webkitRelativePath || file.name
            const filenameOrDirectoryName = relativePathOfFile.split('/')[0]
            return `${currentPath}/${filenameOrDirectoryName}`.replace(
              /^\//,
              ''
            )
          }

          /*
           *  @return: 正在上传的目录名
           */
          const getDirectoryName = files => {
            const file = files[0]
            const relativeFilePath = file.webkitRelativePath || file.name
            return relativeFilePath.split('/')[0]
          }

          /*
           * 上传文件时的重名文件
           * 上传文件夹时，重名文件夹的所有文件
           */
          const repetitiveFiles = uploadingFiles.filter(file =>
            jobStore.fileTree.filterFirstNode(
              node => node.path === getTargetPath(file)
            )
          )
          const nonRepetitiveFiles = uploadingFiles.filter(
            file =>
              !jobStore.fileTree.filterFirstNode(
                node => node.path === getTargetPath(file)
              )
          )

          const filesToBeUploaded = [...nonRepetitiveFiles]

          if (repetitiveFiles.length > 0) {
            if (isUploadingDirectory === true) {
              const directoryName = getDirectoryName(repetitiveFiles)
              // 选中的文件夹将被覆盖
              const selectedDirectory = await showFailure({
                actionName: '上传',
                items: [
                  {
                    isFile: false,
                    name: directoryName
                  }
                ]
              })
              if (selectedDirectory.length > 0) {
                await jobStore.draft.deleteFile(
                  `${currentPath}/${directoryName}`
                )
                filesToBeUploaded.push(...repetitiveFiles)
              }
            } else if (isUploadingDirectory === false) {
              // 选中的文件将被覆盖
              const selectedFiles = await showFailure({
                actionName: '上传',
                items: repetitiveFiles.map(file => ({
                  isFile: true,
                  name: file.name,
                  uid: file.uid
                }))
              })
              if (selectedFiles.length > 0) {
                await jobStore.draft.deleteFiles(
                  selectedFiles.map(files => `${currentPath}/${files.name}`)
                )
                filesToBeUploaded.push(...selectedFiles)
              }
            }
          }

          if (filesToBeUploaded.length > 0) {
            store.expandedKeysSet.add(id)
            // 触发显示 dropdown
            EE.emit(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, { visible: true })
          }

          return filesToBeUploaded.map(item => item.uid)
        },
        onChange: ({ file, origin }) => {
          if (origin !== UPLOAD_ID) {
            return
          }

          if (file.status === 'done') {
            // 有文件上传完成，发消息，check 是否要关闭 dropdown
            EE.emit(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, { visible: false })
          }

          if (!createdMap.current[file.uid]) {
            createdMap.current[file.uid] = true
            jobStore.fileTree.uploadLocalFile(id, file)
          }
          const jobFile = jobStore.fileTree.filterFirstNode(
            item => item.uid === file.uid
          )
          if (!jobFile) return
          jobFile.update({
            status: file.status,
            percent: file.percent
          })
        },
        onRemove: async file => {
          const f = jobStore.fileTree.filterFirstNode(
            item => item.uid === file.uid
          )
          if (f) {
            await jobStore.draft.deleteFile(f.path)
            jobStore.fileTree.removeFirstNode(item => item.path === f.path)
          }
        }
      })
    },

    // import from file manager
    copyFilesFromFileManager: async (
      id: string,
      destinationDirectory: string
    ) => {
      const currentPath = jobStore.fileTree.filterFirstNode(
        node => node.id === id
      ).path
      const selectedFilePaths = await showFileSelector(false)
      let directoryOfFileManager = ''
      const duplicateNodes = []
      const mapFromOriginalToDestination = selectedFilePaths.reduce(
        (pathsObj, selectedFilePath) => {
          if (!directoryOfFileManager) {
            const pathFractions = selectedFilePath.split('/')
            pathFractions.pop()
            directoryOfFileManager = pathFractions.join('/').replace(/^\//, '')
          }
          const filename = selectedFilePath.split('/').pop()
          const targetPath = `${destinationDirectory}/${filename}`.replace(
            /^\//,
            ''
          )

          const duplicateNode = jobStore.fileTree.filterFirstNode(
            node => node.path === targetPath
          )
          if (duplicateNode) {
            duplicateNodes.push(duplicateNode)
          } else {
            pathsObj[selectedFilePath] = targetPath
          }
          return pathsObj
        },
        {}
      )

      console.log(
        'mapFromOriginalToDestination: ',
        mapFromOriginalToDestination
      )
      if (duplicateNodes.length > 0) {
        const nodesWillBeReplaced = await showFailure({
          actionName: '移动',
          items: duplicateNodes
        })
        if (nodesWillBeReplaced.length > 0) {
          // delete first
          await jobStore.draft.deleteFiles(
            nodesWillBeReplaced.map(item => item.path),
            jobStore.isCross,
            jobStore.isCloud,
          )

          // add to mapFromOriginalToDestination
          nodesWillBeReplaced.reduce((pathsObj, nodeWillBeReplaced) => {
            pathsObj[`${directoryOfFileManager}/${nodeWillBeReplaced.name}`] =
              nodeWillBeReplaced.path
            return pathsObj
          }, mapFromOriginalToDestination)
        }
      }

      if (Object.keys(mapFromOriginalToDestination).length > 0) {
        if (!jobStore.tempDirPath) await jobStore.fetchTempDirPath()

        const lsFilePath = [jobStore.tempDirPath, ...currentPath.split('/')]
          .filter(item => !!item)
          .join('/')
        const data = await jobStore.draft.copyFromCommon(
          mapFromOriginalToDestination,
          lsFilePath
        )
        store.expandedKeysSet.add(id)
        jobStore.fileTree.uploadCommonFiles(data,jobStore.tempDirPath)
      }
    },

    // drag and drop file
    onDrop: async (dragKey, dropKey) => {
      const dragNode = jobStore.fileTree.filterFirstNode(
        item => item.id === dragKey
      )
      const dropNode = jobStore.fileTree.filterFirstNode(
        item => item.id === dropKey
      )
      store.expandedKeysSet.add(dropKey)
      if (dragNode.isFile && dragNode.status !== 'done') return
      if (!dropNode.isFile && dragNode.parent !== dropNode) {
        const duplicateNode = dropNode.getDuplicate(dragNode)
        if (duplicateNode) {
          // you can't move sub file/directory to override parent file/directory
          if (dragNode.path.startsWith(duplicateNode?.path)) {
            message.error('无法覆盖同名父文件，请重命名文件再移动')
            return
          }
          const coverNodes = await showFailure({
            actionName: '移动',
            items: [duplicateNode]
          })
          if (coverNodes.length > 0) {
            await jobStore.draft.deleteFile(coverNodes[0].path)
          } else {
            return
          }
        }
        await jobStore.draft.moveFile(dragNode.path, dropNode.path)
        jobStore.fileTree.removeFirstNode(item => item.id === dragKey)
        dropNode.unshift(dragNode)
      }
    }
  }))

  React.useEffect(() => {
    store.expandedKeysSet = new Set(
      jobStore.fileTree.filterNodes(node => !node.isFile).map(item => item.id)
    )
  }, [jobStore.fileTree])

  return store
}

const state = createStore(useModel)

export const Provider = state.Provider
export const Context = state.Context
export const useStore = state.useStore
