/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { Context, useModel, useStore } from './store'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Mask } from '@/components'
import styled from 'styled-components'
import { Info } from './Info'
import { getUrlParams } from '@/utils'
import FileMGT from '@/pages/FileMGT'
import { StatusStep } from './StatusStep'
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
    }
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
    const interval = setInterval(() => {
      store.refresh()
    }, 5000)
    return () => clearInterval(interval)
  }, [job])

  return state.initialLoading ? (
    <Mask.Spin></Mask.Spin>
  ) : (
    <>
      <Info />
      {job.timelines?.length > 0 && <StatusStep steps={job.timelines} />}
      <FileMGT job={job} />
    </>
  )
})

export default function JobDetailWithStore() {
  const id = getUrlParams()?.jobId
  const model = useModel({ id })

  return (
    <Context.Provider value={model}>
      <StyledLayout>
        <JobDetail />
      </StyledLayout>
    </Context.Provider>
  )
}
