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

export function Filterable() {
  const [selectedNames, setSelectedNames] = useState([])
  const [selectedAges, setSelectedAges] = useState([])
  const [selectedAddresses, setSelectedAddresses] = useState([])

  const visibleData = useMemo(() => {
    let finalData = [...data]
    if (selectedNames.length > 0) {
      finalData = finalData.filter(item => selectedNames.includes(item.name))
    }
    if (selectedAges.length > 0) {
      finalData = finalData.filter(item => selectedAges.includes(item.age))
    }
    if (selectedAddresses.length > 0) {
      finalData = finalData.filter(item =>
        selectedAddresses.includes(item.address)
      )
    }

    return finalData
  }, [selectedNames, selectedAges, selectedAddresses])

  return (
    <Table
      props={{ autoHeight: true, data: visibleData, rowKey: 'name' }}
      columns={[
        {
          header: 'Name',
          dataKey: 'name',
          filter: {
            onChange: keys => {
              setSelectedNames(keys)
            },
            items: Array.from(new Set(data.map(item => item.name))).map(
              name => ({
                key: name,
                name,
              })
            ),
          },
          props: {
            flexGrow: 1,
          },
        },
        {
          header: 'Age',
          dataKey: 'age',
          filter: {
            onChange: keys => {
              setSelectedAges(keys)
            },
            items: Array.from(new Set(data.map(item => item.age))).map(age => ({
              key: age,
              name: age,
            })),
          },
          props: {
            flexGrow: 1,
          },
        },
        {
          header: 'Address',
          dataKey: 'address',
          filter: {
            onChange: keys => {
              setSelectedAddresses(keys)
            },
            items: Array.from(new Set(data.map(item => item.address))).map(
              address => ({
                name: address,
                key: address,
              })
            ),
          },
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
