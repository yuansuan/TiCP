import * as React from 'react'
import styled from 'styled-components'
import { observer } from 'mobx-react'
import GauageChart from './ClusterGauge'
import { Icon } from '@/components'
import { roundNumber } from '@/utils/formatter'
import { observable, computed, action } from 'mobx'
import { Popover, Radio } from 'antd'
import { PRECISION } from '@/domain/common'
// import { history } from '@/utils'
import { enlarge } from './Enlarge'
import { Scrollbars } from 'react-custom-scrollbars'
// import { currentUser } from '@/domain'
import { fullScreenInfo } from './fullScreenInfos'
import Guider from '@/components/Rack/Guider'
import { colors } from '@/components/Rack/const'
import { LinkTo } from './LinkTo'

const getColor = value => colors[Math.floor(value)] || '#999'

const getColorByType = (type, ut, memUsage) => {
  if (type === 'all') {
    return getColor(Math.max(ut, isNaN(memUsage) ? -1 : memUsage))
  } else if (type === 'ut') {
    return getColor(ut)
  } else {
    return getColor(memUsage)
  }
}

const toNodeAttr = (val, num) => {
  if (val === -1) return '--'
  if (typeof val === 'string') return roundNumber(val + 0, num)
  if (typeof val === 'number') return roundNumber(val, num)
  return val
}
// '#10E617', '#E6D610', '#E61010'
const okColor = 'linear-gradient(to right, #10E617, #E6D610, #E61010)'
const notOkColor = '#999'
const highLightColor = '#1a6eba'
const normalColor = '#797979'

const PopoverContentWrapper = styled.div`
  .attr {
    display: flex;
    .key {
      width: 160px;
    }
    .value {
      font-weight: 600;
    }
  }
`

const PopoverTitleWrapper = styled.div`
  display: flex;
  align-items: center;
`

const NodeWrapper = styled.div`
  .nodeInfo {
    display: flex;
    flex-wrap: wrap;
    overflow-y: auto;

    .nodeDetail {
      display: flex;
      align-content: center;
      align-items: center;
      width: 200px;
      height: fit-content;

      .ant-popover-inner {
        padding: 0 20px;
        .ant-popover-title {
          padding: 5px 0;
        }
        .ant-popover-inner-content {
          padding: 12px 0;
        }
      }

      .content {
        display: flex;
        flex-direction: column;
        width: 150px;
        .nodeName {
          font-weight: 700;
          font-style: normal;
          font-size: 12px;
          color: #515151;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          margin-bottom: 0px;
          cursor: pointer;
        }
        .otherInfos {
          margin-bottom: 10px;
          > p {
            margin: 0px;
            .infoLabel {
              font-size: 12px;
              color: #797979;
            }
          }
        }
      }
    }
  }
`
interface IWrapper {
  isFullScreen: boolean
}

const Wrapper = styled.div<IWrapper>`
  display: flex;
  width: 100%;
  height: ${props =>
    props.isFullScreen ? `${fullScreenInfo.firstRow}px` : '250px'};

  .cluster {
    width: 25%;
    border-right: 1px dashed rgba(240, 242, 245, 1);

    .title {
      display: inline;
      .engName {
        font-size: ${props =>
          props.isFullScreen ? `${fullScreenInfo.baseSize}px` : '14px'};
        color: #1a6eba;
        padding-left: 15px;
        cursor: pointer;
      }
    }

    .cpuTime {
      color: #333333;
      font-size: ${props =>
        props.isFullScreen ? `${fullScreenInfo.baseSize * 0.8}px` : '12px'};
      .number {
        color: #a1a1a1;
        padding-right: 5px;
      }
    }
    .chart {
      width: 150px;
      margin: 12px auto 0 auto;
    }
  }

  .node {
    padding-left: 25px;
    width: 75%;

    .status {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .title {
        padding-right: 15px;
      }

      .filter {
        display: flex;
        justify-content: space-between;

        .type {
          display: flex;
          margin: 0 10px;

          .ant-radio-group {
            font-size: 10px;
            padding: 0 10px;
          }

          .guide {
            position: relative;
            top: 3px;
            display: flex;
            font-size: 12px;
            .text {
            }
          }
        }

        .btn {
          font-size: ${props =>
            props.isFullScreen ? `${fullScreenInfo.baseSize * 0.8}px` : '12px'};
          padding-left: 10px;
          cursor: pointer;
        }
        .circle {
          width: 7px;
          height: 8px;
          border-radius: 50%;
          display: inline-block;
          margin-right: 8px;
        }
        .ok {
          background: ${okColor};
        }
        .notOk {
          background: ${notOkColor};
        }
      }
    }
    .nodeInfo {
      padding-top: 20px;
    }
  }
`

