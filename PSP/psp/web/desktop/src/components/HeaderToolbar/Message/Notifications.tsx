/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { List, Empty } from 'antd'
import { observer, useLocalStore } from 'mobx-react-lite'
import { format } from 'timeago.js'
import { lastMessages } from '@/domain'
import { StyledPanel } from './style'

export const Notifications = observer(function Notifications() {
  const { dataSource } = useLocalStore(() => ({
    get dataSource() {
      return lastMessages.list.map(item => ({
        ...item,
        body: item.body,
        message: item.message
      }))
    }
  }))

  async function read(item) {
    await item.read()
    lastMessages.fetchUnreadCount()
  }

  function description(item) {
    return item.message
  }

  return (
    <StyledPanel>
      <List
        bordered
        locale={{ emptyText: <Empty description='暂无消息通知' /> }}
        dataSource={dataSource}
        renderItem={item => (
          <List.Item key={item.id} className='item'>
            <List.Item.Meta
              title={
                <div className='title'>
                  <div>{item?.title}</div>
                  <div className='time'>
                    {format(item.create_time.toString(), 'zh_CN')}
                  </div>
                  <div className='actions'>
                    {item.state === 2 ? (
                      <span className='isRead'>已读</span>
                    ) : (
                      <span className='notRead' onClick={() => read(item)}>
                        标为已读
                      </span>
                    )}
                  </div>
                </div>
              }
              description={description(item)}
            />
          </List.Item>
        )}
      />
    </StyledPanel>
  )
})
