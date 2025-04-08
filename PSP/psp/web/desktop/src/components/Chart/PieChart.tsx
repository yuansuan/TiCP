import React from 'react'
import { Chart, Geom, Tooltip, Coord, Label, Legend, Guide } from 'bizcharts'
import { ChartTitle } from './style'
import { PRECISION } from '@/domain/common'
import { roundNumber } from '@/utils/formatter'
import { date } from './commonShowDates'

interface Item {
  key: string
  value: number
}

interface Ilegend {
  position?: bizcharts.LegendPositionType
  offsetX?: number
  offsetY?: number
}

interface IProps {
  data: Item[]
  title?: string
  height?: number
  width?: number
  unit?: string
  color?: (key: string) => string
  legend?: Ilegend
  theme?: string
  hideLabel?: boolean
  hideLegend?: boolean
  padding?: number[] | string
  useGuide?: boolean
  guideHtml?: string
  tooltipFomatter?: (key: string, value: number) => any
  showDates?: any
}

export default class PieChart extends React.Component<IProps> {
  render() {
    const {
      data,
      title,
      height,
      width,
      unit,
      color,
      legend,
      theme,
      hideLabel,
      hideLegend,
      padding,
      useGuide,
      guideHtml,
      tooltipFomatter,
      showDates,
    } = this.props

    const pieColor = color ? ['key', color] : 'key'

    const sum = data.reduce((s, d) => s + d.value, 0)

    const tooltipTemp = `<li data-index={index}>
      <span style="background-color:{color};width:8px;height:8px;border-radius:50%;display:inline-block;margin-right:8px;"></span>
      {name}: {value} ${unit || ''}
    </li>`

    const defaultTooltipFormatter = (key, value) => ({
      name: key,
      value: roundNumber(value, PRECISION),
    })

    return (
      <Chart
        padding={padding || ['15%', '10%']}
        data={data}
        width={width}
        height={data.length !== 0 ? height || 400 : height - 30}
        forceFit
        theme={theme || 'default'}>
        {title && <ChartTitle>{title}</ChartTitle>}
        {showDates && date(showDates)}
        {data.length === 0 ? (
          <div style={{ textAlign: 'center' }}>暂无数据</div>
        ) : (
          <>
            <Coord type='theta' />
            {!hideLegend && (
              <Legend
                position={legend?.position || 'bottom'}
                offsetX={legend?.offsetX || 0}
                offsetY={legend?.offsetY || 0}
              />
            )}
            {useGuide && (
              <Guide>
                <Guide.Html
                  position={['105%', '30%']}
                  html={guideHtml}
                  alignX='middle'
                  alignY='middle'
                />
              </Guide>
            )}
            <Tooltip showTitle={false} itemTpl={tooltipTemp} />
            <Geom
              type='intervalStack'
              position='value'
              color={pieColor}
              tooltip={[
                'key*value',
                tooltipFomatter || defaultTooltipFormatter,
              ]}>
              {!hideLabel && (
                <Label
                  content='key'
                  formatter={(text, item, index) => {
                    let point = item.point
                    let value = point['value']
                    let percent =
                      roundNumber((value / sum) * 100, PRECISION) + '%'
                    return text + ' ' + (sum === 0 ? 0 : percent)
                  }}
                />
              )}
            </Geom>
          </>
        )}
      </Chart>
    )
  }
}
