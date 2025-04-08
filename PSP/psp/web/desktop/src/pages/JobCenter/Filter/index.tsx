/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useCallback, useEffect } from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Form, Input, DatePicker } from 'antd'
import { useStore, initialQuery } from '../store'
import debounce from 'lodash/debounce'
import moment from 'moment'
import { Button } from '@/components'
import { GeneralDatePickerRange } from '@/constant'
import { jobCenterServer } from '@/server'
import { Selector } from './Selector'
import { runInAction } from 'mobx'

const StyledLayout = styled.div`
  display: flex;
  padding: 20px 20px 0px;

  > .left {
    flex: 1;

    .ant-form-inline .ant-form-item {
      margin-bottom: 16px;
    }
  }

  > .right {
    width: 60px;
  }
`

const { useForm } = Form
const { RangePicker } = DatePicker

export const Filter = observer(function Filter() {
  const store = useStore()
  const [form] = useForm()
  const state = useLocalStore(() => ({
    loading: false,
    setLoading(flag) {
      this.loading = flag
    },
    filters: {
      user: [],
      state: []
    },
    setFilters(filters) {
      this.filters = filters
    }
  }))

  useEffect(() => {
    state.setLoading(true)
    jobCenterServer
      .getFilters()
      .then(({ data }) => {
        state.setFilters({
          user: data.user_filters,
          state: data.state_filters.map(item => ({
            key: item.key,
            name: {
              1: '运行中',
              2: '已成功',
              3: '出错',
              4: '取消',
              7: '排队中'
            }[item.name]
          }))
        })
      })
      .finally(() => {
        state.setLoading(false)
      })
  }, [])

  // fullfill form with initial query
  useEffect(() => {
    const { start_seconds, end_seconds } = store.query
    form.setFieldsValue({
      ...store.query,
      submit_time_range: [
        start_seconds ? moment.unix(+start_seconds) : undefined,
        end_seconds ? moment.unix(+end_seconds) : undefined
      ]
    })
  }, [])

  function submit() {
    form.submit()
  }

  const debounceSubmit = useCallback(
    debounce(function () {
      form.submit()
    }, 300),
    []
  )

  function onFinish(values) {
    const submitTime = values['submit_time_range']

    runInAction(() => {
      store.setQuery({
        job_id: values['job_id'].trim(),
        user_filters: values['user_filters'],
        state_filters: values['state_filters'],
        start_seconds: submitTime && submitTime[0]?.unix(),
        end_seconds: submitTime && submitTime[1]?.unix()
      })
      store.setPageIndex(1)
    })
  }

  function reset() {
    form.setFieldsValue({
      ...initialQuery,
      submit_time_range: [
        initialQuery.start_seconds
          ? moment.unix(+initialQuery.start_seconds)
          : undefined,
        initialQuery.end_seconds
          ? moment.unix(+initialQuery.end_seconds)
          : undefined
      ]
    })
    submit()
  }

  return (
    <StyledLayout>
      <div className='left'>
        <Form
          form={form}
          layout='inline'
          labelCol={{ span: 6 }}
          wrapperCol={{ span: 14 }}
          onFinish={onFinish}>
          <Form.Item label='作业编号' name='job_id'>
            <Input
              style={{ width: 200 }}
              placeholder='输入作业编号搜索'
              onChange={debounceSubmit}
            />
          </Form.Item>
          <Form.Item label='创建人' name='user_filters'>
            <Selector
              loading={state.loading}
              filters={state.filters.user}
              onChange={submit}
            />
          </Form.Item>
          <Form.Item label='作业状态' name='state_filters'>
            <Selector
              loading={state.loading}
              filters={state.filters.state}
              onChange={submit}
            />
          </Form.Item>
          <Form.Item label='提交时间' name='submit_time_range'>
            <RangePicker
              style={{ width: 360 }}
              ranges={GeneralDatePickerRange}
              showTime={{
                defaultValue: [
                  moment('00:00:00', 'HH:mm:ss'),
                  moment('23:59:59', 'HH:mm:ss')
                ]
              }}
              onChange={submit}
            />
          </Form.Item>
        </Form>
      </div>
      <div className='right'>
        <Button type='default' onClick={reset}>
          重置
        </Button>
      </div>
    </StyledLayout>
  )
})
