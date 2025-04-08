/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { webConfig } from '@/domain'
import { StatisticsComponent } from './Statistics'
import { LiveChatComponent } from './LiveChat'
import { observer } from 'mobx-react-lite'

export const WebConfigComponent = observer(function WebConfigComponent() {
  return (
    <>
      {webConfig.statisticsEnabled && <StatisticsComponent />}
      {webConfig.liveChatId && <LiveChatComponent />}
    </>
  )
})
