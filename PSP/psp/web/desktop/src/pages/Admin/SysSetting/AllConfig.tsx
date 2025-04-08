import * as React from 'react'
import { observer } from 'mobx-react'
import { Collapse, Spin, Tooltip } from 'antd'
import sysConfig from '@/domain/SysConfig'
import MailServerConfig from './MailServerConfig'
import { observable, computed, action } from 'mobx'
import { INSTALL_TYPE } from '@/utils/const'
import JobConfig from './JobConfig'
import MailNotificationConfig from './MailNotificationConfig'
import ThreeMembersConfig from './ThreeMembersMgrConfig'
import Mark from 'mark.js'
import { PanelHeaderWrapper } from './style'

const { Panel } = Collapse
const markOpts = {
  exclude: ['[data-nomark]']
}

@observer
export default class AllConfig extends React.Component<any> {
  ref = null
  markInstances = {}

  constructor(props) {
    super(props)
    this.ref = React.createRef()
  }

  @observable loading = true
  @observable search = ''
  @observable panels = [
    {
      label: '本地作业',
      PanelComponent: JobConfig,
      isShow: true,
      key: 'JobConf',
      totals: 0
    },
    {
      label: '邮件服务器设置',
      PanelComponent: MailServerConfig,
      isShow: () => {
        return true
      },
      key: 'mailConf',
      totals: 0
    },
    {
      label: '邮件信息提醒设置',
      PanelComponent: MailNotificationConfig,
      isShow: () => {
        return true
      },
      key: 'MailNotificationConfig',
      totals: 0
    },
    {
      label: '三员管理配置',
      PanelComponent: ThreeMembersConfig,
      isShow: () => {
        return sysConfig.enableThreeMemberMgr
      },
      key: 'ThreeMembersConfig',
      totals: 0
    }
  ]

  @observable activeKeys: string[] = this.panels.map(p => p.key)

  @computed
  get pagePanels() {
    return this.panels.filter(panel => {
      // 处理 isShow 为函数或布尔值的情况
      return typeof panel.isShow === 'function' ? panel.isShow() : panel.isShow
    })
  }

  @computed
  get filterPanels() {
    return this.pagePanels.filter(
      panel => panel.totals > 0 || this.props.filterKey === ''
    )
  }

  get isAIO() {
    return sysConfig.installType === INSTALL_TYPE.aio
  }

  markSearchKey(value, p) {
    const markInstance = this.markInstances[p.label]

    if (!markInstance) return
    if (value === '') {
      markInstance.unmark()
      p.totals = 0
    } else {
      markInstance.unmark({
        done: () => {
          markInstance.mark(value, {
            done: nums => {
              p.totals = nums
            },
            ...markOpts
          })
        }
      })
    }
  }

  @action
  markSearchKeyAll() {
    this.panels.map((p, index) => {
      this.markSearchKey(this.props.filterKey, this.panels[index])
    })
  }

  setMarkInstances() {
    this.markInstances = this.filterPanels.reduce((pre, cur, index) => {
      pre[cur.label] = new Mark(
        this.ref.current.querySelectorAll(`.ant-collapse-item`)[index]
      )
      return pre
    }, {})
  }

  async componentDidMount() {
    await Promise.all([
      sysConfig.fetchMailServerConfig(),
      sysConfig.fetchMailInfoConfig()
    ])
    this.loading = false
  }

  componentDidUpdate(prevProps) {
    if (prevProps.filterKey !== this.props.filterKey) {
      this.markSearchKeyAll()

      setTimeout(() => {
        const markIns = new Mark(this.ref.current)
        if (!markIns) return
        if (this.props.filterKey === '') {
          markIns.unmark()
        } else {
          markIns.unmark({
            done: () => {
              markIns.mark(this.props.filterKey)
            },
            ...markOpts
          })
        }
      }, 0)
    }
  }

  componentWillUnmount() {
    this.markInstances = null
  }

  render() {
    return (
      <>
        {this.loading ? (
          <div className='loading'>
            <Spin />
          </div>
        ) : (
          <div ref={this.ref}>
            <Collapse
              bordered={false}
              onChange={keys => {
                this.activeKeys = keys as string[]
              }}
              activeKey={this.activeKeys}>
              {this.filterPanels.map((panel, index) => {
                const PanelComponent = panel.PanelComponent

                const Header = (
                  <PanelHeaderWrapper>
                    <Tooltip
                      getPopupContainer={() =>
                        this.ref.current.querySelectorAll(
                          `.ant-collapse-item .ant-collapse-header`
                        )[index]
                      }
                      placement={'right'}
                      visible={
                        !this.activeKeys.includes(panel.key) &&
                        panel.totals !== 0
                      }
                      title={
                        panel.totals !== 0 ? `${panel.totals}条结果` : null
                      }>
                      {panel.label}
                    </Tooltip>
                  </PanelHeaderWrapper>
                )

                return (
                  <Panel
                    forceRender
                    header={Header}
                    key={panel.key}
                    data-id={panel.key}>
                    <PanelComponent />
                  </Panel>
                )
              })}
            </Collapse>
          </div>
        )}
      </>
    )
  }
}
