/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { Tooltip } from 'antd'

const Wrapper = styled.div`
  display: flex;
  align-items: center;

  .point,
  .point::before,
  .point::after {
    position: absolute;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    content: '';
  }

  .point::before {
    animation: ${props => (props.color === '#63B03D' ? 'scale' : '')} 3s
      infinite;
  }

  .point::after {
    animation: ${props => (props.color === '#63B03D' ? 'scale2' : '')} 3s
      infinite;
  }

  @keyframes scale {
    0% {
      transform: scale(1);
      opacity: 0.9;
    }

    100% {
      transform: scale(2);
      opacity: 0;
    }
  }

  @keyframes scale2 {
    0% {
      transform: scale(1);
      opacity: 0.9;
    }

    100% {
      transform: scale(2);
      opacity: 0;
    }
  }

  .point,
  .point::before,
  .point::after {
    /* 设置颜色 */
    background-color: ${props => props.color};
  }

  .icon {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    margin-right: 4px;
    position: relative;
    box-sizing: content-box;
  }

  .text {
    margin-left: 4px;
  }

  .icon-right {
    margin-left: 4px;

    .anticon {
      height: 12px;
      width: 12px;
      position: absolute;
      top: 21px;
    }
  }
`

interface IProps {
  type: string
}

const colorMap = {
  success: {
    color: '#63B03D'
  },
  close: {
    color: '#ccc'
  },
  error: {
    color: '#EF5350'
  }
}

export function Status(props: IProps) {
  return (
    <Wrapper color={colorMap[props.type]?.color || '#63B03D'}>
      <Tooltip title={`${props.type === 'close' ? '中止' : '运行中'}`}>
        <div
          className='point'
          style={{
            background: colorMap[props.type]?.color
          }}></div>
      </Tooltip>
    </Wrapper>
  )
}
