import * as React from 'react'
import { Tabs, message } from 'antd'
import { computed } from 'mobx'
import { observer } from 'mobx-react'
import Node from '@/domain/NodeMG/NodeManager/Node'
import { nodeManager, NodeActionLabel } from '@/domain/NodeMG'
import { sysConfig } from '@/domain'

import { BackButton } from '@/components'
import { history } from '@/utils'
// import JobListByHost from '@/components/JobList/JobListByHost'
import { Button, Modal } from '@/components'
import { Wrapper, SummaryWrapper, TabAreaWrapper } from './style'
import { NodeDetail } from '../common'
import {
  PBSCanOpenStatus,
  PBSCanCloseStatus,
  LSFCanCloseStatus,
  LSFCanOpenStatus
} from '../commonStatus'

const { TabPane } = Tabs

@observer
export default class NodeDetailWrapper extends React.Component<any, any> {
  state = {
    activeTabVal: 'NodeDetail'
  }

  onTabChange = activeTabVal => {
    this.setState({ activeTabVal })
  }

  private tabs = [
    {
      key: '基本信息',
      val: 'NodeDetail',
      component: NodeDetail
    }
    // {
    //   key: '节点作业',
    //   val: 'NodeJobs',
    //   component: JobListByHost
    // }
  ]

  get hostname() {
    const { name } = this.props.match.params
    return name
  }

  @computed
  get node(): Node | any {
    const { name } = this.props.match.params
    return nodeManager.nodeList.filter(n => n.node_name === name)[0]
  }

  backNodeListPage = () => {
    history.push('/node')
  }

  async componentDidMount() {
    await nodeManager.getNodeList()
  }

  private open = () => {
    Modal.show({
      title: '接受作业',
      content: '确定改变所选机器为接受作业状态吗？',
      onOk: async () => {
        await nodeManager.operate([this.hostname], NodeActionLabel.open)
        message.success('操作正在进行，请稍后等待刷新', 5)
      }
    })
  }

  private close = () => {
    Modal.show({
      title: '拒绝作业',
      content: '确定改变所选机器为拒绝作业状态吗？',
      onOk: async () => {
        await nodeManager.operate([this.hostname], NodeActionLabel.close)
        message.success('操作正在进行，请稍后等待刷新', 5)
      }
    })
  }
  get isEnableOpen() {
    return  this.node?.status.includes(PBSCanOpenStatus)
  }

  get isEnableClose() {
    return PBSCanCloseStatus.some(n => this.node?.status === n)
  }

  render() {
    return (
      <Wrapper>
        <BackButton
          onClick={this.backNodeListPage}
          title='返回集群管理'
          style={{
            fontSize: 20
          }}>
          <div>
            <span className='title'>集群详情</span>
            <span className='name'> {this.hostname}</span>
          </div>
        </BackButton>

        <SummaryWrapper>
          <div>
            <Button
              style={{ marginRight: 20 }}
              ghost
              type='primary'
              disabled={!this.isEnableOpen}
              onClick={this.open}>
              接受作业
            </Button>
            <Button
              ghost
              type='primary'
              disabled={!this.isEnableClose}
              onClick={this.close}>
              拒绝作业
            </Button>
          </div>
        </SummaryWrapper>
        <TabAreaWrapper>
          <Tabs
            defaultActiveKey='NodeDetail'
            size='large'
            onChange={this.onTabChange}>
            {this.tabs.map(
              tab =>
                tab && (
                  <TabPane tab={tab.key} key={tab.val}>
                    {/* 通过这种方式强制刷新 tab pane */}
                    {this.state.activeTabVal === tab.val ? (
                      <tab.component
                        node={this.node || {}}
                        nodeName={this.hostname}
                        hasActions={{
                          buttons: false,
                          search: true
                        }}
                        isPopupDetail={true}
                      />
                    ) : (
                      <div />
                    )}
                  </TabPane>
                )
            )}
          </Tabs>
        </TabAreaWrapper>
      </Wrapper>
    )
  }
}
