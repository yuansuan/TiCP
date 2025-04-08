import React, { useState, useMemo, useEffect, useCallback } from 'react'
import { Button, message, Switch, Tooltip } from 'antd'
import Label from '@/components/Label'
import { Icon, Table } from '@/components'
import { ConfigWrapper } from './style'
import { Modal } from '@/components'
import AddWLForm from './AddWLForm'
import styled from 'styled-components'
import whiteList from '@/domain/WhiteList'
import { Http } from '@/utils'

const TooltipWrapper = styled.div`
  > p {
    margin-bottom: 0;
  }
`
const SwitchStyle = {
  padding: '0 5px',
  marginLeft: 30,
}

const TableLinkBtn = styled(Button)`
  padding: 0;
`
export default function WhitelistConfig() {
  const [loading, setLoading] = useState(false)
  const [ruleLoading, setRuleLoading] = useState(false)
  const [open, setOpen] = useState(false)

  const toggleSwitch = useCallback(async value => {
    try {
      setRuleLoading(true)
      await Http.put('/sysconfig/ruleconfig', { enable: value })
      setOpen(value)
      message.success(`白名单规则${value ? '已开启' : '已关闭'}`)
    } finally {
      setRuleLoading(false)
    }
  }, [])

  const columns = useMemo(() => {
    return [
      {
        props: {
          flexGrow: 1,
        },
        header: 'IP地址',
        dataKey: 'ip',
      },
      {
        props: {
          flexGrow: 1,
        },
        header: '用户名',
        dataKey: 'username',
      },
      {
        props: {
          flexGrow: 1,
        },
        header: '创建时间',
        dataKey: 'time',
      },
      {
        header: '操作',
        props: {
          width: 150,
        },
        cell: {
          props: {
            dataKey: 'operator',
          },
          render: ({ rowData }) => {
            const id = rowData.id

            return (
              <>
                <TableLinkBtn
                  type='link'
                  disabled={!open}
                  onClick={() => deleteRule(id)}>
                  删除
                </TableLinkBtn>
              </>
            )
          },
        },
      },
    ]
  }, [open])

  const addRule = () => {
    Modal.show({
      title: '添加规则',
      closable: false,
      footer: null,
      content: ({ onCancel, onOk }) => {
        const ok = async () => {
          onOk()
          await getList()
        }
        return <AddWLForm onCancel={onCancel} onOk={ok} />
      },
      width: 600,
    })
  }

  const deleteRule = async (id: string) => {
    await Modal.showConfirm({
      content: `确定删除这条规则吗？`,
    })

    const res = await Http.delete('/sysconfig/whitelist', { params: { id } })

    if (res.success) {
      message.success('删除规则成功')
      await getList()
    } else {
      message.error('删除规则失败')
    }
  }

  const getList = async () => {
    try {
      setLoading(true)
      await whiteList.getWhileList()
    } finally {
      setLoading(false)
    }
  }

  const getDownloadRuleConfig = async () => {
    const res = await Http.get('sysconfig/ruleconfig')
    setOpen(res.data.download.enable)
  }

  useEffect(() => {
    Promise.all([getDownloadRuleConfig(), getList()])
  }, [])

  return (
    <ConfigWrapper>
      <div className='item'>
        <Label align={'left'}>开启白名单规则</Label>
        <Switch
          style={SwitchStyle}
          checked={open}
          disabled={ruleLoading}
          checkedChildren='关闭'
          unCheckedChildren='开启'
          onChange={toggleSwitch}
        />
        <Tooltip
          className='tooltip'
          placement='right'
          title={
            <TooltipWrapper>
              设置规则如下：
              <p>1. 设置IP地址未设置用户名时，允许所有用户使用该IP下载。</p>
              <p>2. 设置用户名未设置IP地址时，允许该用户使用任意IP下载。</p>
              <p>3. 设置IP地址和用户名时，只允许该用户在指定IP下载。</p>
              <p>4. 未设置IP地址且未设置用户名，不允许任何IP和用户名下载。</p>
            </TooltipWrapper>
          }>
          <Icon style={{ padding: '0 5px' }} type={'help-circle'} />
        </Tooltip>
      </div>
      <div className='item btn'>
        <Button disabled={!open} onClick={addRule}>
          添加规则
        </Button>
      </div>
      <div className='item'>
        <Table
          columns={columns as any}
          props={{
            data: whiteList.list,
            width: 800,
            ...(whiteList.list.length <= 20
              ? { autoHeight: true }
              : { height: 54 * 21 }),
            rowKey: 'ip',
            loading: loading,
            locale: {
              emptyMessage: '没有数据',
              loading: '数据加载中...',
            },
          }}
        />
      </div>
    </ConfigWrapper>
  )
}
