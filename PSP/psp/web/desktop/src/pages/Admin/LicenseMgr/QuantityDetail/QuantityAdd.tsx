import React from 'react'
import { Form, Input, Col, Row, InputNumber, Button, message } from 'antd'
import { Modal } from '@/components'
import { FormListFieldData } from 'antd/es/form/FormList'
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons'
import styled from 'styled-components'
import { Http } from '@/utils'
const colProps = { labelCol: { span: 6 }, wrapperCol: { span: 14, offset: 1 } }

const formItemLayoutWithOutLabel = {
  labelCol: { span: 6 },
  wrapperCol: { span: 14, offset: 7 }
}

export const StyledLayout = styled.div`
  padding: 30px;

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
interface IProps {
  moduleConf: any
  licenseId: string
  onCancel: () => void
  refresh: (id) => Promise<any> | void
}

export const QuantityAdd = ({
  moduleConf,
  licenseId,
  refresh,
  onCancel
}: IProps) => {
  const [form] = Form.useForm(null)
  function onFinish(values) {
    const { module_conf } = values

    let module_nums = form.getFieldValue('module_nums')

    const newModuleConf = module_conf.map((name, index) => {
      return {
        module_name: name,
        total: module_nums[index] || 0
      }
    })
    if (moduleConf[0].module_name) {
      Http.put(`/moduleConfigs/${moduleConf[0].id}`, {
        id: moduleConf[0].id,
        license_id: licenseId,
        module_name: newModuleConf[0].module_name,
        total: newModuleConf[0].total
      }).then(res => {
        if (res.success) {
          refresh(licenseId)
          onCancel()
        }
      })
    } else {
      Http.post('/moduleConfigs', {
        license_id: licenseId,
        module_name: newModuleConf[0].module_name,
        total: newModuleConf[0].total
      }).then(res => {
        if (res.success) {
          refresh(licenseId)
          onCancel()
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
      form.validateFields(['module_conf'])
      return
    }

    const list = keys.filter(Boolean)

    if (list.length > 1 && list.length !== new Set(list).size) {
      form.validateFields(['module_conf'])
      // message.error('模块名称不能重复')
      return
    }

    form.setFieldsValue({
      [type]: [...new Set([...list, undefined])]
    })

    // // module_nums
    // form.setFieldsValue({
    //   module_nums: [...form.getFieldValue('module_nums'), 0]
    // })
  }

  return (
    <StyledLayout>
      <div className='form'>
        <Form
          name='config_form'
          form={form}
          {...colProps}
          initialValues={{
            module_conf: moduleConf.map(conf => conf.module_name) || [
              undefined
            ],
            module_nums: moduleConf.map(conf => conf.total) || [0]
          }}
          onFinish={onFinish}>
          <Form.List
            name='module_conf'
            rules={[
              {
                validator: async (_, module_conf) => {
                  if (
                    module_conf.some(item => item === '' || item === undefined)
                  ) {
                    return Promise.reject(new Error('模块名称不能为空'))
                  }

                  const list = module_conf.filter(Boolean)

                  if (list.length > 1 && list.length !== new Set(list).size) {
                    return Promise.reject(new Error('模块名称不能重复'))
                  }

                  if (
                    !module_conf ||
                    module_conf.length < 1 ||
                    module_conf.filter(Boolean).length === 0
                  ) {
                    return Promise.reject(new Error('请至少添加一个模块'))
                  }
                  if (module_conf.filter(Boolean).length === 0) {
                    return Promise.reject(new Error('请至少添加一个模块'))
                  }
                }
              }
            ]}>
            {(fields: FormListFieldData[], { add, remove }, { errors }) => (
              <>
                {fields.map((field, index) => {
                  const module_conf = form.getFieldValue('module_conf')
                  const module_nums = form.getFieldValue('module_nums')

                  return (
                    <Form.Item
                      {...(index === 0 ? colProps : formItemLayoutWithOutLabel)}
                      label={index === 0 ? '软件模块配置' : ''}
                      required
                      shouldUpdate
                      key={field.key}>
                      <Row gutter={5} align={'middle'} style={{ margin: 0 }}>
                        <Col>
                          <Form.Item
                            {...field}
                            noStyle
                            required
                            rules={[
                              { required: true, message: '模块名称不能为空' }
                            ]}>
                            <Input
                              value={module_conf[index] || null}
                              placeholder='请输入模块名称'
                              onChange={e => {
                                editLine(
                                  e.target.value.trim(),
                                  index,
                                  'module_conf'
                                )
                              }}
                            />
                          </Form.Item>
                        </Col>
                        <Col>
                          <Form.Item noStyle>
                            <InputNumber
                              value={module_nums[index] || 0}
                              placeholder='请输入license数量'
                              min={0}
                              onChange={value => {
                                editLine(value, index, 'module_nums')
                              }}
                            />
                          </Form.Item>
                        </Col>
                        {/* <Col>
                        <MinusCircleOutlined
                          className='dynamic-delete-button'
                          onClick={() => {
                            if (fields.length <= 1) {
                              message.error('请至少添加一个模块')
                              return
                            }
                            remove(field.name)
                          }}
                        />
                      </Col> */}
                      </Row>
                    </Form.Item>
                  )
                })}

                {/* <Form.Item {...formItemLayoutWithOutLabel}>
                <Button
                  type='dashed'
                  onClick={() => {
                    addLine('module_conf')
                  }}
                  style={{ width: '100%', marginBottom: 10 }}
                  icon={<PlusOutlined />}>
                  添加
                </Button>
                <Form.ErrorList errors={errors} />
              </Form.Item> */}
              </>
            )}
          </Form.List>
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
}
