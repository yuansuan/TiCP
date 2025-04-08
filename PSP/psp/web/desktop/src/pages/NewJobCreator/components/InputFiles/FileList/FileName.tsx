/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Tooltip, message } from 'antd'
import { observer } from 'mobx-react-lite'
import { EditableText, Icon } from '@/components'
import { HoverableIcon } from '@/components'
import { ChildNode } from '@/domain/JobBuilder/JobDirectory'
import { EditableName } from './EditableName'
import { FileNameStyle } from './style'
import { buryPoint } from '@/utils'
interface IProps {
  record: ChildNode
  createFolder: any
  upload: any
  isExpand: boolean
  copyFilesFromFileManager: any
}

export const FileName = observer((props: IProps) => {
  const { record, createFolder, upload, isExpand, copyFilesFromFileManager } =
    props

  const model = EditableText.useModel({
    defaultEditing: !record.isRoot && record.name === '',
    defaultValue: record.name
  })

  const iconType = record.isFile
    ? 'file_table'
    : isExpand
    ? 'folder_open'
    : 'folder_close'

  const rename = () => {
    model.setEditing(true)
  }

  const fileActions = [
    <Tooltip title='重命名' key='rename'>
      <span>
        <HoverableIcon type='edit' onClick={rename} />
      </span>
    </Tooltip>
  ]

  const dirActions = [
    <Tooltip title='新建文件夹' key='newDir'>
      <span>
        <HoverableIcon
          type='new_folder'
          onClick={() => {
            buryPoint({
              category: '作业提交',
              action: '新建文件夹'
            })
            createFolder(record.id)
          }}
        />
      </span>
    </Tooltip>,
    <Tooltip title='上传文件夹' key='uploadDir'>
      <span>
        <HoverableIcon
          type='folder_upload'
          onClick={() => {
            buryPoint({
              category: '作业提交',
              action: '上传文件夹'
            })
            upload(record.id, true)
          }}
        />
      </span>
    </Tooltip>,
    <Tooltip title='上传文件' key='uploadFile'>
      <span>
        <HoverableIcon
          type='file_upload'
          onClick={() => {
            buryPoint({
              category: '作业提交',
              action: '上传文件'
            })
            upload(record.id)
          }}
        />
      </span>
    </Tooltip>,
    <Tooltip title='我的文件' key='myFile'>
      <span>
        <HoverableIcon
          type='folder_mine'
          onClick={() => {
            buryPoint({
              category: '作业提交',
              action: '我的文件'
            })
            copyFilesFromFileManager(record.id, record.path)
          }}
        />
      </span>
    </Tooltip>,
    <Tooltip title='重命名' key='rename'>
      <span>
        <HoverableIcon type='edit' onClick={rename} />
      </span>
    </Tooltip>
  ]

  const rootActions = dirActions.slice(0, -1)

  const getActions = () => {
    if (model.editing) return []
    if (record.isFile) {
      return (record as any).status === 'done' ? fileActions : []
    }
    if (record.isRoot) {
      return rootActions
    }
    return dirActions
  }

  return (
    <FileNameStyle>
      <Icon className='icon-prefix' type={iconType} />
      <div className='filename-text'>
        {record.isRoot ? (
          'HomeFile'
        ) : (
          <EditableName node={record} model={model} />
        )}
        <div className={record.isRoot ? 'actions root' : 'actions'}>
          {getActions()}
        </div>
      </div>
    </FileNameStyle>
  )
})
