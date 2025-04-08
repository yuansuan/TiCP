/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { action, computed, observable } from 'mobx'
import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { uploader } from '@/domain'
import { Http } from '@/utils'
import { currentUser } from '@/domain'
import { JobDirectory } from '@/domain/JobBuilder/JobDirectory'
import { Modal } from '@/components'
import { message } from 'antd'
import {
  showFailure,
  showFileSelector,
  showDirSelector
} from '@/components/NewFileMGT'
import { v4 as uuid } from 'uuid'
import { useStore as useJobStore } from '../../store'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'
import { escapeRegExp } from '@/utils'

const UPLOAD_ID = uuid()

export type Props = {
  fileList: any[]
  deleteAction: (path: string) => void
  setMainFileKeysAction: (keys: string[]) => void
}

type Store = {
  loading: boolean
}

export function useModel() {
  const createdMap = React.useRef({})
  const jobStore = useJobStore()

  const store = useLocalStore(() => ({
    loading: false,
    update(data: Partial<Store>) {
      Object.assign(store, data)
    },
    expandedKeysSet: new Set(
      jobStore.fileTree.filterNodes(node => !node.isFile).map(item => item.id)
    ),
    get expandedKeys() {
      return [...this.expandedKeysSet]
    },
    onExpand: async (keys, data) => {
      if (
        jobStore.tempDirPath !== '' &&
        !data?.isFile &&
        new Set(keys).has(data.id)
      ) {
        const fileData = await jobStore.draft.getFileList({
          path: !!data?.realCommonPathPrefix
            ? jobStore.tempDirPath + '/' + data?.realCommonPathPrefix
            : jobStore.tempDirPath + '/',
          cross: jobStore.isTempDirPath,
          is_cloud: jobStore.isCloud
        })
        jobStore.fileTree.uploadCommonFiles(
          fileData,
          jobStore.tempDirPath,
          new Set(jobStore.getResubmitParam().main_files),
          jobStore.isTempDirPath
        )
      }
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
        const tmpPath = jobStore.isTempDirPath
          ? path
          : path.replace(new RegExp(escapeRegExp(jobStore.tempDirPath)), '')
        const newPath = (jobStore.tempDirPath + '/' + path).replace(
          /\/\//g,
          '/'
        )

        if (node.status === 'done' || node.status === undefined) {
          await jobStore.draft.deleteFile(
            newPath,
            jobStore.isTempDirPath,
            jobStore.isCloud
          )
        }

        if (uid) {
          uploader.remove(uid)
        }
        message.success('文件删除成功')
      }
      jobStore.fileTree.removeFirstNode(item => item.id === id)
    },

    // upload
    upload: async (id: string, isUploadingDirectory = false, setBtnLoading) => {
      const currentPath = jobStore.fileTree.filterFirstNode(
        node => node.id === id
      ).path

      const uniqueID = uuid()
      uploader.upload({
        origin: UPLOAD_ID,
        by: 'chunk',
        multiple: true,
        action: '/storage/upload',
        httpAdapter: Http,
        data: {
          directory: isUploadingDirectory,
          // _uid: uniqueID,
          dir: jobStore.isTempDirPath
            ? currentPath
            : currentPath.replace(
              new RegExp(escapeRegExp(jobStore.tempDirPath)),
              ''
            ),
          user_name: currentUser.name,
          is_cloud: jobStore.isCloud
        },
        directory: isUploadingDirectory,
        beforeUpload: async uploadingFiles => {
          // 上传前检查是否创建临时目录
          if (!jobStore.tempDirPath) {
            try {
              setBtnLoading(true)
              await jobStore.fetchTempDirPath()
            } finally {
              setBtnLoading(false)
            }
          }
          // 追加属性
          uploader.props.data['tempDirPath'] = jobStore.tempDirPath
          uploader.props.data['cross'] = jobStore.isTempDirPath
          // 遇到同名文件（夹），让用户挑选覆盖哪些原有文件（夹）
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
                const deletePath = `${jobStore.tempDirPath}/${
                  jobStore.isTempDirPath
                    ? currentPath
                    : currentPath.replace(
                      new RegExp(escapeRegExp(jobStore.tempDirPath)),
                      ''
                    )
                }/${directoryName}`
                await jobStore.draft.deleteFile(
                  deletePath,
                  jobStore.isTempDirPath,
                  jobStore.isCloud
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
                const deletePath = `${jobStore.tempDirPath}/${
                  jobStore.isTempDirPath
                    ? currentPath
                    : currentPath.replace(
                      new RegExp(escapeRegExp(jobStore.tempDirPath)),
                      ''
                    )
                }`

                await jobStore.draft.deleteFiles(
                  selectedFiles.map(files => `${deletePath}/${files.name}`),
                  jobStore.isTempDirPath,
                  jobStore.isCloud
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
            EE.emit(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, {
              visible: false
            })
            // if (!jobStore.isTempDirPath) jobStore.fetchJobTree()
          }

          if (!createdMap.current[file.uid]) {
            createdMap.current[file.uid] = true
            jobStore.fileTree.uploadLocalFile(
              id,
              file,
              jobStore.isTempDirPath,
              jobStore.tempDirPath
            )
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
          // if (f) {
          //   await jobStore.draft.deleteFile(f.path)
          //   jobStore.fileTree.removeFirstNode(item => item.path === f.path)
          // }
        }
      })
    },

    // import from file manager
    copyOrUploadFilesFromFileManager: async (
      id: string,
      destinationDirectory: string
    ) => {
      const currentPath = jobStore.fileTree.filterFirstNode(
        node => node.id === id
      ).path
      const selectedFiles = await showFileSelector(true, jobStore.isTempDirPath)
      const selectedFilePaths = selectedFiles.map(file => file.path)
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

      if (duplicateNodes.length > 0) {
        const nodesWillBeReplaced = await showFailure({
          actionName: '移动',
          items: duplicateNodes
        })
        if (nodesWillBeReplaced.length > 0) {
          // delete first
          await jobStore.draft.deleteFiles(
            nodesWillBeReplaced.map(item => {
              const tmpPath = jobStore.isTempDirPath
                ? item.path
                : item.path.replace(
                  new RegExp(escapeRegExp(jobStore.tempDirPath)),
                  ''
                )
              return jobStore.tempDirPath + '/' + tmpPath
            }),
            jobStore.isTempDirPath,
            jobStore.isCloud
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
        console.log('current tmp path: ', jobStore.isTempDirPath, jobStore.tempDirPath)
        if (jobStore.isTempDirPath && !jobStore.tempDirPath) {
          await jobStore.fetchTempDirPath()
        } else {
          let selectWorkDirPath =
            selectedFiles
              ?.find(file => !file.isFile)
              ?.path.replace(/^\.\/+/, '')
              ?.replace(/\/+$/, '') || ''
          if (selectWorkDirPath === '') {
            return
          }
          jobStore.setTempDirPath(selectWorkDirPath)
        }
        const newDestPath =
          jobStore.tempDirPath + '/' + currentPath.replace(/\/+$/, '')

        try {
          store.update({ loading: true })
          const data = await jobStore.draft.copyOrUploadFromCommon(
            mapFromOriginalToDestination,
            newDestPath,
            jobStore.isTempDirPath,
            jobStore.isCloud,
            selectedFiles
          )

          EE.once(
            EE_CUSTOM_EVENT.SERVER_FILE_TO_SUPERCOMPUTING,
            async ({ file_status }) => {
              console.log('file_status=========+>: ', file_status)
              if (file_status === 'success') {
                const fileData = await jobStore.draft.getFileList({
                  path: newDestPath,
                  cross: jobStore.isTempDirPath,
                  is_cloud: jobStore.isCloud
                })
                store.expandedKeysSet.add(id)
                let selectFileSet = new Set<String>()
                if (jobStore.jobBuildMode === 'resubmit') {
                  selectFileSet = new Set<String>(
                    jobStore.getResubmitParam()?.main_files
                  )
                }
                jobStore.fileTree.uploadCommonFiles(
                  fileData,
                  jobStore.tempDirPath,
                  selectFileSet,
                  jobStore.isTempDirPath,
                  id
                )
              }
            }
          )
          if (!jobStore.isCloud) {
            store.update({ loading: false })
            store.expandedKeysSet.add(id)
            let selectFileSet = new Set<String>()
            if (jobStore.jobBuildMode === 'resubmit') {
              selectFileSet = new Set<String>(
                jobStore.getResubmitParam()?.main_files
              )
            }
            jobStore.fileTree.uploadCommonFiles(
              data,
              jobStore.tempDirPath,
              selectFileSet,
              jobStore.isTempDirPath,
              id
            )
          }
        } finally {
          store.update({ loading: false })
        }
      }
    },

    changeWorkdir: async () => {
      const dir = await showDirSelector({
        disabledCheckedPaths: ['.'],
        title: '指定工作目录'
      })
      jobStore.setTempDirPath(dir)
      jobStore.setIsTempDirPath(false)
      // 切换成功后，刷新 job tree
      store.update({ loading: true })
      await jobStore.fetchJobTree()
      store.update({ loading: false })
      message.info('切换工作目录成功')
    },

    clearWorkdir: async () => {
      await jobStore.fetchTempDirPath()
      jobStore.setIsTempDirPath(true)
      // jobStore.setTempDirPath(dir)
      // 切换成功后，刷新 job tree
      await jobStore.fetchJobTree()
      message.info('恢复到临时工作目录成功')
    },

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

        // const addTempDirPath =
        await jobStore.draft.moveFile(
          dragNode.path,
          dropNode.path,
          jobStore.isCross,
          jobStore.isCloud
        )
        jobStore.fileTree.removeFirstNode(item => item.id === dragKey)
        dropNode.unshift(dragNode)
      }
    },

    mainFilesExpand: () => {
      if (jobStore.jobBuildMode === 'resubmit') {
        ;(jobStore.getResubmitParam().main_files || []).forEach(filePath => {
          if (filePath.indexOf('/') === -1) {
            return
          }
          const paths = ['', ...filePath.replace(/\/[^\/]*$/, '').split('/')]

          let currentPath = ''
          ;(async currentPath => {
            console.time('mainFilesExpand Time')
            store.loading = true
            const promises = []
            let currentPathSet = new Set()
            for (let i = 0; i < paths.length; i++) {
              currentPath =
                currentPath === '' ? paths[i] : currentPath + '/' + paths[i]
              promises.push(
                jobStore.draft.getFileList({
                  path: jobStore.tempDirPath + '/' + currentPath,
                  cross: jobStore.isTempDirPath,
                  is_cloud: jobStore.isCloud
                })
              )
              currentPathSet.add(
                i === paths.length - 1
                  ? currentPath === ''
                    ? paths[i]
                    : currentPath + '/' + paths[i]
                  : currentPath === ''
                    ? paths[i + 1]
                    : currentPath + '/' + paths[i + 1]
              )
            }

            const fileDataArray = await Promise.all(promises)
            const expandKeys = jobStore.fileTree.uploadCommonFiles(
              fileDataArray.flat(),
              jobStore.tempDirPath,
              new Set(jobStore.getResubmitParam().main_files),
              jobStore.isTempDirPath,
              null,
              currentPathSet
            )
            store.expandedKeysSet = new Set([
              ...store.expandedKeysSet,
              ...expandKeys
            ])
            store.loading = false
            console.timeEnd('mainFilesExpand Time')
          })(currentPath)
          currentPath = ''
        })
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
