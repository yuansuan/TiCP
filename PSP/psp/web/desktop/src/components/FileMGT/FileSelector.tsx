/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Suite, useModel, ModelProps } from '.'
import { Modal, Button } from '@/components'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'

const StyledLayout = styled.div`
  height: 100%;
  padding: 20px 20px 72px 20px;
  box-sizing: border-box;

  > .body {
    height: 100%;
  }

  > .footer {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};
    padding: 10px 0;
  }
`

type Props = {
  onCancel: () => void
  onOk: (paths: string[]) => void
  modelProps?: ModelProps
  isMoreInfo?: boolean
}

export const FileSelector = observer(function FileSelector({
  onCancel,
  onOk,
  modelProps,
  isMoreInfo = false
}: Props) {
  const model = useModel(
    modelProps || {
      widgets: () => ['testFile']
    }
  )

  function confirm() {
    const { selectedKeys, dir } = model
    const nodes = dir.filterNodes(item => selectedKeys.includes(item.id))

    onOk(isMoreInfo ? nodes : nodes.map(item => item.path))
  }

  return (
    <StyledLayout>
      <div className='body'>
        <Suite model={model} jobManger={true} />
      </div>
      <Modal.Footer
        className='footer'
        onCancel={onCancel}
        OkButton={
          <Button
            type='primary'
            disabled={model.selectedKeys.length === 0 ? '未选择文件' : false}
            onClick={confirm}>
            确认
          </Button>
        }
      />
    </StyledLayout>
  )
})

export const showFileSelector = async function (isMoreInfo = false) {
  return await Modal.show({
    title: '我的文件',
    width: 1200,
    bodyStyle: {
      backgroundColor: 'white',
      height: 650,
      padding: 0
    },
    footer: null,
    content: ({ onCancel, onOk }) => (
      <FileSelector onCancel={onCancel} onOk={onOk} isMoreInfo={isMoreInfo} />
    )
  })
}
