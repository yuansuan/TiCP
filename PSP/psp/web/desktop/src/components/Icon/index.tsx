/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { forwardRef } from 'react'
import { IconComponentProps } from '@ant-design/icons/lib/components/Icon'
import { StyledIcon } from './style'
// import './iconfont'

type YSIconProps = IconComponentProps & {
  type?: string
  disabled?: boolean
}

function YSIcon(props: YSIconProps, ref) {
  const { type, disabled, ...finalProps } = props

  // map disabled props to className
  if (disabled) {
    finalProps.className = [
      ...new Set([...(finalProps.className || '').split(' '), 'disabled']),
    ].join(' ')
  }

  // add ysicon class
  finalProps.className = [
    ...new Set([...(finalProps.className || '').split(' '), 'ysicon']),
  ].join(' ')

  // map type to component
  if (type && !finalProps.component) {
    finalProps.component = () => (
      <svg width='1em' height='1em' fill='currentColor'>
        <use href={`#${type}`} />
      </svg>
    )
  }

  return <StyledIcon ref={ref} {...finalProps} />
}

export default forwardRef(YSIcon)
