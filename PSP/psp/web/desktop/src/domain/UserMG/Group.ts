import { Http } from '@/utils'
import { action, observable, runInAction, computed } from 'mobx'

import { GroupList, PermList } from '@/domain/UserMG'

import { approve_status_map } from './const'

export interface IRequest {
  id: number
  name: string
  perm: Object
  roles: number[]
  users: number[]
  approve_status: number
}

interface IGroup {
  id: number
  name: string
  permList: PermList
  roles: number[]
  users: number[]
  approve_status: number
}

export default class Group implements IGroup {
  @observable public id = -1
  @observable public name = ''
  @observable public roles = []
  @observable public permList = new PermList({})
  @observable public users = []
  @observable public approve_status = -1

  constructor(props?: IRequest) {
    props && this.init(props)
  }

  @computed
  get approve_status_str() {
    return approve_status_map.get(this.approve_status) || '--'
  }

  @action
  public init = (props: IRequest) => {
    Object.assign(this, {
      id: props.id,
      name: props.name,
      users: props.users || [],
      roles: props.roles || [],
      permList: new PermList(props.perm || {}),
      approve_status: props.approve_status,
    })
  }

  @action
  public fetch = () =>
    Http.get(`/group/${this.id}`).then(res => {
      this.init(res.data)
      return res
    })

  @action
  public update = ({ name, users, roles, roleNames, userNames }) => {
    return Http.put(`/group/${this.id}`, {
      name,
      users,
      roles,
      roleNames,
      userNames,
    }).then(res => {
      runInAction(() => {
        GroupList.fetch()
      })
      return res
    })
  }

  public toRequest = (): IRequest => ({
    id: this.id,
    name: this.name,
    roles: this.roles || [],
    users: this.users || [],
    perm: {},
    approve_status: this.approve_status,
  })
}
