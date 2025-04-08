/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Button, Modal } from '@/components'
import { observer } from 'mobx-react-lite'
import { Form, Select, InputNumber, message, Divider } from 'antd'
import { StyledOperators, FormWrapper, DepEmptyWrapper } from './style'
import { currentUser, env } from '@/domain'
import { RowData } from './Type'
import { companyServer } from '@/server'
import { SPACE_MGR_ROLE_NAMES } from '@/constant'
import { useStore } from './model'

const { Option } = Select
interface IProps {
  rowData: RowData
  refreshMembers: () => Promise<any>
  departments: any
}

interface ModalProps {
  onCancel: () => void
  onOk: () => void | Promise<void>
  rowData: RowData
  refreshMembers: () => void
  departments: any
}

const SettingModal = observer(
  ({ onCancel, onOk, rowData, refreshMembers, departments }: ModalProps) => {
    const [form] = Form.useForm()
    const onFinish = values => {
      const department_id = values['department_id']

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

      const params = env.company.isOpenDepMgr
        ? {
            user_id: rowData.user_id,
            role_id: values.role_id,
            consume_limit: values.consume_limit,
            department_id: department_id
          }
        : {
            user_id: rowData.user_id,
            role_id: values.role_id,
            consume_limit: values.consume_limit
          }

      companyServer.configUser(params).then(() => {
        message.success('设置成功')
        onOk()
        refreshMembers()
      })
    }
    const formItemLayout = {
      labelCol: { span: 6, offset: 1 },
      wrapperCol: { span: 14 }
    }
    return (
      <FormWrapper>
        <Form
          {...formItemLayout}
          form={form}
          onFinish={onFinish}
          initialValues={{
            role_id: rowData?.role_list[0]?.id,
            consume_limit: rowData['consume_limit'],
            department_id: rowData?.department?.id
          }}>
          <Form.Item label='姓名'>
            <span>{rowData.real_name}</span>
          </Form.Item>
          <Form.Item label='角色' name='role_id'>
            <Select>
              {env.company.roles
                .filter(role => role.name !== '超级管理员')
                .filter(role =>
                  env.isSpaceManager
                    ? SPACE_MGR_ROLE_NAMES.includes(role.name)
                    : true
                )
                .map(role => (
                  <Option key={role.id} value={role.id}>
                    {role.name}
                  </Option>
                ))}
            </Select>
          </Form.Item>
          <Form.Item
            label='每月消费限额'
            name='consume_limit'
            rules={[{ type: 'number', min: 0.01, message: '最小值 0.01' }]}>
            <InputNumber placeholder='单位：元' precision={2} step={0.01} />
          </Form.Item>
        </Form>
        <Modal.Footer onCancel={onCancel} onOk={form.submit} />
      </FormWrapper>
    )
  }
)

const MemberAction = observer(
  ({ rowData, refreshMembers, departments }: IProps) => {
    const store = useStore()
    const onDelete = async item => {
      await Modal.showConfirm({
        title: '删除成员',
        content: `确认要删除“${
          item.real_name !== '' ? item.real_name : item.phone
        }”吗？`
      })
      companyServer
        .delete({
          user_id: item.user_id,
          company_name: env.company?.name,
          company_id: env.company?.id
        })
        .then(res => {
          const data = res.data
          if (data.msg === 'success') {
            message.success('用户删除成功')
            if (item.user_id === currentUser.id) {
              location.replace('/')
            }
          } else {
            data.is_run_job
              ? message.error('有正在运行的作业，无法删除该用户')
              : data.is_open_app
              ? message.error('有正在运行的应用，无法删除该用户')
              : message.error('删除用户失败')
          }
          refreshMembers()
        })
    }

    const onSetting = async item => {
      await Modal.show({
        title: '设置',
        footer: null,
        content: ({ onCancel, onOk }) => (
          <SettingModal
            onCancel={onCancel}
            onOk={onOk}
            rowData={rowData}
            refreshMembers={refreshMembers}
            departments={departments}
          />
        )
      })
    }

    const isDisabled = () => {
      if (rowData.role_list[0].name === '超级管理员') {
        return true
      }

      if (store.currentUserRole === '管理员') {
        if (rowData.role_list[0].name === '管理员') {
          return rowData['user_id'] === currentUser.id ? false : true
        }
      }
      return false
    }

    return (
      <StyledOperators>
        <>
          <Button
            type='link'
            onClick={() => onSetting(rowData)}
            title={
              env.isSpaceManager &&
              !SPACE_MGR_ROLE_NAMES.includes(rowData['roles'])
                ? '空间管理员不能对其它空间管理员和管理员进行操作'
                : ''
            }
            disabled={
              env.isSpaceManager
                ? !SPACE_MGR_ROLE_NAMES.includes(rowData['roles'])
                : false || isDisabled()
            }>
            设置
          </Button>
          <Divider type='vertical' />
          <Button
            type='link'
            onClick={() => onDelete(rowData)}
            title={
              env.isSpaceManager &&
              !SPACE_MGR_ROLE_NAMES.includes(rowData['roles'])
                ? '空间管理员不能对其它空间管理员和管理员进行操作'
                : ''
            }
            disabled={
              env.isSpaceManager
                ? !SPACE_MGR_ROLE_NAMES.includes(rowData['roles'])
                : false || isDisabled()
            }>
            删除
          </Button>
        </>
      </StyledOperators>
    )
  }
)

export default MemberAction
