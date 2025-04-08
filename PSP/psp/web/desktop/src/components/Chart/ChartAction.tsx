import React from 'react'
import { Button } from 'antd'
import { ChartActionWrapper } from './style'

export default function ChartAction({
  chartData,
  exportImage,
  exportExcel,
  otherBtn,
}) {
  const btnStyle = {
    padding: 0,
    margin: '0 15px',
  }
  return (
    <ChartActionWrapper>
      {otherBtn && otherBtn()}
      <Button
        style={btnStyle}
        disabled={!chartData}
        onClick={exportImage}
        type='link'>
        保存图片
      </Button>
      <Button
        style={btnStyle}
        disabled={!chartData}
        onClick={exportExcel}
        type='link'>
        数据导出
      </Button>
    </ChartActionWrapper>
  )
}
