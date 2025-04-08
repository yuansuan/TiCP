/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useRef, useState } from 'react'
import styled from 'styled-components'
import { Provider, useStore } from './store'
import { observer } from 'mobx-react-lite'
import { Page } from '@/components'
import { useDidUpdate } from '@/utils/hooks'
import { JobList, Context, useModel } from './JobList'
import { Filter } from './Filter'

const StyledLayout = styled.div`
  height: calc(100vh - 100px);
  width: 100%;
  background: #fff;

  .areaSelectWrap {
    display: none;
    padding: 10px 20px;
    border-bottom: 6px solid #f5f5f5;
    background: #fff;
    > div {
      display: flex;
      align-items: center;
    }
  }
`

const JobListPage = observer(function JobListPage(props) {
  const store = useStore()
  const jobListStore = useModel(store)
  const ref = useRef(null)
  const [height, setHeight] = useState(800)

  useEffect(() => {
    store.refresh()
    const resizeObserver = new ResizeObserver(entries => {
      for (let entry of entries) {
        setHeight(entry.contentRect.height)
      }
    })

    resizeObserver.observe(ref.current)

    // hack: 处理Table首次加载 bug
    setTimeout(() => {
      if (ref.current) ref.current.style.paddingRight = 1 + 'px'
    }, 3000)
    return () => resizeObserver.disconnect()
  }, [])

  useEffect(() => {
    store.refresh()
  }, [store.query, store.pageIndex, store.pageSize])

  // 如果当前分页大于1且当前页不存在数据，则跳转查询第一页数据
  useDidUpdate(() => {
    if (store.model.list.length === 0 && store.pageIndex > 1) {
      store.setPageIndex(1)
    }
  }, [store.model.list.length])

  useEffect(() => {
    let handler = setInterval(() => {
      store.refresh()
    }, 5000)
    return () => clearTimeout(handler)
  }, [store.model.list])

  return (
    <StyledLayout id='job_manager' ref={ref}>
      {/* 中间的筛选区 */}
      <Page header={null}>
        <Filter />
        <Context.Provider value={jobListStore}>
          {/* 3个按钮+表格+分页 */}
          <JobList height={height - 200} />
        </Context.Provider>
      </Page>
    </StyledLayout>
  )
})

export default function PageWithStore(props) {
  return (
    <Provider>
      <JobListPage {...props} />
    </Provider>
  )
}
