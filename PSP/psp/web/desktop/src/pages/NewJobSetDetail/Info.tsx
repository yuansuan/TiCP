/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { useStore } from './store'
import { getComputeType, getDisplayRunTime } from '@/utils'

const StyledLayout = styled.div`
  padding: 0 20px 0 20px;
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
    color: #000;
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
  max-width: 300px;
`

export const Info = observer(function Info() {
  const store = useStore()
  const { jobSet } = store

  return (
    <StyledLayout>
      <div className='header'>
        <div className='title'>基本信息</div>
      </div>
      <div className='content'>
        <InfoItem>项目编号: {jobSet.project_id}</InfoItem>
        <InfoItem>项目名称: {jobSet.project_name}</InfoItem>
        <InfoItem>作业类型: {getComputeType(jobSet.job_type)}</InfoItem>
        <InfoItem>作业集编号: {jobSet.job_set_id}</InfoItem>
        <InfoItem>作业集名称: {jobSet.job_set_name}</InfoItem>
        <InfoItem>
          运行时长：{getDisplayRunTime(Number(jobSet?.exec_duration))}
        </InfoItem>
        <InfoItem>应用编号: {jobSet.app_id}</InfoItem>
        <InfoItem>应用名称: {jobSet.app_name}</InfoItem>
        <InfoItem>作业数量: {jobSet.job_count}</InfoItem>
        <InfoItem>成功数量: {jobSet.success_count}</InfoItem>
        <InfoItem>失败数量: {jobSet.failure_count}</InfoItem>
        <InfoItem>用户编号: {jobSet.user_id}</InfoItem>
        <InfoItem>用户名称: {jobSet.user_name}</InfoItem>
        <InfoItem>开始时间: {jobSet?.start_time ?? '--'}</InfoItem>
        <InfoItem>结束时间: {jobSet?.end_time ?? '--'}</InfoItem>
      </div>
    </StyledLayout>
  )
})
