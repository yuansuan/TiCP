/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { Tooltip, message } from 'antd'
import { observer } from 'mobx-react-lite'
import { EditableText, Icon, Modal } from '@/components'
import { BaseDirectory } from '@/utils/FileSystem'
import { useStore } from '../store'
import { validateFilename } from '@/utils/Validator'
import { useLayoutRect } from '@/utils/hooks'
import styled from 'styled-components'
import { Http } from '@/utils'
import { env } from '@/domain'

const StyledLayout = styled.div`
  display: inline-block;
  width: calc(100% - 24px);

  &:hover {
    > .menu_toolbar {
      visibility: visible;
    }
  }

  > .main {
    display: inline-block;
    vertical-align: middle;
  }

  > .menu_toolbar {
    visibility: hidden;
    display: inline-block;
    vertical-align: middle;
    color: ${({ theme }) => theme.primaryColor};

    > .delete {
      color: ${({ theme }) => theme.errorColor};
    }
  }
`

type Props = {
  node: BaseDirectory
}

export const Title = observer(function Title({ node }: Props) {
  const store = useStore()
  const {
    nodeId,
    server,
    dir,
    setNodeId,
    selectedKeys,
    setSelectedKeys,
    getWidget
  } = store

  const editModel = EditableText.useModel({
    defaultEditing: !node.path,
    defaultValue: node.path ? node.name : '未命名文件夹'
  })
  const [toolbarRect, toolbarRef, toolbarResize] = useLayoutRect()
  useEffect(() => {
    toolbarResize()
  }, [editModel.editing])

  function beforeConfirmRename(node: BaseDirectory, name) {
    if (!name) {
      return '文件名不能为空'
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

  async function confirmRename(node, name) {
    const parentNode = node.parent
    const path = `${parentNode.path}/${name}`
      .replace(/\/+/, '/')
      .replace(/^\//, '')

    // create node
    if (!node.path) {
      await server.mkdir(path)
      await server.sync(parentNode)
      message.success('创建成功')
    } else if (node.path !== path) {
      // rename
      await server.move([[node.path, path]])
      await server.sync(parentNode)

      // {

      //     Http.post('/filerecord/record', {
      //       type: 4,
      //       info: {
      //         storage_size: node.size,
      //         file_name: name,
      //         file_type: node.isFile ? 1 : 2
      //       }
      //     })
      // }

      message.success('重命名成功')
    } else {
      return
    }

    const newNode = parentNode.filterFirstNode(item => item.path === path)
    setNodeId(newNode.id)
  }

  function cancelRename(node) {
    if (!node.path) {
      node.parent.removeFirstNode(item => item.id === node.id)
    } else {
      editModel.setValue(node.name)
    }
  }

  async function deleteNode(node) {
    await Modal.showConfirm({
      title: '删除目录',
      content: `确认要删除该${node.name}目录吗？`
    })

    const parentNode = node.parent
    await server.delete([node.path])

    parentNode.removeFirstNode(item => item.id === node.id)
    message.success('目录删除成功')

    // delete selected node
    if (node.id === nodeId) {
      store.setNodeId(parentNode.id)
    }
    // delete child node
    if (parentNode.id === nodeId) {
      const removedNode = dir.removeFirstNode(item => item.path === node.path)
      setSelectedKeys(selectedKeys.filter(key => key !== removedNode.id))
    }
  }

  return (
    <StyledLayout>
      <div
        className='main'
        style={{
          width: `calc(100% - ${toolbarRect.width}px)`
        }}>
        <EditableText
          showEdit={false}
          model={editModel}
          beforeConfirm={name => beforeConfirmRename(node, name)}
          onConfirm={name => confirmRename(node, name)}
          onCancel={() => cancelRename(node)}
        />
      </div>
      <div
        className='menu_toolbar'
        ref={toolbarRef}
        onClick={e => e.stopPropagation()}>
        {!editModel.editing && (
          <>
            {getWidget('rename') || (
              <Tooltip title='重命名'>
                <Icon
                  type='rename'
                  onClick={() => editModel.setEditing(true)}
                />
              </Tooltip>
            )}
            {getWidget('delete') || (
              <Tooltip title='删除'>
                <Icon
                  className='delete'
                  type='cancel'
                  onClick={e => deleteNode(node)}
                />
              </Tooltip>
            )}
          </>
        )}
      </div>
    </StyledLayout>
  )
})
