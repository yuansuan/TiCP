/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useMemo } from 'react'
import { FileAction } from './FileAction'
import { FileActionsStyle } from './style'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'
import { downloadTestFile, buryPoint } from '@/utils'
import { env } from '@/domain'
import { getUrlParams } from '@/utils/Validator'
import { useStore as useJobStore } from '../../../store'

export const FileActions = observer(() => {
  const currentPath = window.localStorage.getItem('CURRENTROUTERPATH')
  const { mode } = useMemo(() => getUrlParams(), [currentPath])
  const store = useStore()
  const jobStore = useJobStore()

  const fileActions = React.useMemo(() => {
    const temp = [
      {
        title: '创建文件夹',
        icon: 'new_folder',
        visible: true,

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
        visible: true,

        onClick: () => {
          buryPoint({
            category: '作业提交',
            action: '上传文件夹'
          })
          store.upload(jobStore.fileTree.id, true)
        }
      },
      {
        title: '上传文件',
        icon: 'file_upload',
        visible: true,
        onClick: () => {
          buryPoint({
            category: '作业提交',
            action: '上传文件'
          })
          store.upload(jobStore.fileTree.id)
        }
      },
      {
        title: '我的文件',
        icon: 'folder_mine',
        visible: jobStore.is_cloud ? false : true,
        onClick: () => {
          buryPoint({
            category: '作业提交',
            action: '我的文件'
          })
          store.copyFilesFromFileManager(
            jobStore.fileTree.id,
            jobStore.fileTree.path
          )
        }
      }
    ]

    return temp
  }, [jobStore.fileTree])

  return (
    <FileActionsStyle className='fileUploaderToolbar'>
      {fileActions
        .filter(item => item.visible === false)
        .map((item, index) => (
          <FileAction icon={item.icon} onClick={item.onClick} key={index}>
            {item.title}
          </FileAction>
        ))}
    </FileActionsStyle>
  )
})
