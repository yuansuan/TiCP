import React from 'react'
import { useEffect, useRef } from 'react'
import * as echarts from 'echarts'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'

const ProjectJobChart = observer(() => {
  const chartRef = useRef(null)
  const store = useStore()
  const { projectJobList } = store

  useEffect(() => {
    const userData = projectJobList.users
    const projectData = projectJobList.projects
    const jobData = projectJobList.jobs
    const cputimesData = projectJobList.cputimes.map(item => Number(item).toFixed(3))

    const chart = echarts.init(chartRef.current)
    const option = {
      animation: false, 
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
          crossStyle: {
            color: '#999'
          }
        }
      },
      grid: {
        left: '2%',
        right: '3%',
        top: '15%',
        bottom: '-5%',
        containLabel: true
      },
      legend: {
        itemWidth: 14,
        itemHeight: 5,
        itemGap: 13,
        data: ['项目成员', '提交作业', '使用核时'],
        textStyle: {
          fontSize: 12,
          color: 'rgba(255,255,255,.6)'
        }
      },
      xAxis: [
        {
          type: 'category',
          data: projectData,
          axisLabel: {
            color: 'rgba(255,255,255,.6)',
            fontSize: '12px',
            interval: 0,
            formatter: function (value, index) {
              return value.replace(/(.{1,12})/g, '$1\n')
            }
          },
          axisLine: {
            show: false
          }
        }
      ],
      yAxis: [
        {
          type: 'value',
          name: '个数',
          nameTextStyle: {
            color: 'rgba(255,255,255,.6)',
            fontSize: '12px'
          },
          axisLabel: {
            formatter: '{value}',
            color: 'rgba(255,255,255,.6)',
            fontSize: '12px'
          }
        },
        {
          type: 'value',
          name: '核时',
          nameTextStyle: {
            color: 'rgba(255,255,255,.6)',
            fontSize: '12px'
          },
          axisLabel: {
            formatter: '{value}',
            color: 'rgba(255,255,255,.6)',
            fontSize: '12px'
          },
          splitLine: {
            lineStyle: {
              color: 'rgba(255,255,255,.6)'
            }
          }
        }
      ],
      series: [
        {
          name: '项目成员',
          type: 'bar',
          barWidth: 25,
          tooltip: {
            valueFormatter: function (value) {
              return value + ' 个'
            }
          },
          itemStyle: {
            borderRadius: [15, 15, 0, 0]
          },
          data: userData
        },
        {
          name: '提交作业',
          type: 'bar',
          barWidth: 25,
          tooltip: {
            valueFormatter: function (value) {
              return value + ' 个'
            }
          },
          itemStyle: {
            borderRadius: [15, 15, 0, 0]
          },
          data: jobData
        },
        {
          name: '使用核时',
          type: 'line',
          yAxisIndex: 1,
          smooth: true,
          tooltip: {
            valueFormatter: function (value) {
              return value + ' 核时'
            }
          },
          data: cputimesData
        }
      ]
    }

    chart.setOption({
      loadingOption: {
        show: false
      },
      ...option
    })
    window.addEventListener('resize', () => {
      chart.resize()
    })

    return () => {
      chart.dispose()
    }
  }, [projectJobList])

  return <div ref={chartRef} style={{ width: '100%', height: '100%' }}></div>
})

export default ProjectJobChart
