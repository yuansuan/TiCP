import React, { useState, useEffect } from 'react'
import { useLocalStore } from 'mobx-react-lite'
import { observer } from 'mobx-react'
import { RoleList, UserList } from '@/domain/UserMG'
import UserPreview from '../UserPreview'
import { ListQuery } from '../../utils'
import Operators from './Operators'
import State from './State'
import { PanelWrapper } from './style'
import { formatDateFromMilliSecWithTimeZone } from '@/utils/formatter'
import { Table, Modal } from '@/components'
import currentUser from '@/domain/User'
import { sysConfig } from '@/domain'

interface IProps {
  width: number
  loading: boolean
  height: number

  listQuery: ListQuery
  store: typeof UserList
  updateTotal: (number) => void
}

export default observer(function UserPanel(props: IProps) {
  const { width, loading, listQuery, store, updateTotal, height } = props
  const state = useLocalStore(() => ({
    get dataSource() {
      return store.enabledUsers.map(item => ({
        ...item,
        approve_status_str: item.approve_status_str,
        user_id: item.user_id,
        roles: item.roles.map(r =>
          RoleList.list.get(r) ? RoleList.list.get(r).name : '-'
        )
      }))
    }
  }))

  useEffect(() => {
    updateTotal(store.totalEnabledUsers)
  }, [store.enabledUsers])

  useEffect(() => {
    const { page, pageSize, query } = listQuery
    store.filter = {
      query,
      page: {
        index: page,
        size: pageSize
      }
    }
    store.fetch()
  }, [listQuery.query, listQuery.page, listQuery.pageSize])
  const columns = () => {
    let all = [
      {
        props:
          currentUser.authType !== 'local'
            ? {
                resizable: true,
                width: 200
              }
            : { flexGrow: 1, minWidth: 200 },
        header: '登录名称',
        cell: {
          props: {
            dataKey: 'name'
          },
          render: ({ rowData, dataKey }) => (
            <div
              title={rowData[dataKey]}
              className='userName'
              onClick={() =>
                Modal.show({
                  title: '预览用户',
                  bodyStyle: {
                    height: 710,
                    overflow: 'auto',
                    background: '#F0F5FD'
                  },
                  width: 900,
                  footer: null,
                  content: () => (
                    <UserPreview
                      user={UserList.enabledUserMap.get(rowData.id)}
                    />
                  )
                })
              }>
              <span className='name'>{rowData[dataKey]}</span>
            </div>
          )
        }
      },
      {
        props: {
          width: 200,
          resizable: true
        },
        header: '用户角色',
        cell: {
          props: {
            dataKey: 'roles'
          },
          render: ({ rowData }) => {
            return rowData.roles.length === 0 ? (
              <div>--</div>
            ) : (
              <div className='roleName' title={rowData.roles?.join(';')}>
                <span className='name'>{rowData.roles?.join(';')}</span>
              </div>
            )
          }
        }
      },

      {
        props: {
          width: 200,
          resizable: true
        },
        header: '邮箱',
        dataKey: 'email'
      },
      {
        props: {
          width: 200,
          resizable: true
        },
        header: '电话',
        dataKey: 'mobile'
      },
      {
        props: {
          width: 200
        },
        header: '创建时间',
        cell: {
          dataKey: 'created_at',
          render: ({ rowData }) =>
            formatDateFromMilliSecWithTimeZone(rowData['created_at'])
        }
      },
      // {
      //   props: {
      //     fixed: 'right',
      //     width: 100
      //   },
      //   header: '审批状态',
      //   cell: {
      //     dataKey: 'approve_status_str',
      //     render: ({ rowData }) => {
      //       const { approve_status_str } = rowData
      //       return (
      //         <ApproveStatus
      //           targetId={rowData.id}
      //           targetType={'USER'}
      //           data={rowData}
      //           callback={() => UserList.fetch()}>
      //           {approve_status_str === '--' ? (
      //             <span
      //               style={{
      //                 display: 'flex',
      //                 justifyContent: 'center'
      //               }}>
      //               --
      //             </span>
      //           ) : (
      //             <Button type='link'>{approve_status_str}</Button>
      //           )}
      //         </ApproveStatus>
      //       )
      //     }
      //   }
      // },
      {
        props: {
          fixed: 'right',
          width: 100
        },
        header: '状态',
        cell: {
          dataKey: 'enabled',
          render: ({ rowData }) => <State rowData={rowData} />
        }
      },
      {
        props: {
          fixed: 'right',
          minWidth: 120,
          flexGrow: 1
        },
        header: '操作',
        cell: {
          render: ({ rowData }) => <Operators rowData={rowData} />
        }
      }
    ]

    return all
  }

  return (
    <PanelWrapper>
      <Table
        props={{
          height: height,
          data: state.dataSource || [],
          rowKey: 'id',
          loading
        }}
        columns={columns() as any}
      />
    </PanelWrapper>
  )
})
