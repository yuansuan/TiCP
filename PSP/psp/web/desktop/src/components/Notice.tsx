/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { observer, useLocalStore } from 'mobx-react-lite'
import styled from 'styled-components'
import { noticeList, env } from '@/domain'
import { Alert } from 'antd'
import TextLoop from 'react-text-loop'
import { Http } from '@/utils'

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
  display: flex;

  > .right {
    margin: 0 10px;
    margin-left: auto;

    .link {
      text-decoration: underline;
      cursor: pointer;
    }
  }
`

export const Notice = observer(function Notice() {
  const state = useLocalStore(() => ({
    interval: 6000,
    setInterval(interval) {
      this.interval = interval
    },
    loading: false,
    setLoading(flag) {
      this.loading = flag
    }
  }))

  function onMouseEnter() {
    state.setInterval(0)
  }

  function onMouseLeave() {
    state.setInterval(6000)
  }

  async function neverNotify(noticeId) {
    if (state.loading) {
      return
    }

    try {
      state.setLoading(true)
      await Http.post('/notice/report', {
        notice_id: noticeId
      })
      await noticeList.fetch()
    } finally {
      state.setLoading(false)
    }
  }

  // priority -> 0: '弹窗', 1: '横幅'
  const getTopBarNotice = () => {
    return noticeList.list.filter(
      n => n.priority === 1 && n.company_ids.indexOf(env.company?.id) !== -1
    )
  }

  return (
    <StyledLayout>
      {getTopBarNotice().length > 0 && (
        <Alert
          // hack: trigger remount when noticeList's length change
          key={getTopBarNotice().length}
          closable={true}
          banner={true}
          onMouseEnter={onMouseEnter}
          onMouseLeave={onMouseLeave}
          message={
            <TextLoop className='textLoop' mask interval={state.interval}>
              {getTopBarNotice().map(banner => (
                <StyledMessage>
                  <div>{`${banner.title}：${banner.content}`}</div>
                  <div className='right'>
                    <span
                      className='link'
                      onClick={() => neverNotify(banner.id)}>
                      不再提醒
                    </span>
                  </div>
                </StyledMessage>
              ))}
            </TextLoop>
          }
        />
      )}
    </StyledLayout>
  )
})
