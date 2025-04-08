/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState } from 'react'
import Chooser from '..'

export function Controlled() {
  const [selectedKeys, setSelectedKeys] = useState<Array<string | number>>([
    '1',
    '2',
  ])

  return (
    <Chooser
      style={{
        margin: 20,
      }}
      items={[
        ['1', 'item01'],
        ['2', 'item02'],
        ['3', 'item03'],
        ['4', 'item04'],
        ['5', 'item05'],
        ['6', 'item06'],
      ]}
      selectedKeys={selectedKeys}
      updateSelectedKeys={keys => setSelectedKeys(keys)}
    />
  )
}
