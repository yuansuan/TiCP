import React from 'react'
import styled from 'styled-components'
import { Modal, Button } from '@/components'
import { Http, Validator } from '@/utils'
import { observer, useLocalStore } from 'mobx-react-lite'
import { message, Form, Select, Switch, Input, Upload } from 'antd'
import { Software } from '@/domain/VIsIBV/Software'
import { UploadOutlined } from '@ant-design/icons'
import { OPERATING_SYSTEM_PLATFORM_OF_ADD } from '@/domain/VIsIBV'

const colProps = { labelCol: { span: 5 }, wrapperCol: { span: 16 } }
const StyledLayout = styled.div`
  padding: 20px;
  .iconImage {
    width: 92px;
    height: 46px;
  }
  .footer {
    position: absolute;
    bottom: 20px;
    right: 10px;
    padding: 10px;
  }
`

type Props = {
  softwareItem?: Software
  onCancel?: () => void
  onOk?: () => void
}

function getBase64(img, callback) {
  const reader = new FileReader()
  reader.addEventListener('loadend', () => callback(reader.result))
  reader.readAsDataURL(img)
}
export default observer(function SoftwareEditor({
  softwareItem,
  onCancel,
  onOk
}: Props) {
  const [form] = Form.useForm()
  const state = useLocalStore(() => ({
    loading: false,
    setLoading(loading) {
      this.loading = loading
    },
    imgSrc: softwareItem?.icon,
    setImgSrc(imgSrc) {
      this.imgSrc = imgSrc
    }
  }))

  const handleDoubleClick = () => {
    const init_script = form.getFieldValue('init_script')
    const desc = form.getFieldValue('desc')
    if (init_script) form.setFieldsValue({ init_script: init_script })
    if (desc) form.setFieldsValue({ desc: desc })
  }

  async function onFinish(values) {
    try {
      state.setLoading(true)
      softwareItem
        ? await Http.put('/vis/software', {
            ...values,
            // gpu_desired: true, //默认就是开启GPU的
            id: softwareItem.id,
            icon: state.imgSrc
          })
        : await Http.post('/vis/software', {
            ...values,
            icon: state.imgSrc
          })
      onOk()
      message.success(`镜像${softwareItem ? '编辑' : '添加'}成功`)
    } finally {
      state.setLoading(false)
    }
  }

  const uploadIcon = ({ file }: any) => {
    if (file.type.match('image.*')) {
      getBase64(file, state.setImgSrc)
    }
  }

  const beforeUpload = file => {
    const isImage = ['image/jpg', 'image/jpeg', 'image/png'].includes(file.type)
    if (!isImage) {
      message.error('只能上传图片文件')
    }
    const sizeLimit = file.size / 1024 / 1024 < 2
    if (!sizeLimit) {
      message.error('只能上传小于 2MB 的图片')
    }
    return isImage && sizeLimit
  }

  return (
    <StyledLayout>
      <Form
        form={form}
        onFinish={onFinish}
        {...colProps}
        initialValues={{ ...softwareItem, gpu_desired: true }}>
        <Form.Item
          label='镜像名称'
          name='name'
          rules={[
            {
              required: true,
              validator: (_, value) =>
                Validator.validateInput(_, value, '镜像名称', true)
            }
          ]}>
          <Input />
        </Form.Item>
        <Form.Item label='GPU' name='gpu_desired' valuePropName='checked'>
          <Switch />
        </Form.Item>

        <Form.Item
          label='操作平台'
          name='platform'
          rules={[{ required: true, message: '操作平台不能为空' }]}>
          <Select>
            {Object.entries(OPERATING_SYSTEM_PLATFORM_OF_ADD).map(
              ([id, name]) => (
                <Select.Option key={name} value={name}>
                  {name}
                </Select.Option>
              )
            )}
          </Select>
        </Form.Item>
        <Form.Item
          label='镜像ID'
          name='image_id'
          rules={[
            {
              required: true,
              validator: (_, value) =>
                Validator.validateInput(_, value, '镜像ID', true)
            }
          ]}>
          <Input />
        </Form.Item>

        <Form.Item
          label='镜像描述'
          name='desc'
          rules={[
            {
              validator: (_, value) =>
                Validator.validateDesc(_, value, '镜像描述', false)
            }
          ]}>
          <Input.TextArea rows={4} onDoubleClick={handleDoubleClick} />
        </Form.Item>
        <Form.Item
          label='脚本内容'
          name='init_script'
          rules={[
            {
              validator: (_, value) =>
                Validator.validateScriptText(_, value, '脚本内容', false)
            }
          ]}>
          <Input.TextArea rows={6} onDoubleClick={handleDoubleClick} />
        </Form.Item>

        <Form.Item label='图标' tooltip='仅支持jpe、jpeg、png图片'>
          <Upload
            listType='picture-card'
            showUploadList={false}
            beforeUpload={beforeUpload}
            customRequest={uploadIcon}>
            <div>
              {state.imgSrc ? (
                <img className='iconImage' src={state.imgSrc} />
              ) : null}
            </div>
          </Upload>
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
