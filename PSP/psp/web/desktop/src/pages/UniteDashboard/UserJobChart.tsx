import * as React from 'react'
import { Chart, Geom, Axis, Tooltip, Coord, Legend, Label } from 'bizcharts'
import DataSet from '@antv/data-set'

interface Item {
  key: string
  value: number
}

interface IProps {
  data: Item[]
  padding?: any
  color?: string
}

export default class UserJobChart extends React.Component<IProps> {
  render() {
    const { data, padding, color } = this.props

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
          return index < 5
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

    const scale = {
      value: {
        alias: '作业数',
      },
    }
    return (
      <Chart height={180} data={dv} padding={padding || 'auto'} scale={scale}>
        {data.length === 0 ? (
          <div style={{ textAlign: 'center' }}>暂无数据</div>
        ) : (
          <>
            <Coord transpose />
            <Axis
              name='key'
              label={{
                offset: 12,
              }}
            />
            <Axis name='value' line={{ fill: '#ffffff' }} />
            <Tooltip />
            <Geom type='interval' position='key*value' color={color} />
          </>
        )}
      </Chart>
    )
  }
}
