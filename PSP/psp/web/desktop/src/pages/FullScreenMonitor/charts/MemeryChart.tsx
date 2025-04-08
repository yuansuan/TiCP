import React from 'react'
import { useEffect, useRef } from 'react'
import * as echarts from 'echarts'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'
import { formatTimestamp } from '@/utils/formatter'

const MemeryChart = observer(() => {
  const chartRef = useRef(null)
  const store = useStore()
  const { cpuAvgList, memAvgList } = store

  useEffect(() => {
    const ts = cpuAvgList?.map(item => formatTimestamp(item.t)) || [0]
    const cpuvs = cpuAvgList?.map(item => item.v.toFixed(3)) || [0]
    const memvs = memAvgList?.map(item => item.v.toFixed(3)) || [0]

    const chart = echarts.init(chartRef.current)
    const option = {
      animation: false,
      tooltip: {
        trigger: 'axis'
      },
      legend: {
        icon: 'rect',
        itemWidth: 14,
        itemHeight: 5,
        itemGap: 13,
        data: ['CPU', '内存'],
        right: '4%',
        textStyle: {
          fontSize: 12,
          color: 'rgba(255,255,255,.6)'
        }
      },
      grid: {
        left: '2%',
        right: '3%',
        top: '15%',
        bottom: '0%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        axisLabel: {
          color: 'rgba(255,255,255,.6)',
          fontSize: '12px'
        },
        axisLine: {
          show: false
        },
        data: ts
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: '{value}%',
          color: 'rgba(255,255,255,.6)',
          fontSize: '12px'
        },
        splitLine: {
          lineStyle: {
            color: 'rgba(255,255,255,.6)'
          }
        }
      },
      series: [
        {
          name: 'CPU',
          type: 'line',
          symbol: 'circle',
          showSymbol: false,
          smooth: true,
          data: cpuvs
        },
        {
          name: '内存',
          type: 'line',
          symbol: 'circle',
          showSymbol: false,
          smooth: true,
          data: memvs
        }
      ]
    }

    chart.setOption(option)
    window.addEventListener('resize', () => {
      chart.resize()
    })

    return () => {
      chart.dispose()
    }
  }, [cpuAvgList])

  return <div ref={chartRef} style={{ width: '100%', height: '100%' }}></div>
})

export default MemeryChart
