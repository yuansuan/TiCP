/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import Icon from '..'

export const Custom = () => (
  <div>
    <div>
      fontSize: <Icon type='rename' style={{ fontSize: 12 }} />
      <Icon type='rename' style={{ fontSize: 14 }} />
      <Icon type='rename' style={{ fontSize: 16 }} />
      <Icon type='rename' style={{ fontSize: 18 }} />
      <Icon type='rename' style={{ fontSize: 20 }} />
      <Icon type='rename' style={{ fontSize: 25 }} />
    </div>
    <div>
      color: <Icon type='rename' />
      <Icon type='rename' style={{ color: 'pink' }} />
      <Icon type='rename' style={{ color: 'red' }} />
      <Icon type='rename' style={{ color: 'green' }} />
      <Icon type='rename' style={{ color: 'blue' }} />
      <Icon type='rename' style={{ color: 'gray' }} />
    </div>
  </div>
)
