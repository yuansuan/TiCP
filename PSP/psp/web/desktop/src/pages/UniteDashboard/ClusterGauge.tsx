import React from 'react'
import { Chart, Axis, Coord, Geom, Guide, Shape } from 'bizcharts'

const { Html, Arc } = Guide
const okColor = '#1A6EBA'
const notOkColor = '#E7870E'

Shape.registerShape('point', 'pointer', {
  drawShape(cfg, group) {
    let point = cfg.points[0] // 获取第一个标记点
    // @ts-ignore
    point = this.parsePoint(point)
    // @ts-ignore
    const center = this.parsePoint({
      // 获取极坐标系下画布中心点
      x: 0,
      y: 0,
    })
    // 绘制指针
    group.addShape('line', {
      attrs: {
        x1: center.x,
        y1: center.y + 10,
        x2: point.x,
        y2: point.y,
        stroke: cfg.color,
        lineWidth: 3,
        lineCap: 'round',
      },
    })

    return group.addShape('circle', {
      attrs: {
        x: center.x,
        y: center.y + 10,
        r: 5,
        stroke: cfg.color,
        lineWidth: 3,
        fill: '#fff',
      },
    })
  },
})

interface IProps {
  data: any[]
  color?: string
  height?: number
  min?: number
  max?: number
}

export default class GaugeChart extends React.Component<IProps> {
  render() {
    const nodes = this.props.data[0]?.value || 0
    const notOkNodes = this.props.data[1]?.value || 0
    const data = [{ value: nodes }]
    const { height, min, max } = this.props

    const cols = {
      value: {
        min: min || 0,
        max: max || nodes + notOkNodes,
        nice: false,
      },
    }

    return (
      <Chart
        height={height || 160}
        width={100}
        data={data}
        scale={cols}
        padding={[-30, 0, 10, 0]}>
        <Coord
          type='polar'
          startAngle={(-10 / 8) * Math.PI}
          endAngle={(2 / 8) * Math.PI}
          radius={0.78}
        />
        <Axis name='value' visible={false} />
        <Axis name='1' visible={false} />
        <Guide>
          <Arc
            start={[0, 0.965]}
            end={[nodes + notOkNodes, 0.965]}
            style={{
              stroke: notOkColor,
              lineWidth: 18,
            }}
          />
          <Arc
            top
            start={[0, 0.965]}
            end={[nodes, 0.965]}
            style={{
              stroke: this.props.color || okColor,
              lineWidth: 18,
            }}
          />
          <Html
            position={['50%', '105%']}
            html={`<div style="width: 400px;text-align: center; color: #797979;font-size: 14px;">可用节点/总节点</br>${nodes}&nbsp;/&nbsp;${
              nodes + notOkNodes
            }</div>`}
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
