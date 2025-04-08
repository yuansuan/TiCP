/*
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer } from 'mobx-react-lite'
import { SearchSelect, Button } from '@/components'
import { useStore } from '../store'
import { Form, Row, Col } from 'antd'
import { SESSION_STATUS_MAP } from '@/domain/Vis'
import styled from 'styled-components'

const Style = styled.div`
  .ant-select-selector {
    padding-right: 24px;
  }
  .actions {
    button {
      margin-right: 10px;
    }
  }
`

export const Toolbar = observer(function Toolbar() {
  const store = useStore()

  const hardwareList = store.hardware.hardwareList?.map(item => ({
    key: item?.id,
    name: item?.name
  }))
  const softwareList = store.software.softwareList?.map(item => ({
    key: item?.id,
    name: item?.name
  }))

  
  function onFilter() {
    store.setSessionPageIndex(1)
    store.fetchSessionList()
  }

  const [form] = Form.useForm()

  function reset() {
    const resetState = {
      statuses: [],
      user_id: null,
      hardware_ids: [],
      software_ids: [],
      project_ids: []
    }
    form.setFieldsValue(resetState)
    store.changeSearchItems({}, resetState)
    store.setSessionPageIndex(1)
    store.setSessionPageSize(10)
    store.fetchSessionList()
  }

  function valuesChange(changedValues, allValues) {
    store.changeSearchItems(changedValues, allValues)
    onFilter()
    return true
  }

  return (
    <Style>
      <Form
        form={form}
        name='advanced_search'
        className='ant-advanced-search-form'
        labelCol={{ span: 8 }}
        initialValues={{
          statuses: store?.statuses,
          hardware_ids: store.hardware_ids,
          software_ids: store.software_ids,
          project_ids: store.project_ids
        }}
        onValuesChange={(changedValues, allValues) => {
          valuesChange(changedValues, allValues)
        }}>
        <Row gutter={22}>
          <Col span={5}>
            <Form.Item name='statuses' label='会话状态'>
              <SearchSelect
                mode="multiple"
                placeholder='所有会话状态'
                showArrow={true}
                allowClear={true}
                options={Object.entries(SESSION_STATUS_MAP).map(
                  ([key, name]) => ({
                    key,
                    name
                  })
                )}
              />
            </Form.Item>
          </Col>
          <Col span={5}>
            <Form.Item name='hardware_ids' label='实例名称'>
              <SearchSelect
                mode="multiple"
                placeholder='所有实例'
                showArrow={true}
                allowClear={true}
                options={hardwareList}
              />
            </Form.Item>
          </Col>
          <Col span={5}>
            <Form.Item name='software_ids' label='镜像名称'>
              <SearchSelect
                mode="multiple"
                placeholder='所有镜像'
                showArrow={true}
                allowClear={true}
                options={softwareList}
              />
            </Form.Item>
          </Col>
          <Col span={5}>
            <Form.Item name='project_ids' label='项目名称'>
              <SearchSelect
                mode="multiple"
                placeholder='所有项目'
                showArrow={true}
                allowClear={true}
                options={store.projects}
              />
            </Form.Item>
          </Col>
          <Col span={2}>
            <div className='actions'>
              {/* <Button onClick={onFilter}>查询</Button> */}
              <Button onClick={() => reset()}>重置</Button>
            </div>
          </Col>
        </Row>
        <Row gutter={24}>
          {/* <Col span={7}>
            <Form.Item name='user_id' label='创建者'>
              <SearchSelect
                placeholder='请选择'
                showArrow={true}
                allowClear={true}
                caseSensitive={false}
                loading={userListLoading}
                options={userList}
              />
            </Form.Item>
          </Col> */}
        </Row>
      </Form>
    </Style>
  )
})
