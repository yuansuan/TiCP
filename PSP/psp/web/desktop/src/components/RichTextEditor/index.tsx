import React from 'react'
import BraftEditor, { ControlType } from 'braft-editor'
import 'braft-editor/dist/index.css'
import { StyledEditor } from './style'

export const EditorControls: ControlType[] = [
  'font-size',
  'text-color',
  'bold',
  'italic',
  'underline',
  'strike-through',
  'text-indent',
  'link',
  'headings',
  'list-ul',
  'list-ol',
  'blockquote',
  'code',
  'media'
]

export interface IProps {
  value?: string
  contentStyle?: React.CSSProperties
  placeholder?: string
  onChange?: (editorState: any) => void
}

export class RichTextEditor extends React.Component<IProps> {
  static createEditorState = BraftEditor.createEditorState

  render() {
    const { ...rest } = this.props
    return (
      <StyledEditor>
        <BraftEditor controls={EditorControls} {...rest} />
      </StyledEditor>
    )
  }
}
