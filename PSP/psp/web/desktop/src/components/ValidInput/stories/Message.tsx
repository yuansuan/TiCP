/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState } from 'react'
import { message } from 'antd'
import ValidInput from '..'

export function Message() {
  const [value, setValue] = useState('')

  function onChange(e) {
    setValue(e.target.value)
  }

  return (
    <div style={{ margin: 20 }}>
      <div style={{ color: 'gray' }}>允许输入 0-5 位数字或字母组成的字符串</div>
      <ValidInput
        style={{ marginTop: 10, width: 200 }}
        value={value}
        onChange={onChange}
        validator={value => {
          const valid = /^[0-9a-zA_Z]{0,5}$/.test(value)

          if (!valid) {
            message.error('只允许输入 0-5 位数字或字母组成的字符串')
          }

          return valid
        }}
      />
    </div>
  )
}
