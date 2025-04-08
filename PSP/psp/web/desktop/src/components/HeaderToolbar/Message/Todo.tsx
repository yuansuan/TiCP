/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { List, Empty } from 'antd'
import { ChromeOutlined } from '@ant-design/icons'
import { Button, Modal } from '@/components'
import { observer } from 'mobx-react-lite'
import { format } from 'timeago.js'
import { recordList } from '@/domain'
import { StyledPanel } from './style'
import {ShareFileContent} from './Share'
import {Http} from '@/utils'

interface IProps {
  hideDropdown: () => void
}

export const Todo = observer(function Todo({ hideDropdown }: IProps) {

  async function getShareItem (id) {
    hideDropdown()
    const {data} =await Http.get('/storage/share/get',{
      params: {
        id
      }
    })
    await Modal.show({
      title: '保存方式',
      footer:null,
      content: ({onOk,onCancel}) => (
        <ShareFileContent onOk={onOk} onCancel={onCancel} {...data}/>
      )
    })
  }
  return (
    <StyledPanel>
      <List
        bordered
        locale={{ emptyText: <Empty description='暂无未处理分享' /> }}
        dataSource={recordList.list}
        renderItem={item => (
          <List.Item key={item.id} className='item'>
            <List.Item.Meta
              avatar={<ChromeOutlined />}
              title={
                <div className='title'>
                  <div>通知</div>
                  <div className='time'>
                    {format(item.share_time.toString(), 'zh_CN')}
                  </div>
                  <div className='actions'>
                  {item.state === 2 ? (
                  <span className='isRead'>已读</span>
                  ) : (
                    <Button type='link' onClick={() => {
                      item.readShare(item.id)
                      getShareItem(item.id)
                      }}>
                      立即查看
                    </Button>
                  )}
                  </div>
                </div>
              }
              description={
                item?.content
              }
            />
          </List.Item>
        )}
      />
    </StyledPanel>
  )
})
