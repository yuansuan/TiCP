/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const Wrapper = styled.div`
  display: flex;
  padding: 10px 20px;
  width: 100%;

  .name {
    > svg {
      margin-right: 3px;
    }

    &.slave {
      padding-left: 15px;
    }
  }

  .progress {
    padding: 0 5px;
  }

  .operate {
    .icon {
      margin: 0 5px;
    }
  }

  .rs-table-row-expanded {
    padding: 0;
  }
`

export const SlaveFileTable = styled.div`
  .rs-table {
    border: none;
  }
`
