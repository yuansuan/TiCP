/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Button, Icon } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { message, Spin } from 'antd'
import { Props } from './index'

const textBox = require('@/assets/images/textBox.svg')
const phone = require('@/assets/images/mobile.svg')
const StyledDiv = styled.div<{ textBox: string; phone: string }>`
  > .modal-content {
    display: flex;
    flex-direction: row;
    > .modal-left-content {
      padding: 20px 0 0 28px;
      flex: 1;
      div {
        margin-top: 24px;
      }
      > .title {
        font-size: 14px;
      }
      > .tip {
        width: 218px;
        height: 53px;
        background: url(${textBox}) no-repeat;
        font-size: 20px;
        font-weight: bold;
        line-height: 48px;
        text-align: center;
      }
      > .subtitle {
        color: #999999;
        font-size: 10px;
      }
    }
    > .modal-right-content {
      position: relative;
      flex: 1;
      display: flex;
      flex-direction: column;
      background: url(${phone}) no-repeat;
      width: 278px;
      height: 281px;
      justify-content: center;
      align-items: center;
      > .title {
        font-size: 18px;
        font-weight: bold;
        margin-top: 35px;
      }
      > .code {
        margin-top: 18px;
        width: 155px;
        height: 155px;
        position: relative;
        img {
          width: 155px;
          height: 155px;
          position: absolute;
        }
        .mask {
          width: 155px;
          height: 155px;
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
      > .footer {
        position: absolute;
        display: flex;
        flex-direction: column;
        align-items: center;
        bottom: -80px;
      }
    }
  }
`

export const JobQrCodeModal = observer(function JobQrCodeModal({
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
    <StyledDiv textBox={textBox} phone={phone}>
      <div className='modal-content'>
        <div className='modal-left-content'>
          <div className='title'>您可以接收到</div>
          <div className='tip'>作业计算完成通知</div>
          <div className='tip'>作业计算异常通知</div>
          <div className='subtitle'>
            您还可以通过平台“个人设置”模块选择“微信通知”绑定
          </div>
        </div>
        <div className='modal-right-content'>
          <div className='title'>开启微信实时通知</div>
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
          <div className='footer'>
            <Button
              type='primary'
              style={{ width: '88px' }}
              onClick={() => {
                onCancel()
                afterCancel && afterCancel()
              }}>
              {descriptionNode || <span>我已绑定</span>}
            </Button>
            {hideOk ? null : (
              <a
                style={{
                  textDecoration: 'none',
                  color: '#999999',
                  border: 'none',
                  fontSize: '12px',
                  marginTop: '10px',
                }}
                onClick={() => {
                  onOk()
                  afterOk && afterOk()
                }}>
                本次登录不再提示
              </a>
            )}
          </div>
        </div>
      </div>
    </StyledDiv>
  )
})
