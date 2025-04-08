/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledColumnManager = styled.div`
  min-width: 150px;

  .item {
    display: flex;
    font-size: 12px;
    align-items: center;

    &:hover {
      background: #f2f6ff;

      .move {
        visibility: visible;
      }
    }

    .name {
      margin-left: 8px;
      width: 60px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .move {
      visibility: hidden;
      margin-left: auto;

      > * {
        cursor: pointer;

        &:hover {
          color: #63a9ff;
        }

        &:first-child {
          margin-right: 8px;
        }
      }
    }
  }
`
