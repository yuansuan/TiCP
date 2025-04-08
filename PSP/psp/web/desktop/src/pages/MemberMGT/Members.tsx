/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Pagination } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Table } from '@/components'
import { StyledMembers } from './style'
import MemberAction from './MemberAction'
import { RowData } from './Type'
import { useStore } from '@/pages/MemberMGT/model'
import { env } from '@/domain'

export const Members = observer(function Members() {
  const store = useStore()

  const state = useLocalStore(() => ({
    get dataSource() {
      const searchKey = store.query.key
      const members = store.members
      let data: RowData[] = [...(members && members.list)].map(item => ({
        ...item,
        roles: item.role_list.map(item => item.name).join(','),
        create_time: item.create_time.toString(),
        update_time: item.update_time.toString(),
        last_login_time: item.last_login_time.toString(),
        department_name: item.department?.name || '--'
      }))

      if (searchKey) {
        data = data.filter(
          item =>
            item.phone.includes(searchKey) || item.real_name.includes(searchKey)
        )
      }

      return data
    }
  }))

  function onSort({ sortType, sortKey }) {
    store.setQuery({
      sortKey,
      sortType
    })
  }

  function onPageChange(current, pageSize) {
    store.setPage(current, pageSize)
  }

  return (
    <StyledMembers>
      <Table
        props={{
          data: state.dataSource,
          rowKey: 'user_id',
          autoHeight: true
        }}
        columns={[
          {
            header: '成员ID',
            props: {
              width: 150
            },
            dataKey: 'user_id'
          },
          {
            header: '姓名',
            props: {
              flexGrow: 1,
              minWidth: 100
            },
            dataKey: 'real_name'
          },
          {
            header: '手机号',
            props: {
              width: 150
            },
            dataKey: 'phone'
          },
          env.company.isOpenDepMgr && {
            header: '部门',
            props: {
              width: 200
            },
            dataKey: 'department_name'
          },
          {
            header: '加入时间',
            props: {
              width: 200
            },
            dataKey: 'create_time',
            sorter: onSort
          },
          {
            header: '最近登录时间',
            props: {
              width: 200
            },
            dataKey: 'last_login_time',
            sorter: onSort
          },
          {
            header: '角色',

            props: {
              width: 120,
              resizable: true
            },
            dataKey: 'roles'
          },
          {
            header: '消费限额(元)',
            props: {
              width: 120
            },
            cell: {
              props: {
                dataKey: 'consume_limit'
              },
              render: ({ rowData, dataKey }) => (
                <span>{rowData[dataKey] || '-'}</span>
              )
            }
          },
          {
            header: '操作',
            props: {
              width: 120
            },
            cell: {
              props: {
                dataKey: 'id'
              },
              render: ({ rowData }) => (
                <MemberAction
                  rowData={rowData}
                  refreshMembers={store.fetch}
                  departments={store.departmentList}
                />
              )
            }
          }
        ]}
      />

      <Pagination
        className='pagination'
        showQuickJumper
        showSizeChanger
        pageSize={store.page_size}
        current={store.page_index}
        total={store.total}
        onChange={onPageChange}
      />
    </StyledMembers>
  )
})
