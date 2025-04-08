/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Icon } from '@/components'
import styled from 'styled-components'

const HoverableIconStyle = styled.span`
  display: flex;
  align-items: center;

  span:first-child {
    display: inline-block;
  }
  span:nth-child(2) {
    display: none;
  }

  &:hover {
    span:first-child {
      display: none;
    }
    span:nth-child(2) {
      display: inline-block;
    }
  }

  .anticon {
    font-size: 18px;
    padding: 0 9px;
    box-sizing: content-box;
    cursor: pointer;
    color: ${props => props.theme.primaryColor};
    z-index: 9999;
  }

  &.folder_mine span.anticon {
    font-size: 22px;
  }
`

interface Props {
  type: string
  onClick?: () => void
  style?: React.CSSProperties
}

export const HoverableIcon = ({ type, onClick, style }: Props) => (
  <HoverableIconStyle className={type} style={style}>
    <Icon type={type} onClick={onClick} />
    <Icon type={type + '_active'} onClick={onClick} />
  </HoverableIconStyle>
)
