import * as React from 'react'
import { Tabs } from 'antd'

import { TabAreaWrapper } from './style'
import JobFileSystem from './JobFileSystem'
import JobSummary from './JobSummary'
import SubJobList from './SubJobList'
import JobInfo from './JobInfo'
import { Title, JobNameWrapper } from './JobSummary/style'
import { Icon } from '@/components'
import { historyJobList, jobList } from '@/domain/JobList'
import { observer } from 'mobx-react'
import { observable } from 'mobx'
import { eventEmitter } from '@/utils'

interface IProps {
  goBack?: () => void
  isHistory?: boolean
  isSubJob?: boolean
  jobId: string
  reSubmitCallback?: () => void
  activeTabKey?: string
  jobType?: string
  reSubmit?: boolean
}

let intervalId = null

const { TabPane } = Tabs

@observer
export default class JobDetail extends React.Component<IProps, any> {
  @observable activeTabkey = this.props.activeTabKey || '1'
  @observable jobInfo: any

  private _tabs = [
    {
      label: '基本信息',
      key: '1',
      component: JobInfo
    },
    {
      label: '作业文件',
      key: '2',
      component: JobFileSystem
    },
    {
      label: '子作业列表',
      key: '3',
      component: SubJobList
    }
  ]

  get tabs() {
    if (this.jobInfo.jobType === 'JOB' || this.props.isSubJob) {
      return this._tabs.slice(0, 2)
    } else {
      return this._tabs
    }
  }

  onTabChange = activeTabkey => {
    if (activeTabkey === '1') this.updateJobDetailsInfo()
    this.activeTabkey = activeTabkey
  }

  get id() {
    const { isSubJob, jobId } = this.props
    if (isSubJob) {
      const res = jobId.match(/(\d+)\[(\d+)\]/)

      return {
        jobId: res[1],
        jobArrayIndex: res[2]
      }
    } else {
      return { jobId, jobArrayIndex: null }
    }
  }

  async componentDidMount() {
    await this.updateJobDetailsInfo()

    if (this.jobInfo.isPulling) {
      if (!intervalId) {
        intervalId = setInterval(() => {
          this.updateJobDetailsInfo()
        }, 2000)
      }
    } else {
      clearInterval(intervalId)
      intervalId = null
    }

    eventEmitter.on(`JOB_BURST_${this.props.jobId}`, () => {
      this.updateJobDetailsInfo()
    })
  }

  async componentDidUpdate(prevProps) {
    if (prevProps.jobId !== this.props.jobId) {
      await this.updateJobDetailsInfo()
    }

    if (prevProps.activeTabKey !== this.props.activeTabKey) {
      this.activeTabkey = this.props.activeTabKey
    }

    if (this.jobInfo.isPulling) {
      if (!intervalId) {
        intervalId = setInterval(() => {
          this.updateJobDetailsInfo()
        }, 2000)
      }
    } else {
      clearInterval(intervalId)
      intervalId = null
    }
  }

  componentWillUnmount() {
    eventEmitter.off(`JOB_BURST_${this.props.jobId}`)
    clearInterval(intervalId)
    intervalId = null
  }

  updateJobDetailsInfo = async () => {
    if (this.props.isSubJob) {
      if (this.props.isHistory) {
        await historyJobList.fetchListDetail({
          job_id: this.id.jobId,
          job_array_index: this.id.jobArrayIndex
        })
        this.jobInfo = historyJobList.detail
      } else {
        await jobList.fetchListDetail({
          job_id: this.id.jobId,
          job_array_index: this.id.jobArrayIndex
        })
        this.jobInfo = jobList.detail
      }
    } else {
      if (this.props.isHistory) {
        await historyJobList.fetchListDetail({
          job_id: this.id.jobId,
          job_array_index: 0
        })
        this.jobInfo = historyJobList.detail
      } else {
        await jobList.fetchListDetail({
          job_id: this.id.jobId,
          job_array_index: 0
        })
        this.jobInfo = jobList.detail
      }
    }
  }

  renderHeader = jobInfo => {
    return (
      <>
        <Title>作业详情</Title>
        <JobNameWrapper>
          {jobInfo.destClusterName && (
            <Icon
              style={{ color: '#1A6EBA', padding: '0 10px' }}
              type='cloud1'
            />
          )}
          {jobInfo.jobName}
        </JobNameWrapper>
      </>
    )
  }

  render() {
    const { jobInfo, activeTabkey } = this
    if (!jobInfo) return null

    return (
      <>
        <PageHeader
          label={
            this.props.goBack ? (
              <div
                onClick={this.props.goBack}
                title='返回作业中心'
                style={{
                  fontSize: 20
                }}>
                {this.renderHeader(jobInfo)}
              </div>
            ) : (
              this.renderHeader(jobInfo)
            )
          }
        />
        <JobSummary
          isHistory={this.props.isHistory}
          jobInfo={jobInfo}
          reSubmit={this.props.reSubmit}
          reSubmitCallback={this.props.reSubmitCallback}
          updateJobDetailsInfo={this.updateJobDetailsInfo}
        />
        <TabAreaWrapper>
          <Tabs
            activeKey={this.activeTabkey}
            size='large'
            onChange={this.onTabChange}>
            {this.tabs.map(tab => (
              <TabPane tab={tab.label} key={tab.key}>
                <div
                  className='tabArea'
                  style={{ overflowY: tab.key === '3' ? 'hidden' : 'auto' }}>
                  {activeTabkey === tab.key ? (
                    <tab.component
                      jobInfo={jobInfo}
                      jobDir={jobInfo.jobDir}
                      jobId={jobInfo.jobId}
                      isHistory={this.props.isHistory}
                      workspaceId={jobInfo.workspaceId}
                    />
                  ) : null}
                </div>
              </TabPane>
            ))}
          </Tabs>
        </TabAreaWrapper>
      </>
    )
  }
}
