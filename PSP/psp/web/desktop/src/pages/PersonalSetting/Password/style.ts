/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledEditor = styled.div`
  .body {
    padding-bottom: 40px;
  }

  .footer {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    padding: 10px 0;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};
  }
`
