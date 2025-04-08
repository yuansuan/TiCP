/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import { Tabs } from 'antd'
import { Header } from './Header'
import { StyledLayout, StyledHeader, StyledContent } from './style'

const { TabPane } = Tabs

interface ITabContent {
  tabName: string
  tabKey?: string
  content: React.ReactNode | null
}

interface IProps {
  header?: React.ReactNode | null
  children?: React.ReactNode | null
  contentStyle?: React.CSSProperties
  tabConfig?: {
    tabContentList: ITabContent[]
    defaultActiveKey?: string
    onChange?: (activeKey: string) => void
  }
}

export const Page = (props: IProps) => {
  const {
    header = (
      <StyledHeader>
        <Header />
      </StyledHeader>
    ),
    tabConfig
  } = props

  return (
    <StyledLayout>
      {header}
      <StyledContent style={props?.contentStyle}>
        {tabConfig && tabConfig.tabContentList ? (
          <Tabs
            tabBarStyle={{
              height: '44px'
            }}
            defaultActiveKey={tabConfig.defaultActiveKey}
            onChange={tabConfig.onChange}>
            {tabConfig.tabContentList.map(tabContent => (
              <TabPane
                tab={tabContent.tabName}
                key={tabContent.tabKey || Math.random().toString()}>
                {tabContent.content}
              </TabPane>
            ))}
          </Tabs>
        ) : (
          props.children
        )}
      </StyledContent>
    </StyledLayout>
  )
}
