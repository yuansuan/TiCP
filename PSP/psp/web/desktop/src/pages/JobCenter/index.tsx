/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import styled from 'styled-components'
import { Context, useModel, useStore } from './store'
import React, { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { Filter } from './Filter'
import { JobList } from './JobList'
import { useDidUpdate } from '@/utils/hooks'
import { Toolbar } from './Toolbar'
import { pageStateStore } from '@/utils'

const StyledLayout = styled.div`
  margin: 20px;
  background: white;

  .body {
    padding: 16px 20px;
  }
`

const JobCenter = observer(function JobCenter() {
  const store = useStore()

  useEffect(() => {
    store.refresh()
    // reserve query && page state
    pageStateStore.setByPath({
      query: store.query,
      pageIndex: store.pageIndex,
      pageSize: store.pageSize,
    })
  }, [store.query, store.pageIndex, store.pageSize])

  // 如果当前分页大于1且当前页不存在数据，则跳转查询第一页数据
  useDidUpdate(() => {
    if (store.model.list.length === 0 && store.pageIndex > 1) {
      store.setPageIndex(1)
    }
  }, [store.model.list.length])

  return (
    <StyledLayout>
      <Filter />
      <Toolbar />

      <div className='body'>
        <JobList />
      </div>
    </StyledLayout>
  )
})

export default function PageWithStore() {
  const model = useModel()

  return (
    <Context.Provider value={model}>
      <JobCenter />
    </Context.Provider>
  )
}
