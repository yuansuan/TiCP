/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Icon } from '@/components'
import { FileActionStyle } from './style'

interface Props {
  children: string
  icon?: string
  onClick: () => void
}

export const FileAction = ({ children, icon, onClick }: Props) => (
  <FileActionStyle onClick={onClick}>
    {icon && (
      <>
        <Icon type={icon} className={icon} />
        <Icon type={icon + '_active'} className={icon} />
      </>
    )}
    {children}
  </FileActionStyle>
)
