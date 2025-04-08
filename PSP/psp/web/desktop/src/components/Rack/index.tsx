import * as React from 'react'
import styled from 'styled-components'
import { rowHeight, rackWidth } from './const'
import Server, { IServer } from './Server'

interface IRackBodyProps {
  rowHeight?: number
  height?: number // 42 U 机架高度
  rackWidth?: number
}

const RackWrapper = styled.div`
  position: relative;
  margin: 10px;
`

const RackBody = styled.div<IRackBodyProps>`
  display: flex;
  flex-direction: column-reverse;
  width: ${rackWidth}px;
  border: 1px solid #000;
  border-bottom: 0px solid #000;

  .rackRow {
    display: flex;
    box-sizing: border-box;
    height: ${props => props.rowHeight}px;
    border-bottom: 1px solid #000;
    .number {
      width: 20px;
      text-align: center;
      border-right: 1px solid #000;
    }
    .server {
      background: #8fbfbf;
      width: ${rackWidth - 20}px;
    }
  }
`

function RackRow({ rowNumber }) {
  return (
    <div className='rackRow'>
      <div className='number'>{rowNumber}</div>
      <div className='server'></div>
    </div>
  )
}

interface IRack {
  height: number
  servers: IServer[]
}

export default class Rack extends React.Component<IRack> {
  render() {
    const { height, servers } = this.props
    const totalRackRows = new Array(height).fill(0)

    return (
      <RackWrapper>
        <>
          {servers.map((server, index) => (
            <Server
              key={index}
              rackHeight={height}
              start={server.start}
              height={server.height}
              name={server.name}
              nodes={server.nodes}
            />
          ))}
        </>

        <RackBody height={height} rowHeight={rowHeight}>
          {totalRackRows.map((v, index) => (
            <RackRow key={index} rowNumber={index + 1} />
          ))}
        </RackBody>
      </RackWrapper>
    )
  }
}
