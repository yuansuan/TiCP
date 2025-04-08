import React, { useState, useEffect } from 'react'
import { Button, Modal } from '@/components'
import { observer } from 'mobx-react-lite'
import { Form, Input, Transfer, Tag, message, DatePicker, Select } from 'antd'
import { FormWrapper } from './style'
import Project from '@/domain/ProjectMG/Project'
import { PROJECT_STATE_MAP, adminProjectMG as projectMG} from '@/domain/ProjectMG'
import { PlusOutlined } from '@ant-design/icons'
import { toJS } from 'mobx'
import { DatePicker_FORMAT, DatePicker_SHOWTIME_FORMAT } from '@/constant'
import { currentUser } from '@/domain'

const { RangePicker } = DatePicker


interface FormProps {
  onCancel: () => void
  onOk: () => void | Promise<void>
  rowData?: Project
  isAdding: boolean
  isPreview?: boolean
  refresh: () => void
}

interface UserSelectorProps {
  onCancel: () => void
  onOk: (users: any) => void | Promise<void>
  users: any[]
  disabledKeys: any[]
}

export const EditingForm = observer(
  ({ onCancel, onOk, rowData, refresh, isAdding, isPreview }: FormProps) => {
    if (isAdding) {
      // 初始化为空 Project 对象
      rowData = new Project({
        project_owner_id: currentUser.id,
        members: [{user_id: currentUser.id, user_name: currentUser.name}],
      })
    }

    const [newUsers, setNewUsers] = useState([])

    useEffect(() => {
      setNewUsers(toJS(rowData?.members || []))
    }, [])

    const [form] = Form.useForm()

    const onFinish = values => {
      if (isPreview) {
        onOk()
        return
      }

      if (isAdding) {
        projectMG
          .add({
            ...values,
            start_time: values?.times[0]?.unix(),
            end_time: values?.times[1]?.unix(),
            members: newUsers.map(u => u.user_id)
          })
          .then(() => {
            message.success('新增项目成功')
            onOk()
            refresh()
          })
      } else {
        projectMG
          .edit({
            ...values,
            start_time: values?.times[0]?.unix(),
            end_time: values?.times[1]?.unix(),
            project_id: rowData.id,
            project_owner_id: rowData.project_owner_id,
            members: newUsers.map(u => u.user_id)
          })
          .then(() => {
            message.success('编辑项目成功')
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
            project_name: rowData.project_name,
            comment: rowData.comment,
            users: rowData['users'],
            times: [rowData.start_time_momnet, rowData.end_time_momnet]
          }}>
          <Form.Item
            label='项目名称'
            name='project_name'
            rules={[
              { required: true, message: '项目名称不能为空' },
              // 项目名称只能包含字母数字下划线, 最大长度不能超过64位
              { pattern: /^(?!_)[a-zA-Z0-9_]{1,64}$/, message: '项目名称只能包含字母数字下划线，并且不能以下划线开头, 最大长度不能超过64位' }
            ]}>
            {isPreview ? <span>{rowData.project_name}</span> : <Input disabled={!isAdding}/>}
          </Form.Item>
          {
            (isPreview || isAdding) && (
              <Form.Item
                label='项目管理员'
                name='project_owner_name'>
                <span>{rowData.project_owner_name || currentUser.name}</span>
              </Form.Item>
            )
          }
          {
            isPreview && (
              <>
                <Form.Item
                  label='项目状态'
                  name='state'>
                  <span>{PROJECT_STATE_MAP[rowData.state]}</span>
                </Form.Item>
              </>
            )
          }
          <Form.Item
            label='项目周期'
            name='times'
            rules={[
              { required: true, message: '项目周期不能为空' },
              {
                validator: (rule, value) => {
                  if (value.length < 2 || value.some((item) => !item)) {
                    return Promise.reject("项目周期不能为空");
                  }
                  return Promise.resolve();
                }
              }
            ]}>
            {isPreview ? <span>{`${rowData.start_time} - ${rowData.end_time}`}</span> : <RangePicker
              showTime={{ format: DatePicker_SHOWTIME_FORMAT }}
              format={DatePicker_FORMAT}
            />}
          </Form.Item>
          <Form.Item label='部门成员'>
            {newUsers?.map(item => {
              return (
                <Tag key={item.user_id}>
                  {item.user_name}
                </Tag>
              )
            })}
            {!isPreview && (
              <Button
                type='primary'
                shape='circle'
                icon={<PlusOutlined rev={'none'}/>}
                onClick={() => openUserSelector({
                  project_owner_id: rowData.project_owner_id,
                  members: newUsers,
                }, (users) => setNewUsers(users))}
                size={'small'}>
                选择成员
              </Button>
            )}
          </Form.Item>
          <Form.Item
            label='项目备注'
            name='comment'
            rules={[{ type: 'string', max: 200 }]}>
            {isPreview ? (
              <span>{rowData.comment}</span>
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


export const openUserSelector = async (item, callback?) => {
  await Modal.show({
    title: '选择项目成员',
    width: 600,
    footer: null,
    content: ({ onCancel, onOk }) => {
      const OK = users => {
        callback && callback(users)
        onOk()
      }
      return <UserSelector onCancel={onCancel} onOk={OK} users={item.members || []} disabledKeys={[item.project_owner_id]}/>
    }
  })
}

export const UserSelector = observer(
  ({ onCancel, onOk, users, disabledKeys }: UserSelectorProps) => {
    const [oldTargetKeys, setOldTargetKeys] = useState([])
    const [allData, setAllData] = useState([])
    const [targetKeys, setTargetKeys] = useState([])
    const [selectedKeys, setSelectedKeys] = useState([])

    const globalConfig = JSON.parse(localStorage.getItem('GlobalConfig') || '{}')

    const onChange = (nextTargetKeys, direction, moveKeys) => {
      if (nextTargetKeys.length > 100) {
        message.error('单个项目最多支持 100 个成员')
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
        const res = await projectMG.getALLUsers()
        const allUsers = res?.data?.map(u => {
          return {user_id: u.key, user_name: u.title, disabled: disabledKeys.includes(u.key)}
        })  || []
        setAllData(allUsers)
        setTargetKeys(users.map(user => user.user_id) || [])
        setOldTargetKeys(users.map(user => user.user_id) || [])
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
          render={item => item.user_name}
        />
        <Modal.Footer
          onCancel={onCancel}
          onOk={() => {
            const removedUserIds = oldTargetKeys.filter(item => !targetKeys.includes(item))
            const removedUser = allData.filter(item => removedUserIds.includes(item.user_id))

            if (removedUser && removedUser.length !== 0) {
              Modal.confirm({
                title: '移除成员',
                content: (
                  <>
                    <p>{`确定要从当前项目中移除成员 ${removedUser.map(u => u.user_name).join(',')} 吗？`}</p>
                    { globalConfig?.enable_visual &&
                      <p>
                       注意：移除成员在当前项目中开启的3D云应用会话将被自动删除，请提前通知备份数据。
                      </p>
                    }
                  </>
                ),
                onOk: () => {
                  onOk(allData.filter(item => targetKeys.includes(item.user_id)))
                }
              })
            } else {
              onOk(allData.filter(item => targetKeys.includes(item.user_id)))
            }
          }}
        />
      </>
    )
  }
)

export const OwnerForm = observer(
  ({ onCancel, onOk, rowData, refresh}: Partial<FormProps>) => {

    const [form] = Form.useForm()
    const [projectMgrs, setProjectMgrs] = useState([])

    const onFinish = values => {
      if (values.new_owner_id === rowData.project_owner_id) {
        message.error('请指定新的项目管理员')
        return
      } else {
        projectMG
          .changeOwner({
            project_id: rowData.id,
            new_owner_id: values.new_owner_id
          })
          .then(() => {
            message.success('项目所有者转移成功')
            onOk()
            refresh()
          })
      }
    }

    useEffect(() => {
      (async () => {
        const res = await projectMG.getALLUsersWithProjectMgrPerm() 
        setProjectMgrs(res?.data?.map(u => ({
          user_id: u.key,
          user_name: u.title
        })).filter(o => o.user_id !== rowData.project_owner_id) || [])
      })()
    }, [])

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
            project_name: rowData.project_name,
            new_owner_id: null,
          }}>
          <Form.Item
            label='项目名称'
            name='project_name'
            rules={[{ required: true, message: '项目名称不能为空' }]}>
            <span>{rowData.project_name}</span>
          </Form.Item>
          <Form.Item
            label='当前项目管理员'
            name='project_owner_name'>
            <span>{rowData.project_owner_name || currentUser.name}</span>
          </Form.Item>
          <Form.Item
            label='新项目管理员'
            name='new_owner_id'
            tooltip={'新项目管理员必须具有项目管理(项目管理员)权限，这里只显示具有该权限的用户'}
            rules={[{ required: true, message: '新项目管理员不能为空' }]}>
            <Select placeholder="请指定新的项目管理员">
              {
                projectMgrs.map(u => <Select.Option value={u.user_id}>{u.user_name}</Select.Option>)
              }
            </Select>
          </Form.Item>
        </Form>
        <Modal.Footer onCancel={onCancel} onOk={form.submit} />
      </FormWrapper>
    )
  }
) 

export const onChangeOwner = async (rowData, refresh) => {
  await rowData.getDetail()

  await Modal.show({
    title: '管理员转移',
    width: 680,
    footer: null,
    content: ({ onCancel, onOk }) => (
      <OwnerForm onCancel={onCancel} onOk={onOk} rowData={rowData} refresh={refresh}/>
    )
  })
}

export const onPreview = async (rowData) => {
  
  await rowData.getDetail()

  await Modal.show({
    title: '预览项目',
    width: 680,
    footer: null,
    content: ({ onCancel, onOk }) => (
      <EditingForm
        isAdding={false}
        isPreview={true}
        onCancel={onCancel}
        onOk={onOk}
        rowData={rowData}
        refresh={() => {}}
      />
    )
  })
}

export const onAdd = async (refresh) => {
  await Modal.show({
    title: '添加项目',
    width: 680,
    footer: null,
    content: ({ onCancel, onOk }) => (
      <EditingForm
        isAdding={true}
        onCancel={onCancel}
        onOk={onOk}
        refresh={refresh || (() => {})}
      />
    )
  })
}

export const onEdit = async (rowData, refresh) => {
  await rowData.getDetail()
  await Modal.show({
    title: '编辑项目',
    width: 680,
    footer: null,
    content: ({ onCancel, onOk }) => (
      <EditingForm
        isAdding={false}
        onCancel={onCancel}
        onOk={onOk}
        rowData={rowData}
        refresh={refresh || (() => {})}
      />
    )
  })
}