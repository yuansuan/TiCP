import * as React from 'react'
import { InfoCircleOutlined }  from '@ant-design/icons'

interface ShowIconProps {
  isShowSummary: boolean
  onClick: () => void
}

export function IconControlSummary(props: ShowIconProps) {
  const { isShowSummary, onClick } = props

  return (
    <InfoCircleOutlined
      rev="icon-control-summary"
      title={isShowSummary ? '隐藏信息' : '显示信息'}
      style={{
        color: isShowSummary ? 'blue' : '#aaa',
        margin: '0 15px',
        fontSize: 20,
      }}
      onClick={() => onClick()}
    />
  )
}
