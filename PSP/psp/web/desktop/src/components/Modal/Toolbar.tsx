/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { Tooltip } from 'antd'

const StyledLayout = styled.div`
  position: absolute;
  top: 0px;
  right: 56px;
  height: 56px;
  display: flex;
  align-items: center;

  > .actionItem {
    display: flex;
    cursor: pointer;
    height: 100%;
    padding: 0 10px;
    align-items: center;

    &:hover {
      background-color: ${({ theme }) => theme.backgroundColorBase};
    }
  }
`

type ActionItem = {
  tip?: string
  slot: React.ReactNode
}

type Props = {
  className?: string
  style?: React.CSSProperties
  actions?: ActionItem[]
  children?: any
}

export function Toolbar({ children, className, style, actions = [] }: Props) {
  return (
    <StyledLayout className={className} style={style}>
      {children ||
        actions.map(({ tip, slot }) =>
          tip ? (
            <Tooltip title={tip}>
              <div className='actionItem'>{slot}</div>
            </Tooltip>
          ) : (
            <div className='actionItem'>{slot}</div>
          )
        )}
    </StyledLayout>
  )
}
