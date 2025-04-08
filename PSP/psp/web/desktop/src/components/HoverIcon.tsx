/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { Icon } from '@/components'

const HoverIconStyle = styled.div`
  height: 100%;
  padding: 0 8px;

  margin: 0 5px;

  display: flex;
  align-items: center;

  cursor: pointer;

  span.anticon {
    font-size: 18px;
    &.active {
      display: none;
    }
    &.not-active {
      display: inline-block;
    }
  }
  &:hover {
    background-color: #f6f8fa;
    span.anticon {
      &.active {
        display: inline-block;
      }
      &.not-active {
        display: none;
      }
    }
  }
`

interface HoverIconProps {
  type: string
  onClick?: () => void
}

export const HoverIcon = ({ type, onClick, ...rest }: HoverIconProps) => {
  return (
    <HoverIconStyle onClick={onClick} className='hover-icon' {...rest}>
      <Icon
        type={type + '_active'}
        className='active'
        style={{ color: '#0034b4' }}
      />
      <Icon
        type={type}
        className='not-active'
        style={{ color: 'rgba(102,102,102,0.99' }}
      />
    </HoverIconStyle>
  )
}
