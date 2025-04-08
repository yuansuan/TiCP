import React, { useState } from 'react'
import { Descriptions, Input, Button, message} from 'antd'
import { FormWrapper } from './style'
import { Label } from '@/components'

export function Form({ rowData, onOk, onCancel }) {
  const [reason, setReason] = useState('')

  const submit = () => {
    if (!reason.trim()) {
      message.error('拒绝原因不能为空')
      return 
    }
    onOk && onOk(reason)
  }

  return (
    <FormWrapper>
      <Descriptions title='' column={1} style={{ margin: '0 0 20px 0' }}>
      <Descriptions.Item label={<Label width={120}>申请人</Label>}>
          <div className='item'>
            {rowData.application_name}
          </div>
        </Descriptions.Item>
        <Descriptions.Item label={<Label width={120}>申请内容</Label>}>
          <div className='item'>
            {rowData.content}
          </div>
        </Descriptions.Item>
        <Descriptions.Item label={<Label required width={120}>拒绝原因</Label>}>
          <div className='item'>
            <Input.TextArea
              className='formItem'
              placeholder='请输入拒绝原因'
              maxLength={100}
              style={{ height: 164 }}
              value={reason}
              onChange={e => {
                setReason(e.target.value)
              }}
            />
          </div>
        </Descriptions.Item>
      </Descriptions>
      <div className='footer'>
        <div className='footerMain'>
          <Button type='primary' onClick={submit}>
            确认
          </Button>
          <Button onClick={onCancel}>取消</Button>
        </div>
      </div>
    </FormWrapper>
  )
}
