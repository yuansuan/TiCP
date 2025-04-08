/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'

export const StyledLayout = styled.div`
  position: absolute;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  z-index: 999;
  display: flex;
  align-items: center;
  justify-content: center;
  animation-name: ${({ animation = {} }: any) => animation.name || ''};
  animation-duration: ${({ animation = {} }: any) =>
    animation.duration || '1s'};
  animation-delay: ${({ animation = {} }: any) => animation.delay || 0};
  animation-direction: ${({ animation = {} }: any) => animation.direction};
  animation-fill-mode: ${({ animation = {} }: any) => animation.fill};
`
