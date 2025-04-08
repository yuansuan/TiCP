/* Copyright (C) 2016-present, Yuansuan.cn */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { AppList } from './AppList'
import { Tabs } from 'antd'
import { useStore } from './store'

const StyledLayout = styled.div`
  .ant-tabs-top > .ant-tabs-nav {
    margin: 0;
  }

  .ant-tabs-content {
    padding-top: 10px;
    background-color: white;
  }
`

const { TabPane } = Tabs

export const AppListTabs = observer(function AppListTabs() {
  const store = useStore()

  // 作业重提交时，根据软件的 is_trail 属性自动切换tab
  useEffect(() => {
    const is_trial = store.data.currentApp?.is_trial
    // 当软件数量为0，那么currentApp为空，此时不根据软件的 is_trial 属性自动切换tab
    if (store.data.currentApp === undefined) return

    if (is_trial && !store.is_trial) {
      store.setTabKey('trial')
    } else if (!is_trial && store.is_trial) {
      store.setTabKey('formal')
    }
  }, [store.currentAppId])

  return (
    <StyledLayout>
      <Tabs
        type='card'
        activeKey={store.tabKey}
        onChange={key => {
          store.setTabKey(key)
          store.data.currentApp =
            key === 'trial'
              ? store.apps.filter(item => item.is_trial)[0]
              : store.apps.filter(item => !item.is_trial)[0]
        }}>
        <TabPane key='formal' tab='正式软件'></TabPane>
        <TabPane key='trial' tab='试用软件'></TabPane>
      </Tabs>
      {store.tabKey === 'formal' && <AppList is_trial={false} />}
      {store.tabKey === 'trial' && <AppList is_trial={true} />}
    </StyledLayout>
  )
})
