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
    const oldPath = node.path.replace(/\/$/, '')

    // create tempDir
    if (!jobStore.tempDirPath) await jobStore.fetchTempDirPath()

    if (node.name !== '') {
      // 重命名
      const newPath = oldPath.split('/').slice(0, -1).concat(value).join('/')
      if (oldPath !== newPath) {
        await server.move({ srcPaths: oldPath, destPath: newPath })
      }
    } else {
      // 新建
      const newPath = `${jobStore.tempDirPath + '/'}${oldPath}/${value}`
      await server.mkdir(newPath)
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
