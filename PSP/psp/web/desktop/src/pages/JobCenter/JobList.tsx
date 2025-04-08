/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { formatAmount } from '@/utils'
import { useStore } from './store'
import { Table } from '@/components'
import { JobStatus } from '@/components'
import { Pagination } from 'antd'
import { JobName } from './JobName'

const StyledLayout = styled.div`
  .pagination {
    text-align: right;
    margin: 20px;
  }
`

export const JobList = observer(function JobCenterList() {
  const store = useStore()
  const state = useLocalStore(() => ({
    get dataSource() {
      return store.model.list.map(item => {
        // 核时
        function getCoreTime() {
          if (!item.runtime.cpu_time || !item.runtime.resource_assign?.cpus) {
            return '--'
          }

          return item.runtime
            ? (
                ((item.runtime.cpu_time || 1) / 3600) *
                (item.runtime.resource_assign?.cpus || 1)
              ).toFixed(2)
            : '--'
        }

        // 运行时长
        function getDisplayRunTime(j) {
          if (!j.runtime?.cpu_time || j.runtime?.cup_time === 0) return null

          const hour = Math.floor(j.runtime.cpu_time / 3600)
            .toString()
            .padStart(2, '0')
          const minute = Math.floor((j.runtime.cpu_time % 3600) / 60)
            .toString()
            .padStart(2, '0')
          const second = (j.runtime.cpu_time % 60).toString().padStart(2, '0')
          return `${hour}:${minute}:${second}`
        }

        return {
          ...item,
          deleteable: item.deleteable || '--',
          cancelable: item.cancelable || '--',
          downloadable: item.downloadable || '--',
          residualVisible: !!item.runtime?.have_residual || '--',
          displayCpus: item.resource_usage.cpus || 0,
          displayRunTime: getDisplayRunTime(item) || '--',
          displayCreateTime: item.create_time.toString() || '--',
          displayStartTime: item.runtime?.start_time?.toString() || '--',
          displayEndTime: item.runtime?.end_time?.toString() || '--',
          displayState:
            {
              1: '运行中',
              2: '已成功',
              3: '出错',
              4: '取消',
              7: '排队中',
            }[item.display_state] || '--',
          displayAmount:
            item.amount === undefined || item.amount === 0
              ? '--'
              : formatAmount(item.amount),
          displayCoreTime: getCoreTime() || '--',
        }
      })
    },
  }))

  function onPageChange(index, size) {
    store.setPageIndex(index)
    store.setPageSize(size)
  }

  return (
    <StyledLayout>
      <Table
        tableId='job_center_list_table'
        defaultConfig={[
          { key: 'name', active: true },
          { key: 'id', active: true },
          { key: 'user_name', active: true },
          { key: 'tier_name', active: true },
          { key: 'app_name', active: true },
          { key: 'app_version', active: true },
          { key: 'app_type', active: true },
          { key: 'displayCpus', active: true },
          { key: 'displayRunTime', active: true },
          { key: 'displayCreateTime', active: true },
          { key: 'displayStartTime', active: true },
          { key: 'displayEndTime', active: true },
          { key: 'displayState', active: true },
          { key: 'displayCoreTime', active: true },
          { key: 'displayAmount', active: true },
        ]}
        props={{
          data: state.dataSource,
          rowKey: 'id',
          autoHeight: true,
          loading: store.loading,
        }}
        columns={[
          {
            header: '名称',
            props: {
              resizable: true,
              width: 200,
              fixed: 'left',
            },
            dataKey: 'name',
            cell: {
              props: {
                dataKey: 'name',
              },
              render: ({ rowData }) => (
                <JobName
                  id={rowData['id']}
                  project_id={rowData['project_id']}
                  name={rowData['name']}
                />
              ),
            },
          },
          {
            header: '作业编号',
            props: {
              width: 150,
              fixed: 'left',
            },
            dataKey: 'id',
            cell: {
              props: {
                dataKey: 'id',
              },
            },
          },
          {
            header: '创建人',
            props: {
              width: 200,
            },
            dataKey: 'user_name',
            cell: {
              render: ({ rowData }) => {
                return <div>{rowData.user_name}</div>
              },
            },
          },
          {
            header: '算力资源',
            props: {
              width: 120,
            },
            cell: {
              props: {
                dataKey: 'tier_name',
              },
            },
          },
          {
            header: '软件名称',
            props: {
              width: 120,
            },
            dataKey: 'app_name',
          },
          {
            header: '软件版本',
            props: {
              width: 120,
            },
            dataKey: 'app_version',
          },
          {
            header: '软件类型',
            props: {
              width: 120,
            },
            dataKey: 'app_type',
          },
          {
            header: '核数',
            props: {
              width: 100,
            },
            dataKey: 'displayCores',
          },
          {
            header: '提交时间',
            props: {
              width: 200,
            },
            dataKey: 'displayCreateTime',
          },
          {
            header: '运行时间',
            props: {
              width: 200,
            },
            dataKey: 'displayStartTime',
          },
          {
            header: '结束时间',
            props: {
              width: 200,
            },
            dataKey: 'displayEndTime',
          },
          {
            header: '运行时长',
            props: {
              width: 140,
            },
            dataKey: 'displayRunTime',
          },
          {
            header: '作业状态',
            props: {
              width: 120,
            },
            dataKey: 'displayState',
            cell: {
              props: {
                dataKey: 'displayState',
              },
              render: ({ rowData }) => {
                return <JobStatus showDropDown={false} job={rowData} />
              },
            },
          },
          {
            header: '用量(核时)',
            props: {
              width: 120,
              fixed: 'right',
            },
            dataKey: 'displayCoreTime',
          },
          {
            header: '消费金额(元)',
            props: {
              width: 150,
              fixed: 'right',
            },
            dataKey: 'displayAmount',
          },
        ]}
      />

      <Pagination
        className='pagination'
        showQuickJumper
        showSizeChanger
        pageSize={store.pageSize}
        pageSizeOptions={['10', '25', '50']}
        current={store.pageIndex}
        total={store.model.page_ctx.total}
        onChange={onPageChange}
      />
    </StyledLayout>
  )
})
