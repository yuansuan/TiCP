import React from 'react'
import { observer } from 'mobx-react'
import { Steps, Progress } from 'antd'
import styled from 'styled-components'
import moment from 'moment'
import { LoadingOutlined } from '@ant-design/icons'
import { JOB_STEP_STATE } from '@/constant'
const StyledLayout = styled.div`
  padding: 0 20px 20px;
  background: #fff;
  display: flex;
  .ant-steps-item-content {
    .ant-steps-item-title {
      font-size: 14px;
      white-space: normal !important;
    }
  }
`
interface IProps {
  steps: Array<{ name: string; progress: number; time: string }>
}
export const StatusStep = observer(({ steps = [] }: IProps) => {
  const findLocalBursting = steps.find(item => item.name === 'LocalBursting')
  const isBurstFinish = steps.find(item => item.progress === 100)
  const LocalBurstingIndex = steps.findIndex(
    item => item.name === 'LocalBursting'
  )
  const CloudCompleted = steps.findIndex(item => item.name === 'CloudCompleted')
  function getName(name) {
    switch (name) {
      case 'CloudSubmitted':
        if (findLocalBursting) {
          return '爆发完成'
        } else {
          return JOB_STEP_STATE[name]
        }
      case 'CloudRunning':
        if (CloudCompleted) {
          return '已运行'
        } else {
          return JOB_STEP_STATE[name]
        }
      case 'LocalBursting':
        if (steps.length - 1 > LocalBurstingIndex) {
          return '已爆发'
        } else {
          return JOB_STEP_STATE[name]
        }
      default:
        return JOB_STEP_STATE[name]
    }
  }

  function getStatus(name) {
    if (
      name === 'LocalBurstFailed' ||
      name === 'CloudFailed' ||
      name === 'LocalFailed' ||
      name === 'LocalBurstFailed'
    ) {
      return 'error'
    }
    return 'finish'
  }
  return (
    <StyledLayout>
      <Steps
        current={steps.length}
        status={JOB_STEP_STATE[steps[steps.length - 1]?.name]}>
        {steps.map(item => {
          return (
            <Steps.Step
              status={getStatus(item.name)}
              title={getName(item.name)}
              subTitle={
                item.time ? moment(item.time).format('YYYY-MM-DD HH:mm:ss') : ''
              }
              icon={
                item.progress !== -1 && !isBurstFinish ? (
                  <Progress type='circle' percent={item.progress} width={35} />
                ) : null
              }
            />
          )
        })}
      </Steps>
    </StyledLayout>
  )
})
