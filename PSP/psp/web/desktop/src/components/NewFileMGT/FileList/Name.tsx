/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useMemo, useEffect } from 'react'
import { Tooltip, message } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { EditableText, Icon, Mask } from '@/components'
import { BaseDirectory, BaseFile } from '@/utils/FileSystem'
import { useStore } from '../store'
import { validateFilename } from '@/utils/Validator'
import { useLayoutRect } from '@/utils/hooks'
import { SelectApp } from './SelectApp'
import styled from 'styled-components'
import { env, visualConfig } from '@/domain'
import { previewImage, showTextEditor } from '@/components'
import { newBoxServer } from '@/server'
import { Previewer } from './Previewer'

const StyledLayout = styled.div`
  display: flex;

  > .body {
    display: flex;

    > .icon {
      display: flex;
      align-items: center;
      color: ${({ theme }) => theme.disabledColor};
    }
  }

  > .toolbar {
    visibility: hidden;
    position: absolute;
    background-color: ${({ theme }) => theme.backgroundColorHover};
    right: 0;
    display: flex;
    align-items: center;
    height: 100%;

    > * {
      margin: 0 6px;
    }
  }
`

type Node = BaseDirectory | BaseFile

type Props = {
  nodeId: string
}

export const Name = observer(function Name({ nodeId }: Props) {
  const store = useStore()
  const { server, dir, dirTree, getWidget, isWidgetVisible } = store
  const node = useMemo(
    () => dir.filterFirstNode(item => item.id === nodeId),
    [dir, nodeId]
  )
  const state = useLocalStore(() => ({
    loading: false,
    setLoading(flag) {
      this.loading = flag
    },
    url: '',
    setUrl(url) {
      this.url = url
    }
  }))
  const editModel = EditableText.useModel({
    defaultValue: node.name,
    defaultEditing: false
  })
  const [rect, ref, resize] = useLayoutRect()

  useEffect(() => {
    setTimeout(() => {
      resize()
    }, 0)
  }, [editModel.editing])

  function beforeConfirmRename(node: Node, name) {
    if (!name) {
      return '节点名称不能为空'
    }

    const validation = validateFilename(name)
    if (validation !== true) {
      return validation
    }

    if (node.parent.getDuplicate({ id: node.id, name })) {
      return `${name} 已存在`
    }

    return true
  }

  async function confirmRename(node: Node, name) {
    const parentNode = node.parent
    const path = `${parentNode.path}/${name}`
      .replace(/\/+/, '/')
      .replace(/^\//, '')

    if (node.path === path) {
      return
    }

    // rename
    await server.move({ srcPaths: node.path, destPath: path })
    // sync tree menu
    await server.sync(
      dirTree.filterFirstNode(item => item.path === parentNode.path)
    )
    const data = await server.stat(path)
    node.update({
      path,
      name: data.name,
      mtime: data.mtime,
      type: data.type
    })

    message.success('重命名成功')
  }

  function cancelRename() {
    editModel.setValue(node.name)
  }

  const isImage = type => /gif|jpe?g|tiff|png|webp|bmp$/i.test(type)

  const previewDisable = !isWidgetVisible('preview')

  async function onClick() {
    if (!node.isFile) {
      const menuNode = dirTree.filterFirstNode(item => item.path === node.path)
      store.setNodeId(menuNode?.id)
    } else if (isImage(node.type)) {
      if (previewDisable) {
        return
      }
      // preview image
      try {
        state.setLoading(true)
        const url =
          state.url ||
          (await server.getFileUrl([node.path], [true], [node.size], true))
        state.setUrl(url)
        previewImage({ fileName: node.name, src: url })
      } finally {
        state.setLoading(false)
      }
    } else {
      !previewDisable &&
        showTextEditor({
          path: node.path,
          fileInfo: {
            ...node
          },
          readonly: true,
          boxServerUtil: newBoxServer
        })
    }
  }

  return (
    <StyledLayout className='fileName'>
      {state.loading && <Mask.Spin />}
      <div style={{ width: '100%' }} className='body' onClick={onClick}>
        <div className='icon'>
          {!node.isFile ? (
            <Icon type='folder_close' />
          ) : isImage(node.type) ? (
            <Icon type='image_table' />
          ) : (
            <Icon type='file_table' />
          )}
        </div>
        <div className='editor' style={{ width: 'calc(100% - 24px)' }}>
          <EditableText
            style={{ width: '100%' }}
            model={editModel}
            showEdit={false}
            beforeConfirm={name => beforeConfirmRename(node, name)}
            onConfirm={name => confirmRename(node, name)}
            onCancel={cancelRename}
          />
        </div>
      </div>
      <div ref={ref} className='toolbar'>
        {!editModel.editing && (
          <>
            {isImage(node.type) && !previewDisable && (
              <Previewer nodeId={nodeId} />
            )}
            {getWidget('cloud-app') ||
              (!!visualConfig.showVisualizeApp && node.isFile && (
                <SelectApp node={node} />
              ))}
            {getWidget('rename') || (
              <Tooltip title='重命名'>
                <Icon
                  type='rename'
                  onClick={() => editModel.setEditing(true)}
                />
              </Tooltip>
            )}
          </>
        )}
      </div>
    </StyledLayout>
  )
})
