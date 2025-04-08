import React, { useRef, useState, useEffect } from 'react'
import { Wrapper, ActionButtonWrapper, JobBaseInfoWrapper } from './style'

import { Tooltip } from 'antd'
import { lsfJobBaseInfoFields, pbsJobBaseInfoFields } from './common'
import { hasSimpleJobState } from '../../JobList/common'
import { JobActionButtonGroup } from '@/components'
import { getJobStatus } from '../../JobList/common'

import { isContentOverflow, Http } from '@/utils'
import { sysConfig } from '@/domain'
import { currentUser } from '@/domain'

function Value({ jobInfo, infoKey }) {
  const [isOverflow, setIsOverflow] = useState(false)
  const divEl = useRef(null)

  useEffect(() => {
    setIsOverflow(isContentOverflow(divEl.current, jobInfo[infoKey]))
  })

  return (
    <Tooltip title={isOverflow ? jobInfo[infoKey] : undefined}>
      <div ref={divEl} className='value'>
        {infoKey === 'jobStatus' ? (
          getJobStatus(jobInfo)
        ) : (
          <>{jobInfo[infoKey] || '--'}</>
        )}
      </div>
    </Tooltip>
  )
}

export default class JobSummary extends React.Component<any, any> {
  // feature: "monitor"
  startVNCMonitor = async cwd => {
    return Http.get('/vnc/monitor', {
      params: {
        cwd,
        username: currentUser.name
      }
    })
  }

  startVNC = async (rowData, startMonitor = false) => {
    const { jobId: id, jobDir: cwd } = rowData
    if (startMonitor) {
      await this.startVNCMonitor(cwd)
      // 判断是否
      await this.openVNC(id, cwd)
    } else {
      await this.openVNC(id, cwd)
    }
  }

  openVNC = async (id, cwd) => {
    window.open(`/vnc.html?type=job&id=${id}&cwd=${cwd}`, `${id}`)
  }

  render() {
    const { jobInfo, reSubmit } = this.props

    return (
      <Wrapper>
        <ActionButtonWrapper>
          <JobActionButtonGroup
            isHistory={this.props.isHistory}
            selectedItems={[jobInfo]}
            buttonsOption={{
              reSubmit
            }}
            reSubmitCallback={this.props.reSubmitCallback}
            operateCallback={this.props.updateJobDetailsInfo}
          />
        </ActionButtonWrapper>
        <JobBaseInfoWrapper>
          {Object.entries(
            sysConfig.schedulerName === 'lsf'
              ? lsfJobBaseInfoFields
              : pbsJobBaseInfoFields
          ).map(([key, label]) => {
            return (
              <div className='field' key={key}>
                <div className='label'>{label}</div>
                {key === 'jobId' ? (
                  <>
                    <div> {jobInfo[key]}</div>
                    {jobInfo.vncId > 0 &&
                      hasSimpleJobState(jobInfo, state => state === 'Running')}
                    {jobInfo.feature === 'monitor' &&
                      hasSimpleJobState(jobInfo, state => state === 'Running')}
                  </>
                ) : (
                  <Value jobInfo={jobInfo} infoKey={key} />
                )}
              </div>
            )
          })}
        </JobBaseInfoWrapper>
      </Wrapper>
    )
  }
}
