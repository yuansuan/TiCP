/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Form, Select, DatePicker } from 'antd'
import { useStore } from './store'
import { Button, Icon } from '@/components'
import { Http } from '@/utils'
import { billUserServer } from '@/server'
import moment from 'moment'
import { env } from '@/domain'
import { NewJobLayout } from './List'

const StyledLayout = styled.div`
  display: flex;
  flex-direction: column;
  > .top {
    display: flex;
    margin-bottom: 10px;
    > .left {
      display: inline-flex;
      align-items: center;
      justify-content: center;
    }
    > .right {
      margin-left: auto;
    }
  }
`

export const Toolbar = observer(function Toolbar() {
  const store = useStore()
  const [form] = Form.useForm()
  const state = useLocalStore(() => ({
    mdseList: [],
    mdseLoading: false,
    setMdseLoading(flag) {
      this.mdseLoading = flag
    },
    setMdseList(list) {
      this.mdseList = list
    }
  }))
  function onFinish(values) {
    store.update({
      queryKey: {
        types: values?.types,
        merchandise_id: values?.merchandise_id,
        billing_month: values?.month?.format('YYYY-MM')
      },
      pageIndex: 1,
      pageSize: 10
    })
  }

  useEffect(() => {
    Http.get('/company/merchandise', {
      params: {
        out_resource_types: store.queryKey.types,
        company_id: !env.isPersonal ? env?.company?.id : '1'
      }
    }).then(({ data }) => {
      state.setMdseList(data)
    })
  }, [store.queryKey.types])

  function submit() {
    form.submit()
  }

  async function exportFile() {
    await billUserServer
      .export({
        ...store.queryKey
      })
      .then(response => {
        const url = window.URL.createObjectURL(new Blob([response.data]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute(
          'download',
          `账单详情_${moment().format('YYYYMMDD')}.xlsx`
        )
        document.body.appendChild(link)
        link.click()
      })
  }

  function valuesChange(changedValues, allValues) {
    changedValues.types &&
      (((store.queryKey.merchandise_id = ''),
      (store.queryKey.types = changedValues.types)),
      form.setFieldsValue({ merchandise_id: '' }))
  }

  return (
    <StyledLayout>
      <div className='top'>
        <div className='left'>
          <Form
            form={form}
            layout='inline'
            onFinish={onFinish}
            onValuesChange={(changedValues, allValues) => {
              valuesChange(changedValues, allValues)
            }}
            initialValues={{
              ...store.queryKey,
              month: null
            }}>
            <Form.Item label='账单类型' name='types'>
              <Select
                showSearch
                allowClear
                mode='multiple'
                placeholder='请选择账单类型'
                filterOption={(input, option) =>
                  option.children.toLowerCase().indexOf(input.toLowerCase()) >=
                  0
                }
                style={{ width: 240 }}>
                <Select.Option value={1}>计算作业</Select.Option>
                <Select.Option value={5}>3D云应用-软件</Select.Option>
                <Select.Option value={6}>3D云应用-硬件</Select.Option>
                <Select.Option value={7}>
                  <NewJobLayout>
                    计算作业
                    <Icon type='nys-new' className='icon' />
                  </NewJobLayout>
                </Select.Option>
                <Select.Option value={103}>3D云应用套餐</Select.Option>
              </Select>
            </Form.Item>
            <Form.Item label='应用名称' name='merchandise_id'>
              <Select
                showSearch
                value={store.queryKey.merchandise_id}
                placeholder='请选择应用名称'
                loading={state.mdseLoading}
                filterOption={(input, option) =>
                  option.children.toLowerCase().indexOf(input.toLowerCase()) >=
                  0
                }
                style={{ width: 240 }}>
                <Select.Option value=''>全部</Select.Option>
                {state.mdseList.map(merchandise => (
                  <Select.Option value={merchandise?.id} key={merchandise?.id}>
                    {merchandise?.name}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
            <Form.Item label='作业提交时间' name='month'>
              <DatePicker.MonthPicker />
            </Form.Item>
            <Form.Item>
              <Button loading={store.loading} onClick={submit}>
                查询
              </Button>
            </Form.Item>
            <Form.Item>
              <Button
                loading={store.loading}
                onClick={() => {
                  form.setFieldsValue({
                    merchandise_id: '',
                    types: [],
                    month: null
                  })
                  state.setMdseList([])
                  store.update({
                    queryKey: {
                      types: [],
                      merchandise_id: '',
                      billing_month: ''
                    },
                    pageIndex: 1,
                    pageSize: 10
                  })
                }}>
                重置
              </Button>
            </Form.Item>
          </Form>
        </div>
        <div className='right'>
          共{store.model?.page_ctx?.total}项条目,消费
          {(store.model?.total_amount / 100000).toFixed(2)}元,退款总额
          {(store.model?.total_refund_amount / 100000).toFixed(2)}元 |
          <Button
            disabled={store.model.list.length <= 0}
            style={{ padding: '5px' }}
            type='link'
            onClick={exportFile}>
            导出
          </Button>
        </div>
      </div>
    </StyledLayout>
  )
})
