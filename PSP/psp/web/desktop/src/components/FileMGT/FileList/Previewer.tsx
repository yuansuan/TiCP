/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useMemo } from 'react'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useStore } from '../store'
import { Icon } from '@/components'
import { previewImage } from '@/components'
import { Tooltip } from 'antd'

type Props = {
  nodeId: string
}

export const Previewer = observer(
  function Previewer({ nodeId }: Props, ref: any) {
    const { dir, server } = useStore()
    const state = useLocalStore(() => ({
      url: '',
      setUrl(url) {
        this.url = url
      }
    }))
    const node = useMemo(
      () => dir.filterFirstNode(item => item.id === nodeId),
      [dir, nodeId]
    )

    async function preview() {
      try {
        if (node) {
          const url =
            state.url ||
            (await server.getFileUrl([node.path], [true], [node.size], true))
          state.setUrl(url)
          previewImage({ fileName: node.name, src: url })
        }
      } finally {
      }
    }

    return (
      <Tooltip title='预览图片'>
        <Icon ref={ref} type='img_preview' onClick={preview}>
          预览
        </Icon>
      </Tooltip>
    )
  },
  {
    forwardRef: true
  }
)
