/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer, useLocalStore } from 'mobx-react-lite'
import styled from 'styled-components'
import { noticeList, env } from '@/domain'
import TextLoop from 'react-text-loop'
import { Http } from '@/utils'
import { Modal } from '@/components'
import ReactMarkdown from 'react-markdown'
import 'github-markdown-css/github-markdown-light.css'

const StyledLayout = styled.div`
  .textLoop {
    width: 100%;
    > div {
      > div {
        width: 100%;
      }
    }
  }
`

const StyledMessage = styled.div`
  .markdown-body {
    box-sizing: border-box;
    min-width: 200px;
    max-width: 980px;
    margin: 0 auto;
  }
`
type IProps = {
  onOk: () => void
  onCancel: () => void
}

export const PopupNotice = observer(function Notice(props: IProps) {
  const { onCancel, onOk } = props

  const state = useLocalStore(() => ({
    interval: 30000,
    setInterval(interval) {
      this.interval = interval
    }
  }))

  function onMouseEnter() {
    state.setInterval(0)
  }

  function onMouseLeave() {
    state.setInterval(30000)
  }

  function neverNotify(noticeId) {
    return Http.post('/notice/report', {
      notice_id: noticeId
    })
  }

  // priority -> 0: '弹窗', 1: '横幅'
  const getPopupNotice = () => {
    return noticeList.list.filter(
      n => n.priority === 0 && n.company_ids.indexOf(env.company.id) !== -1
    )
  }

  return (
    <StyledLayout onMouseEnter={onMouseEnter} onMouseLeave={onMouseLeave}>
      {getPopupNotice().length > 0 && (
        <TextLoop
          className='textLoop'
          mask
          interval={state.interval}
          noWrap={false}>
          {getPopupNotice().map(banner => (
            <StyledMessage key={banner.id}>
              <div className='markdown-body'>
                <h3>{banner.title}</h3>
                <ReactMarkdown children={banner.content} />
              </div>
            </StyledMessage>
          ))}
        </TextLoop>
      )}
      <Modal.Footer
        onCancel={onCancel}
        cancelText={'关闭'}
        okText={'不再提醒'}
        onOk={async () => {
          try {
            await Promise.all(
              getPopupNotice().map(async b => neverNotify(b.id))
            )
          } finally {
            noticeList.fetch()
            onOk()
          }
        }}
      />
    </StyledLayout>
  )
})