const nodeAttrs = [
  // 又不要展示该属性
  // {
  //   key: 'node_status',
  //   label: '节点状态',
  // },
  {
    key: 'cpu_percent',
    label: 'CPU利用率',
    formatter: val => `${toNodeAttr(val, PRECISION)} %`,
  },
  {
    key: 'available_mem',
    label: '内存(可用/最大)',
    formatter: (val, all?) =>
      `${toNodeAttr(val/1024, PRECISION)}/${toNodeAttr(
        all['max_mem']/1024,
        PRECISION
      )} GB`,
  },
  {
    key: 'free_swap',
    label: 'swap(可用/最大)',
    formatter: (val, all?) =>
      `${toNodeAttr(val/1024, PRECISION)}/${toNodeAttr(
        all['max_swap']/1024,
        PRECISION
      )} GB`,
  },
  {
    key: 'free_tmp',
    label: 'tmp(可用/最大)',
    formatter: (val, all?) =>
      `${toNodeAttr(val/1024, PRECISION)}/${toNodeAttr(
        all['max_tmp']/1024,
        PRECISION
      )} GB`,
  },
  {
    key: 'write_throughput',
    label: '磁盘吞吐量(读/写)',
    formatter: (val, all?) => `${toNodeAttr(all['read_throughput']/1024, PRECISION)}/${toNodeAttr(val/1024, PRECISION)} MB/s`,
  },
  {
    key: 'r1m',
    label: '负载(r1m/r5m/r15m)',
    formatter: (val, all?) =>
      `${toNodeAttr(val, PRECISION)}/${toNodeAttr(
        all['r5m'],
        PRECISION
      )}/${toNodeAttr(all['r15m'], PRECISION)}`,
  },
  // {
  //   key: 'pg',
  //   label: '页面速率',
  //   formatter: (val, all?) => `${toNodeAttr(val, PRECISION)} PG/s`,
  // },
  // {
  //   key: 'it',
  //   label: '空闲时间',
  //   formatter: val => `${toNodeAttr(val, PRECISION)} min`,
  // },
  {
    key: 'scheduler_status',
    label: '调度器状态',
  },
  // {
  //   key: 'ls',
  //   label: '登录用户数',
  // },
  // {
  //   key: 'n_disk',
  //   label: '磁盘数量',
  //   formatter: val => `${toNodeAttr(val, 0)}`,
  // },
]

function PopoverTitle({ options, colorByType }) {
  let memUsage =
    ((options['max_mem'] - options['available_mem']) / options['max_mem']) * 100

  let okColor = getColorByType(colorByType, options['cpu_percent'], memUsage)

  let color = options['node_status'] === 'Up' ? okColor : notOkColor

  return (
    <PopoverTitleWrapper>
      <Icon type='zhuji' style={{ color: color, fontSize: 24, padding: 4 }} />
      <span>{options['node_name']}</span>
    </PopoverTitleWrapper>
  )
}

