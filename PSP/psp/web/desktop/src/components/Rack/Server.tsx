import * as React from 'react'
import styled from 'styled-components'
import { rowHeight, rackWidth } from './const'
import { Tooltip } from 'antd'
import Node, { INode } from './Node'
import { Icon } from '@/components'

interface IServerWrapper {
  height: number
  rackHeight: number
  start: number
}

const ServerWrapper = styled.div<IServerWrapper>`
  position: absolute;
  width: ${rackWidth - 20 - 2}px;
  background: #408080;
  height: ${props => rowHeight * props.height - 1}px;
  left: 21px;
  top: ${props =>
    (props.rackHeight - props.start - props.height + 1) * rowHeight + 1}px;
  text-align: center;
  display: flex;

  .name {
    position: absolute;
    right: -1px;
    top: ${props => (rowHeight * (props.height - 1)) / 2}px;

    .mark {
      color: #eee;
      padding: 0 0 0 6px;
      width: 10px;
    }
  }

  .nodes {
    display: flex;
    z-index: 1;
    align-items: flex-end;
    flex-wrap: wrap-reverse;
    align-self: flex-end;
  }
`

export interface IServer {
  height: number
  start: number
  name: string
  nodes: INode[]
}

export default function Server({ rackHeight, height, start, name, nodes }) {
  return (
    <ServerWrapper rackHeight={rackHeight} height={height} start={start}>
      <div className='name'>
        {name}
        <Tooltip title={`${name}, 高度 ${height}U`}>
          <span className='mark'>
            <Icon type='shebei' />
          </span>
        </Tooltip>
      </div>
      <div className='nodes'>
        {nodes.map(({ name, value }, index) => (
          <Node key={index} name={name} value={value} />
        ))}
      </div>
    </ServerWrapper>
  )
}
