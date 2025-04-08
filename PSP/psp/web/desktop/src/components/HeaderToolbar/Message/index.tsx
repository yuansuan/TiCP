/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { Badge, Dropdown, Tabs, Tooltip } from 'antd'
import { useDispatch } from 'react-redux'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Icon } from '@/components'
import { lastMessages,recordList } from '@/domain'
import { StyledLayout, StyledOverlay } from './style'
import { Notifications } from './Notifications'
import { Todo } from './Todo'
import { buryPoint, history } from '@/utils'
import styled from 'styled-components'

const { TabPane } = Tabs

const Style = styled.div`
  flex: 1 0 auto;
`

export const Message = observer(function Message() {
  const state = useLocalStore(() => ({
    visible: false,
    hovered: false,
    setHovered(flag) {
      this.hovered = flag
    },
    setVisible(visible) {
      this.visible = visible
    },
    tipVisible: false,
    setTipVisible(visible) {
      this.tipVisible = visible
    },
    activeKey: 'notification',
    setActiveKey(key) {
      this.activeKey = key
    },
    get isEmpty() {
      return lastMessages.list.length === 0 && recordList.list.length === 0
    }
  }))
  const dispatch = useDispatch()
  const { visible, hovered, activeKey, isEmpty } = state

  function refresh() {
    lastMessages.fetchLast()
    recordList.fetchLast()

  }

  useEffect(() => {
    if (visible) {
      refresh()
    }
  }, [visible])

  // 展开下拉菜单，隐藏 tip
  useEffect(() => {
    if (visible) {
      state.setTipVisible(false)
    }
  }, [visible])

  function hideDropdown() {
    state.setVisible(false)
  }

  function showMore() {
    hideDropdown()
    window.localStorage.setItem('CURRENTROUTERPATH', '/messages?tab=messages')
    dispatch({
      type: 'MESSAGES',
      payload: 'togg'
    })
  }
  const readAll = async () => {
    if( activeKey === 'notification'){
      await lastMessages.readAll()
      await lastMessages.fetchLast()
    }else if(activeKey ==='todo'){
      await recordList.readAll()
      await recordList.fetchLast()
    }
  }

  const count = lastMessages.unreadCount + recordList.unhandledCount

  return (
    <Dropdown
      visible={visible}
      onVisibleChange={state.setVisible.bind(state)}
      placement='bottomRight'
      overlay={
        <StyledOverlay>
          <div className='body' onClick={e => e.stopPropagation()}>
            <Tabs
              activeKey={activeKey}
              onChange={state.setActiveKey.bind(state)}>
              <TabPane
                key='notification'
                tab={
                  <Badge count={lastMessages.unreadCount}>
                    <span className='tabName'>通知</span>
                  </Badge>
                }>
                <Notifications />
              </TabPane>
              <TabPane
                key='todo'
                tab={
                  <Badge count={recordList.unhandledCount}>
                    <span className='tabName'>分享通知</span>
                  </Badge>
                }>
                <Todo hideDropdown={hideDropdown} />
              </TabPane>
            </Tabs>
          </div>
          {!isEmpty && (
            <div className='footer'>
              <div className='read' onClick={readAll}>
                全部已读
              </div>
              <div className='link' onClick={showMore}>
                查看更多
              </div>
            </div>
          )}
        </StyledOverlay>
      }
      trigger={['click']}>
      <StyledLayout
        onMouseEnter={() => {
          state.setHovered(true)
        }}
        onMouseLeave={() => {
          state.setHovered(false)
        }}
        onClick={() => {
          buryPoint({
            category: '导航栏',
            action: '通知'
          })
        }}>
        <Tooltip
          onVisibleChange={visible => {
            state.setTipVisible(visible)
          }}
          visible={state.tipVisible}
          title='消息通知'>
          <Style>
            <Icon
              type={hovered || visible ? 'message_active' : 'message_default'}
            />
            <Badge offset={[4, 0]} count={count} />
          </Style>
        </Tooltip>
      </StyledLayout>
    </Dropdown>
  )
})
