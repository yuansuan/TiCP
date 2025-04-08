/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { Context, useModel, useStore } from './store'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Mask } from '@/components'
import styled from 'styled-components'
import { Info } from './Info'
import { Empty } from 'antd'
import { FileList } from './FileList'
import { useParams } from 'react-router'

export const StyledLayout = styled.div`
  padding: 20px;
`

export const StyledEmpty = styled.div`
  height: 200px;
  display: flex;
  justify-content: center;
  align-items: center;
`

export const JobDetail = observer(function JobDetail() {
  const store = useStore()
  const { job } = store
  const state = useLocalStore(() => ({
    initialLoading: false,
    setInitialLoading(flag) {
      this.initialLoading = flag
    },
  }))

  useEffect(() => {
    try {
      state.setInitialLoading(true)
      store.refresh()
    } finally {
      state.setInitialLoading(false)
    }
  }, [])

  useEffect(() => {
    const interval = setInterval(() => store.refresh(), 5000)

    return () => clearInterval(interval)
  }, [])

  return state.initialLoading ? (
    <Mask.Spin></Mask.Spin>
  ) : job.is_deleted ? (
    <StyledEmpty>
      <Empty description='作业已删除' />
    </StyledEmpty>
  ) : (
    <>
      <Info />
      <FileList />
    </>
  )
})

export default function JobDetailWithStore() {
  const { id } = useParams<{ id: string }>()
  const model = useModel({ id })

  return (
    <Context.Provider value={model}>
      <StyledLayout>
        <JobDetail />
      </StyledLayout>
    </Context.Provider>
  )
}
