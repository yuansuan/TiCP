import React from 'react'
import styled from 'styled-components'
import { Chart, Geom, Axis, Tooltip, Legend, Guide } from 'bizcharts'
import { ChartTitle } from './style'
import Slider from 'bizcharts-plugin-slider'
import DataSet from '@antv/data-set'
import { PRECISION } from '@/domain/common'
import { roundNumber } from '@/utils/formatter'
import { date } from './commonShowDates'

const Wrapper = styled.div``

const paddingMap = [15, 16, 16, 16, 14, 12, 10]

interface DataItem {
  t: number
  v: number
}

interface NamedItem {
  n: string
  d: DataItem[]
}

interface IProps {
  data: NamedItem[]
  lineColor?: string
  title?: string
  unit?: string
  height?: number
  min?: number
  max?: number
  type?: 'line' | 'area'
  timeFormat?: string
  slider?: boolean
  start?: number
  end?: number
  showMarker?: boolean
  showSummary?: boolean
  theme?: string
  hideLegend?: boolean
  summaryPosition?: any
  summaryWidth?: string
  padding?: number[] | string
  showDates?: any
}

// TODO showMarker 与 slider 功能冲突

export default class LineChart extends React.Component<IProps> {
  chartIns = null

  names = this.props.data.map(d => d?.n)

  minVals = this.names.reduce((pre, curr) => {
    pre[curr] = { t: -1, v: Number.MAX_SAFE_INTEGER, n: curr }
    return pre
  }, {})

  maxVals = this.names.reduce((pre, curr) => {
    pre[curr] = { t: -1, v: Number.MIN_SAFE_INTEGER, n: curr }
    return pre
  }, {})

  sumVals = this.names.reduce((pre, curr) => {
    pre[curr] = { t: -1, v: 0, n: curr }
    return pre
  }, {})

  pointsVals = this.names.reduce((pre, curr) => {
    pre[curr] = { t: -1, v: 0, n: curr }
    return pre
  }, {})

  ds = new DataSet({
    state: {
      start: this.props.start || this.props.data[0]?.d[0]?.t,
      end:
        this.props.end ||
        this.props.data[0]?.d[this.props.data[0]?.d.length - 1]?.t
    }
  })

  originView = this.ds.createView('origin')

