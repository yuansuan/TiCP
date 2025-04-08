/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { Tree, Tooltip, message } from 'antd'
import { useStore, Context, useModel } from './store'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Icon, Modal, Button, EditableText } from '@/components'
import { BaseDirectory } from '@/utils/FileSystem'
import { validateFilename } from '@/utils/Validator'
import { getUrlParams } from '@/utils'
const { DirectoryTree } = Tree
const StyledLayout = styled.div`
  .ant-tree.ant-tree-directory {
    .ant-tree-treenode {
      &.ant-tree-treenode-disabled {
        .ant-tree-node-content-wrapper {
          .ant-tree-iconEle svg {
            fill: ${({ theme }) => theme.disabledColor};
          }
        }
      }

      .ant-tree-iconEle {
        vertical-align: middle;
      }

      .ant-tree-node-content-wrapper {
        .ant-tree-iconEle svg {
          fill: ${({ theme }) => theme.primaryColor};
        }

        &.ant-tree-node-selected {
          .ant-tree-iconEle svg {
            fill: white;
          }
        }
      }
    }
  }

  > .main {
    width: 100%;
    overflow: auto;
    padding-bottom: 40px;

    .nodeName {
      display: inline-block;
      width: calc(100% - 24px);

      > .main {
        display: inline-block;
        vertical-align: middle;
      }

      > .toolbar {
        display: inline-block;
        vertical-align: middle;

        > .creator {
          margin-left: 8px;
          transform: translateY(2px);

          &:hover {
            color: ${props => props.theme.primaryColor};
          }
        }
      }
    }
  }

  > .footer {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    padding: 10px 0;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};
  }
`

type Props = {
  disabledPaths?: string[]
  onCancel: () => void
  onOk: (path: string) => void | Promise<void>
}

export const DirSelector = observer(function DirSelector({
  disabledPaths = [],
  onCancel,
  onOk
}: Props) {
  const store = useStore()
  const { dirTree, server } = store
  const isCloud = getUrlParams()?.isCloud && JSON.parse(getUrlParams()?.isCloud)
  const state = useLocalStore(() => ({
    expandedKeys: [],
    setExpandedKeys(keys) {
      this.expandedKeys = [...keys]
    },
    selectedKeys: [],
    setSelectedKeys(keys) {
      this.selectedKeys = [...keys]
    },
    get okDisable() {
      if (!this.selectedKeys[0]) {
        return '请选择路径'
      }
      return false
    },
    get disabledKeys() {
      return dirTree
        .filterNodes(item => disabledPaths.includes(item.path))
        .map(item => item.id)
    },
    loadedKeys: [],
    setLoadedKeys(keys) {
      this.loadedKeys = keys
    }
  }))
  const { expandedKeys, selectedKeys } = state

  function onSelect(keys) {
    state.setSelectedKeys(keys)
  }

  function onExpand(keys) {
    state.setExpandedKeys(keys)
  }

  async function createNode(node) {
    state.setExpandedKeys([...new Set([...expandedKeys, node.id])])
    node.unshift(new BaseDirectory())
  }

  function beforeConfirmRename(node: BaseDirectory, name) {
    if (!name) {
      message.error('目录名称不能为空')
      return false
    }

    const validation = validateFilename(name)
    if (validation !== true) {
      message.error(validation)
      return false
    }

    if (node.parent.getDuplicate({ id: node.id, name })) {
      message.error(`${name} 已存在`)
      return false
    }

    return true
  }

  async function confirmRename(node, name) {
    const parentNode = node.parent
    const path = `${parentNode.path}/${name}`
      .replace(/\/+/, '/')
      .replace(/^\//, '')

    await server.mkdir(path)
    await server.sync(parentNode)
    message.success('创建成功')
  }

  function cancelRename(node) {
    if (!node.path) {
      node.parent.removeFirstNode(item => item.id === node.id)
    }
  }

  function mapNode(node) {
    const { isFile, children } = node
    let finalChildren =
      !isFile && (children || []).filter(item => !item.isFile).map(mapNode)

    // hack: always show expand icon
    if (!isFile && finalChildren.length === 0) {
      finalChildren = [
        { key: `__mock__${node.id}`, style: { display: 'none' } }
      ]
    }

    const disabled = state.disabledKeys.includes(node.id)

    return {
      title: (
        <div className='nodeName'>
          <div className='main'>
            {node.path && <span>{node.name}</span>}
            {!node.path && (
              <EditableText
                style={{ display: 'inline-block', width: '200px' }}
                defaultEditing={true}
                defaultValue={node.name}
                showEdit={false}
                beforeConfirm={name => beforeConfirmRename(node, name)}
                onConfirm={name => confirmRename(node, name)}
                onCancel={() => cancelRename(node)}
              />
            )}
          </div>
          <div className='toolbar'>
            {!disabled && node.path && (
              <Tooltip title='新建文件夹'>
                <Icon
                  className='creator'
                  type='folder_new'
                  onClick={e => {
                    e.stopPropagation()
                    createNode(node)
                  }}
                />
              </Tooltip>
            )}
          </div>
        </div>
      ),
      key: node.id,
      children: !disabled && finalChildren,
      disabled,
      icon: expandedKeys.includes(node.id) ? (
        <Icon type='folder_open' />
      ) : (
        <Icon type='folder_close' />
      )
    }
  }

  async function onConfirm() {
    const id = state.selectedKeys[0]
    const node = dirTree.filterFirstNode(item => item.id === id)
    await onOk(node.path.replace(/^\//, '') || '/')
  }

  async function onLoadData(props) {
    const { key } = props
    const node = dirTree.filterFirstNode(item => item.id === key)
    const dir = await server.fetch(node.path, false, isCloud)
    node.children = dir.children
    state.setLoadedKeys([...new Set([...state.loadedKeys, node.id])])
  }

  return (
    <StyledLayout>
      <div className='main'>
        <DirectoryTree
          loadData={onLoadData}
          selectedKeys={selectedKeys}
          expandedKeys={expandedKeys}
          onExpand={onExpand}
          onSelect={onSelect}
          expandAction='doubleClick'
          treeData={dirTree.children.map(item => ({
            ...mapNode(item),
            icon: expandedKeys.includes(item.id) ? (
              <Icon type='folder_open' />
            ) : (
              <Icon type='folder_close' />
            )
          }))}
        />
      </div>
      <Modal.Footer
        className='footer'
        onCancel={onCancel}
        OkButton={
          <Button type='primary' disabled={state.okDisable} onClick={onConfirm}>
            确认
          </Button>
        }
      />
    </StyledLayout>
  )
})

function SelectorWithStore({
  store,
  onCancel,
  onOk,
  disabledPaths = []
}: Props & {
  store?: ReturnType<typeof useStore>
}) {
  const defaultStore = useModel()
  let finalStore = store
  if (!finalStore) {
    finalStore = defaultStore
    defaultStore.initDirTree()
  }

  return (
    <Context.Provider value={finalStore}>
      <DirSelector
        onCancel={onCancel}
        onOk={onOk}
        disabledPaths={disabledPaths}
      />
    </Context.Provider>
  )
}

export const showDirSelector = (options?: {
  disabledPaths?: string[]
  store?: ReturnType<typeof useStore>
}) => {
  return Modal.show({
    title: '选择目录',
    footer: null,
    content: ({ onCancel, onOk }) => (
      <SelectorWithStore {...options} onCancel={onCancel} onOk={onOk} />
    )
  })
}
