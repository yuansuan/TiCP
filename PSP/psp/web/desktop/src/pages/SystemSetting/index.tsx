/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useState } from 'react'
import { Page } from '@/components'
import General from './General'
import { WechatBindSection } from './WechatBindSection'

export default function Index() {
  const [key, setKey] = useState('1')
  return (
    <Page
      header={null}
      tabConfig={{
        tabContentList: [
          {
            tabName: '设置',
            tabKey: '1',
            content: key === '1' && <General />,
          },
          {
            tabName: '余额通知',
            tabKey: '2',
            content: key === '2' && <WechatBindSection />,
          },
        ],
        defaultActiveKey: key,
        onChange: v => setKey(v),
      }}
    />
  )
}
