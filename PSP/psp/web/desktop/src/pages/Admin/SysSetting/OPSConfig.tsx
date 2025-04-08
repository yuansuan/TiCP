import React, { useState, useEffect, useCallback } from 'react'
import styled from 'styled-components'
import { InputNumber, Switch, Tooltip, message } from 'antd'
import Label from '@/components/Label'
import { Icon } from '@/components'
import { ConfigWrapper } from './style'
import { Http } from '@/utils'
import debounce from 'lodash.debounce'
import sysConfig from '@/domain/SysConfig'

const Tips = styled.div`
  border-radius: 3px;
  background: #e3f4ff;
  border: 1px solid #0090fa;
  color: #0090fa;
  width: 800px;
  padding: 10px;
  margin: 10px 0px;
`

const itemStyle = {
  padding: '0 5px',
  marginRight: 30,
}

export default function OPSConfig() {
  const [port, setPort] = useState(null)
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)

  const toggleSwitch = useCallback(
    async value => {
      if (port) {
        try {
          setLoading(true)
          const res = await sysConfig.fetchFirewallConfig()

          if (res.data.level === 'all') {
            await Http.put('/sysconfig/optstatus', {
              status: value ? 'Open' : 'Close',
            })
            setOpen(value)
            message.success(`远程运维状态${value ? '已开启' : '已关闭'}`)
          } else {
            message.error('开启远程运维之前, 防火墙策略必须设置为不限制')
          }
        } finally {
          setLoading(false)
        }
      } else {
        message.error('开启远程运维之前, 远程运维端口不能为空')
      }
    },
    [port]
  )

  const updatePort = useCallback(async value => {
    try {
      setLoading(true)
      await Http.put('/sysconfig/optport', { port: value })
      setPort(value)
    } finally {
      setLoading(false)
    }
  }, [])

  const debounceUpdatePort = debounce(port => {
    updatePort(port)
  }, 600)

  useEffect(() => {
    const getOptInfo = async () => {
      const res = await Http.get('/sysconfig/optinfo')
      if (res.data?.id) {
        setOpen(res.data.status === 0 ? false : true)
        setPort(res.data.port)
      }
    }
    getOptInfo()
  }, [])

  return (
    <ConfigWrapper>
      <Tips>
        注意: 开启远程运维开关后，远程运维端口不能再编辑，防火墙策略不能再设置。
      </Tips>
      <div className='item'>
        <Label align={'left'}>远程运维端口</Label>
        <InputNumber
          size='small'
          disabled={open}
          min={1}
          max={9999}
          step={1}
          value={port}
          precision={0}
          parser={value => {
            if (isNaN(parseInt(value))) {
              return null
            } else {
              const v = parseInt(value)
              if (v > 9999) return 9999
              if (v < 1) return 1
              return v
            }
          }}
          onChange={debounceUpdatePort}
        />
        <Tooltip title='远程运维端口范围为 1 ~ 9999'>
          <Icon style={itemStyle} type={'help-circle'} />
        </Tooltip>
        <Switch
          style={itemStyle}
          checked={open}
          disabled={loading}
          checkedChildren='关闭'
          unCheckedChildren='开启'
          onChange={toggleSwitch}
        />
        <span>
          远程运维状态: <span data-nomark>{open ? `已开启` : `已关闭`}</span>
        </span>
      </div>
    </ConfigWrapper>
  )
}
