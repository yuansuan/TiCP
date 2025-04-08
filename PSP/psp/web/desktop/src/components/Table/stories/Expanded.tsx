/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Tag as AntTag, Divider as AntDivider } from 'antd'
import Table from '..'
import { useState } from '@storybook/addons'

const data = [
  {
    key: '1',
    name: 'John Brown',
    age: 32,
    address: 'New York No. 1 Lake Park',
    tags: ['nice', 'developer']
  },
  {
    key: '2',
    name: 'Jim Green',
    age: 42,
    address: 'London No. 1 Lake Park',
    tags: ['loser']
  },
  {
    key: '3',
    name: 'Joe Black',
    age: 32,
    address: 'Sidney No. 1 Lake Park',
    tags: ['cool', 'teacher']
  }
]

const Basic = () => (
  <Table
    props={{ autoHeight: true, data }}
    columns={[
      {
        header: 'Name',
        dataKey: 'name',
        props: {
          flexGrow: 1
        }
      },
      {
        header: 'Age',
        dataKey: 'age',
        props: {
          flexGrow: 1
        }
      },
      {
        header: 'Address',
        dataKey: 'address',
        props: {
          flexGrow: 1
        }
      },
      {
        header: 'Tags',

        props: {
          flexGrow: 1
        },
        cell: {
          props: {
            dataKey: 'tags'
          },
          render: ({ rowData, dataKey }) => {
            const tags = rowData[dataKey]

            return (
              <span>
                {tags.map(tag => {
                  let color = tag.length > 5 ? 'geekblue' : 'green'
                  if (tag === 'loser') {
                    color = 'volcano'
                  }
                  return (
                    <AntTag color={color} key={tag}>
                      {tag.toUpperCase()}
                    </AntTag>
                  )
                })}
              </span>
            )
          }
        }
      },
      {
        header: 'Action',
        props: {
          flexGrow: 1
        },
        cell: {
          props: {
            dataKey: 'action'
          },
          render: ({ rowData }) => {
            return (
              <span>
                <a>Invite {rowData.name}</a>
                <AntDivider type='vertical' />
                <a>Delete</a>
              </span>
            )
          }
        }
      }
    ]}
  />
)

export const Expanded = function Expanded() {
  const [expand, setExpand] = useState([] as any[])
  const [isExpand, setIsExpand] = useState(false)

  async function handleExpanded(key) {
    const nextExpandedRowKeys = expand
    if (!expand.includes(key)) {
      nextExpandedRowKeys.push(key)
    } else {
      const index = expand.indexOf(key)
      nextExpandedRowKeys.splice(index, 1)
    }
    setExpand([...nextExpandedRowKeys])
  }

  return (
    <Table
      props={{
        autoHeight: true,
        data,
        rowKey: 'key',
        expandedRowKeys: expand,
        rowExpandedHeight: 250,
        onRowClick: data => handleExpanded(data['key']),
        renderRowExpanded: rowData => {
          return <Basic />
        }
      }}
      columns={[
        {
          header: 'Name',
          dataKey: 'name',
          props: {
            flexGrow: 1
          },
          cell: {
            props: {
              dataKey: 'name'
            },
            render: ({ rowData, dataKey }) => (
              <div>
                {expand?.some(key => key === rowData['key']) ? '-  ' : '+  '}
                {rowData[dataKey]}
              </div>
            )
          }
        },
        {
          header: 'Age',
          dataKey: 'age',
          props: {
            flexGrow: 1
          }
        },
        {
          header: 'Address',
          dataKey: 'address',
          props: {
            flexGrow: 1
          }
        },
        {
          header: 'Tags',

          props: {
            flexGrow: 1
          },
          cell: {
            props: {
              dataKey: 'tags'
            },
            render: ({ rowData, dataKey }) => {
              const tags = rowData[dataKey]

              return (
                <span>
                  {tags.map(tag => {
                    let color = tag.length > 5 ? 'geekblue' : 'green'
                    if (tag === 'loser') {
                      color = 'volcano'
                    }
                    return (
                      <AntTag color={color} key={tag}>
                        {tag.toUpperCase()}
                      </AntTag>
                    )
                  })}
                </span>
              )
            }
          }
        },
        {
          header: 'Action',
          props: {
            flexGrow: 1
          },
          cell: {
            props: {
              dataKey: 'action'
            },
            render: ({ rowData }) => {
              return (
                <span>
                  <a>Invite {rowData.name}</a>
                  <AntDivider type='vertical' />
                  <a>Delete</a>
                </span>
              )
            }
          }
        }
      ]}
    />
  )
}
