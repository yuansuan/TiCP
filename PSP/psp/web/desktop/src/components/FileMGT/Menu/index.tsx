/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Tree, message, Input } from 'antd'
import { useStore } from '../store'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Button, Icon } from '@/components'
import { Title } from './Title'
import { BaseDirectory } from '@/utils/FileSystem'
import { useResize, env } from '@/domain'
import difference from 'lodash/difference'
import { getUrlParams } from '@/utils'
const { DirectoryTree } = Tree
const StyledLayout = styled.div`
  height: 100%;
  display: flex;
  flex-direction: column;

  .ant-tree.ant-tree-directory .ant-tree-treenode {
    overflow: hidden;

    .ant-tree-node-content-wrapper {
      .ant-tree-iconEle svg {
        fill: ${({ theme }) => theme.primaryColor};
      }

      &.ant-tree-node-selected {
        .ant-tree-iconEle svg {
          fill: white;
        }

        .menu_toolbar {
          color: white;
        }
      }
    }
  }

  > .menu_toolbar {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    background-color: ${({ theme }) => theme.backgroundColorBase};
    padding: 4px;
    z-index: 99;
    > * {
      width: 100%;
      margin: 2px 0;
    }

    .ant-btn {
      width: 100%;
    }
  }

  > .menu_main {
    flex: 1;
    width: 100%;
    overflow: auto;
  }
`

type Props = {
  isSyncToLocal:boolean
}

export const Menu = observer(function Menu({isSyncToLocal}: Props) {
  const store = useStore()
  const isCloud = getUrlParams()?.isCloud && JSON.parse(getUrlParams()?.isCloud)

  const { dirTree, currentNode, getWidget, server } = store
  const state = useLocalStore(() => ({
    expandedKeys: [],
    loadedKeys: [],
    searchKey: '',
    update(
      props: Partial<{
        expandedKeys: string[]
        loadedKeys: string[]
        searchKey: string
      }>
    ) {
      Object.assign(this, props)
    }
  }))
  const [rect, ref] = useResize()
  const { expandedKeys, searchKey } = state

  // auto expand ancestral nodes when node is selected
  useEffect(() => {
    const keys = []
    let node = currentNode
    while (node?.parent) {
      keys.push(currentNode.parent.id)
      node = node.parent
    }

    if (keys.length > 0 && difference(keys, state.expandedKeys)) {
      state.update({
        expandedKeys: [...new Set([...state.expandedKeys, ...keys])]
      })
    }
  }, [currentNode?.id])

  function updateNodeId(keys) {
    // do nothing when select newly created node
    const node = dirTree.filterFirstNode(item => item.id === keys[0])
    if (!node?.path) {
      return
    }
    store.setNodeId(keys[0])
  }

  function onExpand(keys, { expanded }) {
    state.update({
      expandedKeys: keys
    })
    if (!expanded) {
    }
  }

  async function createNode() {
    if (!currentNode) {
      message.error('父节点不存在')
      return
    }

    state.update({
      expandedKeys: [...new Set([...expandedKeys, currentNode.id])]
    })
    currentNode.unshift(new BaseDirectory())
  }

  function mapNode(node, editable = true) {
    const { isFile, children } = node
    let finalChildren =
      !isFile &&
      (children || []).filter(item => !item.isFile).map(item => mapNode(item))

    // hack: always show expand icon
    if (!isFile && finalChildren.length === 0) {
      finalChildren = [
        { key: `__mock__${node.id}`, style: { display: 'none' } }
      ]
    }

    return {
      title: editable ? <Title node={node} /> : node.name,
      key: node.id,
      children: finalChildren,
      icon: expandedKeys.includes(node.id) ? (
        <Icon type='folder_open' />
      ) : (
        <Icon type='folder_close' />
      )
    }
  }

  async function onLoadData(props) {
    const { key } = props
    const node = dirTree.filterFirstNode(item => item.id === key) || {}
    const rawCloud = isSyncToLocal ? false : isCloud
    const dir = await server.fetch(node?.path, false, rawCloud)
    node.children = dir.children
    state.update({
      loadedKeys: [...new Set([...state.loadedKeys, node.id])]
    })
  }

  return (
    <StyledLayout>
      <div className='menu_toolbar'>
        <div className='search'>
          <Input
            allowClear
            placeholder='搜索'
            maxLength={64}
            value={searchKey}
            onChange={e => {
              state.update({
                searchKey: e.target.value,
                expandedKeys: [
                  ...new Set([...state.expandedKeys, dirTree.children?.[0]?.id])
                ]
              })
            }}
          />
        </div>
        {/* {getWidget('newDir') || (
          <Button
            icon='folder_new'
            onClick={createNode}
            disabled={!currentNode && '请选择节点'}>
            新建文件夹
          </Button>
        )} */}
      </div>
      <div className='menu_main' ref={ref}>
        <DirectoryTree
          loadData={onLoadData}
          height={rect.height}
          selectedKeys={[currentNode?.id]}
          expandedKeys={expandedKeys}
          onExpand={onExpand}
          onSelect={updateNodeId}
          treeData={dirTree.children.map(item =>
            mapNode(
              {
                ...item,
                children: item.children.filter(item => {
                  if (!searchKey) {
                    return true
                  }
                  return item.name
                    .toLowerCase()
                    .includes(searchKey.toLowerCase())
                })
              },
              false
            )
          )}
          expandAction='doubleClick'
        />
      </div>
    </StyledLayout>
  )
})
