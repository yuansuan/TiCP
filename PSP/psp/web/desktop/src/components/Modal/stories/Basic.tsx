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

export const Basic = () => (
  <StyledContainer>
    <Button type='primary' onClick={() => Modal.showConfirm()}>
      确认弹窗
    </Button>
    <Button type='primary' onClick={() => Modal.show()}>
      模态弹窗
    </Button>
  </StyledContainer>
)
