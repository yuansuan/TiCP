import React from 'react'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Form, Input } from 'antd'
import { Button, Modal } from '@/components'
import styled from 'styled-components'
import { vis } from '@/domain'

const StyledLayout = styled.div`
  padding: 20px;

  .footer {
    padding: 10px;
  }
`
interface IProps {
  onOk: () => void
  onCancel: () => void
  rowData: any
}

export const CloseSessionForm = observer(function CloseSessionForm({
  rowData,
  onOk,
  onCancel
}: IProps) {
  const [form] = Form.useForm()
  const state = useLocalStore(() => ({
    loading: false,
    setLoadig(loading) {
      this.loading = loading
    }
  }))

  async function confirm(values) {
    try {
      state.setLoadig(true)
      // await vis.closeSession({session_id: rowData?.id, reason: values['reason']})
      await vis.closeSession(rowData?.id, values['exit_reason'])
      onOk()
    } finally {
      state.setLoadig(false)
    }
  }

  return (
    <StyledLayout>
      <Form form={form} onFinish={confirm}>
        <Form.Item
          label='关闭原因'
          name='exit_reason'
          rules={[
            { type: 'string', max: 50 },
            { required: true, message: '关闭原因不能为空' }
          ]}>
          <Input.TextArea
            style={{ width: 450 }}
            rows={4}
            maxLength={50}
            showCount={{
              formatter: ({ count, maxLength }) =>
                `还可输入${maxLength - count}字`
            }}
          />
        </Form.Item>
        <Modal.Footer
          className='footer'
          onCancel={onCancel}
          OkButton={
            <Button
              type='primary'
              loading={state.loading}
              onClick={form.submit}>
              确认
            </Button>
          }
        />
      </Form>
    </StyledLayout>
  )
})
