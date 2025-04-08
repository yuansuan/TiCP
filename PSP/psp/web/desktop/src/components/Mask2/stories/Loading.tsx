/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import Mask from '..'

export const Loading = () => (
  <div style={{ position: 'relative', height: 200 }}>
    <Mask.Spin theme='white' />
  </div>
)
