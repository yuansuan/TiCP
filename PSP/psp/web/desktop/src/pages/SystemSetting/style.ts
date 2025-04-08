/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledLayout = styled.div`
  padding: 0 20px 20px 20px;
  section {
    margin-bottom: 20px;

    h1 {
      font-size: 16px;
      font-family: 'PingFangSC-Medium';
      color: #333333;
    }

    > .section-bottom {
      padding-left: 20px;

      > .row {
        display: flex;
        flex-flow: row nowrap;
        align-items: center;

        font-size: 14px;
        margin-bottom: 4px;
        height: 24px;
      }
    }
  }
`