  render() {
    const {
      data,
      lineColor,
      title,
      unit,
      height,
      min,
      max,
      type,
      timeFormat,
      slider,
      showMarker,
      theme,
      hideLegend,
      padding,
      showSummary,
      summaryPosition,
      showDates,
      summaryWidth,
    } = this.props

    const chartData = data.reduce((d, item) => {
      const items = item.d.map(t => ({ ...t, n: item.n }))
      return [...d, ...items]
    }, [])

    const dv = this.originView.source(chartData)

    dv.transform({
      type: 'filter',
      callback: obj => {
        const time = obj.t
        const value = obj.v
        const name = obj.n

        // 计算最大值和最小值
        if (time >= this.ds.state.start && time <= this.ds.state.end) {
          if (this.maxVals[name].v < value) {
            this.maxVals[name].v = value
            this.maxVals[name].t = time
          }
          if (this.minVals[name].v > value) {
            this.minVals[name].v = value
            this.minVals[name].t = time
          }

          this.pointsVals[name].v += 1
          this.sumVals[name].v += value
        }
        return time >= this.ds.state.start && time <= this.ds.state.end
      }
    })

    const scale = {
      t: {
        type: 'time',
        tickCount: 4,
        mask: timeFormat || 'MM-DD HH:mm'
      },
      v: {
        minTickInterval: 1,
        alias: unit ? `${unit} ` : ' ',
        min: min,
        max: max
      }
    }

    const tooltipTemp = `<li data-index={index}>
      <span style="background-color:{color};width:8px;height:8px;border-radius:50%;display:inline-block;margin-right:8px;"></span>
      {name}: {value} ${unit || ''}
    </li>`

    return (
      <Wrapper>
        <Chart
          theme={theme || 'default'}
          height={height || 400}
          data={slider ? dv : chartData}
          scale={scale}
          forceFit
          padding={padding || 'auto'}>
          {title && <ChartTitle>{title}</ChartTitle>}
          {showDates && date(showDates)}
          {showSummary && <div style={{ height: 25 }}></div>}
          <Axis name='t' />
          <Axis
            name='v'
            label={{
              formatter: val => `${val}${unit}`
            }}
          />
          <Tooltip
            itemTpl={tooltipTemp}
            crosshairs={{
              type: 'y'
            }}
            inPlot={false}
          />
          {!hideLegend && <Legend marker='hyphen' />}
          <Geom
            type={type || 'line'}
            position='t*v'
            color={lineColor || 'n'}
            tooltip={[
              'n*t*v',
              (n, _, v) => {
                return {
                  name: n,
                  value: `${roundNumber(v, PRECISION)}`
                }
              }
            ]}
          />
          <Guide>
            {showSummary && (
              <Guide.Html
                position={
                  summaryPosition || [
                    `${this.names.length * paddingMap[this.names.length]}%`,
                    '-8%'
                  ]
                }
                html={(xScale, yScale) => {
                  return `<div style="font-size:12px;border-radius: 5px; display: flex;flex-wrap: wrap;justify-content: center; width: ${summaryWidth || '100%'};">
                    ${this.names
                      .map(
                        n =>
                          `<span> <b>&nbsp;${n}</b> 最大值: ${
                            this.maxVals[n]?.v === Number.MIN_SAFE_INTEGER
                              ? '--'
                              : roundNumber(
                                  Number(this.maxVals[n]?.v),
                                  PRECISION
                                )
                          }${unit || ''} 最小值: ${
                            this.minVals[n]?.v === Number.MAX_SAFE_INTEGER
                              ? '--'
                              : roundNumber(
                                  Number(this.minVals[n]?.v),
                                  PRECISION
                                )
                          }${unit || ''} 平均值: ${
                            this.pointsVals[n]?.v === 0
                              ? '--'
                              : roundNumber(
                                  Number(
                                    this.sumVals[n]?.v / this.pointsVals[n]?.v
                                  ),
                                  PRECISION
                                )
                          }${unit || ''}</span>`
                      )
                      .join('')}
                    </div>`
                }}
              />
            )}
            {showMarker && (
              <>
                {this.names.map((n, index) => (
                  <>
                    <Guide.DataMarker
                      key={`${n}max`}
                      position={(x, y) => {
                        return this.maxVals[n]
                      }}
                      lineLength={(index + 1) * 30}
                      content={`${n}最大值`}
                    />
                    <Guide.DataMarker
                      key={`${n}min`}
                      lineLength={(index + 1) * 20}
                      position={() => this.minVals[n]}
                      content={`${n}最小值`}
                    />
                  </>
                ))}
              </>
            )}
          </Guide>
        </Chart>
        <div style={{ width: 'calc(100% - 30px)' }}>
          {slider && (
            <Slider
              width='auto'
              height={26}
              padding={[20, 80, 95]}
              start={this.ds.state.start}
              end={this.ds.state.end}
              xAxis='t'
              yAxis='v'
              scales={{
                t: {
                  type: 'time',
                  tickCount: 10,
                  mask: timeFormat || 'MM-DD HH:mm'
                }
              }}
              data={dv}
              backgroundChart={{
                type: type || 'line'
              }}
              onChange={this.onChange}
            />
          )}
        </div>
      </Wrapper>
    )
  }

  onChange = obj => {
    // reset minVal and maxVal
    this.minVals = this.names.reduce((pre, curr) => {
      pre[curr] = { t: -1, v: Number.MAX_SAFE_INTEGER, n: curr }
      return pre
    }, {})
    this.maxVals = this.names.reduce((pre, curr) => {
      pre[curr] = { t: -1, v: Number.MIN_SAFE_INTEGER, n: curr }
      return pre
    }, {})

    this.sumVals = this.names.reduce((pre, curr) => {
      pre[curr] = { t: -1, v: 0, n: curr }
      return pre
    }, {})

    this.pointsVals = this.names.reduce((pre, curr) => {
      pre[curr] = { t: -1, v: 0, n: curr }
      return pre
    }, {})

    const { startValue, endValue } = obj
    this.ds.setState('start', startValue)
    this.ds.setState('end', endValue)
  }
}
