/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { notification } from 'antd'
import { ArgsProps } from 'antd/lib/notification'
import {
  CloseCircleFilled,
  InfoCircleFilled,
  CheckCircleFilled,
  ExclamationCircleFilled,
} from '@ant-design/icons'

notification.config({
  placement: 'topRight',
  top: 0,
})

const style = {
  minWidth: 400,
  maxWidth: 568,
  boxShadow: 'unset',
}

export default {
  error: ({ ...props }: ArgsProps) => {
    return new Promise(resolve => {
      notification.error({
        ...props,
        duration: 0,
        onClose: () => {
          resolve()
        },
        className: 'notification',
        style: {
          ...style,
          backgroundColor: '#fff1f0',
          border: '1px solid #ffa39e',
          margin: 0,
        },
        icon: <CloseCircleFilled style={{ color: '#f5222e' }} />,
      })
    })
  },
  warning: ({ ...props }: ArgsProps) => {
    return new Promise(resolve => {
      notification.warning({
        ...props,
        duration: 0,
        onClose: () => {
          resolve()
        },
        className: 'notification',
        style: {
          ...style,
          backgroundColor: '#fffbe6',
          border: '1px solid #ffe58f',
        },
        icon: <InfoCircleFilled style={{ color: '#f9bf02' }} />,
      })
    })
  },
  info: ({ ...props }: ArgsProps) => {
    return new Promise(resolve => {
      notification.warning({
        ...props,
        duration: 0,
        onClose: () => {
          resolve()
        },
        className: 'notification',
        style: {
          ...style,
          backgroundColor: '#E6F7FF',
          border: '1px solid #91D5FF',
        },
        icon: <ExclamationCircleFilled style={{ color: '#398FE9' }} />,
      })
    })
  },
  success: ({ ...props }: ArgsProps) => {
    notification.success({
      ...props,
      duration: 0,
      className: 'notification',
      style: {
        ...style,
        backgroundColor: '#F6FFED',
        border: '1px solid #B7EB8F',
      },
      icon: <CheckCircleFilled style={{ color: '#52C51A' }} />,
    })
  },
}
