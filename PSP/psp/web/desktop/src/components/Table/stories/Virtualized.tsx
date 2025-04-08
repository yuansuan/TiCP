/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Tag as AntTag, Divider as AntDivider } from 'antd'
import Table from '..'

const data = Array.from({ length: 10000 }).map((item, index) => ({
  key: index,
  name: `name_${index}`,
  age: Math.ceil(Math.random() * 90 + 10),
  address: `address_${index}`,
  tags: ['nice', 'developer'],
}))

export const Virtualized = () => (
  <Table
    props={{
      data,
      rowKey: 'name',
      // enable virtualized
      virtualized: true,
      height: 300,
    }}
    columns={[
      {
        header: 'Name',
        dataKey: 'name',
        props: {
          flexGrow: 1,
        },
      },
      {
        header: 'Age',
        dataKey: 'age',
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
