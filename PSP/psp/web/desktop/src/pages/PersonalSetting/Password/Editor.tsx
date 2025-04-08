/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Form, Input, message } from 'antd'
import { Button, Modal } from '@/components'
import { StyledEditor } from './style'
import { currentUser, env } from '@/domain'
import { validatePwd } from '@ys/utils'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Captcha } from '@/components'
import { userServer } from '@/server'

const colProps = { labelCol: { span: 4 }, wrapperCol: { span: 14 } }

interface IProps {
  onOk: () => Promise<any> | void
  onCancel: () => Promise<any> | void
}

export const Editor = observer(function Editor(props: IProps) {
  const store = useLocalStore(() => ({
    loading: false,
    updateLoading(loading) {
      this.loading = loading
    },
    captcha: '',
    updateCaptcha(captcha) {
      this.captcha = captcha
    }
  }))

  const [form] = Form.useForm()

  const onFinish = async values => {
    try {
      store.updateLoading(true)
      const oldpwd = values['oldpwd']
      const npwd = values['npwd']

      await userServer.updatePassword({
        password: oldpwd,
        newPassword: npwd,
        name: currentUser.name
      })

      message.success('密码修改成功')
      const { onOk } = props
      onOk()

      // 密码修改成功需要重新登录
      env.logout()
    } finally {
      store.updateLoading(false)
    }
  }

  const submit = () => {
    form.submit()
  }

  const { onCancel } = props

  return (
    <StyledEditor>
      <div className='body'>
        <Form form={form} onFinish={onFinish} {...colProps}>
          {/* <Form.Item label='手机号'>
            <span>{currentUser.mobile}</span>
          </Form.Item>
          <Form.Item
            name='captcha'
            label='验证码'
            rules={[
              {
                required: true,
                message: '验证码不能为空',
              },
            ]}>
            <Captcha
              phone={currentUser.phone}
              skip_check
              captcha={store.captcha}
              setCaptcha={code => {
                store.updateCaptcha(code)
                form.setFieldsValue({
                  captcha: code,
                })
              }}
            />
          </Form.Item> */}
          <Form.Item
            name='oldpwd'
            label='原密码'
            rules={[
              {
                required: true,
                message: '密码不能为空'
              }
              // {
              //   pattern: validatePwd.reg,
              //   message: '请输入 8-16 位包含数字和字母的非空字符'
              // }
            ]}>
            <Input.Password placeholder='请输入原密码' />
          </Form.Item>
          <Form.Item
            name='npwd'
            label='新密码'
            rules={[
              {
                required: true,
                message: '密码不能为空'
              }
              // {
              //   pattern: validatePwd.reg,
              //   message: '请输入 8-16 位包含数字和字母的非空字符',
              // },
            ]}>
            <Input.Password placeholder='请输入新密码' />
          </Form.Item>
          <Form.Item
            name='rpwd'
            label='确认密码'
            dependencies={['npwd']}
            rules={[
              {
                required: true,
                message: '确认密码不能为空'
              },
              ({ getFieldValue }) => ({
                validator(rule, value) {
                  if (!value || getFieldValue('npwd') === value) {
                    return Promise.resolve()
                  }
                  return Promise.reject('两次密码输入不一致')
                }
              })
            ]}>
            <Input.Password placeholder='请确认新密码' />
          </Form.Item>
        </Form>
        <Modal.Footer
          className='footer'
          onCancel={onCancel}
          OkButton={
            <Button type='primary' loading={store.loading} onClick={submit}>
              确认
            </Button>
          }
        />
      </div>
    </StyledEditor>
  )
})
