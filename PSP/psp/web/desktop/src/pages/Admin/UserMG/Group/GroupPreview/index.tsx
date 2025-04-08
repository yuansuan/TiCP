import React from 'react'
import { observer } from 'mobx-react'
import { computed, observable, action } from 'mobx'
import { Spin } from 'antd'

import { Group, UserList, RoleList } from '@/domain/UserMG'
import { GroupDetail, StyledLoading } from './style'
import { Section, PermPreview, RadiusItem } from '../../components'

interface IProps {
  group: Group
}

@observer
export default class GroupPreview extends React.Component<IProps> {
  @observable loading = false
  @action
  updateLoading = loading => (this.loading = loading)

  async componentDidMount() {
    this.updateLoading(true)
    await Promise.all([this.props.group.fetch(), UserList.fetch()])
    this.updateLoading(false)
  }

  @computed
  get roleNames() {
    return (
      this.props.group.roles.map(r =>
        RoleList.list.get(r) ? RoleList.list.get(r).name : ''
      ) || []
    ).filter(rl => rl)
  }

  @computed
  get userNames() {
    const allNames =
      this.props.group.users.map(u =>
        UserList.enabledUserMap.get(u)
          ? UserList.enabledUserMap.get(u).name
          : ''
      ) || []
    // for some cases group member is deleted, but user id is still returned by api
    return allNames.filter(al => al)
  }

  render() {
    const { group } = this.props
    const { loading } = this

    if (loading) {
      return (
        <StyledLoading>
          <Spin />
        </StyledLoading>
      )
    }

    return (
      <GroupDetail>
        <div className='groupName' title={group.name}>
          用户组名称：{group.name}
        </div>
        <Section title='用户：'>
          <RadiusItem itemList={this.userNames} />
        </Section>
        <Section title='角色：'>
          <RadiusItem itemList={this.roleNames} />
        </Section>
        <PermPreview perms={group.permList} />
      </GroupDetail>
    )
  }
}
