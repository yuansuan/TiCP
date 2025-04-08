/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import {
  Form,
  Input,
  DatePicker,
  Button,
  message,
  InputNumber,
  Switch,
  Row,
  Col,
  Select
} from 'antd'
import { observer } from 'mobx-react-lite'
import { FormListFieldData } from 'antd/es/form/FormList'
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons'
import moment from 'moment'
import { CollectorTypeLabel } from '@/domain/LicenseMgr/LicenseInfo'
import { Http } from '@/utils'

const { RangePicker } = DatePicker

export const StyledLayout = styled.div`
  padding: 20px;
  overflow-y: hidden;

  .form {
    width: 660px;
  }

  .footer {
    display: flex;
    justify-content: flex-end;
    padding: 20px;
    .btn {
      margin: 10px;
    }
  }
`

const colProps = { labelCol: { span: 6 }, wrapperCol: { span: 14, offset: 1 } }

const formItemLayoutWithOutLabel = {
  labelCol: { span: 6 },
  wrapperCol: { span: 14, offset: 7 }
}

interface IProps {
  licenseConfig?: any
  onSubmit: (values) => void
  onCancel: () => void
}

export default observer(function ConfigForm(props: IProps) {
  const [form] = Form.useForm(null)
  const { licenseConfig, onSubmit, onCancel } = props

  if (!licenseConfig) {
    form.resetFields()
  }

  function onFinish(values) {
    const { time } = values
    const start_time = moment(time[0] || '').format('YYYY-MM-DD HH:mm:ss')
    const end_time = moment(time[1] || '').format('YYYY-MM-DD HH:mm:ss')
    delete values.time
    if (licenseConfig.license_name) {
      Http.put(`licenseInfos/${licenseConfig?.id}`, {
        ...values,
        start_time,
        end_time,
        id: licenseConfig?.id,
        manager_id: licenseConfig?.manager_id
      }).then(res => {
        if (res.success) {
          onSubmit('ok')
        }
      })
    } else {
      Http.post('licenseInfos', {
        ...values,
        start_time,
        end_time,
        manager_id: licenseConfig.id
      }).then(res => {
        if (res.success) {
          onSubmit('ok')
        }
      })
    }
  }

  function editLine(value, index, type) {
    const keys = form.getFieldValue(type)

    keys.splice(index, 1, value)

    form.setFieldsValue({
      [type]: keys
    })
  }

  function addLine(type) {
    const keys = form.getFieldValue(type)

    if (keys.some(item => item === '' || item === undefined)) {
      form.validateFields(['allowable_hpc_endpoints'])
      return
    }

    const list = keys.filter(Boolean)

    if (list.length > 1 && list.length !== new Set(list).size) {
      form.validateFields(['allowable_hpc_endpoints'])
      // message.error('模块名称不能重复')
      return
    }

    form.setFieldsValue({
      [type]: [...new Set([...list, undefined])]
    })
  }

  return (
    <StyledLayout>
      <div className='form'>
        <Form
          name='config_form'
          form={form}
          {...colProps}
          initialValues={{
            license_env_var: licenseConfig?.license_env_var || '',
            allowable_hpc_endpoints: licenseConfig?.allowable_hpc_endpoints ?? [
              ''
            ],
            mac_addr: licenseConfig?.mac_addr || '',
            tool_path: licenseConfig?.tool_path || '',
            license_url: licenseConfig?.license_url || '',
            port: licenseConfig?.port || 0,
            license_num: licenseConfig?.license_num || '',
            license_name: licenseConfig?.license_name || '',
            // module_conf: licenseConfig?.module_conf.map(
            //   conf => conf.moduleName
            // ) || [undefined],
            // module_nums: licenseConfig?.module_conf.map(conf => conf.num) || [
            //   0
            // ],
            weight: licenseConfig?.weight || 0,
            auth:
              licenseConfig?.auth === undefined ? true : licenseConfig?.auth,
            time:
              licenseConfig?.begin_time && licenseConfig?.end_time
                ? [licenseConfig?.begin_time, licenseConfig?.end_time]
                : null,
            collector_type: licenseConfig?.collector_type ?? ''
          }}
          onFinish={onFinish}>
          <Form.Item
            label='许可证名字'
            name='license_name'
            required
            rules={[{ required: true, message: '许可证名字不能为空' }]}>
            <Input placeholder='请输入许可证名字' />
          </Form.Item>
          <Form.Item
            label='许可证变量'
            name='license_env_var'
            required
            rules={[{ required: true, message: '许可证变量不能为空' }]}>
            <Input placeholder='请输入许可证变量' />
          </Form.Item>
          <Form.Item
            label='Mac地址'
            name='mac_addr'
            rules={[
              // { required: true, message: 'Mac地址不能为空' },
              {
                pattern: /^([0-9a-fA-F]{2})(([/\s:-][0-9a-fA-F]{2}){5})$/,
                message: 'Mac地址格式不正确'
              }
            ]}>
            <Input placeholder='请输入Mac地址' />
          </Form.Item>
          <Form.Item
            label='许可证服务器地址'
            name='license_url'
            required
            rules={[{ required: true, message: '许可证服务器地址不能为空' }]}>
            <Input placeholder='请输入许可证服务器地址' />
          </Form.Item>
          <Form.Item
            label='端口'
            name='port'
            required
            tooltip={'范围 0 ~ 65535'}
            rules={[{ required: true, message: '端口不能为空' }]}>
            <InputNumber
              placeholder='请输入 0 ~ 65535 中的数字'
              min={0}
              max={65535}
            />
          </Form.Item>
          <Form.Item
            label='许可证序列号'
            name='license_num'
            rules={[
              // { required: true, message: '许可证序列号不能为空' }
            ]}>
            <Input placeholder='请输入许可证序列号' />
          </Form.Item>
          {/* <Form.List
            name='allowable_hpc_endpoints'
            rules={[
              {
                validator: async (_, allowable_hpc_endpoints) => {
                  if (
                    allowable_hpc_endpoints.some(
                      item => item === '' || item === undefined
                    )
                  ) {
                    return Promise.reject(new Error('标准计算服务不能为空'))
                  }

                  const list = allowable_hpc_endpoints.filter(Boolean)

                  if (list.length > 1 && list.length !== new Set(list).size) {
                    return Promise.reject(new Error('标准计算服务不能重复'))
                  }

                  if (
                    !allowable_hpc_endpoints ||
                    allowable_hpc_endpoints.length < 1 ||
                    allowable_hpc_endpoints.filter(Boolean).length === 0
                  ) {
                    return Promise.reject(
                      new Error('请至少添加一个标准计算服务')
                    )
                  }
                  if (allowable_hpc_endpoints.filter(Boolean).length === 0) {
                    return Promise.reject(
                      new Error('请至少添加一个标准计算服务')
                    )
                  }
                }
              }
            ]}>
            {(fields: FormListFieldData[], { add, remove }, { errors }) => (
              <>
                {fields.map((field, index) => {
                  const allowable_hpc_endpoints = form.getFieldValue(
                    'allowable_hpc_endpoints'
                  )

                  return (
                    <Form.Item
                      {...(index === 0 ? colProps : formItemLayoutWithOutLabel)}
                      label={index === 0 ? '标准计算服务' : ''}
                      required
                      shouldUpdate
                      key={field.key}>
                      <Row gutter={5} align={'middle'} style={{ margin: 0 }}>
                        <Col>
                          <Form.Item {...field} noStyle required>
                            <Input
                              style={{ width: 200 }}
                              value={allowable_hpc_endpoints[index] || null}
                              placeholder='请输入标准计算服务地址'
                              onChange={e => {
                                editLine(
                                  e.target.value.trim(),
                                  index,
                                  'allowable_hpc_endpoints'
                                )
                              }}
                            />
                          </Form.Item>
                        </Col>

                        <Col>
                          <MinusCircleOutlined
                            className='dynamic-delete-button'
                            onClick={() => {
                              if (fields.length <= 1) {
                                message.error('请至少添加一个标准计算服务')
                                return
                              }
                              remove(field.name)
                            }}
                          />
                        </Col>
                      </Row>
                    </Form.Item>
                  )
                })}

                <Form.Item {...formItemLayoutWithOutLabel}>
                  <Button
                    type='dashed'
                    onClick={() => {
                      addLine('allowable_hpc_endpoints')
                    }}
                    style={{ width: '100%', marginBottom: 10 }}
                    icon={<PlusOutlined />}>
                    添加
                  </Button>
                  <Form.ErrorList errors={errors} />
                </Form.Item>
              </>
            )}
          </Form.List> */}
          <Form.Item
            label='调度优先级'
            name='weight'
            tooltip={'请输入一个整数, 数字越小优先级越高'}>
            <InputNumber
              placeholder='请输入一个自然数, 数字越小优先级越高, 0的优先级最高'
              min={0}
            />
          </Form.Item>
          <Form.Item
            label='有效时间'
            required
            name='time'
            rules={[{ required: true, message: '有效时间不能为空' }]}>
            <RangePicker showTime />
          </Form.Item>
          <Form.Item label='是否授权' name='auth' valuePropName='checked'>
            <Switch />
          </Form.Item>
          <Form.Item
            label='许可证服务器'
            name='collector_type'
            required
            rules={[
              {
                required: true,
                message: '请选择许可证服务器'
              }
            ]}>
            <Select>
              {CollectorTypeLabel.map((label, index) => (
                <Select.Option key={index} value={label}>
                  {label}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            label='工具路径'
            name='tool_path'
            required
            rules={[{ required: true, message: '路径不能为空' }]}>
            <Input placeholder='请输入路径' />
          </Form.Item>
        </Form>
        <div className='footer'>
          <Button className='btn' onClick={() => onCancel()}>
            关闭
          </Button>
          <Button className='btn' type='primary' onClick={() => form.submit()}>
            保存
          </Button>
        </div>
      </div>
    </StyledLayout>
  )
})
