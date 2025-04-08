/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Modal, Button } from '@/components'
import styled from 'styled-components'
import { observer } from 'mobx-react'
import { observable, action } from 'mobx'
import { Form, Input } from 'antd'
import { FormInstance } from 'antd/lib/form'
import { userServer } from '@/server'

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

interface IProps {
  onOk: (password?: string) => Promise<any> | void
  onCancel: () => Promise<any> | void
}

@observer
export class PasswordChecker extends React.Component<IProps> {
  formRef = React.createRef<FormInstance>()
  @observable loading = false
  @action
  updateLoading = loading => {
    this.loading = loading
  }

  onFinish = async values => {
    try {
      this.updateLoading(true)
      const password = values['password']
      await userServer.checkPassword(password)
      this.props.onOk(password)
    } finally {
      this.updateLoading(false)
    }
  }

  submit = () => {
    this.formRef.current.submit()
  }

  render() {
    const { onCancel } = this.props

    return (
      <StyledLayout>
        <div className='body'>
          <Form ref={this.formRef} onFinish={this.onFinish}>
            <Form.Item
              name='password'
              rules={[
                {
                  required: true,
                  message: '密码不能为空',
                },
              ]}
              label='密码'
              {...colProps}>
              <Input.Password autoFocus placeholder='请输入密码' />
            </Form.Item>
          </Form>
          <Modal.Footer
            className='footer'
            onCancel={onCancel}
            OkButton={
              <Button
                type='primary'
                loading={this.loading}
                onClick={this.submit}>
                下一步
              </Button>
            }
          />
        </div>
      </StyledLayout>
    )
  }
}

export const checkPassword = () =>
  Modal.show({
    title: '密码验证',
    content: ({ onCancel, onOk }) => (
      <PasswordChecker onCancel={onCancel} onOk={onOk} />
    ),
    footer: null,
  })
