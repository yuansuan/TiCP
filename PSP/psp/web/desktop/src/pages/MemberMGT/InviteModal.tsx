/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { Form, Input, Select, message } from 'antd'
import { Modal } from '@/components'
import { env } from '@/domain'
import { companyServer } from '@/server'
import { SPACE_MGR_ROLE_NAMES } from '@/constant'
import { DepEmptyWrapper } from './style'

const StyledLayout = styled.div`
  padding-bottom: 40px;

  > .footer {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    padding: 10px;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};
  }
`

const { useForm } = Form
const colProps = { labelCol: { span: 4 }, wrapperCol: { span: 14 } }

type Props = {
  onCancel: () => void
  onOk: (result: any[]) => void
  departments: any
}

export const InviteModal = observer(function InviteModal({
  onCancel,
  onOk,
  departments
}: Props) {
  const [form] = useForm()
  const rolesList = env.company.roles.filter(role => role.name !== '超级管理员')
  const defaultRole = rolesList[0]?.id

  async function onFinish(values) {
    const department_id = values['department']
    const phones = (values['phones'] || '')
      .split(/[\n,\,,，,、]/)
      .map(item => item.trim())
      .filter(Boolean)

    if (env.company.isOpenDepMgr) {
      if (departments.length === 0) {
        message.error('您已启用部门管理，请先前往部门管理新建部门')
        return
      }

      if (!department_id) {
        message.error('请选择部门')
        return
      }
    }

    if (phones.length === 0) {
      message.error('请输入手机号')
      return
    }

    if (phones.length > 200) {
      message.error('最多添加 200 个手机号')
      return
    }

    const params = env.company.isOpenDepMgr
      ? {
          role_id: values['role'],
          phone_list: [...new Set(phones)],
          department_id: department_id
        }
      : {
          role_id: values['role'],
          phone_list: [...new Set(phones)]
        }

    const { data } = await companyServer.batchInvite(params)
    onOk(data)
  }

  function onSubmit() {
    form.submit()
  }

  return (
    <StyledLayout>
      <div className='body'>
        <Form
          form={form}
          onFinish={onFinish}
          {...colProps}
          initialValues={{
            role: defaultRole
          }}>
          <Form.Item name='role' label='选择角色' required>
            <Select>
              {rolesList
                .filter(role => role.type === 1)
                .map(role => (
                  <Select.Option key={role.id} value={role.id}>
                    {role.name}
                  </Select.Option>
                ))}
            </Select>
          </Form.Item>

          <Form.Item
            name='phones'
            label='手机号'
            required
            help={
              <div>
                <div>
                  *可输入多个手机号，支持换行、逗号、顿号分隔，单次最多新增200个
                </div>
                <div>*已加入其他企业的手机号不能被邀请</div>
              </div>
            }>
            <Input.TextArea style={{ height: 85 }} />
          </Form.Item>
        </Form>
      </div>
      <Modal.Footer
        className='footer'
        onCancel={onCancel}
        onOk={() => {
          onSubmit()
        }}
      />
    </StyledLayout>
  )
})
