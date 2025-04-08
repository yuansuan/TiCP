import React from 'react'
import { observer } from 'mobx-react'
import { Role } from '@/domain/UserMG'
import { observable, action } from 'mobx'
import { Spin } from 'antd'

import { RolePreviewWrapper, StyledLoading } from './style'
import { PermPreview } from '../../components'
import BasicInfo from './BasicInfo'

interface IProps {
  role: Role
}

@observer
export default class Preview extends React.Component<IProps> {
  @observable loading = false
  @action
  updateLoading = loading => (this.loading = loading)

  async componentDidMount() {
    this.updateLoading(true)
    await this.props.role.fetch()
    this.updateLoading(false)
  }

  render() {
    const { loading } = this
    const { role } = this.props

    if (loading) {
      return (
        <StyledLoading>
          <Spin />
        </StyledLoading>
      )
    }

    return (
      <RolePreviewWrapper>
        <BasicInfo key='title' title='角色名称：' className='roleName'>
          <span title={role.name}>{role.name}</span>
        </BasicInfo>
        <BasicInfo key='comment' title='角色描述：'>
          {role.comment}
        </BasicInfo>
        <PermPreview perms={role.permList} />
      </RolePreviewWrapper>
    )
  }
}
