/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import RichTextEditor from '..'

export function Basic() {
  function onChangeContent(editorState) {
    const value = editorState.toHTML()
    console.log(value)
  }

  return (
    <div style={{ height: 300 }}>
      <RichTextEditor
        contentStyle={{ maxHeight: 200 }}
        defaultValue={RichTextEditor.createEditorState('hello world')}
        onChange={onChangeContent}
      />
    </div>
  )
}
