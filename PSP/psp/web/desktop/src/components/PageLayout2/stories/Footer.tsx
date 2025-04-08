/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Button } from 'antd'
import PageLayout from '..'
import styled from 'styled-components'

const StyledBtn = styled(Button)`
  font-family: PingFangSC-Regular;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.65);
  text-align: left;
  background-color: rgba(24, 144, 255, 0.13);
  width: 100%;
  overflow: hidden;
`

export const Footer = () => (
  <PageLayout
    SiderFooter={
      <StyledBtn
        className='sider-footer-btn'
        type='link'
        onClick={() => (window.location.href = '/')}>
        Menu-Footer组件
      </StyledBtn>
    }
  />
)
