// switch muli line chart by type

import React from 'react'
import { Chart, Geom, Axis, Tooltip, Legend } from 'bizcharts'
import { ChartTitle, SwitchKeyWrapper } from './style'

interface Item {
  timestamp: number
  value: string | number
  type: string
}

interface ItemsBykey {
  [propName: string]: Item[]
}

interface IProps {
  data: ItemsBykey
  keys: string[]
  title?: string
  unit?: string
  height?: number
  min?: number
  max?: number
  timeFormat?: string
}

export default class LineChartByType extends React.Component<IProps> {
  state = {
    key: this.props.keys[0],
    data: this.props.data[this.props.keys[0]],
  }

  onChange = key => {
    this.setState({
      key,
      data: this.props.data[key],
    })
  }

  render() {
    const { title, keys, unit, height, min, max, timeFormat } = this.props

    const scale = {
      timestamp: {
        type: 'time',
        tickCount: 4,
        mask: timeFormat || 'MM-DD HH:mm',
      },
      value: {
        minTickInterval: 1,
        alias: unit ? `${unit} ` : ' ',
        min: min,
        max: max,
      },
    }

    return (
      <Chart
        height={height || 400}
        data={this.state.data || []}
        scale={scale}
        forceFit>
        {title && <ChartTitle>{title}</ChartTitle>}
        {keys && (
          <SwitchKeyWrapper>
            {keys.map(k => (
              <span
                key={k}
                style={{ color: this.state.key === k ? 'blue' : '#888' }}
                onClick={() => {
                  this.onChange(k)
                }}>
                {k}
              </span>
            ))}
          </SwitchKeyWrapper>
        )}
        <Legend />
        <Axis name='timestamp' />
        <Axis
          name='value'
          label={{
            formatter: val => `${val}${unit || ''}`,
          }}
        />
        <Tooltip
          crosshairs={{
            type: 'y',
          }}
          inPlot={false}
        />
        <Geom type='line' position='timestamp*value' color={'type'} />
      </Chart>
    )
  }
}
