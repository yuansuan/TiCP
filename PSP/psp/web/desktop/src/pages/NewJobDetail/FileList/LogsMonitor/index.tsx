/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useRef, useEffect } from 'react'
import styled from 'styled-components'
import { Modal, CodeEditor, Mask, Icon } from '@/components'
import { useLocalStore, observer } from 'mobx-react-lite'
import { Toolbar } from './Toolbar'
import { ConnectionFactory } from '@/utils/WebSocket'
import { Status } from './Status'

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

type Props = {
  path: string
  readonly?: boolean
  display_state: number
  id: string // jobId
  onCancel: () => void
  onOk: () => void
}

const host = location.host
const wsUrl = `wss://${host}/ws/monitor_job_output`
const factory = new ConnectionFactory(wsUrl)

export const LogsMonitor = observer(function Monitor({ readonly, id }: Props) {
  const logsRef = useRef('')
  let conn = null
  const ref = useRef(undefined)
  const state = useLocalStore(() => ({
    loading: false,
    file: null,
    reversed: true,
    cursor: 0,
    isOpenLogsModal: false,
    logsStatus: '',
    setLogsStatus(status) {
      this.logsStatus = status
    },
    setLoading(flag) {
      this.loading = flag
    },
    setReversed(flag) {
      this.reversed = flag
    }
  }))
  const { loading, reversed } = state

  useEffect(() => {
    try {
      state.setLoading(true)
      conn = factory.create(`job_id=${id}`)
      conn.onReceive(data => {
        if (data.endsWith('completed')) {
          state.setLogsStatus('close')
        }
        logsRef.current += data + '\n'

        if (!logsRef.current) {
          state.setLogsStatus('error')
          return undefined
        } else {
          state.setLogsStatus('success')
          setEditorContent(logsRef.current)
          // resetScrollTop()
        }
      })
      conn.onClose(() => {
        console.log('ws连接断开')
        state.setLogsStatus('close')
      })
    } finally {
      state.setLoading(false)
    }

    return () => {
      conn.close()
      logsRef.current = ''
      console.log('close ws conn when destory Modal')
    }
  }, [])

  useEffect(() => {
    // locate bottom
    if (reversed) {
      const scrollHeight = ref.current?.editor.getScrollHeight()
      ref.current?.editor.setScrollTop(scrollHeight)
    } else {
      // locate top
      ref.current?.editor.setScrollTop(0)
    }
  }, [reversed])

  function setEditorContent(content) {
    return ref.current?.setValue(content || '')
  }

  function resetScrollTop() {
    // locate bottom
    if (reversed) {
      const scrollHeight = ref.current?.editor.getScrollHeight()
      ref.current?.editor.setScrollTop(scrollHeight)
    } else {
      // locate top
      ref.current?.editor.setScrollTop(0)
    }
  }

  function backward() {
    state.setReversed(true)
  }

  function forward() {
    state.setReversed(false)
  }

  function refresh() {
    try {
      state.setLoading(true)
      conn = factory.create(`job_id=${id}`)
      conn.onReceive(data => {
        if (data.endsWith('completed')) {
          // 实时日志全部返回完成
          state.setLogsStatus('close')
        }
        logsRef.current += data + '\n'

        if (!logsRef.current) {
          state.setLogsStatus('error')
          return undefined
        } else {
          state.setLogsStatus('success')
          setEditorContent(logsRef.current)
          // resetScrollTop()
        }
      })
    } finally {
      state.setLoading(false)
    }
  }

  const clearScreen = async () => {
    setEditorContent('')
    logsRef.current = ''
  }

  return (
    <StyledLayout data-readonly={readonly}>
      <div
        style={{
          position: 'absolute',
          top: '28px',
          left: '10px'
        }}>
        <Status type={state.logsStatus} />
      </div>
      <Modal.Toolbar>
        <Toolbar
          editor={ref.current?.editor}
          refresh={refresh}
          clearScreen={clearScreen}
          showRefreshAction={state.logsStatus === 'close'}
        />
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
    </StyledLayout>
  )
})

export function showLogsMonitor(props: Omit<Props, 'onCancel' | 'onOk'>) {
  return Modal.show({
    title: `实时日志 ${props.path.split('/').pop()}`,
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
      <LogsMonitor onCancel={onCancel} onOk={onOk} {...props} />
    )
  })
}
