/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import Icon from '..'

export const Basic = () => (
  <div>
    <div>
      origin：
      <Icon type='rename' />
    </div>
    <div>
      disabled: <Icon type='rename' disabled />
      <Icon type='rename' className='disabled' />
    </div>
  </div>
)
