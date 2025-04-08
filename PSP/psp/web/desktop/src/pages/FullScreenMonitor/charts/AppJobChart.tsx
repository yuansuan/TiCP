import React from 'react'
import { useEffect, useRef } from 'react'
import * as echarts from 'echarts'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'

const AppJobChart = observer(() => {
  const chartRef = useRef(null)
  const store = useStore()
  const { appJobList } = store

  useEffect(() => {
    let legend = []
    let seriesData = []
    appJobList?.app_jobs?.map(item => {
      legend.push(item.app_name)
      seriesData.push({
        name: item.app_name,
        value: item.num
      })
    })

    const chart = echarts.init(chartRef.current)
    const option = {
      animation: false, 
      tooltip: {
        trigger: 'item',
        formatter: '{b} : {d}% {c}个'
      },
      grid: {
        left: '0%',
        right: '0%',
        top: '10%',
        bottom: '0%',
        containLabel: true
      },
      legend: {
        orient: 'vertical',
        show: true,
        right: '5%',
        y: 'center',
        itemWidth: 3,
        itemHeight: 25,
        itemGap: 10,
        textStyle: {
          color: '#7a8c9f',
          fontSize: 12,
          lineHeight: 15,
          rich: {
            percent: {
              color: '#fff',
              fontSize: 12
            }
          }
        },
        data: legend,
        formatter: function (name: string) {
          if (seriesData.length) {
            const item = seriesData.filter(item => item.name === name)[0]
            return `{name|${name}}\r\n{value|已提交 ${item.value} 个作业}`
          }
          return `{name|${name}}\r\n{value|已提交 0 个作业}`
        }
      },
      series: [
        {
          type: 'pie',
          radius: ['50%', '80%'],
          center: ['30%', '50%'],
          data: seriesData,
          label: {
            show: false
          }
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
  }, [appJobList])

  return <div ref={chartRef} style={{ width: '100%', height: '100%' }}></div>
})

export default AppJobChart
