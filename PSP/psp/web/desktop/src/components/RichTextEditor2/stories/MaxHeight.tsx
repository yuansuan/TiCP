/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import RichTextEditor from '..'

export function MaxHeight() {
  return (
    <div>
      <div style={{ color: 'gray', padding: 10 }}>
        可以缩放页面动态控制编辑器最大高度
      </div>
      <RichTextEditor
        style={{ height: 300 }}
        autoHeight
        defaultValue={RichTextEditor.createEditorState('hello world')}
      />
    </div>
  )
}
