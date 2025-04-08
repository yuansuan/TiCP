/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { Page } from '@/components'
import { StyledLayout } from './style'
import { Context, useModel, useStore } from './model'
import { Departments } from './Departments'
import { Toolbar } from './ToolBar'

const Content = observer(function Content() {
  const store = useStore()

  useEffect(() => {
    store.fetch()
  }, [store.list.page_size, store.list.page_index, store.list.name])

  return (
    <Page header={null}>
      <StyledLayout>
        <div className='main'>
          <Toolbar />
          <div className='body'>
            <Departments />
          </div>
        </div>
      </StyledLayout>
    </Page>
  )
})

export default observer(function Index() {
  const store = useModel()

  return (
    <Context.Provider value={store}>
      <Content />
    </Context.Provider>
  )
})
