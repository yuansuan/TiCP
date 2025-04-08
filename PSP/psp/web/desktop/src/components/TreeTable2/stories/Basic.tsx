/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState } from 'react'
import { TreeTable, Button } from '../..'

const columns = [
  {
    title: '文件名',
    key: 'filename',
    width: 600,
  },
  {
    title: '操作',
    key: 'actions',
    render() {
      return <Button>操作</Button>
    },
  },
]

function moveIdToParentId(data, id, parentId) {
  const res = [...data]
  const filterFirstItem = (d, filter) => {
    let r = null
    d.some(item => {
      if (filter(item)) {
        r = item
        return true
      }
      if (!item.children) {
        return false
      }
      r = filterFirstItem(item.children, filter)
      return !!r
    })
    return r
  }

  const child = filterFirstItem(res, item => item.id === id)

  // rm self
  const parent = filterFirstItem(res, item => {
    if (!item.children) return false
    return item.children.some(c => c.id === id)
  })
  const p = parent ? parent.children : res
  const index = p.findIndex(r => r.id === id)
  p.splice(index, 1)

  // move to target
  const targetParent = filterFirstItem(res, item => item.id === parentId)
  if (child && targetParent && targetParent.children) {
    targetParent.children.unshift(child)
  }
  return res
}

export function Basic() {
  const [data, setData] = useState([
    {
      id: '1',
      filename: 'dir_a',
      filesize: 1024,
      status: 'done',
      isDirectory: true,
      children: [
        {
          id: '1-1',
          filename: 'dir_a_a',
          filesize: 32,
          status: 'done',
          isDirectory: true,
          children: [
            {
              id: '1-1-1',
              isDirectory: false,
              filename: 'file_a_a_a',
              filesize: 10032,
              status: 'done',
            },
          ],
        },
        {
          id: '1-2',
          isDirectory: false,
          filename: 'file_a_a',
          filesize: 32,
          status: 'done',
        },
      ],
    },
    {
      id: '2',
      filename: 'file_b',
      filesize: 32,
      status: 'done',
    },
  ])

  const onDragEnd = (dragKey, dropKey) => {
    const d = moveIdToParentId(data, dragKey, dropKey)
    setData(d)
  }

  return (
    <div className='App'>
      <TreeTable
        rowKey='id'
        dataSource={data}
        columns={columns}
        defaultExpandAll={true}
        onDragEnd={onDragEnd}
      />
    </div>
  )
}
