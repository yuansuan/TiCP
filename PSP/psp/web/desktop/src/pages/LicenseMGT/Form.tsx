import React from 'react'
import { Button } from '@/components'
import { observer } from 'mobx-react-lite'
import { Form, Input } from 'antd'
import { FormWrapper } from './style'

interface IProps {
  rowData?: {
    url: ''
    user: ''
    passwd: ''
  }
}

export const LoginForm = observer((props: IProps) => {
  const { rowData } = props
  const [form] = Form.useForm()

  const onFinish = values => {
    const { url, user, passwd } = values
    const token = btoa(`${user}:${passwd}`)
    localStorage.setItem('licenseMgrInfos', `${url}:${token}`)
    window.open(`${url}?token=${token}`, '_blank')
  }

  const formItemLayout = {
    labelCol: { span: 6, offset: 1 },
    wrapperCol: { span: 14 },
  }

  return (
    <FormWrapper>
      <h3>注意: 该浏览器第一次访问许可证管理系统, 需要进行登陆验证</h3>
      <Form
        {...formItemLayout}
        form={form}
        onFinish={onFinish}
        initialValues={{
          url: rowData?.url || '',
          user: '',
          passwd: '',
        }}>
        <Form.Item
          label='许可证管理系统URL'
          name='url'
          rules={[
            { required: true, message: '不能为空' },
            {
              pattern: /^(http|https):\/\/([\w.]+\/?)\S*/,
              message: 'url格式不对',
            },
          ]}>
          <Input />
        </Form.Item>
        <Form.Item
          label='登陆账户'
          name='user'
          rules={[{ required: true, message: '不能为空' }]}>
          <Input />
        </Form.Item>
        <Form.Item
          label='登录密码'
          name='passwd'
          rules={[{ required: true, message: '不能为空' }]}>
          <Input.Password />
        </Form.Item>
      </Form>
      <Button onClick={form.submit}>登陆</Button>
    </FormWrapper>
  )
})
