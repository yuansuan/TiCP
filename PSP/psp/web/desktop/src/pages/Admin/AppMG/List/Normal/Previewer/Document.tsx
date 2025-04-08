import * as React from 'react'
import styled from 'styled-components'
import { RichTextEditor } from '@/components'

const Wrapper = styled.div`
  height: 100%;
  padding-bottom: 50px;

  iframe {
    width: 100%;
    height: 100%;
    min-height: 600px;
    border: none;
  }
`

interface IProps {
  helpDoc: { type: string; value: string }
}

export default class Document extends React.Component<IProps> {
  render() {
    const { type, value } = this.props.helpDoc

    return (
      <Wrapper>
        {type === 'url' && value ? <iframe src={value} /> : <div />}
        {type === 'content' ? (
          <RichTextEditor
            value={RichTextEditor.createEditorState(value)}
            readOnly={true}
            controls={[]}
          />
        ) : null}
      </Wrapper>
    )
  }
}
