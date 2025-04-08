/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledRoot = styled.div`
  .ant-back-top {
    bottom: 80px;
  }
  .ant-pagination {
    > .ant-pagination-item-active {
      background-color: ${({ theme }) => theme.primaryColor};

      > a {
        /* color: white; */
      }
    }
  }
`
