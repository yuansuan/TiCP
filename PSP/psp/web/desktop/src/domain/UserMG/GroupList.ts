/**
 * @module GroupList
 * @description total group list
 */
import { action, observable, computed, runInAction } from 'mobx'

import { Http } from '@/utils'
import Group from './Group'

export class GroupList {
  @observable public list = new Map<number, Group>()
  @observable public totalGroups = 0

  @computed
  get groupList() {
    return [...this.list.values()]
  }

  public get = id => this.list.get(id)

  @action
  public fetch = () =>
    Http.get('/group').then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data.groups.map(item => [item.id, new Group(item)])
        )
        this.totalGroups = res.data.total
      })
      return res
    })

  @action
  public add = ({ name, roles, users, roleNames, userNames }) =>
    Http.post('/group', {
      name,
      roles,
      users,
      roleNames,
      userNames,
    }).then(res => {
      runInAction(() => {
        if (res.data?.isAskRequest) {
          this.fetch()
        } else {
          this.list = new Map(
            res.data.groups.map(item => [item.id, new Group(item)])
          )
          this.totalGroups = res.data.total
        }
      })
      return res
    })

  @action
  public delete = (id, name, body) =>
    Http.delete(`/group/${id}?name=${name}`, { data: { ...body } }).then(
      res => {
        runInAction(() => {
          if (res.data?.isAskRequest) {
            this.fetch()
          } else {
            this.list = new Map(
              res.data.groups.map(item => [item.id, new Group(item)])
            )
            this.totalGroups = res.data.total
          }
        })
        return res
      }
    )
}

export default new GroupList()
