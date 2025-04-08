/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import SearchSelect from '..'

export const Basic = () => (
  <SearchSelect
    style={{ margin: 20, width: 120 }}
    options={[
      [1, 'option_01'],
      [2, 'OPTION_02'],
      [3, 'option_03'],
      [4, 'OPTION_04'],
      [5, 'option_05'],
    ]}
  />
)
