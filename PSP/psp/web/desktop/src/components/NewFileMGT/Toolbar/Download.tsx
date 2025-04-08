/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Button } from '@/components'
import { observer } from 'mobx-react-lite'
import { useStore } from '../store'

type Props = {
  disabled: boolean | string
}

export const Download = observer(function Download({ disabled }: Props) {
  const store = useStore()
  const { selectedKeys, dir, server } = store

  async function download() {
    const nodes = dir.filterNodes(item => selectedKeys.includes(item.id))
    const paths = nodes.map(item => item.path)
    const types = nodes.map(item => item.isFile)
    const sizes = nodes.map(item => item.size)

    await server.download(paths, types, sizes)
  }

  return (
    <Button
      disabled={disabled}
      onClick={download}
>
      下载
    </Button>
  )
})
