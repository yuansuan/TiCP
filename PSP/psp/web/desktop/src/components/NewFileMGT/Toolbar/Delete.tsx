/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import { Button, Modal } from '@/components'
import { useStore } from '../store'
import { message } from 'antd'
import { Http } from '@/utils'
import { env } from '@/domain'
import SelectedFiles from '@/components/SelectItems'
type Props = {
  disabled?: string | boolean
}

export const Delete = observer(function Delete({ disabled }: Props) {
  const store = useStore()
  const { selectedKeys, server } = store
  const [fetch] = store.useRefresh()

  async function deleteNodes(ids: string[]) {
    const nodes = store.dir.filterNodes(item => ids.includes(item.id))
    await Modal.showConfirm({
      title: '删除文件',
      content: <SelectedFiles selectedName={nodes.map(item => item.name)}/>
    })

    await server.delete(nodes.map(item => item.path))

    const size = nodes
      .map(item => item.size)
      .reduce((previousValue, currentValue) => previousValue + currentValue, 0)
    const names = nodes.map(item => item.path.split('/').pop())
    const name =
      nodes.length === 1
        ? names[0]
        : `[批量删除]${
            nodes.length > 2
              ? names.slice(0, 2).join(',') + '等'
              : names.join(',')
          }`

    await fetch()
    message.success('文件删除成功')
  }

  async function onDelete() {
    await deleteNodes(selectedKeys)
  }

  return (
    <Button disabled={disabled} onClick={onDelete}>
      删除
    </Button>
  )
})
