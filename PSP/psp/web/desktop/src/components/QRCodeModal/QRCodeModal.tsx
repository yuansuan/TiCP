/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Modal, Button, Icon } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { message, Spin } from 'antd'
import { Props } from './index'
const StyledDiv = styled.div`
  display: flex;
  flex-direction: column;

  .modal-top-content {
    margin-bottom: 12px;
  }

  div {
    display: flex;
    flex-direction: column;
    font-size: 14px;
    span {
      font-size: 14px;
      line-height: 22px;
      .bold {
        font-weight: bold;
      }
    }
  }

  .code {
    width: 240px;
    height: 240px;
    position: relative;
    margin: auto;
    img {
      width: 240px;
      height: 240px;
      position: absolute;
    }
    .mask {
      width: 240px;
      height: 240px;
      position: absolute;
      background-color: #000000;
      opacity: 0.8;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 18px;
      color: #ffffff;
      flex-direction: row;
      cursor: pointer;
      .refresh {
        width: 24px;
        height: 24px;
        color: #ffffff;
      }
    }
  }

  .main {
    margin: 20px auto;
  }
`

export const QRCodeModal = observer(function QRCodeModal({
  onOk,
  onCancel,
  descriptionNode,
  fetchQRCodeFunc,
  validConfig: { getCheckValue, pollingInterval, afterCheckedFunc },
  hideOk = false,
  afterOk,
  afterCancel,
}: Props & { onOk: () => void; onCancel: () => void }) {
  const store = useLocalStore(() => ({
    imgWeChatUrl: null,
    setImgWeChatUrl(url) {
      this.imgWeChatUrl = url
    },
    expireTime: null,
    setExpireTime(count) {
      this.expireTime = count * 1000
    },
    refresh: false,
    setRefresh(status) {
      this.refresh = status
    },
    countdown: null,
    codePolling: null,
  }))

  const clearTimers = () => {
    clearInterval(store.codePolling)
    clearTimeout(store.countdown)
  }

  const fetchQrCode = async () => {
    const {
      data: { qrcodeUrl, expireSeconds },
    } = await fetchQRCodeFunc()
    store.setImgWeChatUrl(qrcodeUrl)
    store.setExpireTime(expireSeconds || 3000)

    store.countdown = setTimeout(() => {
      store.setRefresh(true)
      clearTimeout(store.countdown)
    }, store.expireTime)
  }

  const pollingValid = async () => {
    store.codePolling = setInterval(async () => {
      // 非falsy则提示绑定成功 https://developer.mozilla.org/zh-CN/docs/Glossary/Falsy
      if (getCheckValue && !!(await getCheckValue())) {
        clearInterval(store.codePolling)
        message.success('绑定成功')
        onOk()
        afterCheckedFunc && (await afterCheckedFunc())
      }
    }, pollingInterval || 1000)
  }

  useEffect(() => {
    fetchQrCode()
    pollingValid()

    return () => {
      clearTimers()
    }
  }, [])

  return (
    <StyledDiv>
      <div className='modal-top-content'>
        {descriptionNode || (
          <>
            <span>
              使用微信扫描以下二维码，关注<span className='bold'>“远算云”</span>
              公众号。
            </span>
            <span>
              通过关注公众号，您能够及时收到作业的状态通知，方便您更好的使用远算云平台。
            </span>
            <span>您也可以通过【个人设置-作业通知微信号】绑定您的微信号</span>
          </>
        )}
      </div>
      <div className='code'>
        {store.imgWeChatUrl ? (
          <img src={store.imgWeChatUrl} alt='二维码' />
        ) : (
          <div
            style={{
              width: 240,
              height: 240,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              background: '#fafafa',
            }}>
            <Spin size='large' />
          </div>
        )}
        {store.refresh && (
          <div
            className='mask'
            onClick={() => {
              clearTimeout(store.countdown)
              fetchQrCode().then(() => store.setRefresh(false))
            }}>
            <Icon className='refresh' type='revert' />
            点击刷新
          </div>
        )}
      </div>

      <Modal.Footer
        className='footer'
        CancelButton={
          <Button
            type='primary'
            onClick={() => {
              onCancel()
              afterCancel && afterCancel()
            }}>
            我知道了
          </Button>
        }
        OkButton={
          hideOk ? null : (
            <Button
              style={{ color: '#999999', border: 'none' }}
              onClick={() => {
                onOk()
                afterOk && afterOk()
              }}>
              不再提示
            </Button>
          )
        }
      />
    </StyledDiv>
  )
})
