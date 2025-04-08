/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const Wrapper = styled.div``

export const StyledTable = styled.table`
  width: 100%;
  border-collapse: collapse;

  th,
  td {
    padding: 16px;
    border-bottom: 1px solid #f0f0f0;
  }

  thead {
    th {
      text-align: left;
      background: #fafafa;
    }
  }

  tbody > tr:hover > td {
    background: #f5f5f5;
  }

  th.padding,
  td.padding {
    padding: 0;
  }
`

export const LeftPadding = styled.td<{ level: number; indent?: number }>`
  width: ${props => props.level * (props.indent || 24) + 'px'};
  padding: 0 !important;
`
