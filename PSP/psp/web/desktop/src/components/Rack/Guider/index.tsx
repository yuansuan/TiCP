import * as React from 'react'
import styled from 'styled-components'
import { colors, rowHeight } from '../const'

interface IGuiderWraper {
  height: number
  horizontal?: boolean
}

const GuiderWraper = styled.div<IGuiderWraper>`
  transform: ${props => (props.horizontal ? `rotate(270deg)` : `rotate(0deg)`)};
  transform-origin: top;
  margin: 10px;
  display: flex;
  .colorBar {
  }
  .textArea {
    height: ${props => props.height}px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    .text {
      writing-mode: vertical-rl;
    }
  }
`

interface IBlock {
  color: string
  height: number
}

const Block = styled.div<IBlock>`
  background: ${props => props.color};
  height: ${props => props.height}px;
  width: 10px;
`

export default function Guider({
  title,
  height,
  horizontal = false,
  showText = true,
  style = {},
}) {
  return (
    <GuiderWraper
      height={height * rowHeight}
      horizontal={horizontal}
      style={style}>
      <div className='colorBar'>
        {colors.map(color => (
          <Block
            key={color}
            color={color}
            height={(height * rowHeight) / 100}
          />
        ))}
      </div>
      {showText && (
        <div className='textArea'>
          <div className='text'>0%</div>
          <div className='text'>{title}</div>
          <div className='text'>100%</div>
        </div>
      )}
    </GuiderWraper>
  )
}
