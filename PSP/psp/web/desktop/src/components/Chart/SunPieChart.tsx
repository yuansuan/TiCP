import React from 'react'
import { Chart, Geom, Tooltip, Coord, Label, View, Legend } from 'bizcharts'
import DataSet from '@antv/data-set'
import { ChartTitle } from './style'
import { PRECISION } from '@/domain/common'
import { roundNumber } from '@/utils/formatter'

interface Item {
  name: string
  type: string
  value: number
}

interface IProps {
  data: Item[]
  title?: string
  height?: number
}

export default class SunPieChart extends React.Component<IProps> {
  render() {
    const { DataView } = DataSet

    const { data, title, height } = this.props

    const dv = new DataView()

    dv.source(data).transform({
      type: 'percent',
      field: 'value',
      dimension: 'type',
      as: 'percent'
    })

    const cols = {
      percent: {
        formatter: val => {
          val = `${roundNumber(val * 100, PRECISION)}%`
          return val
        }
      }
    }
    const dv1 = new DataView()

    dv1.source(data).transform({
      type: 'percent',
      field: 'value',
      dimension: 'name',
      as: 'percent'
    })

    return (
      <Chart
        height={height || 400}
        data={dv}
        scale={cols}
        padding={'auto'}
        forceFit>
        {title && <ChartTitle>{title}</ChartTitle>}
        {data.length === 0 && (
          <div
            style={{ position: 'relative', top: '48%', textAlign: 'center' }}>
            暂无数据
          </div>
        )}
        <Coord type='theta' radius={0.75} />
        <Tooltip
          showTitle={false}
          itemTpl='<li><span style="background-color:{color};" class="g2-tooltip-marker"></span>{name}: {value}</li>'
        />
        <Geom
          type='intervalStack'
          position='percent'
          color='type'
          tooltip={[
            'type*percent',
            (item, percent) => {
              percent = `${roundNumber(percent * 100, PRECISION)}%`
              return {
                name: item,
                value: percent
              }
            }
          ]}
          style={{
            lineWidth: 1,
            stroke: '#fff'
          }}
          select={false}>
          <Label content='type' offset={-10} type='treemap' />
        </Geom>
        <Legend name='type' position='left' offsetX={200} />
        <Legend name='name' />
        <View data={dv1} scale={cols}>
          <Coord type='theta' radius={0.9} innerRadius={0.75 / 0.9} />
          <Geom
            type='intervalStack'
            position='percent'
            color={'name'}
            tooltip={[
              'name*percent',
              (item, percent) => {
                percent = `${roundNumber(percent * 100, PRECISION)}%`
                return {
                  name: item,
                  value: percent
                }
              }
            ]}
            style={{
              lineWidth: 1,
              stroke: '#fff'
            }}
            select={false}>
            <Label content='name' />
          </Geom>
        </View>
      </Chart>
    )
  }
}
