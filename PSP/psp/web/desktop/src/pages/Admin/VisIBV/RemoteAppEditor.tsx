import React, { useState } from 'react'
import styled from 'styled-components'
import { Form, Button, Space, Input, message, Switch, Tooltip } from 'antd'
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons'
import { Http } from '@/utils'
import { Modal } from '@/components'

const { TextArea } = Input

const Wrapper = styled.div`
  padding: 20px 5px 5px 20px;
  overflow-y: hidden;

  .form {
    width: 100%;
    height: 380px;
    overflow-y: auto;
  }

  .footer {
    position: absolute;
    right: 10px;
    bottom: 0px;
    .btn {
      margin: 10px;
    }
  }
`

type RemoteApp = {
  software_id: string
  desc: string
  base_url: string
  name: string
  dir: string
  args: string
  logo: string
  disable_gfx: boolean
}

interface IProps {
  remoteAppList?: RemoteApp[]
  onOk: () => void
  onCancel: () => void
  softwareId: string
}

const initAppData = {
  base_url: null,
  name: null,
  dir: null,
  args: null,
  logo: null,
  disable_gfx: true
}

const colProps = { labelCol: { span: 5 } }
const formItemStyle = { width: 540 }

export function RemoteAppEditor({
  softwareId,
  remoteAppList,
  onOk,
  onCancel
}: IProps) {
  const [submiting, setSubmiting] = useState(false)

  const [form] = Form.useForm(null)

  if (!remoteAppList) {
    form.resetFields()
  }

  const addApp = body => {
    return Http.post('/vis/software/remote/app', {
      ...body,
      software_id: softwareId
    })
  }

  const updateApp = (remoteAppId, body) => {
    return Http.put(
      `/vis/software/remote/app/${remoteAppId}`,
      {
        ...body,
        software_id: softwareId
      },
      {}
    )
  }

  const deleteApp = remoteAppId => {
    return Http.delete(`/vis/software/remote/app/${remoteAppId}`, {})
  }

  async function onFinish(values) {
    let appList = form.getFieldValue('app_list')

    try {
      if (remoteAppList.length !== 0 && appList.length === 0) {
        await Modal.showConfirm({
          title: '确认',
          content: '确认删除所有配置的应用吗？'
        })
      }

      if (remoteAppList.length === 0 && appList.length === 0) {
        onOk()
        return
      }

      setSubmiting(true)
      message.info('远程应用配置中，请稍后......')

      await Promise.all(
        appList.map(async app => {
          if (app.id) {
            // update
            return await updateApp(app.id, app)
          } else {
            // add
            return await addApp(app)
          }
        })
      )

      const appIdSet = new Set(appList.map(app => app.id))

      const delAppList = remoteAppList.filter(item => !appIdSet.has(item.id))

      await Promise.all(
        delAppList.map(async app => {
          return await deleteApp(app.id)
        })
      )

      message.success('远程应用配置成功')
      onOk()
    } finally {
      setSubmiting(false)
    }
  }

  return (
    <Wrapper>
      <div className='form'>
        <Form
          name='remote_app_list_form'
          form={form}
          {...colProps}
          initialValues={{
            app_list: remoteAppList.length !== 0 ? remoteAppList : [initAppData]
          }}
          onFinish={onFinish}>
          <Form.List name='app_list'>
            {(fields, { add, remove }) => (
              <>
                {fields.length > 0 ? (
                  fields.map(({ key, name, ...restField }) => (
                    <Space
                      key={key}
                      style={{
                        display: 'flex',
                        flexDirection: 'column',
                        marginBottom: 8
                      }}
                      align='baseline'>
                      <Space>{<h3>远程应用{Number(name) + 1}</h3>}</Space>
                      <Space
                        style={{ display: 'flex', alignItems: 'baseline' }}>
                        <Form.Item
                          {...restField}
                          label='应用名称'
                          style={formItemStyle}
                          name={[name, 'name']}
                          rules={[{ required: true, message: '输入应用名称' }]}>
                          <Input maxLength={64} placeholder='应用名称' />
                        </Form.Item>
                        <Button
                          type='ghost'
                          onClick={() => remove(name)}
                          icon={<MinusCircleOutlined />}></Button>
                      </Space>
                      <Space>
                        <Form.Item
                          {...restField}
                          label='应用参数'
                          style={formItemStyle}
                          name={[name, 'args']}>
                          <TextArea
                            maxLength={512}
                            style={{ height: 80, resize: 'none' }}
                            placeholder='应用参数'
                          />
                        </Form.Item>
                      </Space>
                      <Space>
                        <Form.Item
                          {...restField}
                          style={formItemStyle}
                          label='基础URL'
                          name={[name, 'base_url']}
                          rules={[{ required: true, message: '输入基础URL' }]}>
                          <Input maxLength={64} placeholder='基础URL' />
                        </Form.Item>
                      </Space>
                      <Space>
                        <Form.Item
                          {...restField}
                          label='仅图像传输'
                          name={[name, 'disable_gfx']}
                          style={formItemStyle}
                          valuePropName='checked'>
                          <Switch
                            checkedChildren='开启'
                            unCheckedChildren='关闭'
                          />
                        </Form.Item>
                      </Space>
                    </Space>
                  ))
                ) : (
                  <>无应用配置</>
                )}
                {/* <Form.Item>
                  <Tooltip
                    title={fields.length >= 1 && '目前仅支持单个应用配置'}>
                    <Button
                      type='dashed'
                      disabled={fields.length >= 1}
                      style={{ width: 540 }}
                      onClick={() => {
                        form
                          .validateFields()
                          .then(res => {
                            add(initAppData)
                          })
                          .catch(e => {
                            message.error('请确保已添加应用信息校验通过')
                          })
                      }}
                      block
                      icon={<PlusOutlined />}>
                      添加应用
                    </Button>
                  </Tooltip>
                </Form.Item> */}
              </>
            )}
          </Form.List>
        </Form>
      </div>
      <div className='footer'>
        <Button className='btn' onClick={() => onCancel()}>
          取消
        </Button>
        <Button
          loading={submiting}
          className='btn'
          type='primary'
          onClick={() => form.submit()}>
          配置
        </Button>
      </div>
    </Wrapper>
  )
}
