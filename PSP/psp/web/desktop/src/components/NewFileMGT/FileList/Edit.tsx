/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Button } from 'antd'
import { showTextEditor } from '@/components/TextEditor'
import { newBoxServer } from '@/server'
import { useStore } from '../store'
import styled from 'styled-components'

type Props = {
  nodeId: string
  readonly?: boolean
  disabled?: boolean
}

const TableLinkBtn = styled(Button)`
  padding: 0;
`
export function Edit({ nodeId, disabled, readonly = false }: Props) {
  const store = useStore()
  const { dir } = store
  const node = dir.filterFirstNode(item => item.id === nodeId)

  function edit() {
    showTextEditor({
      path: node.path,
      fileInfo: {
        ...node
      },
      readonly,
      boxServerUtil: newBoxServer
    })
  }

  return (
    <TableLinkBtn type='link' disabled={disabled} onClick={edit}>
      {readonly ? '查看' : '编辑'}
    </TableLinkBtn>
  )
}
