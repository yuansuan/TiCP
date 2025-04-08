/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { Page } from '@/components'
import { StyledLayout } from './style'
import { Members } from './Members'
import { Context, useModel, useStore } from './model'
import { Toolbar } from './Toolbar'

const Content = observer(function Content() {
  const store = useStore()

  useEffect(() => {
    store.fetch()
  }, [store.page_size, store.page_index, store.query])

  return (
    <Page header={null}>
      <StyledLayout>
        <div className='main'>
          <Toolbar />

          <div className='body'>
            <Members />
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
