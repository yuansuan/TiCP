import { observable, runInAction, action } from 'mobx'
import currentUser from '@/domain/User'
import { Http } from '@/utils'
import Company, { CompanyStatus } from './Company'

interface ICompanyList {
  list: Map<string, Company>
}

interface ICompanyBody {
  name: string
  contact: string
  remark: string
  phone: string
}

export default class CompanyList implements ICompanyList {
  @observable name = ''
  @observable status = CompanyStatus.COMPANY_UNKNOWN
  @observable list = new Map()
  @observable index = 1
  @observable size = 10
  @observable totals = 0
  // TODO
  // @observable orderBy = 'CreateTime'
  // @observable orderAsc = false // 升序降序

  @action
  updateIndex(current: number) {
    this.index = current
  }

  @action
  updateSize(current: number, size: number) {
    this.index = current
    this.size = size
  }

  // @action
  // updateOrder(orderBy, orderAsc) {
  //   this.orderBy = orderBy
  //   this.orderAsc = orderAsc
  // }

  get = id => this.list.get(id)

  fetch = () => {
    let url = `/company/list?index=${this.index}&size=${this.size}&status=${this.status}&name=${this.name}`

    return Http.get(url).then(res => {
      runInAction(() => {
        this.list = new Map(
          res.data?.list?.map(item => {
            return [item.id, new Company(item)]
          })
        )
        this.totals = res.data?.total || 0
      })
      return res
    })
  }

  // 删除
  delete = name => {}

  // 更新
  update = (id, body: ICompanyBody) => {
    return Http.put(`/company/${id}`, {
      ...body,
      modify_uid: currentUser.user_id,
      modify_name: currentUser.name,
    })
  }

  // 创建
  create = (body: ICompanyBody) => {
    return Http.post('/company', {
      ...body,
      create_uid: currentUser.user_id,
      create_name: currentUser.name,
    })
  };

  *[Symbol.iterator]() {
    yield* this.list.values()
  }
}
