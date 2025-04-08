/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Pagination } from 'antd'
import { observer } from 'mobx-react-lite'
import { Context, useModel, useStore } from './store'
import { Toolbar } from './toolbar'
import { List } from './List'

const StyledLayout = styled.div`
  margin: 20px 20px;
  padding: 20px 20px;
  background-color: #ffffff;

  .main {
    padding: 20px 0;
    display: flex;
    flex-direction: column;
    margin-bottom: 20px;

    .pagination {
      margin: 20px 0 20px auto;
    }
  }
`

const BillForUserList = observer(function BillForUserList() {
  const store = useStore()
  const { model } = store
  const { page_ctx } = model

  useEffect(() => {
    store.fetch()
  }, [store.queryKey, store.pageSize, store.pageIndex])

  function onPageChange(index, size) {
    store.update({
      pageIndex: index,
      pageSize: size
    })
  }
  return (
    <StyledLayout>
      <Toolbar />
      <div className='main'>
        <List />
        {page_ctx?.total > 0 && (
          <Pagination
            className='pagination'
            showQuickJumper
            showSizeChanger
            pageSize={store.pageSize}
            current={store.pageIndex}
            total={page_ctx?.total}
            onChange={onPageChange}
          />
        )}
      </div>
    </StyledLayout>
  )
})

export default function Page() {
  const model = useModel()

  return (
    <Context.Provider value={model}>
      <BillForUserList />
    </Context.Provider>
  )
}
