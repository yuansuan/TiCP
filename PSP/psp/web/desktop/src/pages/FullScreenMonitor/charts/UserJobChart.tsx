import React from 'react'
import { useEffect, useRef } from 'react'
import * as echarts from 'echarts'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'

const UserJobChart = observer(() => {
  const chartRef = useRef(null)
  const store = useStore()
  const { userJobList } = store

  useEffect(() => {
    let maxData = []
    let yData = []
    let data = []
    let seriesData = []
    let color = ['#6BF1BF', '#C7F895', '#E6D349', '#F8A065', '#FF6B5F']

    const maxNum = Math.max(...userJobList.map(item => item.num));
    userJobList?.forEach(item => {
      yData.push(item.user_name)
      data.push(item.num)
      maxData.push(maxNum)
    });
    yData.reverse()
    data.reverse()

    const dataSize = data.length
    data?.forEach((item, index) => {
      seriesData.push({
        name: '',
        value: item,
        itemStyle: {
          color: color[5-(dataSize-index)],
          borderRadius: 12
        }
      })
    })
    
    const chart = echarts.init(chartRef.current)
    const option = {
      animation: false,
      legend: {
        show: false
      },
      grid: {
        left: '2%',
        right: '0%',
        top: '5%',
        bottom: '0%',
        containLabel: true
      },
      xAxis: {
        type: 'value',
        axisTick: {
          show: false
        },
        axisLine: {
          show: false
        },
        splitLine: {
          show: false
        },
        axisLabel: {
          show: false
        }
      },
      yAxis: [
        {
          type: 'category',
          axisTick: {
            show: false
          },
          axisLine: {
            show: false,
            lineStyle: {
              color: '#363e83'
            }
          },
          axisLabel: {
            inside: false,
            color: 'rgba(255,255,255,.6)',
            fontWeight: '700',
            fontSize: 14
          },
          data: yData
        },
        {
          type: 'category',
          axisLine: {
            show: false
          },
          axisTick: {
            show: false
          },
          axisLabel: {
            show: false
          },
          splitArea: {
            show: false
          },
          splitLine: {
            show: false
          },
          data: yData
        },
        {
          type: 'category',
          axisLine: {
            show: false
          },
          axisTick: {
            show: false
          },
          axisLabel: {
            show: false
          },
          splitArea: {
            show: false
          },
          splitLine: {
            show: false
          },
          data: yData
        }
      ],
      series: [
        {
          name: '',
          type: 'bar',
          stack: '1',
          yAxisIndex: 0,
          data: seriesData,
          barWidth: 25,
          label: {
            show: true,
            color: '#3752ca',
            fontSize: 14,
            fontWeight: 700,
            padding: [0, 0, 0, 0],
            position: 'inside',
            distance: 0,
            formatter: function (params) {
              return data[params.dataIndex]
            }
          },
          z: 3
        },
        {
          name: '',
          type: 'bar',
          yAxisIndex: 2,
          data: maxData,
          barWidth: 25,
          itemStyle: {
            color: 'rgba(20, 29, 98, 0)',
            borderRadius: 12
          },

          z: 5
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
  }, [userJobList])

  return <div ref={chartRef} style={{ width: '100%', height: '100%' }}></div>
})

export default UserJobChart
