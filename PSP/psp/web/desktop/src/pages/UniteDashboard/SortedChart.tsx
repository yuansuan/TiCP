import * as React from 'react'
import { Chart, Geom, Axis, Tooltip, Coord, Legend, Label } from 'bizcharts'
import DataSet from '@antv/data-set'

interface Item {
  name: string
  [propName: string]: any
}

interface IProps {
  data: Item[]
  fields: string[]
  padding?: any
  height?: number
  unit?: string
}

export default class SortedChart extends React.Component<IProps> {
  render() {
    const { data, fields, padding, height, unit } = this.props

    const ds = new DataSet()
    const dv = ds.createView().source(data)
    dv.source(data)
    dv.transform({
      type: 'fold',
      fields: fields,
      key: 'key',
      value: 'value',
    })

    let totals = {}

    dv.rows.forEach(row => {
      const { key, value } = row
      totals[key] = totals[key] || 0
      totals[key] += value || 0 
    })

    dv.rows.map(d => {
      d.total = totals[d.key]
    })

    const scale = {
      value: {
        formatter: val => `${val}${unit}`,
        tickCount: 4,
      },
      total: {
        formatter: val => `${val}${unit}`,
      },
    }

    return (
      <Chart
        height={height || 150}
        width={390}
        data={fields.length === 0 ? [] : dv}
        scale={scale}
        padding={padding || 'auto'}>
        {data.length === 0 ? (
          <div style={{ textAlign: 'center' }}>暂无数据</div>
        ) : (
          <>
            <Coord transpose />
            <Legend position='top-bottom' offsetX={-130} />
            <Axis name='value' />
            <Axis name='key' label={{ offset: 8 }} />
            <Tooltip />
            <Geom
              type='intervalStack'
              position='key*value'
              color={['name', ['#42a3fc', '#54c877']]}>
              <Label
                content={'total'}
                offset={2}
                formatter={(text, item, index) => {
                  // 仅显示 最上面一组的 label 达成总数显示需求
                  if (item._origin.name === '已使用') {
                    return null
                  }
                  // 显示总数
                  return text
                }}
              />
            </Geom>
          </>
        )}
      </Chart>
    )
  }
}
