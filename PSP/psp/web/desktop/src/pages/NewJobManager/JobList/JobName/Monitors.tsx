import React, { useEffect, useState } from 'react'
import { observer } from 'mobx-react-lite'
import { Tabs, Spin } from 'antd'
import { MonitorChart } from './MonitorChart'
import { Residual } from './Residual'
import { CloudGraphic } from './CloudGraphic'
import { Http as JobLogHttp } from '@/utils/JobLogHttp'

const { TabPane } = Tabs

type Props = {
  id: string // jobId
  residualVisible: boolean
  monitorVisible: boolean
  cloudGraphicVisible?: boolean
  jobState: number
  projectId: string
  userId: string
  jobRuntimeId: string
}

let preChartsData = [] // persist history data

export const Monitors = observer(
  ({
    id,
    residualVisible,
    monitorVisible,
    cloudGraphicVisible,
    jobState,
    projectId,
    userId,
    jobRuntimeId
  }: Props) => {
    const [chartsData, setChartsData] = useState([])
    const [graphicData, setGraphicData] = useState([])

    let initKey = 'Residual'

    if (residualVisible) {
      initKey = 'Residual'
    } else if (cloudGraphicVisible) {
      initKey = 'C0'
    } else if (monitorVisible) {
      initKey = 'M0'
    }

    const [activeKey, setActiveKey] = useState(initKey)

    const [errGraphic, setErrGraphic] = useState(null)
    const [err, setErr] = useState(null)

    const setData = chartsData => {
      chartsData.sort((a, b) => a.key - b.key)
      preChartsData = [...chartsData]
      setChartsData(chartsData)
    }

    const refreshCloudGraphicData = async () => {
      let res = null
      try {
        res = await JobLogHttp.get('/job/snapshots', {
          params: {
            job_id: id
          }
        })

        if (res.success) {
          const keys = Object.keys(res.data?.snapshots)
          setGraphicData(
            keys.map(k => ({
              key: k,
              items: res.data?.snapshots[k].map(
                path =>
                  `/job/snapshot?job_id=${id}&path=${path}`
              )
            }))
          )
        } else {
          setErrGraphic(true)
        }
      } catch (e) {
        setErrGraphic(true)
      } finally {
        setErrGraphic(true)
      }
      return res
    }

    useEffect(() => {
      if (!cloudGraphicVisible) return undefined

      let intervalId = null

      ;(async () => {
        const res = await refreshCloudGraphicData()

        if (res.success) {
          intervalId = setInterval(async () => {
            if (jobState === 2 || jobState === 3 || jobState === 4) {
              intervalId && clearInterval(intervalId)
              return
            }
            await refreshCloudGraphicData()
          }, 5000)
        }
      })()

      return () => {
        intervalId && clearInterval(intervalId)
      }
    }, [id])

    const renderMonitorChartTab = chartsData => {
      return chartsData.length !== 0 ? (
        <>
          {
            <TabPane tab={'监控项'} key={'M0'}>
              <MonitorChart data={chartsData} />
            </TabPane>
          }
        </>
      ) : (
        <>
          <TabPane tab={'监控项'} key={'M0'}>
            {err ? (
              '无法获取监控数据，请检查相关配置文件'
            ) : (
              <>
                监控数据加载中，请耐心等待 <Spin size='small' />
              </>
            )}
          </TabPane>
        </>
      )
    }

    const renderCloudGraphicTab = graphicData => {
      console.log('graphicData: ', graphicData);
      return graphicData.length !== 0 ? (
        <>
          {graphicData.map((data, index) => (
            <TabPane tab={`云图项: ${data.key}`} key={'C' + index}>
              <CloudGraphic key={data.key} data={data.items} />
            </TabPane>
          ))}
        </>
      ) : (
        <>
          <TabPane tab={'云图项'} key={'C0'}>
            {errGraphic ? (
              '无法获取云图，请检查相关配置文件'
            ) : (
              <>
                云图加载中，请耐心等待 <Spin size='small' />
              </>
            )}
          </TabPane>
        </>
      )
    }

    return (
      <Tabs
        activeKey={activeKey}
        onChange={activeKey => setActiveKey(activeKey)}>
        {/* {monitorVisible && renderMonitorChartTab(chartsData)} */}
        {residualVisible && (
          <TabPane tab={'残差图'} key={'Residual'}>
            <Residual id={id} />
          </TabPane>
        )}
        {cloudGraphicVisible && renderCloudGraphicTab(graphicData)}
      </Tabs>
    )
  }
)
