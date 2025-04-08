/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useRef, useEffect } from 'react'
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

export function PercentWidth() {
  const [width, setWidth] = useState(0)
  const containerRef = useRef(undefined)

  useEffect(() => {
    setWidth(containerRef.current.clientWidth)
    window.addEventListener('resize', onWindowResize)

    return () => {
      window.removeEventListener('resize', onWindowResize)
    }
  }, [])

  function onWindowResize() {
    setWidth(containerRef.current.clientWidth)
  }

  return (
    <div ref={containerRef}>
      <Table
        props={{
          autoHeight: true,
          data,
          width,
        }}
        columns={[
          {
            header: 'Name',
            dataKey: 'name',
            props: {
              width: '50%',
            },
          },
          {
            header: 'Age',
            dataKey: 'age',
            props: {
              width: 100,
            },
          },
          {
            header: 'Address',
            dataKey: 'address',
            props: {
              width: 200,
            },
          },
          {
            header: 'Tags',
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
    </div>
  )
}
