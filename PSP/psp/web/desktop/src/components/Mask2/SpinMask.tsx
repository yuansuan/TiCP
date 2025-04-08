/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Spin } from 'antd'
import { SpinProps } from 'antd/lib/spin'
import { keyframes } from 'styled-components'

import Mask, { MaskProps } from './index'

const fadeIn = keyframes`
  0% { opacity: 0; }
  50% {opacity: 0;}
  100% { opacity: 1; }
`

type Props = MaskProps & {
  spinProps?: SpinProps
}

export function SpinMask(props: Props) {
  const { spinProps, children, ...rest } = props

  return (
    <Mask {...rest} animation={{ name: fadeIn, duration: '1s' }}>
      {children || <Spin {...spinProps} />}
    </Mask>
  )
}
