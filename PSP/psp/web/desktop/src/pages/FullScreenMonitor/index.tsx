import React, { createContext, useContext, useEffect, useState } from 'react'
import moment from 'moment'
import { ScreenFullMonitorWrapper } from './style'
import { FullscreenExitOutlined } from '@ant-design/icons'
import { Context, useModel, useStore } from './store'
import { observer } from 'mobx-react-lite'
import MemeryChart from './charts/MemeryChart'
import ClusterChart from './charts/ClusterChart'
import NodeChart from './charts/NodeChart'
import IOChart from './charts/IOChart'
import ProjectJobChart from './charts/ProjectJobChart'
import StorageChart from './charts/StorageChart'
import UserJobChart from './charts/UserJobChart'
import AppJobChart from './charts/AppJobChart'

export const FullScreenContext = createContext(() => {})

export const ScreenFullMonitor = observer(() => {
  const dateFormat = 'YYYY年MM月DD日-HH时mm分ss秒'
  const bgImg = require('@/assets/images/monitor/bg.jpg')
  const headBgImg = require('@/assets/images/monitor/head_bg.png')
  const lineImg = require('@/assets/images/monitor/line.png')
  const mapImg = require('@/assets/images/monitor/map.png')
  const maskImg = require('@/assets/images/monitor/mask.png')

  const store = useStore()
  const { jobStatusList, clusterInfo } = store
  const [currentTime, setCurrentTime] = useState(moment().format(dateFormat))

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(moment().format(dateFormat))
    }, 1000)

    return () => {
      clearInterval(timer)
    }
  }, [])

  useEffect(() => {
    store.refresh()

    const timer = setInterval(() => {
      store.refresh()
    }, 600000)

    return () => {
      clearInterval(timer)
    }
  }, [])

  const exitFullScreen = useContext(FullScreenContext)

  const exitFullScreenHandler = () => {
    if (document.exitFullscreen) {
      document
        .exitFullscreen()
        .then(() => {
          exitFullScreen()
        })
        .catch(err => {
          console.error(
            `Error attempting to disable full-screen mode: ${err.message} (${err.name})`
          )
        })
    }
  }

  return (
    <ScreenFullMonitorWrapper
      bgImg={bgImg}
      headBgImg={headBgImg}
      lineImg={lineImg}
      mapImg={mapImg}
      maskImg={maskImg}>
      <div className='dataVis'>
        <header className='header'>
          <h1>仿真云平台总览大屏</h1>
          <div className='showTime'>当前时间：{currentTime}</div>
          <div className='fullScreen'>
            <div className='button'>
              <FullscreenExitOutlined
                onClick={exitFullScreenHandler}
                title='退出'
                rev='true'
              />
            </div>
          </div>
        </header>

        <section className='content'>
          <section className='column'>
            <div className='panel'>
              <h2>CPU/内存利用率</h2>
              <div className='chart'>
                <MemeryChart />
              </div>
              <div className='panelFooter'></div>
            </div>

            <div className='panel'>
              <h2>磁盘IO速率</h2>
              <div className='chart'>
                <IOChart />
              </div>
              <div className='panelFooter'></div>
            </div>

            <div className='panel'>
              <h2>存储用量</h2>
              <div className='chart'>
                <StorageChart />
              </div>
              <div className='panelFooter'></div>
            </div>
          </section>

          <section className='column'>
            <div className='summary'>
              <div className='summaryHd'>
                <ul>
                  {['Running', 'Pending', 'Completed', 'Failed'].map(status => (
                    <li key={status}>
                      {jobStatusList?.list?.find(item => item.status === status)
                        ?.job_count || 0}
                    </li>
                  ))}
                </ul>
              </div>

              <div className='summaryBd'>
                <ul>
                  <li>运行作业</li>
                  <li>等待作业</li>
                  <li>完成作业</li>
                  <li>失败作业</li>
                </ul>
              </div>
            </div>

            <div className='cluster'>
              <div className='sphere' />
              <div className='mask' />
              <div className='smallBox leftBottom'>
                <h3>集群节点</h3>
                <ul>
                  <li>
                    总体节点：<span>{clusterInfo?.totalNodeNum}</span>
                  </li>
                  <li>
                    可用节点：<span>{clusterInfo?.availableNodeNum}</span>
                  </li>
                </ul>
              </div>
              <div className='smallBox rightBottom'>
                <h3>集群核数</h3>
                <ul>
                  <li>
                    总体核数：<span>{clusterInfo?.cores}</span>
                  </li>
                  <li>
                    空闲核数：<span>{clusterInfo?.freeCores}</span>
                  </li>
                </ul>
              </div>
              <ClusterChart />
            </div>

            <div className='panel'>
              <h2>
                节点列表
                <ul>
                  <li>
                    <span className='circle notOk'></span>
                    <span>不可用</span>
                  </li>
                  <li>
                    <span className='circle ok'></span>
                    <span>可用</span>
                  </li>
                </ul>
              </h2>
              <div className='chart'>
                <NodeChart />
              </div>
              <div className='panelFooter'></div>
            </div>
          </section>

          <section className='column'>
            <div className='panel'>
              <h2>用户作业数-TOP5</h2>
              <div className='chart'>
                <UserJobChart />
              </div>
              <div className='panelFooter'></div>
            </div>

            <div className='panel'>
              <h2>项目(成员|作业|核时)-TOP5</h2>
              <div className='chart'>
                <ProjectJobChart />
              </div>
              <div className='panelFooter'></div>
            </div>

            <div className='panel'>
              <h2>应用作业数-TOP5</h2>
              <div className='chart'>
                <AppJobChart />
              </div>
              <div className='panelFooter'></div>
            </div>
          </section>
        </section>
      </div>
    </ScreenFullMonitorWrapper>
  )
})

export default function ScreenFullMonitorWithStore() {
  const model = useModel()

  return (
    <Context.Provider value={model}>
      <ScreenFullMonitor />
    </Context.Provider>
  )
}
