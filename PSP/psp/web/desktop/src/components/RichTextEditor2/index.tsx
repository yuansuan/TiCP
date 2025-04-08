/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import BraftEditor, { BraftEditorProps } from 'braft-editor'
import 'braft-editor/dist/index.css'
import ReactDOM from 'react-dom'
import { debounce } from 'lodash'

export * from 'braft-editor'

type IProps = BraftEditorProps & {
  autoHeight?: boolean
}

export default class RichTextEditor extends React.Component<IProps> {
  static createEditorState = BraftEditor.createEditorState

  state = {
    maxHeight: undefined,
  }

  ref = null
  onResize = this.props.autoHeight
    ? debounce(
        () => {
          const editorNode = ReactDOM.findDOMNode(this.ref) as Element
          const toolbarHeight = editorNode.querySelector('.bf-controlbar')
            .clientHeight

          this.setState({
            maxHeight: editorNode.clientHeight - toolbarHeight,
          })
        },
        300,
        { leading: true }
      )
    : null

  componentDidMount() {
    if (this.onResize) {
      this.onResize()
      window.addEventListener('resize', this.onResize)
    }
  }

  componentWillUnmount() {
    if (this.onResize) {
      window.removeEventListener('resize', this.onResize)
    }
  }

  render() {
    const { maxHeight } = this.state
    const { autoHeight, contentStyle = {}, ...rest } = this.props

    return (
      <BraftEditor
        ref={ref => (this.ref = ref)}
        contentStyle={{ height: 'auto', ...contentStyle, maxHeight }}
        excludeControls={['fullscreen']}
        {...rest}
      />
    )
  }
}
