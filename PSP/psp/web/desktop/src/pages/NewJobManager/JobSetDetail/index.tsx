/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import { InfoBlock, InfoItem, JobListWrapper, Wrapper } from './style'
import { observer } from 'mobx-react-lite'
import { useStore, Provider } from './store'
import { JobList, Context, useModel } from '../JobList'
import { getUrlParams } from '@/utils/Validator'
const JobSetDetail = observer(function JobSetDetail() {
  const params = getUrlParams()
  const store = useStore()
  const jobListStore = useModel(store)
  const { jobSet } = store

  useEffect(() => {
    store.refresh()
  }, [store.pageIndex, store.pageSize])

  // useEffect(() => {
  //   let handler = setTimeout(function pull() {
  //     store.refresh().finally(() => {
  //       handler = setTimeout(pull, 5000)
  //     })
  //   }, 5000)

  //   return () => clearTimeout(handler)
  // }, [])

  return (
    <Wrapper>
      <InfoBlock>
        <div className='header'>
          <div className='title'>基本信息</div>
        </div>
        <div className='content'>
          <InfoItem title={jobSet.name}>作业集名称：{jobSet.name}</InfoItem>
          <InfoItem title={jobSet.count + ''}>作业数：{jobSet.count}</InfoItem>
          <InfoItem title={jobSet.user_name}>
            创建者：{jobSet.user_name}
          </InfoItem>
          <InfoItem title={jobSet.is_batch_job ? '是' : '否'}>
            是否批量作业：{jobSet.is_batch_job ? '是' : '否'}
          </InfoItem>
          <InfoItem title={jobSet.finish_time?.toString()}>
            完成时间：{jobSet.finish_time?.toString()}
          </InfoItem>
        </div>
      </InfoBlock>

      <JobListWrapper>
        <div className='item'>
          <Context.Provider value={jobListStore}>
            <JobList showJobSetName={false} />
          </Context.Provider>
        </div>
      </JobListWrapper>
    </Wrapper>
  )
})

export default function JobSetDetailWithStore() {
  return (
    <Provider>
      <JobSetDetail />
    </Provider>
  )
}
