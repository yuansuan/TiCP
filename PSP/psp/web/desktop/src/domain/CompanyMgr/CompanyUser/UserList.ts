import { observable, runInAction, action } from 'mobx'
import { Http } from '@/utils'
import { User, CompanyUserStatus } from './User'

interface ICompanyUser {
  companyId: string
  userId?: string
  isAdmin: boolean
  status?: string
}
interface IUserList {
  list: Map<string, User>
}

export default class UserList implements IUserList {
  @observable list = new Map()
  @observable index = 1
  @observable size = 10
  @observable totals = 0
  @observable status = CompanyUserStatus.UNKNOWN
  @observable key = ''

  @action
  updateIndex(current: number) {
    this.index = current
  }

  @action
  updateSize(current: number, size: number) {
    this.index = current
    this.size = size
  }

  get = id => this.list.get(id)

  fetch = companyId => {
    let url = `/company/user/list?index=${this.index}&size=${this.size}&companyId=${companyId}&status=${this.status}&key=${this.key}`

    return Http.get(url).then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data?.list?.map(item => {
            return [item.user_id, new User(item)]
          })
        )
        this.totals = res.data?.total || 0
      })
      return res
    })
  }

  addCompanyUser = (body: ICompanyUser) => {
    return Http.post(`/company/user`, { ...body })
  }

  updateCompanyUser = (uid: string, body: ICompanyUser) => {
    return Http.put(`/company/user/${uid}`, { ...body })
  }

  deleteCompanyUser = (uid: string, body: ICompanyUser) => {
    return Http.put(`/company/user/${uid}`, {
      ...body,
      status: CompanyUserStatus.DELETED
    })
  };

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
