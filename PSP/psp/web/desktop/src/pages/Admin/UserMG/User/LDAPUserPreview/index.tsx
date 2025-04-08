import * as React from 'react'
import { computed, observable, action, toJS } from 'mobx'
import { observer } from 'mobx-react'
import { Spin } from 'antd'

import BasicInfo from './BasicInfo'
import { Http } from '@/utils'
import { GroupList, RoleList, User, PermList } from '@/domain/UserMG'
import { Section, RadiusItem } from '../../components'
import { UserEditorWrapper, StyledLoading } from './style'
import uniqBy from 'lodash/uniqBy'

interface IProps {
  user: User
  noFetch?: boolean
}

@observer
export default class LDAPUserPreview extends React.Component<IProps> {
  @observable loading = false
  @observable permList = new PermList({})

  @action
  updateLoading = loading => (this.loading = loading)

  async componentDidMount() {
    if (!this.props.noFetch) {
      this.updateLoading(true)
      await this.props.user.fetch()
      this.updateLoading(false)
    } else {
      const allRes = await Promise.all(
        this.roles.map(async roleId => {
          const res = await Http.get(`/role/${roleId}`)
          return res.data.perm
        })
      )

      // 合并
      let perms = allRes.slice(0, -1).reduce(
        (pre, cur) => {
          pre['local_app'] = [...pre['local_app'], ...cur['local_app']].filter(
            r => r.has
          )

          pre['system'] = [...pre['system'], ...cur['system']].filter(
            r => r.has
          )

          pre['visual_software'] = [
            ...pre['visual_software'],
            ...cur['visual_software']
          ].filter(r => r.has)


          pre['cloud_app'] = [
            ...pre['cloud_app'],
            ...cur['cloud_app']
          ].filter(r => r.has)

          return pre
        },
        { local_app: [], system: [], visual_software: [],cloud_app:[] }
      )

      const last = allRes[allRes.length - 1]
      // 去重
      perms = {
        local_app: uniqBy([...perms['local_app'], ...last['local_app']], 'id'),
        visual_software: uniqBy(
          [...perms['visual_software'], ...last['visual_software']],
          'id'
        ),
        system: uniqBy([...perms['system'], ...last['system']], 'id'),
        cloud_app: uniqBy([...perms['cloud_app'], ...last['cloud_app']], 'id')
      }

      this.permList = new PermList(perms)
    }
  }

  @computed
  get userRoles() {
    return this.props.user.roles.map(u => RoleList.list.get(u).name) || []
  }

  @computed
  get roles() {
    const roles = this.props.user.roles

    return [...roles]
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
        <Section title='角色'>
          <RadiusItem itemList={this.userRoles} />
        </Section>
      </UserEditorWrapper>
    )
  }
}
