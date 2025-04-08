import React, { useRef, useEffect, useState } from 'react'
import styled from 'styled-components'
import { Tooltip } from 'antd'
import { isContentOverflow } from '@/utils'

interface InfoCellProps {
  infoKey: string
  infoVal: string | React.ReactNode
  width?: string
  infoValTip?: string
  infoKeyTip?: string
}

const Wrapper = styled.div`
  display: flex;
  float: left;
  margin-bottom: 10px;
  font-size: 14px;
  line-height: 18px;
  white-space: nowrap;

  .value {
    overflow: hidden;
    text-overflow: ellipsis;
  }
`
export default function(props: InfoCellProps) {
  const { infoKey, infoVal, width, infoValTip, infoKeyTip } = props
  const [isOverflow, setIsOverflow] = useState(false)
  const divEl = useRef(null)

  useEffect(() => {
    setIsOverflow(isContentOverflow(divEl.current, infoVal))
  })

  return (
    <Wrapper style={{ width: width }}>
      <Tooltip title={infoKeyTip}>
        <div style={{ color: 'rgba(0,0,0,0.4)' }}>{infoKey}:&emsp;</div>
      </Tooltip>
      <Tooltip title={isOverflow ? infoValTip : undefined}>
        <div ref={divEl} className='value' style={{ color: 'rgba(0,0,0,0.8)' }}>
          {infoVal}
        </div>
      </Tooltip>
    </Wrapper>
  )
}
