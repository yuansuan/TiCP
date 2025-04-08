/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useEffect } from 'react'
import styled from 'styled-components'
import { Button, Icon, Modal } from '@/components'
import { Http } from '@/utils'
import { message, Tooltip } from 'antd'
import BindForm from './BindForm'
import sysConfig from '@/domain/SysConfig'
import { QuestionCircleOutlined } from '@ant-design/icons'
import { eventEmitter, IEventData } from '@/utils'

const ysCloud = require('@/assets/images/logo.svg')

export const CLOUD_PLATFORM = [
  { label: '泛超算云', img: ysCloud, color: '#0144e7' }
]

export const CLOUD_COLORS = {
  泛超算云: '#0144e7'
}

const BurstWrapper = styled.div`
  margin-left: -20px;
  display: flex;
  flex-direction: column;

  .tips {
    border-radius: 3px;
    background: #e3f4ff;
    border: 1px solid #0090fa;
    color: #0090fa;
    width: 800px;
    padding: 10px;
    margin: 0 20px;
  }

  .cards {
    display: flex;
  }
`

interface ICloudCardWrapper {
  checked: boolean
}

const CloudCardWrapper = styled.div<ICloudCardWrapper>`
  display: flex;
  width: 300px;
  margin: 20px;
  border-radius: 10px;
  border: 1px solid ${props => (props.checked ? '#3490ff' : '#bcbcbc')};

  &:hover {
    border: 1px solid #3490ff;
  }

  .logo {
    width: 45%;
    border-radius: 10px 0px 0px 10px;
    height: 100px;
    background: ${props => (props.checked ? '#e3f4ff' : '#f4f4f4')};
    display: flex;
    padding: 5px;
    align-items: center;
  }

  .content {
    position: relative;
    display: flex;
    padding: 5px 0;
    align-items: center;
    flex-direction: column;
    justify-content: space-around;
    width: 55%;

    .cardName {
      font-size: 16px;

      .anticon {
        margin-left: 2px;
      }
    }

    .msg {
      font-size: 8px;
    }

    .checkBox {
      position: absolute;
      right: 10px;
      top: 5px;
      font-size: 18px;
    }

    .checked {
      color: #3490ff;
    }
  }
`

const SyncContentWrapper = styled.div`
  display: inline;
`

const SyncContentTitleWrapper = styled.div`
  .sync-btn-div {
    margin-top: 10px;
    display: flex;
    justify-content: flex-end;
  }
`

function SyncContent() {
  const [hour, setHour] = useState(0)

  useEffect(() => {
    async function fetch() {
      const { data } = await Http.get('/bindcloud/frequency')
      setHour(data.frequency)
    }

    fetch()
  }, [])

  const syncApps = async () => {
    message.warn('泛超算云应用同步中，请稍候')
    await Http.put('/bindcloud/sync-apps')
    message.success('泛超算云应用同步成功')
  }

  return (
    <SyncContentWrapper>
      <Tooltip
        title={
          <SyncContentTitleWrapper>
            <div>当前云端应用每 {hour} 小时同步一次</div>
            <div className='sync-btn-div'>
              <Button onClick={syncApps}>立即同步</Button>
            </div>
          </SyncContentTitleWrapper>
        }>
        <QuestionCircleOutlined />
      </Tooltip>
    </SyncContentWrapper>
  )
}

function CloudCard({
  cardLogo,
  cardName,
  cardMsg,
  cardValue,
  cardStatus,
  onChange,
  onBind,
  onUnBind,
  loading
}) {
  return (
    <CloudCardWrapper
      checked={cardValue === cardName}
      onClick={() => onChange(cardName)}>
      <div className='logo'>
        <img src={cardLogo} alt={cardName} width='40px' />
        <span
          style={{
            paddingLeft: 5,
            fontSize: 18,
            fontWeight: 'bold',
            color: CLOUD_COLORS[cardName]
          }}>
          {cardName}
        </span>
      </div>
      <div className='content'>
        <div className='cardName'>
          {cardName}
          {cardStatus && <SyncContent />}
          <Icon
            type='check-filled'
            className='checkBox'
            style={{ color: cardValue === cardName ? '#3490ff' : '#fff' }}
          />
        </div>
        <div className='msg'>
          状态: <span data-nomark>{loading ? '获取中...' : cardMsg}</span>
        </div>
        <div className='btns'>
          {cardStatus === true ? (
            <Button onClick={onUnBind}>解绑</Button>
          ) : (
            <Button onClick={onBind}>绑定</Button>
          )}
        </div>
      </div>
    </CloudCardWrapper>
  )
}

export default function BurstConfig() {
  let [cloud, setCloud] = useState('泛超算云')
  let [bindStatus, setBindStatus] = useState(false)
  let [msg, setMsg] = useState('')
  let [loading, setLoading] = useState(false)

  const getBindStatus = async () => {
    try {
      setLoading(true)
      const res = await Http.get('/bindcloud/status', { timeout: 0 })
      const { isBind, msg } = res.data
      setBindStatus(isBind)
      setMsg(msg)
      return res
    } finally {
      setLoading(false)
    }
  }

  const onBind = () => {
    if (sysConfig.firewallConfig.level === 'none') {
      message.error(
        '防火墙策略为完全限制，对外不能做任何访问，请修改防火墙策略后，再执行绑定'
      )
      return
    }

    Modal.show({
      title: '绑定泛超算云',
      closable: false,
      footer: null,
      content: ({ onCancel, onOk }) => {
        const ok = async () => {
          try {
            await getBindStatus()
            // 绑定后，通知刷新账户
            eventEmitter.emit(`REFRESH_ACCOUNT`, {} as IEventData)
          } finally {
            onOk()
          }
        }
        return <BindForm onCancel={onCancel} onOk={ok} />
      },
      width: 600
    })
  }

  const onUnBind = async () => {
    if (sysConfig.firewallConfig.level === 'none') {
      message.error(
        '防火墙策略为完全限制，对外不能做任何访问，请修改防火墙策略后，再执行解绑。'
      )
      return
    }

    await Modal.showConfirm({
      content: `确定解绑泛超算云？`
    })

    Http.delete('/bindcloud/unbind').then(async res => {
      if (res.data.success) {
        message.info('解绑泛超算云成功')
        await getBindStatus()
        // 绑定后，通知刷新账户
        eventEmitter.emit(`REFRESH_ACCOUNT`, {} as IEventData)
      } else {
        message.error('解绑泛超算云失败')
      }
    })
  }

  // useEffect(() => {
  //   getBindStatus()
  // }, [])

  return (
    <>
      <BurstWrapper>
        <div className='tips'>
          通过云服务, 您可以在本地算力不足的情况下, 将作业调度上云,
          让作业低成本，高效，完全的运行。
        </div>
        <div className='cards'>
          {CLOUD_PLATFORM.map(c => (
            <CloudCard
              key={c.label}
              loading={loading}
              onChange={value => setCloud(value)}
              cardLogo={c.img}
              cardName={c.label}
              cardValue={cloud}
              cardMsg={msg}
              cardStatus={bindStatus}
              onBind={onBind}
              onUnBind={onUnBind}
            />
          ))}
        </div>
      </BurstWrapper>
    </>
  )
}
