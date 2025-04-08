/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import Mask from '..'

export const Theme = () => (
  <div
    style={{
      position: 'relative',
      height: 200,
      display: 'flex',
    }}>
    <div style={{ position: 'relative', flex: 1 }}>
      <Mask theme='black'>It's a black mask.</Mask>
    </div>
    <div style={{ position: 'relative', flex: 1 }}>
      <Mask theme='white'>It's a white mask.</Mask>
    </div>
  </div>
)
