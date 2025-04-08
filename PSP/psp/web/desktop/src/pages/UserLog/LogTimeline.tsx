import React from 'react'
import { Timeline } from 'antd'
import { Icon } from '@/components'
import { replaceCustomTags } from '@/utils/formatter'

export function LogTimeline({ logs, iconMap }) {
  return (
    <Timeline mode='left'>
      {logs.length > 0 ? (
        logs.map(log => {
          return (
            <Timeline.Item
              key={log.id}
              dot={
                <Icon
                  type={iconMap[log.operate_type]}
                  style={{ fontSize: '16px' }}
                />
              }>
              {log.operate_time_str} - {log.ip_address} -{' '}
              {replaceCustomTags(log.operate_content)}
            </Timeline.Item>
          )
        })
      ) : (
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            width: '100wh',
            height: 'calc(100vh - 360px)',
          }}>
          暂无数据
        </div>
      )}
    </Timeline>
  )
}
