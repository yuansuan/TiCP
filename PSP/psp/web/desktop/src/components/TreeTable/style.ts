/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const Wrapper = styled.div``

export const StyledTable = styled.table`
  width: 100%;
  height: 100%;
  border-collapse: collapse;
  table-layout: fixed;

  .rs-table {
    border: 0;
  }

  .rs-table-hover .rs-table-row {
    &:hover {
      background: ${({ theme }) => theme.backgroundColorHover};

      .rs-table-cell {
        background: ${({ theme }) => theme.backgroundColorHover};
      }

      .rs-table-cell-group {
        background: ${({ theme }) => theme.backgroundColorHover};
      }
    }
  }

  .rs-table-row {
    border-bottom-color: #e8e8e8;

    .rs-table-row-selected {
      .rs-table-cell {
        background: ${({ theme }) => theme.backgroundColorHover};
      }
    }

    .rs-table-row-header {
      background: #f3f5f8;

      .rs-table-cell {
        background: #f3f5f8;
      }
    }
  }

  .rs-table-cell-content {
    font-size: 14px;
    color: rgba(0, 0, 0, 0.65);
  }

  .rs-table-cell .rs-table-cell-content {
    padding: 0px 4px;

    .ant-btn-link {
      padding: 0px 4px;
    }
  }
`

export const LeftPadding = styled.td<{ level: number; indent?: number }>`
  width: ${props => props.level * (props.indent || 24) + 'px'};
  padding: 0 !important;
`
