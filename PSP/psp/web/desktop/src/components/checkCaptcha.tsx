/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Modal, Button } from '@/components'
import { Form } from 'antd'
import { currentUser } from '@/domain'
import { Captcha } from '@/components'
import { Http } from '@/utils'

const StyledLayout = styled.div`
  .body {
    padding-bottom: 40px;
  }

  .footer {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    padding: 10px 0;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};
  }
`

const colProps = { labelCol: { span: 4 }, wrapperCol: { span: 12 } }
const { useForm } = Form

type Props = {
  onOk: (token?: string) => Promise<any> | void
  onCancel: () => Promise<any> | void
}

export const CaptchaChecker = observer(function CaptchaChecker({
  onCancel,
  onOk,
}: Props) {
  const state = useLocalStore(() => ({
    loading: false,
    setLoading(flag) {
      this.loading = flag
    },
    captcha: '',
    setCaptcha(captcha) {
      this.captcha = captcha
    },
  }))
  const [form] = useForm()

  async function onFinish() {
    try {
      state.setLoading(true)
      const {
        data: { token },
      } = await Http.post('/captcha', {
        phone: currentUser.phone,
        captcha: state.captcha,
      })

      onOk(token)
    } finally {
      state.setLoading(false)
    }
  }

  function submit() {
    form.submit()
  }

  return (
    <StyledLayout>
      <div className='body'>
        <Form form={form} {...colProps} onFinish={onFinish}>
          <Form.Item label='手机号'>
            <span>{currentUser.phone}</span>
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
      </div>
      <Modal.Footer
        className='footer'
        onCancel={onCancel}
        OkButton={
          <Button type='primary' loading={state.loading} onClick={submit}>
            下一步
          </Button>
        }
      />
    </StyledLayout>
  )
})

export const checkCaptcha = () =>
  Modal.show({
    title: '手机验证',
    content: ({ onCancel, onOk }) => (
      <CaptchaChecker onCancel={onCancel} onOk={onOk} />
    ),
    footer: null,
  })
