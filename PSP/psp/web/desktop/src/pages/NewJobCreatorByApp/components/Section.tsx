/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { SectionStyle } from './style'

interface IProps {
  title?: React.ReactNode
  icon?: React.ReactNode
  toolbar?: React.ReactNode
  className?: string
  children?: any
}

export const Section = ({
  title,
  icon,
  children,
  toolbar,
  ...props
}: IProps) => (
  <SectionStyle {...props}>
    <div className='section-header'>
      <div className='left'>{title}</div>
      <div className='right'>{toolbar}</div>
    </div>
    <div className='section-content'>{children}</div>
  </SectionStyle>
)
