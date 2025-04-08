/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import ReactDOM, { render } from 'react-dom'
import { isReactComponent } from '@/components/utils'

const StyledDiv = styled.div`
  position: fixed;
  top: 0;
  display: flex;
  height: 100vh;
  width: 100vw;
  justify-content: center;
  touch-action: none;
  z-index: 999;

  > .mask {
    position: relative;
    right: 0;
    top: 0;
    background-color: rgba(0, 0, 0, 0.5);
    width: 100vw;
    height: 100vh;
    margin: auto;
    display: flex;
    justify-content: center;
    align-items: center;
  }
`

export const Mask = function ShareMask({
  close,
  content,
  closable = true
}: Props) {
  let finalContent

  // React.ComponentType
  if (isReactComponent(content)) {
    const Content = content as any
    finalContent = <Content onClose={close} />
  } else if (typeof content === 'string') {
    finalContent = <span style={{ wordBreak: 'break-all' }}>{content}</span>
  } else {
    // null or React.ReactNode
    finalContent = content
  }

  return (
    <StyledDiv onClick={closable ? close : undefined}>
      <div className='mask'>{finalContent}</div>
    </StyledDiv>
  )
}

type Props = {
  children?: any
  close: () => void
  closable?: boolean
  content:
    | React.ReactNode
    | React.ComponentType<{
        onClose?: () => void
      }>
}

export function showMask({ content, ...props }) {
  const div = document.createElement('div')
  document.body.appendChild(div)

  function close() {
    let unmountResult = ReactDOM.unmountComponentAtNode(div)

    if (unmountResult && div.parentNode) {
      div.parentNode.removeChild(div)
    }
  }

  render(<Mask close={close} {...props} content={content} />, div)
}
