/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react-lite'
import { useStore } from './store'
import { Tooltip } from 'antd'
import { JobStatus, Status } from '@/components'
import { DATASTATEMAP } from '@/constant'
import { getDisplayRunTime, getComputeType } from '@/utils'

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
  const { job } = store

  return (
    <StyledLayout>
      <div className='header'>
        <div className='title'>基本信息</div>
      </div>
      <div className='content'>
        <InfoItem title={job.name}>作业名称：{job.name}</InfoItem>
        <InfoItem title={job.app_name}>应用名称：{job.app_name}</InfoItem>
        <InfoItem>用户名称：{job.user_name}</InfoItem>
        <InfoItem className={'job-status'}>
          计算状态：
          <JobStatus job={job as any} />
        </InfoItem>
        <InfoItem>作业类型：{getComputeType(job.type)}</InfoItem>
        <InfoItem>
          调度编号：
          {job.isCloud ? job?.out_job_id ?? '--' : job?.real_job_id ?? '--'}
        </InfoItem>
        <InfoItem>分配核数：{job.cpus_alloc ?? '--'}</InfoItem>
        <InfoItem>分配内存：{job.mem_alloc ?? '--'}</InfoItem>
        {!job.isCloud && <InfoItem>队列名称：{job.queue ?? '--'}</InfoItem>}
        <InfoItem>
          项目名称：{job?.project_name ?? '--'}
        </InfoItem>
        {!job.isCloud && (
          <InfoItem>
            运行节点：
            <Tooltip title={job?.exec_hosts}>{job.exec_hosts ?? '--'}</Tooltip>
          </InfoItem>
        )}
        {!job.isCloud && (
          <InfoItem>节点总数：{job.exec_host_num ?? '--'}</InfoItem>
        )}
        <InfoItem>
          运行时长：{getDisplayRunTime(Number(job?.exec_duration))}
        </InfoItem>
        {job.isCloud && (
          <InfoItem className={'job-status'}>
            数据状态：
            <Status
              text={DATASTATEMAP[job?.data_state]?.text}
              type={DATASTATEMAP[job?.data_state]?.type}
            />
          </InfoItem>
        )}
        <InfoItem>
          工作目录：
          <Tooltip title={job?.work_dir}>{job.work_dir ?? '--'}</Tooltip>
        </InfoItem>
        {!job.isCloud && <InfoItem>退出码：{job.exit_code ?? '--'}</InfoItem>}
        <InfoItem>提交时间：{job?.submit_time ?? '--'}</InfoItem>
        <InfoItem>运行时间：{job?.start_time ?? '--'}</InfoItem>
        <InfoItem>结束时间：{job?.end_time ?? '--'}</InfoItem>
      </div>
    </StyledLayout>
  )
})
