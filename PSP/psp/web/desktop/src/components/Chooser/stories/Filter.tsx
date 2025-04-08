/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import Chooser from '..'

export const Filter = () => (
  <Chooser
    style={{
      margin: 20,
    }}
    items={[
      { name: 'item01', key: '1' },
      { name: 'item02', key: '2' },
      { name: 'item03', key: '3' },
      { name: 'item04', key: '4' },
      { name: 'item05', key: '5' },
      { name: 'item06', key: '6' },
    ]}
    filter={(item, keywords) => {
      return !keywords || item.name === keywords
    }}
  />
)
