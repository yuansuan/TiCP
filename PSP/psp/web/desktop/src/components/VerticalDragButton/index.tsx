import React, { useRef, useState, useCallback } from 'react'
import styled from 'styled-components'
import { Tooltip } from 'antd'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'
interface StyledVerticalDragButtonProps {
  isDragging: boolean
  dragY: number
}
const fileTransport = require('@/assets/images/fileTransport.png')

const StyledVerticalDragButton = styled.div<StyledVerticalDragButtonProps>`
  position: fixed;
  right: 0;
  top: ${props => `${props.dragY}px`};
  bottom: 40px;
  width: 20px;
  height: 50px;
  background: ${props => (props.isDragging ? '#0056b3' : '#4193F7')};
  transform: perspective(0.5em) rotateY(-3deg);
  z-index: 9999999999999999999;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: ${props => (props.isDragging ? 'grabbing' : 'grab')};
  user-select: none;
  > img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }
`

const VerticalDragButton = () => {
  const buttonRef = useRef(null)
  const [isDragging, setIsDragging] = useState(false)
  const [buttonTop, setButtonTop] = useState(150)

  function stopDefaultEvent(event) {
    event.preventDefault()
    event.stopPropagation()
  }

  const handleButtonClick = event => {
    stopDefaultEvent(event)
    EE.emit(EE_CUSTOM_EVENT.TOGGLE_UPLOAD_DROPDOWN, { visible: true })
  }

  const handleMouseMove = useCallback(
    event => {
      setIsDragging(true)
      stopDefaultEvent(event)
      const offsetY = event.clientY
      const maxY = document.body.clientHeight - 60
      const newY = Math.min(maxY, Math.max(0, offsetY))
      setButtonTop(newY)
    },
    [isDragging]
  )
  const handleMouseDown = event => {
    stopDefaultEvent(event)
    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
  }

  const handleMouseUp = event => {
    setIsDragging(false)
    stopDefaultEvent(event)
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
  }

  return (
    <Tooltip placement='left' title='文件传输窗口'>
      <StyledVerticalDragButton
        ref={buttonRef}
        isDragging={isDragging}
        dragY={buttonTop}
        onMouseDown={handleMouseDown}
        onMouseUp={handleMouseUp}
        onClick={handleButtonClick}>
        <img src={fileTransport} alt='文件传输' />
      </StyledVerticalDragButton>
    </Tooltip>
  )
}

export default VerticalDragButton
