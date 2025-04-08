/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState } from 'react'
import { message } from 'antd'
import EditableText from '..'

export function Validation() {
  const [value, setValue] = useState('hello world')

  return (
    <EditableText
      style={{ width: 200, margin: 20 }}
      beforeConfirm={value => {
        const flag = /^[a-zA-Z0-9]$/.test(value)

        if (!flag) {
          message.error('只允许输入字母和数字')
        }

        return flag
      }}
      help='只允许输入字母和数字'
      defaultEditing={true}
      defaultValue={value}
      onConfirm={setValue}
    />
  )
}
