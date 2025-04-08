/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import CodeEditor from '..'

export const Readonly = () => (
  <div style={{ height: 200 }}>
    <CodeEditor
      language='javascript'
      value='var name = "hello world"'
      readOnly
    />
  </div>
)