function PopoverContent({ options }) {
  return (
    <PopoverContentWrapper>
      {nodeAttrs.map(({ key, label, formatter }) => {
        if (key === 'n_disk' && options['n_disk'] === 0) {
          return null
        } else {
          return (
            <div key={key} className='attr'>
              <div className='key'>{label}</div>
              <div className='value'>
                {formatter ? formatter(options[key], options) : options[key]}
              </div>
            </div>
          )
        }
      })}
    </PopoverContentWrapper>
  )
}

function Node({ options, colorByType }) {
  let memUsage = options['max_mem'] === 0 ? 0 :
    ((options['max_mem'] - options['available_mem']) / options['max_mem']) * 100

  let okColor = getColorByType(colorByType, options['cpu_percent'], memUsage)

  let color = options['node_status'] === 'Up' ? okColor : notOkColor

  // 与后端约定好，如果max_mem为0 memUsageStr 为 '0%'
  const memUsageStr = options['max_mem'] === 0 ? 
    '0%' : (isNaN(memUsage) ? '--%' : `${toNodeAttr(memUsage, PRECISION)}%`)

  const utStr =
    toNodeAttr(options['cpu_percent'], 0) !== '--'
      ? `${toNodeAttr(options['cpu_percent'], PRECISION)}%`
      : `--%`

  return (
    <div className='nodeDetail'>
      <Popover
        mouseEnterDelay={0.3}
        mouseLeaveDelay={0.2}
        placement='topLeft'
        title={<PopoverTitle options={options} colorByType={colorByType} />}
        content={<PopoverContent options={options} />}>
        <Icon type='zhuji' style={{ color: color, fontSize: 50, paddingRight: 8 }} />
      </Popover>
      <div className='content'>
        <p
          className='nodeName'
          // onClick={() => {
          //   if (currentUser.perms.includes('system-node_management')) {
          //     history.push(`/node/${options?.node_name}`)
          //   }
          // }}
          >
          {options['node_name']}
        </p>
        <div className='otherInfos'>
          <p>
            <span className='infoLabel'>CPU利用率: {utStr}</span>
          </p>
          <p>
            <span className='infoLabel'>内存利用率: {memUsageStr}</span>
          </p>
        </div>
      </div>
    </div>
  )
}

interface IPorps {
  nodesInfo: any
  color?: string
  isFullScreen?: boolean
}

const FILTER_STATUS_LABEL = {
  '': '全部节点',
  OK: '可用节点',
  notOK: '不可用节点',
}

@observer
export default class ClusterAndNodeListInfo extends React.Component<IPorps> {
  @observable filterStatus = ''
  @observable scrollTop = 0
  @observable colorByType: 'all' | 'ut' | 'mem' = 'all'

  ref = null
  intervalID = null

  constructor(props) {
    super(props)
    this.ref = React.createRef()
  }

  handleUpdate = () => {
    if (this.ref) {
      const { scrollTop, top } = this.ref.current.getValues()
      this.scrollTop = top === 1 ? 0 : scrollTop
    }
  }

  @action
  onStatusChange = status => {
    this.filterStatus = status
  }

  @computed
  get nodeList() {
    const list = this.props.nodesInfo?.nodeList || []

    const nodeList = list
      .sort((x, y) => {
        return x.node_name >= y.node_name ? -1 : 1
      })
      .sort((x, _) => {
        return x.node_status === 'Up' ? -1 : 1
      })

    if (this.filterStatus) {
      if (this.filterStatus === 'OK') {
        return nodeList.filter(node => node.node_status === 'Up')
      } else if (this.filterStatus === 'notOK') {
        return nodeList.filter(node => node.node_status !== 'Up')
      }
    } else {
      return nodeList
    }
  }

  nodeInfo = () => {
    return (
      <NodeWrapper>
        {this.nodeList.length !== 0 ? (
          <div className='nodeInfo'>
            {this.nodeList?.map(h => {
              return (
                <Node
                  key={h.node_name}
                  options={h}
                  colorByType={this.colorByType}
                />
              )
            })}
          </div>
        ) : (
          <div className='nodeInfo'>无节点信息</div>
        )}
      </NodeWrapper>
    )
  }

