/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Modal } from '@/components'
import { QRCodeModal } from './QRCodeModal'
import { JobQrCodeModal } from './JobQRCodeModal'

export type Props = {
  descriptionNode?: React.ReactNode
  fetchQRCodeFunc: () => Promise<{
    data: {
      qrcodeUrl: string
      expireSeconds: number
    }
  }>
  validConfig: {
    // 轮询判断是否绑定成功的值
    getCheckValue?: () => any
    // 轮询成功后触发此方法
    afterCheckedFunc?: () => void
    // 轮询间隔 可选
    pollingInterval?: number
  }
  afterOk?: () => void
  afterCancel?: () => void
  hideOk?: boolean
  title?: string
}

export async function showJobQRCodeModal(props?: Props) {
  return await Modal.show({
    content: ({ onCancel, onOk }) => (
      <JobQrCodeModal onCancel={onCancel} onOk={onOk} {...props} />
    ),
    width: 600,
    bodyStyle: { width: '100%', height: '400px' },
    zIndex: 1200,
    footer: null,
    closable: false,
    showHeader: false,
  })
}
export async function showQRCodeModal(props?: Props) {
  return await Modal.show({
    title: props?.title || '提示',
    content: ({ onCancel, onOk }) => (
      <QRCodeModal onCancel={onCancel} onOk={onOk} {...props} />
    ),
    width: 600,
    bodyStyle: { width: '100%', height: '100%' },
    zIndex: 1200,
    footer: null,
    closable: false,
  })
}
