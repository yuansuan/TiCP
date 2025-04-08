/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState } from 'react'
import { Input as AntInput } from 'antd'
import RichTextEditor from '..'

export function Controlled() {
  const [value, setValue] = useState(
    RichTextEditor.createEditorState('hello world')
  )

  function onChange(e) {
    setValue(RichTextEditor.createEditorState(e.target.value))
  }

  return (
    <div style={{ height: 300 }}>
      <div style={{ padding: 10, borderBottom: '1px solid #ccc' }}>
        <div style={{ color: 'gray' }}>通过改变输入框内容控制编辑器内容：</div>
        <AntInput value={value.toText()} onChange={onChange} />
      </div>
      <RichTextEditor contentStyle={{ maxHeight: 200 }} value={value} />
    </div>
  )
}
