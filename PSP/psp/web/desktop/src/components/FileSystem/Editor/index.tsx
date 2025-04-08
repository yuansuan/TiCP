import * as React from 'react'
import { Spin, message } from 'antd'
import { observer } from 'mobx-react'
import { observable, action, computed } from 'mobx'
import 'monaco-editor/min/vs/editor/editor.main.css'

import { createMobxStream } from '@/utils'
import { CodeEditor } from '@/components'
import { Button } from '@/components'
import Toolbar from './Toolbar'
import { StyledEditor } from './style'
import { untilDestroyed } from '@/utils/operators'

interface IProps {
  path: string
  fileSize: number
  readOnly?: boolean
  onCancel: () => void
  onOk: () => void
  viewContent: (options: {
    path: string
    offset: number
    len: number
  }) => Promise<string>
  saveContent?: (options: { path: string; content: string }) => Promise<boolean>
  onBeforeRefresh?: () => Promise<number>
}

const CHUNK_SIZE = 1 * 1024 * 1024

@observer
export default class Editor extends React.Component<IProps> {
  @observable loading = false
  @observable saving = false
  @observable content = ''
  @observable isReversed = false
  @observable cursor = 0
  @action
  updateIsReversed = flag => (this.isReversed = flag)
  @action
  updateLoading = loading => (this.loading = loading)
  @action
  updateContent = content => (this.content = content)
  @action
  updateSaving = saving => (this.saving = saving)
  @action
  updateCursor = cursor => (this.cursor = cursor)

  editorRef = null

  componentDidMount() {
    const { readOnly } = this.props
    const { editor } = this.editorRef
    let editorHeight = editor.getLayoutInfo().height

    // reverse: fetch file content
    createMobxStream(() => this.isReversed)
      .pipe(untilDestroyed(this))
      .subscribe(this.freshen)

    // when readonly: fetch content by chunk
    if (readOnly) {
      editor.onDidScrollChange(e => {
        if (this.loading || this.fetchFinished) {
          return
        }

        // positive sequence
        if (!this.isReversed) {
          if (e.scrollTop > e.scrollHeight - editorHeight - 100) {
            this.fetchContent().then(() => {
              editorHeight = editor.getLayoutInfo().height
            })
          }
        } else {
          // inverted order
          if (e.scrollTop < 100) {
            const oldScrollHeight = e.scrollHeight
            this.fetchContent()
              .then(() => {
                editorHeight = editor.getLayoutInfo().height

                // update scrollHeight to locate old position
                const addedHeight = editor.getScrollHeight() - oldScrollHeight
                const scrollTop = editor.getScrollTop()
                editor.setScrollTop(scrollTop + addedHeight)
              })
              .catch(err => console.error(err))
          }
        }
      })
    }
  }

  @computed
  get fetchFinished() {
    return this.cursor * CHUNK_SIZE >= this.props.fileSize
  }

  fetchContent = (initial?: boolean) => {
    const { viewContent } = this.props

    this.updateLoading(true)

    let offset = this.isReversed
      ? this.props.fileSize - CHUNK_SIZE * (this.cursor + 1)
      : CHUNK_SIZE * this.cursor

    // fetch content anew
    if (initial) {
      this.updateContent('')
      this.updateCursor(0)
      offset = this.isReversed
        ? Math.max(this.props.fileSize - CHUNK_SIZE, 0)
        : 0
    }

    // update cursor
    this.updateCursor(this.cursor + 1)

    return viewContent({
      path: this.props.path,
      offset,
      len: CHUNK_SIZE,
    })
      .then(content => {
        let newContent = this.getEditorContent()
        if (this.isReversed) {
          if (initial) {
            newContent = content
          } else {
            newContent = content + newContent
          }
        } else {
          if (initial) {
            newContent = ''
          }
          newContent += content
        }

        this.updateContent(newContent)
      })
      .finally(() => this.updateLoading(false))
  }

  fetchFullContent = () => {
    const { viewContent } = this.props

    this.updateLoading(true)

    return viewContent({
      path: this.props.path,
      offset: 0,
      len: this.props.fileSize,
    })
      .then(this.updateContent)
      .finally(() => this.updateLoading(false))
  }

  getEditorContent = () => {
    return this.editorRef ? this.editorRef.getValue() : ''
  }

  private resetScrollTop = () => {
    const { editor } = this.editorRef

    // locate bottom
    if (this.isReversed) {
      const scrollHeight = editor.getScrollHeight()
      editor.setScrollTop(scrollHeight)
    } else {
      // locate top
      editor.setScrollTop(0)
    }
  }

  private freshen = async () => {

    const { readOnly, onBeforeRefresh } = this.props

    let newSize = null

    if (onBeforeRefresh) {
      newSize = await onBeforeRefresh()
    }
    // when readonly: fetch content by chunk
    if (readOnly) {
      // If fetch finished, don't fetch the content again
      if (!this.fetchFinished || newSize !== this.props.fileSize) {
        this.fetchContent(true).then(this.resetScrollTop)
      } else {
        this.resetScrollTop()
      }
    } else {
      // when edit: fetch content fully
      this.fetchFullContent().then(this.resetScrollTop)
    }
  }

  render() {
    const { readOnly = false, onCancel } = this.props

    return (
      <StyledEditor>
        <div>
          <Toolbar
            isReversed={this.isReversed}
            updateIsReversed={this.updateIsReversed}
            fetchContent={this.fetchContent}
            freshen={this.freshen}
            readOnly={readOnly}
            find={() => {
              this.editorRef &&
                this.editorRef.editor.trigger('', 'actions.find')
            }}
          />
        </div>

        <div className='editorMain'>
          {this.loading && (
            <div className='loading'>
              <Spin />
            </div>
          )}

          <CodeEditor
            ref={ref => (this.editorRef = ref)}
            readOnly={readOnly}
            value={this.content}
          />
        </div>

        {!readOnly && (
          <div className='footer'>
            <div className='footerMain'>
              <Button onClick={onCancel}>取消</Button>
              <Button
                disabled={this.saving}
                type='primary'
                onClick={this.onSave}>
                {this.saving ? '保存中...' : '保存'}
              </Button>
            </div>
          </div>
        )}
      </StyledEditor>
    )
  }

  private onSave = () => {
    const { saveContent } = this.props

    this.updateSaving(true)
    saveContent({
      path: this.props.path,
      content: this.getEditorContent(),
    })
      .then(() => {
        message.success('文件保存成功')
        this.props.onOk()
      })
      .finally(() => {
        this.updateSaving(false)
      })
  }
}
