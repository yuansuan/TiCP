/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import Icon from './Icon'

const Wrapper = styled.div`
  display: flex;
  align-items: center;

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
  text: string
  type: 'success' | 'warn' | 'error'
}

const colorMap = {
  success: {
    color: '#63B03D',
    borderColor: '#D7F9C7'
  },
  warn: {
    color: '#FF9100',
    borderColor: '#FDEFC7'
  },
  error: {
    color: '#EF5350',
    borderColor: '#F9D9D9'
  },
  primary: {
    color: '#2A8FDF',
    borderColor: '#C7E3F9'
  },
  4: {
    color: '#C5C5C5',
    borderColor: '#E6E4E4'
  }
}

export function Status(props: IProps) {
  const getIcon = text => {
    if (text === '运行中') {
      return <Icon type='running' />
    } else if (text === '取消中') {
      return <Icon type='loading' />
    } else {
      return null
    }
  }
  return (
    <Wrapper>
      <div
        className='icon'
        style={{
          background: colorMap[props.type]?.color,
          border: `2px solid ${colorMap[props.type]?.borderColor}`
        }}
      />
      <div className='text'>{props.text ?? '--'}</div>
      <div className='icon-right'>{getIcon(props.text)}</div>
    </Wrapper>
  )
}
