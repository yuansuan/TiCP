/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { createGlobalStyle } from 'styled-components'

export const GlobalStyle = createGlobalStyle<{ showHeader: boolean }>`
  .ant-modal {
    padding: 0;

    &.confirm {
      .ant-modal-header {
        border-bottom: none;
      }
    }

    .ant-modal-header {
      ${({ showHeader }) => !showHeader && `display: none;`}
      .ant-modal-title {
        color: rgba(0, 0, 0, 0.85);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        padding-right: 40px;
      }
    }

    .ant-modal-body {
      background-color: ${({ theme }: { theme: any }) => theme.backgroundColor};
    }
  }
`
