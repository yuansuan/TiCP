/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledLayout = styled.div`
  background: #fff;

  > .header {
    font-size: 14px;
    border-bottom: 1px solid ${props => props.theme.borderColorBase};
    padding: 16px 20px;

    .title {
      font-weight: 500;
      position: relative;
      font-family: PingFangSC-Semibold;
      font-size: 16px;
      color: #333333;
    }
  }

  > .body {
    padding: 16px 20px;

    .date-picker {
      margin-bottom: 16px;
    }

    .Pagination {
      margin-top: 20px;
      text-align: center;
    }
  }
`
