/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { IColumn } from '../types'

interface IProps {
  columns: IColumn[]
}

export function TableHeader(props: IProps) {
  const { columns } = props
  return (
    <>
      <colgroup>
        {columns.map(c => {
          return (
            <col
              key={c.key}
              style={{ width: `${c.width}px`, minWidth: `${c.width}px` }}
            />
          )
        })}
      </colgroup>
      <thead>
        <tr>
          {columns.map(c => {
            let data: any = c.title
            if (c.headerRender) {
              data = c.headerRender(c.title)
            }
            return <th key={c.key}>{data}</th>
          })}
        </tr>
      </thead>
    </>
  )
}
