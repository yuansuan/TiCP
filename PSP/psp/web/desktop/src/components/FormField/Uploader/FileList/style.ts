/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const FileListStyle = styled.div`
  position: relative;
  width: 100%;
  max-height: 460px;
  padding: 0;
  overflow: auto;
  background-color: white;

  span.disabled {
    color: rgba(0, 0, 0, 0.25);
  }

  tr:hover td .actions {
    display: flex;
  }

  a.delete {
    color: #f5222d;
    &:hover {
      color: #f5222d;
    }
  }

  > .mask-wrapper {
    > div {
      background: rgba(0, 0, 0, 0.1);
    }
  }
`
