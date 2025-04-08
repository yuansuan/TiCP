/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect, useState } from 'react'
import styled from 'styled-components'
import { Status } from '@/components'
import { Job } from '@/domain/JobList/Job'
import { Dropdown, Menu, Progress } from 'antd'
import { Icon } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { formatByte } from '@/utils/Validator'
import { JobState, JOB_STATE_ENUM } from './JobState'
import { ALL_JOB_STATES } from '@/constant'
import { jobServer } from '@/server'
import { runInAction } from 'mobx'

enum BackStateMap {
  NONE = 0,
  BACK_PROCESS = 1,
  BACK_FINISH = 2
}

const Wrapper = styled.div``

const StatusWrapper = styled.div`
  padding-left: 0px !important;
`

const StyledMenu = styled(Menu)`
  .content {
    margin: 10px 15px 15px 14px;
    line-height: 22px;

    .description {
      display: inline-block;
      padding-left: 3px;
      font-size: 12px;
    }

    .progress {
      display: inline-block;
      width: 120px;
      padding-left: 12px;
    }

    .info {
      display: inline-block;
      margin-left: 12px;
      font-size: 12px;

      .anticon {
        height: 12px;
        width: 12px;
      }
    }

    .process-line {
      margin-top: 10px;
    }
  }
`

interface IProps {
  job: Job
  className?: string
  showDropDown?: boolean
}

const statusMapping = {
  等待中: 'primary',
  运行中: 'primary',
  已完成: 'success',
  已失败: 'error',
  已提交: 'primary',
  已终止: 'warn',
  已暂停: 'primary', //'warn'
  爆发失败: 'error',
  爆发中: 'primary'
}

const statusArr = [
  {
    descriptions: '调度',
    percent: 0,
    zeroInfo: '0%',
    processInfo: '100mb/s',
    finishInfo: '100%'
  },
  {
    descriptions: '计算',
    percent: 0,
    zeroInfo: '0%',
    processInfo: <Icon type='loading' />,
    finishInfo: '完成'
  },
  {
    descriptions: '回传',
    percent: 0,
    zeroInfo: '0%',
    processInfo: '100mb/s',
    finishInfo: '100%'
  }
]

const cheatPercent = 20

const round = num => Math.round(num * 100) / 100

const SmallProgress = ({ percent }: { percent: number }) => (
  <Progress
    className='process-line'
    strokeWidth={6}
    showInfo={true}
    format={percent < 100 ? percent => `${round(percent)}%` : undefined}
    percent={percent}
    size='small'
  />
)

type LineProps = {
  description: string
  percent: number
  info: React.ReactElement | any
}

const LineBlock = (props: LineProps) => (
  <div>
    <div className='description'>{props.description}</div>
    <div className='progress'>
      <SmallProgress percent={props.percent} />
    </div>
    <div className='info'>{props.info}</div>
  </div>
)

const JobStatus = observer(function JobStatus(props: IProps) {
  const { job, className, showDropDown = false } = props

  const [visible, setVisible] = useState(false)

  const store = useLocalStore(() => ({
    jobState: new JobState({
      state: JOB_STATE_ENUM.UNKNOWN,
      speed: 0,
      progress: 0,
      total_size: 0
    })
  }))

  const getJobState = async () => {
    if (visible) {
      const { data } = await jobServer.getStatus(job.id)
      runInAction(() => {
        store.jobState.update(data)
      })
    }
  }

  function displayDropDown() {
    // 需要显示下弹窗 && (是需要展示的state || 是有回传状态的state)
    return showDropDown && [1, 2].includes(store.jobState.state)
  }

  function displayDropDownItem(statusArr) {
    // 把显示下拉框抽象成一个函数
    const len = statusArr.length
    if ([1, 3].includes(store.jobState.state)) {
      // 如果在上传，计算，回传过程中
      return (
        <Menu.Item key={1}>
          <LineBlock
            description={statusArr[store.jobState.state - 1].descriptions}
            percent={
              // 回传中的百分比永远加上20，即进度百分比显示在80%中，造成一种文件回传和厉害的假象
              store.jobState.state === JOB_STATE_ENUM.BACK
                ? cheatPercent +
                  (store.jobState.progress / 100) * (100 - cheatPercent)
                : store.jobState.progress
            }
            info={`${formatByte(store.jobState.speed)}/s`}
          />
        </Menu.Item>
      )
    } else if ([0, 4].includes(store.jobState.state)) {
      return (
        // 作业状态要么是0 未知 要么是4 已完成
        <Menu.Item key={1}>
          <LineBlock
            description={statusArr[len - 1].descriptions}
            percent={100}
            info={statusArr[len - 1].finishInfo}
          />
        </Menu.Item>
      )
    } else {
      // 如果是计算 不显示计算进度条，因为不准
      return null
    }
  }

  useEffect(() => {
    let interval = null
    if (visible) {
      getJobState()
      interval = setInterval(() => getJobState(), 2000)
    }
    return () => {
      interval && clearInterval(interval)
    }
  }, [visible])

  return (
    <Wrapper>
      {displayDropDown() ? (
        <Dropdown
          visible={visible}
          onVisibleChange={setVisible}
          placement='bottomCenter'
          overlayStyle={{
            width: '275px'
          }}
          overlay={
            <StyledMenu>
              <div className='content'>
                {displayDropDownItem(statusArr)}

                <Menu.Item>
                  <p
                    style={{
                      fontSize: '12px',
                      paddingTop: '10px',
                      marginLeft: '3px'
                    }}>
                    {(store.jobState.state === 0 ||
                      store.jobState.state === 1) &&
                      '上传文件大小'}
                    {(store.jobState.state === 2 ||
                      store.jobState.state === 3 ||
                      store.jobState.state === 4) &&
                      '回传文件大小'}
                    ：{formatByte(store.jobState?.total_size)}
                  </p>
                </Menu.Item>
              </div>
            </StyledMenu>
          }>
          <StatusWrapper className={className}>
            <Status
              text={ALL_JOB_STATES[job.state]}
              type={statusMapping[ALL_JOB_STATES[job.state]]}
            />
          </StatusWrapper>
        </Dropdown>
      ) : (
        <StatusWrapper className={className}>
          <Status
            text={ALL_JOB_STATES[job.state]}
            type={statusMapping[ALL_JOB_STATES[job.state]]}
          />
        </StatusWrapper>
      )}
    </Wrapper>
  )
})
export default JobStatus
