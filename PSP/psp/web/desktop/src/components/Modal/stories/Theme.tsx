/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Modal, Button } from '../..'
import styled from 'styled-components'

const StyledExample = styled.div`
  margin: 10px;

  button {
    margin: 5px;
  }
`

const StyledContent = styled.div`
  color: ${({ theme }) => theme.primaryColor};
`

export function Theme() {
  function show() {
    Modal.theme = {
      primaryColor: 'red',
      primaryHighlightColor: 'red',
      secondaryColor: 'pink',
    }
    Modal.show({
      content: <StyledContent>配置弹窗主题</StyledContent>,
      footer: ({ onCancel }) => (
        <Button type='secondary' onClick={onCancel}>
          取消
        </Button>
      ),
    }).finally(() => {
      Modal.theme = null
    })
  }

  return (
    <StyledExample>
      <Button type='primary' onClick={show}>
        模态弹窗
      </Button>
    </StyledExample>
  )
}
