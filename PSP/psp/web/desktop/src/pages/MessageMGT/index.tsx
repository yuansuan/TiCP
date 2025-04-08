/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState } from 'react'
import styled from 'styled-components'
import { Tabs, Input } from 'antd'
import { Page } from '@/components'
import { Messages } from './Messages'
import { getUrlParams } from '@/utils/Validator'
import {ShareList} from './ShareList'

const { TabPane } = Tabs

const StyledLayout = styled.div`
  .ant-tabs-nav-wrap {
    padding: 20px 0 0 20px;
  }
`

export default function MessageMGT() {
  const [tab, updateTab] = useState(
    (getUrlParams().tab as string) || 'messages'
    )
    const [searchKey, setSearchKey] = useState('')

  return (
    <Page header={null}>
      <StyledLayout>
        <Tabs
          activeKey={tab}
          onChange={updateTab}
          tabBarExtraContent={
            <Input.Search
              placeholder='请输入关键词搜索'
              onChange={e => setSearchKey(e.target.value)}
            />
          }>
          <TabPane key='messages' tab='消息通知'>
            <Messages  visible={tab === 'messages'} searchKey={searchKey} />
          </TabPane>
          <TabPane key='share' tab='分享通知'>
            <ShareList visible={tab === 'share'}  searchKey={searchKey}/>
          </TabPane>
        </Tabs>
      </StyledLayout>
    </Page>
  )
}
