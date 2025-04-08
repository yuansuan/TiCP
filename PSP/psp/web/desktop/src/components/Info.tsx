/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'
import React from 'react'
import { InfoCircleOutlined } from '@ant-design/icons'

const InfoStyle = styled.div`
  border: 1px solid #8edcff;
  padding: 10px;
  display: flex;
  align-items: flex-start;
  background-color: #e3f8ff;

  p {
    margin: 0;
    color: #6b7478;
  }

  span.info-circle {
    margin-right: 10px;
    color: #3f91ff;
    font-size: 1.5em;
  }
`

interface InfoProps {
  content: string
  className?: string
}

export const Info = ({ content, ...rest }: InfoProps) => (
  <InfoStyle {...rest}>
    <InfoCircleOutlined className={'info-circle'} />
    <p>{content}</p>
  </InfoStyle>
)
