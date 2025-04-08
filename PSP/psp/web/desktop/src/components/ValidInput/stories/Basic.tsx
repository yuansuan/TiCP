/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import ValidInput from '..'

export const Basic = () => (
  <div style={{ margin: 20 }}>
    <div style={{ color: 'gray' }}>允许输入 0-5 位数字或字母组成的字符串</div>
    <ValidInput
      style={{ marginTop: 10, width: 200 }}
      validator={/^[0-9a-zA_Z]{0,5}$/}
    />
  </div>
)
