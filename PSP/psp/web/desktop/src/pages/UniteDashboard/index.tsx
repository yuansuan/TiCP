import * as React from 'react'
import { observer } from 'mobx-react'
import { observable, action, runInAction } from 'mobx'
import moment from 'moment'
import { Spin } from 'antd'
import { FullscreenOutlined, FullscreenExitOutlined } from '@ant-design/icons'
import { sysConfig } from '@/domain'
import { Wrapper } from './style'
import ClusterAndNodeListInfo from './ClusterAndNodeList'
import AppJobNumInfo from './AppJobNum'
import JobStatusInfo from './JobStatus'
import SortedChart from './SortedChart'
import UserJobChart from './UserJobChart'
import DashboardData from '@/domain/Dashboard'
import LineChart from '@/components/Chart/LineChart'
import Worker from '../../worker/dashboard.worker'
import { LinkTo } from './LinkTo'
import { fullScreenInfo } from './fullScreenInfos'
import ScreenFullMonitor, { FullScreenContext } from '@/pages/FullScreenMonitor'

const _eventName = 'dashboard'

const chartColor = '#1A6EBA'

@observer
export default class Dashboard extends React.Component<any> {
  @observable fullScreen = false
  @observable loading = false
  @observable jobInfoKey = Date.now()
  @observable dateRange: '1h' | '24h' = '24h'

  @observable onLineUsers = 0

  @observable resourceInfo = {
    userDatas: [],
    metric_cpu_ut_avg: [],
    metric_mem_ut_avg: [],
    metric_io_ut_avg: []
  }

  @observable nodesInfo = {
    clusterInfo: {
      clusterName: '',
      usedCores: 0,
      freeCores: 0,
      cores: 0,
      totalNodeNum: 0,
      availableNodeNum: 0
    },
    disks: { data: [], fields: [] },
    nodeList: []
  }

  @observable applicationInfo = {
    app_jobs: [],
    app_total: 0
  }

  @observable userJobInfo = []

  get dates() {
    return [moment().subtract(1, 'days'), moment()].map(m => m.valueOf())
  }

  worker = null
  ref = null
  screenRef = null

  constructor(props) {
    super(props)
    this.ref = React.createRef()
    this.screenRef = React.createRef()
  }

  fullscreenHandler = () => {
    fullScreenInfo.height = document.body.clientHeight
    fullScreenInfo.width = document.body.clientWidth
    fullScreenInfo.isFullScreen = document.fullscreenElement ? true : false
  }

  async componentDidMount() {
    try {
      await this.getAllInfos()
    } finally {
      this.loading = false
    }

    this.worker = new Worker()

    console.debug('send message: 启动 Dashboard interval')

    this.worker.postMessage({
      eventName: _eventName,
      eventData: {
        userId: localStorage.getItem('userId') || ''
      }
    })

    this.worker.addEventListener('message', event => {
      const { eventName, eventData } = event.data
      const { res } = eventData

      if (eventName === _eventName) {
        console.debug('receive Message:', res)

        runInAction(() => {
          this.nodesInfo = res[0].data || this.nodesInfo
          this.applicationInfo = res[1].data || this.applicationInfo
          this.resourceInfo = res[2].data || this.resourceInfo
          this.onLineUsers = res[3].data || this.onLineUsers
          this.userJobInfo = res[4]?.data || []
          this.jobInfoKey = Date.now()
        })
      }
    })

    document.addEventListener('fullscreenchange', this.fullscreenHandler)
  }

  componentWillUnmount() {
    this.worker?.terminate()
    document.removeEventListener('fullscreenchange', this.fullscreenHandler)
  }

  setDataByType = (res, type) => {
    res.status === 'fulfilled' ? (this[type] = res.value.data) : () => {}
  }

  setAllData = res => {
    runInAction(() => {
      this.setDataByType(res[0], 'nodesInfo')
      this.setDataByType(res[1], 'applicationInfo')
      this.setDataByType(res[2], 'resourceInfo')
      this.setDataByType(res[3], 'onLineUsers')
      this.setDataByType(res[4], 'userJobInfo')
    })
  }

  @action
  getAllInfos = async () => {
    //@ts-ignore
    this.loading = true
    const res = await Promise.allSettled([
      DashboardData.getDashboardInfo('ClUSTER_INFO', []),
      DashboardData.getDashboardInfo('SOFTWARE_INFO', this.dates),
      DashboardData.getDashboardInfo('RESOURCE_INFO', this.dates),
      DashboardData.getDashboardInfo('ONLINE_USERS', []),
      DashboardData.getDashboardInfo('USER_JOB_INFO', this.dates)
    ])
    this.loading = false

    this.setAllData(res)
  }

  toggleFullscreen = () => {
    if (this.loading) return

    const ele = this.screenRef.current

    if (!document.fullscreenElement) {
      ele
        .requestFullscreen()
        .then((this.fullScreen = true))
        .catch(err => {
          console.error(
            `Error attempting to enable full-screen mode: ${err.message} (${err.name})`
          )
        })
    } else {
      document.exitFullscreen()
    }
  }

