import { action, computed, observable } from 'mobx'
import { Timestamp } from '@/utils'
import { PermList, RoleList } from '@/domain/UserMG'
import moment from 'moment'

import { approve_status_map } from './const'

interface ILDAPUserExtAttr {
  // LDAP enabled user additional attributes
  display_name: string
  expired_date: Timestamp
  home_dir: string
  shell: string
}

class LDAPUserExtAttr implements ILDAPUserExtAttr {
  @observable display_name = ''
  @observable expired_date = null
  @observable home_dir = ''
  @observable shell = '/bin/bash'

  constructor(obj) {
    Object.assign(this, {
      ...obj,
      expired_date: new Timestamp(obj.expired_date)
    })
  }
}
interface IUserOrgData {
  id: number
  client_id: string
  name: string
  email: string
  mobile: string
  enabled: boolean
  is_internal: boolean
  created_at: number
  user_id: string
  account_id: string
  real_name: string
  perm: Object
  roles: number[]
  // three members
  approve_status: number
}
interface IUserOrg {
  id: number
  clientId: string
  name: string
  enabled: boolean
  mobile: string
  email: string
  isInternal: boolean
  permList: PermList
  created_at: number
  user_id: string
  accountId: string
  realName: string
  roles: number[]
  // three members
  approve_status: number
}

export default class UserOfOrg implements IUserOrg {
  public readonly id
  @observable public clientId = ''
  @observable public user_id: string
  @observable public accountId = ''
  @observable public name = ''
  @observable public email = ''
  @observable public enabled = true
  @observable public mobile = ''
  @observable public created_at = null
  @observable public isInternal = false
  @observable public permList = new PermList({})
  @observable public realName = ''
  @observable public roles = []
  @observable public roleNames = []
  @observable approve_status = -1

  constructor(props?: Partial<IUserOrgData>) {
    props && this.init(props)
  }

  @action
  public init = (props: Partial<IUserOrgData>) => {
    Object.assign(this, {
      id: Number(props.id),
      user_id: props.user_id,
      accountId: props.account_id,
      clientId: props.client_id,
      name: props.name,
      enabled: props.enabled,
      email: props.email,
      mobile: props.mobile,
      isInternal: props.is_internal,
      created_at: props.created_at
        ? moment(props.created_at).format('YYYY-MM-DD HH:mm:ss')
        : '--',
      permList: new PermList(props.perm || {}),
      roles: props.roles || [],
      roleNames: (props.roles || []).map(r =>
        RoleList.list.get(r) ? RoleList.list.get(r).name : '-'
      ),
      realName: props.real_name,
      approve_status: props.approve_status
    })
  }

  @computed
  get approve_status_str() {
    return approve_status_map.get(this.approve_status) || '--'
  }
}
