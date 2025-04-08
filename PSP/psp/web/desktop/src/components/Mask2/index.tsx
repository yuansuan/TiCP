/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Keyframes } from 'styled-components'
import { StyledLayout } from './style'
import { SpinMask } from './SpinMask'

export type MaskProps = {
  style?: React.CSSProperties
  theme?: 'white' | 'black'
  animation?: EffectTiming & {
    name: Keyframes
  }
  children?: React.ReactNode
}

function Mask(props: MaskProps) {
  const { children, theme, style, ...rest } = props
  const finalStyle = {
    ...(theme &&
      (theme === 'black'
        ? {
            backgroundColor: 'rgba(0,0,0,0.6)',
            color: 'white',
          }
        : theme === 'white'
        ? {
            backgroundColor: 'rgba(255,255,255,0.6)',
            color: 'black',
          }
        : null)),
    ...style,
  }

  return (
    <StyledLayout style={finalStyle} {...rest}>
      {children}
    </StyledLayout>
  )
}

Mask.Spin = SpinMask

export default Mask
