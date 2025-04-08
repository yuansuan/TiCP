/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState } from 'react'
import PartialList from '..'
import { Input as AntInput } from 'antd'

export function Controlled() {
  const [maxWidth, setMaxWidth] = useState(100)

  function onChange(e) {
    setMaxWidth(e.target.value)
  }

  return (
    <div>
      <div style={{ padding: 10, borderBottom: '1px solid #ccc' }}>
        <span style={{ color: 'gray' }}>自定义 maxWidth：</span>
        <AntInput style={{ width: 200 }} value={maxWidth} onChange={onChange} />
      </div>
      <PartialList
        style={{ margin: 20 }}
        maxWidth={maxWidth}
        items={[
          0,
          1,
          2,
          3,
          4,
          5,
          6,
          7,
          8,
          9,
          10,
          11,
          12,
          13,
          14,
          15,
          16,
          17,
          18,
        ]}
      />
    </div>
  )
}
