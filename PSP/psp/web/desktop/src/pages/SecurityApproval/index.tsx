import * as React from 'react'
import { observer } from 'mobx-react'
import { observable } from 'mobx'

import { Wrapper } from './style'
import { Tabs } from 'antd'
import AllApproveList from './AllApproveList'
import { currentUser } from '@/domain'

@observer
export default class SecurityApproval extends React.Component<any> {

  tabs = [
    {
      label: '操作申请列表',
      key: 'requstList',
      type: 'all',
      perm: currentUser.hasSysMgrPerm
    },
    {
      label: '待审批列表',
      key: 'unapprovedList',
      type: 'unapproved',
      perm: currentUser.hasSecurityApprovalPerm
    },
    {
      label: '已审批列表',
      key: 'approvedList',
      type: 'approved',
      perm: currentUser.hasSecurityApprovalPerm
    },
  ].filter(tab => tab.perm)

  @observable
  activeKey = this.tabs[0].key

  render() {
    return (
      <Wrapper>
        <div className='body'>
          <Tabs activeKey={this.activeKey} onChange={key => this.activeKey = key} >
            {this.tabs.map(tab => {
              return (
                <Tabs.TabPane tab={tab.label} key={tab.key}>
                  { this.activeKey === tab.key && <AllApproveList type={tab.type} /> }
                </Tabs.TabPane>
              )
            })}
          </Tabs>
        </div>
      </Wrapper>
    )
  }
}
