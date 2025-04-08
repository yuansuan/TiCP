/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Table } from '@/components'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useStore } from '../store'
import { Actions } from './Actions'

const StyledLayout = styled.div`
  .active {
    &::before {
      content: '';
      display: inline-block;
      width: 12px;
      height: 12px;
      border-radius: 12px;
      background: #52c41a;
      border: 2px solid #d7f9c7;
      margin-right: 8px;
    }
  }
  .inactive {
    &::before {
      content: '';
      display: inline-block;
      width: 12px;
      height: 12px;
      border-radius: 12px;
      background: #c5c5c5;
      border: 2px solid#E6E4E4;
      margin-right: 8px;
    }
  }
`

export const List = observer(function List() {
  const store = useStore()
  const { dataSource } = useLocalStore(() => ({
    get dataSource() {
      return store.model.list.map(item => ({
        ...item,
      }))
    },
  }))

  return (
    <StyledLayout>
      <Table
        columns={[
          {
            header: '软件名称',
            dataKey: 'app_name',
            props: {
              flexGrow: 1,
            },
          },
          {
            header: '版本',
            dataKey: 'version',
            props: {
              flexGrow: 1,
            },
          },
          {
            header: '状态',
            props: {
              flexGrow: 1,
            },
            cell: {
              props: {
                dataKey: 'active',
              },
              render({ rowData, dataKey }) {
                return (
                  <div className={rowData[dataKey] ? 'active' : 'inactive'}>
                    {rowData[dataKey] ? '已激活' : '未激活'}
                  </div>
                )
              },
            },
          },
          {
            header: '操作',
            cell: {
              props: {
                dataKey: 'merchandise_id',
              },
              render({ rowData, dataKey }) {
                return <Actions merchandiseId={rowData[dataKey]} />
              },
            },
          },
        ]}
        props={{
          autoHeight: true,
          data: dataSource,
          loading: store.fetching,
        }}
      />
    </StyledLayout>
  )
})
