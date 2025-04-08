import React from 'react'
import { Chart, Geom, Axis, Tooltip, Legend, Coord } from 'bizcharts'
import DataSet from '@antv/data-set'
import { ChartTitle } from './style'
import { PRECISION } from '@/domain/common'
import { roundNumber } from '@/utils/formatter'
import { date } from './commonShowDates'

interface Item {
  name: string
  [propName: string]: any
}

interface IProps {
  data: Item[]
  fields: string[]
  title?: string
  unit?: string
  height?: number
  min?: number
  max?: number
  type?: 'horizontal' | 'vertical'
  minTickInterval?: number
  tooltipFomatter?: (key: string, value: number) => any
  showDates?: any
}

export default class ColumnChart extends React.Component<IProps> {
  render() {
    const {
      data,
      fields,
      title,
      unit,
      height,
      min,
      max,
      type,
      minTickInterval,
      tooltipFomatter,
      showDates,
    } = this.props

    const ds = new DataSet()
    const dv = ds.createView().source(data)
    dv.transform({
      type: 'fold',
      fields: fields,
      // 展开字段集
      key: 'key',
      // key字段
      value: 'value', // value字段
    }).transform({
      type: 'filter',
      callback(row, index) {
        return row.value !== undefined
      }
    })

    console.log(dv)

    const scale = {
      value: {
        alias: unit ? `${unit} ` : ' ',
        min: min,
        max: max,
        minTickInterval: minTickInterval,
      },
    }

    const tooltipTemp = `<li data-index={index}>
      <span style="background-color:{color};width:8px;height:8px;border-radius:50%;display:inline-block;margin-right:8px;"></span>
      {name}: {value} ${unit || ''}
    </li>`
    const defaultTooltipFormatter = (key, value) => ({
      name: key,
      value: roundNumber(value, PRECISION),
    })
    return (
      <div>
        <Chart
          height={height || 400}
          data={fields.length === 0 ? [] : dv}
          scale={scale}
          padding={'auto'}
          forceFit>
          {title && <ChartTitle>{title}</ChartTitle>}
          {showDates && date(showDates)}
          {type === 'horizontal' && <Coord transpose />}
          <Legend />
          <Axis name='key' />
          <Axis
            position={type === 'horizontal' ? 'right' : 'left'}
            name='value'
            label={{
              formatter: val => `${val}${unit}`,
            }}
          />
          <Tooltip
            itemTpl={tooltipTemp}
            crosshairs={{
              type: 'y',
            }}
          />
          <Geom
            type='intervalStack'
            position='key*value'
            color={'name'}
            tooltip={['name*value', tooltipFomatter || defaultTooltipFormatter]}
          />
        </Chart>
      </div>
    )
  }
}
