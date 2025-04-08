import React from 'react'
import { useEffect, useRef } from 'react'
import * as echarts from 'echarts'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'
import { formatTimestamp } from '@/utils/formatter'

const IOChart = observer(() => {
  const chartRef = useRef(null)
  const store = useStore()
  const { ioAvgList } = store

  useEffect(() => {
    const writerIOs = []
    const readerIOs = []
    ioAvgList?.forEach(item => {
      if (item.n === '写速率') {
        writerIOs.push(item.d)
      } else if (item.n === '读速率') {
        readerIOs.push(item.d)
      }
    })

    const ts = readerIOs[0]?.map(item => formatTimestamp(item.t)) || [0]
    const readerIOvs = readerIOs[0]?.map(item => (item.v / 1024).toFixed(3)) || [0]
    const writerIOvs = writerIOs[0]?.map(item => (item.v / 1024).toFixed(3)) || [0]

    const chart = echarts.init(chartRef.current)
    const option = {
      animation: false, 
      grid: {
        left: '2%',
        right: '4%',
        top: '15%',
        bottom: '0%',
        containLabel: true
      },
      tooltip: {
        trigger: 'axis'
      },
      legend: {
        icon: 'rect',
        itemWidth: 14,
        itemHeight: 5,
        itemGap: 13,
        right: '3%',
        textStyle: {
          fontSize: 12,
          color: 'rgba(255,255,255,.6)'
        },
        data: ['写速率', '读速率']
      },
      xAxis: [
        {
          type: 'category',
          boundaryGap: false,
          axisLabel: {
            color: 'rgba(255,255,255,.6)'
          },
          axisLine: {
            show: false
          },
          data: ts
        }
      ],
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: '{value}MB/s',
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
          name: '写速率',
          type: 'line',
          symbol: 'circle',
          showSymbol: false,
          smooth: true,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 1, 0, 0, [
              {
                offset: 0,
                color: 'rgba(7,44,90,0.5)'
              },
              {
                offset: 1,
                color: 'rgba(0,146,246,0.9)'
              }
            ])
          },
          markPoint: {
            itemStyle: {
              color: 'red'
            }
          },
          data: writerIOvs
        },
        {
          name: '读速率',
          type: 'line',
          symbol: 'circle',
          showSymbol: false,
          smooth: true,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 1, 0, 0, [
              {
                offset: 0,
                color: 'rgba(7,44,90,0.5)'
              },
              {
                offset: 1,
                color: 'rgba(114,144,89,0.9)'
              }
            ])
          },
          data: readerIOvs
        }
      ]
    };

    chart.setOption(option)
    window.addEventListener('resize', () => {
      chart.resize()
    })

    return () => {
      chart.dispose()
    }
  }, [ioAvgList])

  return <div ref={chartRef} style={{ width: '100%', height: '100%' }}></div>
})

export default IOChart
