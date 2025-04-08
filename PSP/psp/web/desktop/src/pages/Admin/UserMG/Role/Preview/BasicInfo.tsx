import React from 'react'
import { BasicInfoWrapper } from './style'

interface IProps {
  title: string
  children: React.ReactNode | string
  className?: string
}

export default function BasicInfo({ title, children, className }: IProps) {
  return (
    <BasicInfoWrapper>
      <span className='header'>{title}</span>
      <div className={`body ${className}`}>{children}</div>
    </BasicInfoWrapper>
  )
}
