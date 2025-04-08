/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import { editor } from 'monaco-editor/esm/vs/editor/editor.api.js'
import 'monaco-editor/min/vs/editor/editor.main.css'
import 'monaco-editor/esm/vs/editor/contrib/find/findController.js'
import 'monaco-editor/esm/vs/basic-languages/shell/shell.contribution'
import 'monaco-editor/esm/vs/basic-languages/yaml/yaml.contribution'
import 'monaco-editor/esm/vs/basic-languages/javascript/javascript.contribution'

type languages = 'shell' | 'javascript' | 'yaml'

export interface CodeEditorProps {
  value?: string
  readOnly?: boolean
  onChange?: (e: { detail; value: string }) => void
  language?: languages
}

export default class CodeEditor extends React.Component<CodeEditorProps> {
  el = null // 编辑器实例化的dom节点
  editor = null // 编辑器实例
  model = null // 编辑器绑定的text model
  disposable = null // 内容改变事件disposable

  componentWillReceiveProps(np: CodeEditorProps) {
    if (np.value && np.value !== this.props.value) {
      const scrollTop = this.editor.getScrollTop()
      this.setValue(np.value)
      this.editor.setScrollTop(scrollTop)
    }
  }

  componentDidMount() {
    const { readOnly = false } = this.props

    this.editor = editor.create(this.el, {
      minimap: { enabled: false }, // 关闭右侧mini地图
      fontSize: 14,
      value: this.props.value,
      lineNumbersMinChars: 7,
      scrollBeyondLastLine: false,
      readOnly,
    })
    this.model = this.editor.getModel()

    this.editor.focus()
    if (this.props.language) {
      editor.setModelLanguage(this.model, this.props.language)
    }

    const { onChange } = this.props
    if (onChange) {
      this.disposable = this.editor.onDidChangeModelContent(e => {
        onChange({
          detail: e,
          value: this.getValue(),
        })
      })
    }
  }

  componentWillUnmount() {
    this.editor && this.editor.dispose()
    this.disposable && this.disposable.dispose()
  }

  setValue = (value: string) => this.editor.setValue(value)

  getValue = () => this.model.getValue(editor.EndOfLinePreference.LF)

  render() {
    return (
      <div
        style={{
          height: '100%',
          fontSize: '14px',
        }}
        ref={el => (this.el = el)}
      />
    )
  }
}
