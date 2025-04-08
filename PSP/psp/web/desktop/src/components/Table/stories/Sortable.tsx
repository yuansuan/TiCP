/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useMemo } from 'react'
import { Tag as AntTag, Divider as AntDivider } from 'antd'
import Table from '..'

const data = [
  {
    key: '1',
    name: 'John Brown',
    age: 32,
    address: 'New York No. 1 Lake Park',
    tags: ['nice', 'developer'],
  },
  {
    key: '2',
    name: 'Jim Green',
    age: 42,
    address: 'London No. 1 Lake Park',
    tags: ['loser'],
  },
  {
    key: '3',
    name: 'Joe Black',
    age: 32,
    address: 'Sidney No. 1 Lake Park',
    tags: ['cool', 'teacher'],
  },
]

export function Sortable() {
  const [sortType, setSortType] = useState('')
  const [sortKey, setSortKey] = useState('')
  const sortedData = useMemo(() => {
    // Array.sort will mutate original array, create a copy before sort
    const originalData = [...data]

    if (sortKey === 'name' && sortType) {
      return originalData.sort((x, y) => {
        if (sortType === 'asc') {
          return x.name.localeCompare(y.name)
        } else {
          return y.name.localeCompare(x.name)
        }
      })
    } else if (sortKey === 'age' && sortType) {
      return originalData.sort((x, y) => {
        if (sortType === 'asc') {
          return x.age - y.age
        } else {
          return y.age - x.age
        }
      })
    } else {
      return originalData
    }
  }, [sortType, sortKey])

  function onSort({ sortType, sortKey }) {
    setSortKey(sortKey)
    setSortType(sortType)
  }

  return (
    <Table
      props={{ autoHeight: true, data: sortedData, rowKey: 'name' }}
      columns={[
        {
          header: 'Name',
          dataKey: 'name',
          sorter: onSort,
          props: {
            flexGrow: 1,
          },
        },
        {
          header: 'Age',
          dataKey: 'age',
          sorter: onSort,
          props: {
            flexGrow: 1,
          },
        },
        {
          header: 'Address',
          dataKey: 'address',
          props: {
            flexGrow: 1,
          },
        },
        {
          header: 'Tags',
          props: {
            flexGrow: 1,
          },
          cell: {
            props: {
              dataKey: 'tags',
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
            },
          },
        },
        {
          header: 'Action',
          props: {
            flexGrow: 1,
          },
          cell: {
            props: {
              dataKey: 'action',
            },
            render: ({ rowData }) => {
              return (
                <span>
                  <a>Invite {rowData.name}</a>
                  <AntDivider type='vertical' />
                  <a>Delete</a>
                </span>
              )
            },
          },
        },
      ]}
    />
  )
}
