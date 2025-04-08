/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import Icon from '../Icon'
import ModalFooter, { FooterProps } from './Footer'

const StyledLayout = styled.div`
  > .title {
    display: flex;
    align-items: center;

    > .icon {
      color: ${({ theme }) => theme.warningColor};

      .anticon {
        font-size: 30px;
      }
    }

    > .text {
      font-size: 16px;
      font-weight: bolder;
      margin-left: 16px;
      color: rgba(0, 0, 0, 0.85);
    }
  }

  > .content {
    margin: 12px 20px 34px 54px;
    font-size: 14px;
    color: rgba(0, 0, 0, 0.65);
    line-height: 22px;
  }

  > .footer {
    margin-bottom: 24px;
    margin-right: 8px;
  }
`

type Props = {
  title?: string
  content?:
    | string
    | React.ReactNode
    | React.ComponentType<{
        onCancel?: (data?: any) => void | Promise<any>
        onOk?: (data?: any) => void | Promise<any>
      }>
} & FooterProps

export function ConfirmContent({
  title = '确认弹窗',
  content = '确认执行该操作吗？',
  ...props
}: Props) {
  return (
    <StyledLayout>
      <div className='title'>
        <div className='icon'>
          <Icon type='question_circle' />
        </div>
        <div className='text'>{title}</div>
      </div>
      <div className='content'>{content}</div>
      <div className='footer'>
        <ModalFooter {...props} />
      </div>
    </StyledLayout>
  )
}
