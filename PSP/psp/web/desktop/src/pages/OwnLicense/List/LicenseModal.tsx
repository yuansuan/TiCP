/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Form, Input } from 'antd'
import { Modal, Button, Mask } from '@/components'
import { Http } from '@/utils'
import { useStore } from '../store'
import { License } from '../store/Model/license'

const StyledLayout = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;

  > .main {
    padding: 16px 40px;
    flex: 1;
  }

  .ant-form-item-label {
    text-align: left;
  }

  .ant-form-item-control-input {
    width: 224px;
  }

  > .footer {
    padding: 10px;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};
  }
`

const layout = {
  labelCol: { span: 5 },
  wrapperCol: { span: 16 },
}

interface Props {
  merchandiseId: string
  onOk: () => void
  onCancel: () => void
}
export const LicenseModal = observer(function LicenseModal({
  merchandiseId,
  onOk,
  onCancel,
}: Props) {
  const { useForm } = Form
  const [form] = useForm()
  const store = useStore()
  const state = useLocalStore(() => ({
    loading: false,
    license: new License(),
    setLoading(status) {
      this.loading = status
    },
    updateloading: false,
    setUpdateLoading(status) {
      this.updateloading = status
    },
  }))
  const { license } = state

  useEffect(() => {
    state.setLoading(true)
    license
      .fetch(merchandiseId)
      .then(() => {
        form.setFieldsValue({
          ip: license.ip,
          license_port: license.license_port,
          extra_port: license.extra_port,
          provider_port: license.provider_port,
        })
      })
      .finally(() => {
        state.setLoading(false)
      })
  }, [])

  async function updateLicense(data) {
    try {
      state.setUpdateLoading(true)
      await Http.put('/ownLicense/license', {
        license: data,
        license_id: license.id,
        merchandise_id: merchandiseId,
      })
    } finally {
      state.setUpdateLoading(false)
    }
  }

  const onFinish = async values => {
    await updateLicense(values)
    store.fetch(store.params)
    onOk()
  }
  return (
    <StyledLayout>
      <div className='main'>
        {state.loading && <Mask.Spin />}
        <Form form={form} {...layout} onFinish={onFinish}>
          <Form.Item
            label='IP'
            name='ip'
            rules={[{ required: true, message: '请输入IP地址！' }]}>
            <Input />
          </Form.Item>

          <Form.Item
            label='许可证端口'
            name='license_port'
            rules={[{ required: true, message: '请输入许可证端口' }]}>
            <Input />
          </Form.Item>

          <Form.Item
            label='额外端口'
            name='extra_port'
            rules={[{ required: true, message: '请输入额外端口' }]}>
            <Input />
          </Form.Item>

          <Form.Item
            label='供应商端口'
            name='provider_port'
            rules={[{ required: true, message: '请输入供应商端口' }]}>
            <Input />
          </Form.Item>
        </Form>
      </div>

      <Modal.Footer
        className='footer'
        OkButton={() => (
          <Button
            disabled={state.updateloading}
            type='primary'
            onClick={() => {
              form.submit()
            }}>
            确认
          </Button>
        )}
        onCancel={onCancel}
      />
    </StyledLayout>
  )
})
