import { Http } from '@/utils'
import { action, observable, runInAction, computed } from 'mobx'

import Role from './Role'

export class RoleList {
  @observable public list = new Map()
  @observable public totalRoles = 0
  @observable filter = {
    name_filter: '',
    page: {
      index: 1,
      size: 10
    }
  }
  @computed
  get roleList() {
    return [...this.list.values()]
  }

  public get = id => {
    return this.list.get(id)
  }

  @action
  public fetch = () => {
    return Http.post('/role/query', {
      desc: false,
      order_by: 'name',
      ...this.filter
    }).then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data.roles.map(role => [role.id, new Role(role)])
        )
        this.totalRoles = res.data.total
      })

      const newres = res.data.roles.filter(role => role.type == 1)

      return res
    })
  }

  @action
  public add = (params: { name: string; comment: string; perms: number[] }) =>
    Http.post('/role/add', params).then(res => {
      this.fetch()

      return res
    })

  @action
  public delete = (id, name) =>
    Http.delete(`/role/delete`, {
      params: {
        id,
      }
    }).then(res => {
      this.fetch()
      return res
    })
  @action
  public setDefaultRole = async id =>
    await Http.put('/role/setLdapUserDefRole', {
      id,
    }).then(res => {
      this.fetch()
      return res
    })
}

export default new RoleList()
