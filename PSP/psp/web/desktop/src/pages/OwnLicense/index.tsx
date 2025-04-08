/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { useStore, useModel, Context } from './store'
import { List } from './List'
import { observer } from 'mobx-react-lite'
import { Toolbar } from './Toolbar'

const StyledLayout = styled.div`
  padding: 16px 20px;
  background-color: #ffffff;
  margin: 12px 16px;

  .main {
    padding-bottom: 50px;
    display: flex;
    flex-direction: column;
    background-color: #ffffff;
    .pagination {
      margin: 20px auto;
    }
  }
`

const ListPage = observer(function ListPage() {
  const store = useStore()

  useEffect(() => {
    store.fetch(store.params)
  }, [store.params])

  return (
    <StyledLayout>
      <Toolbar />
      <div className='main'>
        <List />
      </div>
    </StyledLayout>
  )
})

export default function ListPageWithStore() {
  const model = useModel()

  return (
    <Context.Provider value={model}>
      <ListPage />
    </Context.Provider>
  )
}
