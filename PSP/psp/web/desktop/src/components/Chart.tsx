/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import {
  Axis,
  Chart as BizCharts,
  Coordinate,
  Line,
  Point,
  Tooltip
} from 'bizcharts'
import {
  IViewProps,
  Marker,
  ColorString,
  FieldString
} from 'bizcharts/src/interface'
import { ColorAttrCallback } from '@antv/g2/lib/interface'

const StyledLayout = styled.div``

type Props = IViewProps & {
  config?: {
    showYLabel?: boolean
    markerType?: Marker
    offsetY?: number
  }
  children?: React.ReactElement
  position: {
    x: string
    y: string
    color?:
      | ColorString
      | FieldString
      | [FieldString, ColorString | ColorString[] | ColorAttrCallback]
  }
}

export function Chart({ position: { x, y, color }, ...props }: Props) {
  return (
    <StyledLayout>
      <BizCharts {...props}>
        <Point position={`${x}*${y}`} color={color} shape='circle' />
        <Axis
          name='value'
          grid={{
            line: {
              style: {
                stroke: '#ccc',
                lineDash: [1, 1]
              }
            }
          }}
          label={{ offset: 20 }}
        />
        <Line position={`${x}*${y}`} color={color} />
        <Tooltip shared showCrosshairs />
        <Coordinate />
        {props.children}
      </BizCharts>
    </StyledLayout>
  )
}
