/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useState, useMemo } from 'react'
import { FileAction } from './FileAction'
import { FileActionsStyle } from './style'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'
import { buryPoint, getUrlParams } from '@/utils'
import { useStore as useJobStore } from '../../../store'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'
import { FolderOutlined, FolderFilled } from '@ant-design/icons'

export const FileActions = observer(() => {
  const currentPath = window.localStorage.getItem('CURRENTROUTERPATH')
  const { mode } = useMemo(() => getUrlParams(), [currentPath])
  const [fileLoading, setFileLoading] = useState(false)
  const [folderLoading, setFolderLoading] = useState(false)
  const [myFileLoading, seMyFileLoading] = useState(false)
  const [workdirLoading, setWorkdirLoading] = useState(false)
  const store = useStore()
  const jobStore = useJobStore()
  const fileActions = React.useMemo(() => {
    const temp = [
      {
        title: '创建文件夹',
        icon: 'new_folder',
        visible: jobStore.isTempDirPath,
        loading: false,
        onClick: () => {
          buryPoint({
            category: '作业提交',
            action: '新建文件夹'
          })
          store.createFolder(jobStore.fileTree.id)
        }
      },
      {
        title: '上传文件夹',
        icon: 'folder_upload',
        visible: jobStore.isTempDirPath,
        loading: folderLoading,
        onClick: () => {
          store.upload(jobStore.fileTree.id, true, setFolderLoading)
        }
      },
      {
        title: '上传文件',
        icon: 'file_upload',
        visible: jobStore.isTempDirPath,
        loading: fileLoading,
        onClick: () => {
          store.upload(jobStore.fileTree.id, false, setFileLoading)
        }
      },
      {
        title: '我的文件',
        icon: 'folder_mine',
        visible: true,
        loading: myFileLoading,
        onClick: () => {
          jobStore.tempDirPath = ''
          jobStore.resetFileTree()
          store.copyOrUploadFilesFromFileManager(
            jobStore.fileTree.id,
            jobStore.fileTree.path
          )
        }
      },
      {
        title: '指定工作目录',
        icon: () => (
          <>
            <FolderOutlined rev={'none'} />
            <FolderFilled rev={'none'} style={{ color: 'blue' }} />
          </>
        ),
        visible: localStorage.getItem('psp.workdir.enable'),
        loading: workdirLoading,
        // tips: '指定工作目录后，对目录中文件的操作是不可逆的，请谨慎操作！',
        onClick: () => {
          store.changeWorkdir()
        }
      }
      // {
      //   title: '恢复临时工作目录',
      //   icon: 'folder_mine',
      //   visible: true,
      //   loading: workdirLoading,
      //   onClick: () => {
      //     store.clearWorkdir()
      //   }
      // }
    ]

    return temp
  }, [
    jobStore.fileTree,
    fileLoading,
    folderLoading,
    myFileLoading,
    workdirLoading
  ])

  useEffect(() => {
    const handler = ({ taskKey }) => {
      seMyFileLoading(!!taskKey)
    }
    EE.on(EE_CUSTOM_EVENT.SUPERCOMPUTING_TASKKEY, handler)

    return () => {
      EE.off(EE_CUSTOM_EVENT.SUPERCOMPUTING_TASKKEY, handler)
    }
  }, [])
  return (
    <FileActionsStyle className="fileUploaderToolbar">
      {fileActions
        .filter(item => item.visible)
        .map((item, index) => (
          <FileAction
            icon={item.icon}
            onClick={item.onClick}
            key={index}
            warningTip={item.tips}
            loading={item.loading}>
            {item.title}
          </FileAction>
        ))}
    </FileActionsStyle>
  )
})
