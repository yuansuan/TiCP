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
import { jobServer } from '@/server'
import { currentUser } from '@/domain'
import { Selector } from './Selector'
import { runInAction } from 'mobx'
import { JOB_DISPLAY_STATE } from '@/constant'
import { DatePicker_FORMAT, DatePicker_SHOWTIME_FORMAT } from '@/constant'
import { sysConfig } from '@/domain'
import { getComputeType } from '@/utils'

const StyledLayout = styled.div`
  display: flex;
  padding: 20px 20px 0px;

  > .left {
    flex: 1;

    .ant-form-inline .ant-form-item {
      margin-bottom: 16px;

      .ant-form-item-label {
        width: 80px;
        text-align: right;
      }
    }
  }

  > .right {
    width: 60px;
  }

  margin-bottom: -16px;
`

const { useForm } = Form
const { RangePicker } = DatePicker

type Props = {
  showJobSetName?: boolean
}

export const Filter = observer(function Filter(props: Props) {
  const store = useStore()
  const [form] = useForm()
  const { isPersonalJobManager } = currentUser
  const state = useLocalStore(() => ({
    loading: false,
    setLoading(flag) {
      this.loading = flag
    },
    filters: {
      apps: [],
      users: [],
      states: [],
      queues: [],
      projects: [],
      job_types: [],
      job_sets: []
    },

    setFilters(filters) {
      Object.assign(this.filters, filters)
    }
  }))

  useEffect(() => {
    state.setLoading(true)
    jobServer
      .getAllFilters()
      .then(
        ({ app_names, user_names, queue_names, projects, job_set_names }) => {
          state.setFilters({
            apps:
              app_names?.filter(Boolean).map(v => ({ key: v, name: v })) || [],
            users:
              user_names?.filter(Boolean).map(v => ({ key: v, name: v })) || [],
            queues:
              queue_names?.filter(Boolean).map(v => ({ key: v, name: v })) ||
              [],
            projects:
              projects
                ?.filter(Boolean)
                .map(v => ({ key: v.id, name: v.name })) || [],
            states: Object.entries(JOB_DISPLAY_STATE).map(item => ({
              key: item[0],
              name: item[1]
            })),
            job_sets:
              job_set_names?.filter(Boolean).map(v => ({ key: v, name: v })) ||
              []
          })
        }
      )
      .finally(() => {
        state.setLoading(false)
      })
  }, [])

  const handleDropdownVisibleChange = async (open, optType) => {
    switch (optType) {
      case 'app_names':
        if (open) {
          const res = await jobServer.getAppNames()
          if (res?.data?.length > 0) {
            state.setFilters({
              apps:
                res.data?.filter(Boolean).map(v => ({ key: v, name: v })) || []
            })
          }
        }
        break
      case 'user_names':
        if (open) {
          const res = await jobServer.getUserNames()
          if (res?.data?.length > 0) {
            state.setFilters({
              users:
                res.data.filter(Boolean).map(v => ({ key: v, name: v })) || []
            })
          }
        }
        break
      case 'queues':
        if (open) {
          const res = await jobServer.getQueueNames()
          if (res?.data?.length > 0) {
            state.setFilters({
              queues:
                res.data?.filter(Boolean).map(v => ({ key: v, name: v })) || []
            })
          }
        }
        break
      case 'projects':
        if (open) {
          const res = await jobServer.getProjects()
          if (res?.data?.projects?.length > 0) {
            state.setFilters({
              projects:
                res.data?.projects
                  .filter(Boolean)
                  .map(v => ({ key: v.id, name: v.name })) || []
            })
          }
        }
        break
      case 'job_set_names':
        if (open) {
          const res = await jobServer.getJobSetNames()
          if (res?.data?.length > 0) {
            state.setFilters({
              job_sets:
                res.data
                  ?.filter(Boolean)
                  .map(v => ({ key: v.id, name: v.name })) || []
            })
          }
        }
        break
      default:
        break
    }
  }

  // fullfill form with initial query
  useEffect(() => {
    const { start_time, end_time } = store.query
    form.setFieldsValue({
      ...store.query,
      submit_time_range: [
        start_time ? moment.unix(+start_time) : undefined,
        end_time ? moment.unix(+end_time) : undefined
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
    let newStates = values['states']
    if (
      values['states'].length > 0 &&
      values['states'].includes('Failed') &&
      !values['states'].includes('BurstFailed')
    ) {
      newStates = values['states'].concat('BurstFailed')
    }
    runInAction(() => {
      store.setQuery({
        job_id: values['job_id'],
        job_name: values['job_name'],
        user_names: values['user_names'],
        app_names: values['app_names'],
        job_types: values['job_types'],
        states: newStates,
        queues: values['queues'],
        project_ids: values['project_ids'],
        // is_admin: currentUser.hasSysMgrPerm, // 后端同学自己处理了
        start_time: submitTime && submitTime[0]?.unix(),
        end_time: submitTime && submitTime[1]?.unix(),
        job_set_id: values['job_set_id'],
        job_set_names: values['job_set_names']
      })
      store.setPageIndex(1)
    })
  }

  function reset() {
    form.setFieldsValue({
      ...initialQuery,
      submit_time_range: [undefined, undefined]
    })
    submit()
  }
  return (
    <StyledLayout className='job-search-bar'>
      <div className='left'>
        <Form form={form} layout='inline' onFinish={onFinish}>
          <Form.Item label='作业编号' name='job_id'>
            <Input
              style={{ width: 220 }}
              placeholder='输入作业编号搜索'
              onChange={debounceSubmit}
            />
          </Form.Item>
          <Form.Item label='作业名称' name='job_name'>
            <Input
              style={{ width: 220 }}
              placeholder='输入作业名称搜索'
              onChange={debounceSubmit}
            />
          </Form.Item>
          {/* 个人作业管理不展示 */}
          {!isPersonalJobManager && (
            <Form.Item label='用户名称' name='user_names'>
              <Selector
                loading={state.loading}
                filters={state.filters.users}
                onChange={submit}
                onDropdownVisibleChange={open =>
                  handleDropdownVisibleChange(open, 'user_names')
                }
              />
            </Form.Item>
          )}
          <Form.Item label='应用名称' name='app_names'>
            <Selector
              loading={state.loading}
              filters={state.filters.apps}
              onChange={submit}
              onDropdownVisibleChange={open =>
                handleDropdownVisibleChange(open, 'app_names')
              }
            />
          </Form.Item>
          <Form.Item label='队列名称' name='queues'>
            <Selector
              loading={state.loading}
              filters={state.filters.queues}
              onChange={submit}
              onDropdownVisibleChange={open =>
                handleDropdownVisibleChange(open, 'queues')
              }
            />
          </Form.Item>
          <Form.Item label='项目名称' name='project_ids'>
            <Selector
              loading={state.loading}
              filters={state.filters.projects}
              onChange={submit}
              onDropdownVisibleChange={open =>
                handleDropdownVisibleChange(open, 'projects')
              }
            />
          </Form.Item>
          <Form.Item label='计算状态' name='states'>
            <Selector
              loading={state.loading}
              filters={state.filters.states}
              onChange={submit}
            />
          </Form.Item>
          <Form.Item label='作业类型' name='job_types'>
            <Selector
              loading={state.loading}
              filters={[
                { key: 'local', name: getComputeType('local') },
                {
                  key: 'cloud',
                  name: getComputeType('cloud') || '云端'
                }
              ]}
              onChange={submit}
            />
          </Form.Item>
          <Form.Item label='作业集编号' name='job_set_id'>
            <Input
              style={{ width: 220 }}
              placeholder='输入作业集编号搜索'
              onChange={debounceSubmit}
            />
          </Form.Item>
          <Form.Item label='作业集名称' name='job_set_names'>
            <Selector
              loading={state.loading}
              filters={state.filters.job_sets}
              onChange={submit}
              onDropdownVisibleChange={open =>
                handleDropdownVisibleChange(open, 'job_set_names')
              }
            />
          </Form.Item>
          <Form.Item label='提交时间' name='submit_time_range'>
            <RangePicker
              ranges={GeneralDatePickerRange}
              style={{ width: 400 }}
              format={DatePicker_FORMAT}
              showTime={{
                format: DatePicker_SHOWTIME_FORMAT,
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
