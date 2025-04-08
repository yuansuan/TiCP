import React, { useState, useEffect } from 'react'
import { Descriptions, Select, Button, message } from 'antd'
import { Label } from '@/components'
import styled from 'styled-components'
import sysConfig from '@/domain/SysConfig'

const Option = Select.Option

const FormWrapper = styled.div`
  padding: 10px;

  .item {
    display: flex;
    flex-direction: column;
    width: 350px;
  }

  .formItem {
    width: 300px;
  }

  .ant-descriptions-item {
    display: flex;
  }

  .ant-descriptions-item-label {
    padding-top: 5px;
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


export function Form({ onOk, onCancel }) {
  const [id, setId] = useState('')
  const [options, setOptions] = useState([])

  const submit = () => {
    if (!id) {
      message.error('审批人不能为空')
      return
    }
    const name = options.find(o => o.id === id)?.name
    onOk && onOk({id, name})
  }

  useEffect(() => {
    fetch('/api/v1/user/optionList?filterPerm=8')
      .then(response => response.json())  
      .then(res => {
        const opts = res?.data?.map(d => ({id: d.key, name: d.title})) || []
        setOptions(opts)
        setId(sysConfig?.threeMemberMgrConfig?.defaultApprover?.id || opts?.[0]?.id || '')
      })  
  }, [])

  return (
    <FormWrapper>
      <Descriptions title='' column={1} style={{ margin: '0 0 20px 0' }}>
        <Descriptions.Item label={<Label required width={120}>审批人</Label>}>
          <div className='item'>
            <Select style={{width: 200}} value={id} 
              onChange={(v) => { 
                setId(v)
              }}
              notFoundContent={<p>请联系超级管理员，创建审批人（安全管理员）</p>}
              >
                {
                  options.map(o => <Option value={o.id} key={o.id}>
                    {o.name}
                  </Option>)
                }
            </Select>
          </div>
        </Descriptions.Item>
      </Descriptions>
      <div className='footer'>
        <div className='footerMain'>
          <Button type='primary' onClick={submit}>
            发起操作申请
          </Button>
          <Button onClick={onCancel}>取消</Button>
        </div>
      </div>
    </FormWrapper>
  )
}
