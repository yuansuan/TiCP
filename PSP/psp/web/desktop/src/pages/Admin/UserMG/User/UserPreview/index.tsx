import * as React from 'react'
import { computed, observable, action } from 'mobx'
import { observer } from 'mobx-react'
import { Spin } from 'antd'

import { RoleList, User } from '@/domain/UserMG'
import { Section, RadiusItem } from '../../components'
import BasicInfo from './BasicInfo'
import { UserEditorWrapper, StyledLoading } from './style'

interface IProps {
  user: User
}

@observer
export default class UserPreview extends React.Component<IProps> {
  @observable loading = false
  @action
  updateLoading = loading => (this.loading = loading)

  async componentDidMount() {
    this.updateLoading(true)
    await this.props.user.fetch()
    this.updateLoading(false)
  }

  @computed
  get userRoles() {
    return this.props.user.roles.map(u => RoleList.list.get(u).name) || []
  }

  render() {
    const { loading } = this
    const { user } = this.props

    if (loading) {
      return (
        <StyledLoading>
          <Spin />
        </StyledLoading>
      )
    }

    return (
      <UserEditorWrapper>
        <BasicInfo user={user} />
        {/* <Section title='所在用户组：'>
          <RadiusItem itemList={this.userGroups} />
        </Section> */}
        <Section title='用户角色：'>
          <RadiusItem itemList={this.userRoles} />
        </Section>
      </UserEditorWrapper>
    )
  }
}
