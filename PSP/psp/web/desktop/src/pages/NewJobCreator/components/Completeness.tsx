/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { CheckOutlined } from '@ant-design/icons'
import { CompletedStyle, UnCompletedStyle } from './style'

interface Props {
  isCompleted: boolean
}

export const Completeness = ({ isCompleted }: Props) =>
  isCompleted ? (
    <CompletedStyle>
      <CheckOutlined />
    </CompletedStyle>
  ) : (
    <UnCompletedStyle />
  )
