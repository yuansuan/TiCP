import React, { useState, useEffect } from 'react'
import { observer } from 'mobx-react'
import { Descriptions, Button, Transfer } from 'antd'
import { Label } from '@/components'
import { Http } from '@/utils'

import styled from 'styled-components'

export const FormWrapper = styled.div`
  padding: 10px;

  .item {
    display: flex;
    flex-direction: column;
  }

  .formItem {
    width: 400px;
  }

  .ant-descriptions-item {
  }

  .footer {
    position: absolute;
    display: flex;
    bottom: 0px;
    right: 0;
    width: 100%;
    line-height: 64px;
    height: 64px;
    background: white;

    .footerMain {
      margin-left: auto;
      margin-right: 8px;

      button {
        margin: 0 8px;
      }
    }
  }
`

export const Tips = styled.span`
  font-family: PingFangSC-Regular;
  font-size: 12px;
  color: #999999;
  line-height: 22px;
`

export const BindForm = observer(props => {
  const oldBindWSIds = props.data['WS_list'].map(ws => ws.id.toString()) || []
  const [adding, setAdding] = useState(false)
  const [WSList, setWSList] = useState([])
  const [targetWSKeys, setTargetWSKeys] = useState(oldBindWSIds)
  const [selectedWSKeys, setSelectedWSKeys] = useState([])


  const submit = async () => {
    try {
      setAdding(true)
      await props.onOk({
        oldBindWSIds,
        newBindWSIds: targetWSKeys,
        id: props.data['id'],
      })
    } finally {
      setAdding(false)
    }
  }

  useEffect(() => {
    // 获取workstation列表
    (async () => {
      const res = await Http.get('/visual/workstationList', {baseURL: './'})
      setWSList(res.data.map(ws => ({key: ws.id.toString(), name: ws.name})))
    })()
  }, [])

  const handleChange = (nextTargetKeys, direction, moveKeys) => {
    setTargetWSKeys(nextTargetKeys)
  }

  const handleSelectChange = (sourceSelectedKeys, targetSelectedKeys) => {
    setSelectedWSKeys([...sourceSelectedKeys, ...targetSelectedKeys])
  }

  return (
    <FormWrapper>
      <Descriptions title='' column={1} style={{ margin: '0 0 50px 0' }}>
        <Descriptions.Item label={<Label required>软件名称</Label>}>
          <div className='item'>
            {props.data['name']}
          </div>
        </Descriptions.Item>
        <Descriptions.Item label={<Label required>工作站</Label>}>
          <Transfer
            dataSource={WSList}
            titles={['未关联工作站', '已关联工作站']}
            targetKeys={targetWSKeys}
            selectedKeys={selectedWSKeys}
            onChange={handleChange}
            onSelectChange={handleSelectChange}
            render={item => item.name}
          />
        </Descriptions.Item>
      </Descriptions>
      <div className='footer'>
        <div className='footerMain'>
          <Button type='primary' loading={adding} onClick={submit}>
            确认
          </Button>
          <Button onClick={props.onCancel}>取消</Button>
        </div>
      </div>
    </FormWrapper>
  )
})
