import * as React from 'react'
import { Tabs } from 'antd'

// import { JobDetail } from '@/components'
import { Wrapper } from './style'

interface IProps {
  debugInfo: string
  jobId: string
  reSubmitCallback: () => void
}

const { TabPane } = Tabs

export default class TestResult extends React.Component<IProps> {
  public render() {
    const { debugInfo, jobId, reSubmitCallback } = this.props

    return (
      <Wrapper>
        <Tabs defaultActiveKey='debugInfo' animated={false}>
          <TabPane key='debugInfo' tab='调试信息'>
            <div className='debugInfo'>{debugInfo}</div>
          </TabPane>
          <TabPane key='jobDetail' tab='任务详情'>
            {/* <div className='jobDetail'>
              <JobDetail jobId={jobId} reSubmitCallback={reSubmitCallback} />
            </div> */}
          </TabPane>
        </Tabs>
      </Wrapper>
    )
  }
}
