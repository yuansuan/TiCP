import React from 'react'
import { Chart, Axis, Coord, Geom, Guide, Shape } from 'bizcharts'
import { PRECISION } from '@/domain/common'
import { roundNumber } from '@/utils/formatter'

const { Html, Arc } = Guide

Shape.registerShape('point', 'pointer', {
  drawShape(cfg, group) {
    let point = cfg.points[0] // 获取第一个标记点
    // @ts-ignore
    point = this.parsePoint(point)
    // @ts-ignore
    const center = this.parsePoint({
      // 获取极坐标系下画布中心点
      x: 0,
      y: 0
    })
    // 绘制指针
    group.addShape('line', {
      attrs: {
        x1: center.x,
        y1: center.y,
        x2: point.x,
        y2: point.y,
        stroke: cfg.color,
        lineWidth: 5,
        lineCap: 'round'
      }
    })

    return group.addShape('circle', {
      attrs: {
        x: center.x,
        y: center.y,
        r: 12,
        stroke: cfg.color,
        lineWidth: 4.5,
        fill: '#fff'
      }
    })
  }
})

interface IData {
  value: number
}

interface IProps {
  data: IData[]
  title?: string
  unit?: string
  height?: number
  min?: number
  max?: number
}

export default class GaugeChart extends React.Component<IProps> {
  componentDidMount() {}

  render() {
    const { data, title, unit, height, min, max } = this.props

    const cols = {
      value: {
        min: min || 0,
        max: max || 100,
        tickInterval: 10,
        nice: false
      }
    }

    return (
      <Chart
        height={height || 400}
        data={data}
        scale={cols}
        padding={[0, 0, 30, 0]}
        forceFit>
        <Coord
          type='polar'
          startAngle={(-9 / 8) * Math.PI}
          endAngle={(1 / 8) * Math.PI}
          radius={0.75}
        />
        <Axis
          name='value'
          zIndex={2}
          line={null}
          label={{
            offset: -16,
            formatter: val => {
              return `${val}${unit || ''}`
            },
            textStyle: {
              fontSize: 12,
              textAlign: 'center',
              textBaseline: 'middle'
            }
          }}
          subTickCount={4}
          subTickLine={{
            length: -8,
            stroke: '#fff',
            strokeOpacity: 1
          }}
          tickLine={{
            length: -18,
            stroke: '#fff',
            strokeOpacity: 1
          }}
        />
        <Axis name='1' visible={false} />
        <Guide>
          <Arc
            start={[0, 0.965]}
            end={[100, 0.965]}
            style={{
              // 底灰色
              stroke: '#CBCBCB',
              lineWidth: 18
            }}
          />
          <Arc
            top
            start={[0, 0.965]}
            end={[data[0].value, 0.965]}
            style={{
              stroke: '#1890FF',
              lineWidth: 18
            }}
          />
          <Html
            position={['50%', '5%']}
            html={`<div style="width: 200px;text-align: center;">
                    <p style="font-size: 14px;margin: 0;font-weight: 500;">
                      ${title}
                    </p>
                  </div>`}
          />
          <Html
            position={['50%', '80%']}
            html={`<div style="width: 200px;text-align: center;font-size: 12px!important;">
                    <p style="font-size: 16px;color: rgba(0,0,0,0.85);margin: 0;margin-top: 20px">
                    ${
                      data[0].value === '-'
                        ? '-'
                        : roundNumber(data[0].value, PRECISION)
                    }%
                    </p>
                  </div>`}
          />
        </Guide>
        <Geom
          type='point'
          position='value*1'
          shape='pointer'
          color='#1890FF'
          active={false}
          style={{ stroke: '#fff', lineWidth: 1 }}
        />
      </Chart>
    )
  }
}
