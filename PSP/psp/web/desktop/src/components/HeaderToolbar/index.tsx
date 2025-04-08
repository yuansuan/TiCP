/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { UserInfo } from './UserInfo'
import { Message } from './Message'
import styled from 'styled-components'

const HeaderToolbarStyle = styled.div`
  display: flex;
  height: 100%;
  justify-content: end;
  > .right {
    display: flex;
    position: absolute;
    left: 20px;
    top: 12px;
    /* margin-left: auto; */
    color: #282626;

    > .AreaSelectContainer {
      display: flex;
      justify-content: center;
      align-items: center;
      .anticon {
        font-size: 12px;
        margin: 0;
      }
    }

    .anticon {
      font-size: 16px;
      color: #000;
      margin: 4px;

      &.active {
        color: ${props => props.theme.primaryColor};
      }
    }

    > * {
      cursor: pointer;

      &:hover,
      &.ant-dropdown-open {
        background-color: #fafafa;
        border-radius: 4px;
      }
    }
  }
`

export const HeaderToolbar = observer(function HeaderToolbar() {
  return (
    <HeaderToolbarStyle>
      <div className='right'>
        <Message />
        <UserInfo type='inside' />
      </div>
    </HeaderToolbarStyle>
  )
})

export { UserInfo } from './UserInfo'
export { Message } from './Message'
export { Uploader } from './Uploader'
export { Balance } from './Balance'
