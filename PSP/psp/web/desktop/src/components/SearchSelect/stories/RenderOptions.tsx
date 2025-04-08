/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Select as AntSelect } from 'antd'
import SearchSelect from '..'

export const RenderOptions = () => (
  <SearchSelect
    style={{ margin: 20, width: 150 }}
    renderOptions={options =>
      options.map(item => (
        <AntSelect.OptGroup key={item.key}>
          <AntSelect.Option key={item.key} value={item.name}>
            {item.name}
          </AntSelect.Option>
        </AntSelect.OptGroup>
      ))
    }
    options={[
      { key: 1, name: 'option_01' },
      { key: 2, name: 'OPTION_02' },
      { key: 3, name: 'option_03' },
      { key: 4, name: 'OPTION_04' },
      { key: 5, name: 'option_05' },
    ]}
  />
)
