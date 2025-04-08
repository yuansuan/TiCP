import * as React from 'react'
import { Chart, Geom, Axis, Tooltip, Coord, Legend, Label } from 'bizcharts'
import { ChartTitle } from './style'
import DataSet from '@antv/data-set'
import { PRECISION } from '@/domain/common'
import { roundNumber } from '@/utils/formatter'
import { date } from './commonShowDates'

interface Item {
  key: string
  value: number
}

interface IProps {
  data: Item[]
  title?: string
  height?: number
  unit?: string
  top?: number
  theme?: string
  hideLegend?: boolean
  barColor?: string
  valueAxis?: string
  padding?: any
  tooltipFomatter?: (key: string, value: number) => any
  showDates?: any
}

export default class SortedChart extends React.Component<IProps> {
  render() {
    const {
      data,
      title,
      height,
      unit,
      top,
      hideLegend,
      theme,
      barColor,
      valueAxis,
      padding,
      tooltipFomatter,
      showDates
    } = this.props

    const ds = new DataSet()
    const dv = ds.createView().source(data)
    dv.source(data)
    dv.transform({
      type: 'sort-by',
      fields: ['value'],
      order: 'DESC'
    })
      .transform({
        type: 'filter',
        callback(row, index) {
          return index < (top || 10)
        }
      })
      .transform({
        type: 'sort-by',
        fields: ['key'],
        order: 'DESC'
      })
      .transform({
        type: 'sort-by',
        fields: ['value'],
        order: 'ASC'
      })

    const tooltipTemp = `<li data-index={index}>
      <span style="background-color:{color};width:8px;height:8px;border-radius:50%;display:inline-block;margin-right:8px;"></span>
      {name}: {value} ${unit || ''}
    </li>`

    const defaultTooltipFormatter = (key, value) => ({
      name: key,
      value: roundNumber(value, PRECISION)
    })

    const scale = {
      value: {
        alias: valueAxis || '完成作业数'
      }
    }
    return (
      <Chart
        theme={theme || 'default'}
        height={data.length !== 0 ? height || 400 : height - 30}
        scale={scale}
        data={dv}
        padding={padding || 'auto'}
        forceFit>
        {title && <ChartTitle>{title}</ChartTitle>}
        {showDates && date(showDates)}
        {data.length === 0 ? (
          <div style={{ textAlign: 'center' }}>暂无数据</div>
        ) : (
          <>
            <Coord transpose />
            {!hideLegend && <Legend />}
            <Axis
              name='key'
              label={{
                offset: 12
              }}
            />
            <Axis name='value' visible={false} />
            <Tooltip
              itemTpl={tooltipTemp}
              crosshairs={{
                type: 'y'
              }}
            />
            <Geom
              type='interval'
              position='key*value'
              color={barColor || 'key'}
              tooltip={[
                'key*value',
                tooltipFomatter || defaultTooltipFormatter
              ]}>
              <Label
                content={[
                  'key*value',
                  (key, value) => {
                    return `${roundNumber(value, PRECISION)}${unit || ''}`
                  }
                ]}
                offset={5}
              />
            </Geom>
          </>
        )}
      </Chart>
    )
  }
}
