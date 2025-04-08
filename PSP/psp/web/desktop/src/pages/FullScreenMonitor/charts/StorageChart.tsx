import React from 'react'
import { useEffect, useRef } from 'react'
import * as echarts from 'echarts'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'

const arrayToMap = (arr: Record<string, number>[]) => {
  return arr.reduce((map, item) => {
    let usedNum: number
    let freeNum: number
    const size = Object.values(item)[0]
    const key = Object.keys(item)[0] as string
    const nameVal = String(Object.values(item)[1])

    if (map.has(key)) {
      let data = map.get(key)
      nameVal === '已使用' ? (data.used = size) : (data.free = size)
    } else {
      nameVal === '已使用' ? (usedNum = size) : (freeNum = size)
      map.set(key, {
        used: usedNum,
        free: freeNum
      })
    }

    return map
  }, new Map<string, { used: number; free: number }>())
}

const StorageChart = observer(() => {
  const chartRef = useRef(null)
  const store = useStore()
  const { diskList } = store

  useEffect(() => {
    let usageRates = []
    let titleName = []
    let totals = []

    const diskMap = arrayToMap(diskList?.data)
    diskList?.fields.forEach(item => {
      titleName.push(item)

      let disk = diskMap.get(item)
      let totalSize = disk.used + disk.free
      totals.push(totalSize)

      let usageRate = Number((disk.used / totalSize).toFixed(2)) * 100
      usageRates.push(usageRate)
    })

    let defaultColor = ['#1089E7', '#8B78F6', '#56D0E3', '#ef6a6a', '#F8B448']
    const chart = echarts.init(chartRef.current)
    const option = {
      animation: false,
      grid: {
        left: '2%',
        right: '3%',
        top: '10%',
        bottom: '0%',
        containLabel: true
      },
      xAxis: {
        show: false
      },
      yAxis: [
        {
          show: true,
          data: titleName,
          inverse: true,
          axisLine: {
            show: false
          },
          splitLine: {
            show: false
          },
          axisTick: {
            show: false
          },
          axisLabel: {
            color: 'rgba(255,255,255,.6)',
            fontWeight: '700',
            fontSize: 14
          }
        },
        {
          show: true,
          inverse: true,
          data: totals,
          axisLabel: {
            color: 'rgba(255,255,255,.6)',
            fontWeight: '700',
            fontSize: 14,
            formatter: function (value: number) {
              return value + 'GB'
            }
          },
          axisLine: {
            show: false
          },
          splitLine: {
            show: false
          },
          axisTick: {
            show: false
          }
        }
      ],
      series: [
        {
          name: '条',
          type: 'bar',
          yAxisIndex: 0,
          data: usageRates,
          barWidth: 30,
          itemStyle: {
            borderRadius: 30,
            color: function (params) {
              return defaultColor[params.dataIndex]
            }
          },
          label: {
            show: true,
            position: 'inside',
            formatter: '{c}%',
            color: '#fff'
          }
        },
        {
          name: '框',
          type: 'bar',
          yAxisIndex: 1,
          barGap: '-100%',
          data: [100, 100, 100, 100, 100],
          barWidth: 40,
          itemStyle: {
            color: 'none',
            borderColor: '#00c1de',
            borderWidth: 3,
            borderRadius: 15
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
  }, [diskList])

  return <div ref={chartRef} style={{ width: '100%', height: '100%' }}></div>
})

export default StorageChart
