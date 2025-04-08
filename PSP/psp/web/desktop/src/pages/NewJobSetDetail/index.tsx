/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useRef, useState } from 'react'
import { Context, useModel, useStore } from './store'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Mask } from '@/components'
import styled from 'styled-components'
import { Info } from './Info'
import { getUrlParams } from '@/utils'
import { JobList } from '../NewJobManager/JobList'
import { useModel as useJobListModel } from '../NewJobManager/JobList/store'
import { Context as JobListContext } from '../NewJobManager/JobList'

export const StyledLayout = styled.div`
  background: #fff;
  justify-content: center;
  align-items: center;
`

export const JobSetDetail = observer(function JobSetDetail() {
  const store = useStore()
  const ref = useRef(null)
  const jobListStore = useJobListModel(store)
  const [height, setHeight] = useState(720)

  const { jobSet, model } = store
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

      const resizeObserver = new ResizeObserver(entries => {
        for (let entry of entries) {
          setHeight(entry.contentRect.height)
        }
      })

      resizeObserver.observe(ref.current)

      // hack: 处理Table首次加载 bug
      setTimeout(() => {
        if (ref.current) {
          ref.current.style.paddingRight = 1 + 'px'
          ref.current.style.height = height - 120 + 'px'
        }
      }, 3000)

      return () => resizeObserver.disconnect()
    } finally {
      state.setInitialLoading(false)
    }
  }, [])

  useEffect(() => {
    const interval = setInterval(() => {
      store.refresh()
    }, 5000)
    return () => clearInterval(interval)
  }, [jobSet, model])

  return state.initialLoading ? (
    <Mask.Spin></Mask.Spin>
  ) : (
    <>
      <Info />
      <div ref={ref}>
        <JobListContext.Provider value={jobListStore}>
          <JobList showJobSetName={false} hidePagination height={height - 50} />
        </JobListContext.Provider>
      </div>
    </>
  )
})

export default function JobSetDetailWithStore() {
  const id = getUrlParams()?.jobSetId
  const model = useModel({ id })

  return (
    <Context.Provider value={model}>
      <StyledLayout>
        <JobSetDetail />
      </StyledLayout>
    </Context.Provider>
  )
}
