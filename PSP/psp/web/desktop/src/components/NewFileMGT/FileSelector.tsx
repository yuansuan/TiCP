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
  isTemDirPath?: boolean
}

export const FileSelector = observer(function FileSelector({
  onCancel,
  onOk,
  modelProps,
  isMoreInfo = false,
  isTemDirPath = false
}: Props) {
  const model = useModel(
    modelProps || {
      widgets: () => ['testFile']
    }
  )

  function confirm() {
    const { selectedKeys, dir } = model
    const nodes = dir.filterNodes(item => selectedKeys.includes(item.id))

    if (isTemDirPath || (nodes.length === 1 && !nodes[0]?.isFile)) {
      onOk(isMoreInfo ? nodes : nodes.map(item => item.path))
      return
    }

    Modal.confirm({
      title: '工作目录选择提醒',
      content:
        '原地计算模式下只能选择一个工作目录, 多选情况下默认选取展示顺序第一个文件目录数据内容! 请确保选择的是文件目录!',
      okText: '确认',
      cancelText: '取消',
      onOk: () => {
        onOk(isMoreInfo ? nodes : nodes.map(item => item.path))
      }
    })
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

export const showFileSelector = async function (
  isMoreInfo = false,
  isTemDirPath = false
) {
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
      <FileSelector
        onCancel={onCancel}
        onOk={onOk}
        isMoreInfo={isMoreInfo}
        isTemDirPath={isTemDirPath}
      />
    )
  })
}
