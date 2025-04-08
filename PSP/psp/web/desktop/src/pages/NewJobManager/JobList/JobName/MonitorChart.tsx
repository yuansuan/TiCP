import React, { useEffect, useRef } from 'react'
import { Chart, Geom, Slider, Legend } from 'bizcharts'

let initialSlider = [0, 1]
let initChartIns = null

export const MonitorChart = ({ data }) => {
  const sliderRef = useRef(initialSlider)
  const chartInsRef = useRef(initChartIns)
  const chartData = []
  data.forEach(y => {
    y.items.forEach(d => {
      chartData.push({
        x: d.kv[0],
        name: y.key,
        y: d.kv[1]
      })
    })
  })
  function handleSliderChange([start, end]) {
    sliderRef.current = [start, end]
    const c = chartInsRef.current
    const sliderController = c.getController('slider')
    const xScale = c.getXScale()
    const xValues = xScale.values

    const { minText, maxText } = sliderController.getMinMaxText(start, end)
    const sIdx = xValues.indexOf(minText)
    const eIdx = xValues.indexOf(maxText)
    c.views.forEach(view => {
      view.filter(xScale.field, val => {
        const idx = xValues.indexOf(val)
        return idx >= sIdx && idx <= eIdx
      })
    })
  }

  useEffect(() => {
    if (chartInsRef.current) {
      // hack slider sliderchange event invalid when new data come
      chartInsRef.current
        .getController('slider')
        .slider.component.off('sliderchange', handleSliderChange)

      chartInsRef.current
        .getController('slider')
        .slider.component.on('sliderchange', handleSliderChange)
    }

    return () => {
      if (chartInsRef.current) {
        chartInsRef.current
          .getController('slider')
          .slider.component.off('sliderchange', handleSliderChange)
      }
    }
  }, [data.length])

  useEffect(() => {
    return () => {
      if (chartInsRef.current) {
        chartInsRef.current
          .getController('slider')
          .slider.component.off('sliderchange', handleSliderChange)

        // clear chart instance
        chartInsRef.current = null
        sliderRef.current = [0, 1]
      }
    }
  }, [])

  return (
    <Chart
      autoFit
      height={600}
      data={chartData}
      onGetG2Instance={c => {
        chartInsRef.current = c
        const sliderController = c.getController('slider')
        sliderController.slider.component.on('sliderchange', handleSliderChange)
      }}>
      <Legend
        marker={{
          symbol: 'square',
          style: {
            fill: null
          }
        }}
        position='bottom'
      />
      <Geom
        type='line'
        position='x*y'
        color='name'
        animate={{
          appear: {
            animation: 'pahtIn'
          },
          enter: {
            animation: 'fadeIn'
          }
        }}
      />
      <Slider start={sliderRef.current[0]} end={sliderRef.current[1]} />
    </Chart>
  )
}
