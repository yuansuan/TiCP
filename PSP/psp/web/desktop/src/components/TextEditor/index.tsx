/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useRef, useEffect, useCallback } from 'react'
import styled from 'styled-components'
import { Modal, Button, CodeEditor, Mask, Icon } from '@/components'
import { useLocalStore, observer } from 'mobx-react-lite'
import { Toolbar } from './Toolbar'
import { message } from 'antd'
import { FileServer } from '@/server'

const StyledLayout = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;

  > .body {
    border: 1px solid ${({ theme }) => theme.borderColorBase};
    height: 100%;
    padding: 20px;
    padding-bottom: ${props => (props['data-readonly'] ? '20px' : '72px')};
    box-sizing: border-box;
    background-color: ${({ theme }) => theme.backgroundColorBase};
    position: relative;

    > .suspension {
      position: absolute;
      bottom: 85px;
      right: 40px;

      > div {
        display: flex;
        justify-content: center;
        align-items: center;
        width: 40px;
        height: 40px;
        margin: 5px 0;
        background-color: ${({ theme }) => theme.primaryColor};
        border-radius: 2px;

        > .anticon {
          color: white;
          font-size: 40px;

          &:hover {
            background-color: rgba(0, 113, 252, 0.93);
          }
        }
      }
    }
  }

  > .footer {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    padding: 10px 0;
    border-top: 1px solid ${({ theme }) => theme.borderColorBase};
  }
`

const CHUNK_SIZE = 1 * 1024 * 1024

type Props = {
  path: string
  fileInfo: any
  readonly?: boolean
  onCancel: () => void
  onOk: () => void
  boxServerUtil: FileServer // 请求util
}

export const TextEditor = observer(function Editor({
  fileInfo,
  readonly,
  onCancel,
  onOk,
  boxServerUtil
}: Props) {
  const ref = useRef(undefined)
  const state = useLocalStore(() => ({
    loading: false,
    file: fileInfo || null,
    reversed: true,
    cursor: 0,
    setFile(file) {
      this.file = fileInfo
    },
    setLoading(flag) {
      this.loading = flag
    },
    setCursor(cursor) {
      this.cursor = cursor
    },
    setReversed(flag) {
      this.reversed = flag
    },
    get fetchFinished() {
      return this.file ? this.cursor * CHUNK_SIZE >= this.file.size : false
    }
  }))
  const { loading, reversed, cursor, fetchFinished, file } = state

  // NOTE: 查看文件
  const fetchContent = useCallback(
    async function fetchContent(initial?: boolean) {
      try {
        state.setLoading(true)

        let offset = reversed
          ? file?.size - CHUNK_SIZE * (cursor + 1)
          : CHUNK_SIZE * cursor

        let newContent = getEditorContent()

        // fetch content anew
        if (initial) {
          newContent = ''
          state.setCursor(0)
          offset = reversed ? Math.max(file?.size - CHUNK_SIZE, 0) : 0
        }

        // update cursor
        state.setCursor(cursor + 1)

        const content = await boxServerUtil.getContent({
          path: file.path,
          offset,
          cross: file?.cross ? true : false,
          is_cloud: file?.is_cloud ? true : false,
          len: CHUNK_SIZE,
          user_name: file?.user_name
        })

        if (reversed) {
          newContent = content + newContent
        } else {
          newContent += content
        }

        setEditorContent(newContent)
      } finally {
        state.setLoading(false)
      }
    },
    [file, reversed, cursor]
  )

  useEffect(() => {
    if (file) {
      refresh()
    }
  }, [file, reversed])

  useEffect(() => {
    const { editor } = ref.current
    let editorHeight = editor.getLayoutInfo().height

    // when readonly: fetch content by chunk on scroll
    if (!readonly) {
      return undefined
    }
    const disposer = editor.onDidScrollChange(async e => {
      if (loading || fetchFinished) {
        return
      }

      // positive sequence
      if (!reversed) {
        if (e.scrollTop > e.scrollHeight - editorHeight - 100) {
          await fetchContent()
          editorHeight = editor.getLayoutInfo().height
        }
      } else {
        // inverted order
        if (e.scrollTop < 100) {
          const oldScrollHeight = e.scrollHeight
          await fetchContent()
          editorHeight = editor.getLayoutInfo().height

          // update scrollHeight to locate old position
          const addedHeight = editor.getScrollHeight() - oldScrollHeight
          const scrollTop = editor.getScrollTop()
          editor.setScrollTop(scrollTop + addedHeight)
        }
      }
    })

    return () => {
      disposer.dispose()
    }
  }, [loading, fetchFinished, reversed, fetchContent])

  function getEditorContent() {
    return ref.current?.getValue() || ''
  }

  function setEditorContent(content) {
    return ref.current.setValue(content || '')
  }

  async function fetchFullContent() {
    try {
      state.setLoading(true)

      const content = await boxServerUtil.getContent({
        path: file.path,
        offset: 0,
        cross: file?.cross ? true : false,
        is_cloud: file?.is_cloud ? true : false,
        length: file.size,
        user_name: file?.user_name
      })
      setEditorContent(content)
    } finally {
      state.setLoading(false)
    }
  }

  function resetScrollTop() {
    if(!ref?.current) return
    const { editor } = ref.current

    // locate bottom
    if (reversed) {
      const scrollHeight = editor.getScrollHeight()
      editor.setScrollTop(scrollHeight)
    } else {
      // locate top
      editor.setScrollTop(0)
    }
  }

  async function refresh() {
    // when readonly: fetch content by chunk
    if (readonly) {
      await fetchContent(true)
      resetScrollTop()
    } else {
      // when edit: fetch content fully
      await fetchFullContent()
      resetScrollTop()
    }
  }

  async function onConfirm() {
    let uploadParams = {
      file: new window.File([getEditorContent()], file.name),
      path: file.path
    }

    // uploadParams['isEdit'] = true
    await boxServerUtil.upload(uploadParams)
    message.success('保存成功')
    onOk()
  }

  function backward() {
    state.setReversed(true)
  }

  function forward() {
    state.setReversed(false)
  }

  return (
    <StyledLayout data-readonly={readonly}>
      <Modal.Toolbar>
        <Toolbar editor={ref.current?.editor} refresh={refresh} />
      </Modal.Toolbar>
      <div className='body'>
        {loading && <Mask.Spin />}
        <CodeEditor ref={ref} readOnly={readonly} />
        <div className='suspension'>
          <div>
            <Icon type='align_top' onClick={forward} />
          </div>
          <div>
            <Icon
              type='align_top'
              style={{ transform: 'rotateX(180deg)' }}
              onClick={backward}
            />
          </div>
        </div>
      </div>
      {!readonly && (
        <Modal.Footer
          className='footer'
          onCancel={onCancel}
          OkButton={
            <Button type='primary' onClick={onConfirm}>
              确认
            </Button>
          }
        />
      )}
    </StyledLayout>
  )
})

export function showTextEditor(props: Omit<Props, 'onCancel' | 'onOk'>) {
  return Modal.show({
    title: `${props.readonly ? '查看' : '编辑'}文件 ${props.path
      .split('/')
      .pop()}`,
    width: 800,
    className: '__fileEditor__',
    style: {
      padding: 0
    },
    bodyStyle: {
      height: 600,
      padding: 0
    },
    footer: null,
    content: ({ onCancel, onOk }) => (
      <TextEditor onCancel={onCancel} onOk={onOk} {...props} />
    )
  })
}
