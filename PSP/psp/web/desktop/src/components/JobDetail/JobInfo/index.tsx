import * as React from 'react'
import { Wrapper } from './style'
import { lsf, transformer, pbspro, remotepbs } from './fieldName'
import InfoList from './InfoList'
// import { JobDetail } from '@/domain/JobList'
import { sysConfig } from '@/domain'

interface JobInfoProps {
  jobInfo: any
}

export default class JobInfo extends React.Component<JobInfoProps> {
  render() {
    const { jobInfo } = this.props
    if (!jobInfo) return null

    return (
      <Wrapper>
        <InfoList
          template={
            jobInfo.destClusterName
              ? remotepbs
              : sysConfig.schedulerName === 'lsf'
              ? lsf
              : pbspro
          }
          dataSource={jobInfo}
          transformer={transformer}
        />
      </Wrapper>
    )
  }
}
