/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Modal, Button } from '../..'
import styled from 'styled-components'

const StyledContainer = styled.div`
  margin: 10px;

  button {
    margin: 5px;
  }
`

export const Async = () => (
  <StyledContainer>
    <Button
      type='primary'
      onClick={() =>
        Modal.showConfirm({
          onCancel: () =>
            new Promise(resolve => {
              setTimeout(resolve, 1000)
            }),
          onOk: () =>
            new Promise(resolve => {
              setTimeout(resolve, 1000)
            }),
        })
      }>
      async
    </Button>

    <Button
      type='primary'
      onClick={() =>
        Modal.showConfirm({
          onCancel: () =>
            new Promise(resolve => {
              setTimeout(resolve, 1000)
            }),
          onOk: () =>
            new Promise(resolve => {
              setTimeout(resolve, 1000)
            }),
        })
      }>
      Promise
    </Button>
  </StyledContainer>
)
