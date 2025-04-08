import React from 'react'
import styled from 'styled-components'
import { Modal, Button } from '@/components'
import { Http, Validator } from '@/utils'
import { observer, useLocalStore } from 'mobx-react-lite'
import { message, Form, Select, Input, InputNumber } from 'antd'
import { Hardware } from '@/domain/VIsIBV/Hardware'

const colProps = { labelCol: { span: 6 }, wrapperCol: { span: 16 } }
const StyledLayout = styled.div`
  padding: 20px;
`

type Props = {
  hardwareItem?: Hardware
  onCancel?: () => void
  onOk?: () => void
}

export default observer(function HardwareEditor({
  hardwareItem,
  onCancel,
  onOk
}: Props) {
  const [form] = Form.useForm()
  const state = useLocalStore(() => ({
    loading: false,
    setLoading(loading) {
      this.loading = loading
    }
  }))

  async function onFinish(values) {
    try {
      state.setLoading(true)
      hardwareItem
        ? await Http.put('/vis/hardware', {
            ...values,
            id: hardwareItem.id
          })
        : await Http.post(
            '/vis/hardware',
            {
              ...values,
              GpuModel: values.gpu_model,
              CpuModel: values.cpu_model
            },
            {}
          )
      onOk()
      message.success(`实例${hardwareItem ? '编辑' : '添加'}成功`)
    } finally {
      state.setLoading(false)
    }
  }

  return (
    <StyledLayout>
      <Form
        form={form}
        onFinish={onFinish}
        {...colProps}
        initialValues={{
          ...hardwareItem
        }}>
        <Form.Item
          label='实例名称'
          name='name'
          rules={[
            {
              required: true,
              validator: (_, value) =>
                Validator.validateInput(_, value, '实例名称', true)
            }
          ]}>
          <Input />
        </Form.Item>
        <Form.Item
          label='实例类型'
          name='instance_type'
          rules={[
            {
              required: true,
              validator: (_, value) =>
                Validator.validateInput(_, value, '实例类型', true)
            }
          ]}>
          <Input />
        </Form.Item>
        <Form.Item
          label='实例机型系列'
          name='instance_family'
          rules={[
            {
              required: true,
              validator: (_, value) =>
                Validator.validateInput(_, value, '实例机型', true)
            }
          ]}>
          <Input />
        </Form.Item>

        <Form.Item label='实例的最大带宽' required>
          <Form.Item
            name='network'
            noStyle
            rules={[{ required: true, message: '实例的最大带宽不能为空' }]}>
            <InputNumber
              style={{ width: 150 }}
              min={1}
              parser={text =>
                text && Math.round(Number(text.replace(/[^0-9.]+/g, '')))
              }
            />
          </Form.Item>{' '}
          Mbps
        </Form.Item>
        <Form.Item label='CPU核数' required>
          <Form.Item
            name='cpu'
            noStyle
            rules={[{ required: true, message: 'CPU核数不能为空' }]}>
            <InputNumber
              min={1}
              parser={text =>
                text && Math.round(Number(text.replace(/[^0-9.]+/g, '')))
              }
            />
          </Form.Item>{' '}
        </Form.Item>
        <Form.Item label='GPU数量' required>
          <Form.Item
            noStyle
            name='gpu'
            rules={[{ required: true, message: 'GPU数量不能为空' }]}>
            <InputNumber
              min={0}
              parser={text =>
                text && Math.round(Number(text.replace(/[^0-9.]+/g, '')))
              }
            />
          </Form.Item>{' '}
        </Form.Item>
        <Form.Item label='内存' required>
          <Form.Item
            noStyle
            name='mem'
            rules={[{ required: true, message: '内存不能为空' }]}>
            <InputNumber
              min={1}
              parser={text =>
                text && Math.round(Number(text.replace(/[^0-9.]+/g, '')))
              }
            />
          </Form.Item>{' '}
          GB
        </Form.Item>

        <Form.Item
          label='实例描述'
          name='desc'
          rules={[
            {
              validator: (_, value) =>
                Validator.validateDesc(_, value, '实例描述', false)
            }
          ]}>
          <Input.TextArea maxLength={256} rows={4} />
        </Form.Item>
        <Modal.Footer
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
