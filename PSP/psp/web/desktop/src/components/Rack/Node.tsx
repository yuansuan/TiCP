import * as React from 'react'
import { Tooltip } from 'antd'
import styled from 'styled-components'
import { colors, rowHeight } from './const'

interface INodeWrapper {
  color: string
}

const NodeWrapper = styled.div<INodeWrapper>`
  width: ${rowHeight - 4}px;
  height: ${rowHeight - 4}px;
  background: ${props => props.color};
  margin: 1px;
`

export interface INode {
  name: string
  value: number
}

export default function Node({ name, value }) {
  return (
    <Tooltip placement='right' title={`${name}, CPU使用率: ${value}%`}>
      <NodeWrapper color={colors[value]} />
    </Tooltip>
  )
}
