import { action, computed, observable, runInAction } from 'mobx'
import { Http } from '@/utils'
import { User } from '@/domain/UserMG'

export class UserList {
  @observable public enabledUserMap = new Map()
  @observable public disabledUserMap = new Map()
  @observable public totalEnabledUsers = 0
  @observable public totalDisabledUsers = 0
  @observable public successUsername = []
  @observable public failureUsername = []
  @observable public filter = {
    query: '',
    page: {
      index: 1,
      size: 10
    }
  }
  @computed
  get enabledUsers() {
    return [...this.enabledUserMap.values()]
  }

  @computed
  get disabledUsers() {
    return [...this.disabledUserMap.values()]
  }

  public get = id => {
    return this.enabledUserMap.get(id)
  }

  @action
  public fetch = () => {
    const url = '/user/query'
    return Http.post(url, {
      desc: true,
      // enabled: true, //展示激活的用户
      order: 'created_at',
      ...this.filter
    }).then(res => {
      runInAction(() => {
        const { user_obj, total } = res.data
        this.enabledUserMap = new Map(
          user_obj.map(user => {
            return [user.id, new User(user)]
          })
        )
        this.totalEnabledUsers = total
      })

      return res
    })
  }

  @action
  public fetchDisabled = () =>
    Http.get('/user/sys').then(res => {
      runInAction(() => {
        this.disabledUserMap = new Map(
          res.data.users.map(user => [user.id, new User(user)])
        )
        this.totalDisabledUsers = res.data.total
      })

      return res
    })

  @action
  public addUser = body => {
    return Http.post('/user/add', body).then(res => {
      // runInAction(() => {
      //   const { user } = res.data
      //   const { users, total } = user
      //   this.enabledUserMap = new Map(
      //     users.map(user => [user.id, new User(user)])
      //   )
      //   this.totalEnabledUsers = total
      // })
      return res
    })
  }

  @action
  public updateLDAP = (id, body) => {
    return Http.put(`/user/ldap/${id}`, body).then(res => {
      runInAction(() => {
        if (!res.data?.isAskRequest) {
          const { users, total } = res.data
          this.enabledUserMap = new Map(
            users.map(user => [user.id, new User(user)])
          )
          this.totalEnabledUsers = total
        }
      })

      return res
    })
  }

  @action
  public add = (id, roles = [], groups = []) =>
    Http.post(`/user/${id}`, { roles, groups }).then(res => {
      runInAction(() => {
        const { users, total } = res.data
        this.enabledUserMap = new Map(
          users.map(user => [user.id, new User(user)])
        )
        this.totalEnabledUsers = total
      })
      return res
    })

  @action
  public addAll = list =>
    Http.post(`/user/batch/adduser`, {
      list: list
    }).then(res => {
      runInAction(() => {
        this.fetch().then(resList => {
          const { users, total } = resList.data
          this.enabledUserMap = new Map(
            users.map(user => [user.id, new User(user)])
          )
          this.totalEnabledUsers = total
        })
      })
      return res
    })

  @action
  public getUser = id =>
    Http.get(`/user/get`, {
      params: { id }
    }).then(res => {
      runInAction(() => {
        this.fetch()
      })
      return res.data
    })

  @action
  public inactive = id =>
    Http.put(`/user/inactive`, { id }).then(res => {
      runInAction(() => {
        this.fetch()
      })
      return res
    })
  @action
  public active = id =>
    Http.put(`/user/active`, { id }).then(res => {
      runInAction(() => {
        this.fetch()
      })
      return res
    })

  @action
  public delete = id =>
    Http.delete(`/user/delete`, { params: { id: id } }).then(res => {
      runInAction(() => {
        this.fetch()
      })
      return res
    })

  @action
  public resetPasswd = (id, name, body) => {
    return Http.put(`/user/ldap/resetpwd/${id}?name=${name}`, body).then(
      res => {
        runInAction(() => {
          this.fetch()
        })
        return res
      }
    )
  }

  @action
  public resetPwd = user_id => {
    return Http.post(`/user/resetPassword`, {
      user_id
    })
  }
}

export default new UserList()
