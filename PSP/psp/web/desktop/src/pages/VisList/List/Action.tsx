/* Copyright (C) 2016-present, Yuansuan.cn */
import React, { useEffect,useState } from 'react'
import { observer } from 'mobx-react-lite'
import { Select, Form } from 'antd'
import { Button } from '@/components'
import { runInAction } from 'mobx'
import { ListActionWrapper } from '../style'
import { useStore } from './store'

import {
  Hardware,
  Software,
  ListSessionRequest,
  SESSION_STATUS_MAP
} from '@/domain/Vis'
import { currentUser } from '@/domain'

const { Option } = Select
const { useForm } = Form

const resetValue: ListSessionRequest = new ListSessionRequest({
  statuses: [],
  hardware_ids: [],
  software_ids: [],
  project_ids: [],
  user_name: ''
})

const selectStyle = { width: '200px' }

interface IProps {
  disabledCreate?: string
  value: ListSessionRequest
  hardware?: Array<Hardware>
  software?: Array<Software>
  onCreate: () => void
  onSubmit: (value: ListSessionRequest) => void
}

export const Action = observer(
  ({
    onSubmit,
    onCreate,
    value,
    hardware,
    software,
    disabledCreate
  }: IProps) => {
    const store = useStore()
    const { vis } = store
    const [form] = useForm()
    const [projects, setProjects] = useState([])
    const filterParams = vis.filterQuery
    const reset = () => {
      vis.setFilterParams({
        statuses: [],
        hardware_ids: [],
        software_ids: [],
        project_ids: [],
        user_name: '',
        page_index: 1,
        page_size: 20
      })
      form.setFieldsValue({ ...resetValue })
      submit()
    }

    useEffect(() => {
      (async() => {
        const res = await vis.getProjects(false)
        setProjects(res?.data?.projects || [])
      })()
      form.setFieldsValue({ ...resetValue })
      onSubmit(value)
    }, [])

    const onFinish = (values: any) => {
      runInAction(() => {
        values.page_index = 1
        values.page_size = 20
        values.status = values.status || null
        onSubmit(values)
        vis.setFilterParams(values)
      })
    }

    const submit = () => {
      form.submit()
    }

    return (
      <ListActionWrapper>
        <Form
          form={form}
          initialValues={{
            statuses: filterParams?.statuses,
            hardware_ids: filterParams?.hardware_ids,
            software_ids: filterParams?.software_ids,
            project_ids: filterParams?.project_ids
          }}
          layout='inline'
          onFinish={onFinish}
          className='item'>
          <Form.Item label='会话状态' name='statuses'>
            <Select
              className={'status'}
              value={value.statuses}
              mode='multiple'
              onChange={submit}
              allowClear={true}
              placeholder='所有会话状态'
              style={selectStyle}>
              {Object.entries(SESSION_STATUS_MAP).map(([key, name]) => (
                <Option value={key} key={name}>
                  {SESSION_STATUS_MAP[key]}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item label='实例名称' name='hardware_ids'>
            <Select
              className={'hardware_id'}
              value={value.hardware_ids}
              onChange={submit}
              mode='multiple'
              allowClear={true}
              placeholder='所有实例'
              style={selectStyle}>
              {(hardware || []).map(h => (
                <Option value={h.id} key={h.id}>
                  {h.name}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item label='镜像名称' name='software_ids'>
            <Select
              className={'software_id'}
              value={value.software_ids}
              onChange={submit}
              mode='multiple'
              allowClear={true}
              placeholder='所有镜像'
              style={selectStyle}>
              {(software || []).map(s => (
                <Option value={s.id} key={s.id}>
                  {s.name}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item label='项目名称' name='project_ids'>
            <Select
              className={'project_name'}
              value={value.project_ids}
              onChange={submit}
              mode='multiple'
              allowClear={true}
              placeholder='所有项目'
              style={selectStyle}>
              {(projects || []).map(n => (
                <Option value={n.id} key={n.id}>
                  {n.name}
                </Option>
              ))}
            </Select>
          </Form.Item>
        </Form>
        <div className='item'>
          <Button className={'btn'} onClick={reset}>
            重置
          </Button>
          <Button
            className={'btn'}
            disabled={disabledCreate}
            onClick={onCreate}
            type='primary'>
            创建会话
          </Button>
        </div>
      </ListActionWrapper>
    )
  }
)
