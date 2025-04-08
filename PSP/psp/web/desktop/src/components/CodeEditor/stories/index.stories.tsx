/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState, useRef } from 'react'
import { Story, Meta } from '@storybook/react/types-6-0'
import { message } from 'antd'
import CodeEditor, { CodeEditorProps } from '..'
import mdx from './doc.mdx'
import { Button } from '../..'

export default {
  title: 'components/CodeEditor',
  component: CodeEditor,
  parameters: {
    docs: {
      page: mdx,
    },
  },
} as Meta

const Template: Story<CodeEditorProps> = props => (
  <div style={{ height: 200 }}>
    <CodeEditor value='var name = "hello world"' {...props} />
  </div>
)

export const Value = Template.bind({})
Value.args = {
  language: 'javascript',
  value: 'var name = "hello world"',
}

export const Readonly = Template.bind({})
Readonly.args = {
  readOnly: true,
}

export const Controlled = function Controlled() {
  const [value, setValue] = useState('var name = "hello world"')

  return (
    <div style={{ height: 200 }}>
      <CodeEditor
        language='javascript'
        value={value}
        onChange={({ detail, value }) => {
          setValue(value)
        }}
      />
    </div>
  )
}

export const Imperative = function Imperative() {
  const editorRef = useRef(undefined)

  function getValue() {
    message.info(editorRef.current.getValue())
  }

  return (
    <div>
      <Button style={{ margin: 10 }} type='primary' onClick={getValue}>
        获取编辑器内容
      </Button>
      <div style={{ height: 200 }}>
        <CodeEditor
          ref={editorRef}
          language='javascript'
          value='var name = "hello world"'
        />
      </div>
    </div>
  )
}