  exitFullScreenHandler = () => {
    this.fullScreen = false
  }

  render() {
    const ICON = fullScreenInfo.isFullScreen
      ? FullscreenExitOutlined
      : FullscreenOutlined

    return (
      <div ref={this.screenRef}>
        {this.fullScreen ? (
          <FullScreenContext.Provider value={this.exitFullScreenHandler}>
            <ScreenFullMonitor />
          </FullScreenContext.Provider>
        ) : (
          <Wrapper ref={this.ref} isFullScreen={fullScreenInfo.isFullScreen}>
            <Spin tip='数据加载中...' spinning={this.loading}>
              <div
                className='pageTitle'
                style={{ justifyContent: 'flex-end', lineHeight: '20px' }}>
                <div className='name' style={{ display: 'none' }}>
                  {sysConfig.getPageHeader() || '集群监控'}
                </div>
                <div className='right'>
                  <div className='online'>
                    在线用户数:&nbsp;
                    <LinkTo
                      routerPath={'/sys/onlineUser'}
                      render={props => (
                        <span className='num' onClick={() => props.goTo()}>
                          &nbsp;{this.onLineUsers || 0}人
                        </span>
                      )}
                    />
                  </div>
                  <ICON
                    rev={'fullScreenInfo'}
                    title={
                      fullScreenInfo.isFullScreen
                        ? '退出全屏模式'
                        : '进入全屏模式'
                    }
                    className='btn'
                    disabled={true}
                    style={{
                      cursor: this.loading ? 'default' : 'pointer'
                    }}
                    onClick={this.toggleFullscreen}
                  />
                </div>
              </div>
              <div className='gridBody'>
                <div className='item head'>
                  <ClusterAndNodeListInfo
                    isFullScreen={fullScreenInfo.isFullScreen}
                    nodesInfo={this.nodesInfo}
                    color={chartColor}
                  />
                </div>
                <div className='item'>
                  <div className='title'>共享存储</div>
                  <SortedChart
                    padding={['auto', '15%', 'auto', 'auto'] as any}
                    data={this.nodesInfo.disks?.data}
                    fields={this.nodesInfo.disks?.fields}
                    height={fullScreenInfo.chartHeight}
                    unit={'GB'}
                  />
                </div>
                <div className='item'>
                  <div className='title util'>
                    CPU利用率
                    <div className='detail'>过去24小时CPU利用率</div>
                  </div>
                  <LineChart
                    height={fullScreenInfo.chartHeight}
                    data={this.resourceInfo?.metric_cpu_ut_avg || []}
                    unit={'%'}
                    timeFormat={'HH:mm'}
                    hideLegend
                    lineColor={'#0E75C8'}
                    min={0}
                  />
                </div>
                <div className='item'>
                  <div className='title util'>
                    内存利用率
                    <div className='detail'>过去24小时平均内存利用率</div>
                  </div>
                  <LineChart
                    height={fullScreenInfo.chartHeight}
                    data={this.resourceInfo?.metric_mem_ut_avg || []}
                    unit={'%'}
                    timeFormat={'HH:mm'}
                    hideLegend
                    lineColor={'#7ED321'}
                    min={0}
                  />
                </div>
                <div className='item footer'>
                  <JobStatusInfo
                    key={this.jobInfoKey}
                    dateRange={this.dateRange}
                    onDateRangeChange={dateRange => {
                      this.dateRange = dateRange
                    }}
                  />
                </div>
                <div className='item'>
                  <div className='title util'>
                    磁盘IO速率
                    <div className='detail'>过去24小时磁盘IO速率</div>
                  </div>
                  <LineChart
                    height={fullScreenInfo.chartHeight}
                    data={this.resourceInfo?.metric_io_ut_avg || []}
                    unit={'KB/s'}
                    timeFormat={'HH:mm'}
                    hideLegend
                  />
                </div>
                {!fullScreenInfo.isFullScreen && (
                  <>
                    <div className='item footer'>
                      <AppJobNumInfo
                        data={
                          this.applicationInfo?.app_jobs?.map((a, index) => ({
                            top5: index + 1,
                            key: a.app_name,
                            value: a.num
                          })) || []
                        }
                        total={this.applicationInfo?.app_total || 0}
                      />
                    </div>
                    <div className='item'>
                      <div className='title util'>
                        用户作业数
                        <div className='detail'>过去24小时用户作业数(TOP5)</div>
                      </div>
                      <UserJobChart
                        data={
                          this.userJobInfo?.map(u => ({
                            key: u.user_name,
                            value: u.num
                          })) || []
                        }
                        color={chartColor}
                      />
                    </div>
                  </>
                )}
              </div>
            </Spin>
          </Wrapper>
        )}
      </div>
    )
  }
}
