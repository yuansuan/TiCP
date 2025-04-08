/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Pagination } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Table, Modal, Button } from '@/components'
import { StyledDepartments } from './style'
import { useStore } from '@/pages/DepartmentMGT/model'
import DepartmentAction from './DepartmentAction'
import { Tag, Tooltip } from 'antd'
import { EditingModal } from './DepartmentAction'

export const Departments = observer(function Members() {
  const store = useStore()

  const state = useLocalStore(() => ({
    get dataSource() {
      const data = store.list.list
      return data
    }
  }))

  function onPageChange(current, pageSize) {
    store.list.setPage(current, pageSize)
  }

  const renderUsers = users => {
    const TagList = color => (
      <span style={{ margin: 5 }}>
        {users.map(u => (
          <Tag color={color} key={u.user_id} title={`手机号：${u.phone}`}>
            {u.user_name || u.real_name || u.phone}
          </Tag>
        ))}
      </span>
    )

    return (
      <span>
        {users.length !== 0 ? (
          <>
            <Tooltip placement='left' title={TagList('')} color={'gray'}>
              <span>人员总数 {users.length} 个:</span>
            </Tooltip>
            {TagList('')}
          </>
        ) : (
          '-'
        )}
      </span>
    )
  }

  const preview = rowData =>
    Modal.show({
      title: '部门信息',
      content: ({ onCancel, onOk }) => (
        <EditingModal
          isPreview={true}
          isAdding={false}
          rowData={rowData}
          onCancel={onCancel}
          onOk={onOk}
          refresh={store.fetch}
        />
      ),
      footer: null
    })

  return (
    <StyledDepartments>
      <Table
        props={{
          data: state.dataSource,
          rowKey: 'id',
          autoHeight: true
        }}
        columns={[
          {
            header: '部门名称',
            props: {
              flexGrow: 1,
              fixed: 'left',
              align: 'center',
              minWidth: 100
            },
            cell: {
              props: {
                dataKey: 'name'
              },
              render: ({ rowData, dataKey }) => (
                <Button type='link' onClick={() => preview(rowData)}>
                  {rowData[dataKey]}
                </Button>
              )
            }
          },
          {
            header: '人员',
            props: {
              flexGrow: 1,
              minWidth: 500
            },
            cell: {
              props: {
                dataKey: 'users'
              },
              render: ({ rowData, dataKey }) => renderUsers(rowData[dataKey])
            }
          },
          {
            header: '部门备注',
            props: {
              width: 400
            },
            dataKey: 'remark'
          },
          {
            header: '操作',
            props: {
              width: 120,
              fixed: 'right',
              align: 'center'
            },
            cell: {
              props: {
                dataKey: 'id',
                align: 'center'
              },
              render: ({ rowData }) => (
                <DepartmentAction rowData={rowData} refresh={store.fetch} />
              )
            }
          }
        ]}
      />

      <Pagination
        className='pagination'
        showQuickJumper
        showSizeChanger
        pageSize={store.list.page_size}
        current={store.list.page_index}
        total={store.list.total}
        onChange={onPageChange}
      />
    </StyledDepartments>
  )
})
