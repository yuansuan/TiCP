/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { useStore } from './store'
import { JobStatus } from '@/components'
import { scList } from '@/domain'

const StyledLayout = styled.div`
  padding: 0 20px 20px;
  background: #fff;

  .header {
    height: 46px;
    line-height: 46px;
    font-size: 14px;
    color: #333;
    border-bottom: 1px solid ${props => props.theme.borderColorBase};

    .title {
      width: 88px;
      text-align: center;
      font-weight: 500;
      position: relative;
    }
  }

  .content {
    display: grid;
    grid-template-columns: repeat(auto-fill, 246px);
    grid-row-gap: 12px;
    overflow-x: hidden;
    padding: 16px;

    > .job-status {
      display: flex;
    }
  }
`

const InfoItem = styled.div`
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  align-items: center;
  max-width: 240px;
`

export const Info = observer(function Info() {
  const store = useStore()
  const { job } = store
  const state = useLocalStore(() => ({
    get appTier() {
      return scList.getName(job?.sc_id) || '--'
    },
  }))

  return (
    <StyledLayout>
      <div className='header'>
        <div className='title'>基本信息</div>
      </div>
      <div className='content'>
        <InfoItem title={job.name}>名称：{job.name}</InfoItem>
        <InfoItem title={job.app_name}>应用：{job.app_name}</InfoItem>
        <InfoItem>创建人：{job.user_name}</InfoItem>
        <InfoItem className={'job-status'}>
          作业状态：
          <JobStatus job={job as any} />
        </InfoItem>
        <InfoItem>算力资源：{state.appTier}</InfoItem>
        <InfoItem>节点数：{job.resource_usage.nodes}</InfoItem>
        <InfoItem>核数：{job.resource_usage.cpus}</InfoItem>
        <InfoItem>运行时长：{job.displayRunTime}</InfoItem>
        <InfoItem>提交时间：{job.create_time?.toString()}</InfoItem>
        <InfoItem>开始时间：{job?.runtime?.start_time?.toString()}</InfoItem>
        <InfoItem>结束时间：{job?.runtime?.end_time?.toString()}</InfoItem>
      </div>
    </StyledLayout>
  )
})
