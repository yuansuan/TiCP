/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledLayout = styled.div`
  display: flex;
  flex-direction: column;

  > .info {
    display: flex;
    flex-direction: row-reverse;
    align-items: center;
    > div {
      padding-right: 1em;
      border-right: 1px solid rgba(0, 0, 0, 0.09);
    }
  }

  > .footer {
    display: flex;
    flex: 1;
    justify-content: center;
    align-items: center;
    margin: 10px;
  }
`

export const StatusWrapper = styled.div`
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
