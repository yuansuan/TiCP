// switch muli line chart by type

import React from 'react'
import { Chart, Geom, Tooltip } from 'bizcharts'
import { Button } from '@/components'
import styled from 'styled-components'
import { observable, action, computed } from 'mobx'
import { observer } from 'mobx-react'
import moment from 'moment'
import DashboardData from '@/domain/Dashboard'
import { history } from '@/utils'
import { fullScreenInfo } from './fullScreenInfos'

const DATE_RANGE = {
  '1h': () => [moment().subtract(1, 'h'), moment()].map(m => m.valueOf()),
  '24h': () => [moment().subtract(1, 'days'), moment()].map(m => m.valueOf()),
}

const Wrapper = styled.div<{ isFullScreen: boolean }>`
  .job {
    display: flex;
    justify-content: space-between;

    .btn {
      margin: 0 10px;
      color: #aaa;
    }
    .ant-btn-primary {
      background-color: #1a6eba;
      border-color: #1a6eba;
    }
  }

  .chartTable {
    display: flex;

    .state {
      width: 25%;
      padding: 30px 20px 15px 0;
      text-align: center;
      .name {
        font-size: ${props =>
          props.isFullScreen ? `${fullScreenInfo.baseSize * 0.8}px` : '16px'};
        color: #333333;
        font-weight: 650;
        text-align: left;
      }
      .number {
        font-size: ${props =>
          props.isFullScreen ? `${fullScreenInfo.baseSize * 1.2}px` : '22px'};
        color: #333333;
        padding: 10px 0;
        font-weight: 700;
        text-align: left;
        cursor: pointer;
      }
      .desc {
        font-weight: 650;
        font-size: ${props =>
          props.isFullScreen ? `${fullScreenInfo.baseSize * 0.6}px` : '12px'};
        color: #bcbcbc;
        text-align: center;
      }
    }
    .g2-tooltip {
      pointer-events: none !important;
    }
  }
`

interface IProps {
  dateRange: '1h' | '24h'
  onDateRangeChange: (dateRange: '1h' | '24h') => void
}

@observer
export default class JobStatusInfo extends React.Component<IProps> {
  @observable
  jobData = {
    jobResLatest: [],
    jobResRange: [],
  }

  @computed
  get jobStatusNum() {
    const temp = {
      runningJobs: 0,
      pendingJobs: 0,
      suspendedJobs: 0,
      completedJobs: 0,
      failedJobs: 0,
    }

    this.jobData.jobResLatest?.map(row => {
      row.status === 'Running' && (temp.runningJobs = row.job_count)
      row.status === 'Pending' && (temp.pendingJobs = row.job_count)
      row.status === 'Suspended' && (temp.suspendedJobs = row.job_count)
      row.status === 'Completed' && (temp.completedJobs = row.job_count)
      row.status === 'Failed' && (temp.failedJobs = row.job_count)
    })

    return temp
  }

  getJobInfo = async dateRange => {
    const res = await DashboardData.getDashboardInfo(
      'JOB_INFO',
      DATE_RANGE[dateRange]()
    )
    this.jobData = res.data
  }

  async componentDidMount() {
    this.getJobInfo(this.props.dateRange)
  }

  @action
  private changeDateRange = dateRange => {
    this.props.onDateRangeChange(dateRange)
    this.getJobInfo(dateRange)
  }

  private statusChart = status => {
    const dates = DATE_RANGE[this.props.dateRange]()

    const data =
      this.jobData?.jobResRange?.filter(item => item.status === status) || []

    const jobData =
      data.length === 0
        ? [
            { timestamp: dates[0], status: status, job_count: 0 },
            { timestamp: dates[1], status: status, job_count: 0 },
          ]
        : data
    return jobData
  }

  scale = {
    timestamp: {
      type: 'time',
      mask: 'HH:mm',
    },
    job_count: {
      alias: '作业数',
    },
  }

  private renderRangeChart = status => {
    return (
      <Chart
        height={fullScreenInfo.onlyLineChartHeight}
        padding={[30, 'auto', 5, 'auto'] as any}
        data={this.statusChart(status)}
        scale={this.scale}>
        <Tooltip />
        <Geom
          type='line'
          position='timestamp*job_count'
          size={1}
          shape='smooth'
        />
      </Chart>
    )
  }

  render() {
    const { failedJobs, pendingJobs, suspendedJobs, runningJobs, completedJobs } =
      this.jobStatusNum

    const jobStatusFields = [
      {
        status: 'Failed',
        label: '失败',
        number: failedJobs,
        color: '#E02020',
        isSusOrRun: false, //是否是等待状态或者是运行状态
      },
      {
        status: 'Pending',
        label: '等待',
        number: pendingJobs,
        color: '#9013FE',
        isSusOrRun: true,
      },
      // {
      //   status: 'Suspended',
      //   label: '暂停',
      //   number: suspendedJobs,
      //   color: '#F7B500',
      //   isSusOrRun: false,
      // },
      {
        status: 'Running',
        label: '运行',
        number: runningJobs,
        color: '#1890FF',
        isSusOrRun: true,
      },
      {
        status: 'Completed',
        label: '完成',
        number: completedJobs,
        color: '#52C41A',
        isSusOrRun: false,
      },
    ]

    return (
      <Wrapper isFullScreen={fullScreenInfo.isFullScreen}>
        <div className='job'>
          <div className='title'>作业状态 
            <div className='detail' style={{paddingLeft: 20}}>过去24小时作业状态</div>
          </div>
          <div>
            {' '}
            <Button
              style={{display: 'none'}}
              size='small'
              type={this.props.dateRange === '1h' ? 'primary' : 'default'}
              onClick={() => this.changeDateRange('1h')}>
              过去1小时
            </Button>
            <Button
              style={{display: 'none'}}
              size='small'
              className='btn'
              type={this.props.dateRange === '1h' ? 'default' : 'primary'}
              onClick={() => this.changeDateRange('24h')}>
              过去24小时
            </Button>
          </div>
        </div>
        <div className='chartTable'>
          {jobStatusFields.map((job, index) => {
            return (
              <div className='state' key={index}>
                <div className='name'>{job.label}</div>
                <div
                  className='number'
                  style={{ color: job.color }}
                  onClick={() => history.push('/job')}>
                  {job.number}
                </div>
                <span className='desc'>
                  {job.isSusOrRun
                    ? `当前${job.label}的作业`
                    : `过去${this.props.dateRange === '1h' ? 1 : 24}小时${
                        job.label
                      }的作业`}
                </span>
                <div>{this.renderRangeChart(job.status)}</div>
              </div>
            )
          })}
        </div>
      </Wrapper>
    )
  }
}
