/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useEffect } from 'react'
import { Button, Modal } from '@/components'
import { observer } from 'mobx-react-lite'
import { Form, Input, Transfer, Tag, message, Modal as AntModal } from 'antd'
import { StyledOperators, FormWrapper } from './style'
import { RowData } from './Type'
import { Department } from '@/domain/DepartmentList/Department'
import { PlusOutlined } from '@ant-design/icons'
import { companyServer, departmentServer } from '@/server'
import { env } from '@/domain'
import { toJS } from 'mobx'

interface IProps {
  rowData: RowData
  refresh: () => Promise<any>
}

interface ModalProps {
  onCancel: () => void
  onOk: () => void | Promise<void>
  rowData?: RowData
  isAdding: boolean
  isPreview?: boolean
  refresh: () => void
}

interface UserSelectorProps {
  onCancel: () => void
  onOk: (users: any) => void | Promise<void>
  users: any[]
}

export const UserSelector = observer(
  ({ onCancel, onOk, users }: UserSelectorProps) => {
    const [allData, setAllData] = useState([])
    const [targetKeys, setTargetKeys] = useState([])
    const [selectedKeys, setSelectedKeys] = useState([])

    const onChange = (nextTargetKeys, direction, moveKeys) => {
      if (nextTargetKeys.length > 1000) {
        message.error('单个部门最多支持 1000 个成员')
        return
      }
      setTargetKeys(nextTargetKeys)
    }

    const onSelectChange = (sourceSelectedKeys, targetSelectedKeys) => {
      setSelectedKeys([...sourceSelectedKeys, ...targetSelectedKeys])
    }

    const onScroll = (direction, e) => {}

    useEffect(() => {
      ;(async () => {
        const res = await companyServer.queryUsers({
          company_id: env.company.id,
          status: 1, // normal
          page_index: 1,
          page_size: 1000
        })

        setAllData(res.data.list)
        setTargetKeys(users.map(user => user.user_id))
      })()
    }, [])

    return (
      <>
        <Transfer
          style={{ marginBottom: 10 }}
          dataSource={allData}
          titles={['所有成员', '已选成员']}
          targetKeys={targetKeys}
          selectedKeys={selectedKeys}
          onChange={onChange}
          onSelectChange={onSelectChange}
          onScroll={onScroll}
          rowKey={record => record.user_id}
          render={item => item.user_name || item.real_name || item.phone}
        />
        <Modal.Footer
          onCancel={onCancel}
          onOk={() => {
            onOk(allData.filter(item => targetKeys.includes(item.user_id)))
          }}
        />
      </>
    )
  }
)

export const EditingModal = observer(
  ({ onCancel, onOk, rowData, refresh, isAdding, isPreview }: ModalProps) => {
    if (isAdding) {
      // 初始化为空 Department 对象
      rowData = new Department()
    }

    const [newUsers, setNewUsers] = useState([])

    useEffect(() => {
      setNewUsers(toJS(rowData.users || []))
    }, [])

    const onOpen = async item => {
      await Modal.show({
        title: '选择部门成员',
        footer: null,
        content: ({ onCancel, onOk }) => {
          const OK = users => {
            setNewUsers(users)
            onOk()
          }
          return <UserSelector onCancel={onCancel} onOk={OK} users={newUsers} />
        }
      })
    }

    const [form] = Form.useForm()

    const onFinish = values => {
      if (isPreview) {
        onOk()
        return
      }

      if (isAdding) {
        departmentServer
          .create({
            ...values,
            user_ids: newUsers.map(u => u.user_id),
            old_user_ids: toJS(rowData.users).map(u => u.user_id)
          })
          .then(() => {
            message.success('新增部门成功')
            onOk()
            refresh()
          })
      } else {
        departmentServer
          .edit({
            ...values,
            department_id: rowData.id,
            user_ids: newUsers.map(u => u.user_id),
            old_user_ids: toJS(rowData.users).map(u => u.user_id)
          })
          .then(() => {
            message.success('编辑部门成功')
            onOk()
            refresh()
          })
      }
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
            name: rowData.name,
            remark: rowData.remark,
            users: rowData['users']
          }}>
          <Form.Item
            label='部门名称'
            name='name'
            rules={[{ required: true, message: '部门名称不能为空' }]}>
            {isPreview ? <span>{rowData.name}</span> : <Input />}
          </Form.Item>
          <Form.Item label='部门成员'>
            {newUsers?.map(item => {
              return (
                <Tag key={item.user_id} title={`手机号：${item.phone}`}>
                  {item.user_name || item.real_name || item.phone}
                </Tag>
              )
            })}
            {!isPreview && (
              <Button
                type='primary'
                shape='circle'
                icon={<PlusOutlined />}
                onClick={onOpen}
                size={'small'}>
                选择成员
              </Button>
            )}
          </Form.Item>
          <Form.Item
            label='部门备注'
            name='remark'
            rules={[{ type: 'string', max: 200 }]}>
            {isPreview ? (
              <span>{rowData.remark}</span>
            ) : (
              <Input.TextArea
                maxLength={200}
                showCount={{
                  formatter: ({ count, maxLength }) =>
                    `还可输入${maxLength - count}字`
                }}
              />
            )}
          </Form.Item>
        </Form>
        <Modal.Footer onCancel={onCancel} onOk={form.submit} />
      </FormWrapper>
    )
  }
)

const DepartmentAction = observer(({ rowData, refresh }: IProps) => {
  const [visible, setVisible] = useState(false)
  const onEdit = async item => {
    await Modal.show({
      title: '编辑部门',
      footer: null,
      content: ({ onCancel, onOk }) => (
        <EditingModal
          isAdding={false}
          onCancel={onCancel}
          onOk={onOk}
          rowData={rowData}
          refresh={refresh}
        />
      )
    })
  }

  const onDelete = item => {
    const { id: department_id, name } = item
    AntModal.confirm({
      title: '删除部门',
      centered: true,
      visible,
      content: `请确认是否删除【${name}】部门？`,
      okText: '确认',
      cancelText: '取消',
      onOk() {
        return new Promise((resolve, reject) => {
          departmentServer
            .delete({ department_id })
            .then(({ data }) => {
              if (data?.msg === 'success') {
                message.success('部门删除成功')
                refresh()
              } else {
                message.error('删除部门失败', data?.msg)
              }
              resolve(data)
            })
            .finally(() => {
              setVisible(false)
            })
        })
      },
      onCancel() {}
    })
  }

  return (
    <StyledOperators>
      <Button type='link' onClick={() => onEdit(rowData)}>
        编辑
      </Button>
      <Button type='link' onClick={() => onDelete(rowData)}>
        删除
      </Button>
    </StyledOperators>
  )
})

export default DepartmentAction
