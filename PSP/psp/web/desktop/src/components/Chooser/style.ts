/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledChooser = styled.div`
  .search {
    width: 150px;
  }

  .allSelector {
    padding: 8px 0;
    border-bottom: 1px solid ${({ theme }) => theme.borderColorBase};
  }

  .list {
    margin-top: 5px;
    max-height: 240px;
    overflow: auto;

    .item {
      margin: 5px 0;

      .name {
        display: inline-block;
        vertical-align: middle;
        margin-left: 5px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        width: calc(100% - 21px);
      }
    }
  }
`
