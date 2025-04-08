import React from 'react'
import CodeMirror from '@uiw/react-codemirror'
import { EditorView } from "@codemirror/view";
import { loadLanguage } from '@uiw/codemirror-extensions-langs';

interface IProps {
  onChange?: (value: string) => void
  value: string
  height: string
  width?: string
  readOnly?: boolean
  lang?: string
}

function CodeEditor({onChange, value, height, readOnly, width='100%', lang='powershell'}: IProps) {
  const handleChange = React.useCallback((value, viewUpdate) => {
    onChange(value)
  }, [])

  return (
    <CodeMirror
      readOnly={readOnly}
      value={value}
      height={height}
      width={width}
      extensions={[
        loadLanguage(lang as any),
        EditorView.lineWrapping,
      ]}
      onChange={handleChange}
    />
  )
}
export default CodeEditor