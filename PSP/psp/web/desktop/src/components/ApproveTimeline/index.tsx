import React from 'react'
import { Timeline } from 'antd'

export function ApproveTimeline({ approve }) {
  return approve.id ? (
    <Timeline
      pending={approve['status'] === 1 ? '等待审批' : false}
      style={{ maxWidth: 360 }}>
      <Timeline.Item>
        <p>{approve['create_time_str']}</p>
        <p>{approve['application_name']}发起申请
           {approve['content']}
        </p>
      </Timeline.Item>
      {approve['status'] === 2 && (
        <Timeline.Item color={'green'}>
          <p>{approve['approve_time_str']}</p>
          <p>
            {approve['approve_user_name']}同意申请
            {approve['suggest'] ? `, ${approve['suggest']}` : ''}
          </p>
        </Timeline.Item>
      )}
      {approve['status'] === 3 && (
        <Timeline.Item color={'red'}>
          <p>{approve['approve_time_str']}</p>
          <p>
            {approve['approve_user_name']}拒绝申请, 拒绝原因:{' '}
            {approve['suggest'] || '无'}
          </p>
        </Timeline.Item>
      )}
      {approve['status'] === 4 && (
        <Timeline.Item color={'gray'}>
          <p>{approve['approve_time_str']}</p>
          <p>
            {approve['application_name']}撤销了申请
          </p>
        </Timeline.Item>
      )}
      {approve['status'] === 5 && (
        <Timeline.Item color={'gray'}>
          <p>{approve['approve_time_str']}</p>
          <p>
            审批操作执行失败
          </p>
        </Timeline.Item>
      )}
    </Timeline>
  ) : (
    <>暂无数据</>
  )
}
