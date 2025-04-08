/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Divider, message, Button } from 'antd'
import { Table, Modal, Icon } from '@/components'
import { useLocalStore, observer } from 'mobx-react-lite'
import { useStore } from '../store'
import { formatUnixTime } from '@/utils'
import { formatByte } from '@/utils/Validator'
import { Edit } from './Edit'
import { Name } from './Name'
import { EDITABLE_SIZE } from '@/constant'

const StyledLayout = styled.div`
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;

  > .list_body {
    flex: 1;
    padding: 0 20px;
    .rs-table-row:hover,
    .rs-table-row-selected {
      .fileName {
        color: ${({ theme }) => theme.primaryColor};

        .icon {
          color: ${({ theme }) => theme.primaryColor};
        }

        .toolbar {
          visibility: visible;
        }
      }
    }
  }

  > .footer {
    padding: 6px 20px;
    font-size: 12px;
    color: rgba(0, 0, 0, 0.45);
    background-color: #f6f8fa;

    > .icon {
      font-size: inherit;
      display: inline-block;
      width: 18px;
      padding-right: 5px;
    }
  }
`
const TableLinkBtn = styled(Button)`
  padding: 0;
`

type Props = {
  width?: number
  height?: number
}

export const FileList = observer(function FileList({ width, height }: Props) {
  const store = useStore()
  const { selectedKeys, server, getWidget, isWidgetVisible } = store
  const [fetch, loading] = store.useRefresh()

  const state = useLocalStore(() => ({
    get dataSource() {
      const { searchKey, dir } = store
      const { selectedTypes } = this
      let files = [...(dir?.children || [])].filter(item => {
        let flag = true

        // hack: filter new directory
        if (!item.path) {
          return false
        }

        if (searchKey) {
          flag = flag && item.name.includes(searchKey)
        }

        if (selectedTypes.length > 0) {
          flag = flag && selectedTypes.includes(item.type)
        }

        return flag
      })
      if (this.sortType && this.sortKey) {
        switch (this.sortKey) {
          // sort by name
          case 'name': {
            files = files.sort((x, y) => {
              if (this.sortType === 'asc') {
                return x.name.localeCompare(y.name)
              } else {
                return y.name.localeCompare(x.name)
              }
            })
            break
          }
          // sort by size
          case 'size': {
            files = files.sort((x, y) => {
              if (this.sortType === 'asc') {
                return x.size - y.size
              } else {
                return y.size - x.size
              }
            })
            break
          }
          // sort by modifiedTime
          case 'mtime': {
            files = files.sort((x, y) => {
              const xTime = new Date(x.mtime).getTime()
              const yTime = new Date(y.mtime).getTime()

              return this.sortType === 'asc' ? xTime - yTime : yTime - xTime
            })
            break
          }
        }
      }

      return files.map(item => ({
        ...item,
        size: item.isFile ? formatByte(item.size || 0) : '--',
        editable: item.isFile && item.size <= EDITABLE_SIZE,
        mtime: formatUnixTime(item.mtime)
      }))
    },
    get types() {
      const { dir } = store
      return [
        ...new Set(
          [...(dir?.children || [])]
            .map(item => item.type)
            .filter(item => !!item)
        )
      ]
    },
    get fileCount() {
      return this.dataSource.filter(item => item.isFile).length
    },
    get dirCount() {
      return this.dataSource.filter(item => !item.isFile).length
    },
    sortKey: undefined,
    setSortKey(key) {
      this.sortKey = key
    },
    sortType: undefined,
    setSortType(type) {
      this.sortType = type
    },
    selectedTypes: [],
    setSelectedTypes(types) {
      this.selectedTypes = [...types]
    }
  }))
  const { dataSource, types, fileCount, dirCount } = state

  const deletable = isWidgetVisible('delete')
  const editable = isWidgetVisible('edit')

  useEffect(() => {
    fetch()
  }, [fetch])

  async function deleteNode(id: string) {
    const node = store.dir.filterFirstNode(item => item.id === id)
    await Modal.showConfirm({
      title: '删除文件',
      content: `确认要删除文件 ${node.name} 吗`
    })

    await server.delete([node.path])

    await fetch()
    message.success('文件删除成功')
  }

  return (
    <StyledLayout>
      <div className='list_body'>
        <Table
          {...(getWidget('custom-column') || {
            tableId: 'file_mgt_table',
            defaultConfig: [
              ['name', true],
              ['type', true],
              ['size', true],
              ['mtime', true],
              ['operator', true]
            ]
          })}
          props={{
            width,
            height,
            data: state.dataSource || [],
            loading,
            rowKey: 'id'
            // virtualized: true
          }}
          rowSelection={{
            selectedKeys,
            onChange(keys) {
              store.setSelectedKeys(keys)
            }
          }}
          columns={[
            {
              header: '文件名称',
              props: {
                width: 192.75,
                minWidth: 350,
                resizable: true
              },
              sorter: ({ sortType, sortKey }) => {
                state.setSortKey(sortKey)
                state.setSortType(sortType)
              },
              cell: {
                props: {
                  dataKey: 'name'
                },
                render({ rowData }) {
                  return <Name key={rowData['id']} nodeId={rowData['id']} />
                }
              }
            },
            {
              header: '类型',
              dataKey: 'type',
              props: {
                width: 192.75
              },
              filter: {
                items: types.map(item => ({ key: item, name: item })),
                onChange(keys) {
                  state.setSelectedTypes(keys)
                }
              }
            },
            {
              header: '文件大小',
              dataKey: 'size',
              props: {
                width: 192.75
              },
              sorter: ({ sortType, sortKey }) => {
                state.setSortKey(sortKey)
                state.setSortType(sortType)
              }
            },
            {
              header: '修改时间',
              dataKey: 'mtime',
              props: {
                width: 192.75
              },
              sorter: ({ sortType, sortKey }) => {
                state.setSortKey(sortKey)
                state.setSortType(sortType)
              }
            },
            ...(deletable || editable
              ? [
                  {
                    header: '操作',
                    props: {
                      width: 150
                    },
                    cell: {
                      props: {
                        dataKey: 'operator'
                      },
                      render: ({ rowData }) => {
                        const id = rowData.id
                        const type = rowData.type
                        const isImage = /gif|jpe?g|tiff|png|webp|bmp$/i.test(
                          type
                        )

                        return (
                          <>
                            {getWidget('delete') || (
                              <TableLinkBtn
                                type='link'
                                onClick={() => deleteNode(id)}>
                                删除
                              </TableLinkBtn>
                            )}
                            {deletable && editable && (
                              <Divider type='vertical' />
                            )}
                            {getWidget('edit') || (
                              <Edit
                                disabled={!rowData.isFile || isImage}
                                nodeId={id}
                                readonly={!rowData.editable}
                              />
                            )}
                          </>
                        )
                      }
                    }
                  }
                ]
              : [])
          ]}
        />
      </div>
      <div className='footer'>
        <Icon className='icon' type='drawer' />
        当前目录包含{dirCount}个文件夹，{fileCount}个文件
      </div>
    </StyledLayout>
  )
})
