import React from 'react'
import styled from 'styled-components'
import { Spin } from 'antd'

const Wrapper = styled.div`
  width: 100%;
  position: absolute;
  z-index: -1;
  display: flex;
  .loading {
    position: absolute;
    left: 50%;
    top: 50%;
    z-index: 1;
    transform: translate(-50%, -50%);
  }
`

export default function Loading({ loading, message, style, ...props }) {
  const wrapperStyle = style || {}
  if (loading) wrapperStyle['zIndex'] = 2
  return (
    <Wrapper {...props} style={wrapperStyle}>
      <Spin tip={message} className='loading' spinning={loading} />
    </Wrapper>
  )
}
