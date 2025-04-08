/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import { EditableText } from '@/components'
import { ChildNode } from '@/domain/JobBuilder/JobDirectory'
import { JobFile } from '@/domain/JobBuilder/JobFile'
import { validateFilename } from '@/utils/Validator'
import { EditableNameStyle } from './style'
import { serverFactory } from '@/components/NewFileMGT/store/common'
import { newBoxServer } from '@/server'
import { useStore as JobStore } from '../../../store'

interface IProps {
  node: ChildNode
  model: ReturnType<typeof EditableText.useModel>
}

export const EditableName = observer((props: IProps) => {
  const { node, model } = props
  const server = serverFactory(newBoxServer)
  const jobStore = JobStore()
  const onCancel = () => {
    if (node.name === '') {
      node.parent.removeFirstNode(item => item.id === node.id)
    } else {
      model.setError('')
      model.setValue(node.name)
    }
  }

  const beforeConfirm = async (value: string) => {
    const validation = validateFilename(value)
    if (validation !== true) {
      return Promise.reject(validation)
    }

    const duplicateNode = node.parent.getDuplicate({ ...node, name: value })
    if (duplicateNode) {
      return Promise.reject(`${value}已存在`)
    }
    // create tempDir
    if (!jobStore.tempDirPath) await jobStore.fetchTempDirPath()
    const newTempDirPath =
      jobStore.tempDirPath && jobStore.tempDirPath.endsWith('/')
        ? jobStore.tempDirPath
        : jobStore.tempDirPath + '/'

    const tmpPath = jobStore.isTempDirPath
      ? node.path
      : node.path.replace(new RegExp(jobStore.tempDirPath), '')

    const oldPath = (newTempDirPath + tmpPath.replace(/\/$/, '')).replace(
      /\/\//g,
      '/'
    )

    if (node.name !== '') {
      // 重命名
      const newPath = oldPath.split('/').slice(0, -1).concat(value).join('/')
      if (oldPath !== newPath) {
        await server.rename({
          path: oldPath,
          newName: newPath,
          cross: jobStore.isTempDirPath,
          is_cloud: jobStore.isCloud
        })
      }
    } else {
      // 新建
      const newPath = `${oldPath}/${value}`.replace(/\/\//g, '/')

      await server.mkdir(newPath, jobStore.isTempDirPath, jobStore.isCloud)
    }

    return Promise.resolve()
  }

  const onConfirm = value => {
    node.update({ name: value })
  }

  return (
    <EditableNameStyle>
      {node.isFile && (node as JobFile).status !== 'done' ? (
        node.name
      ) : (
        <EditableText
          style={{ fontSize: 14 }}
          onCancel={onCancel}
          beforeConfirm={beforeConfirm}
          onConfirm={onConfirm}
          showEdit={false}
          model={model}
        />
      )}
    </EditableNameStyle>
  )
})
