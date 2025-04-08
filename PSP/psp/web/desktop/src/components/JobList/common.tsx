import React from 'react'
import styled from 'styled-components'
import { Tooltip, Progress } from 'antd'
import { StatsBall } from '@/components'
import {
  jobStatusColumnFields,
  jobStatusColumnMap,
  CloudJobStatusMap,
  CloudJobStatus,
} from '@/domain/JobList'

interface IJobStatusWrapper {
  width: number
}

const JobStatusWrapper = styled.div<IJobStatusWrapper>`
  display: flex;
  width: ${props => props.width}px;
`

interface ISubJobStatus {
  color: string
  width: number
}

// TODO ARRAY JOB 和 SUB_ARRAY JOB 是否支持 VNC
export const hasSimpleJobState = (job, callback) => {
  if (job.jobType === 'JOB') {
    const jobStatus = jobStatusColumnFields.filter(
      o => job._jobStatus[o.key] === 1
    )
    return callback(jobStatus[0]?.status)
  } else {
    return false
  }
}

export const hasJobState = (job, callback) => {
  if (job.jobType === 'JOB' || job.jobType === 'SUB_ARRAY') {
    const jobStatus = jobStatusColumnFields.filter(
      o => job._jobStatus[o.key] === 1
    )
    return callback(jobStatus[0]?.status)
  } else {
    // array job 类似于操作多个 Job
    return true
  }
}

const SubJobStatus = styled.div<ISubJobStatus>`
  background: ${props => props.color};
  width: ${props => props.width}px;
  height: 15px;
`

// new design
function ArrayJobStatus(props) {
  let { run, susp, pend, exited, done } = props
  let sum = run + susp + pend + exited + done

  let tipMsg = jobStatusColumnFields.map(
    o => `${jobStatusColumnMap[o.key].label}: ${props[o.key]}`
  )

  return (
    <Tooltip title={tipMsg.join(', ')}>
      <JobStatusWrapper width={140}>
        {jobStatusColumnFields.map(o => (
          <SubJobStatus
            key={o.key}
            color={o.color}
            width={(props[o.key] * 140) / sum}
          />
        ))}
      </JobStatusWrapper>
    </Tooltip>
  )
}

function TableStatusBall(props) {
  const { color, label } = props
  return (
    <div
      style={{
        width: '150px',
        height: 28,
        lineHeight: '28px',
      }}>
      <StatsBall color={color}>{label}</StatsBall>
    </div>
  )
}

export const getJobStatus = rowData => {
  // 非 Array Job
  if (rowData.jobType === 'JOB' || rowData.workloadType === 'SUB_ARRAY') {
    // 云作业
    if (rowData.cloudJobStatus !== null) {
      if (rowData.cloudJobStatus === CloudJobStatus.init) {
        return (
          <TableStatusBall
            color={'#ccc'}
            label={CloudJobStatusMap[rowData.cloudJobStatus]}
          />
        )
      } else if (
        rowData.cloudJobStatus === CloudJobStatus.uploading ||
        rowData.cloudJobStatus === CloudJobStatus.upload_done
      ) {
        const percent = (rowData.uploadProgress * 100) / rowData.uploadTotalSize
        return (
          <div
            style={{
              width: '150px',
              height: 28,
              padding: '0 20px 0 10px',
              lineHeight: '28px',
            }}>
            <Progress percent={Number(percent.toFixed(1))} size='small' />
          </div>
        )
      } else {
        // 已经开始上云，直接使用 job_status
        // 由于时间差，job_status 有可能是 Suspended
        if (rowData.jobStatus === 'Suspended') rowData.jobStatus = 'Pending'

        const jobStatus = jobStatusColumnFields.filter(
          o => o.status === rowData.jobStatus
        )

        return (
          <TableStatusBall
            color={jobStatus[0]?.color}
            label={jobStatus[0]?.label}
          />
        )
      }
    }

    // 本地作业
    const jobStatus = jobStatusColumnFields.filter(
      o => rowData._jobStatus[o.key] === 1
    )

    return (
      <TableStatusBall
        color={jobStatus[0]?.color}
        label={jobStatus[0]?.label}
      />
    )
  } else {
    // array job 给一个分布占比
    return <ArrayJobStatus {...rowData._jobStatus} />
  }
}
