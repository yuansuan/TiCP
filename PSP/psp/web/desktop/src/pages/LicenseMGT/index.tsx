/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { Page } from '@/components'
import { StyledLayout } from './style'
import { Context, useModel } from './model'
// import { LoginForm } from './Form'

const Content = observer(function Content() {
  useEffect(() => {
    const url = 'https://simplant.t3caic.com/_lm/#/login?auto=true'
    window.open(url, '_blank')
  }, [])

  return (
    <Page header={null}>
      <StyledLayout>
        <div className='main'>欢迎使用许可证管理系统</div>
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
