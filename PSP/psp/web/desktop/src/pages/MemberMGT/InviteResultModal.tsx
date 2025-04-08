/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useMemo } from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { Table } from '@/components'

const StyledLayout = styled.div`
  overflow: auto;

  > .stat {
    margin: 0 0 12px;
    text-align: right;
  }

  .status {
    padding-left: 14px;
    position: relative;

    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 44%;
      width: 10px;
      height: 10px;
      background-color: #52c41a;
      border: 2px solid #d7f9c7;
      border-radius: 50%;
    }

    &.warning {
      &::before {
        background-color: #ffa726;
        border-color: #f9d9d9;
      }
    }

    &.error {
      &::before {
        background-color: #f5222d;
        border-color: #f9d9d9;
      }
    }
  }
`

type Props = {
  list: any[]
}

export const InviteResultModal = observer(function InviteResultModal({
  list
}: Props) {
  const dataSource = useMemo(
    () =>
      list.map(item => ({
        ...item,
        message: {
          170011: '用户已邀请',
          170012: '手机号格式错误',
          170013: '用户已加入其他企业',
          2: '未知错误',
          0: '成功'
        }[item.code]
      })),
    []
  )

  return (
    <StyledLayout>
      <div className='stat'>
        共计邀请{dataSource.length}个用户，成功
        {dataSource.filter(item => item.code === 0).length}个，用户已存在
        {dataSource.filter(item => item.code === 170011).length}个，手机号错误
        {dataSource.filter(item => item.code === 170012).length}个
      </div>
      <Table
        props={{
          data: dataSource,
          height: 250
        }}
        columns={[
          {
            header: '手机号',
            dataKey: 'phone',
            props: {
              flexGrow: 1
            }
          },
          {
            header: '状态',
            props: {
              flexGrow: 1
            },
            cell: {
              props: {
                dataKey: 'message'
              },
              render({ rowData, dataKey }) {
                return (
                  <div
                    className={`status ${
                      {
                        170012: 'error',
                        2: 'error',
                        170011: 'warning',
                        0: 'success',
                        170013: 'error'
                      }[rowData['code']]
                    }`}>
                    {rowData[dataKey]}
                  </div>
                )
              }
            }
          }
        ]}
      />
    </StyledLayout>
  )
})
