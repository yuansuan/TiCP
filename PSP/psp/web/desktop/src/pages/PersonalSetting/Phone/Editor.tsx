/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { message, Form, Input } from 'antd'
import { Modal, Button } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { StyledEditor } from './style'
import { currentUser } from '@/domain'
import { validatePhone } from '@ys/utils'
import { Captcha } from '@/components'
import { userServer } from '@/server'
import { runInAction } from 'mobx'

interface IProps {
  token: string
  onOk: () => Promise<any> | void
  onCancel: () => Promise<any> | void
}

const colProps = { labelCol: { span: 4 }, wrapperCol: { span: 12 } }

export const Editor = observer(function Editor({
  token,
  onOk,
  onCancel,
}: IProps) {
  const [form] = Form.useForm()

  const state = useLocalStore(() => ({
    phone: '',
    updatePhone(phone) {
      this.phone = phone
    },
    captcha: '',
    setCaptcha(captcha) {
      this.captcha = captcha
    },
    loading: false,
    updateLoading(loading) {
      this.loading = loading
    },
  }))

  const onFinish = async () => {
    try {
      state.updateLoading(true)
      await userServer.updatePhone({
        phone: state.phone,
        captcha: state.captcha,
        token,
        oldPhone: currentUser.phone,
      })
      runInAction(() => {
        currentUser.update({
          phone: state.phone,
        })
      })
      message.success('手机号绑定成功')
      onOk()
    } finally {
      state.updateLoading(false)
    }
  }

  const submit = () => {
    form.submit()
  }

  return (
    <StyledEditor>
      <div className='body'>
        <Form form={form} onFinish={onFinish} {...colProps}>
          <Form.Item
            name='phone'
            label='新手机号'
            rules={[
              {
                required: true,
                message: '新手机号不能为空',
              },
              {
                pattern: validatePhone.reg,
                message: '新手机号格式不正确',
              },
            ]}>
            <Input
              autoFocus
              placeholder='请输入手机号'
              value={state.phone}
              onChange={e => state.updatePhone(e.target.value.trim())}
            />
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
              phone={state.phone}
              captcha={state.captcha}
              setCaptcha={code => {
                state.setCaptcha(code)
                form.setFieldsValue({
                  captcha: code,
                })
              }}
            />
          </Form.Item>
        </Form>
        <Modal.Footer
          className='footer'
          onCancel={onCancel}
          OkButton={
            <Button type='primary' loading={state.loading} onClick={submit}>
              确认
            </Button>
          }
        />
      </div>
    </StyledEditor>
  )
})
