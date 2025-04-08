/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { Modal } from '@/components'
import { env, currentUser } from '@/domain'
import { Spin } from 'antd'
const reportBg = require('@/assets/images/report_bg.jpg')

const StyledLayout = styled.div`
  display: flex;
  flex-direction: column;

  .qr-modal-content {
    width: 700px;
    height: 400px;
    position: relative;
    display: flex;
    flex-direction: column;
    background: url(${reportBg}) no-repeat;
    background-size: contain;
    justify-content: center;
    align-items: center;
  }

  .code {
    width: 100%;
    height: 100%;
    margin: 0 auto;
    img {
      width: 180px;
      height: 180px;
      position: absolute;
      top: 110px;
      left: 66px;
    }
  }
`

export const WXQrModal = observer(() => {
  const [imgData, setImgData] = useState('')

  let count = Number(window.localStorage.getItem('show_yearly'))
  count += 1
  useEffect(() => {
    getQrImageUrl(currentUser.id)
    window.localStorage.setItem('show_yearly', count || '0')
  }, [])

  async function getQrImageUrl(userId) {
    await env.getWxQrCode(userId).then(data => {
      let imgUrl = 'data:image/png;base64,' + data.image
      setImgData(imgUrl)
    })
  }

  return (
    <StyledLayout>
      <div className='qr-modal-content'>
        <div className='code'>
          {imgData ? (
            <img src={imgData} alt='二维码' />
          ) : (
            <div
              style={{
                width: 180,
                height: 180,
                position: 'absolute',
                top: 110,
                left: 66,
                background: '#fafafa',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
              }}>
              <Spin size='large' />
            </div>
          )}
        </div>
      </div>
    </StyledLayout>
  )
})

export const showWXQrModal = () => {
  return Modal.show({
    title: '年度工作报告',
    footer: null,
    width: 700,
    bodyStyle: {
      height: 400,
      padding: 0
    },
    onCancel: next => {
      next()
    },
    content: <WXQrModal />
  })
}
