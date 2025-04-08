/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Button } from '@/components'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'
import { message } from 'antd'
import { showDirSelector } from '../DirSelector'
import { showFailure } from '../Failure'

type Props = {
  disabled?: string | boolean
}

export const Move = observer(function Move({ disabled }: Props) {
  const store = useStore()
  const { selectedKeys, server, dir, dirTree } = store
  const [refresh] = store.useRefresh()

  async function moveTo(path) {
    const nodes = dir.filterNodes(item => selectedKeys.includes(item.id))

    // check duplicate
    const targetDir = await server.fetch(path)
    const rejectedNodes = []
    const resolvedNodes = []

    nodes.forEach(item => {
      if (targetDir.getDuplicate({ id: undefined, name: item.name })) {
        rejectedNodes.push(item)
      } else {
        resolvedNodes.push(item)
      }
    })

    const destMoveNodes = [...resolvedNodes]
    if (rejectedNodes.length > 0) {
      const coverNodes = await showFailure({
        actionName: '移动',
        items: rejectedNodes
      })
      if (coverNodes.length > 0) {
        // coverNodes 中的要删除
        await server.delete(coverNodes.map(item => `${path}/${item.name}`))
        destMoveNodes.push(...coverNodes)
      }
    }

    if (destMoveNodes.length > 0) {
      await server.move(
        destMoveNodes.map(item => [item.path, `${path}/${item.name}`])
      )
      await refresh()
      const targetNode = dirTree.filterFirstNode(item => item.path === path, {
        self: true
      })

      const { children } = await server.fetch(targetNode.path)
      targetNode.children = children
      message.success('文件移动成功')
    }
  }

  async function onMove() {
    const selectedNodes = dir.filterNodes(item =>
      selectedKeys.includes(item.id)
    )

    const path = await showDirSelector({
      disabledPaths: selectedNodes.map(item => item.path)
    })
    await moveTo(path)
  }

  return (
    <Button disabled={disabled} onClick={onMove}>
      移动
    </Button>
  )
})
