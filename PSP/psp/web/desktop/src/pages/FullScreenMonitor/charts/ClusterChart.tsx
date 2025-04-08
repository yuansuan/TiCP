import React from 'react'
import { useEffect, useRef } from 'react'
import * as echarts from 'echarts'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'

const ClusterChart = observer(() => {
  const chartRef = useRef(null)
  const store = useStore()
  const { clusterInfo } = store
  
  useEffect(() => {
    let value = Number((clusterInfo?.availableNodeNum / clusterInfo?.totalNodeNum).toFixed(2)) * 100 || 100
    let angle = 0
    const chart = echarts.init(chartRef.current)
    const option = {
      animation: false,
      title: {
        text: '{a|' + value + '}{c|%\r\n集群可用率}',
        x: 'center',
        y: 'center',
        textStyle: {
          rich: {
            a: {
              fontSize: 48,
              color: '#29EEF3'
            },
            c: {
              fontSize: 20,
              color: '#ffffff',
              padding: [5, 0]
            }
          }
        }
      },

      series: [
        // 紫色
        {
          name: 'ring5',
          type: 'custom',
          coordinateSystem: 'none',
          renderItem: function (params, api) {
            return {
              type: 'arc',
              shape: {
                cx: api.getWidth() / 2,
                cy: api.getHeight() / 2,
                r: (Math.min(api.getWidth(), api.getHeight()) / 2) * 0.9,
                startAngle: ((0 + angle) * Math.PI) / 180,
                endAngle: ((90 + angle) * Math.PI) / 180
              },
              style: {
                stroke: '#8383FA',
                fill: 'transparent',
                lineWidth: 1.5
              },
              silent: true
            }
          },
          data: [0]
        },
        {
          name: 'ring5', //紫点
          type: 'custom',
          coordinateSystem: 'none',
          renderItem: function (params, api) {
            let x0 = api.getWidth() / 2
            let y0 = api.getHeight() / 2
            let r = (Math.min(api.getWidth(), api.getHeight()) / 2) * 0.9
            let point = getCirlPoint(x0, y0, r, 90 + angle)
            return {
              type: 'circle',
              shape: {
                cx: point.x,
                cy: point.y,
                r: 4
              },
              style: {
                stroke: '#8450F9', //绿
                fill: '#8450F9'
              },
              silent: true
            }
          },
          data: [0]
        },
        // 蓝色
        {
          name: 'ring5',
          type: 'custom',
          coordinateSystem: 'none',
          renderItem: function (params, api) {
            return {
              type: 'arc',
              shape: {
                cx: api.getWidth() / 2,
                cy: api.getHeight() / 2,
                r: (Math.min(api.getWidth(), api.getHeight()) / 2) * 0.9,
                startAngle: ((180 + angle) * Math.PI) / 180,
                endAngle: ((270 + angle) * Math.PI) / 180
              },
              style: {
                stroke: '#4386FA',
                fill: 'transparent',
                lineWidth: 1.5
              },
              silent: true
            }
          },
          data: [0]
        },
        {
          name: 'ring5', // 蓝点
          type: 'custom',
          coordinateSystem: 'none',
          renderItem: function (params, api) {
            let x0 = api.getWidth() / 2
            let y0 = api.getHeight() / 2
            let r = (Math.min(api.getWidth(), api.getHeight()) / 2) * 0.9
            let point = getCirlPoint(x0, y0, r, 180 + angle)
            return {
              type: 'circle',
              shape: {
                cx: point.x,
                cy: point.y,
                r: 4
              },
              style: {
                stroke: '#4386FA',
                fill: '#4386FA'
              },
              silent: true
            }
          },
          data: [0]
        },
        {
          name: 'ring5',
          type: 'custom',
          coordinateSystem: 'none',
          renderItem: function (params, api) {
            return {
              type: 'arc',
              shape: {
                cx: api.getWidth() / 2,
                cy: api.getHeight() / 2,
                r: (Math.min(api.getWidth(), api.getHeight()) / 2) * 1,
                startAngle: ((270 + -angle) * Math.PI) / 180,
                endAngle: ((40 + -angle) * Math.PI) / 180
              },
              style: {
                stroke: '#0CD3DB',
                fill: 'transparent',
                lineWidth: 1.5
              },
              silent: true
            }
          },
          data: [0]
        },
        // 橘色
        {
          name: 'ring5',
          type: 'custom',
          coordinateSystem: 'none',
          renderItem: function (params, api) {
            return {
              type: 'arc',
              shape: {
                cx: api.getWidth() / 2,
                cy: api.getHeight() / 2,
                r: (Math.min(api.getWidth(), api.getHeight()) / 2) * 1,
                startAngle: ((90 + -angle) * Math.PI) / 180,
                endAngle: ((220 + -angle) * Math.PI) / 180
              },
              style: {
                stroke: '#FF8E89',
                fill: 'transparent',
                lineWidth: 1.5
              },
              silent: true
            }
          },
          data: [0]
        },
        {
          name: 'ring5', // 橘点
          type: 'custom',
          coordinateSystem: 'none',
          renderItem: function (params, api) {
            let x0 = api.getWidth() / 2
            let y0 = api.getHeight() / 2
            let r = (Math.min(api.getWidth(), api.getHeight()) / 2) * 1
            let point = getCirlPoint(x0, y0, r, 90 + -angle)
            return {
              type: 'circle',
              shape: {
                cx: point.x,
                cy: point.y,
                r: 4
              },
              style: {
                stroke: '#FF8E89',
                fill: '#FF8E89'
              },
              silent: true
            }
          },
          data: [0]
        },
        {
          name: 'ring5',
          type: 'custom',
          coordinateSystem: 'none',
          renderItem: function (params, api) {
            let x0 = api.getWidth() / 2
            let y0 = api.getHeight() / 2
            let r = (Math.min(api.getWidth(), api.getHeight()) / 2) * 1
            let point = getCirlPoint(x0, y0, r, 270 + -angle)
            return {
              type: 'circle',
              shape: {
                cx: point.x,
                cy: point.y,
                r: 4
              },
              style: {
                stroke: '#0CD3DB',
                fill: '#0CD3DB'
              },
              silent: true
            }
          },
          data: [0]
        },
        {
          name: 'pie',
          type: 'pie',
          radius: ['75%', '58%'],
          silent: true,
          clockwise: true,
          startAngle: 90,
          z: 0,
          zlevel: 0,
          label: {
            position: 'center'
          },
          data: [
            {
              value: value,
              name: '',
              itemStyle: {
                color: {
                  // 完成的圆环的颜色
                  colorStops: [
                    {
                      offset: 0,
                      color: '#A098FC' // 0% 处的颜色
                    },
                    {
                      offset: 0.3,
                      color: '#4386FA' // 0% 处的颜色
                    },
                    {
                      offset: 0.6,
                      color: '#4FADFD' // 0% 处的颜色
                    },
                    {
                      offset: 0.8,
                      color: '#0CD3DB' // 100% 处的颜色
                    },
                    {
                      offset: 1,
                      color: '#646CF9' // 100% 处的颜色
                    }
                  ]
                }
              }
            },
            {
              value: 100 - value,
              name: '',
              label: {
                show: false
              },
              itemStyle: {
                color: '#173164'
              }
            }
          ]
        },
        {
          name: 'pie',
          type: 'pie',
          radius: ['40%', '48%'],
          silent: true,
          clockwise: true,
          startAngle: 270,
          z: 0,
          zlevel: 0,
          label: {
            position: 'center'
          },
          data: [
            {
              value: value,
              name: '',
              itemStyle: {
                color: {
                  // 完成的圆环的颜色
                  colorStops: [
                    {
                      offset: 0,
                      color: '#00EDF3' // 0% 处的颜色
                    },
                    {
                      offset: 1,
                      color: '#646CF9' // 100% 处的颜色
                    }
                  ]
                }
              }
            },
            {
              value: 100 - value,
              name: '',
              label: {
                show: false
              },
              itemStyle: {
                color: '#173164'
              }
            }
          ]
        }
      ]
    }

    function getCirlPoint(x0, y0, r, angle) {
      let x1 = x0 + r * Math.cos((angle * Math.PI) / 180)
      let y1 = y0 + r * Math.sin((angle * Math.PI) / 180)
      return {
        x: x1,
        y: y1
      }
    }

    function draw() {
      angle = angle + 3
      chart.setOption(option, true)
    }

    const timer = setInterval(function () {
      draw()
    }, 50)

    window.addEventListener('resize', () => {
      chart.resize()
    })

    return () => {
      clearInterval(timer)
      chart.dispose()
    }
  }, [clusterInfo])

  return <div ref={chartRef} style={{ width: '100%', height: '100%' }}></div>
})

export default ClusterChart
