/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { Modal } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { StyledItem } from '../style'
import { showJobQRCodeModal } from '@/components'
import moment from 'moment'
import { message } from 'antd'
import { userServer } from '@/server'

export const Wechat = observer(function WeChat() {
  const store = useLocalStore(() => ({
    bindInfo: {
      notification_activated: undefined,
      wechat_nickname: undefined,
      wechat_openid: '',
    },
    setBindInfo(v) {
      this.bindInfo = {
        ...this.bindInfo,
        ...v,
      }
    },
  }))

  useEffect(() => {
    fetchBindInfo()
  }, [])

  async function fetchBindInfo() {
    const { data } = await userServer.checkWxBind('job')

    store.setBindInfo(data)
    return !!data?.notification_activated
  }

  const showCode = async () => {
    await showJobQRCodeModal({
      descriptionNode: <span>我知道了</span>,
      fetchQRCodeFunc: async () => {
        return await userServer.getWxCode('job')
      },
      validConfig: {
        getCheckValue: fetchBindInfo,
      },
      afterOk: () => {
        localStorage.setItem('showQrCode', `${moment().format('YMD')}`)
      },
      hideOk: true,
    })
  }

  const cancelCode = async () => {
    await Modal.show({
      title: '解绑提示',
      width: 600,
      content: '解绑后您将无法及时收到作业的状态通知，确定要解除绑定吗？',
      onCancel: close => {
        close()
      },
      onOk: async close => {
        await userServer.unbindWx('job', store.bindInfo.wechat_openid)
        message.success('解绑成功')
        await fetchBindInfo()
        close()
      },
    })
  }

  return (
    <StyledItem>
      <label>微信通知：</label>
      <span className='text'>
        {!!store.bindInfo.notification_activated
          ? store.bindInfo.wechat_nickname
          : '暂未绑定'}
      </span>
      <div className='right'>
        {!!store.bindInfo.notification_activated ? (
          <span className='edit' onClick={cancelCode}>
            点击解绑
          </span>
        ) : (
          <span className='edit' onClick={showCode}>
            点击绑定
          </span>
        )}
      </div>
    </StyledItem>
  )
})
