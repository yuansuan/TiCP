import * as React from 'react'
import styled from 'styled-components'
import { Radio, Input } from 'antd'
import { observer } from 'mobx-react'
import { observable, action } from 'mobx'
import { RichTextEditor } from '@/components'

const Wrapper = styled.div`
  padding: 20px;
  height: 100%;
  padding-bottom: 30px;
  display: flex;
  flex-direction: column;
  background-color: #f0f5fd;

  > div {
    > div {
      margin: 5px 0;
    }
  }

  .urlEditor {
    margin: 10px 0;
  }

  .contentEditor {
    min-width: 860px;
    max-width: 1136px;
    height: calc(100% - 11px);
    display: flex;
    flex-direction: column;
    overflow: hidden;

    .editor {
      border: 1px solid #d9d9d9;
      background-color: '#eee';
      height: calc(100% - 21px);

      .bf-content {
        height: calc(100% - 92px);
      }
    }
  }

  .hidden {
    display: none;
  }
`

interface IProps {
  helpDoc: {
    type: string
    value: string
  }
}

@observer
export default class Document extends React.Component<IProps> {
  @observable url = ''
  @observable content = ''
  @action
  updateType = type => (this.props.helpDoc.type = type)
  @action
  updateUrl = url => (this.url = url)
  @action
  updateContent = content => (this.content = content)

  constructor(props) {
    super(props)

    const { type, value } = props.helpDoc
    if (type === 'url') {
      this.updateUrl(value)
    } else if (type === 'content') {
      this.updateContent(value)
    }
  }

  @action
  onChangeUrl = e => {
    const { value } = e.target
    this.updateUrl(value)

    const { helpDoc } = this.props
    if (helpDoc.type === 'url') {
      helpDoc.value = value
    }
  }

  @action
  onChangeContent = editorState => {
    if (!editorState) {
      return
    }
    const value = editorState.toHTML()
    // this.updateContent(value)

    const { helpDoc } = this.props
    if (helpDoc.type === 'content') {
      helpDoc.value = value
    }
  }

  render() {
    const { type } = this.props.helpDoc

    return (
      <Wrapper>
        <div className='urlEditor'>
          <div>
            <Radio
              checked={type === 'url'}
              onChange={e => {
                if (e.target.checked) {
                  this.updateType('url')
                }
              }}>
              请填写外部帮助文档的链接
            </Radio>
          </div>
          <div className={type === 'url' ? '' : 'hidden'}>
            <Input
              maxLength={64}
              value={this.url}
              onChange={this.onChangeUrl}
            />
          </div>
        </div>
        <div className='contentEditor'>
          <div>
            <Radio
              checked={type === 'content'}
              onChange={e => {
                if (e.target.checked) {
                  this.updateType('content')
                }
              }}>
              自定义帮助文档
            </Radio>
          </div>
          <div className={`editor ${type === 'content' ? '' : 'hidden'}`}>
            <RichTextEditor
              value={RichTextEditor.createEditorState(this.content)}
              onChange={this.onChangeContent}
            />
          </div>
        </div>
      </Wrapper>
    )
  }
}
