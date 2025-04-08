import { Modal, Table, Button } from '@/components'
import { message, Pagination } from 'antd'
import { toJS } from 'mobx'
import { observer, useLocalStore } from 'mobx-react'
import React, { useEffect } from 'react'
import styled from 'styled-components'
import { useStore } from '../store'
import { Edit } from './Edit'
import { UserPreview } from './Preview'
import { RoleList } from '@/domain/UserMG'
import { currentUser, sysConfig } from '@/domain'
import { State } from './State'
import ApproveStatus from '../../components/ApproveStatus'

const StyledLayout = styled.div`
  > .body {
    padding: 0;

    .userName {
      color: ${props => props.theme.primaryHighlightColor};
      cursor: pointer;

      .name {
        width: 95%;
        margin-left: 6px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }

    .roleName {
      display: flex;
      align-items: center;

      .name {
        width: 95%;
        margin-left: 6px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }

    .action {
      cursor: pointer;

      &:hover {
        color: ${props => props.theme.primaryHighlightColor};
      }
    }

    .disabled {
      color: #ccc;
      cursor: not-allowed;

      &:hover {
        color: #ccc;
      }
    }
  }
`

type Props = {
  width?: number
  height?: number
}
export const UserList = observer(function UserList({ width, height }: Props) {
  const store = useStore()
  const { selectedKeys, userList } = store
  const [fetch, loading] = store.getUserList()

  const state = useLocalStore(() => ({
    get dataSource() {
      const { userList } = store
      let list = userList || []
      return list.map(item => ({
        ...item,
        approve_status_str: item.approve_status_str
      }))
    }
  }))

  const { dataSource } = state

  useEffect(() => {
    fetch()
  }, [fetch])

  useEffect(() => {
    RoleList.fetch()
  }, [])

  function onPageChange(current, pageSize) {
    store.setPage(current, pageSize)
  }

  function onPageSizeChange(current, size) {
    store.setPage(current, size)
  }

  function isInternalOfSelf(rowData) {
    return currentUser.id === rowData.id || rowData.isInternal
  }

  return (
    <StyledLayout>
      <div className='body'>
        <Table
          props={{
            width,
            height,
            data: dataSource || [],
            loading,
            rowKey: 'id'
          }}
          rowSelection={{
            selectedRowKeys: selectedKeys,
            onSelectInvert() {
              store.setSelectedKeys([])
            },
            onSelectAll(keys) {
              let _selectedList = toJS(
                userList.filter(u => keys.includes(u.id))
              )

              let _finalSelectedlist = _selectedList.filter(
                u => !isInternalOfSelf(u)
              )

              if (selectedKeys.length >= _finalSelectedlist.length) {
                store.setSelectedKeys([])
              } else {
                const isSelectInternalOfSelf = _selectedList.some(u =>
                  isInternalOfSelf(u)
                )

                if (isSelectInternalOfSelf) {
                  message.error(`无法对自身或者内置用户进行选择操作`)
                }
                store.setSelectedKeys(_finalSelectedlist.map(u => u.id))
              }
            },
            onSelect(key, checked) {
              if (checked) {
                let user = toJS(userList.filter(u => key === u.id))[0]

                if (isInternalOfSelf(user)) {
                  message.error(`无法对自身或者内置用户进行选择操作`)
                  return
                }

                store.setSelectedKeys([...selectedKeys, key])
              } else {
                const index = selectedKeys.indexOf(key)
                selectedKeys.splice(index, 1)
                store.setSelectedKeys([...selectedKeys])
              }
            }
          }}
          columns={[
            {
              header: '登录名称',
              props: {
                width: 150,
                resizable: true
              },
              sorter: ({ sortType, sortKey }) => {
                store.setOrder(sortType === 'asc' ? true : false, sortKey)
              },
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
                        content: <UserPreview user={rowData} />
                      })
                    }>
                    <span className='name'>{rowData[dataKey]}</span>
                  </div>
                )
              }
            },
            {
              header: '用户名称',
              dataKey: 'name',
              props: {
                width: 150,
                resizable: true
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
                  return rowData.roleNames?.length === 0 ? (
                    <div>--</div>
                  ) : (
                    <div
                      className='roleName'
                      title={rowData.roleNames?.join(';')}>
                      <span className='name'>
                        {rowData.roleNames?.join(';')}
                      </span>
                    </div>
                  )
                }
              }
            },
            ...(sysConfig.enableThreeMembers
              ? [
                  {
                    header: '审批状态',
                    props: {
                      resizable: true,
                      width: 100
                    },
                    cell: {
                      props: {
                        dataKey: 'approve_status_str'
                      },
                      render: ({ rowData }) => {
                        const { approve_status_str } = rowData
                        return (
                          <ApproveStatus
                            targetId={rowData.id}
                            targetType={'USER'}
                            data={rowData}
                            callback={() => fetch()}>
                            {approve_status_str === '--' ? (
                              <span
                                style={{
                                  display: 'flex',
                                  justifyContent: 'center'
                                }}>
                                --
                              </span>
                            ) : (
                              <Button type='link'>{approve_status_str}</Button>
                            )}
                          </ApproveStatus>
                        )
                      }
                    }
                  }
                ]
              : []),
            {
              header: '状态',
              props: {
                width: 100,
                resizable: true
              },
              cell: {
                props: {
                  dataKey: 'enabled'
                },
                render: ({ rowData }) => <State user={rowData} />
              }
            },
            {
              props: {
                width: 150,
                resizable: true
              },
              header: '电话',
              dataKey: 'mobile'
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
                width: 100
              },
              header: '操作',
              cell: {
                props: {
                  dataKey: 'name'
                },
                render: ({ rowData, dataKey }) => (
                  <div
                    title={rowData[dataKey]}
                    className='userName'
                    onClick={() => {
                      if (isInternalOfSelf(rowData)) {
                        return
                      }
                      if (
                        rowData.approve_status === 0 &&
                        sysConfig.enableThreeMembers
                      ) {
                        message.warn(
                          `用户${rowData.name}有未完成的审批，请等待审批结束`
                        )
                        return
                      }
                      Modal.show({
                        title: '编辑用户',
                        bodyStyle: {
                          height: 710,
                          background: '#F0F5FD'
                        },
                        width: 1200,
                        footer: null,
                        content: ({ onCancel, onOk }) => (
                          <Edit
                            onCancel={onCancel}
                            onOk={() => {
                              fetch()
                              onOk()
                            }}
                            user={rowData}
                          />
                        )
                      })
                    }}>
                    <span
                      style={{
                        display: 'true'
                      }}
                      className={`action ${
                        isInternalOfSelf(rowData) ? 'disabled' : ''
                      }`}>
                      编辑
                    </span>
                  </div>
                )
              }
            }
          ]}
        />
      </div>
      <div style={{ textAlign: 'center', padding: 10 }}>
        <Pagination
          showSizeChanger
          showQuickJumper
          pageSize={store.page_size}
          current={store.page_index}
          total={store.total}
          onChange={onPageChange}
          onShowSizeChange={onPageSizeChange}
        />
      </div>
    </StyledLayout>
  )
})