  highLightByStatus = (status: string) => {
    return {
      color: this.filterStatus === status ? highLightColor : normalColor,
    }
  }

  componentDidMount() {
    this.intervalID = setInterval(() => {
      if (fullScreenInfo.isFullScreen)
        this.ref?.current.scrollTop(this.scrollTop + 3)
    }, 100)
  }

  componentWillUnmount() {
    clearInterval(this.intervalID)
  }

  render() {
    const { nodesInfo } = this.props

    const nodesStatus = [
      {key: "可用", value: nodesInfo.clusterInfo.availableNodeNum}, 
      {key: "不可用", value: nodesInfo.clusterInfo.totalNodeNum - nodesInfo.clusterInfo.availableNodeNum}
    ];

    return (
      <Wrapper isFullScreen={fullScreenInfo.isFullScreen}>
        <div className='cluster'>
          <div className='title'>
            集群
            <LinkTo routerPath={'/sys/node'} render={props => (
              <span className='engName' onClick={() => props.goTo()}>{nodesInfo.clusterInfo?.clusterName || ' '}</span>
            )} />
          </div>
          <div className='cpuTime'>
            总核数： <span className='number'>{nodesInfo.clusterInfo?.cores || 0}个</span>
            <span className='number'>｜已用{nodesInfo.clusterInfo?.usedCores || 0}个</span>
            <span className='number'>
              ｜可用{nodesInfo.clusterInfo.freeCores || 0}个
            </span>
          </div>
          <div className='chart'>
            <GauageChart
              height={fullScreenInfo.chartHeight}
              data={nodesStatus}
              color={this.props.color}
            />
          </div>
        </div>
        <div className='node'>
          <div className='status'>
            <div className='title'>节点列表</div>
            <div className='filter'>
              <div className='type'>
                <Radio.Group
                  name='type'
                  value={this.colorByType}
                  onChange={e => (this.colorByType = e.target.value)}
                  size={'small'}>
                  <Radio.Button value={'all'}>全部</Radio.Button>
                  <Radio.Button value={'ut'}>CPU利用率</Radio.Button>
                  <Radio.Button value={'mem'}>内存利用率</Radio.Button>
                </Radio.Group>
                <div className='guide'>
                  0%
                  <div style={{ width: 5 * 20 }}>
                    <Guider
                      height={5}
                      horizontal={true}
                      showText={false}
                      title={''}
                      style={{
                        height: 0,
                        transform: 'scale(0.8) rotate(270deg)',
                        position: 'absolute',
                        left: 15,
                        top: -1,
                      }}
                    />
                  </div>
                  100%
                </div>
              </div>
              |
              <div>
                <span className='btn' onClick={() => this.onStatusChange('')}>
                  <span style={this.highLightByStatus('')}>全部节点</span>
                </span>
                <span className='btn' onClick={() => this.onStatusChange('OK')}>
                  <span className='circle ok'></span>
                  <span style={this.highLightByStatus('OK')}>可用</span>
                </span>
                <span
                  className='btn'
                  onClick={() => this.onStatusChange('notOK')}>
                  <span className='circle notOk'></span>
                  <span style={this.highLightByStatus('notOK')}>不可用</span>
                </span>
                {!this.props.isFullScreen &&
                  enlarge(
                    this.nodeInfo(),
                    FILTER_STATUS_LABEL[this.filterStatus]
                  )}
              </div>
            </div>
          </div>
          <Scrollbars
            ref={this.ref}
            onUpdate={this.handleUpdate}
            autoHeight
            autoHeightMax={
              this.props.isFullScreen
                ? `${fullScreenInfo.firstRow - fullScreenInfo.baseSize * 5}px`
                : '210px'
            }>
            {this.nodeInfo()}
          </Scrollbars>
        </div>
      </Wrapper>
    )
  }
}
