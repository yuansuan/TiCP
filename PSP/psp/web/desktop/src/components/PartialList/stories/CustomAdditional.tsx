/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Popover as AntPopover } from 'antd'
import styled from 'styled-components'
import { PartialList, Icon } from '../..'

const StyledAdditionalList = styled.div`
  max-width: 145px;

  .title {
    font-family: 'PingFangSC-Medium';
    font-size: 14px;
    color: #000000;
    letter-spacing: 0;
    margin-bottom: 5px;
  }

  .list {
    display: flex;
    flex-wrap: wrap;
  }
`

function AdditionalList({ items, title }) {
  return (
    <StyledAdditionalList>
      {title ? <div className='title'>{title}</div> : null}
      <div className='list'>
        {items.map((item, index) => (
          <span className='item' key={index}>
            {item}
            {index === items.length - 1 ? null : '；'}
          </span>
        ))}
      </div>
    </StyledAdditionalList>
  )
}

const items = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18]

export const CustomAdditional = () => (
  <PartialList
    style={{ width: 200, margin: 20 }}
    items={items}
    Additional={({ visibleIndex }) => (
      <AntPopover
        trigger='hover'
        content={
          <AdditionalList
            title='自定义标题'
            items={items.slice(Math.max(visibleIndex, 0))}
          />
        }
        placement='bottom'>
        <Icon type='more' />
      </AntPopover>
    )}
  />
)
