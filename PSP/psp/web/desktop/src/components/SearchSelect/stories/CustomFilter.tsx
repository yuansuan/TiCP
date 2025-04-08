/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import SearchSelect from '..'

export const CustomFilter = () => (
  <SearchSelect
    style={{ margin: 20, width: 120 }}
    filter={(searchKey, item) => {
      return item.name.toLowerCase().includes(searchKey.toLowerCase())
    }}
    options={[
      { key: 1, name: 'option_01' },
      { key: 2, name: 'OPTION_02' },
      { key: 3, name: 'option_03' },
      { key: 4, name: 'OPTION_04' },
      { key: 5, name: 'option_05' },
    ]}
  />
)
