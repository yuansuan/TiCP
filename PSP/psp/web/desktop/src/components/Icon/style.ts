/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'
import Icon from '@ant-design/icons'

export const StyledIcon: typeof Icon = styled(Icon)`
  width: 1em;
  height: 1em;
  font-size: 24px;

  > svg {
    width: 100%;
    height: 100%;
    fill: currentColor;
  }

  &.disabled {
    cursor: not-allowed;
    pointer-events: none;
    opacity: 0.5;
  }
`
