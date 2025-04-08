import { observable, runInAction, action } from 'mobx'
import { Http } from '@/utils'
import currentUser from '@/domain/User'
import {
  CompanyMerchandiseStatus,
  CompanyMerchandise
} from './CompanyMerchandise'

interface ICompanyMerchandise {
  company_id: string
  merchandise_id: string
}

interface IUserList {
  list: Map<string, CompanyMerchandise>
}

export default class UserList implements IUserList {
  @observable list = new Map()
  @observable index = 1
  @observable size = 10
  @observable totals = 0
  @observable status = CompanyMerchandiseStatus.STATE_UNKNOWN
  @observable keyword = ''
  @observable company_id
  @observable company_name

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

  fetch = () => {
    let url = `/company/merchandise/list?index=${this.index}&size=${
      this.size
    }&company_id=${this.company_id || ''}&state=${this.status}&keyword=${
      this.keyword
    }`

    return Http.get(url).then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data?.list?.map(item => {
            return [item.id, new CompanyMerchandise(item)]
          })
        )
        this.totals = res.data?.total || 0
      })
      return res
    })
  }

  addCompanyMerchandise = (body: ICompanyMerchandise) => {
    return Http.post('/company/merchandise', {
      ...body,
      create_uid: currentUser?.user_id || currentUser?.id
    })
  }

  enableMerchandise = mid => {
    return Http.put(
      `/company/merchandise/status/${mid}?action=${CompanyMerchandiseStatus.STATE_ONLINE}`
    )
  }

  disableMerchandise = mid => {
    return Http.put(
      `/company/merchandise/status/${mid}?action=${CompanyMerchandiseStatus.STATE_OFFLINE}`
    )
  };

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
