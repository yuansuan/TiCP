import * as React from 'react'
import { message, Tabs } from 'antd'
import { UserMGWrapper } from './style'
import User from './User'
import Role from './Role'
import { Organization } from './Organization'
import { currentUser, sysConfig } from '@/domain'
import { history } from '@/utils'

const { TabPane } = Tabs

const tabs =
  currentUser.authType === 'ad'
    ? [
        {
          name: '用户组织结构',
          component: Organization
        },
        {
          name: '角色',
          component: Role
        }
      ]
    : [
        {
          name: '用户',
          component: User
        },

        {
          name: '角色',
          component: Role
        }
      ]

export default class UserMG extends React.Component {
  state = {
    key: tabs[0].name
  }

  render() {
    const { key } = this.state
    return (
      <UserMGWrapper>
        <div className='body'>
          <Tabs activeKey={key} onChange={e => this.setState({ key: e })}>
            {tabs.map(tab => (
              <TabPane tab={tab.name} key={tab.name} />
            ))}
          </Tabs>
          {tabs.map(
            tab => key === tab.name && <tab.component key={tab.name} />
          )}
        </div>
      </UserMGWrapper>
    )
  }
}
